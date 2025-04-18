package waterfall

import "fmt"

// calculateLayout_4_3T1B calculates the 3 Top, 1 Bottom layout.
func calculateLayout_4_3T1B(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 4), Dimensions: make([][]float64, 4)}
	if len(ARs) != 4 || len(types) != 4 {
		return layout, fmt.Errorf("3T1B layout requires 4 ARs and types")
	}
	AR0, AR1, AR2, AR3 := ARs[0], ARs[1], ARs[2], ARs[3]

	// Top row (0, 1, 2)
	topRowW := AW - 2*spacing
	topARSum := AR0 + AR1 + AR2
	H_top := 0.0
	if topRowW > 1e-6 && topARSum > 1e-6 {
		H_top = topRowW / topARSum
	} else {
		return layout, fmt.Errorf("cannot calculate 3T1B top row height (W=%.2f, ARSum=%.2f)", topRowW, topARSum)
	}
	if H_top <= 1e-6 {
		return layout, fmt.Errorf("3T1B calculated zero height for top row")
	}
	W0 := H_top * AR0
	W1 := H_top * AR1
	W2 := H_top * AR2

	// Bottom pic (3)
	W3 := AW
	H3 := 0.0
	if AR3 > 1e-6 {
		H3 = W3 / AR3
	} else {
		return layout, fmt.Errorf("3T1B calculated zero height for bottom pic")
	}
	if H3 <= 1e-6 {
		return layout, fmt.Errorf("3T1B calculated zero height for bottom pic")
	}

	layout.TotalHeight = H_top + spacing + H3
	layout.TotalWidth = AW

	currentX := 0.0
	layout.Positions[0] = []float64{currentX, 0}
	layout.Dimensions[0] = []float64{W0, H_top}
	currentX += W0 + spacing
	layout.Positions[1] = []float64{currentX, 0}
	layout.Dimensions[1] = []float64{W1, H_top}
	currentX += W1 + spacing
	layout.Positions[2] = []float64{currentX, 0}
	layout.Dimensions[2] = []float64{W2, H_top}

	layout.Positions[3] = []float64{0, H_top + spacing}
	layout.Dimensions[3] = []float64{W3, H3}

	// Minimum height check removed - will be done in the main function after scaling.

	return layout, nil
}
