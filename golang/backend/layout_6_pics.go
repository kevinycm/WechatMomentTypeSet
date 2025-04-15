package backend

import (
	"fmt"
	"math"
)

// --- Calculation Helper Functions (Local to layout_6_pics.go) ---

// calculateLayout_6_3T3B calculates the 3 Top, 3 Bottom layout.
func calculateLayout_6_3T3B(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 6), Dimensions: make([][]float64, 6)}
	if len(ARs) != 6 || len(types) != 6 {
		return layout, fmt.Errorf("3T3B layout requires 6 ARs and types")
	}

	// Row 1 (Pics 0, 1, 2)
	widths1, height1, err1 := calculateRowLayout(ARs[0:3], AW, spacing) // Use global helper
	if err1 != nil {
		return layout, fmt.Errorf("failed to calculate top row for 3T3B: %w", err1)
	}
	W0, W1, W2 := widths1[0], widths1[1], widths1[2]

	// Row 2 (Pics 3, 4, 5)
	widths2, height2, err2 := calculateRowLayout(ARs[3:6], AW, spacing) // Use global helper
	if err2 != nil {
		return layout, fmt.Errorf("failed to calculate bottom row for 3T3B: %w", err2)
	}
	W3, W4, W5 := widths2[0], widths2[1], widths2[2]

	layout.TotalHeight = height1 + spacing + height2
	layout.TotalWidth = AW

	// Positions and Dimensions
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
	currentX += W4 + spacing
	layout.Positions[5] = []float64{currentX, bottomY}
	layout.Dimensions[5] = []float64{W5, height2}

	return layout, nil
}

// calculateLayout_6_2T2M2B calculates the 2 Top, 2 Middle, 2 Bottom layout.
func calculateLayout_6_2T2M2B(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 6), Dimensions: make([][]float64, 6)}
	if len(ARs) != 6 || len(types) != 6 {
		return layout, fmt.Errorf("2T2M2B layout requires 6 ARs and types")
	}

	// Row 1 (Pics 0, 1)
	widths1, height1, err1 := calculateRowLayout(ARs[0:2], AW, spacing)
	if err1 != nil {
		return layout, fmt.Errorf("failed to calculate top row for 2T2M2B: %w", err1)
	}
	W0, W1 := widths1[0], widths1[1]

	// Row 2 (Pics 2, 3)
	widths2, height2, err2 := calculateRowLayout(ARs[2:4], AW, spacing)
	if err2 != nil {
		return layout, fmt.Errorf("failed to calculate middle row for 2T2M2B: %w", err2)
	}
	W2, W3 := widths2[0], widths2[1]

	// Row 3 (Pics 4, 5)
	widths3, height3, err3 := calculateRowLayout(ARs[4:6], AW, spacing)
	if err3 != nil {
		return layout, fmt.Errorf("failed to calculate bottom row for 2T2M2B: %w", err3)
	}
	W4, W5 := widths3[0], widths3[1]

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
	layout.Positions[5] = []float64{W4 + spacing, yRow3}
	layout.Dimensions[5] = []float64{W5, height3}

	return layout, nil
}

// calculateLayout_6_3L3R calculates the 3 Left, 3 Right Stacked layout.
func calculateLayout_6_3L3R(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 6), Dimensions: make([][]float64, 6)}
	if len(ARs) != 6 || len(types) != 6 {
		return layout, fmt.Errorf("3L3R layout requires 6 ARs and types")
	}
	AR0, AR1, AR2 := ARs[0], ARs[1], ARs[2] // Left stack
	AR3, AR4, AR5 := ARs[3], ARs[4], ARs[5] // Right stack

	invSumL, invSumR := 0.0, 0.0
	for i := 0; i < 3; i++ {
		if ARs[i] > 1e-6 {
			invSumL += 1.0 / ARs[i]
		} else {
			return layout, fmt.Errorf("invalid AR in left stack for 3L3R")
		}
		if ARs[i+3] > 1e-6 {
			invSumR += 1.0 / ARs[i+3]
		} else {
			return layout, fmt.Errorf("invalid AR in right stack for 3L3R")
		}
	}

	denominator := 1.0
	if invSumR > 1e-6 {
		denominator += invSumL / invSumR
	} else if invSumL > 1e-6 {
		// If right side has infinite AR sum (zero invSumR), right width WR -> 0?
		// This case is tricky, implies right side is infinitely tall. Let's assume error.
		return layout, fmt.Errorf("right stack AR sum invalid for 3L3R")
	}

	WL := 0.0
	if denominator > 1e-6 {
		WL = (AW - spacing) / denominator
	} else {
		return layout, fmt.Errorf("3L3R cannot solve for WL (denominator zero)")
	}

	if WL <= 1e-6 || WL >= AW-spacing {
		return layout, fmt.Errorf("3L3R geometry infeasible (WL=%.2f)", WL)
	}
	WR := AW - spacing - WL
	if WR <= 1e-6 {
		return layout, fmt.Errorf("3L3R geometry infeasible (WR=%.2f)", WR)
	}

	H := 0.0
	if invSumL > 1e-6 {
		H = WL*invSumL + 2*spacing
	} else if invSumR > 1e-6 {
		H = WR*invSumR + 2*spacing
	} else {
		return layout, fmt.Errorf("3L3R cannot determine height (both stacks have zero AR sum)")
	}

	if H <= 1e-6 {
		return layout, fmt.Errorf("3L3R calculated zero total height")
	}

	H0 := WL / AR0
	H1 := WL / AR1
	H2 := WL / AR2
	H3 := WR / AR3
	H4 := WR / AR4
	H5 := WR / AR5

	// Verify heights add up (optional check, might cause issues with float precision)
	leftStackH := H0 + H1 + H2 + 2*spacing
	rightStackH := H3 + H4 + H5 + 2*spacing
	if math.Abs(leftStackH-H) > 1e-3 || math.Abs(rightStackH-H) > 1e-3 {
		fmt.Printf("Warning: 3L3R height mismatch H=%.2f, HL=%.2f, HR=%.2f. Using calculated H.\n", H, leftStackH, rightStackH)
		// Don't adjust H here, stick with the geometrically derived total H
	}

	layout.TotalHeight = H
	layout.TotalWidth = AW

	// Left Column
	currentY := 0.0
	layout.Positions[0] = []float64{0, currentY}
	layout.Dimensions[0] = []float64{WL, H0}
	currentY += H0 + spacing
	layout.Positions[1] = []float64{0, currentY}
	layout.Dimensions[1] = []float64{WL, H1}
	currentY += H1 + spacing
	layout.Positions[2] = []float64{0, currentY}
	layout.Dimensions[2] = []float64{WL, H2}

	// Right Column
	currentY = 0.0
	rightX := WL + spacing
	layout.Positions[3] = []float64{rightX, currentY}
	layout.Dimensions[3] = []float64{WR, H3}
	currentY += H3 + spacing
	layout.Positions[4] = []float64{rightX, currentY}
	layout.Dimensions[4] = []float64{WR, H4}
	currentY += H4 + spacing
	layout.Positions[5] = []float64{rightX, currentY}
	layout.Dimensions[5] = []float64{WR, H5}

	return layout, nil
}

// calculateLayout_6_1T2M3B calculates the 1 Top, 2 Middle, 3 Bottom layout.
func calculateLayout_6_1T2M3B(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 6), Dimensions: make([][]float64, 6)}
	if len(ARs) != 6 || len(types) != 6 {
		return layout, fmt.Errorf("1T2M3B layout requires 6 ARs and types")
	}

	// Row 1 (Pic 0)
	W0 := AW
	AR0 := ARs[0]
	height1 := 0.0
	if AR0 > 1e-6 {
		height1 = W0 / AR0
	} else {
		return layout, fmt.Errorf("invalid AR for top pic in 1T2M3B")
	}
	if height1 <= 1e-6 {
		return layout, fmt.Errorf("zero height for top pic in 1T2M3B")
	}

	// Row 2 (Pics 1, 2)
	widths2, height2, err2 := calculateRowLayout(ARs[1:3], AW, spacing)
	if err2 != nil {
		return layout, fmt.Errorf("failed middle row for 1T2M3B: %w", err2)
	}
	W1, W2 := widths2[0], widths2[1]

	// Row 3 (Pics 3, 4, 5)
	widths3, height3, err3 := calculateRowLayout(ARs[3:6], AW, spacing)
	if err3 != nil {
		return layout, fmt.Errorf("failed bottom row for 1T2M3B: %w", err3)
	}
	W3, W4, W5 := widths3[0], widths3[1], widths3[2]

	layout.TotalHeight = height1 + spacing + height2 + spacing + height3
	layout.TotalWidth = AW

	layout.Positions[0] = []float64{0, 0}
	layout.Dimensions[0] = []float64{W0, height1}
	yRow2 := height1 + spacing
	layout.Positions[1] = []float64{0, yRow2}
	layout.Dimensions[1] = []float64{W1, height2}
	layout.Positions[2] = []float64{W1 + spacing, yRow2}
	layout.Dimensions[2] = []float64{W2, height2}
	yRow3 := yRow2 + height2 + spacing
	currentX := 0.0
	layout.Positions[3] = []float64{currentX, yRow3}
	layout.Dimensions[3] = []float64{W3, height3}
	currentX += W3 + spacing
	layout.Positions[4] = []float64{currentX, yRow3}
	layout.Dimensions[4] = []float64{W4, height3}
	currentX += W4 + spacing
	layout.Positions[5] = []float64{currentX, yRow3}
	layout.Dimensions[5] = []float64{W5, height3}

	return layout, nil
}

// calculateLayout_6_3T2M1B calculates the 3 Top, 2 Middle, 1 Bottom layout.
func calculateLayout_6_3T2M1B(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 6), Dimensions: make([][]float64, 6)}
	if len(ARs) != 6 || len(types) != 6 {
		return layout, fmt.Errorf("3T2M1B layout requires 6 ARs and types")
	}

	// Row 1 (Pics 0, 1, 2)
	widths1, height1, err1 := calculateRowLayout(ARs[0:3], AW, spacing)
	if err1 != nil {
		return layout, fmt.Errorf("failed top row for 3T2M1B: %w", err1)
	}
	W0, W1, W2 := widths1[0], widths1[1], widths1[2]

	// Row 2 (Pics 3, 4)
	widths2, height2, err2 := calculateRowLayout(ARs[3:5], AW, spacing)
	if err2 != nil {
		return layout, fmt.Errorf("failed middle row for 3T2M1B: %w", err2)
	}
	W3, W4 := widths2[0], widths2[1]

	// Row 3 (Pic 5)
	W5 := AW
	AR5 := ARs[5]
	height3 := 0.0
	if AR5 > 1e-6 {
		height3 = W5 / AR5
	} else {
		return layout, fmt.Errorf("invalid AR for bottom pic in 3T2M1B")
	}
	if height3 <= 1e-6 {
		return layout, fmt.Errorf("zero height for bottom pic in 3T2M1B")
	}

	layout.TotalHeight = height1 + spacing + height2 + spacing + height3
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
	yRow2 := height1 + spacing
	layout.Positions[3] = []float64{0, yRow2}
	layout.Dimensions[3] = []float64{W3, height2}
	layout.Positions[4] = []float64{W3 + spacing, yRow2}
	layout.Dimensions[4] = []float64{W4, height2}
	yRow3 := yRow2 + height2 + spacing
	layout.Positions[5] = []float64{0, yRow3}
	layout.Dimensions[5] = []float64{W5, height3}

	return layout, nil
}

// calculateLayout_6_1T3M2B calculates the 1 Top, 3 Middle, 2 Bottom layout.
func calculateLayout_6_1T3M2B(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 6), Dimensions: make([][]float64, 6)}
	if len(ARs) != 6 || len(types) != 6 {
		return layout, fmt.Errorf("1T3M2B layout requires 6 ARs and types")
	}

	// Row 1 (Pic 0)
	W0 := AW
	AR0 := ARs[0]
	height1 := 0.0
	if AR0 > 1e-6 {
		height1 = W0 / AR0
	} else {
		return layout, fmt.Errorf("invalid AR for top pic in 1T3M2B")
	}
	if height1 <= 1e-6 {
		return layout, fmt.Errorf("zero height for top pic in 1T3M2B")
	}

	// Row 2 (Pics 1, 2, 3)
	widths2, height2, err2 := calculateRowLayout(ARs[1:4], AW, spacing)
	if err2 != nil {
		return layout, fmt.Errorf("failed middle row for 1T3M2B: %w", err2)
	}
	W1, W2, W3 := widths2[0], widths2[1], widths2[2]

	// Row 3 (Pics 4, 5)
	widths3, height3, err3 := calculateRowLayout(ARs[4:6], AW, spacing)
	if err3 != nil {
		return layout, fmt.Errorf("failed bottom row for 1T3M2B: %w", err3)
	}
	W4, W5 := widths3[0], widths3[1]

	layout.TotalHeight = height1 + spacing + height2 + spacing + height3
	layout.TotalWidth = AW

	layout.Positions[0] = []float64{0, 0}
	layout.Dimensions[0] = []float64{W0, height1}
	yRow2 := height1 + spacing
	currentX := 0.0
	layout.Positions[1] = []float64{currentX, yRow2}
	layout.Dimensions[1] = []float64{W1, height2}
	currentX += W1 + spacing
	layout.Positions[2] = []float64{currentX, yRow2}
	layout.Dimensions[2] = []float64{W2, height2}
	currentX += W2 + spacing
	layout.Positions[3] = []float64{currentX, yRow2}
	layout.Dimensions[3] = []float64{W3, height2}
	yRow3 := yRow2 + height2 + spacing
	layout.Positions[4] = []float64{0, yRow3}
	layout.Dimensions[4] = []float64{W4, height3}
	layout.Positions[5] = []float64{W4 + spacing, yRow3}
	layout.Dimensions[5] = []float64{W5, height3}

	return layout, nil
}

// calculateLayout_6_2T3M1B calculates the 2 Top, 3 Middle, 1 Bottom layout.
func calculateLayout_6_2T3M1B(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 6), Dimensions: make([][]float64, 6)}
	if len(ARs) != 6 || len(types) != 6 {
		return layout, fmt.Errorf("2T3M1B layout requires 6 ARs and types")
	}

	// Row 1 (Pics 0, 1)
	widths1, height1, err1 := calculateRowLayout(ARs[0:2], AW, spacing)
	if err1 != nil {
		return layout, fmt.Errorf("failed top row for 2T3M1B: %w", err1)
	}
	W0, W1 := widths1[0], widths1[1]

	// Row 2 (Pics 2, 3, 4)
	widths2, height2, err2 := calculateRowLayout(ARs[2:5], AW, spacing)
	if err2 != nil {
		return layout, fmt.Errorf("failed middle row for 2T3M1B: %w", err2)
	}
	W2, W3, W4 := widths2[0], widths2[1], widths2[2]

	// Row 3 (Pic 5)
	W5 := AW
	AR5 := ARs[5]
	height3 := 0.0
	if AR5 > 1e-6 {
		height3 = W5 / AR5
	} else {
		return layout, fmt.Errorf("invalid AR for bottom pic in 2T3M1B")
	}
	if height3 <= 1e-6 {
		return layout, fmt.Errorf("zero height for bottom pic in 2T3M1B")
	}

	layout.TotalHeight = height1 + spacing + height2 + spacing + height3
	layout.TotalWidth = AW

	layout.Positions[0] = []float64{0, 0}
	layout.Dimensions[0] = []float64{W0, height1}
	layout.Positions[1] = []float64{W0 + spacing, 0}
	layout.Dimensions[1] = []float64{W1, height1}
	yRow2 := height1 + spacing
	currentX := 0.0
	layout.Positions[2] = []float64{currentX, yRow2}
	layout.Dimensions[2] = []float64{W2, height2}
	currentX += W2 + spacing
	layout.Positions[3] = []float64{currentX, yRow2}
	layout.Dimensions[3] = []float64{W3, height2}
	currentX += W3 + spacing
	layout.Positions[4] = []float64{currentX, yRow2}
	layout.Dimensions[4] = []float64{W4, height2}
	yRow3 := yRow2 + height2 + spacing
	layout.Positions[5] = []float64{0, yRow3}
	layout.Dimensions[5] = []float64{W5, height3}

	return layout, nil
}

// --- Main Calculation Function for 6 Pictures ---
func (e *ContinuousLayoutEngine) calculateSixPicturesLayout(pictures []Picture, layoutAvailableHeight float64) (TemplateLayout, error) {
	if len(pictures) != 6 {
		return TemplateLayout{}, fmt.Errorf("incorrect number of pictures for 6-pic layout: %d", len(pictures))
	}

	spacing := e.imageSpacing
	AW := e.availableWidth

	// Get Aspect Ratios (W/H) and Types
	ARs := make([]float64, 6)
	types := make([]string, 6)
	validARs := true
	for i, pic := range pictures {
		if pic.Height > 0 && pic.Width > 0 {
			ARs[i] = float64(pic.Width) / float64(pic.Height)
			types[i] = getPictureType(ARs[i]) // Use global helper
		} else {
			ARs[i] = 1.0 // Default AR
			types[i] = "unknown"
			validARs = false
			fmt.Printf("Warning: Invalid dimensions for picture %d in 6-pic layout.\n", i)
		}
	}

	if !validARs {
		return TemplateLayout{}, fmt.Errorf("invalid dimensions encountered in 6-pic layout")
	}

	// --- Define Layout Calculation Functions Map ---
	type calcFuncType func(*ContinuousLayoutEngine, []float64, []string, float64, float64) (TemplateLayout, error)
	possibleLayouts := map[string]calcFuncType{
		"3T3B":   calculateLayout_6_3T3B,
		"2T2M2B": calculateLayout_6_2T2M2B,
		"3L3R":   calculateLayout_6_3L3R,
		"1T2M3B": calculateLayout_6_1T2M3B,
		"3T2M1B": calculateLayout_6_3T2M1B,
		"1T3M2B": calculateLayout_6_1T3M2B,
		"2T3M1B": calculateLayout_6_2T3M1B,
	}

	// --- Store results from all layout attempts ---
	validLayouts := make(map[string]TemplateLayout) // Layouts meeting strict minimums
	layoutAreas := make(map[string]float64)         // Areas for strictly valid layouts
	// scaledLayouts := make(map[string]TemplateLayout)   // Optional: Store all scaled layouts for fallback
	// layoutViolationFactors := make(map[string]float64) // Optional: Store violation factors for fallback
	var firstCalcError error

	// --- Calculate and Evaluate All Layouts ---
	for name, calcFunc := range possibleLayouts {
		layout, err := calcFunc(e, ARs, types, AW, spacing)
		if err != nil {
			fmt.Printf("Debug: Error calculating initial 6-pic layout %s: %v\n", name, err)
			if firstCalcError == nil {
				firstCalcError = fmt.Errorf("initial 6-pic layout %s: %w", name, err)
			}
			continue // Skip this layout
		}

		// --- Scale Layout if Needed ---
		scale := 1.0
		if layout.TotalHeight > layoutAvailableHeight {
			if layout.TotalHeight > 1e-6 {
				scale = layoutAvailableHeight / layout.TotalHeight
				scaledLayout := TemplateLayout{
					Positions:   make([][]float64, len(layout.Positions)),
					Dimensions:  make([][]float64, len(layout.Dimensions)),
					TotalHeight: layout.TotalHeight * scale,
					TotalWidth:  layout.TotalWidth, // Assume width stays AW
				}
				for i := range layout.Positions {
					if len(layout.Positions[i]) == 2 {
						scaledLayout.Positions[i] = []float64{layout.Positions[i][0] * scale, layout.Positions[i][1] * scale}
					}
					if len(layout.Dimensions[i]) == 2 {
						scaledLayout.Dimensions[i] = []float64{layout.Dimensions[i][0] * scale, layout.Dimensions[i][1] * scale}
					}
				}
				layout = scaledLayout // Use the scaled layout
			} else {
				fmt.Printf("Debug: Layout %s has zero/tiny height, skipping scaling.\n", name)
				continue
			}
		}
		// scaledLayouts[name] = layout // Optional Store

		// --- Check Minimum Heights After Scaling ---
		meetsScaledMin := true
		// maxViolationFactor := 1.0 // Optional Store
		for i, picType := range types {
			requiredMinHeight := getRequiredMinHeight(e, picType) // Use global helper
			if i < len(layout.Dimensions) && len(layout.Dimensions[i]) == 2 {
				actualHeight := layout.Dimensions[i][1]
				if actualHeight < requiredMinHeight {
					meetsScaledMin = false
					// Optional: Calculate violation factor
					// if actualHeight > 1e-6 { ... } else { maxViolationFactor = math.Inf(1) }
					break // Exit check early if one fails
				}
			} else {
				fmt.Printf("Warning: Invalid dimensions data for layout %s, picture %d\n", name, i)
				meetsScaledMin = false
				// maxViolationFactor = math.Inf(1) // Optional Store
				break
			}
		}
		// layoutViolationFactors[name] = maxViolationFactor // Optional Store

		// --- Store Strictly Valid Layout and Calculate Area ---
		if meetsScaledMin {
			validLayouts[name] = layout
			totalArea := 0.0
			for _, dim := range layout.Dimensions {
				if len(dim) == 2 {
					totalArea += dim[0] * dim[1]
				}
			}
			layoutAreas[name] = totalArea
			fmt.Printf("Debug: 6-Pic Layout %s valid (Scale: %.2f), Area: %.2f\n", name, scale, totalArea)
		} else {
			fmt.Printf("Debug: 6-Pic Layout %s failed minimum height check (Scale: %.2f).\n", name, scale)
		}
	}

	// --- Select Best Layout or Signal Split/New Page ---
	if len(validLayouts) > 0 {
		bestLayoutName := ""
		maxArea := -1.0
		for name, area := range layoutAreas {
			if area > maxArea {
				maxArea = area
				bestLayoutName = name
			}
		}
		fmt.Printf("Debug: Selected best fitting valid 6-pic layout: %s (Area: %.2f)\n", bestLayoutName, maxArea)
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
			fmt.Println("Debug: No fitting layout for 6 pics with wide/tall images. Signaling force_new_page.")
			return TemplateLayout{}, fmt.Errorf("force_new_page")
		} else {
			fmt.Println("Debug: No fitting layout for 6 pics (no wide/tall). Signaling split_required.")
			return TemplateLayout{}, fmt.Errorf("split_required")
		}
		// Optional: Implement fallback logic using scaledLayouts and layoutViolationFactors if needed
		// Currently, if no layout is strictly valid, we signal error.
	}
}

// Helper function to get minimum height based on type (should be global)
/*
func getRequiredMinHeight(e *ContinuousLayoutEngine, picType string) float64 {
	switch picType {
	case "wide": return e.minWideHeight
	case "tall": return e.minTallHeight
	case "landscape": return e.minLandscapeHeight
	case "portrait": return e.minPortraitHeight
	default: return e.minLandscapeHeight // Fallback
	}
}
*/

// Helper function to calculate row layout (should be global)
/*
func calculateRowLayout(ARs []float64, AW, spacing float64) (widths []float64, height float64, err error) {
	numPicsInRow := len(ARs)
	if numPicsInRow < 1 { return nil, 0, fmt.Errorf("zero pictures in row") }
	totalSpacing := float64(numPicsInRow-1) * spacing
	rowAvailableWidth := AW - totalSpacing
	if rowAvailableWidth <= 1e-6 { return nil, 0, fmt.Errorf("row width too small") }
	totalARSum := 0.0
	for _, ar := range ARs { if ar <= 1e-6 { return nil, 0, fmt.Errorf("invalid AR") }; totalARSum += ar }
	if totalARSum <= 1e-6 { return nil, 0, fmt.Errorf("AR sum too small") }
	height = rowAvailableWidth / totalARSum
	if height <= 1e-6 { return nil, 0, fmt.Errorf("calculated height too small") }
	widths = make([]float64, numPicsInRow)
	for i, ar := range ARs { widths[i] = height * ar }
	return widths, height, nil
}
*/
