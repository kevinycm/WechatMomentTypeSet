package waterfall

import (
	"fmt"
	"math"
)

// calculateLayout_4_3L1R calculates the 3 Left Stacked, 1 Right layout.
// Mirror image of 1L3R.
func calculateLayout_4_3L1R(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 4), Dimensions: make([][]float64, 4)}
	if len(ARs) != 4 || len(types) != 4 {
		return layout, fmt.Errorf("3L1R layout requires 4 ARs and types")
	}
	AR0, AR1, AR2, AR3 := ARs[0], ARs[1], ARs[2], ARs[3] // 0,1,2 are left stack; 3 is right

	// Similar algebraic approach as 1L3R, solving for WL (width of left stack)
	invSumL := 0.0
	if AR0 > 1e-6 {
		invSumL += 1.0 / AR0
	}
	if AR1 > 1e-6 {
		invSumL += 1.0 / AR1
	}
	if AR2 > 1e-6 {
		invSumL += 1.0 / AR2
	}

	invAR3 := 0.0
	if AR3 > 1e-6 {
		invAR3 = 1.0 / AR3
	}

	numerator := AW*invAR3 - spacing*invAR3 - 2*spacing
	denominator := invSumL + invAR3

	WL := 0.0
	if denominator > 1e-6 {
		WL = numerator / denominator
	} else {
		return layout, fmt.Errorf("3L1R cannot solve for WL (denominator zero)")
	}

	if WL <= 1e-6 || WL >= AW-spacing { // Check if WL is valid
		return layout, fmt.Errorf("3L1R geometry infeasible (WL=%.2f)", WL)
	}
	W3 := AW - spacing - WL
	if W3 <= 1e-6 {
		return layout, fmt.Errorf("3L1R geometry infeasible (W3=%.2f)", W3)
	}

	H := 0.0
	if AR3 > 1e-6 {
		H = W3 / AR3
	} else {
		H = WL*invSumL + 2*spacing
	} // Calc H
	if H <= 1e-6 {
		return layout, fmt.Errorf("3L1R calculated zero total height")
	}

	H0, H1, H2 := 0.0, 0.0, 0.0
	if AR0 > 1e-6 {
		H0 = WL / AR0
	} else {
		return layout, fmt.Errorf("3L1R zero height pic 0")
	}
	if AR1 > 1e-6 {
		H1 = WL / AR1
	} else {
		return layout, fmt.Errorf("3L1R zero height pic 1")
	}
	if AR2 > 1e-6 {
		H2 = WL / AR2
	} else {
		return layout, fmt.Errorf("3L1R zero height pic 2")
	}
	if H0 <= 1e-6 || H1 <= 1e-6 || H2 <= 1e-6 {
		return layout, fmt.Errorf("3L1R calculated zero height in left stack")
	}

	// Verify calculated height matches estimate (within tolerance)
	leftStackH := H0 + H1 + H2 + 2*spacing
	if math.Abs(leftStackH-H) > 1e-3 {
		fmt.Printf("Warning: 3L1R height mismatch H=%.2f, H_stack=%.2f. Adjusting.\n", H, leftStackH)
		H = leftStackH // Favor the stack height calculation
		W3 = H * AR3
	}

	layout.TotalHeight = H
	layout.TotalWidth = AW

	currentY := 0.0
	layout.Positions[0] = []float64{0, currentY}
	layout.Dimensions[0] = []float64{WL, H0}
	currentY += H0 + spacing
	layout.Positions[1] = []float64{0, currentY}
	layout.Dimensions[1] = []float64{WL, H1}
	currentY += H1 + spacing
	layout.Positions[2] = []float64{0, currentY}
	layout.Dimensions[2] = []float64{WL, H2}

	layout.Positions[3] = []float64{WL + spacing, 0}
	layout.Dimensions[3] = []float64{W3, H}

	// Minimum height check removed - will be done in the main function after scaling.

	return layout, nil
}
