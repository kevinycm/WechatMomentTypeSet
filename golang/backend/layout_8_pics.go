package backend

import (
	"fmt"
	"math"
)

// --- 8-Picture Layout Calculation Functions ---

// calculateLayout_8_3T3M2B: 3 Top, 3 Middle, 2 Bottom
func calculateLayout_8_3T3M2B(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 8), Dimensions: make([][]float64, 8)}
	if len(ARs) != 8 || len(types) != 8 {
		return layout, fmt.Errorf("3T3M2B layout requires 8 ARs/types")
	}

	widths1, h1, err1 := calculateRowLayout(ARs[0:3], AW, spacing)
	if err1 != nil {
		return layout, fmt.Errorf("row1 failed for 3T3M2B: %w", err1)
	}
	widths2, h2, err2 := calculateRowLayout(ARs[3:6], AW, spacing)
	if err2 != nil {
		return layout, fmt.Errorf("row2 failed for 3T3M2B: %w", err2)
	}
	widths3, h3, err3 := calculateRowLayout(ARs[6:8], AW, spacing)
	if err3 != nil {
		return layout, fmt.Errorf("row3 failed for 3T3M2B: %w", err3)
	}

	layout.TotalHeight = h1 + spacing + h2 + spacing + h3
	layout.TotalWidth = AW

	y := 0.0
	currentX := 0.0
	// Row 1
	for i := 0; i < 3; i++ {
		layout.Positions[i] = []float64{currentX, y}
		layout.Dimensions[i] = []float64{widths1[i], h1}
		currentX += widths1[i] + spacing
	}
	// Row 2
	y += h1 + spacing
	currentX = 0.0
	for i := 0; i < 3; i++ {
		layout.Positions[i+3] = []float64{currentX, y}
		layout.Dimensions[i+3] = []float64{widths2[i], h2}
		currentX += widths2[i] + spacing
	}
	// Row 3
	y += h2 + spacing
	currentX = 0.0
	for i := 0; i < 2; i++ {
		layout.Positions[i+6] = []float64{currentX, y}
		layout.Dimensions[i+6] = []float64{widths3[i], h3}
		currentX += widths3[i] + spacing
	}
	return layout, nil
}

// calculateLayout_8_2T3M3B: 2 Top, 3 Middle, 3 Bottom
func calculateLayout_8_2T3M3B(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 8), Dimensions: make([][]float64, 8)}
	if len(ARs) != 8 || len(types) != 8 {
		return layout, fmt.Errorf("2T3M3B layout requires 8 ARs/types")
	}

	widths1, h1, err1 := calculateRowLayout(ARs[0:2], AW, spacing)
	if err1 != nil {
		return layout, fmt.Errorf("row1 failed for 2T3M3B: %w", err1)
	}
	widths2, h2, err2 := calculateRowLayout(ARs[2:5], AW, spacing)
	if err2 != nil {
		return layout, fmt.Errorf("row2 failed for 2T3M3B: %w", err2)
	}
	widths3, h3, err3 := calculateRowLayout(ARs[5:8], AW, spacing)
	if err3 != nil {
		return layout, fmt.Errorf("row3 failed for 2T3M3B: %w", err3)
	}

	layout.TotalHeight = h1 + spacing + h2 + spacing + h3
	layout.TotalWidth = AW

	y := 0.0
	currentX := 0.0
	// Row 1
	for i := 0; i < 2; i++ {
		layout.Positions[i] = []float64{currentX, y}
		layout.Dimensions[i] = []float64{widths1[i], h1}
		currentX += widths1[i] + spacing
	}
	// Row 2
	y += h1 + spacing
	currentX = 0.0
	for i := 0; i < 3; i++ {
		layout.Positions[i+2] = []float64{currentX, y}
		layout.Dimensions[i+2] = []float64{widths2[i], h2}
		currentX += widths2[i] + spacing
	}
	// Row 3
	y += h2 + spacing
	currentX = 0.0
	for i := 0; i < 3; i++ {
		layout.Positions[i+5] = []float64{currentX, y}
		layout.Dimensions[i+5] = []float64{widths3[i], h3}
		currentX += widths3[i] + spacing
	}
	return layout, nil
}

// calculateLayout_8_3T2M3B: 3 Top, 2 Middle, 3 Bottom
func calculateLayout_8_3T2M3B(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 8), Dimensions: make([][]float64, 8)}
	if len(ARs) != 8 || len(types) != 8 {
		return layout, fmt.Errorf("3T2M3B layout requires 8 ARs/types")
	}

	widths1, h1, err1 := calculateRowLayout(ARs[0:3], AW, spacing)
	if err1 != nil {
		return layout, fmt.Errorf("row1 failed for 3T2M3B: %w", err1)
	}
	widths2, h2, err2 := calculateRowLayout(ARs[3:5], AW, spacing)
	if err2 != nil {
		return layout, fmt.Errorf("row2 failed for 3T2M3B: %w", err2)
	}
	widths3, h3, err3 := calculateRowLayout(ARs[5:8], AW, spacing)
	if err3 != nil {
		return layout, fmt.Errorf("row3 failed for 3T2M3B: %w", err3)
	}

	layout.TotalHeight = h1 + spacing + h2 + spacing + h3
	layout.TotalWidth = AW

	y := 0.0
	currentX := 0.0
	// Row 1
	for i := 0; i < 3; i++ {
		layout.Positions[i] = []float64{currentX, y}
		layout.Dimensions[i] = []float64{widths1[i], h1}
		currentX += widths1[i] + spacing
	}
	// Row 2
	y += h1 + spacing
	currentX = 0.0
	for i := 0; i < 2; i++ {
		layout.Positions[i+3] = []float64{currentX, y}
		layout.Dimensions[i+3] = []float64{widths2[i], h2}
		currentX += widths2[i] + spacing
	}
	// Row 3
	y += h2 + spacing
	currentX = 0.0
	for i := 0; i < 3; i++ {
		layout.Positions[i+5] = []float64{currentX, y}
		layout.Dimensions[i+5] = []float64{widths3[i], h3}
		currentX += widths3[i] + spacing
	}
	return layout, nil
}

// calculateLayout_8_2T2M2M2B: 2 Top, 2 Mid1, 2 Mid2, 2 Bottom
func calculateLayout_8_2T2M2M2B(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 8), Dimensions: make([][]float64, 8)}
	if len(ARs) != 8 || len(types) != 8 {
		return layout, fmt.Errorf("2T2M2M2B layout requires 8 ARs/types")
	}

	widths1, h1, err1 := calculateRowLayout(ARs[0:2], AW, spacing)
	if err1 != nil {
		return layout, fmt.Errorf("row1 failed for 2T2M2M2B: %w", err1)
	}
	widths2, h2, err2 := calculateRowLayout(ARs[2:4], AW, spacing)
	if err2 != nil {
		return layout, fmt.Errorf("row2 failed for 2T2M2M2B: %w", err2)
	}
	widths3, h3, err3 := calculateRowLayout(ARs[4:6], AW, spacing)
	if err3 != nil {
		return layout, fmt.Errorf("row3 failed for 2T2M2M2B: %w", err3)
	}
	widths4, h4, err4 := calculateRowLayout(ARs[6:8], AW, spacing)
	if err4 != nil {
		return layout, fmt.Errorf("row4 failed for 2T2M2M2B: %w", err4)
	}

	layout.TotalHeight = h1 + spacing + h2 + spacing + h3 + spacing + h4
	layout.TotalWidth = AW

	y := 0.0
	currentX := 0.0
	// Row 1
	for i := 0; i < 2; i++ {
		layout.Positions[i] = []float64{currentX, y}
		layout.Dimensions[i] = []float64{widths1[i], h1}
		currentX += widths1[i] + spacing
	}
	// Row 2
	y += h1 + spacing
	currentX = 0.0
	for i := 0; i < 2; i++ {
		layout.Positions[i+2] = []float64{currentX, y}
		layout.Dimensions[i+2] = []float64{widths2[i], h2}
		currentX += widths2[i] + spacing
	}
	// Row 3
	y += h2 + spacing
	currentX = 0.0
	for i := 0; i < 2; i++ {
		layout.Positions[i+4] = []float64{currentX, y}
		layout.Dimensions[i+4] = []float64{widths3[i], h3}
		currentX += widths3[i] + spacing
	}
	// Row 4
	y += h3 + spacing
	currentX = 0.0
	for i := 0; i < 2; i++ {
		layout.Positions[i+6] = []float64{currentX, y}
		layout.Dimensions[i+6] = []float64{widths4[i], h4}
		currentX += widths4[i] + spacing
	}

	return layout, nil
}

// calculateLayout_8_3T2M2M1B: 3 Top, 2 Mid1, 2 Mid2, 1 Bottom
func calculateLayout_8_3T2M2M1B(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 8), Dimensions: make([][]float64, 8)}
	if len(ARs) != 8 || len(types) != 8 {
		return layout, fmt.Errorf("3T2M2M1B layout requires 8 ARs/types")
	}

	widths1, h1, err1 := calculateRowLayout(ARs[0:3], AW, spacing)
	if err1 != nil {
		return layout, fmt.Errorf("row1 failed for 3T2M2M1B: %w", err1)
	}
	widths2, h2, err2 := calculateRowLayout(ARs[3:5], AW, spacing)
	if err2 != nil {
		return layout, fmt.Errorf("row2 failed for 3T2M2M1B: %w", err2)
	}
	widths3, h3, err3 := calculateRowLayout(ARs[5:7], AW, spacing)
	if err3 != nil {
		return layout, fmt.Errorf("row3 failed for 3T2M2M1B: %w", err3)
	}
	// Row 4 (Single Pic)
	w4 := AW
	ar4 := ARs[7]
	h4 := 0.0
	if ar4 <= 1e-6 {
		return layout, fmt.Errorf("invalid AR for bottom pic in 3T2M2M1B")
	}
	h4 = w4 / ar4
	if h4 <= 1e-6 {
		return layout, fmt.Errorf("zero height for bottom pic in 3T2M2M1B")
	}

	layout.TotalHeight = h1 + spacing + h2 + spacing + h3 + spacing + h4
	layout.TotalWidth = AW

	y := 0.0
	currentX := 0.0
	// Row 1
	for i := 0; i < 3; i++ {
		layout.Positions[i] = []float64{currentX, y}
		layout.Dimensions[i] = []float64{widths1[i], h1}
		currentX += widths1[i] + spacing
	}
	// Row 2
	y += h1 + spacing
	currentX = 0.0
	for i := 0; i < 2; i++ {
		layout.Positions[i+3] = []float64{currentX, y}
		layout.Dimensions[i+3] = []float64{widths2[i], h2}
		currentX += widths2[i] + spacing
	}
	// Row 3
	y += h2 + spacing
	currentX = 0.0
	for i := 0; i < 2; i++ {
		layout.Positions[i+5] = []float64{currentX, y}
		layout.Dimensions[i+5] = []float64{widths3[i], h3}
		currentX += widths3[i] + spacing
	}
	// Row 4
	y += h3 + spacing
	layout.Positions[7] = []float64{0, y}
	layout.Dimensions[7] = []float64{w4, h4}

	return layout, nil
}

// calculateLayout_8_1T2M2M3B: 1 Top, 2 Mid1, 2 Mid2, 3 Bottom
func calculateLayout_8_1T2M2M3B(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 8), Dimensions: make([][]float64, 8)}
	if len(ARs) != 8 || len(types) != 8 {
		return layout, fmt.Errorf("1T2M2M3B layout requires 8 ARs/types")
	}

	// Row 1 (Single Pic)
	w1 := AW
	ar1 := ARs[0]
	h1 := 0.0
	if ar1 <= 1e-6 {
		return layout, fmt.Errorf("invalid AR for top pic in 1T2M2M3B")
	}
	h1 = w1 / ar1
	if h1 <= 1e-6 {
		return layout, fmt.Errorf("zero height for top pic in 1T2M2M3B")
	}

	widths2, h2, err2 := calculateRowLayout(ARs[1:3], AW, spacing)
	if err2 != nil {
		return layout, fmt.Errorf("row2 failed for 1T2M2M3B: %w", err2)
	}
	widths3, h3, err3 := calculateRowLayout(ARs[3:5], AW, spacing)
	if err3 != nil {
		return layout, fmt.Errorf("row3 failed for 1T2M2M3B: %w", err3)
	}
	widths4, h4, err4 := calculateRowLayout(ARs[5:8], AW, spacing)
	if err4 != nil {
		return layout, fmt.Errorf("row4 failed for 1T2M2M3B: %w", err4)
	}

	layout.TotalHeight = h1 + spacing + h2 + spacing + h3 + spacing + h4
	layout.TotalWidth = AW

	y := 0.0
	// Row 1
	layout.Positions[0] = []float64{0, y}
	layout.Dimensions[0] = []float64{w1, h1}
	// Row 2
	y += h1 + spacing
	currentX := 0.0
	for i := 0; i < 2; i++ {
		layout.Positions[i+1] = []float64{currentX, y}
		layout.Dimensions[i+1] = []float64{widths2[i], h2}
		currentX += widths2[i] + spacing
	}
	// Row 3
	y += h2 + spacing
	currentX = 0.0
	for i := 0; i < 2; i++ {
		layout.Positions[i+3] = []float64{currentX, y}
		layout.Dimensions[i+3] = []float64{widths3[i], h3}
		currentX += widths3[i] + spacing
	}
	// Row 4
	y += h3 + spacing
	currentX = 0.0
	for i := 0; i < 3; i++ {
		layout.Positions[i+5] = []float64{currentX, y}
		layout.Dimensions[i+5] = []float64{widths4[i], h4}
		currentX += widths4[i] + spacing
	}

	return layout, nil
}

// calculateLayout_8_2T3M2M1B: 2 Top, 3 Mid1, 2 Mid2, 1 Bottom
func calculateLayout_8_2T3M2M1B(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 8), Dimensions: make([][]float64, 8)}
	if len(ARs) != 8 || len(types) != 8 {
		return layout, fmt.Errorf("2T3M2M1B layout requires 8 ARs/types")
	}

	widths1, h1, err1 := calculateRowLayout(ARs[0:2], AW, spacing)
	if err1 != nil {
		return layout, fmt.Errorf("row1 failed for 2T3M2M1B: %w", err1)
	}
	widths2, h2, err2 := calculateRowLayout(ARs[2:5], AW, spacing)
	if err2 != nil {
		return layout, fmt.Errorf("row2 failed for 2T3M2M1B: %w", err2)
	}
	widths3, h3, err3 := calculateRowLayout(ARs[5:7], AW, spacing)
	if err3 != nil {
		return layout, fmt.Errorf("row3 failed for 2T3M2M1B: %w", err3)
	}
	// Row 4 (Single Pic)
	w4 := AW
	ar4 := ARs[7]
	h4 := 0.0
	if ar4 <= 1e-6 {
		return layout, fmt.Errorf("invalid AR for bottom pic in 2T3M2M1B")
	}
	h4 = w4 / ar4
	if h4 <= 1e-6 {
		return layout, fmt.Errorf("zero height for bottom pic in 2T3M2M1B")
	}

	layout.TotalHeight = h1 + spacing + h2 + spacing + h3 + spacing + h4
	layout.TotalWidth = AW

	y := 0.0
	currentX := 0.0
	// Row 1
	for i := 0; i < 2; i++ {
		layout.Positions[i] = []float64{currentX, y}
		layout.Dimensions[i] = []float64{widths1[i], h1}
		currentX += widths1[i] + spacing
	}
	// Row 2
	y += h1 + spacing
	currentX = 0.0
	for i := 0; i < 3; i++ {
		layout.Positions[i+2] = []float64{currentX, y}
		layout.Dimensions[i+2] = []float64{widths2[i], h2}
		currentX += widths2[i] + spacing
	}
	// Row 3
	y += h2 + spacing
	currentX = 0.0
	for i := 0; i < 2; i++ {
		layout.Positions[i+5] = []float64{currentX, y}
		layout.Dimensions[i+5] = []float64{widths3[i], h3}
		currentX += widths3[i] + spacing
	}
	// Row 4
	y += h3 + spacing
	layout.Positions[7] = []float64{0, y}
	layout.Dimensions[7] = []float64{w4, h4}

	return layout, nil
}

// calculateLayout_8_1T2M3M2B: 1 Top, 2 Mid1, 3 Mid2, 2 Bottom
func calculateLayout_8_1T2M3M2B(e *ContinuousLayoutEngine, ARs []float64, types []string, AW, spacing float64) (TemplateLayout, error) {
	layout := TemplateLayout{Positions: make([][]float64, 8), Dimensions: make([][]float64, 8)}
	if len(ARs) != 8 || len(types) != 8 {
		return layout, fmt.Errorf("1T2M3M2B layout requires 8 ARs/types")
	}

	// Row 1 (Single Pic)
	w1 := AW
	ar1 := ARs[0]
	h1 := 0.0
	if ar1 <= 1e-6 {
		return layout, fmt.Errorf("invalid AR for top pic in 1T2M3M2B")
	}
	h1 = w1 / ar1
	if h1 <= 1e-6 {
		return layout, fmt.Errorf("zero height for top pic in 1T2M3M2B")
	}

	widths2, h2, err2 := calculateRowLayout(ARs[1:3], AW, spacing)
	if err2 != nil {
		return layout, fmt.Errorf("row2 failed for 1T2M3M2B: %w", err2)
	}
	widths3, h3, err3 := calculateRowLayout(ARs[3:6], AW, spacing)
	if err3 != nil {
		return layout, fmt.Errorf("row3 failed for 1T2M3M2B: %w", err3)
	}
	widths4, h4, err4 := calculateRowLayout(ARs[6:8], AW, spacing)
	if err4 != nil {
		return layout, fmt.Errorf("row4 failed for 1T2M3M2B: %w", err4)
	}

	layout.TotalHeight = h1 + spacing + h2 + spacing + h3 + spacing + h4
	layout.TotalWidth = AW

	y := 0.0
	// Row 1
	layout.Positions[0] = []float64{0, y}
	layout.Dimensions[0] = []float64{w1, h1}
	// Row 2
	y += h1 + spacing
	currentX := 0.0
	for i := 0; i < 2; i++ {
		layout.Positions[i+1] = []float64{currentX, y}
		layout.Dimensions[i+1] = []float64{widths2[i], h2}
		currentX += widths2[i] + spacing
	}
	// Row 3
	y += h2 + spacing
	currentX = 0.0
	for i := 0; i < 3; i++ {
		layout.Positions[i+3] = []float64{currentX, y}
		layout.Dimensions[i+3] = []float64{widths3[i], h3}
		currentX += widths3[i] + spacing
	}
	// Row 4
	y += h3 + spacing
	currentX = 0.0
	for i := 0; i < 2; i++ {
		layout.Positions[i+6] = []float64{currentX, y}
		layout.Dimensions[i+6] = []float64{widths4[i], h4}
		currentX += widths4[i] + spacing
	}

	return layout, nil
}

// --- Main Calculation Function for 8 Pictures ---
func (e *ContinuousLayoutEngine) calculateEightPicturesLayout(pictures []Picture, layoutAvailableHeight float64) (TemplateLayout, error) {
	numPics := 8
	if len(pictures) != numPics {
		return TemplateLayout{}, fmt.Errorf("incorrect number of pictures for %d-pic layout: %d", numPics, len(pictures))
	}

	spacing := e.imageSpacing
	AW := e.availableWidth

	// Get Aspect Ratios (W/H) and Types
	ARs := make([]float64, numPics)
	types := make([]string, numPics)
	validARs := true
	for i, pic := range pictures {
		if pic.Height > 0 && pic.Width > 0 {
			ARs[i] = float64(pic.Width) / float64(pic.Height)
			types[i] = getPictureType(ARs[i])
		} else {
			ARs[i] = 1.0 // Default AR
			types[i] = "unknown"
			validARs = false
			fmt.Printf("Warning: Invalid dimensions for picture %d in %d-pic layout.\n", i, numPics)
		}
	}
	if !validARs {
		return TemplateLayout{}, fmt.Errorf("invalid dimensions encountered in %d-pic layout", numPics)
	}

	// --- Define Layout Calculation Functions Map ---
	type calcFuncType func(*ContinuousLayoutEngine, []float64, []string, float64, float64) (TemplateLayout, error)
	possibleLayouts := map[string]calcFuncType{
		"3T3M2B":   calculateLayout_8_3T3M2B,
		"2T3M3B":   calculateLayout_8_2T3M3B,
		"3T2M3B":   calculateLayout_8_3T2M3B,
		"2T2M2M2B": calculateLayout_8_2T2M2M2B,
		"3T2M2M1B": calculateLayout_8_3T2M2M1B,
		"1T2M2M3B": calculateLayout_8_1T2M2M3B,
		"2T3M2M1B": calculateLayout_8_2T3M2M1B,
		"1T2M3M2B": calculateLayout_8_1T2M3M2B,
	}

	// --- Store results from all layout attempts ---
	validLayouts := make(map[string]TemplateLayout)
	layoutAreas := make(map[string]float64)
	scaledLayouts := make(map[string]TemplateLayout)
	layoutViolationFactors := make(map[string]float64)
	var firstCalcError error

	// --- Calculate and Evaluate All Layouts ---
	for name, calcFunc := range possibleLayouts {
		layout, err := calcFunc(e, ARs, types, AW, spacing)
		if err != nil {
			fmt.Printf("Debug: Error calculating initial %d-pic layout %s: %v\n", numPics, name, err)
			if firstCalcError == nil {
				firstCalcError = fmt.Errorf("initial %d-pic layout %s: %w", numPics, name, err)
			}
			layoutViolationFactors[name] = math.Inf(1) // Mark as non-viable
			continue
		}

		// --- Scale Layout if Needed ---
		scale := 1.0
		if layout.TotalHeight > layoutAvailableHeight {
			if layout.TotalHeight > 1e-6 {
				scale = layoutAvailableHeight / layout.TotalHeight
				scaledLayout := TemplateLayout{
					Positions:   make([][]float64, len(layout.Positions)),
					Dimensions:  make([][]float64, len(layout.Dimensions)),
					TotalHeight: layout.TotalHeight * scale,
					TotalWidth:  layout.TotalWidth, // Assuming layout maintains AW
				}
				for i := range layout.Positions {
					if len(layout.Positions[i]) == 2 {
						scaledLayout.Positions[i] = []float64{layout.Positions[i][0] * scale, layout.Positions[i][1] * scale}
					}
					if len(layout.Dimensions[i]) == 2 {
						scaledLayout.Dimensions[i] = []float64{layout.Dimensions[i][0] * scale, layout.Dimensions[i][1] * scale}
					}
				}
				layout = scaledLayout // Use the scaled layout
			} else {
				fmt.Printf("Debug: %d-Pic Layout %s has zero/tiny height, skipping scaling.\n", numPics, name)
				layoutViolationFactors[name] = math.Inf(1) // Mark as non-viable
				continue
			}
		}
		scaledLayouts[name] = layout // Store the final (potentially scaled) layout

		// --- Check Minimum Heights After Scaling & Calculate Violation Factor ---
		meetsScaledMin := true
		maxViolationFactor := 1.0
		for i, picType := range types {
			requiredMinHeight := getRequiredMinHeight(e, picType, numPics) // Use numPics
			if i < len(layout.Dimensions) && len(layout.Dimensions[i]) == 2 {
				actualHeight := layout.Dimensions[i][1]
				if actualHeight < requiredMinHeight {
					meetsScaledMin = false
					if actualHeight > 1e-6 {
						violationRatio := requiredMinHeight / actualHeight
						if violationRatio > maxViolationFactor {
							maxViolationFactor = violationRatio
						}
					} else {
						maxViolationFactor = math.Inf(1)
					}
				}
			} else {
				fmt.Printf("Warning: Invalid dimensions data for %d-pic layout %s, picture %d\n", numPics, name, i)
				meetsScaledMin = false
				maxViolationFactor = math.Inf(1)
				break
			}
		}
		layoutViolationFactors[name] = maxViolationFactor

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
			fmt.Printf("Debug: %d-Pic Layout %s valid (Scale: %.2f), Area: %.2f\n", numPics, name, scale, totalArea)
		} else {
			fmt.Printf("Debug: %d-Pic Layout %s failed minimum height check (Scale: %.2f, ViolationFactor: %.2f).\n", numPics, name, scale, maxViolationFactor)
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
		fmt.Printf("Debug: Selected best fitting valid %d-pic layout: %s (Area: %.2f)\n", numPics, bestLayoutName, maxArea)
		return validLayouts[bestLayoutName], nil
	} else {
		// Fallback logic
		hasWideOrTall := false
		for _, picType := range types {
			if picType == "wide" || picType == "tall" {
				hasWideOrTall = true
				break
			}
		}

		if hasWideOrTall {
			fmt.Printf("Debug: No fitting layout for %d pics with wide/tall images. Signaling force_new_page.\n", numPics)
			return TemplateLayout{}, fmt.Errorf("force_new_page")
		} else {
			fmt.Printf("Debug: No fitting layout for %d pics (no wide/tall). Signaling split_required.\n", numPics)
			return TemplateLayout{}, fmt.Errorf("split_required")
		}
		// TODO: Consider if fallback using the layout with minimum violation factor is needed later.
	}
}
