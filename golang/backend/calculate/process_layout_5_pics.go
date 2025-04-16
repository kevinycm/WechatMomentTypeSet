package calculate

import "fmt"

// processLayoutForFivePictures handles the layout calculation and placement for exactly 5 pictures.
// It includes logic for handling split_required (to 3+2) and force_new_page signals.
func (e *ContinuousLayoutEngine) processLayoutForFivePictures(pictures []Picture, layoutAvailableHeight float64) float64 {
	numPics := 5
	if len(pictures) != numPics {
		fmt.Printf("Error: processLayoutForFivePictures called with %d pictures. Skipping.\n", len(pictures))
		return 0
	}

	layoutInfo, err := e.calculateFivePicturesLayout(pictures, layoutAvailableHeight)

	// --- Special handling for 5-pic ---
	if err != nil {
		switch err.Error() {
		case "force_new_page":
			fmt.Println("Info: Forcing new page for 5-picture layout due to wide/tall images not fitting.")
			e.newPage()
			newLayoutAvailableHeight := (e.marginTop + e.availableHeight) - e.currentY // Recalculate for new page
			// Retry calculation on the new page
			fmt.Printf("Debug (Case 5 - Retry): Attempting calculation on new page. Available height: %.2f\n", newLayoutAvailableHeight)
			layoutInfo, err = e.calculateFivePicturesLayout(pictures, newLayoutAvailableHeight)
			if err != nil {
				fmt.Printf("Error calculating 5-pic layout even on new page: %v. Skipping.\n", err)
				return 0
			}
			// If retry succeeds, err is nil, flow continues below

		case "split_required":
			fmt.Println("Info: Splitting 5-picture layout into 3+2.")
			// Estimate minimum height for the first group (3 pictures)
			minHeightForThree := GetRequiredMinHeight(e, "landscape", 3) // Use base landscape as estimate
			if layoutAvailableHeight < minHeightForThree {
				fmt.Printf("Error (Case 5 - Split 3+2): Not enough space (%.2f) for first three pictures (min %.2f est.). Skipping.\n", layoutAvailableHeight, minHeightForThree)
				return 0
			}

			// Place first three pictures - Use the specific 3-pic processing function
			fmt.Printf("Debug (Case 5 - Split 3+2): Placing first 3 (pics %d-%d). Available height: %.2f\n", pictures[0].Index, pictures[2].Index, layoutAvailableHeight)
			// IMPORTANT: Call the new specific function, not the old general one.
			heightUsed1 := e.processLayoutForThreePictures(pictures[0:3], layoutAvailableHeight)
			if heightUsed1 <= 1e-6 { // Check if first placement failed
				fmt.Println("Error: Failed to place first three pictures during 5-pic split. Skipping rest.")
				return 0
			}
			// The called function already updated e.currentY

			// Start a new page for the next two
			e.newPage()
			newLayoutAvailableHeight := (e.marginTop + e.availableHeight) - e.currentY

			// Place last two pictures on the new page
			fmt.Printf("Debug (Case 5 - Split 3+2): Placing last 2 (pics %d-%d) on new page. Available height: %.2f\n", pictures[3].Index, pictures[4].Index, newLayoutAvailableHeight)
			heightUsed2 := e.processTwoPicturesLayoutAndPlace(pictures[3:5], newLayoutAvailableHeight) // Assuming this function exists and handles 2 pics correctly
			if heightUsed2 <= 1e-6 {                                                                   // Check if second placement failed
				fmt.Println("Error: Failed to place last two pictures during 5-pic split. First part placed, but operation incomplete.")
				return 0
			}
			// If second placement succeeds, return the height used by the second part on the new page.
			return heightUsed2

		default:
			fmt.Printf("Error calculating layout for 5 pictures: %v. Skipping placement.\n", err)
			return 0
		}
	}

	// If we reach here, err is nil, and layoutInfo is valid for placement.
	// +++ Log Before Placing +++
	fmt.Printf("Debug (Place - 5 Pics): Placing %d pictures. CurrentY before placement: %.2f. Layout TotalHeight: %.2f\n", len(pictures), e.currentY, layoutInfo.TotalHeight)
	// +++ End Log +++
	e.placePicturesInTemplate(pictures, layoutInfo)

	return layoutInfo.TotalHeight
}
