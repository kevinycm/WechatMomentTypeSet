package backend

import (
	"fmt"
	"math"
)

// --- Calculation Helper Functions (Local to layout_4_pics.go) ---

// calculateLayout_4_2x2 calculates the 2x2 grid layout.
// Aims for uniform height within each row.
func calculateLayout_4_2x2(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 4), Dimensions: make([][]float64, 4)}
	if len(ARs) != 4 || len(types) != 4 {
		return layout, fmt.Errorf("2x2 layout requires 4 ARs and types")
	}
	AR0, AR1, AR2, AR3 := ARs[0], ARs[1], ARs[2], ARs[3]

	// Top row calculation
	topRowW := AW - spacing
	topARSum := AR0 + AR1
	H_top := 0.0
	if topRowW > 1e-6 && topARSum > 1e-6 {
		H_top = topRowW / topARSum
	} else {
		return layout, fmt.Errorf("cannot calculate 2x2 top row height (W=%.2f, ARSum=%.2f)", topRowW, topARSum)
	}
	if H_top <= 1e-6 {
		return layout, fmt.Errorf("2x2 calculated zero height for top row")
	}
	W0 := H_top * AR0
	W1 := H_top * AR1

	// Bottom row calculation
	bottomRowW := AW - spacing
	bottomARSum := AR2 + AR3
	H_bottom := 0.0
	if bottomRowW > 1e-6 && bottomARSum > 1e-6 {
		H_bottom = bottomRowW / bottomARSum
	} else {
		return layout, fmt.Errorf("cannot calculate 2x2 bottom row height (W=%.2f, ARSum=%.2f)", bottomRowW, bottomARSum)
	}
	if H_bottom <= 1e-6 {
		return layout, fmt.Errorf("2x2 calculated zero height for bottom row")
	}
	W2 := H_bottom * AR2
	W3 := H_bottom * AR3

	layout.TotalHeight = H_top + spacing + H_bottom
	layout.TotalWidth = AW

	layout.Positions[0] = []float64{0, 0}
	layout.Dimensions[0] = []float64{W0, H_top}
	layout.Positions[1] = []float64{W0 + spacing, 0}
	layout.Dimensions[1] = []float64{W1, H_top}
	layout.Positions[2] = []float64{0, H_top + spacing}
	layout.Dimensions[2] = []float64{W2, H_bottom}
	layout.Positions[3] = []float64{W2 + spacing, H_top + spacing}
	layout.Dimensions[3] = []float64{W3, H_bottom}

	// Minimum height check removed - will be done in the main function after scaling.

	return layout, nil
}

// calculateLayout_4_4Col calculates the 4-pics-in-a-column layout.

// calculateLayout_4_1T3B calculates the 1 Top, 3 Bottom layout.
func calculateLayout_4_1T3B(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 4), Dimensions: make([][]float64, 4)}
	if len(ARs) != 4 || len(types) != 4 {
		return layout, fmt.Errorf("1T3B layout requires 4 ARs and types")
	}
	AR0, AR1, AR2, AR3 := ARs[0], ARs[1], ARs[2], ARs[3]

	// Top pic (0)
	W0 := AW
	H0 := 0.0
	if AR0 > 1e-6 {
		H0 = W0 / AR0
	} else {
		return layout, fmt.Errorf("1T3B calculated zero height for top pic")
	}
	if H0 <= 1e-6 {
		return layout, fmt.Errorf("1T3B calculated zero height for top pic")
	}

	// Bottom row (1, 2, 3)
	bottomRowW := AW - 2*spacing
	bottomARSum := AR1 + AR2 + AR3
	H_bottom := 0.0
	if bottomRowW > 1e-6 && bottomARSum > 1e-6 {
		H_bottom = bottomRowW / bottomARSum
	} else {
		return layout, fmt.Errorf("cannot calculate 1T3B bottom row height (W=%.2f, ARSum=%.2f)", bottomRowW, bottomARSum)
	}
	if H_bottom <= 1e-6 {
		return layout, fmt.Errorf("1T3B calculated zero height for bottom row")
	}
	W1 := H_bottom * AR1
	W2 := H_bottom * AR2
	W3 := H_bottom * AR3

	layout.TotalHeight = H0 + spacing + H_bottom
	layout.TotalWidth = AW

	layout.Positions[0] = []float64{0, 0}
	layout.Dimensions[0] = []float64{W0, H0}

	currentX := 0.0
	bottomY := H0 + spacing
	layout.Positions[1] = []float64{currentX, bottomY}
	layout.Dimensions[1] = []float64{W1, H_bottom}
	currentX += W1 + spacing
	layout.Positions[2] = []float64{currentX, bottomY}
	layout.Dimensions[2] = []float64{W2, H_bottom}
	currentX += W2 + spacing
	layout.Positions[3] = []float64{currentX, bottomY}
	layout.Dimensions[3] = []float64{W3, H_bottom}

	// Minimum height check removed - will be done in the main function after scaling.

	return layout, nil
}

// calculateLayout_4_3T1B calculates the 3 Top, 1 Bottom layout.
func calculateLayout_4_3T1B(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 4), Dimensions: make([][]float64, 4)}
	if len(ARs) != 4 || len(types) != 4 {
		return layout, fmt.Errorf("3T1B layout requires 4 ARs and types")
	}
	AR0, AR1, AR2, AR3 := ARs[0], ARs[1], ARs[2], ARs[3]

	// Top row (0, 1, 2)
	topRowW := AW - 2*spacing
	topARSum := AR0 + AR1 + AR2
	H_top := 0.0
	if topRowW > 1e-6 && topARSum > 1e-6 {
		H_top = topRowW / topARSum
	} else {
		return layout, fmt.Errorf("cannot calculate 3T1B top row height (W=%.2f, ARSum=%.2f)", topRowW, topARSum)
	}
	if H_top <= 1e-6 {
		return layout, fmt.Errorf("3T1B calculated zero height for top row")
	}
	W0 := H_top * AR0
	W1 := H_top * AR1
	W2 := H_top * AR2

	// Bottom pic (3)
	W3 := AW
	H3 := 0.0
	if AR3 > 1e-6 {
		H3 = W3 / AR3
	} else {
		return layout, fmt.Errorf("3T1B calculated zero height for bottom pic")
	}
	if H3 <= 1e-6 {
		return layout, fmt.Errorf("3T1B calculated zero height for bottom pic")
	}

	layout.TotalHeight = H_top + spacing + H3
	layout.TotalWidth = AW

	currentX := 0.0
	layout.Positions[0] = []float64{currentX, 0}
	layout.Dimensions[0] = []float64{W0, H_top}
	currentX += W0 + spacing
	layout.Positions[1] = []float64{currentX, 0}
	layout.Dimensions[1] = []float64{W1, H_top}
	currentX += W1 + spacing
	layout.Positions[2] = []float64{currentX, 0}
	layout.Dimensions[2] = []float64{W2, H_top}

	layout.Positions[3] = []float64{0, H_top + spacing}
	layout.Dimensions[3] = []float64{W3, H3}

	// Minimum height check removed - will be done in the main function after scaling.

	return layout, nil
}

// calculateLayout_4_1L3R calculates the 1 Left, 3 Right Stacked layout.
// This is complex geometrically. Let's try a simplification:
// Assume left pic (0) and right stack (1, 2, 3) share the same total height.
// Calculate height H1, H2, H3 for right stack assuming they fit in some width WR.
// Total right height HR = H1+H2+H3+2*spacing.
// Calculate required width W0 for left pic to have height HR: W0 = HR * AR0.
// Check if W0 + spacing + WR = AW.
// This often requires iterative solving or algebraic manipulation.
// Simpler approach: Treat right side as a single column, calculate its total AR.
func calculateLayout_4_1L3R(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 4), Dimensions: make([][]float64, 4)}
	if len(ARs) != 4 || len(types) != 4 {
		return layout, fmt.Errorf("1L3R layout requires 4 ARs and types")
	}
	AR0, AR1, AR2, AR3 := ARs[0], ARs[1], ARs[2], ARs[3]

	// Estimate right column effective AR (complex, simplifying)
	// Let's fix the height H and calculate widths.
	// We need to solve for H such that W0 + spacing + WR = AW
	// W0 = H * AR0
	// WR needs H1, H2, H3. Assume WR is fixed for a moment.
	// H1 = WR / AR1, H2 = WR / AR2, H3 = WR / AR3
	// H = H1 + H2 + H3 + 2*spacing
	// Substituting: H = WR/AR1 + WR/AR2 + WR/AR3 + 2*spacing
	// H = WR * (1/AR1 + 1/AR2 + 1/AR3) + 2*spacing
	// W0 = (WR * (1/AR1 + 1/AR2 + 1/AR3) + 2*spacing) * AR0
	// (WR * (1/AR1 + 1/AR2 + 1/AR3) + 2*spacing) * AR0 + spacing + WR = AW
	// This gets complicated quickly. Let's use a simpler geometric approach similar to 1L2R.

	// Assume total height H is determined by the left image (pic 0).
	// Let W0 be width of left, WR be width of right column.
	// W0 + spacing + WR = AW.
	// H = W0 / AR0 (if AR0 != 0)
	// Right stack height must also be H.
	// H = H1 + H2 + H3 + 2*spacing
	// H = WR/AR1 + WR/AR2 + WR/AR3 + 2*spacing = WR * (1/AR1 + 1/AR2 + 1/AR3) + 2*spacing
	// W0/AR0 = WR * (1/AR1 + 1/AR2 + 1/AR3) + 2*spacing
	// (AW - WR - spacing) / AR0 = WR*invSumR + 2*spacing where invSumR = 1/AR1+1/AR2+1/AR3
	// AW/AR0 - WR/AR0 - spacing/AR0 = WR*invSumR + 2*spacing
	// AW/AR0 - spacing/AR0 - 2*spacing = WR*invSumR + WR/AR0 = WR * (invSumR + 1/AR0)
	invSumR := 0.0
	if AR1 > 1e-6 {
		invSumR += 1.0 / AR1
	}
	if AR2 > 1e-6 {
		invSumR += 1.0 / AR2
	}
	if AR3 > 1e-6 {
		invSumR += 1.0 / AR3
	}

	invAR0 := 0.0
	if AR0 > 1e-6 {
		invAR0 = 1.0 / AR0
	}

	numerator := AW*invAR0 - spacing*invAR0 - 2*spacing
	denominator := invSumR + invAR0

	WR := 0.0
	if denominator > 1e-6 {
		WR = numerator / denominator
	} else {
		return layout, fmt.Errorf("1L3R cannot solve for WR (denominator zero)")
	}

	if WR <= 1e-6 || WR >= AW-spacing { // Check if WR is valid
		return layout, fmt.Errorf("1L3R geometry infeasible (WR=%.2f)", WR)
	}
	W0 := AW - spacing - WR
	if W0 <= 1e-6 {
		return layout, fmt.Errorf("1L3R geometry infeasible (W0=%.2f)", W0)
	}

	H := 0.0
	if AR0 > 1e-6 {
		H = W0 / AR0
	} else {
		H = WR*invSumR + 2*spacing
	} // Calc H
	if H <= 1e-6 {
		return layout, fmt.Errorf("1L3R calculated zero total height")
	}

	H1, H2, H3 := 0.0, 0.0, 0.0
	if AR1 > 1e-6 {
		H1 = WR / AR1
	} else {
		return layout, fmt.Errorf("1L3R zero height pic 1")
	}
	if AR2 > 1e-6 {
		H2 = WR / AR2
	} else {
		return layout, fmt.Errorf("1L3R zero height pic 2")
	}
	if AR3 > 1e-6 {
		H3 = WR / AR3
	} else {
		return layout, fmt.Errorf("1L3R zero height pic 3")
	}
	if H1 <= 1e-6 || H2 <= 1e-6 || H3 <= 1e-6 {
		return layout, fmt.Errorf("1L3R calculated zero height in right stack")
	}

	// Verify calculated height matches estimate (within tolerance)
	rightStackH := H1 + H2 + H3 + 2*spacing
	if math.Abs(rightStackH-H) > 1e-3 {
		fmt.Printf("Warning: 1L3R height mismatch H=%.2f, H_stack=%.2f. Adjusting.\n", H, rightStackH)
		// Could potentially adjust spacing or rescale, but let's proceed with calculated H for now.
		H = rightStackH // Favor the stack height calculation? Or average?
		W0 = H * AR0
	}

	layout.TotalHeight = H
	layout.TotalWidth = AW

	layout.Positions[0] = []float64{0, 0}
	layout.Dimensions[0] = []float64{W0, H}

	currentY := 0.0
	rightX := W0 + spacing
	layout.Positions[1] = []float64{rightX, currentY}
	layout.Dimensions[1] = []float64{WR, H1}
	currentY += H1 + spacing
	layout.Positions[2] = []float64{rightX, currentY}
	layout.Dimensions[2] = []float64{WR, H2}
	currentY += H2 + spacing
	layout.Positions[3] = []float64{rightX, currentY}
	layout.Dimensions[3] = []float64{WR, H3}

	// Minimum height check removed - will be done in the main function after scaling.

	return layout, nil
}

// calculateLayout_4_3L1R calculates the 3 Left Stacked, 1 Right layout.
// Mirror image of 1L3R.
func calculateLayout_4_3L1R(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 4), Dimensions: make([][]float64, 4)}
	if len(ARs) != 4 || len(types) != 4 {
		return layout, fmt.Errorf("3L1R layout requires 4 ARs and types")
	}
	AR0, AR1, AR2, AR3 := ARs[0], ARs[1], ARs[2], ARs[3] // 0,1,2 are left stack; 3 is right

	// Similar algebraic approach as 1L3R, solving for WL (width of left stack)
	invSumL := 0.0
	if AR0 > 1e-6 {
		invSumL += 1.0 / AR0
	}
	if AR1 > 1e-6 {
		invSumL += 1.0 / AR1
	}
	if AR2 > 1e-6 {
		invSumL += 1.0 / AR2
	}

	invAR3 := 0.0
	if AR3 > 1e-6 {
		invAR3 = 1.0 / AR3
	}

	numerator := AW*invAR3 - spacing*invAR3 - 2*spacing
	denominator := invSumL + invAR3

	WL := 0.0
	if denominator > 1e-6 {
		WL = numerator / denominator
	} else {
		return layout, fmt.Errorf("3L1R cannot solve for WL (denominator zero)")
	}

	if WL <= 1e-6 || WL >= AW-spacing { // Check if WL is valid
		return layout, fmt.Errorf("3L1R geometry infeasible (WL=%.2f)", WL)
	}
	W3 := AW - spacing - WL
	if W3 <= 1e-6 {
		return layout, fmt.Errorf("3L1R geometry infeasible (W3=%.2f)", W3)
	}

	H := 0.0
	if AR3 > 1e-6 {
		H = W3 / AR3
	} else {
		H = WL*invSumL + 2*spacing
	} // Calc H
	if H <= 1e-6 {
		return layout, fmt.Errorf("3L1R calculated zero total height")
	}

	H0, H1, H2 := 0.0, 0.0, 0.0
	if AR0 > 1e-6 {
		H0 = WL / AR0
	} else {
		return layout, fmt.Errorf("3L1R zero height pic 0")
	}
	if AR1 > 1e-6 {
		H1 = WL / AR1
	} else {
		return layout, fmt.Errorf("3L1R zero height pic 1")
	}
	if AR2 > 1e-6 {
		H2 = WL / AR2
	} else {
		return layout, fmt.Errorf("3L1R zero height pic 2")
	}
	if H0 <= 1e-6 || H1 <= 1e-6 || H2 <= 1e-6 {
		return layout, fmt.Errorf("3L1R calculated zero height in left stack")
	}

	// Verify calculated height matches estimate (within tolerance)
	leftStackH := H0 + H1 + H2 + 2*spacing
	if math.Abs(leftStackH-H) > 1e-3 {
		fmt.Printf("Warning: 3L1R height mismatch H=%.2f, H_stack=%.2f. Adjusting.\n", H, leftStackH)
		H = leftStackH // Favor the stack height calculation
		W3 = H * AR3
	}

	layout.TotalHeight = H
	layout.TotalWidth = AW

	currentY := 0.0
	layout.Positions[0] = []float64{0, currentY}
	layout.Dimensions[0] = []float64{WL, H0}
	currentY += H0 + spacing
	layout.Positions[1] = []float64{0, currentY}
	layout.Dimensions[1] = []float64{WL, H1}
	currentY += H1 + spacing
	layout.Positions[2] = []float64{0, currentY}
	layout.Dimensions[2] = []float64{WL, H2}

	layout.Positions[3] = []float64{WL + spacing, 0}
	layout.Dimensions[3] = []float64{W3, H}

	// Minimum height check removed - will be done in the main function after scaling.

	return layout, nil
}

// --- Main Calculation Function for 4 Pictures ---

// calculateFourPicturesLayout implements the logic described in rules 4.1-4.6
func (e *ContinuousLayoutEngine) calculateFourPicturesLayout(pictures []Picture, layoutAvailableHeight float64) (TemplateLayout, error) {
	if len(pictures) != 4 {
		return TemplateLayout{}, fmt.Errorf("incorrect number of pictures for 4-pic layout: %d", len(pictures))
	}

	spacing := e.imageSpacing
	AW := e.availableWidth

	// Get Aspect Ratios (W/H) and Types
	ARs := make([]float64, 4)
	types := make([]string, 4)
	validARs := true
	for i, pic := range pictures {
		if pic.Height > 0 && pic.Width > 0 {
			ARs[i] = float64(pic.Width) / float64(pic.Height)
			types[i] = getPictureType(ARs[i])
		} else {
			ARs[i] = 1.0 // Default AR
			types[i] = "unknown"
			validARs = false
			fmt.Printf("Warning: Invalid dimensions for picture %d in 4-pic layout.\n", i)
		}
	}

	if !validARs {
		return TemplateLayout{}, fmt.Errorf("invalid dimensions encountered in 4-pic layout")
	}

	// --- Define Layout Calculation Functions Map ---
	type calcFuncType func(*ContinuousLayoutEngine, []float64, []string, float64, float64) (TemplateLayout, error)
	possibleLayouts := map[string]calcFuncType{
		"2x2":  calculateLayout_4_2x2,
		"1T3B": calculateLayout_4_1T3B,
		"3T1B": calculateLayout_4_3T1B,
		"1L3R": calculateLayout_4_1L3R,
		"3L1R": calculateLayout_4_3L1R,
	}

	// --- Store results from all layout attempts ---
	validLayouts := make(map[string]TemplateLayout)    // Layouts meeting strict minimums
	layoutAreas := make(map[string]float64)            // Areas for strictly valid layouts
	scaledLayouts := make(map[string]TemplateLayout)   // All calculated & scaled layouts
	layoutViolationFactors := make(map[string]float64) // Max violation factor for each layout
	var firstCalcError error

	// --- Calculate and Evaluate All Layouts ---
	for name, calcFunc := range possibleLayouts {
		// Calculate initial layout
		layout, err := calcFunc(e, ARs, types, AW, spacing)
		if err != nil {
			fmt.Printf("Debug: Error calculating initial layout %s: %v\n", name, err)
			if firstCalcError == nil {
				firstCalcError = fmt.Errorf("initial layout %s: %w", name, err)
			}
			layoutViolationFactors[name] = math.Inf(1) // Mark as non-viable if calculation fails
			continue                                   // Skip to next layout
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
				fmt.Printf("Debug: Layout %s has zero/tiny height, skipping scaling.\n", name)
				layoutViolationFactors[name] = math.Inf(1) // Mark as non-viable
				continue
			}
		}

		scaledLayouts[name] = layout // Store the final (potentially scaled) layout

		// --- Check Minimum Heights After Scaling & Calculate Violation Factor ---
		meetsScaledMin := true
		maxViolationFactor := 1.0
		for i, picType := range types {
			requiredMinHeight := getRequiredMinHeight(e, picType, len(pictures))
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
				fmt.Printf("Warning: Invalid dimensions data for layout %s, picture %d\n", name, i)
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
			fmt.Printf("Debug: Layout %s valid (Scale: %.2f), Area: %.2f\n", name, scale, totalArea)
		} else {
			fmt.Printf("Debug: Layout %s failed minimum height check (Scale: %.2f, ViolationFactor: %.2f).\n", name, scale, maxViolationFactor)
		}
	}

	// --- Select Best Layout or Signal Split/New Page ---
	if len(validLayouts) > 0 {
		// Strategy 1: Strict selection based on max area
		bestLayoutName := ""
		maxArea := -1.0
		for name, area := range layoutAreas {
			if area > maxArea {
				maxArea = area
				bestLayoutName = name
			}
		}
		fmt.Printf("Debug: Selected best fitting valid 4-pic layout: %s (Area: %.2f)\n", bestLayoutName, maxArea)
		return validLayouts[bestLayoutName], nil
	} else {
		// Strategy 2: No standard layout fits, determine if new page or split is needed
		hasWideOrTall := false
		for _, picType := range types {
			if picType == "wide" || picType == "tall" {
				hasWideOrTall = true
				break
			}
		}

		if hasWideOrTall {
			fmt.Println("Debug: No fitting layout for 4 pics with wide/tall images. Signaling force_new_page.")
			return TemplateLayout{}, fmt.Errorf("force_new_page")
		} else {
			fmt.Println("Debug: No fitting layout for 4 pics (no wide/tall). Signaling split_required.")
			return TemplateLayout{}, fmt.Errorf("split_required")
		}
	}
}
