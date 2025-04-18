package waterfall

import "fmt"

// calculateLayout_9_2T2M3M2B: 2 Top, 2 Mid1, 3 Mid2, 2 Bottom
func calculateLayout_9_2T2M3M2B(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 9), Dimensions: make([][]float64, 9)}
	if len(ARs) != 9 || len(types) != 9 {
		return layout, fmt.Errorf("2T2M3M2B requires 9 ARs/types")
	}

	widths1, h1, err1 := calculateRowLayout(ARs[0:2], AW, spacing)
	if err1 != nil {
		return layout, fmt.Errorf("row1 fail: %w", err1)
	}
	widths2, h2, err2 := calculateRowLayout(ARs[2:4], AW, spacing)
	if err2 != nil {
		return layout, fmt.Errorf("row2 fail: %w", err2)
	}
	widths3, h3, err3 := calculateRowLayout(ARs[4:7], AW, spacing)
	if err3 != nil {
		return layout, fmt.Errorf("row3 fail: %w", err3)
	}
	widths4, h4, err4 := calculateRowLayout(ARs[7:9], AW, spacing)
	if err4 != nil {
		return layout, fmt.Errorf("row4 fail: %w", err4)
	}

	layout.TotalHeight = h1 + spacing + h2 + spacing + h3 + spacing + h4
	layout.TotalWidth = AW

	y := 0.0
	currentX := 0.0
	for i := 0; i < 2; i++ {
		layout.Positions[i] = []float64{currentX, y}
		layout.Dimensions[i] = []float64{widths1[i], h1}
		currentX += widths1[i] + spacing
	}
	y += h1 + spacing
	currentX = 0.0
	for i := 0; i < 2; i++ {
		layout.Positions[i+2] = []float64{currentX, y}
		layout.Dimensions[i+2] = []float64{widths2[i], h2}
		currentX += widths2[i] + spacing
	}
	y += h2 + spacing
	currentX = 0.0
	for i := 0; i < 3; i++ {
		layout.Positions[i+4] = []float64{currentX, y}
		layout.Dimensions[i+4] = []float64{widths3[i], h3}
		currentX += widths3[i] + spacing
	}
	y += h3 + spacing
	currentX = 0.0
	for i := 0; i < 2; i++ {
		layout.Positions[i+7] = []float64{currentX, y}
		layout.Dimensions[i+7] = []float64{widths4[i], h4}
		currentX += widths4[i] + spacing
	}

	return layout, nil
}
