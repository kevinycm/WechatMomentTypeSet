package calculate

import (
	"fmt"
	"math"
)

// --- Actual implementation for the refactored function ---
// processSinglePictureLayoutAndPlace calculates and places a single picture
// based on its type (wide, tall, landscape, portrait) and the available height.
// Pagination decisions (using min heights) are made *before* calling this function.
// This function focuses on fitting the image into the given layoutAvailableHeight.
func (e *ContinuousLayoutEngine) processSinglePictureLayoutAndPlace(picture Picture, layoutAvailableHeight float64) float64 {
	// Calculate AR
	aspectRatio := 1.0
	validAR := false
	if picture.Height > 0 && picture.Width > 0 {
		aspectRatio = float64(picture.Width) / float64(picture.Height)
		validAR = true
	} else {
		fmt.Printf("Warning: Invalid dimensions for picture index %d. Using default AR=1.\n", picture.Index)
		// Attempt to use a default size based on generic min height?
	}

	picType := GetPictureType(aspectRatio)

	finalWidth := 0.0
	finalHeight := 0.0

	// Handle invalid AR or zero available height with a fallback
	if !validAR || layoutAvailableHeight <= 1e-6 {
		fmt.Printf("Warning: Using fallback dimensions for Pic %d due to invalid AR or no available height.\n", picture.Index)
		// Use min landscape height for 1 picture as fallback, capped by available height
		finalHeight = math.Min(e.minLandscapeHeights[1], layoutAvailableHeight)
		if finalHeight < 1.0 {
			finalHeight = 1.0
		} // Ensure positive height
		if validAR {
			finalWidth = finalHeight * aspectRatio // Try to respect AR
		} else {
			finalWidth = e.availableWidth // Default width if no AR
		}
		// Clamp width to available width
		if finalWidth > e.availableWidth {
			finalWidth = e.availableWidth
			if validAR {
				finalHeight = finalWidth / aspectRatio // Recalc height if AR exists
			}
		}
		// Ensure final fallback height is not more than originally available
		if finalHeight > layoutAvailableHeight {
			finalHeight = layoutAvailableHeight
		}
		if finalHeight < 1.0 {
			finalHeight = 1.0
		} // Ensure positive height
	} else {
		// Calculate dimensions based on picture type and available space
		switch picType {
		case "wide", "landscape":
			// Rule 1.1, 1.3: Fill width first, check height.
			finalWidth = e.availableWidth
			finalHeight = finalWidth / aspectRatio
			// If calculated height exceeds available, scale down to fit height.
			if finalHeight > layoutAvailableHeight {
				finalHeight = layoutAvailableHeight
				finalWidth = finalHeight * aspectRatio
			}
		case "tall", "portrait", "square", "unknown": // Treat square/unknown like portrait/tall
			// Rule 1.2, 1.4: Fill height first, check width.
			finalHeight = layoutAvailableHeight
			finalWidth = finalHeight * aspectRatio
			// If calculated width exceeds available, scale down to fit width.
			if finalWidth > e.availableWidth {
				finalWidth = e.availableWidth
				finalHeight = finalWidth / aspectRatio
			}
		}

		// Final check: Ensure dimensions are positive after calculations
		if finalWidth < 1.0 {
			finalWidth = 1.0
		}
		if finalHeight < 1.0 {
			finalHeight = 1.0
		}
	}

	// --- Place Picture (Horizontally Centered) ---
	e.placeSinglePicture(picture, finalWidth, finalHeight)

	return finalHeight // Return the actual height used for placement
}

// placeSinglePicture places a single picture, centering it horizontally.
func (e *ContinuousLayoutEngine) placeSinglePicture(pic Picture, width, height float64) {
	// Calculate horizontal centering
	startX := 0.0
	if width < e.availableWidth {
		startX = (e.availableWidth - width) / 2
	}

	// Ensure entry exists
	if len(e.currentPage.Entries) == 0 {
		fmt.Println("Warning: Placing single picture but no entry exists on current page. Creating one.")
		e.currentPage.Entries = append(e.currentPage.Entries, PageEntry{})
	}
	currentEntry := &e.currentPage.Entries[len(e.currentPage.Entries)-1]

	// Calculate absolute coordinates
	// Remember: e.currentY is the starting Y for this element, before adding its height.
	absX0 := e.marginLeft + startX
	absY0 := e.currentY
	absX1 := absX0 + width
	absY1 := absY0 + height

	// Ensure coordinates are valid floats
	if math.IsNaN(absX0) || math.IsInf(absX0, 0) ||
		math.IsNaN(absY0) || math.IsInf(absY0, 0) ||
		math.IsNaN(absX1) || math.IsInf(absX1, 0) ||
		math.IsNaN(absY1) || math.IsInf(absY1, 0) {
		fmt.Printf("Error: Invalid coordinates calculated for picture %d: [%.2f, %.2f], [%.2f, %.2f]\n",
			pic.Index, absX0, absY0, absX1, absY1)
		// Skip appending the picture if coordinates are invalid
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
		Width:  int(math.Round(width)),  // Store final layout width (rounded)
		Height: int(math.Round(height)), // Store final layout height (rounded)
	})
	// currentY updated by the caller (processSinglePictureLayoutAndPlace) using the returned height
}
