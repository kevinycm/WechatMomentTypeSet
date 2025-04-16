package calculate

import "fmt"

// --- Helper function for 1 Left, 2 Right Stacked Template ---
func (e *ContinuousLayoutEngine) calculateLayout_1L2R(ARs []float64, types []string, AW, spacing float64) (TemplateLayout, bool, error) {
	// Added definition based on layout_3_pics.go version
	layout := TemplateLayout{
		Positions:  make([][]float64, 3),
		Dimensions: make([][]float64, 3),
	}
	AR0 := ARs[0]
	denominator := 1.0
	if ARs[1] > 1e-6 {
		denominator += AR0 / ARs[1]
	}
	if ARs[2] > 1e-6 {
		denominator += AR0 / ARs[2]
	}
	WR := 0.0
	if denominator > 1e-6 {
		WR = (AW - spacing*(1.0+AR0)) / denominator
	}
	if WR <= 1e-6 || WR > AW-spacing+1e-6 {
		return layout, false, fmt.Errorf("1L2R geometry infeasible (WR=%.2f)", WR)
	}
	W0 := AW - spacing - WR
	if W0 <= 1e-6 {
		return layout, false, fmt.Errorf("1L2R geometry infeasible (W0=%.2f)", W0)
	}
	H0 := 0.0
	if AR0 > 1e-6 {
		H0 = W0 / AR0
	}
	H1 := 0.0
	if ARs[1] > 1e-6 {
		H1 = WR / ARs[1]
	}
	H2 := 0.0
	if ARs[2] > 1e-6 {
		H2 = WR / ARs[2]
	}
	if H0 <= 1e-6 || H1 <= 1e-6 || H2 <= 1e-6 {
		return layout, false, fmt.Errorf("1L2R calculated zero height")
	}

	layout.TotalHeight = H0
	layout.TotalWidth = AW
	layout.Positions[0] = []float64{0, 0}
	layout.Dimensions[0] = []float64{W0, H0}
	layout.Positions[1] = []float64{W0 + spacing, 0}
	layout.Dimensions[1] = []float64{WR, H1}
	layout.Positions[2] = []float64{W0 + spacing, H1 + spacing}
	layout.Dimensions[2] = []float64{WR, H2}

	// Type-specific minimum height check
	meetsMin := true
	heights := []float64{H0, H1, H2} // Heights for pics 0, 1, 2
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
		if heights[i] < requiredMinHeights[i] {
			meetsMin = false
			break
		}
	}
	// No explicit min width check based on rules

	return layout, meetsMin, nil
}
