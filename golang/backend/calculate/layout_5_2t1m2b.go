package calculate

import "fmt"

// calculateLayout_5_2T1M2B calculates the 2 Top, 1 Middle, 2 Bottom layout.
func calculateLayout_5_2T1M2B(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 5), Dimensions: make([][]float64, 5)}
	if len(ARs) != 5 || len(types) != 5 {
		return layout, fmt.Errorf("2T1M2B layout requires 5 ARs and types")
	}

	// Row 1 (Pics 0, 1)
	widths1, height1, err1 := calculateRowLayout(ARs[0:2], AW, spacing)
	if err1 != nil {
		return layout, fmt.Errorf("failed to calculate top row for 2T1M2B: %w", err1)
	}
	W0, W1 := widths1[0], widths1[1]

	// Row 2 (Pic 2)
	W2 := AW
	AR2 := ARs[2]
	height2 := 0.0
	if AR2 > 1e-6 {
		height2 = W2 / AR2
	} else {
		return layout, fmt.Errorf("invalid aspect ratio (%.2f) for middle picture in 2T1M2B", AR2)
	}
	if height2 <= 1e-6 {
		return layout, fmt.Errorf("calculated zero height (%.2f) for middle picture in 2T1M2B", height2)
	}

	// Row 3 (Pics 3, 4)
	widths3, height3, err3 := calculateRowLayout(ARs[3:5], AW, spacing)
	if err3 != nil {
		return layout, fmt.Errorf("failed to calculate bottom row for 2T1M2B: %w", err3)
	}
	W3, W4 := widths3[0], widths3[1]

	layout.TotalHeight = height1 + spacing + height2 + spacing + height3
	layout.TotalWidth = AW

	// Positions and Dimensions
	layout.Positions[0] = []float64{0, 0}
	layout.Dimensions[0] = []float64{W0, height1}
	layout.Positions[1] = []float64{W0 + spacing, 0}
	layout.Dimensions[1] = []float64{W1, height1}

	yRow2 := height1 + spacing
	layout.Positions[2] = []float64{0, yRow2}
	layout.Dimensions[2] = []float64{W2, height2}

	yRow3 := yRow2 + height2 + spacing
	layout.Positions[3] = []float64{0, yRow3}
	layout.Dimensions[3] = []float64{W3, height3}
	layout.Positions[4] = []float64{W3 + spacing, yRow3}
	layout.Dimensions[4] = []float64{W4, height3}

	return layout, nil
}
