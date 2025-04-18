package waterfall

import (
	"fmt"
	"math"
)

// --- 9-Picture Layout Calculation Functions (Max 3 per row) ---

// --- Main Calculation Function for 9 Pictures ---
func (e *ContinuousLayoutEngine) calculateNinePicturesLayout(pictures []Picture, layoutAvailableHeight float64) (TemplateLayout, error) {
	numPics := 9
	if len(pictures) != numPics {
		return TemplateLayout{}, fmt.Errorf("incorrect number of pictures for %d-pic layout: %d", numPics, len(pictures))
	}

	spacing := e.imageSpacing
	AW := e.availableWidth

	// Get Aspect Ratios (W/H) and Types
	ARs := make([]float64, numPics)
	types := make([]string, numPics)
	validARs := true
	for i, pic := range pictures {
		if pic.Height > 0 && pic.Width > 0 {
			ARs[i] = float64(pic.Width) / float64(pic.Height)
			types[i] = GetPictureType(ARs[i])
		} else {
			ARs[i] = 1.0 // Default AR
			types[i] = "unknown"
			validARs = false
			fmt.Printf("Warning: Invalid dimensions for picture %d in %d-pic layout.\n", i, numPics)
		}
	}
	if !validARs {
		return TemplateLayout{}, fmt.Errorf("invalid dimensions encountered in %d-pic layout", numPics)
	}

	// --- Define Layout Calculation Functions Map ---
	type calcFuncType func(*ContinuousLayoutEngine, []float64, []string, float64, float64) (TemplateLayout, error)
	possibleLayouts := map[string]calcFuncType{
		"3T3M2M1B": calculateLayout_9_3T3M2M1B,
		"3T2M2M2B": calculateLayout_9_3T2M2M2B,
		"2T3M2M2B": calculateLayout_9_2T3M2M2B,
		"2T2M3M2B": calculateLayout_9_2T2M3M2B,
		"2T2M2M3B": calculateLayout_9_2T2M2M3B,
		"3T3M3B":   calculateLayout_9_3T3M3B, // Include 3x3 as well
	}

	// --- Store results from all layout attempts ---
	validLayouts := make(map[string]TemplateLayout)
	layoutAreas := make(map[string]float64)
	scaledLayouts := make(map[string]TemplateLayout)
	layoutViolationFactors := make(map[string]float64)
	var firstCalcError error

	// --- Calculate and Evaluate All Layouts ---
	for name, calcFunc := range possibleLayouts {
		layout, err := calcFunc(e, ARs, types, AW, spacing)
		if err != nil {
			fmt.Printf("Debug: Error calculating initial %d-pic layout %s: %v\n", numPics, name, err)
			if firstCalcError == nil {
				firstCalcError = fmt.Errorf("initial %d-pic layout %s: %w", numPics, name, err)
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
					TotalWidth:  layout.TotalWidth,
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
				fmt.Printf("Debug: %d-Pic Layout %s has zero/tiny height, skipping scaling.\n", numPics, name)
				layoutViolationFactors[name] = math.Inf(1) // Mark as non-viable
				continue
			}
		}
		scaledLayouts[name] = layout // Store the final (potentially scaled) layout

		// --- Check Minimum Heights After Scaling & Calculate Violation Factor ---
		meetsScaledMin := true
		maxViolationFactor := 1.0
		for i, picType := range types {
			requiredMinHeight := GetRequiredMinHeight(e, picType, numPics) // Use numPics
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
				fmt.Printf("Warning: Invalid dimensions data for %d-pic layout %s, picture %d\n", numPics, name, i)
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
			fmt.Printf("Debug: %d-Pic Layout %s valid (Scale: %.2f), Area: %.2f\n", numPics, name, scale, totalArea)
		} else {
			fmt.Printf("Debug: %d-Pic Layout %s failed minimum height check (Scale: %.2f, ViolationFactor: %.2f).\n", numPics, name, scale, maxViolationFactor)
		}
	}

	// --- Select Best Layout or Signal Split/New Page ---
	// +++ Add Detailed Logging Before Selection +++
	fmt.Printf("Debug (Calc 9-Pic): Evaluation complete. NumValidLayouts: %d, FirstCalcError: %v\n", len(validLayouts), firstCalcError)
	for name, violFactor := range layoutViolationFactors {
		fmt.Printf("Debug (Calc 9-Pic): Layout %s -> ViolationFactor: %.2f\n", name, violFactor)
	}
	// +++ End Logging +++
	if len(validLayouts) > 0 {
		bestLayoutName := ""
		maxArea := -1.0
		for name, area := range layoutAreas {
			if area > maxArea {
				maxArea = area
				bestLayoutName = name
			}
		}
		fmt.Printf("Debug (Calc 9-Pic): Selected best fitting valid %d-pic layout: %s (Area: %.2f)\n", numPics, bestLayoutName, maxArea)
		return validLayouts[bestLayoutName], nil
	} else {
		// --- Original Logic: No strictly valid layout found, signal error ---
		fmt.Println("Debug (Calc 9-Pic): No strictly valid layouts found. Signaling error based on picture types.") // Modified log

		// Fallback logic (Original - determine error type)
		hasWideOrTall := false
		for _, picType := range types {
			if picType == "wide" || picType == "tall" {
				hasWideOrTall = true
				break
			}
		}

		if hasWideOrTall {
			fmt.Printf("Debug (Calc 9-Pic): No valid layouts found AND has wide/tall. Signaling force_new_page.\n")
			return TemplateLayout{}, fmt.Errorf("force_new_page")
		} else {
			fmt.Printf("Debug (Calc 9-Pic): No valid layouts found (no wide/tall). Signaling split_required.\n")
			return TemplateLayout{}, fmt.Errorf("split_required")
		}
	}
}
