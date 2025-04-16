package calculate

import "fmt"

// calculateLayout_7_1T2M2M2B calculates the 1T-2M-2M-2B layout.
func calculateLayout_7_1T2M2M2B(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 7), Dimensions: make([][]float64, 7)}
	if len(ARs) != 7 || len(types) != 7 {
		return layout, fmt.Errorf("1T2M2M2B layout requires 7 ARs and types")
	}

	// Row 1 (0)
	W0 := AW
	AR0 := ARs[0]
	height1 := 0.0
	if AR0 > 1e-6 {
		height1 = W0 / AR0
	} else {
		return layout, fmt.Errorf("invalid AR for pic 0 in 1T2M2M2B")
	}
	if height1 <= 1e-6 {
		return layout, fmt.Errorf("zero height for pic 0 in 1T2M2M2B")
	}

	// Row 2 (1, 2)
	widths2, height2, err2 := calculateRowLayout(ARs[1:3], AW, spacing)
	if err2 != nil {
		return layout, fmt.Errorf("failed row 2 for 1T2M2M2B: %w", err2)
	}
	W1, W2 := widths2[0], widths2[1]

	// Row 3 (3, 4)
	widths3, height3, err3 := calculateRowLayout(ARs[3:5], AW, spacing)
	if err3 != nil {
		return layout, fmt.Errorf("failed row 3 for 1T2M2M2B: %w", err3)
	}
	W3, W4 := widths3[0], widths3[1]

	// Row 4 (5, 6)
	widths4, height4, err4 := calculateRowLayout(ARs[5:7], AW, spacing)
	if err4 != nil {
		return layout, fmt.Errorf("failed row 4 for 1T2M2M2B: %w", err4)
	}
	W5, W6 := widths4[0], widths4[1]

	layout.TotalHeight = height1 + spacing + height2 + spacing + height3 + spacing + height4
	layout.TotalWidth = AW

	layout.Positions[0] = []float64{0, 0}
	layout.Dimensions[0] = []float64{W0, height1}
	yRow2 := height1 + spacing
	layout.Positions[1] = []float64{0, yRow2}
	layout.Dimensions[1] = []float64{W1, height2}
	layout.Positions[2] = []float64{W1 + spacing, yRow2}
	layout.Dimensions[2] = []float64{W2, height2}
	yRow3 := yRow2 + height2 + spacing
	layout.Positions[3] = []float64{0, yRow3}
	layout.Dimensions[3] = []float64{W3, height3}
	layout.Positions[4] = []float64{W3 + spacing, yRow3}
	layout.Dimensions[4] = []float64{W4, height3}
	yRow4 := yRow3 + height3 + spacing
	layout.Positions[5] = []float64{0, yRow4}
	layout.Dimensions[5] = []float64{W5, height4}
	layout.Positions[6] = []float64{W5 + spacing, yRow4}
	layout.Dimensions[6] = []float64{W6, height4}

	return layout, nil
}
