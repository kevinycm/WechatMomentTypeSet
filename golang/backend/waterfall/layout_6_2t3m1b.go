package waterfall

import "fmt"

// calculateLayout_6_2T3M1B calculates the 2 Top, 3 Middle, 1 Bottom layout.
func calculateLayout_6_2T3M1B(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 6), Dimensions: make([][]float64, 6)}
	if len(ARs) != 6 || len(types) != 6 {
		return layout, fmt.Errorf("2T3M1B layout requires 6 ARs and types")
	}

	// Row 1 (Pics 0, 1)
	widths1, height1, err1 := calculateRowLayout(ARs[0:2], AW, spacing)
	if err1 != nil {
		return layout, fmt.Errorf("failed top row for 2T3M1B: %w", err1)
	}
	W0, W1 := widths1[0], widths1[1]

	// Row 2 (Pics 2, 3, 4)
	widths2, height2, err2 := calculateRowLayout(ARs[2:5], AW, spacing)
	if err2 != nil {
		return layout, fmt.Errorf("failed middle row for 2T3M1B: %w", err2)
	}
	W2, W3, W4 := widths2[0], widths2[1], widths2[2]

	// Row 3 (Pic 5)
	W5 := AW
	AR5 := ARs[5]
	height3 := 0.0
	if AR5 > 1e-6 {
		height3 = W5 / AR5
	} else {
		return layout, fmt.Errorf("invalid AR for bottom pic in 2T3M1B")
	}
	if height3 <= 1e-6 {
		return layout, fmt.Errorf("zero height for bottom pic in 2T3M1B")
	}

	layout.TotalHeight = height1 + spacing + height2 + spacing + height3
	layout.TotalWidth = AW

	layout.Positions[0] = []float64{0, 0}
	layout.Dimensions[0] = []float64{W0, height1}
	layout.Positions[1] = []float64{W0 + spacing, 0}
	layout.Dimensions[1] = []float64{W1, height1}
	yRow2 := height1 + spacing
	currentX := 0.0
	layout.Positions[2] = []float64{currentX, yRow2}
	layout.Dimensions[2] = []float64{W2, height2}
	currentX += W2 + spacing
	layout.Positions[3] = []float64{currentX, yRow2}
	layout.Dimensions[3] = []float64{W3, height2}
	currentX += W3 + spacing
	layout.Positions[4] = []float64{currentX, yRow2}
	layout.Dimensions[4] = []float64{W4, height2}
	yRow3 := yRow2 + height2 + spacing
	layout.Positions[5] = []float64{0, yRow3}
	layout.Dimensions[5] = []float64{W5, height3}

	return layout, nil
}
