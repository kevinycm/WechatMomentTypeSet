package calculate

import "fmt"

// calculateLayout_7_1T2M3M1B calculates the 1T-2M-3M-1B layout.
func calculateLayout_7_1T2M3M1B(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 7), Dimensions: make([][]float64, 7)}
	if len(ARs) != 7 || len(types) != 7 {
		return layout, fmt.Errorf("1T2M3M1B layout requires 7 ARs and types")
	}

	// Row 1 (0)
	W0 := AW
	AR0 := ARs[0]
	height1 := 0.0
	if AR0 > 1e-6 {
		height1 = W0 / AR0
	} else {
		return layout, fmt.Errorf("invalid AR for pic 0 in 1T2M3M1B")
	}
	if height1 <= 1e-6 {
		return layout, fmt.Errorf("zero height for pic 0 in 1T2M3M1B")
	}

	// Row 2 (1, 2)
	widths2, height2, err2 := calculateRowLayout(ARs[1:3], AW, spacing)
	if err2 != nil {
		return layout, fmt.Errorf("failed row 2 for 1T2M3M1B: %w", err2)
	}
	W1, W2 := widths2[0], widths2[1]

	// Row 3 (3, 4, 5)
	widths3, height3, err3 := calculateRowLayout(ARs[3:6], AW, spacing)
	if err3 != nil {
		return layout, fmt.Errorf("failed row 3 for 1T2M3M1B: %w", err3)
	}
	W3, W4, W5 := widths3[0], widths3[1], widths3[2]

	// Row 4 (6)
	W6 := AW
	AR6 := ARs[6]
	height4 := 0.0
	if AR6 > 1e-6 {
		height4 = W6 / AR6
	} else {
		return layout, fmt.Errorf("invalid AR for pic 6 in 1T2M3M1B")
	}
	if height4 <= 1e-6 {
		return layout, fmt.Errorf("zero height for pic 6 in 1T2M3M1B")
	}

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
	currentX := 0.0
	layout.Positions[3] = []float64{currentX, yRow3}
	layout.Dimensions[3] = []float64{W3, height3}
	currentX += W3 + spacing
	layout.Positions[4] = []float64{currentX, yRow3}
	layout.Dimensions[4] = []float64{W4, height3}
	currentX += W4 + spacing
	layout.Positions[5] = []float64{currentX, yRow3}
	layout.Dimensions[5] = []float64{W5, height3}
	yRow4 := yRow3 + height3 + spacing
	layout.Positions[6] = []float64{0, yRow4}
	layout.Dimensions[6] = []float64{W6, height4}

	return layout, nil
}
