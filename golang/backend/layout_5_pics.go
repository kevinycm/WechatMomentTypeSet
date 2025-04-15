package backend

import (
	"fmt"
)

// --- Helper Functions (Copied/Adapted for 5-pics) ---

// getRequiredMinHeight is defined globally in continuous_layout.go
/*
func getRequiredMinHeight(e *ContinuousLayoutEngine, picType string) float64 {
	switch picType {
	case "wide":
		return e.minWideHeight
	case "tall":
		return e.minTallHeight
	case "landscape":
		return e.minLandscapeHeight
	case "portrait":
		return e.minPortraitHeight
	default: // square, unknown
		return e.minLandscapeHeight // Use landscape as fallback
	}
}
*/

// checkMinHeights checks if all pictures in a calculated layout meet minimum height requirements.
func checkMinHeights(e *ContinuousLayoutEngine, layout TemplateLayout, types []string) bool {
	if len(layout.Dimensions) != len(types) {
		fmt.Printf("Warning: Dimension/type mismatch (%d/%d) in checkMinHeights\n", len(layout.Dimensions), len(types))
		return false
	}
	for i, picType := range types {
		requiredMinHeight := getRequiredMinHeight(e, picType)
		if i < len(layout.Dimensions) && len(layout.Dimensions[i]) == 2 {
			actualHeight := layout.Dimensions[i][1]
			if actualHeight < requiredMinHeight {
				return false // Found a picture that doesn't meet minimum height
			}
		} else {
			fmt.Printf("Warning: Invalid dimensions data for picture %d in checkMinHeights\n", i)
			return false // Invalid data
		}
	}
	return true // All pictures meet minimum height
}

// calculateRowLayout calculates dimensions for a row of pictures aiming for uniform height.
func calculateRowLayout(ARs []float64, AW, spacing float64) (widths []float64, height float64, err error) {
	numPicsInRow := len(ARs)
	if numPicsInRow < 1 {
		return nil, 0, fmt.Errorf("cannot calculate row layout with zero pictures")
	}

	totalSpacing := float64(numPicsInRow-1) * spacing
	rowAvailableWidth := AW - totalSpacing
	if rowAvailableWidth <= 1e-6 {
		return nil, 0, fmt.Errorf("row available width (%.2f) is too small", rowAvailableWidth)
	}

	totalARSum := 0.0
	for _, ar := range ARs {
		if ar <= 1e-6 {
			return nil, 0, fmt.Errorf("invalid aspect ratio (%.2f) encountered in row calculation", ar)
		}
		totalARSum += ar
	}

	if totalARSum <= 1e-6 {
		return nil, 0, fmt.Errorf("total aspect ratio sum (%.2f) is too small for row calculation", totalARSum)
	}

	height = rowAvailableWidth / totalARSum
	if height <= 1e-6 {
		return nil, 0, fmt.Errorf("calculated row height (%.2f) is too small", height)
	}

	widths = make([]float64, numPicsInRow)
	for i, ar := range ARs {
		widths[i] = height * ar
	}

	return widths, height, nil
}

// --- 5-Picture Layout Calculation Functions ---

// calculateLayout_5_2T3B calculates the 2 Top, 3 Bottom layout.
func calculateLayout_5_2T3B(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 5), Dimensions: make([][]float64, 5)}
	if len(ARs) != 5 || len(types) != 5 {
		return layout, fmt.Errorf("2T3B layout requires 5 ARs and types")
	}

	widths1, height1, err1 := calculateRowLayout(ARs[0:2], AW, spacing)
	if err1 != nil {
		return layout, fmt.Errorf("failed to calculate top row for 2T3B: %w", err1)
	}
	W0, W1 := widths1[0], widths1[1]

	widths2, height2, err2 := calculateRowLayout(ARs[2:5], AW, spacing)
	if err2 != nil {
		return layout, fmt.Errorf("failed to calculate bottom row for 2T3B: %w", err2)
	}
	W2, W3, W4 := widths2[0], widths2[1], widths2[2]

	layout.TotalHeight = height1 + spacing + height2
	layout.TotalWidth = AW

	layout.Positions[0] = []float64{0, 0}
	layout.Dimensions[0] = []float64{W0, height1}
	layout.Positions[1] = []float64{W0 + spacing, 0}
	layout.Dimensions[1] = []float64{W1, height1}

	bottomY := height1 + spacing
	currentX := 0.0
	layout.Positions[2] = []float64{currentX, bottomY}
	layout.Dimensions[2] = []float64{W2, height2}
	currentX += W2 + spacing
	layout.Positions[3] = []float64{currentX, bottomY}
	layout.Dimensions[3] = []float64{W3, height2}
	currentX += W3 + spacing
	layout.Positions[4] = []float64{currentX, bottomY}
	layout.Dimensions[4] = []float64{W4, height2}

	return layout, nil
}

// calculateLayout_5_3T2B calculates the 3 Top, 2 Bottom layout.
func calculateLayout_5_3T2B(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 5), Dimensions: make([][]float64, 5)}
	if len(ARs) != 5 || len(types) != 5 {
		return layout, fmt.Errorf("3T2B layout requires 5 ARs and types")
	}

	widths1, height1, err1 := calculateRowLayout(ARs[0:3], AW, spacing)
	if err1 != nil {
		return layout, fmt.Errorf("failed to calculate top row for 3T2B: %w", err1)
	}
	W0, W1, W2 := widths1[0], widths1[1], widths1[2]

	widths2, height2, err2 := calculateRowLayout(ARs[3:5], AW, spacing)
	if err2 != nil {
		return layout, fmt.Errorf("failed to calculate bottom row for 3T2B: %w", err2)
	}
	W3, W4 := widths2[0], widths2[1]

	layout.TotalHeight = height1 + spacing + height2
	layout.TotalWidth = AW

	currentX := 0.0
	layout.Positions[0] = []float64{currentX, 0}
	layout.Dimensions[0] = []float64{W0, height1}
	currentX += W0 + spacing
	layout.Positions[1] = []float64{currentX, 0}
	layout.Dimensions[1] = []float64{W1, height1}
	currentX += W1 + spacing
	layout.Positions[2] = []float64{currentX, 0}
	layout.Dimensions[2] = []float64{W2, height1}

	bottomY := height1 + spacing
	currentX = 0.0
	layout.Positions[3] = []float64{currentX, bottomY}
	layout.Dimensions[3] = []float64{W3, height2}
	currentX += W3 + spacing
	layout.Positions[4] = []float64{currentX, bottomY}
	layout.Dimensions[4] = []float64{W4, height2}

	return layout, nil
}

// --- NEW: Three-Row Layout Functions ---

// calculateLayout_5_2T2M1B calculates the 2 Top, 2 Middle, 1 Bottom layout.
func calculateLayout_5_2T2M1B(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 5), Dimensions: make([][]float64, 5)}
	if len(ARs) != 5 || len(types) != 5 {
		return layout, fmt.Errorf("2T2M1B layout requires 5 ARs and types")
	}

	// Row 1 (Pics 0, 1)
	widths1, height1, err1 := calculateRowLayout(ARs[0:2], AW, spacing)
	if err1 != nil {
		return layout, fmt.Errorf("failed to calculate top row for 2T2M1B: %w", err1)
	}
	W0, W1 := widths1[0], widths1[1]

	// Row 2 (Pics 2, 3)
	widths2, height2, err2 := calculateRowLayout(ARs[2:4], AW, spacing)
	if err2 != nil {
		return layout, fmt.Errorf("failed to calculate middle row for 2T2M1B: %w", err2)
	}
	W2, W3 := widths2[0], widths2[1]

	// Row 3 (Pic 4)
	W4 := AW
	AR4 := ARs[4]
	height3 := 0.0
	if AR4 > 1e-6 {
		height3 = W4 / AR4
	} else {
		return layout, fmt.Errorf("invalid aspect ratio (%.2f) for bottom picture in 2T2M1B", AR4)
	}
	if height3 <= 1e-6 {
		return layout, fmt.Errorf("calculated zero height (%.2f) for bottom picture in 2T2M1B", height3)
	}

	layout.TotalHeight = height1 + spacing + height2 + spacing + height3
	layout.TotalWidth = AW

	// Positions and Dimensions
	layout.Positions[0] = []float64{0, 0}
	layout.Dimensions[0] = []float64{W0, height1}
	layout.Positions[1] = []float64{W0 + spacing, 0}
	layout.Dimensions[1] = []float64{W1, height1}

	yRow2 := height1 + spacing
	layout.Positions[2] = []float64{0, yRow2}
	layout.Dimensions[2] = []float64{W2, height2}
	layout.Positions[3] = []float64{W2 + spacing, yRow2}
	layout.Dimensions[3] = []float64{W3, height2}

	yRow3 := yRow2 + height2 + spacing
	layout.Positions[4] = []float64{0, yRow3}
	layout.Dimensions[4] = []float64{W4, height3}

	return layout, nil
}

// calculateLayout_5_2T1M2B calculates the 2 Top, 1 Middle, 2 Bottom layout.
func calculateLayout_5_2T1M2B(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 5), Dimensions: make([][]float64, 5)}
	if len(ARs) != 5 || len(types) != 5 {
		return layout, fmt.Errorf("2T1M2B layout requires 5 ARs and types")
	}

	// Row 1 (Pics 0, 1)
	widths1, height1, err1 := calculateRowLayout(ARs[0:2], AW, spacing)
	if err1 != nil {
		return layout, fmt.Errorf("failed to calculate top row for 2T1M2B: %w", err1)
	}
	W0, W1 := widths1[0], widths1[1]

	// Row 2 (Pic 2)
	W2 := AW
	AR2 := ARs[2]
	height2 := 0.0
	if AR2 > 1e-6 {
		height2 = W2 / AR2
	} else {
		return layout, fmt.Errorf("invalid aspect ratio (%.2f) for middle picture in 2T1M2B", AR2)
	}
	if height2 <= 1e-6 {
		return layout, fmt.Errorf("calculated zero height (%.2f) for middle picture in 2T1M2B", height2)
	}

	// Row 3 (Pics 3, 4)
	widths3, height3, err3 := calculateRowLayout(ARs[3:5], AW, spacing)
	if err3 != nil {
		return layout, fmt.Errorf("failed to calculate bottom row for 2T1M2B: %w", err3)
	}
	W3, W4 := widths3[0], widths3[1]

	layout.TotalHeight = height1 + spacing + height2 + spacing + height3
	layout.TotalWidth = AW

	// Positions and Dimensions
	layout.Positions[0] = []float64{0, 0}
	layout.Dimensions[0] = []float64{W0, height1}
	layout.Positions[1] = []float64{W0 + spacing, 0}
	layout.Dimensions[1] = []float64{W1, height1}

	yRow2 := height1 + spacing
	layout.Positions[2] = []float64{0, yRow2}
	layout.Dimensions[2] = []float64{W2, height2}

	yRow3 := yRow2 + height2 + spacing
	layout.Positions[3] = []float64{0, yRow3}
	layout.Dimensions[3] = []float64{W3, height3}
	layout.Positions[4] = []float64{W3 + spacing, yRow3}
	layout.Dimensions[4] = []float64{W4, height3}

	return layout, nil
}

// calculateLayout_5_1T2M2B calculates the 1 Top, 2 Middle, 2 Bottom layout.
func calculateLayout_5_1T2M2B(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 5), Dimensions: make([][]float64, 5)}
	if len(ARs) != 5 || len(types) != 5 {
		return layout, fmt.Errorf("1T2M2B layout requires 5 ARs and types")
	}

	// Row 1 (Pic 0)
	W0 := AW
	AR0 := ARs[0]
	height1 := 0.0
	if AR0 > 1e-6 {
		height1 = W0 / AR0
	} else {
		return layout, fmt.Errorf("invalid aspect ratio (%.2f) for top picture in 1T2M2B", AR0)
	}
	if height1 <= 1e-6 {
		return layout, fmt.Errorf("calculated zero height (%.2f) for top picture in 1T2M2B", height1)
	}

	// Row 2 (Pics 1, 2)
	widths2, height2, err2 := calculateRowLayout(ARs[1:3], AW, spacing)
	if err2 != nil {
		return layout, fmt.Errorf("failed to calculate middle row for 1T2M2B: %w", err2)
	}
	W1, W2 := widths2[0], widths2[1]

	// Row 3 (Pics 3, 4)
	widths3, height3, err3 := calculateRowLayout(ARs[3:5], AW, spacing)
	if err3 != nil {
		return layout, fmt.Errorf("failed to calculate bottom row for 1T2M2B: %w", err3)
	}
	W3, W4 := widths3[0], widths3[1]

	layout.TotalHeight = height1 + spacing + height2 + spacing + height3
	layout.TotalWidth = AW

	// Positions and Dimensions
	layout.Positions[0] = []float64{0, 0}
	layout.Dimensions[0] = []float64{W0, height1}

	yRow2 := height1 + spacing
	layout.Positions[1] = []float64{0, yRow2}
	layout.Dimensions[1] = []float64{W1, height2}
	layout.Positions[2] = []float64{W1 + spacing, yRow2}
	layout.Dimensions[2] = []float64{W2, height2}

	yRow3 := yRow2 + height2 + spacing
	layout.Positions[3] = []float64{0, yRow3}
	layout.Dimensions[3] = []float64{W3, height3}
	layout.Positions[4] = []float64{W3 + spacing, yRow3}
	layout.Dimensions[4] = []float64{W4, height3}

	return layout, nil
}

// calculateFivePicturesLayout determines the best layout for 5 pictures.
func (e *ContinuousLayoutEngine) calculateFivePicturesLayout(pictures []Picture, layoutAvailableHeight float64) (TemplateLayout, error) {
	if len(pictures) != 5 {
		return TemplateLayout{}, fmt.Errorf("incorrect number of pictures for 5-pic layout: %d", len(pictures))
	}

	spacing := e.imageSpacing
	AW := e.availableWidth

	ARs := make([]float64, 5)
	types := make([]string, 5)
	validARs := true
	for i, pic := range pictures {
		if pic.Height > 0 && pic.Width > 0 {
			ARs[i] = float64(pic.Width) / float64(pic.Height)
			types[i] = getPictureType(ARs[i]) // Use global function
		} else {
			ARs[i] = 1.0
			types[i] = "unknown"
			validARs = false
			fmt.Printf("Warning: Invalid dimensions for picture %d in 5-pic layout.\n", i)
		}
	}

	if !validARs {
		return TemplateLayout{}, fmt.Errorf("invalid dimensions encountered in 5-pic layout")
	}

	type calcFuncType func(*ContinuousLayoutEngine, []float64, []string, float64, float64) (TemplateLayout, error)
	possibleLayouts := map[string]calcFuncType{
		"2T3B":   calculateLayout_5_2T3B,
		"3T2B":   calculateLayout_5_3T2B,
		"2T2M1B": calculateLayout_5_2T2M1B,
		"2T1M2B": calculateLayout_5_2T1M2B,
		"1T2M2B": calculateLayout_5_1T2M2B,
	}

	validLayouts := make(map[string]TemplateLayout)
	layoutAreas := make(map[string]float64)
	var firstCalcError error

	for name, calcFunc := range possibleLayouts {
		layout, err := calcFunc(e, ARs, types, AW, spacing)
		if err != nil {
			fmt.Printf("Debug: Error calculating initial 5-pic layout %s: %v\n", name, err)
			if firstCalcError == nil {
				firstCalcError = fmt.Errorf("initial 5-pic layout %s: %w", name, err)
			}
			continue
		}

		if layout.TotalHeight <= layoutAvailableHeight && layout.TotalHeight > 1e-6 {
			if checkMinHeights(e, layout, types) {
				validLayouts[name] = layout
				totalArea := 0.0
				for _, dim := range layout.Dimensions {
					if len(dim) == 2 {
						totalArea += dim[0] * dim[1]
					}
				}
				layoutAreas[name] = totalArea
				fmt.Printf("Debug: 5-Pic Layout %s valid and fits. Area: %.2f\n", name, totalArea)
			} else {
				fmt.Printf("Debug: 5-Pic Layout %s fits but failed minimum height check.\n", name)
			}
		} else {
			fmt.Printf("Debug: 5-Pic Layout %s does not fit available height (%.2f > %.2f) or has zero height.\n", name, layout.TotalHeight, layoutAvailableHeight)
		}
	}

	if len(validLayouts) > 0 {
		bestLayoutName := ""
		maxArea := -1.0
		for name, area := range layoutAreas {
			if area > maxArea {
				maxArea = area
				bestLayoutName = name
			}
		}
		fmt.Printf("Debug: Selected best fitting valid 5-pic layout: %s (Area: %.2f)\n", bestLayoutName, maxArea)
		return validLayouts[bestLayoutName], nil
	} else {
		hasWideOrTall := false
		for _, picType := range types {
			if picType == "wide" || picType == "tall" {
				hasWideOrTall = true
				break
			}
		}

		if hasWideOrTall {
			fmt.Println("Debug: No fitting layout for 5 pics with wide/tall images. Signaling force_new_page.")
			return TemplateLayout{}, fmt.Errorf("force_new_page")
		} else {
			fmt.Println("Debug: No fitting layout for 5 pics (no wide/tall). Signaling split_required.")
			return TemplateLayout{}, fmt.Errorf("split_required")
		}
	}
}
