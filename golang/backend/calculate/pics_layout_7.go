package calculate

import (
	"fmt"
	"math"
)

// --- Main Calculation Function for 7 Pictures ---

func (e *ContinuousLayoutEngine) calculateSevenPicturesLayout(pictures []Picture, layoutAvailableHeight float64) (TemplateLayout, error) {
	if len(pictures) != 7 {
		return TemplateLayout{}, fmt.Errorf("incorrect number of pictures for 7-pic layout: %d", len(pictures))
	}

	spacing := e.imageSpacing
	AW := e.availableWidth

	// Get Aspect Ratios (W/H) and Types
	ARs := make([]float64, 7)
	types := make([]string, 7)
	validARs := true
	for i, pic := range pictures {
		if pic.Height > 0 && pic.Width > 0 {
			ARs[i] = float64(pic.Width) / float64(pic.Height)
			types[i] = GetPictureType(ARs[i]) // Use global helper
		} else {
			ARs[i] = 1.0
			types[i] = "unknown"
			validARs = false
			fmt.Printf("Warning: Invalid dimensions for picture %d in 7-pic layout.\n", i)
		}
	}
	if !validARs {
		return TemplateLayout{}, fmt.Errorf("invalid dimensions in 7-pic layout")
	}

	// --- Define Layout Calculation Functions Map (Updated) ---
	type calcFuncType func(*ContinuousLayoutEngine, []float64, []string, float64, float64) (TemplateLayout, error)
	possibleLayouts := map[string]calcFuncType{
		"3T2M2B":   calculateLayout_7_3T2M2B,
		"2T3M2B":   calculateLayout_7_2T3M2B,
		"2T2M3B":   calculateLayout_7_2T2M3B,
		"3T3M1B":   calculateLayout_7_3T3M1B,
		"3L4R":     calculateLayout_7_3L4R,
		"4L3R":     calculateLayout_7_4L3R,
		"2T2M2M1B": calculateLayout_7_2T2M2M1B,
		"1T2M2M2B": calculateLayout_7_1T2M2M2B,
		"3T2M1M1B": calculateLayout_7_3T2M1M1B,
		"1T2M3M1B": calculateLayout_7_1T2M3M1B,
	}

	// --- Store results ---
	validLayouts := make(map[string]TemplateLayout)
	layoutAreas := make(map[string]float64)
	scaledLayouts := make(map[string]TemplateLayout)   // Store all calculated & scaled layouts
	layoutViolationFactors := make(map[string]float64) // Store violation factors for fallback
	var firstCalcError error

	// --- Calculate and Evaluate All Layouts ---
	for name, calcFunc := range possibleLayouts {
		layout, err := calcFunc(e, ARs, types, AW, spacing)
		if err != nil {
			fmt.Printf("Debug: Error calculating initial 7-pic layout %s: %v\n", name, err)
			if firstCalcError == nil {
				firstCalcError = fmt.Errorf("initial 7-pic layout %s: %w", name, err)
			}
			continue
		}

		// Scale Layout if Needed
		scale := 1.0
		if layout.TotalHeight > layoutAvailableHeight {
			if layout.TotalHeight > 1e-6 {
				scale = layoutAvailableHeight / layout.TotalHeight
				scaledLayout := TemplateLayout{TotalHeight: layout.TotalHeight * scale, TotalWidth: layout.TotalWidth}
				scaledLayout.Positions = make([][]float64, len(layout.Positions))
				scaledLayout.Dimensions = make([][]float64, len(layout.Dimensions))
				for i := range layout.Positions {
					// Error check for nil slices if necessary, assuming calcFunc initializes correctly
					if len(layout.Positions[i]) == 2 {
						scaledLayout.Positions[i] = []float64{layout.Positions[i][0] * scale, layout.Positions[i][1] * scale}
					}
					if len(layout.Dimensions[i]) == 2 {
						scaledLayout.Dimensions[i] = []float64{layout.Dimensions[i][0] * scale, layout.Dimensions[i][1] * scale}
					}
				}
				layout = scaledLayout
			} else {
				fmt.Printf("Debug: Layout %s has zero/tiny height, skipping scaling.\n", name)
				continue
			}
		}
		scaledLayouts[name] = layout // Store the final (potentially scaled) layout

		// --- Check Minimum Heights After Scaling & Calculate Violation Factor ---
		meetsScaledMin := true
		maxViolationFactor := 1.0 // Start assuming it meets minimums
		if !CheckMinHeights(e, layout, types, 7) {
			meetsScaledMin = false
			// Calculate violation factor only if it failed
			for i, picType := range types {
				requiredMinHeight := GetRequiredMinHeight(e, picType, len(pictures))
				if i < len(layout.Dimensions) && len(layout.Dimensions[i]) == 2 {
					actualHeight := layout.Dimensions[i][1]
					if actualHeight < requiredMinHeight {
						if actualHeight > 1e-6 {
							violationRatio := requiredMinHeight / actualHeight
							if violationRatio > maxViolationFactor {
								maxViolationFactor = violationRatio
							}
						} else {
							maxViolationFactor = math.Inf(1) // Assign infinite factor if actual height is zero
						}
					}
				} else {
					maxViolationFactor = math.Inf(1) // Treat invalid data as infinite violation
					break
				}
			}
			fmt.Printf("Debug: 7-Pic Layout %s failed minimum height check (Scale: %.2f, ViolationFactor: %.2f).\n", name, scale, maxViolationFactor)
		} else {
			fmt.Printf("Debug: 7-Pic Layout %s passed minimum height check (Scale: %.2f).\n", name, scale)
		}
		layoutViolationFactors[name] = maxViolationFactor // Store violation factor regardless

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
			fmt.Printf("Debug: 7-Pic Layout %s stored as valid. Area: %.2f\n", name, totalArea)
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
		fmt.Printf("Debug: Selected best fitting valid 7-pic layout: %s (Area: %.2f)\n", bestLayoutName, maxArea)
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
			fmt.Println("Debug: No fitting layout for 7 pics with wide/tall images. Signaling force_new_page.")
			return TemplateLayout{}, fmt.Errorf("force_new_page")
		} else {
			fmt.Println("Debug: No fitting layout for 7 pics (no wide/tall). Signaling split_required.")
			return TemplateLayout{}, fmt.Errorf("split_required")
		}
	}
}
