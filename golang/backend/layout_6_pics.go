package backend

import (
	"fmt"
	"math"
)

// calculateSixPicturesLayout defines templates and calculates dimensions for 6 pictures using a 3x2 geometric approach.
func (e *ContinuousLayoutEngine) calculateSixPicturesLayout(pictures []Picture) (TemplateLayout, error) {
	// Template: Geometric 3x2 Grid

	if len(pictures) != 6 {
		return TemplateLayout{}, fmt.Errorf("incorrect number of pictures for 6-pic layout")
	}

	minAllowedHeight := 2500.00
	minAllowedPictureWidth := 1666.67
	spacing := e.imageSpacing
	AW := e.availableWidth

	layout := TemplateLayout{
		Positions:  make([][]float64, 6),
		Dimensions: make([][]float64, 6),
	}

	// Get Aspect Ratios (W/H)
	ARs := make([]float64, 6) // Size 6 for 6 pictures
	for i, pic := range pictures {
		if pic.Height > 0 && pic.Width > 0 {
			ARs[i] = float64(pic.Width) / float64(pic.Height)
		} else {
			ARs[i] = 1.0 // Default AR
		}
	}

	// --- Geometric Calculation (Allocate width based on height needs for 3 columns, 2 rows) ---
	invAR := make([]float64, 6) // Size 6 for 6 pictures
	for i, ar := range ARs {
		if ar > 1e-6 {
			invAR[i] = 1.0 / ar
		}
	}

	// Calculate height factors for each column (2 rows)
	heightFactorLeft := invAR[0] + invAR[3]  // Pic 0 + Pic 3
	heightFactorMid := invAR[1] + invAR[4]   // Pic 1 + Pic 4
	heightFactorRight := invAR[2] + invAR[5] // Pic 2 + Pic 5

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
	if columnsAvailableWidth > 1e-6 && totalInvHfSum > 1e-6 {
		W_left = columnsAvailableWidth * (invHfL / totalInvHfSum)
		W_mid = columnsAvailableWidth * (invHfM / totalInvHfSum)
		W_right = columnsAvailableWidth * (invHfR / totalInvHfSum)
	} else {
		// Fallback: equal width allocation (or could use row-based fallback)
		fmt.Println("Warning: Cannot calculate 3x2 geometry factors, using equal widths.")
		equalWidth := columnsAvailableWidth / 3.0
		if equalWidth < 1.0 {
			equalWidth = 1.0
		}
		W_left, W_mid, W_right = equalWidth, equalWidth, equalWidth
	}

	if W_left <= 0 || W_mid <= 0 || W_right <= 0 {
		// Fallback if calculated widths are invalid
		fmt.Println("Warning: Calculated column widths non-positive, using equal widths.")
		equalWidth := columnsAvailableWidth / 3.0
		if equalWidth < 1.0 {
			equalWidth = 1.0
		}
		W_left, W_mid, W_right = equalWidth, equalWidth, equalWidth
	}

	// Calculate initial heights based on allocated widths
	H0 := W_left * invAR[0]
	H1 := W_mid * invAR[1]
	H2 := W_right * invAR[2]
	H3 := W_left * invAR[3]
	H4 := W_mid * invAR[4]
	H5 := W_right * invAR[5]

	// Determine final row heights
	H_top := math.Max(H0, math.Max(H1, H2))
	H_bottom := math.Max(H3, math.Max(H4, H5))

	// Recalculate widths based on final row heights to maintain AR
	W0 := H_top * ARs[0]
	W1 := H_top * ARs[1]
	W2 := H_top * ARs[2]
	W3 := H_bottom * ARs[3]
	W4 := H_bottom * ARs[4]
	W5 := H_bottom * ARs[5]

	// Scale rows proportionally if their total width exceeds available width
	topWidth := W0 + W1 + W2 + 2*spacing
	if topWidth > AW+1e-3 { // Allow float tolerance
		scale := AW / topWidth
		H_top *= scale
		W0 *= scale
		W1 *= scale
		W2 *= scale
	}
	bottomWidth := W3 + W4 + W5 + 2*spacing
	if bottomWidth > AW+1e-3 { // Allow float tolerance
		scale := AW / bottomWidth
		H_bottom *= scale
		W3 *= scale
		W4 *= scale
		W5 *= scale
	}

	// --- Check Minimums ---
	meetsMin := true
	if H_top < minAllowedHeight || H_bottom < minAllowedHeight {
		meetsMin = false
	}
	widths := []float64{W0, W1, W2, W3, W4, W5} // Use array for iteration
	for _, w := range widths {
		if w < minAllowedPictureWidth {
			meetsMin = false
			break
		}
	}

	if !meetsMin {
		fmt.Println("Warning: 6-picture layout (3x2 Geo) violates minimum dimensions.")
	}

	// --- Populate Layout Struct ---
	layout.TotalHeight = H_top + spacing + H_bottom
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
	layout.Dimensions[3] = []float64{W3, H_bottom}
	currentX += W3 + spacing
	layout.Positions[4] = []float64{currentX, H_top + spacing}
	layout.Dimensions[4] = []float64{W4, H_bottom}
	currentX += W4 + spacing
	layout.Positions[5] = []float64{currentX, H_top + spacing}
	layout.Dimensions[5] = []float64{W5, H_bottom}

	return layout, nil
}
