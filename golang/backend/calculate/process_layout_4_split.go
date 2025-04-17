package calculate

import (
	"fmt"
)

// processFourPicturesWithSplitLogic handles layout for 4 pictures,
// attempting 2+2 split if all 4 don't fit initially.
func (e *ContinuousLayoutEngine) processFourPicturesWithSplitLogic(pictures []Picture, layoutAvailableHeight float64) float64 {
	numPics := 4
	if len(pictures) != numPics {
		fmt.Printf("Error (process4Split): Incorrect number of pictures: %d\n", len(pictures))
		return 0
	}

	// --- Attempt 1: Try placing all 4 on the current page ---
	fmt.Println("Debug (process4Split): Attempting to place all 4 pictures initially.")
	layoutInfo4, err4 := e.calculateFourPicturesLayout(pictures, layoutAvailableHeight)

	if err4 == nil && layoutInfo4.TotalHeight <= layoutAvailableHeight+1e-6 { // Success and fits
		fmt.Println("Debug (process4Split): All 4 fit on current page.")
		e.placePicturesInTemplate(pictures, layoutInfo4)
		return layoutInfo4.TotalHeight
	}

	// --- Attempt 2: Try placing Pics 1-2 on current page, Pics 3-4 on new page ---
	fmt.Printf("Debug (process4Split): All 4 failed/didn't fit (err: %v, H: %.2f > Avail: %.2f). Attempting 2+2 split.\n", err4, layoutInfo4.TotalHeight, layoutAvailableHeight)

	// Try calculating layout for just the first two pictures
	layoutInfo2_1, err2_1 := e.calculateRowLayout(pictures[0:2], "row-of-2", layoutAvailableHeight)

	if err2_1 == nil && layoutInfo2_1.TotalHeight <= layoutAvailableHeight+1e-6 { // Pics 1-2 fit
		fmt.Println("Debug (process4Split): Pics 1-2 fit on current page. Placing them.")
		e.placePicturesInRow(pictures[0:2], layoutInfo2_1)
		heightUsed1 := layoutInfo2_1.TotalHeight
		e.currentY += heightUsed1 // Update Y coordinate *after* placing pics 1-2

		// Now, new page for pics 3-4
		fmt.Println("Debug (process4Split): Creating new page for Pics 3-4.")
		e.newPage()

		newAvailableHeight := (e.marginTop + e.availableHeight) - e.currentY
		fmt.Printf("Debug (process4Split): Placing Pics 3-4 on new page (Page %d). Available H: %.2f\n", e.currentPage.Page, newAvailableHeight)

		// Calculate and place pics 3-4 as a row-of-2
		layoutInfo2_2, err2_2 := e.calculateRowLayout(pictures[2:4], "row-of-2", newAvailableHeight)
		if err2_2 != nil {
			fmt.Printf("Error (process4Split): Failed to calculate layout for Pics 3-4 on new page: %v. Aborting split.\n", err2_2)
			return 0
		}

		if layoutInfo2_2.TotalHeight > newAvailableHeight+1e-6 {
			fmt.Printf("Error (process4Split): Calculated height (%.2f) for Pics 3-4 exceeds available height (%.2f) on new page. Aborting split.\n", layoutInfo2_2.TotalHeight, newAvailableHeight)
			return 0
		}

		e.placePicturesInRow(pictures[2:4], layoutInfo2_2)
		// Return height used on the *last* page
		return layoutInfo2_2.TotalHeight

	} else {
		// --- Attempt 3: Pics 1-2 didn't fit either. Place all 4 on a new page ---
		fmt.Printf("Debug (process4Split): Pics 1-2 also failed/didn't fit (err: %v, H: %.2f > Avail: %.2f). Placing all 4 on new page.\n", err2_1, layoutInfo2_1.TotalHeight, layoutAvailableHeight)
		e.newPage()
		newAvailableHeight := (e.marginTop + e.availableHeight) - e.currentY
		fmt.Printf("Debug (process4Split): Placing all 4 on new page (Page %d). Available H: %.2f\n", e.currentPage.Page, newAvailableHeight)

		// Retry calculation for all 4 on the new page
		layoutInfo4Retry, err4Retry := e.calculateFourPicturesLayout(pictures, newAvailableHeight)
		if err4Retry != nil {
			fmt.Printf("Error (process4Split): Failed to calculate layout for all 4 pics even on new page: %v\n", err4Retry)
			return 0
		}
		if layoutInfo4Retry.TotalHeight > newAvailableHeight+1e-6 {
			fmt.Printf("Error (process4Split): Calculated height (%.2f) for 4 pics exceeds available height (%.2f) on new page.\n", layoutInfo4Retry.TotalHeight, newAvailableHeight)
			return 0
		}

		e.placePicturesInTemplate(pictures, layoutInfo4Retry)
		// Return height used on the new page
		return layoutInfo4Retry.TotalHeight
	}
}
