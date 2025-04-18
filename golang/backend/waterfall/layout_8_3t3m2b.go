package waterfall

import "fmt"

// calculateLayout_8_3T3M2B: 3 Top, 3 Middle, 2 Bottom
func calculateLayout_8_3T3M2B(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 8), Dimensions: make([][]float64, 8)}
	if len(ARs) != 8 || len(types) != 8 {
		return layout, fmt.Errorf("3T3M2B layout requires 8 ARs/types")
	}

	widths1, h1, err1 := calculateRowLayout(ARs[0:3], AW, spacing)
	if err1 != nil {
		return layout, fmt.Errorf("row1 failed for 3T3M2B: %w", err1)
	}
	widths2, h2, err2 := calculateRowLayout(ARs[3:6], AW, spacing)
	if err2 != nil {
		return layout, fmt.Errorf("row2 failed for 3T3M2B: %w", err2)
	}
	widths3, h3, err3 := calculateRowLayout(ARs[6:8], AW, spacing)
	if err3 != nil {
		return layout, fmt.Errorf("row3 failed for 3T3M2B: %w", err3)
	}

	layout.TotalHeight = h1 + spacing + h2 + spacing + h3
	layout.TotalWidth = AW

	y := 0.0
	currentX := 0.0
	// Row 1
	for i := 0; i < 3; i++ {
		layout.Positions[i] = []float64{currentX, y}
		layout.Dimensions[i] = []float64{widths1[i], h1}
		currentX += widths1[i] + spacing
	}
	// Row 2
	y += h1 + spacing
	currentX = 0.0
	for i := 0; i < 3; i++ {
		layout.Positions[i+3] = []float64{currentX, y}
		layout.Dimensions[i+3] = []float64{widths2[i], h2}
		currentX += widths2[i] + spacing
	}
	// Row 3
	y += h2 + spacing
	currentX = 0.0
	for i := 0; i < 2; i++ {
		layout.Positions[i+6] = []float64{currentX, y}
		layout.Dimensions[i+6] = []float64{widths3[i], h3}
		currentX += widths3[i] + spacing
	}
	return layout, nil
}
