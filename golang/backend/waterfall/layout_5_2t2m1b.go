package waterfall

import "fmt"

// calculateLayout_5_2T2M1B calculates the 2 Top, 2 Middle, 1 Bottom layout.
func calculateLayout_5_2T2M1B(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 5), Dimensions: make([][]float64, 5)}
	if len(ARs) != 5 || len(types) != 5 {
		return layout, fmt.Errorf("2T2M1B layout requires 5 ARs and types")
	}

	// Row 1 (Pics 0, 1)
	widths1, height1, err1 := calculateRowLayout(ARs[0:2], AW, spacing)
	if err1 != nil {
		return layout, fmt.Errorf("failed to calculate top row for 2T2M1B: %w", err1)
	}
	W0, W1 := widths1[0], widths1[1]

	// Row 2 (Pics 2, 3)
	widths2, height2, err2 := calculateRowLayout(ARs[2:4], AW, spacing)
	if err2 != nil {
		return layout, fmt.Errorf("failed to calculate middle row for 2T2M1B: %w", err2)
	}
	W2, W3 := widths2[0], widths2[1]

	// Row 3 (Pic 4)
	W4 := AW
	AR4 := ARs[4]
	height3 := 0.0
	if AR4 > 1e-6 {
		height3 = W4 / AR4
	} else {
		return layout, fmt.Errorf("invalid aspect ratio (%.2f) for bottom picture in 2T2M1B", AR4)
	}
	if height3 <= 1e-6 {
		return layout, fmt.Errorf("calculated zero height (%.2f) for bottom picture in 2T2M1B", height3)
	}

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
	layout.Positions[3] = []float64{W2 + spacing, yRow2}
	layout.Dimensions[3] = []float64{W3, height2}

	yRow3 := yRow2 + height2 + spacing
	layout.Positions[4] = []float64{0, yRow3}
	layout.Dimensions[4] = []float64{W4, height3}

	return layout, nil
}
