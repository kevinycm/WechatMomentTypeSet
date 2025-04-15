package backend

import (
	"fmt"
	// "math"
)

// calculateEightPicturesLayout calculates dimensions for 8 pictures (4 Top, 4 Bottom).
func (e *ContinuousLayoutEngine) calculateEightPicturesLayout(pictures []Picture) (TemplateLayout, error) {
	if len(pictures) != 8 {
		return TemplateLayout{}, fmt.Errorf("incorrect number of pictures for 8-pic layout")
	}

	// Define minimums (Rule 3.8)
	minAllowedHeight := 2500.00
	minAllowedPictureWidth := 1666.67
	spacing := e.imageSpacing
	AW := e.availableWidth

	layout := TemplateLayout{
		Positions:  make([][]float64, 8),
		Dimensions: make([][]float64, 8),
	}

	row1Pics := pictures[0:4]
	row2Pics := pictures[4:8]

	widths1, _, height1 := e.calculateUniformRowHeightLayout(row1Pics, AW)
	widths2, _, height2 := e.calculateUniformRowHeightLayout(row2Pics, AW)

	if height1 <= 1e-6 || height2 <= 1e-6 {
		return layout, fmt.Errorf("failed to calculate row layouts for 4T4B (8 pics)")
	}
	W0, W1, W2, W3 := widths1[0], widths1[1], widths1[2], widths1[3]
	W4, W5, W6, W7 := widths2[0], widths2[1], widths2[2], widths2[3]

	// --- Check Minimums ---*
	meetsMin := true
	if height1 < minAllowedHeight || height2 < minAllowedHeight {
		meetsMin = false
	}
	// Combine widths for easier checking
	widths := append(widths1, widths2...)
	for _, w := range widths {
		if w < minAllowedPictureWidth {
			meetsMin = false
			break
		}
	}

	if !meetsMin {
		// Log warning, but still return the calculated layout as fallback
		fmt.Println("Warning: 8-picture layout (4T4B) violates minimum dimensions.")
	}

	// --- Populate Layout Struct ---*
	layout.TotalHeight = height1 + spacing + height2
	layout.TotalWidth = AW
	// Row 1
	currentX := 0.0
	layout.Positions[0] = []float64{currentX, 0}
	layout.Dimensions[0] = []float64{W0, height1}
	currentX += W0 + spacing
	layout.Positions[1] = []float64{currentX, 0}
	layout.Dimensions[1] = []float64{W1, height1}
	currentX += W1 + spacing
	layout.Positions[2] = []float64{currentX, 0}
	layout.Dimensions[2] = []float64{W2, height1}
	currentX += W2 + spacing
	layout.Positions[3] = []float64{currentX, 0}
	layout.Dimensions[3] = []float64{W3, height1}
	// Row 2
	currentX = 0.0
	layout.Positions[4] = []float64{currentX, height1 + spacing}
	layout.Dimensions[4] = []float64{W4, height2}
	currentX += W4 + spacing
	layout.Positions[5] = []float64{currentX, height1 + spacing}
	layout.Dimensions[5] = []float64{W5, height2}
	currentX += W5 + spacing
	layout.Positions[6] = []float64{currentX, height1 + spacing}
	layout.Dimensions[6] = []float64{W6, height2}
	currentX += W6 + spacing
	layout.Positions[7] = []float64{currentX, height1 + spacing}
	layout.Dimensions[7] = []float64{W7, height2}

	return layout, nil // Return calculated layout (even if fallback)
}
