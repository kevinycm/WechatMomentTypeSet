package calculate

import "fmt"

// calculateLayout_8_2T3M2M1B: 2 Top, 3 Mid1, 2 Mid2, 1 Bottom
func calculateLayout_8_2T3M2M1B(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 8), Dimensions: make([][]float64, 8)}
	if len(ARs) != 8 || len(types) != 8 {
		return layout, fmt.Errorf("2T3M2M1B layout requires 8 ARs/types")
	}

	widths1, h1, err1 := calculateRowLayout(ARs[0:2], AW, spacing)
	if err1 != nil {
		return layout, fmt.Errorf("row1 failed for 2T3M2M1B: %w", err1)
	}
	widths2, h2, err2 := calculateRowLayout(ARs[2:5], AW, spacing)
	if err2 != nil {
		return layout, fmt.Errorf("row2 failed for 2T3M2M1B: %w", err2)
	}
	widths3, h3, err3 := calculateRowLayout(ARs[5:7], AW, spacing)
	if err3 != nil {
		return layout, fmt.Errorf("row3 failed for 2T3M2M1B: %w", err3)
	}
	// Row 4 (Single Pic)
	w4 := AW
	ar4 := ARs[7]
	h4 := 0.0
	if ar4 <= 1e-6 {
		return layout, fmt.Errorf("invalid AR for bottom pic in 2T3M2M1B")
	}
	h4 = w4 / ar4
	if h4 <= 1e-6 {
		return layout, fmt.Errorf("zero height for bottom pic in 2T3M2M1B")
	}

	layout.TotalHeight = h1 + spacing + h2 + spacing + h3 + spacing + h4
	layout.TotalWidth = AW

	y := 0.0
	currentX := 0.0
	// Row 1
	for i := 0; i < 2; i++ {
		layout.Positions[i] = []float64{currentX, y}
		layout.Dimensions[i] = []float64{widths1[i], h1}
		currentX += widths1[i] + spacing
	}
	// Row 2
	y += h1 + spacing
	currentX = 0.0
	for i := 0; i < 3; i++ {
		layout.Positions[i+2] = []float64{currentX, y}
		layout.Dimensions[i+2] = []float64{widths2[i], h2}
		currentX += widths2[i] + spacing
	}
	// Row 3
	y += h2 + spacing
	currentX = 0.0
	for i := 0; i < 2; i++ {
		layout.Positions[i+5] = []float64{currentX, y}
		layout.Dimensions[i+5] = []float64{widths3[i], h3}
		currentX += widths3[i] + spacing
	}
	// Row 4
	y += h3 + spacing
	layout.Positions[7] = []float64{0, y}
	layout.Dimensions[7] = []float64{w4, h4}

	return layout, nil
}
