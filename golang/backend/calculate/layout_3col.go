package calculate

import "fmt"

// --- Helper function for 3 Columns (Vertical Stack) ---
func (e *ContinuousLayoutEngine) calculateLayout_3Col(ARs []float64, types []string, AW, spacing float64) (TemplateLayout, bool, error) {
	// Added definition based on layout_3_pics.go version
	layout := TemplateLayout{Positions: make([][]float64, 3), Dimensions: make([][]float64, 3)}
	AR0, AR1, AR2 := ARs[0], ARs[1], ARs[2]
	W0, W1, W2 := AW, AW, AW // Full width for each
	H0, H1, H2 := 0.0, 0.0, 0.0
	if AR0 > 1e-6 {
		H0 = W0 / AR0
	} else {
		return layout, false, fmt.Errorf("3Col calculated zero height for pic 0")
	}
	if AR1 > 1e-6 {
		H1 = W1 / AR1
	} else {
		return layout, false, fmt.Errorf("3Col calculated zero height for pic 1")
	}
	if AR2 > 1e-6 {
		H2 = W2 / AR2
	} else {
		return layout, false, fmt.Errorf("3Col calculated zero height for pic 2")
	}
	if H0 <= 1e-6 || H1 <= 1e-6 || H2 <= 1e-6 {
		return layout, false, fmt.Errorf("3Col calculated zero height")
	}
	heights := []float64{H0, H1, H2}
	layout.TotalHeight = H0 + H1 + H2 + 2*spacing
	layout.TotalWidth = AW
	currentY := 0.0
	for i := 0; i < 3; i++ {
		layout.Positions[i] = []float64{0, currentY}
		layout.Dimensions[i] = []float64{AW, heights[i]}
		if i < 2 {
			currentY += heights[i] + spacing
		}
	}

	// Type-specific minimum height check
	meetsMin := true
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
