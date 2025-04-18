package waterfall

import "fmt"

// calculateLayout_9_3T3M3B: 3 Top, 3 Middle, 3 Bottom (3x3 Grid)
func calculateLayout_9_3T3M3B(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 9), Dimensions: make([][]float64, 9)}
	if len(ARs) != 9 || len(types) != 9 {
		return layout, fmt.Errorf("3T3M3B requires 9 ARs/types")
	}

	widths1, h1, err1 := calculateRowLayout(ARs[0:3], AW, spacing)
	if err1 != nil {
		return layout, fmt.Errorf("row1 fail: %w", err1)
	}
	widths2, h2, err2 := calculateRowLayout(ARs[3:6], AW, spacing)
	if err2 != nil {
		return layout, fmt.Errorf("row2 fail: %w", err2)
	}
	widths3, h3, err3 := calculateRowLayout(ARs[6:9], AW, spacing)
	if err3 != nil {
		return layout, fmt.Errorf("row3 fail: %w", err3)
	}

	layout.TotalHeight = h1 + spacing + h2 + spacing + h3
	layout.TotalWidth = AW

	y := 0.0
	currentX := 0.0
	for i := 0; i < 3; i++ {
		layout.Positions[i] = []float64{currentX, y}
		layout.Dimensions[i] = []float64{widths1[i], h1}
		currentX += widths1[i] + spacing
	}
	y += h1 + spacing
	currentX = 0.0
	for i := 0; i < 3; i++ {
		layout.Positions[i+3] = []float64{currentX, y}
		layout.Dimensions[i+3] = []float64{widths2[i], h2}
		currentX += widths2[i] + spacing
	}
	y += h2 + spacing
	currentX = 0.0
	for i := 0; i < 3; i++ {
		layout.Positions[i+6] = []float64{currentX, y}
		layout.Dimensions[i+6] = []float64{widths3[i], h3}
		currentX += widths3[i] + spacing
	}

	return layout, nil
}
