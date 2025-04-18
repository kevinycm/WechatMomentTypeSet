package waterfall

import "fmt"

// calculateLayout_7_3T1M3B calculates the 3 Top, 1 Middle, 3 Bottom layout.
func calculateLayout_7_3T1M3B(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 7), Dimensions: make([][]float64, 7)}
	if len(ARs) != 7 || len(types) != 7 {
		return layout, fmt.Errorf("3T1M3B layout requires 7 ARs and types")
	}

	// Row 1 (Pics 0, 1, 2)
	widths1, height1, err1 := calculateRowLayout(ARs[0:3], AW, spacing)
	if err1 != nil {
		return layout, fmt.Errorf("failed row 1 for 3T1M3B: %w", err1)
	}
	// Row 2 (Pic 3 - Single)
	W3 := AW
	AR3 := ARs[3]
	height2 := 0.0
	if AR3 > 1e-6 {
		height2 = W3 / AR3
	} else {
		return layout, fmt.Errorf("invalid AR for pic 3 in 3T1M3B")
	}
	if height2 <= 1e-6 {
		return layout, fmt.Errorf("zero height for pic 3 in 3T1M3B")
	}
	// Row 3 (Pics 4, 5, 6)
	widths3, height3, err3 := calculateRowLayout(ARs[4:7], AW, spacing)
	if err3 != nil {
		return layout, fmt.Errorf("failed row 3 for 3T1M3B: %w", err3)
	}

	layout.TotalHeight = height1 + spacing + height2 + spacing + height3
	layout.TotalWidth = AW

	y := 0.0
	currentX := 0.0
	// Row 1
	for i := 0; i < 3; i++ {
		layout.Positions[i] = []float64{currentX, y}
		layout.Dimensions[i] = []float64{widths1[i], height1}
		currentX += widths1[i] + spacing
	}
	// Row 2
	y += height1 + spacing
	layout.Positions[3] = []float64{0, y} // Single picture starts at x=0
	layout.Dimensions[3] = []float64{W3, height2}
	// Row 3
	y += height2 + spacing
	currentX = 0.0
	for i := 0; i < 3; i++ {
		layout.Positions[i+4] = []float64{currentX, y}
		layout.Dimensions[i+4] = []float64{widths3[i], height3}
		currentX += widths3[i] + spacing
	}

	return layout, nil
}
