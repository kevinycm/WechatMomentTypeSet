package calculate

import "fmt"

// --- Helper function for 2 Top, 1 Bottom Full Width Template ---
func (e *ContinuousLayoutEngine) calculateLayout_2T1B(ARs []float64, types []string, AW, spacing float64) (TemplateLayout, bool, error) {
	layout := TemplateLayout{
		Positions:  make([][]float64, 3),
		Dimensions: make([][]float64, 3),
	}
	AR0, AR1, AR2 := ARs[0], ARs[1], ARs[2]

	// Calculate top row height (H_top) based on fitting pics 0 & 1 in available width
	topRowAvailableWidth := AW - spacing
	topTotalARSum := AR0 + AR1
	H_top := 0.0
	if topRowAvailableWidth > 1e-6 && topTotalARSum > 1e-6 {
		H_top = topRowAvailableWidth / topTotalARSum
	} else {
		return layout, false, fmt.Errorf("cannot calculate 2T1B top row height")
	}
	if H_top <= 1e-6 {
		return layout, false, fmt.Errorf("2T1B calculated zero height for top row")
	}

	W0 := H_top * AR0
	W1 := H_top * AR1

	// Pic 2 takes full width at bottom
	W2 := AW
	H2 := 0.0
	if AR2 > 1e-6 {
		H2 = W2 / AR2
	}
	if H2 <= 1e-6 {
		return layout, false, fmt.Errorf("2T1B calculated zero height for bottom picture")
	}

	layout.TotalHeight = H_top + spacing + H2
	layout.TotalWidth = AW

	layout.Positions[0] = []float64{0, 0} // Pic 0 Top Left
	layout.Dimensions[0] = []float64{W0, H_top}
	layout.Positions[1] = []float64{W0 + spacing, 0} // Pic 1 Top Right
	layout.Dimensions[1] = []float64{W1, H_top}
	layout.Positions[2] = []float64{0, H_top + spacing} // Pic 2 Bottom Full
	layout.Dimensions[2] = []float64{W2, H2}

	// Type-specific minimum height check
	meetsMin := true
	// heights := []float64{H_top, H_top, H2} // Heights for pics 0, 1, 2 - Not strictly needed
	requiredMinHeights := make([]float64, 3)
	for i, picType := range types {
		// Use the GetRequiredMinHeight helper function
		requiredMinHeight := GetRequiredMinHeight(e, picType, 3)
		requiredMinHeights[i] = requiredMinHeight // Store it (might be redundant)

		/* // Old switch logic removed
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
		*/

		// Use correct height (H_top for 0, 1; H2 for 2)
		checkHeight := H_top
		if i == 2 {
			checkHeight = H2
		}
		if checkHeight < requiredMinHeights[i] {
			meetsMin = false
			break
		}
	}
	// No explicit min width check based on rules

	return layout, meetsMin, nil
}
