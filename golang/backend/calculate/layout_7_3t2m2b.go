package calculate

import "fmt"

// calculateLayout_7_3T2M2B calculates the 3 Top, 2 Middle, 2 Bottom layout.
func calculateLayout_7_3T2M2B(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 7), Dimensions: make([][]float64, 7)}
	if len(ARs) != 7 || len(types) != 7 {
		return layout, fmt.Errorf("3T2M2B layout requires 7 ARs and types")
	}

	// Row 1 (0, 1, 2)
	widths1, height1, err1 := calculateRowLayout(ARs[0:3], AW, spacing)
	if err1 != nil {
		return layout, fmt.Errorf("failed top row for 3T2M2B: %w", err1)
	}
	W0, W1, W2 := widths1[0], widths1[1], widths1[2]

	// Row 2 (3, 4)
	widths2, height2, err2 := calculateRowLayout(ARs[3:5], AW, spacing)
	if err2 != nil {
		return layout, fmt.Errorf("failed middle row for 3T2M2B: %w", err2)
	}
	W3, W4 := widths2[0], widths2[1]

	// Row 3 (5, 6)
	widths3, height3, err3 := calculateRowLayout(ARs[5:7], AW, spacing)
	if err3 != nil {
		return layout, fmt.Errorf("failed bottom row for 3T2M2B: %w", err3)
	}
	W5, W6 := widths3[0], widths3[1]

	layout.TotalHeight = height1 + spacing + height2 + spacing + height3
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
	yRow2 := height1 + spacing
	layout.Positions[3] = []float64{0, yRow2}
	layout.Dimensions[3] = []float64{W3, height2}
	layout.Positions[4] = []float64{W3 + spacing, yRow2}
	layout.Dimensions[4] = []float64{W4, height2}
	yRow3 := yRow2 + height2 + spacing
	layout.Positions[5] = []float64{0, yRow3}
	layout.Dimensions[5] = []float64{W5, height3}
	layout.Positions[6] = []float64{W5 + spacing, yRow3}
	layout.Dimensions[6] = []float64{W6, height3}

	return layout, nil
}
