package calculate

import (
	"fmt"
	"math"
)

// --- Actual implementation for the refactored function ---
// processTwoPicturesLayoutAndPlace calculates and places two pictures based on their types
// and the rules provided (Up/Down or Left/Right).
// Pagination is handled before this function. It fits the layout into layoutAvailableHeight.
func (e *ContinuousLayoutEngine) processTwoPicturesLayoutAndPlace(pictures []Picture, layoutAvailableHeight float64) float64 {
	if len(pictures) != 2 {
		fmt.Printf("Error: processTwoPicturesLayoutAndPlace called with %d pictures. Skipping.\n", len(pictures))
		return 0
	}

	pic1, pic2 := pictures[0], pictures[1]
	ar1, ar2 := 1.0, 1.0 // Default ARs
	validAR1, validAR2 := false, false
	if pic1.Height > 0 && pic1.Width > 0 {
		ar1 = float64(pic1.Width) / float64(pic1.Height)
		validAR1 = true
	}
	if pic2.Height > 0 && pic2.Width > 0 {
		ar2 = float64(pic2.Width) / float64(pic2.Height)
		validAR2 = true
	}

	// Fallback if any AR is invalid or available height is zero
	if !validAR1 || !validAR2 || layoutAvailableHeight <= 1e-6 {
		fmt.Printf("Warning: Using fallback layout for 2 pictures due to invalid AR or no available height.\n")
		// Fallback to simple left/right layout with min height, scaled if needed
		widths, _, rowHeight := e.calculateUniformRowHeightLayout(pictures, e.availableWidth)

		// Determine the stricter (smaller) minimum height required for 2 pics as fallback threshold
		minRequiredFallbackHeight := math.Min(e.minLandscapeHeights[2], e.minPortraitHeights[2])

		if rowHeight < minRequiredFallbackHeight { // Ensure minimum height in fallback
			rowHeight = minRequiredFallbackHeight
			// Recalculate widths based on forced min height (might exceed available width)
			// Note: This assumes ar1 and ar2 have default values (1.0) if original ARs were invalid.
			widths[0] = rowHeight * ar1
			widths[1] = rowHeight * ar2
			// This fallback doesn't rescale width perfectly, it prioritizes min height.
		}
		finalHeight := rowHeight
		scale := 1.0
		if finalHeight > layoutAvailableHeight {
			scale = layoutAvailableHeight / finalHeight
			finalHeight *= scale
			widths[0] *= scale
			widths[1] *= scale
		}
		e.placePictureRow(pictures, widths, finalHeight)
		return finalHeight
	}

	type1 := GetPictureType(ar1)
	type2 := GetPictureType(ar2)

	finalTotalHeight := 0.0

	// Determine Layout Type based on rules 2.1 - 2.9
	layoutType := ""
	switch {
	// Up/Down Rules (2.1 - 2.4)
	case (type1 == "wide" && type2 == "tall"), (type1 == "tall" && type2 == "wide"):
		layoutType = "up_down"
	case (type1 == "wide" && type2 == "landscape"), (type1 == "landscape" && type2 == "wide"):
		layoutType = "up_down"
	case (type1 == "wide" && type2 == "portrait"), (type1 == "portrait" && type2 == "wide"):
		layoutType = "up_down"
	case (type1 == "tall" && type2 == "landscape"), (type1 == "landscape" && type2 == "tall"):
		layoutType = "up_down"

	// Left/Right Rules (2.5 - 2.9) - Default to Left/Right for other combos
	default:
		layoutType = "left_right"
	}

	switch layoutType {
	case "up_down":
		// Calculate initial heights based on fitting available width
		width1 := e.availableWidth
		height1 := width1 / ar1
		width2 := e.availableWidth
		height2 := width2 / ar2

		totalRequiredHeight := height1 + height2 + e.imageSpacing
		scale := 1.0

		if totalRequiredHeight > layoutAvailableHeight {
			scale = layoutAvailableHeight / totalRequiredHeight
		}

		finalHeight1 := height1 * scale
		finalWidth1 := width1 * scale
		finalHeight2 := height2 * scale
		finalWidth2 := width2 * scale

		// Let's not scale spacing for now, just cap total height
		finalTotalHeight = finalHeight1 + finalHeight2 + e.imageSpacing
		if finalTotalHeight > layoutAvailableHeight {
			// Adjust spacing or heights slightly if rounding/scaling pushed over?
			// Simpler: Trust the scale calculation, the sum should be <= layoutAvailableHeight
			finalTotalHeight = layoutAvailableHeight // Cap at available height
		}

		// Place pictures vertically stacked and centered horizontally
		startY := e.currentY
		e.placeSinglePictureStacked(pic1, finalWidth1, finalHeight1, startY)
		e.placeSinglePictureStacked(pic2, finalWidth2, finalHeight2, startY+finalHeight1+e.imageSpacing)

	case "left_right":
		// Use existing uniform height layout logic
		widths, _, rowHeight := e.calculateUniformRowHeightLayout(pictures, e.availableWidth)
		finalHeight := rowHeight
		scale := 1.0

		if finalHeight > layoutAvailableHeight {
			scale = layoutAvailableHeight / finalHeight
			finalHeight *= scale
			for i := range widths {
				widths[i] *= scale
			}
		}
		finalTotalHeight = finalHeight
		e.placePictureRow(pictures, widths, finalHeight)
	}

	return finalTotalHeight
}

// placeSinglePictureStacked is a helper to place a picture at a specific Y offset
// within a stacked layout, centering it horizontally.
func (e *ContinuousLayoutEngine) placeSinglePictureStacked(pic Picture, width, height, startY float64) {
	// Calculate horizontal centering
	startX := 0.0
	if width < e.availableWidth {
		startX = (e.availableWidth - width) / 2
	}

	// Ensure entry exists
	if len(e.currentPage.Entries) == 0 {
		fmt.Println("Warning: Placing stacked picture but no entry exists on current page. Creating one.")
		e.currentPage.Entries = append(e.currentPage.Entries, PageEntry{})
	}
	currentEntry := &e.currentPage.Entries[len(e.currentPage.Entries)-1]

	// Calculate absolute coordinates using the provided startY
	absX0 := e.marginLeft + startX
	absY0 := startY // Use the provided startY directly
	absX1 := absX0 + width
	absY1 := absY0 + height

	// Ensure coordinates are valid floats
	if math.IsNaN(absX0) || math.IsInf(absX0, 0) ||
		math.IsNaN(absY0) || math.IsInf(absY0, 0) ||
		math.IsNaN(absX1) || math.IsInf(absX1, 0) ||
		math.IsNaN(absY1) || math.IsInf(absY1, 0) {
		fmt.Printf("Error: Invalid coordinates calculated for stacked picture %d: [%.2f, %.2f], [%.2f, %.2f]\n",
			pic.Index, absX0, absY0, absX1, absY1)
		return
	}

	area := [][]float64{
		{absX0, absY0},
		{absX1, absY1},
	}
	currentEntry.Pictures = append(currentEntry.Pictures, Picture{
		Index:  pic.Index,
		Area:   area,
		URL:    pic.URL,
		Width:  int(math.Round(width)),
		Height: int(math.Round(height)),
	})
	// e.currentY is managed by the calling function (processTwoPicturesLayoutAndPlace)
}

// placePictureRow places a single row of pictures, centering it horizontally.
func (e *ContinuousLayoutEngine) placePictureRow(pictures []Picture, widths []float64, rowHeight float64) {
	if len(pictures) == 0 || len(pictures) != len(widths) || rowHeight <= 1e-6 {
		return
	}

	// Calculate total width including spacing for centering (Rule 3.8 implied)
	totalImageWidth := 0.0
	for _, w := range widths {
		totalImageWidth += w
	}
	totalRowWidth := totalImageWidth + e.imageSpacing*float64(len(pictures)-1) // Use e.imageSpacing (Rule 3.11)
	startX := 0.0
	if totalRowWidth < e.availableWidth {
		startX = (e.availableWidth - totalRowWidth) / 2
	}

	// Add pictures to the current entry
	if len(e.currentPage.Entries) == 0 {
		// This should be handled before placePictureRow is called
		fmt.Println("Warning: Placing picture row but no entry exists on current page. Creating one.")
		e.currentPage.Entries = append(e.currentPage.Entries, PageEntry{})
	}
	currentEntry := &e.currentPage.Entries[len(e.currentPage.Entries)-1]
	currentX := startX

	for i, pic := range pictures {
		area := [][]float64{
			{e.marginLeft + currentX, e.currentY},
			{e.marginLeft + currentX + widths[i], e.currentY + rowHeight},
		}
		currentEntry.Pictures = append(currentEntry.Pictures, Picture{
			Index:  pic.Index,
			Area:   area,
			URL:    pic.URL,
			Width:  int(widths[i]), // Store final layout width
			Height: int(rowHeight), // Store final layout height
		})
		if i < len(pictures)-1 {
			currentX += widths[i] + e.imageSpacing // Add spacing between images (Rule 3.11)
		}
	}
	// currentY updated by caller (processTwoPicturesLayoutAndPlace)
}

// calculateUniformRowHeightLayout calculates dimensions for a row aiming for uniform height,
// fitting within the availableWidth. It returns the calculated widths for each picture,
// (unused heights array), and the final uniform row height.
func (e *ContinuousLayoutEngine) calculateUniformRowHeightLayout(pictures []Picture, availableWidth float64) ([]float64, []float64, float64) {
	if len(pictures) == 0 {
		return nil, nil, 0
	}

	numPics := len(pictures)
	aspectRatios := make([]float64, numPics)
	totalAspectRatioSum := 0.0
	validARCount := 0
	for i, pic := range pictures {
		if pic.Height > 0 && pic.Width > 0 {
			aspectRatios[i] = float64(pic.Width) / float64(pic.Height)
			totalAspectRatioSum += aspectRatios[i]
			validARCount++
		} else {
			aspectRatios[i] = 1.0 // Default AR for invalid data
			totalAspectRatioSum += 1.0
			validARCount++
		}
	}

	// Available width for images themselves (excluding spacing)
	rowAvailableWidth := availableWidth - e.imageSpacing*float64(numPics-1)
	if rowAvailableWidth < 1.0 {
		rowAvailableWidth = 1.0
	} // Avoid negative/zero width

	finalRowHeight := 0.0
	if validARCount > 0 && totalAspectRatioSum > 1e-6 {
		// Calculate the height H such that Sum(H * AR_i) = rowAvailableWidth
		finalRowHeight = rowAvailableWidth / totalAspectRatioSum
	} else if validARCount > 0 {
		// Handle case where all ARs are zero or invalid leading to zero sum
		finalRowHeight = math.Min(e.availableHeight, 2500.0) // Use min height or available page height
	} else {
		return nil, nil, 0 // No valid pictures
	}

	// Calculate final widths based on this uniform height
	finalRowWidths := make([]float64, numPics)
	calculatedTotalWidth := 0.0
	for i, ar := range aspectRatios {
		finalRowWidths[i] = finalRowHeight * ar
		calculatedTotalWidth += finalRowWidths[i]
	}

	// Sanity check: Cap height at available page height
	if finalRowHeight > e.availableHeight {
		// scale := e.availableHeight / finalRowHeight // Variable not used
		finalRowHeight = e.availableHeight
		calculatedTotalWidth = 0.0 // Recalculate widths
		for i, ar := range aspectRatios {
			finalRowWidths[i] = finalRowHeight * ar
			calculatedTotalWidth += finalRowWidths[i]
		}
	}

	// Final check: If total width still exceeds available width after height calc/capping,
	// scale down based on width constraint.
	if calculatedTotalWidth > rowAvailableWidth && calculatedTotalWidth > 1e-6 {
		scaleFactor := rowAvailableWidth / calculatedTotalWidth
		finalRowHeight *= scaleFactor
		for i := range finalRowWidths {
			finalRowWidths[i] *= scaleFactor
		}
	}

	// The function signature requires returning heights, but they are all `finalRowHeight`
	finalRowHeights := make([]float64, numPics)
	for i := range finalRowHeights {
		finalRowHeights[i] = finalRowHeight
	}

	return finalRowWidths, finalRowHeights, finalRowHeight
}
