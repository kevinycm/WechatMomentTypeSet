package calculate

import "fmt"

// processLayoutForFourPictures handles the layout calculation and placement for exactly 4 pictures.
// It includes logic for handling split_required (to 2+2) and force_new_page signals.
func (e *ContinuousLayoutEngine) processLayoutForFourPictures(pictures []Picture, layoutAvailableHeight float64) float64 {
	numPics := 4
	if len(pictures) != numPics {
		fmt.Printf("Error: processLayoutForFourPictures called with %d pictures. Skipping.\n", len(pictures))
		return 0
	}

	layoutInfo, err := e.calculateFourPicturesLayout(pictures, layoutAvailableHeight)

	// --- Special handling for 4-pic ---
	if err != nil {
		switch err.Error() {
		case "force_new_page":
			fmt.Println("Info: Forcing new page for 4-picture layout due to wide/tall images not fitting.")
			e.newPage()
			newLayoutAvailableHeight := (e.marginTop + e.availableHeight) - e.currentY // Recalculate for new page
			// Retry calculation on the new page
			fmt.Printf("Debug (Case 4 - Retry): Attempting calculation on new page. Available height: %.2f\n", newLayoutAvailableHeight)
			layoutInfo, err = e.calculateFourPicturesLayout(pictures, newLayoutAvailableHeight)
			if err != nil {
				// If it still fails (e.g., split required even on a full page, or other calc error)
				fmt.Printf("Error calculating 4-pic layout even on new page: %v. Skipping.\n", err)
				return 0
			}
			// If calculation on new page succeeded, err is now nil,
			// and flow will continue to the common placement logic below.

		case "split_required":
			fmt.Println("Info: Splitting 4-picture layout into 2+2.")
			// Estimate minimum height for the first group (2 pictures)
			minHeightForTwo := GetRequiredMinHeight(e, "landscape", 2) // Use base landscape as estimate
			if layoutAvailableHeight < minHeightForTwo {
				fmt.Printf("Error (Case 4 - Split 2+2): Not enough space (%.2f) for first two pictures (min %.2f est.). Skipping.\n", layoutAvailableHeight, minHeightForTwo)
				return 0 // Cannot place even the first part
			}

			// Place first two pictures
			fmt.Printf("Debug (Case 4 - Split 2+2): Placing first 2 (pics %d-%d). Available height: %.2f\n", pictures[0].Index, pictures[1].Index, layoutAvailableHeight)
			heightUsed1 := e.processTwoPicturesLayoutAndPlace(pictures[0:2], layoutAvailableHeight)
			if heightUsed1 <= 1e-6 { // Check if first placement failed
				fmt.Println("Error: Failed to place first two pictures during 4-pic split. Skipping rest.")
				return 0
			}
			// Update Y coordinate *after* successful placement of the first part
			// Note: processTwoPicturesLayoutAndPlace should ideally update e.currentY itself, but review if needed.
			// Assuming it does, we just need to proceed.

			// Start a new page for the next two
			e.newPage()
			newLayoutAvailableHeight := (e.marginTop + e.availableHeight) - e.currentY

			// Place last two pictures on the new page
			fmt.Printf("Debug (Case 4 - Split 2+2): Placing last 2 (pics %d-%d) on new page. Available height: %.2f\n", pictures[2].Index, pictures[3].Index, newLayoutAvailableHeight)
			heightUsed2 := e.processTwoPicturesLayoutAndPlace(pictures[2:4], newLayoutAvailableHeight)
			if heightUsed2 <= 1e-6 { // Check if second placement failed
				fmt.Println("Error: Failed to place last two pictures during 4-pic split. First part placed, but operation incomplete.")
				return 0
			}
			// If second placement succeeds, return the height used by the second part on the new page.
			return heightUsed2

		default:
			// Handle other calculation errors (not split_required or force_new_page)
			fmt.Printf("Error calculating layout for 4 pictures: %v. Skipping placement.\n", err)
			return 0 // Return immediately as we cannot place
		}
	}

	// If we reach here, err is nil, and layoutInfo is valid for placement.
	// This happens if:
	// 1. Initial calculation succeeded without error.
	// 2. force_new_page occurred, but the retry on the new page succeeded.
	// +++ Log Before Placing +++
	fmt.Printf("Debug (Place - 4 Pics): Placing %d pictures. CurrentY before placement: %.2f. Layout TotalHeight: %.2f\n", len(pictures), e.currentY, layoutInfo.TotalHeight)
	// +++ End Log +++
	e.placePicturesInTemplate(pictures, layoutInfo)

	return layoutInfo.TotalHeight
}
