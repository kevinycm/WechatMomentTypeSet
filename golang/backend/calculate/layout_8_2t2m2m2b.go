package calculate

import "fmt"

// calculateLayout_8_2T2M2M2B: 2 Top, 2 Mid1, 2 Mid2, 2 Bottom
func calculateLayout_8_2T2M2M2B(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 8), Dimensions: make([][]float64, 8)}
	if len(ARs) != 8 || len(types) != 8 {
		return layout, fmt.Errorf("2T2M2M2B layout requires 8 ARs/types")
	}

	widths1, h1, err1 := calculateRowLayout(ARs[0:2], AW, spacing)
	if err1 != nil {
		return layout, fmt.Errorf("row1 failed for 2T2M2M2B: %w", err1)
	}
	widths2, h2, err2 := calculateRowLayout(ARs[2:4], AW, spacing)
	if err2 != nil {
		return layout, fmt.Errorf("row2 failed for 2T2M2M2B: %w", err2)
	}
	widths3, h3, err3 := calculateRowLayout(ARs[4:6], AW, spacing)
	if err3 != nil {
		return layout, fmt.Errorf("row3 failed for 2T2M2M2B: %w", err3)
	}
	widths4, h4, err4 := calculateRowLayout(ARs[6:8], AW, spacing)
	if err4 != nil {
		return layout, fmt.Errorf("row4 failed for 2T2M2M2B: %w", err4)
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
	for i := 0; i < 2; i++ {
		layout.Positions[i+2] = []float64{currentX, y}
		layout.Dimensions[i+2] = []float64{widths2[i], h2}
		currentX += widths2[i] + spacing
	}
	// Row 3
	y += h2 + spacing
	currentX = 0.0
	for i := 0; i < 2; i++ {
		layout.Positions[i+4] = []float64{currentX, y}
		layout.Dimensions[i+4] = []float64{widths3[i], h3}
		currentX += widths3[i] + spacing
	}
	// Row 4
	y += h3 + spacing
	currentX = 0.0
	for i := 0; i < 2; i++ {
		layout.Positions[i+6] = []float64{currentX, y}
		layout.Dimensions[i+6] = []float64{widths4[i], h4}
		currentX += widths4[i] + spacing
	}

	return layout, nil
}
