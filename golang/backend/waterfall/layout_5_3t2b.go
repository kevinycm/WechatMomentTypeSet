package waterfall

import "fmt"

// calculateLayout_5_3T2B calculates the 3 Top, 2 Bottom layout.
func calculateLayout_5_3T2B(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 5), Dimensions: make([][]float64, 5)}
	if len(ARs) != 5 || len(types) != 5 {
		return layout, fmt.Errorf("3T2B layout requires 5 ARs and types")
	}

	widths1, height1, err1 := calculateRowLayout(ARs[0:3], AW, spacing)
	if err1 != nil {
		return layout, fmt.Errorf("failed to calculate top row for 3T2B: %w", err1)
	}
	W0, W1, W2 := widths1[0], widths1[1], widths1[2]

	widths2, height2, err2 := calculateRowLayout(ARs[3:5], AW, spacing)
	if err2 != nil {
		return layout, fmt.Errorf("failed to calculate bottom row for 3T2B: %w", err2)
	}
	W3, W4 := widths2[0], widths2[1]

	layout.TotalHeight = height1 + spacing + height2
	layout.TotalWidth = AW

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

	return layout, nil
}
