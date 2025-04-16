package calculate

import "fmt"

// processLayoutForSevenPictures handles the layout calculation and placement for exactly 7 pictures.
// It includes logic for handling split_required (to 3+4) and force_new_page signals.
func (e *ContinuousLayoutEngine) processLayoutForSevenPictures(pictures []Picture, layoutAvailableHeight float64) float64 {
	numPics := 7
	if len(pictures) != numPics {
		fmt.Printf("Error: processLayoutForSevenPictures called with %d pictures. Skipping.\n", len(pictures))
		return 0
	}

	layoutInfo, err := e.calculateSevenPicturesLayout(pictures, layoutAvailableHeight)

	// --- Special handling for 7-pic ---
	if err != nil {
		switch err.Error() {
		case "force_new_page":
			fmt.Println("Info: Forcing new page for 7-picture layout due to wide/tall images not fitting.")
			e.newPage()
			newLayoutAvailableHeight := (e.marginTop + e.availableHeight) - e.currentY // Recalculate for new page
			// Retry calculation on the new page
			fmt.Printf("Debug (Case 7 - Retry): Attempting calculation on new page. Available height: %.2f\n", newLayoutAvailableHeight)
			layoutInfo, err = e.calculateSevenPicturesLayout(pictures, newLayoutAvailableHeight)
			if err != nil {
				fmt.Printf("Error calculating 7-pic layout even on new page: %v. Skipping.\n", err)
				return 0
			}
			// If retry succeeds, err is nil, flow continues below

		case "split_required":
			fmt.Println("Info: Splitting 7-picture layout into 3+4.")
			// Estimate minimum height for the first group (3 pictures)
			minHeightForThree := GetRequiredMinHeight(e, "landscape", 3) // Use base landscape as estimate
			if layoutAvailableHeight < minHeightForThree {
				fmt.Printf("Error (Case 7 - Split 3+4): Not enough space (%.2f) for first three pictures (min %.2f est.). Skipping.\n", layoutAvailableHeight, minHeightForThree)
				return 0
			}

			// Place first three pictures
			fmt.Printf("Debug (Case 7 - Split 3+4): Placing first 3 (pics %d-%d). Available height: %.2f\n", pictures[0].Index, pictures[2].Index, layoutAvailableHeight)
			heightUsed1 := e.processLayoutForThreePictures(pictures[0:3], layoutAvailableHeight)
			if heightUsed1 <= 1e-6 {
				fmt.Println("Error: Failed to place first three pictures during 7-pic split. Skipping rest.")
				return 0
			}
			// The called function already updated e.currentY

			// Start a new page for the next four
			e.newPage()
			newLayoutAvailableHeight := (e.marginTop + e.availableHeight) - e.currentY

			// Place last four pictures on the new page
			fmt.Printf("Debug (Case 7 - Split 3+4): Placing last 4 (pics %d-%d) on new page. Available height: %.2f\n", pictures[3].Index, pictures[6].Index, newLayoutAvailableHeight)
			heightUsed2 := e.processLayoutForFourPictures(pictures[3:7], newLayoutAvailableHeight) // Use the 4-pic specific function
			if heightUsed2 <= 1e-6 {                                                               // Check if second placement failed
				fmt.Println("Error: Failed to place last four pictures during 7-pic split. First part placed, but operation incomplete.")
				return 0
			}
			// If second placement succeeds, return the height used by the second part on the new page.
			return heightUsed2

		default:
			fmt.Printf("Error calculating layout for 7 pictures: %v. Skipping placement.\n", err)
			return 0
		}
	}

	// If we reach here, err is nil, and layoutInfo is valid for placement.
	// +++ Log Before Placing +++
	fmt.Printf("Debug (Place - 7 Pics): Placing %d pictures. CurrentY before placement: %.2f. Layout TotalHeight: %.2f\n", len(pictures), e.currentY, layoutInfo.TotalHeight)
	// +++ End Log +++
	e.placePicturesInTemplate(pictures, layoutInfo)

	return layoutInfo.TotalHeight
}
