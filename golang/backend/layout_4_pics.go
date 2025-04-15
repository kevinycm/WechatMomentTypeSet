package backend

import (
	"fmt"
	"math"
)

// calculateFourPicturesLayout defines templates and calculates dimensions for 4 pictures.
func (e *ContinuousLayoutEngine) calculateFourPicturesLayout(pictures []Picture) (TemplateLayout, error) {
	// Template: Geometric 2x2 Grid Calculation

	if len(pictures) != 4 {
		return TemplateLayout{}, fmt.Errorf("incorrect number of pictures for 4-pic layout")
	}

	minAllowedHeight := 2500.00
	minAllowedPictureWidth := 1666.67
	spacing := e.imageSpacing
	AW := e.availableWidth

	layout := TemplateLayout{
		Positions:  make([][]float64, 4),
		Dimensions: make([][]float64, 4),
	}

	// Get Aspect Ratios (W/H)
	ARs := make([]float64, 4)
	for i, pic := range pictures {
		if pic.Height > 0 && pic.Width > 0 {
			ARs[i] = float64(pic.Width) / float64(pic.Height)
		} else {
			ARs[i] = 1.0 // Default AR
		}
	}
	AR0, AR1, AR2, AR3 := ARs[0], ARs[1], ARs[2], ARs[3]
	var W0, W1, W2, W3 float64 // Declare variables here

	// --- Geometric Calculation (Simplified: Allocate width based on height needs) ---
	// Calculate inverse aspect ratios (H/W) for height calculation
	invAR0 := 0.0
	if AR0 > 1e-6 {
		invAR0 = 1.0 / AR0
	}
	invAR1 := 0.0
	if AR1 > 1e-6 {
		invAR1 = 1.0 / AR1
	}
	invAR2 := 0.0
	if AR2 > 1e-6 {
		invAR2 = 1.0 / AR2
	}
	invAR3 := 0.0
	if AR3 > 1e-6 {
		invAR3 = 1.0 / AR3
	}

	// Calculate total height factor for each column
	heightFactorLeft := invAR0 + invAR2
	heightFactorRight := invAR1 + invAR3

	// Available width for the two columns (excluding middle spacing)
	columnsAvailableWidth := AW - spacing

	W_left := 0.0
	W_right := 0.0
	totalHeightFactor := heightFactorLeft + heightFactorRight

	if columnsAvailableWidth <= 1e-6 || totalHeightFactor <= 1e-6 {
		// Cannot calculate geometry factors
		return TemplateLayout{}, fmt.Errorf("cannot calculate 2x2 geometry factors (zero width or height factor sum)")
	}
	// Allocate width proportionally inverse to the height factor
	// More height needed -> less width allocated, aiming for equal column heights
	W_left = columnsAvailableWidth * (heightFactorRight / totalHeightFactor)
	W_right = columnsAvailableWidth * (heightFactorLeft / totalHeightFactor)

	// Check if calculated widths are valid
	if W_left <= 0 || W_right <= 0 {
		// Calculated widths are invalid
		return TemplateLayout{}, fmt.Errorf("calculated 2x2 column widths non-positive (W_L:%.2f, W_R:%.2f)", W_left, W_right)
	}

	// Calculate heights based on allocated widths and individual ARs
	H0 := W_left * invAR0
	H1 := W_right * invAR1
	H2 := W_left * invAR2
	H3 := W_right * invAR3

	// Determine row heights (use the max height within the conceptual row)
	H_top := math.Max(H0, H1)
	H_bottom := math.Max(H2, H3)

	// Recalculate widths based on the determined row heights to maintain aspect ratio
	// and ensure they fit within the allocated W_left/W_right
	// Use '=' instead of ':=' here since variables are already declared
	W0 = H_top * AR0
	W2 = H_bottom * AR2

	// Scale widths (and heights proportionally) if they overflow the column width
	if W0 > W_left {
		scale := W_left / W0
		W0 = W_left
		H_top = H_top * scale // Keep H_top consistent for scaling check
	}
	// Re-check H1 based on potentially scaled H_top
	H1_check := W_right / AR1
	if H1_check > H_top {
		H_top = H1_check
	} // Adjust H_top if needed
	// Re-calc W0, W1 based on final H_top
	W0 = H_top * AR0
	W1 = H_top * AR1
	// Scale row 1 if total width now exceeds AW
	if W0+spacing+W1 > AW+1e-3 { // Allow small float tolerance
		scale := AW / (W0 + spacing + W1)
		H_top *= scale
		W0 *= scale
		W1 *= scale
	}

	// Repeat for bottom row
	if W2 > W_left {
		scale := W_left / W2
		W2 = W_left
		H_bottom = H_bottom * scale
	}
	H3_check := W_right / AR3
	if H3_check > H_bottom {
		H_bottom = H3_check
	}
	W2 = H_bottom * AR2
	W3 = H_bottom * AR3
	if W2+spacing+W3 > AW+1e-3 {
		scale := AW / (W2 + spacing + W3)
		H_bottom *= scale
		W2 *= scale
		W3 *= scale
	}

	// --- Check Minimums (Rule 3.8) ---
	meetsMin := true
	if H_top < minAllowedHeight || H_bottom < minAllowedHeight {
		meetsMin = false
	}
	if W0 < minAllowedPictureWidth || W1 < minAllowedPictureWidth || W2 < minAllowedPictureWidth || W3 < minAllowedPictureWidth {
		meetsMin = false
	}

	if !meetsMin {
		fmt.Println("Warning: 4-picture layout (2x2 Geo) violates minimum dimensions.")
		// Return calculated layout anyway
	}

	// --- Populate Layout Struct ---
	layout.TotalHeight = H_top + spacing + H_bottom
	layout.TotalWidth = AW // Should fill the width

	// Pic 0 (Top Left)
	layout.Positions[0] = []float64{0, 0}
	layout.Dimensions[0] = []float64{W0, H_top}

	// Pic 1 (Top Right)
	layout.Positions[1] = []float64{W0 + spacing, 0}
	layout.Dimensions[1] = []float64{W1, H_top}

	// Pic 2 (Bottom Left)
	layout.Positions[2] = []float64{0, H_top + spacing}
	layout.Dimensions[2] = []float64{W2, H_bottom}

	// Pic 3 (Bottom Right)
	layout.Positions[3] = []float64{W2 + spacing, H_top + spacing} // X depends on W2
	layout.Dimensions[3] = []float64{W3, H_bottom}

	return layout, nil
}
