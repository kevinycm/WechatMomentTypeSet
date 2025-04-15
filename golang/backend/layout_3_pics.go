package backend

import (
	"fmt"
	"math"
)

// calculateThreePicturesLayout determines the *specific* layout template based on the types
// of the three pictures according to the provided rules (3.1 - 3.62).
// It calculates all possible layouts, scales them, checks minimum heights,
// and selects the best based on area. If no layout meets minimums perfectly,
// it falls back to the layout requiring the minimum relaxation of height constraints.
func (e *ContinuousLayoutEngine) calculateThreePicturesLayout(pictures []Picture, layoutAvailableHeight float64) (TemplateLayout, error) {
	if len(pictures) != 3 {
		return TemplateLayout{}, fmt.Errorf("incorrect number of pictures for 3-pic layout: %d", len(pictures))
	}

	spacing := e.imageSpacing
	AW := e.availableWidth

	// Get Aspect Ratios (W/H) and Types
	ARs := make([]float64, 3)
	types := make([]string, 3)
	validARs := true
	for i, pic := range pictures {
		if pic.Height > 0 && pic.Width > 0 {
			ARs[i] = float64(pic.Width) / float64(pic.Height)
			types[i] = getPictureType(ARs[i])
		} else {
			ARs[i] = 1.0 // Default AR
			types[i] = "unknown"
			validARs = false
			fmt.Printf("Warning: Invalid dimensions for picture %d in 3-pic layout.\n", i) // Assuming pic has an Index field
		}
	}

	if !validARs {
		return TemplateLayout{}, fmt.Errorf("invalid dimensions encountered in 3-pic layout")
	}

	// --- Define Layout Calculation Functions ---
	possibleLayouts := map[string]func([]float64, []string, float64, float64) (TemplateLayout, bool, error){
		"1L2R": e.calculateLayout_1L2R,
		"2L1R": e.calculateLayout_2L1R,
		"1T2B": e.calculateLayout_1T2B,
		"2T1B": e.calculateLayout_2T1B,
		"3Row": e.calculateLayout_3Row,
		"3Col": e.calculateLayout_3Col,
	}

	// --- Store results from all layout attempts ---
	validLayouts := make(map[string]TemplateLayout)    // Layouts meeting strict minimums
	layoutAreas := make(map[string]float64)            // Areas for strictly valid layouts
	scaledLayouts := make(map[string]TemplateLayout)   // All calculated & scaled layouts
	layoutViolationFactors := make(map[string]float64) // Max (requiredMinH / actualH) for each layout
	var firstCalcError error

	// --- Calculate and Evaluate All Layouts ---
	for name, calcFunc := range possibleLayouts {
		// Calculate initial layout
		layout, _, err := calcFunc(ARs, types, AW, spacing) // Initial 'meetsMin' from calcFunc is ignored here, we re-check after scaling
		if err != nil {
			fmt.Printf("Debug: Error calculating initial layout %s: %v\n", name, err)
			if firstCalcError == nil {
				firstCalcError = fmt.Errorf("initial layout %s: %w", name, err)
			}
			continue // Skip this layout if calculation failed
		}

		// --- Scale Layout if Needed (Rule 3.2) ---
		scale := 1.0
		if layout.TotalHeight > layoutAvailableHeight {
			if layout.TotalHeight > 1e-6 { // Avoid division by zero/tiny number
				scale = layoutAvailableHeight / layout.TotalHeight
				// Create a scaled copy to avoid modifying the original layout potentially stored elsewhere
				scaledLayout := TemplateLayout{
					Positions:   make([][]float64, len(layout.Positions)),
					Dimensions:  make([][]float64, len(layout.Dimensions)),
					TotalHeight: layout.TotalHeight * scale,
					TotalWidth:  layout.TotalWidth, // Width scaling might be needed depending on layout type, assume AW for now
				}
				for i := range layout.Positions {
					if len(layout.Positions[i]) == 2 { // Basic sanity check
						scaledLayout.Positions[i] = []float64{layout.Positions[i][0] * scale, layout.Positions[i][1] * scale}
					}
					if len(layout.Dimensions[i]) == 2 { // Basic sanity check
						scaledLayout.Dimensions[i] = []float64{layout.Dimensions[i][0] * scale, layout.Dimensions[i][1] * scale}
					}
				}
				layout = scaledLayout // Use the scaled layout for further checks
			} else {
				// Layout height is zero or tiny, cannot scale meaningfully
				fmt.Printf("Debug: Layout %s has zero/tiny height, skipping scaling.\n", name)
				continue // Skip this layout
			}
		} // else: no scaling needed

		scaledLayouts[name] = layout // Store the (potentially) scaled layout

		// --- Check Minimum Heights After Scaling & Calculate Violation Factor ---
		meetsScaledMin := true
		maxViolationFactor := 1.0 // Start assuming it meets minimums
		requiredMinHeights := make([]float64, 3)
		for i, picType := range types {
			switch picType {
			case "wide":
				requiredMinHeights[i] = e.minWideHeight
			case "tall":
				requiredMinHeights[i] = e.minTallHeight
			case "landscape":
				requiredMinHeights[i] = e.minLandscapeHeight
			case "portrait":
				requiredMinHeights[i] = e.minPortraitHeight
			default:
				requiredMinHeights[i] = e.minLandscapeHeight // Fallback
			}

			if i < len(layout.Dimensions) && len(layout.Dimensions[i]) == 2 {
				actualHeight := layout.Dimensions[i][1]
				requiredMinHeight := getRequiredMinHeight(e, picType, 3) // Pass numPics=3
				if actualHeight < requiredMinHeight {
					meetsScaledMin = false
					if actualHeight > 1e-6 { // Avoid division by zero
						violationRatio := requiredMinHeight / actualHeight
						if violationRatio > maxViolationFactor {
							maxViolationFactor = violationRatio
						}
					} else {
						maxViolationFactor = math.Inf(1) // Assign infinite factor if actual height is zero
					}
				}
			} else {
				fmt.Printf("Warning: Invalid dimensions data for layout %s, picture %d\n", name, i)
				meetsScaledMin = false           // Treat invalid data as not meeting minimums
				maxViolationFactor = math.Inf(1) // Assign infinite factor for invalid data
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

	// --- Select Best Layout ---
	if len(validLayouts) > 0 {
		// Strategy 1: At least one layout met the strict minimums, choose the largest area among them
		bestLayoutName := ""
		maxArea := -1.0
		for name, area := range layoutAreas {
			if area > maxArea {
				maxArea = area
				bestLayoutName = name
			}
		}
		fmt.Printf("Debug: Selected best strictly valid 3-pic layout: %s (Area: %.2f)\n", bestLayoutName, maxArea)
		return validLayouts[bestLayoutName], nil
	} else if len(scaledLayouts) > 0 {
		// Strategy 2: No layout met strict minimums, fallback logic
		bestFallbackName := ""
		minViolationFactor := math.Inf(1)

		// Prioritize non-vertical layouts first
		preferredLayouts := []string{"1L2R", "2L1R", "1T2B", "2T1B", "3Row"}
		foundPreferredFallback := false

		for _, name := range preferredLayouts {
			if factor, exists := layoutViolationFactors[name]; exists {
				if factor < minViolationFactor {
					minViolationFactor = factor
					bestFallbackName = name
					foundPreferredFallback = true
				}
			}
		}

		// If no preferred layout was viable (or only 3Col existed), consider 3Col
		if !foundPreferredFallback || minViolationFactor == math.Inf(1) {
			if factor3Col, exists3Col := layoutViolationFactors["3Col"]; exists3Col {
				// Only choose 3Col if it's better than the best preferred OR if no preferred was found
				if !foundPreferredFallback || factor3Col < minViolationFactor {
					// Check if 3Col itself is viable (factor is not infinite)
					if factor3Col != math.Inf(1) {
						minViolationFactor = factor3Col
						bestFallbackName = "3Col"
					}
				}
			}
		}

		if bestFallbackName != "" && minViolationFactor != math.Inf(1) {
			fmt.Printf("Warning: No standard 3-pic layout was valid. Using fallback layout '%s' which required relaxing minimum heights by a factor of %.2f.\n", bestFallbackName, minViolationFactor)
			return scaledLayouts[bestFallbackName], nil
		} else {
			// Fallback failed, likely due to calculation errors or infinite violation factors for all
			errMsg := "No valid 3-picture layout met minimum dimensions, and fallback selection also failed."
			if firstCalcError != nil {
				errMsg = fmt.Sprintf("%s First calculation error: %v", errMsg, firstCalcError)
			}
			return TemplateLayout{}, fmt.Errorf("%s", errMsg)
		}

	} else {
		// Strategy 3: No layouts could even be calculated initially
		errMsg := "No 3-picture layout could be calculated."
		if firstCalcError != nil {
			errMsg = fmt.Sprintf("%s First calculation error: %v", errMsg, firstCalcError)
		}
		return TemplateLayout{}, fmt.Errorf("%s", errMsg)
	}
}
