package calculate

import "fmt"

// --- Helper function for 3 in a Row Template ---
func (e *ContinuousLayoutEngine) calculateLayout_3Row(ARs []float64, types []string, AW, spacing float64) (TemplateLayout, bool, error) {
	layout := TemplateLayout{
		Positions:  make([][]float64, 3),
		Dimensions: make([][]float64, 3),
	}
	AR0, AR1, AR2 := ARs[0], ARs[1], ARs[2]

	rowAvailableWidth := AW - 2*spacing
	totalARSum := AR0 + AR1 + AR2
	H := 0.0
	if rowAvailableWidth > 1e-6 && totalARSum > 1e-6 {
		H = rowAvailableWidth / totalARSum
	} else {
		return layout, false, fmt.Errorf("cannot calculate 3-in-a-row layout (zero width or AR sum)")
	}
	if H <= 1e-6 {
		return layout, false, fmt.Errorf("3-in-a-row calculated zero height")
	}

	W0 := H * AR0
	W1 := H * AR1
	W2 := H * AR2
	widths := []float64{W0, W1, W2}

	layout.TotalHeight = H
	layout.TotalWidth = AW
	currentX := 0.0
	for i := 0; i < 3; i++ {
		layout.Positions[i] = []float64{currentX, 0}
		layout.Dimensions[i] = []float64{widths[i], H}
		if i < 2 {
			currentX += widths[i] + spacing
		}
	}

	// Type-specific minimum height check
	meetsMin := true
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
		// All pictures have the same height H in this layout
		if H < requiredMinHeights[i] {
			meetsMin = false
			break
		}
	}
	// No explicit min width check based on rules

	return layout, meetsMin, nil
}
