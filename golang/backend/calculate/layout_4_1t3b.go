package calculate

import "fmt"

// calculateLayout_4_1T3B calculates the 1 Top, 3 Bottom layout.
func calculateLayout_4_1T3B(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 4), Dimensions: make([][]float64, 4)}
	if len(ARs) != 4 || len(types) != 4 {
		return layout, fmt.Errorf("1T3B layout requires 4 ARs and types")
	}
	AR0, AR1, AR2, AR3 := ARs[0], ARs[1], ARs[2], ARs[3]

	// Top pic (0)
	W0 := AW
	H0 := 0.0
	if AR0 > 1e-6 {
		H0 = W0 / AR0
	} else {
		return layout, fmt.Errorf("1T3B calculated zero height for top pic")
	}
	if H0 <= 1e-6 {
		return layout, fmt.Errorf("1T3B calculated zero height for top pic")
	}

	// Bottom row (1, 2, 3)
	bottomRowW := AW - 2*spacing
	bottomARSum := AR1 + AR2 + AR3
	H_bottom := 0.0
	if bottomRowW > 1e-6 && bottomARSum > 1e-6 {
		H_bottom = bottomRowW / bottomARSum
	} else {
		return layout, fmt.Errorf("cannot calculate 1T3B bottom row height (W=%.2f, ARSum=%.2f)", bottomRowW, bottomARSum)
	}
	if H_bottom <= 1e-6 {
		return layout, fmt.Errorf("1T3B calculated zero height for bottom row")
	}
	W1 := H_bottom * AR1
	W2 := H_bottom * AR2
	W3 := H_bottom * AR3

	layout.TotalHeight = H0 + spacing + H_bottom
	layout.TotalWidth = AW

	layout.Positions[0] = []float64{0, 0}
	layout.Dimensions[0] = []float64{W0, H0}

	currentX := 0.0
	bottomY := H0 + spacing
	layout.Positions[1] = []float64{currentX, bottomY}
	layout.Dimensions[1] = []float64{W1, H_bottom}
	currentX += W1 + spacing
	layout.Positions[2] = []float64{currentX, bottomY}
	layout.Dimensions[2] = []float64{W2, H_bottom}
	currentX += W2 + spacing
	layout.Positions[3] = []float64{currentX, bottomY}
	layout.Dimensions[3] = []float64{W3, H_bottom}

	// Minimum height check removed - will be done in the main function after scaling.

	return layout, nil
}
