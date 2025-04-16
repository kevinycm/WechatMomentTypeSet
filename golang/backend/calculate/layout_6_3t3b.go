package calculate

import "fmt"

// calculateLayout_6_3T3B calculates the 3 Top, 3 Bottom layout.
func calculateLayout_6_3T3B(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 6), Dimensions: make([][]float64, 6)}
	if len(ARs) != 6 || len(types) != 6 {
		return layout, fmt.Errorf("3T3B layout requires 6 ARs and types")
	}

	// Row 1 (Pics 0, 1, 2)
	widths1, height1, err1 := calculateRowLayout(ARs[0:3], AW, spacing) // Use global helper
	if err1 != nil {
		return layout, fmt.Errorf("failed to calculate top row for 3T3B: %w", err1)
	}
	W0, W1, W2 := widths1[0], widths1[1], widths1[2]

	// Row 2 (Pics 3, 4, 5)
	widths2, height2, err2 := calculateRowLayout(ARs[3:6], AW, spacing) // Use global helper
	if err2 != nil {
		return layout, fmt.Errorf("failed to calculate bottom row for 3T3B: %w", err2)
	}
	W3, W4, W5 := widths2[0], widths2[1], widths2[2]

	layout.TotalHeight = height1 + spacing + height2
	layout.TotalWidth = AW

	// Positions and Dimensions
	currentX := 0.0
	layout.Positions[0] = []float64{currentX, 0}
	layout.Dimensions[0] = []float64{W0, height1}
	currentX += W0 + spacing
	layout.Positions[1] = []float64{currentX, 0}
	layout.Dimensions[1] = []float64{W1, height1}
	currentX += W1 + spacing
	layout.Positions[2] = []float64{currentX, 0}
	layout.Dimensions[2] = []float64{W2, height1}

	bottomY := height1 + spacing
	currentX = 0.0
	layout.Positions[3] = []float64{currentX, bottomY}
	layout.Dimensions[3] = []float64{W3, height2}
	currentX += W3 + spacing
	layout.Positions[4] = []float64{currentX, bottomY}
	layout.Dimensions[4] = []float64{W4, height2}
	currentX += W4 + spacing
	layout.Positions[5] = []float64{currentX, bottomY}
	layout.Dimensions[5] = []float64{W5, height2}

	return layout, nil
}
