package waterfall

import "fmt"

// --- Helper function for 2 Left Stacked, 1 Right Template --- (Mirror of 1L2R)
func (e *ContinuousLayoutEngine) calculateLayout_2L1R(ARs []float64, types []string, AW, spacing float64) (TemplateLayout, bool, error) {
	layout := TemplateLayout{
		Positions:  make([][]float64, 3),
		Dimensions: make([][]float64, 3),
	}
	AR0, AR1, AR2 := ARs[0], ARs[1], ARs[2]

	// Calculate WL based on geometry (similar to WR in 1L2R)
	denominator := 1.0
	if AR0 > 1e-6 {
		denominator += AR2 / AR0
	}
	if AR1 > 1e-6 {
		denominator += AR2 / AR1
	}

	WL := 0.0
	if denominator > 1e-6 {
		WL = (AW - spacing*(1.0+AR2)) / denominator
	}

	if WL <= 1e-6 || WL > AW-spacing+1e-6 {
		return layout, false, fmt.Errorf("2L1R geometry infeasible (WL=%.2f)", WL)
	}

	W2 := AW - spacing - WL
	if W2 <= 1e-6 {
		return layout, false, fmt.Errorf("2L1R geometry infeasible (W2=%.2f)", W2)
	}

	H0 := 0.0
	if AR0 > 1e-6 {
		H0 = WL / AR0
	}
	H1 := 0.0
	if AR1 > 1e-6 {
		H1 = WL / AR1
	}
	H2 := 0.0
	if AR2 > 1e-6 {
		H2 = W2 / AR2
	}

	if H0 <= 1e-6 || H1 <= 1e-6 || H2 <= 1e-6 {
		return layout, false, fmt.Errorf("2L1R calculated zero height")
	}

	layout.TotalHeight = H2
	layout.TotalWidth = AW
	layout.Positions[0] = []float64{0, 0}
	layout.Dimensions[0] = []float64{WL, H0}
	layout.Positions[1] = []float64{0, H0 + spacing}
	layout.Dimensions[1] = []float64{WL, H1}
	layout.Positions[2] = []float64{WL + spacing, 0}
	layout.Dimensions[2] = []float64{W2, H2}

	// Type-specific minimum height check
	meetsMin := true
	heights := []float64{H0, H1, H2} // Heights for pics 0, 1, 2
	requiredMinHeights := make([]float64, 3)
	for i, picType := range types {
		// Use the GetRequiredMinHeight helper function
		requiredMinHeight := GetRequiredMinHeight(e, picType, 3)
		requiredMinHeights[i] = requiredMinHeight // Store it (might be redundant)

		if heights[i] < requiredMinHeights[i] {
			meetsMin = false
			break
		}
	}
	// No explicit min width check based on rules

	return layout, meetsMin, nil
}
