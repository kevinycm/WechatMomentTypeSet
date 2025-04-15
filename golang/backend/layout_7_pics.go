package backend

import (
	"fmt"
	"math"
)

// --- Calculation Helper Functions (Local to layout_7_pics.go) ---

// calculateLayout_7_3T2M2B calculates the 3 Top, 2 Middle, 2 Bottom layout.
func calculateLayout_7_3T2M2B(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 7), Dimensions: make([][]float64, 7)}
	if len(ARs) != 7 || len(types) != 7 {
		return layout, fmt.Errorf("3T2M2B layout requires 7 ARs and types")
	}

	// Row 1 (0, 1, 2)
	widths1, height1, err1 := calculateRowLayout(ARs[0:3], AW, spacing)
	if err1 != nil {
		return layout, fmt.Errorf("failed top row for 3T2M2B: %w", err1)
	}
	W0, W1, W2 := widths1[0], widths1[1], widths1[2]

	// Row 2 (3, 4)
	widths2, height2, err2 := calculateRowLayout(ARs[3:5], AW, spacing)
	if err2 != nil {
		return layout, fmt.Errorf("failed middle row for 3T2M2B: %w", err2)
	}
	W3, W4 := widths2[0], widths2[1]

	// Row 3 (5, 6)
	widths3, height3, err3 := calculateRowLayout(ARs[5:7], AW, spacing)
	if err3 != nil {
		return layout, fmt.Errorf("failed bottom row for 3T2M2B: %w", err3)
	}
	W5, W6 := widths3[0], widths3[1]

	layout.TotalHeight = height1 + spacing + height2 + spacing + height3
	layout.TotalWidth = AW

	currentX := 0.0
	layout.Positions[0] = []float64{currentX, 0}
	layout.Dimensions[0] = []float64{W0, height1}
	currentX += W0 + spacing
	layout.Positions[1] = []float64{currentX, 0}
	layout.Dimensions[1] = []float64{W1, height1}
	currentX += W1 + spacing
	layout.Positions[2] = []float64{currentX, 0}
	layout.Dimensions[2] = []float64{W2, height1}
	yRow2 := height1 + spacing
	layout.Positions[3] = []float64{0, yRow2}
	layout.Dimensions[3] = []float64{W3, height2}
	layout.Positions[4] = []float64{W3 + spacing, yRow2}
	layout.Dimensions[4] = []float64{W4, height2}
	yRow3 := yRow2 + height2 + spacing
	layout.Positions[5] = []float64{0, yRow3}
	layout.Dimensions[5] = []float64{W5, height3}
	layout.Positions[6] = []float64{W5 + spacing, yRow3}
	layout.Dimensions[6] = []float64{W6, height3}

	return layout, nil
}

// calculateLayout_7_2T3M2B calculates the 2 Top, 3 Middle, 2 Bottom layout.
func calculateLayout_7_2T3M2B(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 7), Dimensions: make([][]float64, 7)}
	if len(ARs) != 7 || len(types) != 7 {
		return layout, fmt.Errorf("2T3M2B layout requires 7 ARs and types")
	}

	// Row 1 (0, 1)
	widths1, height1, err1 := calculateRowLayout(ARs[0:2], AW, spacing)
	if err1 != nil {
		return layout, fmt.Errorf("failed top row for 2T3M2B: %w", err1)
	}
	W0, W1 := widths1[0], widths1[1]

	// Row 2 (2, 3, 4)
	widths2, height2, err2 := calculateRowLayout(ARs[2:5], AW, spacing)
	if err2 != nil {
		return layout, fmt.Errorf("failed middle row for 2T3M2B: %w", err2)
	}
	W2, W3, W4 := widths2[0], widths2[1], widths2[2]

	// Row 3 (5, 6)
	widths3, height3, err3 := calculateRowLayout(ARs[5:7], AW, spacing)
	if err3 != nil {
		return layout, fmt.Errorf("failed bottom row for 2T3M2B: %w", err3)
	}
	W5, W6 := widths3[0], widths3[1]

	layout.TotalHeight = height1 + spacing + height2 + spacing + height3
	layout.TotalWidth = AW

	layout.Positions[0] = []float64{0, 0}
	layout.Dimensions[0] = []float64{W0, height1}
	layout.Positions[1] = []float64{W0 + spacing, 0}
	layout.Dimensions[1] = []float64{W1, height1}
	yRow2 := height1 + spacing
	currentX := 0.0
	layout.Positions[2] = []float64{currentX, yRow2}
	layout.Dimensions[2] = []float64{W2, height2}
	currentX += W2 + spacing
	layout.Positions[3] = []float64{currentX, yRow2}
	layout.Dimensions[3] = []float64{W3, height2}
	currentX += W3 + spacing
	layout.Positions[4] = []float64{currentX, yRow2}
	layout.Dimensions[4] = []float64{W4, height2}
	yRow3 := yRow2 + height2 + spacing
	layout.Positions[5] = []float64{0, yRow3}
	layout.Dimensions[5] = []float64{W5, height3}
	layout.Positions[6] = []float64{W5 + spacing, yRow3}
	layout.Dimensions[6] = []float64{W6, height3}

	return layout, nil
}

// calculateLayout_7_2T2M3B calculates the 2 Top, 2 Middle, 3 Bottom layout.
func calculateLayout_7_2T2M3B(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 7), Dimensions: make([][]float64, 7)}
	if len(ARs) != 7 || len(types) != 7 {
		return layout, fmt.Errorf("2T2M3B layout requires 7 ARs and types")
	}

	// Row 1 (0, 1)
	widths1, height1, err1 := calculateRowLayout(ARs[0:2], AW, spacing)
	if err1 != nil {
		return layout, fmt.Errorf("failed top row for 2T2M3B: %w", err1)
	}
	W0, W1 := widths1[0], widths1[1]

	// Row 2 (2, 3)
	widths2, height2, err2 := calculateRowLayout(ARs[2:4], AW, spacing)
	if err2 != nil {
		return layout, fmt.Errorf("failed middle row for 2T2M3B: %w", err2)
	}
	W2, W3 := widths2[0], widths2[1]

	// Row 3 (4, 5, 6)
	widths3, height3, err3 := calculateRowLayout(ARs[4:7], AW, spacing)
	if err3 != nil {
		return layout, fmt.Errorf("failed bottom row for 2T2M3B: %w", err3)
	}
	W4, W5, W6 := widths3[0], widths3[1], widths3[2]

	layout.TotalHeight = height1 + spacing + height2 + spacing + height3
	layout.TotalWidth = AW

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
	currentX := 0.0
	layout.Positions[4] = []float64{currentX, yRow3}
	layout.Dimensions[4] = []float64{W4, height3}
	currentX += W4 + spacing
	layout.Positions[5] = []float64{currentX, yRow3}
	layout.Dimensions[5] = []float64{W5, height3}
	currentX += W5 + spacing
	layout.Positions[6] = []float64{currentX, yRow3}
	layout.Dimensions[6] = []float64{W6, height3}

	return layout, nil
}

// calculateLayout_7_3T3M1B calculates the 3 Top, 3 Middle, 1 Bottom layout.
func calculateLayout_7_3T3M1B(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 7), Dimensions: make([][]float64, 7)}
	if len(ARs) != 7 || len(types) != 7 {
		return layout, fmt.Errorf("3T3M1B layout requires 7 ARs and types")
	}

	// Row 1 (0, 1, 2)
	widths1, height1, err1 := calculateRowLayout(ARs[0:3], AW, spacing)
	if err1 != nil {
		return layout, fmt.Errorf("failed top row for 3T3M1B: %w", err1)
	}
	W0, W1, W2 := widths1[0], widths1[1], widths1[2]

	// Row 2 (3, 4, 5)
	widths2, height2, err2 := calculateRowLayout(ARs[3:6], AW, spacing)
	if err2 != nil {
		return layout, fmt.Errorf("failed middle row for 3T3M1B: %w", err2)
	}
	W3, W4, W5 := widths2[0], widths2[1], widths2[2]

	// Row 3 (6)
	W6 := AW
	AR6 := ARs[6]
	height3 := 0.0
	if AR6 > 1e-6 {
		height3 = W6 / AR6
	} else {
		return layout, fmt.Errorf("invalid AR for bottom pic in 3T3M1B")
	}
	if height3 <= 1e-6 {
		return layout, fmt.Errorf("zero height for bottom pic in 3T3M1B")
	}

	layout.TotalHeight = height1 + spacing + height2 + spacing + height3
	layout.TotalWidth = AW

	currentX := 0.0
	layout.Positions[0] = []float64{currentX, 0}
	layout.Dimensions[0] = []float64{W0, height1}
	currentX += W0 + spacing
	layout.Positions[1] = []float64{currentX, 0}
	layout.Dimensions[1] = []float64{W1, height1}
	currentX += W1 + spacing
	layout.Positions[2] = []float64{currentX, 0}
	layout.Dimensions[2] = []float64{W2, height1}
	yRow2 := height1 + spacing
	currentX = 0.0
	layout.Positions[3] = []float64{currentX, yRow2}
	layout.Dimensions[3] = []float64{W3, height2}
	currentX += W3 + spacing
	layout.Positions[4] = []float64{currentX, yRow2}
	layout.Dimensions[4] = []float64{W4, height2}
	currentX += W4 + spacing
	layout.Positions[5] = []float64{currentX, yRow2}
	layout.Dimensions[5] = []float64{W5, height2}
	yRow3 := yRow2 + height2 + spacing
	layout.Positions[6] = []float64{0, yRow3}
	layout.Dimensions[6] = []float64{W6, height3}

	return layout, nil
}

// calculateLayout_7_3L4R calculates the 3 Left, 4 Right Stacked layout.
func calculateLayout_7_3L4R(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 7), Dimensions: make([][]float64, 7)}
	if len(ARs) != 7 || len(types) != 7 {
		return layout, fmt.Errorf("3L4R layout requires 7 ARs and types")
	}
	// Pics 0,1,2 left; 3,4,5,6 right

	invSumL, invSumR := 0.0, 0.0
	numLeft, numRight := 3, 4
	numLeftSpacing, numRightSpacing := float64(numLeft-1), float64(numRight-1)
	for i := 0; i < numLeft; i++ {
		if ARs[i] > 1e-6 {
			invSumL += 1.0 / ARs[i]
		} else {
			return layout, fmt.Errorf("invalid AR in left stack for 3L4R")
		}
	}
	for i := numLeft; i < numLeft+numRight; i++ {
		if ARs[i] > 1e-6 {
			invSumR += 1.0 / ARs[i]
		} else {
			return layout, fmt.Errorf("invalid AR in right stack for 3L4R")
		}
	}

	// Solve for WL assuming total height H is equal for both stacks
	// H_L = WL * invSumL + numLeftSpacing * spacing
	// H_R = WR * invSumR + numRightSpacing * spacing
	// H_L = H_R => WL * invSumL + numLeftSpacing * spacing = WR * invSumR + numRightSpacing * spacing
	// WR = AW - spacing - WL
	// WL*invSumL + numLeftSpacing*spacing = (AW-spacing-WL)*invSumR + numRightSpacing*spacing
	// WL*invSumL + numLeftSpacing*spacing = AW*invSumR - spacing*invSumR - WL*invSumR + numRightSpacing*spacing
	// WL*invSumL + WL*invSumR = AW*invSumR - spacing*invSumR + numRightSpacing*spacing - numLeftSpacing*spacing
	// WL * (invSumL + invSumR) = AW*invSumR - spacing*invSumR + (numRightSpacing - numLeftSpacing)*spacing
	// WL = (AW*invSumR - spacing*invSumR + (numRightSpacing - numLeftSpacing)*spacing) / (invSumL + invSumR)

	denominator := invSumL + invSumR
	WL := 0.0
	if denominator > 1e-6 {
		numerator := AW*invSumR - spacing*invSumR + (numRightSpacing-numLeftSpacing)*spacing
		WL = numerator / denominator
	} else {
		return layout, fmt.Errorf("3L4R cannot solve for WL (denominator zero)")
	}

	if WL <= 1e-6 || WL >= AW-spacing {
		return layout, fmt.Errorf("3L4R geometry infeasible (WL=%.2f)", WL)
	}
	WR := AW - spacing - WL
	if WR <= 1e-6 {
		return layout, fmt.Errorf("3L4R geometry infeasible (WR=%.2f)", WR)
	}

	H := 0.0 // Total height
	if invSumL > 1e-6 {
		H = WL*invSumL + numLeftSpacing*spacing // +2 spacing for 3 items
	} else if invSumR > 1e-6 {
		H = WR*invSumR + numRightSpacing*spacing // +3 spacing for 4 items
	} else {
		return layout, fmt.Errorf("3L4R cannot determine height")
	}

	if H <= 1e-6 {
		return layout, fmt.Errorf("3L4R calculated zero total height")
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

// --- NEW: Four-Row Layout Functions ---

// calculateLayout_7_2T2M2M1B calculates the 2T-2M-2M-1B layout.
func calculateLayout_7_2T2M2M1B(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 7), Dimensions: make([][]float64, 7)}
	if len(ARs) != 7 || len(types) != 7 {
		return layout, fmt.Errorf("2T2M2M1B layout requires 7 ARs and types")
	}

	// Row 1 (0, 1)
	widths1, height1, err1 := calculateRowLayout(ARs[0:2], AW, spacing)
	if err1 != nil {
		return layout, fmt.Errorf("failed row 1 for 2T2M2M1B: %w", err1)
	}
	W0, W1 := widths1[0], widths1[1]

	// Row 2 (2, 3)
	widths2, height2, err2 := calculateRowLayout(ARs[2:4], AW, spacing)
	if err2 != nil {
		return layout, fmt.Errorf("failed row 2 for 2T2M2M1B: %w", err2)
	}
	W2, W3 := widths2[0], widths2[1]

	// Row 3 (4, 5)
	widths3, height3, err3 := calculateRowLayout(ARs[4:6], AW, spacing)
	if err3 != nil {
		return layout, fmt.Errorf("failed row 3 for 2T2M2M1B: %w", err3)
	}
	W4, W5 := widths3[0], widths3[1]

	// Row 4 (6)
	W6 := AW
	AR6 := ARs[6]
	height4 := 0.0
	if AR6 > 1e-6 {
		height4 = W6 / AR6
	} else {
		return layout, fmt.Errorf("invalid AR for pic 6 in 2T2M2M1B")
	}
	if height4 <= 1e-6 {
		return layout, fmt.Errorf("zero height for pic 6 in 2T2M2M1B")
	}

	layout.TotalHeight = height1 + spacing + height2 + spacing + height3 + spacing + height4
	layout.TotalWidth = AW

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
	layout.Positions[5] = []float64{W4 + spacing, yRow3}
	layout.Dimensions[5] = []float64{W5, height3}
	yRow4 := yRow3 + height3 + spacing
	layout.Positions[6] = []float64{0, yRow4}
	layout.Dimensions[6] = []float64{W6, height4}

	return layout, nil
}

// calculateLayout_7_1T2M2M2B calculates the 1T-2M-2M-2B layout.
func calculateLayout_7_1T2M2M2B(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 7), Dimensions: make([][]float64, 7)}
	if len(ARs) != 7 || len(types) != 7 {
		return layout, fmt.Errorf("1T2M2M2B layout requires 7 ARs and types")
	}

	// Row 1 (0)
	W0 := AW
	AR0 := ARs[0]
	height1 := 0.0
	if AR0 > 1e-6 {
		height1 = W0 / AR0
	} else {
		return layout, fmt.Errorf("invalid AR for pic 0 in 1T2M2M2B")
	}
	if height1 <= 1e-6 {
		return layout, fmt.Errorf("zero height for pic 0 in 1T2M2M2B")
	}

	// Row 2 (1, 2)
	widths2, height2, err2 := calculateRowLayout(ARs[1:3], AW, spacing)
	if err2 != nil {
		return layout, fmt.Errorf("failed row 2 for 1T2M2M2B: %w", err2)
	}
	W1, W2 := widths2[0], widths2[1]

	// Row 3 (3, 4)
	widths3, height3, err3 := calculateRowLayout(ARs[3:5], AW, spacing)
	if err3 != nil {
		return layout, fmt.Errorf("failed row 3 for 1T2M2M2B: %w", err3)
	}
	W3, W4 := widths3[0], widths3[1]

	// Row 4 (5, 6)
	widths4, height4, err4 := calculateRowLayout(ARs[5:7], AW, spacing)
	if err4 != nil {
		return layout, fmt.Errorf("failed row 4 for 1T2M2M2B: %w", err4)
	}
	W5, W6 := widths4[0], widths4[1]

	layout.TotalHeight = height1 + spacing + height2 + spacing + height3 + spacing + height4
	layout.TotalWidth = AW

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
	yRow4 := yRow3 + height3 + spacing
	layout.Positions[5] = []float64{0, yRow4}
	layout.Dimensions[5] = []float64{W5, height4}
	layout.Positions[6] = []float64{W5 + spacing, yRow4}
	layout.Dimensions[6] = []float64{W6, height4}

	return layout, nil
}

// calculateLayout_7_3T2M1M1B calculates the 3T-2M-1M-1B layout.
func calculateLayout_7_3T2M1M1B(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 7), Dimensions: make([][]float64, 7)}
	if len(ARs) != 7 || len(types) != 7 {
		return layout, fmt.Errorf("3T2M1M1B layout requires 7 ARs and types")
	}

	// Row 1 (0, 1, 2)
	widths1, height1, err1 := calculateRowLayout(ARs[0:3], AW, spacing)
	if err1 != nil {
		return layout, fmt.Errorf("failed row 1 for 3T2M1M1B: %w", err1)
	}
	W0, W1, W2 := widths1[0], widths1[1], widths1[2]

	// Row 2 (3, 4)
	widths2, height2, err2 := calculateRowLayout(ARs[3:5], AW, spacing)
	if err2 != nil {
		return layout, fmt.Errorf("failed row 2 for 3T2M1M1B: %w", err2)
	}
	W3, W4 := widths2[0], widths2[1]

	// Row 3 (5)
	W5 := AW
	AR5 := ARs[5]
	height3 := 0.0
	if AR5 > 1e-6 {
		height3 = W5 / AR5
	} else {
		return layout, fmt.Errorf("invalid AR for pic 5 in 3T2M1M1B")
	}
	if height3 <= 1e-6 {
		return layout, fmt.Errorf("zero height for pic 5 in 3T2M1M1B")
	}

	// Row 4 (6)
	W6 := AW
	AR6 := ARs[6]
	height4 := 0.0
	if AR6 > 1e-6 {
		height4 = W6 / AR6
	} else {
		return layout, fmt.Errorf("invalid AR for pic 6 in 3T2M1M1B")
	}
	if height4 <= 1e-6 {
		return layout, fmt.Errorf("zero height for pic 6 in 3T2M1M1B")
	}

	layout.TotalHeight = height1 + spacing + height2 + spacing + height3 + spacing + height4
	layout.TotalWidth = AW

	currentX := 0.0
	layout.Positions[0] = []float64{currentX, 0}
	layout.Dimensions[0] = []float64{W0, height1}
	currentX += W0 + spacing
	layout.Positions[1] = []float64{currentX, 0}
	layout.Dimensions[1] = []float64{W1, height1}
	currentX += W1 + spacing
	layout.Positions[2] = []float64{currentX, 0}
	layout.Dimensions[2] = []float64{W2, height1}
	yRow2 := height1 + spacing
	layout.Positions[3] = []float64{0, yRow2}
	layout.Dimensions[3] = []float64{W3, height2}
	layout.Positions[4] = []float64{W3 + spacing, yRow2}
	layout.Dimensions[4] = []float64{W4, height2}
	yRow3 := yRow2 + height2 + spacing
	layout.Positions[5] = []float64{0, yRow3}
	layout.Dimensions[5] = []float64{W5, height3}
	yRow4 := yRow3 + height3 + spacing
	layout.Positions[6] = []float64{0, yRow4}
	layout.Dimensions[6] = []float64{W6, height4}

	return layout, nil
}

// calculateLayout_7_1T2M3M1B calculates the 1T-2M-3M-1B layout.
func calculateLayout_7_1T2M3M1B(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 7), Dimensions: make([][]float64, 7)}
	if len(ARs) != 7 || len(types) != 7 {
		return layout, fmt.Errorf("1T2M3M1B layout requires 7 ARs and types")
	}

	// Row 1 (0)
	W0 := AW
	AR0 := ARs[0]
	height1 := 0.0
	if AR0 > 1e-6 {
		height1 = W0 / AR0
	} else {
		return layout, fmt.Errorf("invalid AR for pic 0 in 1T2M3M1B")
	}
	if height1 <= 1e-6 {
		return layout, fmt.Errorf("zero height for pic 0 in 1T2M3M1B")
	}

	// Row 2 (1, 2)
	widths2, height2, err2 := calculateRowLayout(ARs[1:3], AW, spacing)
	if err2 != nil {
		return layout, fmt.Errorf("failed row 2 for 1T2M3M1B: %w", err2)
	}
	W1, W2 := widths2[0], widths2[1]

	// Row 3 (3, 4, 5)
	widths3, height3, err3 := calculateRowLayout(ARs[3:6], AW, spacing)
	if err3 != nil {
		return layout, fmt.Errorf("failed row 3 for 1T2M3M1B: %w", err3)
	}
	W3, W4, W5 := widths3[0], widths3[1], widths3[2]

	// Row 4 (6)
	W6 := AW
	AR6 := ARs[6]
	height4 := 0.0
	if AR6 > 1e-6 {
		height4 = W6 / AR6
	} else {
		return layout, fmt.Errorf("invalid AR for pic 6 in 1T2M3M1B")
	}
	if height4 <= 1e-6 {
		return layout, fmt.Errorf("zero height for pic 6 in 1T2M3M1B")
	}

	layout.TotalHeight = height1 + spacing + height2 + spacing + height3 + spacing + height4
	layout.TotalWidth = AW

	layout.Positions[0] = []float64{0, 0}
	layout.Dimensions[0] = []float64{W0, height1}
	yRow2 := height1 + spacing
	layout.Positions[1] = []float64{0, yRow2}
	layout.Dimensions[1] = []float64{W1, height2}
	layout.Positions[2] = []float64{W1 + spacing, yRow2}
	layout.Dimensions[2] = []float64{W2, height2}
	yRow3 := yRow2 + height2 + spacing
	currentX := 0.0
	layout.Positions[3] = []float64{currentX, yRow3}
	layout.Dimensions[3] = []float64{W3, height3}
	currentX += W3 + spacing
	layout.Positions[4] = []float64{currentX, yRow3}
	layout.Dimensions[4] = []float64{W4, height3}
	currentX += W4 + spacing
	layout.Positions[5] = []float64{currentX, yRow3}
	layout.Dimensions[5] = []float64{W5, height3}
	yRow4 := yRow3 + height3 + spacing
	layout.Positions[6] = []float64{0, yRow4}
	layout.Dimensions[6] = []float64{W6, height4}

	return layout, nil
}

// --- Main Calculation Function for 7 Pictures ---

func (e *ContinuousLayoutEngine) calculateSevenPicturesLayout(pictures []Picture, layoutAvailableHeight float64) (TemplateLayout, error) {
	if len(pictures) != 7 {
		return TemplateLayout{}, fmt.Errorf("incorrect number of pictures for 7-pic layout: %d", len(pictures))
	}

	spacing := e.imageSpacing
	AW := e.availableWidth

	// Get Aspect Ratios (W/H) and Types
	ARs := make([]float64, 7)
	types := make([]string, 7)
	validARs := true
	for i, pic := range pictures {
		if pic.Height > 0 && pic.Width > 0 {
			ARs[i] = float64(pic.Width) / float64(pic.Height)
			types[i] = getPictureType(ARs[i]) // Use global helper
		} else {
			ARs[i] = 1.0
			types[i] = "unknown"
			validARs = false
			fmt.Printf("Warning: Invalid dimensions for picture %d in 7-pic layout.\n", i)
		}
	}
	if !validARs {
		return TemplateLayout{}, fmt.Errorf("invalid dimensions in 7-pic layout")
	}

	// --- Define Layout Calculation Functions Map (Updated) ---
	type calcFuncType func(*ContinuousLayoutEngine, []float64, []string, float64, float64) (TemplateLayout, error)
	possibleLayouts := map[string]calcFuncType{
		"3T2M2B":   calculateLayout_7_3T2M2B,
		"2T3M2B":   calculateLayout_7_2T3M2B,
		"2T2M3B":   calculateLayout_7_2T2M3B,
		"3T3M1B":   calculateLayout_7_3T3M1B,
		"3L4R":     calculateLayout_7_3L4R,
		"4L3R":     calculateLayout_7_4L3R,
		"2T2M2M1B": calculateLayout_7_2T2M2M1B,
		"1T2M2M2B": calculateLayout_7_1T2M2M2B,
		"3T2M1M1B": calculateLayout_7_3T2M1M1B,
		"1T2M3M1B": calculateLayout_7_1T2M3M1B,
	}

	// --- Store results ---
	validLayouts := make(map[string]TemplateLayout)
	layoutAreas := make(map[string]float64)
	scaledLayouts := make(map[string]TemplateLayout)   // Store all calculated & scaled layouts
	layoutViolationFactors := make(map[string]float64) // Store violation factors for fallback
	var firstCalcError error

	// --- Calculate and Evaluate All Layouts ---
	for name, calcFunc := range possibleLayouts {
		layout, err := calcFunc(e, ARs, types, AW, spacing)
		if err != nil {
			fmt.Printf("Debug: Error calculating initial 7-pic layout %s: %v\n", name, err)
			if firstCalcError == nil {
				firstCalcError = fmt.Errorf("initial 7-pic layout %s: %w", name, err)
			}
			continue
		}

		// Scale Layout if Needed
		scale := 1.0
		if layout.TotalHeight > layoutAvailableHeight {
			if layout.TotalHeight > 1e-6 {
				scale = layoutAvailableHeight / layout.TotalHeight
				scaledLayout := TemplateLayout{TotalHeight: layout.TotalHeight * scale, TotalWidth: layout.TotalWidth}
				scaledLayout.Positions = make([][]float64, len(layout.Positions))
				scaledLayout.Dimensions = make([][]float64, len(layout.Dimensions))
				for i := range layout.Positions {
					// Error check for nil slices if necessary, assuming calcFunc initializes correctly
					if len(layout.Positions[i]) == 2 {
						scaledLayout.Positions[i] = []float64{layout.Positions[i][0] * scale, layout.Positions[i][1] * scale}
					}
					if len(layout.Dimensions[i]) == 2 {
						scaledLayout.Dimensions[i] = []float64{layout.Dimensions[i][0] * scale, layout.Dimensions[i][1] * scale}
					}
				}
				layout = scaledLayout
			} else {
				fmt.Printf("Debug: Layout %s has zero/tiny height, skipping scaling.\n", name)
				continue
			}
		}
		scaledLayouts[name] = layout // Store the final (potentially scaled) layout

		// --- Check Minimum Heights After Scaling & Calculate Violation Factor ---
		meetsScaledMin := true
		maxViolationFactor := 1.0 // Start assuming it meets minimums
		if !checkMinHeights(e, layout, types, 7) {
			meetsScaledMin = false
			// Calculate violation factor only if it failed
			for i, picType := range types {
				requiredMinHeight := getRequiredMinHeight(e, picType, len(pictures))
				if i < len(layout.Dimensions) && len(layout.Dimensions[i]) == 2 {
					actualHeight := layout.Dimensions[i][1]
					if actualHeight < requiredMinHeight {
						if actualHeight > 1e-6 {
							violationRatio := requiredMinHeight / actualHeight
							if violationRatio > maxViolationFactor {
								maxViolationFactor = violationRatio
							}
						} else {
							maxViolationFactor = math.Inf(1) // Assign infinite factor if actual height is zero
						}
					}
				} else {
					maxViolationFactor = math.Inf(1) // Treat invalid data as infinite violation
					break
				}
			}
			fmt.Printf("Debug: 7-Pic Layout %s failed minimum height check (Scale: %.2f, ViolationFactor: %.2f).\n", name, scale, maxViolationFactor)
		} else {
			fmt.Printf("Debug: 7-Pic Layout %s passed minimum height check (Scale: %.2f).\n", name, scale)
		}
		layoutViolationFactors[name] = maxViolationFactor // Store violation factor regardless

		// --- Store Strictly Valid Layout and Calculate Area ---
		if meetsScaledMin {
			validLayouts[name] = layout
			totalArea := 0.0
			for _, dim := range layout.Dimensions {
				if len(dim) == 2 {
					totalArea += dim[0] * dim[1]
				}
			}
			layoutAreas[name] = totalArea
			fmt.Printf("Debug: 7-Pic Layout %s stored as valid. Area: %.2f\n", name, totalArea)
		}
	}

	// --- Select Best Layout or Signal Split/New Page ---
	if len(validLayouts) > 0 {
		bestLayoutName := ""
		maxArea := -1.0
		for name, area := range layoutAreas {
			if area > maxArea {
				maxArea = area
				bestLayoutName = name
			}
		}
		fmt.Printf("Debug: Selected best fitting valid 7-pic layout: %s (Area: %.2f)\n", bestLayoutName, maxArea)
		return validLayouts[bestLayoutName], nil
	} else {
		hasWideOrTall := false
		for _, picType := range types {
			if picType == "wide" || picType == "tall" {
				hasWideOrTall = true
				break
			}
		}

		if hasWideOrTall {
			fmt.Println("Debug: No fitting layout for 7 pics with wide/tall images. Signaling force_new_page.")
			return TemplateLayout{}, fmt.Errorf("force_new_page")
		} else {
			fmt.Println("Debug: No fitting layout for 7 pics (no wide/tall). Signaling split_required.")
			return TemplateLayout{}, fmt.Errorf("split_required")
		}
	}
}
