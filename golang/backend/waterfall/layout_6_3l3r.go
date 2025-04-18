package waterfall

import (
	"fmt"
	"math"
)

// calculateLayout_6_3L3R calculates the 3 Left, 3 Right Stacked layout.
func calculateLayout_6_3L3R(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 6), Dimensions: make([][]float64, 6)}
	if len(ARs) != 6 || len(types) != 6 {
		return layout, fmt.Errorf("3L3R layout requires 6 ARs and types")
	}
	AR0, AR1, AR2 := ARs[0], ARs[1], ARs[2] // Left stack
	AR3, AR4, AR5 := ARs[3], ARs[4], ARs[5] // Right stack

	invSumL, invSumR := 0.0, 0.0
	for i := 0; i < 3; i++ {
		if ARs[i] > 1e-6 {
			invSumL += 1.0 / ARs[i]
		} else {
			return layout, fmt.Errorf("invalid AR in left stack for 3L3R")
		}
		if ARs[i+3] > 1e-6 {
			invSumR += 1.0 / ARs[i+3]
		} else {
			return layout, fmt.Errorf("invalid AR in right stack for 3L3R")
		}
	}

	denominator := 1.0
	if invSumR > 1e-6 {
		denominator += invSumL / invSumR
	} else if invSumL > 1e-6 {
		// If right side has infinite AR sum (zero invSumR), right width WR -> 0?
		// This case is tricky, implies right side is infinitely tall. Let's assume error.
		return layout, fmt.Errorf("right stack AR sum invalid for 3L3R")
	}

	WL := 0.0
	if denominator > 1e-6 {
		WL = (AW - spacing) / denominator
	} else {
		return layout, fmt.Errorf("3L3R cannot solve for WL (denominator zero)")
	}

	if WL <= 1e-6 || WL >= AW-spacing {
		return layout, fmt.Errorf("3L3R geometry infeasible (WL=%.2f)", WL)
	}
	WR := AW - spacing - WL
	if WR <= 1e-6 {
		return layout, fmt.Errorf("3L3R geometry infeasible (WR=%.2f)", WR)
	}

	H := 0.0
	if invSumL > 1e-6 {
		H = WL*invSumL + 2*spacing
	} else if invSumR > 1e-6 {
		H = WR*invSumR + 2*spacing
	} else {
		return layout, fmt.Errorf("3L3R cannot determine height (both stacks have zero AR sum)")
	}

	if H <= 1e-6 {
		return layout, fmt.Errorf("3L3R calculated zero total height")
	}

	H0 := WL / AR0
	H1 := WL / AR1
	H2 := WL / AR2
	H3 := WR / AR3
	H4 := WR / AR4
	H5 := WR / AR5

	// Verify heights add up (optional check, might cause issues with float precision)
	leftStackH := H0 + H1 + H2 + 2*spacing
	rightStackH := H3 + H4 + H5 + 2*spacing
	if math.Abs(leftStackH-H) > 1e-3 || math.Abs(rightStackH-H) > 1e-3 {
		fmt.Printf("Warning: 3L3R height mismatch H=%.2f, HL=%.2f, HR=%.2f. Using calculated H.\n", H, leftStackH, rightStackH)
		// Don't adjust H here, stick with the geometrically derived total H
	}

	layout.TotalHeight = H
	layout.TotalWidth = AW

	// Left Column
	currentY := 0.0
	layout.Positions[0] = []float64{0, currentY}
	layout.Dimensions[0] = []float64{WL, H0}
	currentY += H0 + spacing
	layout.Positions[1] = []float64{0, currentY}
	layout.Dimensions[1] = []float64{WL, H1}
	currentY += H1 + spacing
	layout.Positions[2] = []float64{0, currentY}
	layout.Dimensions[2] = []float64{WL, H2}

	// Right Column
	currentY = 0.0
	rightX := WL + spacing
	layout.Positions[3] = []float64{rightX, currentY}
	layout.Dimensions[3] = []float64{WR, H3}
	currentY += H3 + spacing
	layout.Positions[4] = []float64{rightX, currentY}
	layout.Dimensions[4] = []float64{WR, H4}
	currentY += H4 + spacing
	layout.Positions[5] = []float64{rightX, currentY}
	layout.Dimensions[5] = []float64{WR, H5}

	return layout, nil
}
