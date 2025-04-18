package waterfall

import "fmt"

// calculateLayout_8_1T2M2M3B: 1 Top, 2 Mid1, 2 Mid2, 3 Bottom
func calculateLayout_8_1T2M2M3B(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 8), Dimensions: make([][]float64, 8)}
	if len(ARs) != 8 || len(types) != 8 {
		return layout, fmt.Errorf("1T2M2M3B layout requires 8 ARs/types")
	}

	// Row 1 (Single Pic)
	w1 := AW
	ar1 := ARs[0]
	h1 := 0.0
	if ar1 <= 1e-6 {
		return layout, fmt.Errorf("invalid AR for top pic in 1T2M2M3B")
	}
	h1 = w1 / ar1
	if h1 <= 1e-6 {
		return layout, fmt.Errorf("zero height for top pic in 1T2M2M3B")
	}

	widths2, h2, err2 := calculateRowLayout(ARs[1:3], AW, spacing)
	if err2 != nil {
		return layout, fmt.Errorf("row2 failed for 1T2M2M3B: %w", err2)
	}
	widths3, h3, err3 := calculateRowLayout(ARs[3:5], AW, spacing)
	if err3 != nil {
		return layout, fmt.Errorf("row3 failed for 1T2M2M3B: %w", err3)
	}
	widths4, h4, err4 := calculateRowLayout(ARs[5:8], AW, spacing)
	if err4 != nil {
		return layout, fmt.Errorf("row4 failed for 1T2M2M3B: %w", err4)
	}

	layout.TotalHeight = h1 + spacing + h2 + spacing + h3 + spacing + h4
	layout.TotalWidth = AW

	y := 0.0
	// Row 1
	layout.Positions[0] = []float64{0, y}
	layout.Dimensions[0] = []float64{w1, h1}
	// Row 2
	y += h1 + spacing
	currentX := 0.0
	for i := 0; i < 2; i++ {
		layout.Positions[i+1] = []float64{currentX, y}
		layout.Dimensions[i+1] = []float64{widths2[i], h2}
		currentX += widths2[i] + spacing
	}
	// Row 3
	y += h2 + spacing
	currentX = 0.0
	for i := 0; i < 2; i++ {
		layout.Positions[i+3] = []float64{currentX, y}
		layout.Dimensions[i+3] = []float64{widths3[i], h3}
		currentX += widths3[i] + spacing
	}
	// Row 4
	y += h3 + spacing
	currentX = 0.0
	for i := 0; i < 3; i++ {
		layout.Positions[i+5] = []float64{currentX, y}
		layout.Dimensions[i+5] = []float64{widths4[i], h4}
		currentX += widths4[i] + spacing
	}

	return layout, nil
}
