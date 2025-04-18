package waterfall

import "fmt"

// --- Helper function for 1 Top, 2 Bottom Stacked Template ---
func (e *ContinuousLayoutEngine) calculateLayout_1T2B(ARs []float64, types []string, AW, spacing float64) (TemplateLayout, bool, error) {
	layout := TemplateLayout{
		Positions:  make([][]float64, 3),
		Dimensions: make([][]float64, 3),
	}
	AR0, AR1, AR2 := ARs[0], ARs[1], ARs[2]

	// Pic 0 takes full width
	W0 := AW
	if W0 <= 1e-6 {
		return layout, false, fmt.Errorf("1T2B available width is zero")
	}
	H0 := 0.0
	if AR0 > 1e-6 {
		H0 = W0 / AR0
	}
	if H0 <= 1e-6 {
		return layout, false, fmt.Errorf("1T2B calculated zero height for top picture")
	}

	// Calculate bottom row height based on fitting pics 1 & 2 in available width
	bottomRowAvailableWidth := AW - spacing
	bottomTotalARSum := AR1 + AR2
	H_bottom := 0.0
	if bottomRowAvailableWidth > 1e-6 && bottomTotalARSum > 1e-6 {
		H_bottom = bottomRowAvailableWidth / bottomTotalARSum
	} else {
		return layout, false, fmt.Errorf("cannot calculate 1T2B bottom row height")
	}
	if H_bottom <= 1e-6 {
		return layout, false, fmt.Errorf("1T2B calculated zero height for bottom row")
	}

	W1 := H_bottom * AR1
	W2 := H_bottom * AR2

	layout.TotalHeight = H0 + spacing + H_bottom
	layout.TotalWidth = AW

	layout.Positions[0] = []float64{0, 0} // Pic 0 Top
	layout.Dimensions[0] = []float64{W0, H0}

	layout.Positions[1] = []float64{0, H0 + spacing} // Pic 1 Bottom Left
	layout.Dimensions[1] = []float64{W1, H_bottom}

	layout.Positions[2] = []float64{W1 + spacing, H0 + spacing} // Pic 2 Bottom Right
	layout.Dimensions[2] = []float64{W2, H_bottom}

	// Type-specific minimum height check
	meetsMin := true
	// heights := []float64{H0, H_bottom, H_bottom} // Heights for pics 0, 1, 2 respectively - Not strictly needed for check logic
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
		// Use correct height for check (H0 for pic 0, H_bottom for pics 1 & 2)
		checkHeight := H0
		if i > 0 {
			checkHeight = H_bottom
		}
		if checkHeight < requiredMinHeights[i] {
			meetsMin = false
			break
		}
	}
	// No explicit min width check based on rules

	return layout, meetsMin, nil
}
