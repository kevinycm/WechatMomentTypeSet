package waterfall

import (
	"fmt"
	"math"
)

// This might also be moved to global if used by other layout files
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

// --- NEW: Three-Row Layout Functions ---

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
			types[i] = GetPictureType(ARs[i]) // Use global function
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
	scaledLayouts := make(map[string]TemplateLayout)   // Store all calculated & scaled layouts
	layoutViolationFactors := make(map[string]float64) // Max violation factor for each layout
	var firstCalcError error

	for name, calcFunc := range possibleLayouts {
		layout, err := calcFunc(e, ARs, types, AW, spacing)
		if err != nil {
			fmt.Printf("Debug: Error calculating initial 5-pic layout %s: %v\n", name, err)
			if firstCalcError == nil {
				firstCalcError = fmt.Errorf("initial 5-pic layout %s: %w", name, err)
			}
			layoutViolationFactors[name] = math.Inf(1) // Mark as non-viable
			continue
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
					TotalWidth:  layout.TotalWidth, // Assuming layout maintains AW
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
				fmt.Printf("Debug: 5-Pic Layout %s has zero/tiny height, skipping scaling.\n", name)
				layoutViolationFactors[name] = math.Inf(1) // Mark as non-viable
				continue
			}
		}

		scaledLayouts[name] = layout // Store the final (potentially scaled) layout

		// --- Check Minimum Heights After Scaling & Calculate Violation Factor ---
		meetsScaledMin := true
		maxViolationFactor := 1.0
		for i, picType := range types {
			requiredMinHeight := GetRequiredMinHeight(e, picType, len(pictures)) // Use numPics=5
			if i < len(layout.Dimensions) && len(layout.Dimensions[i]) == 2 {
				actualHeight := layout.Dimensions[i][1]
				if actualHeight < requiredMinHeight {
					meetsScaledMin = false
					if actualHeight > 1e-6 {
						violationRatio := requiredMinHeight / actualHeight
						if violationRatio > maxViolationFactor {
							maxViolationFactor = violationRatio
						}
					} else {
						maxViolationFactor = math.Inf(1)
					}
				}
			} else {
				fmt.Printf("Warning: Invalid dimensions data for 5-pic layout %s, picture %d\n", name, i)
				meetsScaledMin = false
				maxViolationFactor = math.Inf(1)
				break
			}
		}
		layoutViolationFactors[name] = maxViolationFactor

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
			fmt.Printf("Debug: 5-Pic Layout %s valid (Scale: %.2f), Area: %.2f\n", name, scale, totalArea)
		} else {
			fmt.Printf("Debug: 5-Pic Layout %s failed minimum height check (Scale: %.2f, ViolationFactor: %.2f).\n", name, scale, maxViolationFactor)
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
