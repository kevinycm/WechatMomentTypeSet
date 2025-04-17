package calculate

import "fmt"

// calculateLayout_7_1T3M3B calculates the 1 Top, 3 Middle, 3 Bottom layout.
func calculateLayout_7_1T3M3B(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 7), Dimensions: make([][]float64, 7)}
	if len(ARs) != 7 || len(types) != 7 {
		return layout, fmt.Errorf("1T3M3B layout requires 7 ARs and types")
	}

	// Row 1 (Single Pic)
	W0 := AW
	AR0 := ARs[0]
	height1 := 0.0
	if AR0 > 1e-6 {
		height1 = W0 / AR0
	} else {
		return layout, fmt.Errorf("invalid AR for pic 0 in 1T3M3B")
	}
	if height1 <= 1e-6 {
		return layout, fmt.Errorf("zero height for pic 0 in 1T3M3B")
	}
	widths2, height2, err2 := calculateRowLayout(ARs[1:4], AW, spacing)
	if err2 != nil {
		return layout, fmt.Errorf("failed row 2 for 1T3M3B: %w", err2)
	}
	widths3, height3, err3 := calculateRowLayout(ARs[4:7], AW, spacing)
	if err3 != nil {
		return layout, fmt.Errorf("failed row 3 for 1T3M3B: %w", err3)
	}

	layout.TotalHeight = height1 + spacing + height2 + spacing + height3
	layout.TotalWidth = AW

	layout.Positions[0] = []float64{0, 0}
	layout.Dimensions[0] = []float64{W0, height1}
	yRow2 := height1 + spacing
	currentX := 0.0
	for i := 0; i < 3; i++ {
		layout.Positions[i+1] = []float64{currentX, yRow2}
		layout.Dimensions[i+1] = []float64{widths2[i], height2}
		currentX += widths2[i] + spacing
	}
	yRow3 := yRow2 + height2 + spacing
	currentX = 0.0
	for i := 0; i < 3; i++ {
		layout.Positions[i+4] = []float64{currentX, yRow3}
		layout.Dimensions[i+4] = []float64{widths3[i], height3}
		currentX += widths3[i] + spacing
	}

	return layout, nil
}
