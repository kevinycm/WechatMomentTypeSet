package waterfall

import "fmt"

// calculateLayout_7_4L3R calculates the 4 Left, 3 Right Stacked layout.
func calculateLayout_7_4L3R(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 7), Dimensions: make([][]float64, 7)}
	if len(ARs) != 7 || len(types) != 7 {
		return layout, fmt.Errorf("4L3R layout requires 7 ARs and types")
	}
	// Pics 0,1,2,3 left; 4,5,6 right

	invSumL, invSumR := 0.0, 0.0
	numLeft, numRight := 4, 3
	numLeftSpacing, numRightSpacing := float64(numLeft-1), float64(numRight-1)
	for i := 0; i < numLeft; i++ {
		if ARs[i] > 1e-6 {
			invSumL += 1.0 / ARs[i]
		} else {
			return layout, fmt.Errorf("invalid AR in left stack for 4L3R")
		}
	}
	for i := numLeft; i < numLeft+numRight; i++ {
		if ARs[i] > 1e-6 {
			invSumR += 1.0 / ARs[i]
		} else {
			return layout, fmt.Errorf("invalid AR in right stack for 4L3R")
		}
	}

	// Same solving logic as 3L4R, just different numLeft/numRight
	denominator := invSumL + invSumR
	WL := 0.0
	if denominator > 1e-6 {
		numerator := AW*invSumR - spacing*invSumR + (numRightSpacing-numLeftSpacing)*spacing
		WL = numerator / denominator
	} else {
		return layout, fmt.Errorf("4L3R cannot solve for WL (denominator zero)")
	}

	if WL <= 1e-6 || WL >= AW-spacing {
		return layout, fmt.Errorf("4L3R geometry infeasible (WL=%.2f)", WL)
	}
	WR := AW - spacing - WL
	if WR <= 1e-6 {
		return layout, fmt.Errorf("4L3R geometry infeasible (WR=%.2f)", WR)
	}

	H := 0.0 // Total height
	if invSumL > 1e-6 {
		H = WL*invSumL + numLeftSpacing*spacing // +3 spacing for 4 items
	} else if invSumR > 1e-6 {
		H = WR*invSumR + numRightSpacing*spacing // +2 spacing for 3 items
	} else {
		return layout, fmt.Errorf("4L3R cannot determine height")
	}

	if H <= 1e-6 {
		return layout, fmt.Errorf("4L3R calculated zero total height")
	}

	layout.TotalHeight = H
	layout.TotalWidth = AW

	// Calculate individual heights
	Hs := make([]float64, 7)
	for i := 0; i < numLeft; i++ {
		Hs[i] = WL / ARs[i]
	}
	for i := numLeft; i < numLeft+numRight; i++ {
		Hs[i] = WR / ARs[i]
	}

	// Set positions and dimensions
	currentY := 0.0
	for i := 0; i < numLeft; i++ {
		layout.Positions[i] = []float64{0, currentY}
		layout.Dimensions[i] = []float64{WL, Hs[i]}
		if i < numLeft-1 {
			currentY += Hs[i] + spacing
		}
	}
	currentY = 0.0
	rightX := WL + spacing
	for i := numLeft; i < numLeft+numRight; i++ {
		layout.Positions[i] = []float64{rightX, currentY}
		layout.Dimensions[i] = []float64{WR, Hs[i]}
		if i < numLeft+numRight-1 {
			currentY += Hs[i] + spacing
		}
	}

	return layout, nil
}
