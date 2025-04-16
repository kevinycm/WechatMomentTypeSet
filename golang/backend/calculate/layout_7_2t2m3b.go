package calculate

import "fmt"

// calculateLayout_7_2T2M3B calculates the 2 Top, 2 Middle, 3 Bottom layout.
func calculateLayout_7_2T2M3B(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 7), Dimensions: make([][]float64, 7)}
	if len(ARs) != 7 || len(types) != 7 {
		return layout, fmt.Errorf("2T2M3B layout requires 7 ARs and types")
	}

	// Row 1 (0, 1)
	widths1, height1, err1 := calculateRowLayout(ARs[0:2], AW, spacing)
	if err1 != nil {
		return layout, fmt.Errorf("failed top row for 2T2M3B: %w", err1)
	}
	W0, W1 := widths1[0], widths1[1]

	// Row 2 (2, 3)
	widths2, height2, err2 := calculateRowLayout(ARs[2:4], AW, spacing)
	if err2 != nil {
		return layout, fmt.Errorf("failed middle row for 2T2M3B: %w", err2)
	}
	W2, W3 := widths2[0], widths2[1]

	// Row 3 (4, 5, 6)
	widths3, height3, err3 := calculateRowLayout(ARs[4:7], AW, spacing)
	if err3 != nil {
		return layout, fmt.Errorf("failed bottom row for 2T2M3B: %w", err3)
	}
	W4, W5, W6 := widths3[0], widths3[1], widths3[2]

	layout.TotalHeight = height1 + spacing + height2 + spacing + height3
	layout.TotalWidth = AW

	layout.Positions[0] = []float64{0, 0}
	layout.Dimensions[0] = []float64{W0, height1}
	layout.Positions[1] = []float64{W0 + spacing, 0}
	layout.Dimensions[1] = []float64{W1, height1}
	yRow2 := height1 + spacing
	layout.Positions[2] = []float64{0, yRow2}
	layout.Dimensions[2] = []float64{W2, height2}
	layout.Positions[3] = []float64{W2 + spacing, yRow2}
	layout.Dimensions[3] = []float64{W3, height2}
	yRow3 := yRow2 + height2 + spacing
	currentX := 0.0
	layout.Positions[4] = []float64{currentX, yRow3}
	layout.Dimensions[4] = []float64{W4, height3}
	currentX += W4 + spacing
	layout.Positions[5] = []float64{currentX, yRow3}
	layout.Dimensions[5] = []float64{W5, height3}
	currentX += W5 + spacing
	layout.Positions[6] = []float64{currentX, yRow3}
	layout.Dimensions[6] = []float64{W6, height3}

	return layout, nil
}
