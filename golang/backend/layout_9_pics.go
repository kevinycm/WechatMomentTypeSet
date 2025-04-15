package backend

import (
	"fmt"
	"math"
)

// calculateNinePicturesLayout calculates dimensions for 9 pictures using a 3x3 geometric approach.
func (e *ContinuousLayoutEngine) calculateNinePicturesLayout(pictures []Picture) (TemplateLayout, error) {
	// Template: Geometric 3x3 Grid

	if len(pictures) != 9 {
		return TemplateLayout{}, fmt.Errorf("incorrect number of pictures for 9-pic layout")
	}

	minAllowedHeight := 2500.00
	minAllowedPictureWidth := 1666.67
	spacing := e.imageSpacing
	AW := e.availableWidth

	layout := TemplateLayout{
		Positions:  make([][]float64, 9),
		Dimensions: make([][]float64, 9),
	}

	// Get Aspect Ratios (W/H)
	ARs := make([]float64, 9)
	for i, pic := range pictures {
		if pic.Height > 0 && pic.Width > 0 {
			ARs[i] = float64(pic.Width) / float64(pic.Height)
		} else {
			ARs[i] = 1.0 // Default AR
		}
	}

	// --- Geometric Calculation (Allocate width based on height needs for 3 columns, 3 rows) ---
	invAR := make([]float64, 9)
	for i, ar := range ARs {
		if ar > 1e-6 {
			invAR[i] = 1.0 / ar
		}
	}

	// Calculate height factors for each column (3 rows)
	heightFactorLeft := invAR[0] + invAR[3] + invAR[6]  // Pic 0 + Pic 3 + Pic 6
	heightFactorMid := invAR[1] + invAR[4] + invAR[7]   // Pic 1 + Pic 4 + Pic 7
	heightFactorRight := invAR[2] + invAR[5] + invAR[8] // Pic 2 + Pic 5 + Pic 8

	// Inverse height factors for width allocation
	invHfL := 0.0
	if heightFactorLeft > 1e-6 {
		invHfL = 1.0 / heightFactorLeft
	}
	invHfM := 0.0
	if heightFactorMid > 1e-6 {
		invHfM = 1.0 / heightFactorMid
	}
	invHfR := 0.0
	if heightFactorRight > 1e-6 {
		invHfR = 1.0 / heightFactorRight
	}
	totalInvHfSum := invHfL + invHfM + invHfR

	// Available width for 3 columns
	columnsAvailableWidth := AW - 2*spacing
	W_left, W_mid, W_right := 0.0, 0.0, 0.0

	// Allocate column widths proportionally to inverse height factors
	if columnsAvailableWidth <= 1e-6 || totalInvHfSum <= 1e-6 {
		// Cannot calculate geometric factors
		return TemplateLayout{}, fmt.Errorf("cannot calculate 3x3 geometry factors (zero width or height factor sum)")
	}
	W_left = columnsAvailableWidth * (invHfL / totalInvHfSum)
	W_mid = columnsAvailableWidth * (invHfM / totalInvHfSum)
	W_right = columnsAvailableWidth * (invHfR / totalInvHfSum)

	if W_left <= 0 || W_mid <= 0 || W_right <= 0 {
		// Calculated widths are invalid
		return TemplateLayout{}, fmt.Errorf("calculated 3x3 column widths are non-positive (W_L:%.2f, W_M:%.2f, W_R:%.2f)", W_left, W_mid, W_right)
	}

	// Calculate initial heights based on allocated widths
	H0 := W_left * invAR[0]
	H1 := W_mid * invAR[1]
	H2 := W_right * invAR[2]
	H3 := W_left * invAR[3]
	H4 := W_mid * invAR[4]
	H5 := W_right * invAR[5]
	H6 := W_left * invAR[6]
	H7 := W_mid * invAR[7]
	H8 := W_right * invAR[8]

	// Determine final row heights
	H_top := math.Max(H0, math.Max(H1, H2))
	H_mid_row := math.Max(H3, math.Max(H4, H5))
	H_bottom := math.Max(H6, math.Max(H7, H8))

	// Recalculate widths based on final row heights to maintain AR
	W0 := H_top * ARs[0]
	W1 := H_top * ARs[1]
	W2 := H_top * ARs[2]
	W3 := H_mid_row * ARs[3]
	W4 := H_mid_row * ARs[4]
	W5 := H_mid_row * ARs[5]
	W6 := H_bottom * ARs[6]
	W7 := H_bottom * ARs[7]
	W8 := H_bottom * ARs[8]

	// Scale rows proportionally if their total width exceeds available width
	topWidth := W0 + W1 + W2 + 2*spacing
	if topWidth > AW+1e-3 { // Allow float tolerance
		scale := AW / topWidth
		H_top *= scale
		W0 *= scale
		W1 *= scale
		W2 *= scale
	}
	midWidth := W3 + W4 + W5 + 2*spacing
	if midWidth > AW+1e-3 { // Allow float tolerance
		scale := AW / midWidth
		H_mid_row *= scale
		W3 *= scale
		W4 *= scale
		W5 *= scale
	}
	bottomWidth := W6 + W7 + W8 + 2*spacing
	if bottomWidth > AW+1e-3 { // Allow float tolerance
		scale := AW / bottomWidth
		H_bottom *= scale
		W6 *= scale
		W7 *= scale
		W8 *= scale
	}

	// --- Check Minimums ---
	meetsMin := true
	if H_top < minAllowedHeight || H_mid_row < minAllowedHeight || H_bottom < minAllowedHeight {
		meetsMin = false
	}
	widths := []float64{W0, W1, W2, W3, W4, W5, W6, W7, W8}
	for _, w := range widths {
		if w < minAllowedPictureWidth {
			meetsMin = false
			break
		}
	}

	if !meetsMin {
		fmt.Println("Warning: 9-picture layout (3x3 Geo) violates minimum dimensions.")
		// Return the calculated layout anyway, despite violation
	}

	// --- Populate Layout Struct ---
	layout.TotalHeight = H_top + spacing + H_mid_row + spacing + H_bottom
	layout.TotalWidth = AW

	// Row 1
	currentX := 0.0
	layout.Positions[0] = []float64{currentX, 0}
	layout.Dimensions[0] = []float64{W0, H_top}
	currentX += W0 + spacing
	layout.Positions[1] = []float64{currentX, 0}
	layout.Dimensions[1] = []float64{W1, H_top}
	currentX += W1 + spacing
	layout.Positions[2] = []float64{currentX, 0}
	layout.Dimensions[2] = []float64{W2, H_top}

	// Row 2
	currentX = 0.0
	layout.Positions[3] = []float64{currentX, H_top + spacing}
	layout.Dimensions[3] = []float64{W3, H_mid_row}
	currentX += W3 + spacing
	layout.Positions[4] = []float64{currentX, H_top + spacing}
	layout.Dimensions[4] = []float64{W4, H_mid_row}
	currentX += W4 + spacing
	layout.Positions[5] = []float64{currentX, H_top + spacing}
	layout.Dimensions[5] = []float64{W5, H_mid_row}

	// Row 3
	currentX = 0.0
	layout.Positions[6] = []float64{currentX, H_top + spacing + H_mid_row + spacing}
	layout.Dimensions[6] = []float64{W6, H_bottom}
	currentX += W6 + spacing
	layout.Positions[7] = []float64{currentX, H_top + spacing + H_mid_row + spacing}
	layout.Dimensions[7] = []float64{W7, H_bottom}
	currentX += W7 + spacing
	layout.Positions[8] = []float64{currentX, H_top + spacing + H_mid_row + spacing}
	layout.Dimensions[8] = []float64{W8, H_bottom}

	return layout, nil
}
