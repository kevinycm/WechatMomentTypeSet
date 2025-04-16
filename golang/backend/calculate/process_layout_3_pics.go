package calculate

import "fmt"

// processLayoutForThreePictures handles the layout calculation and placement for exactly 3 pictures.
// It determines the best layout based on available height and minimum requirements.
// Note: The logic for 3 pictures doesn't inherently involve splitting further,
// so it mainly calculates, places, and returns the height used, or 0 on error.
func (e *ContinuousLayoutEngine) processLayoutForThreePictures(pictures []Picture, layoutAvailableHeight float64) float64 {
	numPics := 3
	if len(pictures) != numPics {
		fmt.Printf("Error: processLayoutForThreePictures called with %d pictures. Skipping.\n", len(pictures))
		return 0
	}

	// Calculate the best layout for 3 pictures given the available height.
	layoutInfo, err := e.calculateThreePicturesLayout(pictures, layoutAvailableHeight)
	if err != nil {
		// Check if the error indicates an inability to find *any* valid layout, even fallback.
		// calculateThreePicturesLayout currently returns errors in such cases.
		fmt.Printf("Error calculating layout for 3 pictures: %v. Skipping placement.\n", err)
		// Treat inability to calculate a valid layout as failure for placement.
		return 0 // Return 0 height used on calculation failure
	}

	// If calculation succeeds (even if it's a fallback layout), place it.
	// No need to check for split_required or force_new_page as 3 is the base case.
	// +++ Log Before Placing +++
	fmt.Printf("Debug (Place - 3 Pics): Placing %d pictures. CurrentY before placement: %.2f. Layout TotalHeight: %.2f\n", len(pictures), e.currentY, layoutInfo.TotalHeight)
	// +++ End Log +++
	e.placePicturesInTemplate(pictures, layoutInfo)

	// Return the actual height consumed by the placed layout.
	// The caller (processPictures) will use the engine's updated e.currentY.
	return layoutInfo.TotalHeight
}
