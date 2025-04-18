package waterfall

import "fmt"

// calculateLayout_4_2x2 calculates the 2x2 grid layout.
// Aims for uniform height within each row.
func calculateLayout_4_2x2(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 4), Dimensions: make([][]float64, 4)}
	if len(ARs) != 4 || len(types) != 4 {
		return layout, fmt.Errorf("2x2 layout requires 4 ARs and types")
	}
	AR0, AR1, AR2, AR3 := ARs[0], ARs[1], ARs[2], ARs[3]

	// Top row calculation
	topRowW := AW - spacing
	topARSum := AR0 + AR1
	H_top := 0.0
	if topRowW > 1e-6 && topARSum > 1e-6 {
		H_top = topRowW / topARSum
	} else {
		return layout, fmt.Errorf("cannot calculate 2x2 top row height (W=%.2f, ARSum=%.2f)", topRowW, topARSum)
	}
	if H_top <= 1e-6 {
		return layout, fmt.Errorf("2x2 calculated zero height for top row")
	}
	W0 := H_top * AR0
	W1 := H_top * AR1

	// Bottom row calculation
	bottomRowW := AW - spacing
	bottomARSum := AR2 + AR3
	H_bottom := 0.0
	if bottomRowW > 1e-6 && bottomARSum > 1e-6 {
		H_bottom = bottomRowW / bottomARSum
	} else {
		return layout, fmt.Errorf("cannot calculate 2x2 bottom row height (W=%.2f, ARSum=%.2f)", bottomRowW, bottomARSum)
	}
	if H_bottom <= 1e-6 {
		return layout, fmt.Errorf("2x2 calculated zero height for bottom row")
	}
	W2 := H_bottom * AR2
	W3 := H_bottom * AR3

	layout.TotalHeight = H_top + spacing + H_bottom
	layout.TotalWidth = AW

	layout.Positions[0] = []float64{0, 0}
	layout.Dimensions[0] = []float64{W0, H_top}
	layout.Positions[1] = []float64{W0 + spacing, 0}
	layout.Dimensions[1] = []float64{W1, H_top}
	layout.Positions[2] = []float64{0, H_top + spacing}
	layout.Dimensions[2] = []float64{W2, H_bottom}
	layout.Positions[3] = []float64{W2 + spacing, H_top + spacing}
	layout.Dimensions[3] = []float64{W3, H_bottom}

	// Minimum height check removed - will be done in the main function after scaling.

	return layout, nil
}
