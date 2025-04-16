package calculate

import "fmt"

// calculateLayout_5_1T2M2B calculates the 1 Top, 2 Middle, 2 Bottom layout.
func calculateLayout_5_1T2M2B(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 5), Dimensions: make([][]float64, 5)}
	if len(ARs) != 5 || len(types) != 5 {
		return layout, fmt.Errorf("1T2M2B layout requires 5 ARs and types")
	}

	// Row 1 (Pic 0)
	W0 := AW
	AR0 := ARs[0]
	height1 := 0.0
	if AR0 > 1e-6 {
		height1 = W0 / AR0
	} else {
		return layout, fmt.Errorf("invalid aspect ratio (%.2f) for top picture in 1T2M2B", AR0)
	}
	if height1 <= 1e-6 {
		return layout, fmt.Errorf("calculated zero height (%.2f) for top picture in 1T2M2B", height1)
	}

	// Row 2 (Pics 1, 2)
	widths2, height2, err2 := calculateRowLayout(ARs[1:3], AW, spacing)
	if err2 != nil {
		return layout, fmt.Errorf("failed to calculate middle row for 1T2M2B: %w", err2)
	}
	W1, W2 := widths2[0], widths2[1]

	// Row 3 (Pics 3, 4)
	widths3, height3, err3 := calculateRowLayout(ARs[3:5], AW, spacing)
	if err3 != nil {
		return layout, fmt.Errorf("failed to calculate bottom row for 1T2M2B: %w", err3)
	}
	W3, W4 := widths3[0], widths3[1]

	layout.TotalHeight = height1 + spacing + height2 + spacing + height3
	layout.TotalWidth = AW

	// Positions and Dimensions
	layout.Positions[0] = []float64{0, 0}
	layout.Dimensions[0] = []float64{W0, height1}

	yRow2 := height1 + spacing
	layout.Positions[1] = []float64{0, yRow2}
	layout.Dimensions[1] = []float64{W1, height2}
	layout.Positions[2] = []float64{W1 + spacing, yRow2}
	layout.Dimensions[2] = []float64{W2, height2}

	yRow3 := yRow2 + height2 + spacing
	layout.Positions[3] = []float64{0, yRow3}
	layout.Dimensions[3] = []float64{W3, height3}
	layout.Positions[4] = []float64{W3 + spacing, yRow3}
	layout.Dimensions[4] = []float64{W4, height3}

	return layout, nil
}
