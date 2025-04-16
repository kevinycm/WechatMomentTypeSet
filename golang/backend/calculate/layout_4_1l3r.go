package calculate

import (
	"fmt"
	"math"
)

// calculateLayout_4_1L3R calculates the 1 Left, 3 Right Stacked layout.
// This is complex geometrically. Let's try a simplification:
// Assume left pic (0) and right stack (1, 2, 3) share the same total height.
// Calculate height H1, H2, H3 for right stack assuming they fit in some width WR.
// Total right height HR = H1+H2+H3+2*spacing.
// Calculate required width W0 for left pic to have height HR: W0 = HR * AR0.
// Check if W0 + spacing + WR = AW.
// This often requires iterative solving or algebraic manipulation.
// Simpler approach: Treat right side as a single column, calculate its total AR.
func calculateLayout_4_1L3R(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 4), Dimensions: make([][]float64, 4)}
	if len(ARs) != 4 || len(types) != 4 {
		return layout, fmt.Errorf("1L3R layout requires 4 ARs and types")
	}
	AR0, AR1, AR2, AR3 := ARs[0], ARs[1], ARs[2], ARs[3]

	// Estimate right column effective AR (complex, simplifying)
	// Let's fix the height H and calculate widths.
	// We need to solve for H such that W0 + spacing + WR = AW
	// W0 = H * AR0
	// WR needs H1, H2, H3. Assume WR is fixed for a moment.
	// H1 = WR / AR1, H2 = WR / AR2, H3 = WR / AR3
	// H = H1 + H2 + H3 + 2*spacing
	// Substituting: H = WR/AR1 + WR/AR2 + WR/AR3 + 2*spacing
	// H = WR * (1/AR1 + 1/AR2 + 1/AR3) + 2*spacing
	// W0 = (WR * (1/AR1 + 1/AR2 + 1/AR3) + 2*spacing) * AR0
	// (WR * (1/AR1 + 1/AR2 + 1/AR3) + 2*spacing) * AR0 + spacing + WR = AW
	// This gets complicated quickly. Let's use a simpler geometric approach similar to 1L2R.

	// Assume total height H is determined by the left image (pic 0).
	// Let W0 be width of left, WR be width of right column.
	// W0 + spacing + WR = AW.
	// H = W0 / AR0 (if AR0 != 0)
	// Right stack height must also be H.
	// H = H1 + H2 + H3 + 2*spacing
	// H = WR/AR1 + WR/AR2 + WR/AR3 + 2*spacing = WR * (1/AR1 + 1/AR2 + 1/AR3) + 2*spacing
	// W0/AR0 = WR * (1/AR1 + 1/AR2 + 1/AR3) + 2*spacing
	// (AW - WR - spacing) / AR0 = WR*invSumR + 2*spacing where invSumR = 1/AR1+1/AR2+1/AR3
	// AW/AR0 - WR/AR0 - spacing/AR0 = WR*invSumR + 2*spacing
	// AW/AR0 - spacing/AR0 - 2*spacing = WR*invSumR + WR/AR0 = WR * (invSumR + 1/AR0)
	invSumR := 0.0
	if AR1 > 1e-6 {
		invSumR += 1.0 / AR1
	}
	if AR2 > 1e-6 {
		invSumR += 1.0 / AR2
	}
	if AR3 > 1e-6 {
		invSumR += 1.0 / AR3
	}

	invAR0 := 0.0
	if AR0 > 1e-6 {
		invAR0 = 1.0 / AR0
	}

	numerator := AW*invAR0 - spacing*invAR0 - 2*spacing
	denominator := invSumR + invAR0

	WR := 0.0
	if denominator > 1e-6 {
		WR = numerator / denominator
	} else {
		return layout, fmt.Errorf("1L3R cannot solve for WR (denominator zero)")
	}

	if WR <= 1e-6 || WR >= AW-spacing { // Check if WR is valid
		return layout, fmt.Errorf("1L3R geometry infeasible (WR=%.2f)", WR)
	}
	W0 := AW - spacing - WR
	if W0 <= 1e-6 {
		return layout, fmt.Errorf("1L3R geometry infeasible (W0=%.2f)", W0)
	}

	H := 0.0
	if AR0 > 1e-6 {
		H = W0 / AR0
	} else {
		H = WR*invSumR + 2*spacing
	} // Calc H
	if H <= 1e-6 {
		return layout, fmt.Errorf("1L3R calculated zero total height")
	}

	H1, H2, H3 := 0.0, 0.0, 0.0
	if AR1 > 1e-6 {
		H1 = WR / AR1
	} else {
		return layout, fmt.Errorf("1L3R zero height pic 1")
	}
	if AR2 > 1e-6 {
		H2 = WR / AR2
	} else {
		return layout, fmt.Errorf("1L3R zero height pic 2")
	}
	if AR3 > 1e-6 {
		H3 = WR / AR3
	} else {
		return layout, fmt.Errorf("1L3R zero height pic 3")
	}
	if H1 <= 1e-6 || H2 <= 1e-6 || H3 <= 1e-6 {
		return layout, fmt.Errorf("1L3R calculated zero height in right stack")
	}

	// Verify calculated height matches estimate (within tolerance)
	rightStackH := H1 + H2 + H3 + 2*spacing
	if math.Abs(rightStackH-H) > 1e-3 {
		fmt.Printf("Warning: 1L3R height mismatch H=%.2f, H_stack=%.2f. Adjusting.\n", H, rightStackH)
		// Could potentially adjust spacing or rescale, but let's proceed with calculated H for now.
		H = rightStackH // Favor the stack height calculation? Or average?
		W0 = H * AR0
	}

	layout.TotalHeight = H
	layout.TotalWidth = AW

	layout.Positions[0] = []float64{0, 0}
	layout.Dimensions[0] = []float64{W0, H}

	currentY := 0.0
	rightX := W0 + spacing
	layout.Positions[1] = []float64{rightX, currentY}
	layout.Dimensions[1] = []float64{WR, H1}
	currentY += H1 + spacing
	layout.Positions[2] = []float64{rightX, currentY}
	layout.Dimensions[2] = []float64{WR, H2}
	currentY += H2 + spacing
	layout.Positions[3] = []float64{rightX, currentY}
	layout.Dimensions[3] = []float64{WR, H3}

	// Minimum height check removed - will be done in the main function after scaling.

	return layout, nil
}
