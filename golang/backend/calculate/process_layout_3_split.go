package calculate

import (
	"fmt"
)

// processThreePicturesWithSplitLogic handles layout for 3 pictures,
// attempting 1+2 split if all 3 don't fit initially.
func (e *ContinuousLayoutEngine) processThreePicturesWithSplitLogic(pictures []Picture, layoutAvailableHeight float64) float64 {
	numPics := 3
	if len(pictures) != numPics {
		fmt.Printf("Error (process3Split): Incorrect number of pictures: %d\n", len(pictures))
		return 0
	}

	// --- Attempt 1: Try placing all 3 on the current page ---
	fmt.Println("Debug (process3Split): Attempting to place all 3 pictures initially.")
	layoutInfo3, err3 := e.calculateThreePicturesLayout(pictures, layoutAvailableHeight)

	if err3 == nil && layoutInfo3.TotalHeight <= layoutAvailableHeight+1e-6 { // Success and fits
		fmt.Println("Debug (process3Split): All 3 fit on current page.")
		e.placePicturesInTemplate(pictures, layoutInfo3)
		// Don't update e.currentY here; return height for the caller (processPicturesOldStrategy)
		return layoutInfo3.TotalHeight
	}

	// --- Attempt 2: Try placing Pic 1 on current page, Pics 2+3 on new page ---
	fmt.Println("Debug (process3Split): All 3 failed/didn't fit. Attempting 1+2 split.")

	// Try calculating layout for just the first picture
	// Note: We use calculateRowLayout here as it's designed for single/double/triple rows
	// and handles the case where pic1 might need scaling itself.
	layoutInfo1, err1 := e.calculateRowLayout(pictures[0:1], "row-of-1", layoutAvailableHeight)

	if err1 == nil && layoutInfo1.TotalHeight <= layoutAvailableHeight+1e-6 { // Pic 1 fits
		fmt.Println("Debug (process3Split): Pic 1 fits on current page. Placing it.")
		e.placePicturesInRow(pictures[0:1], layoutInfo1)
		heightUsed1 := layoutInfo1.TotalHeight
		e.currentY += heightUsed1 // Update Y coordinate *after* placing pic 1

		// Now, new page for pics 2 and 3
		fmt.Println("Debug (process3Split): Creating new page for Pics 2 & 3.")
		e.newPage()

		newAvailableHeight := (e.marginTop + e.availableHeight) - e.currentY
		fmt.Printf("Debug (process3Split): Placing Pics 2 & 3 on new page (Page %d). Available H: %.2f\n", e.currentPage.Page, newAvailableHeight)

		// Calculate and place pics 2 & 3 as a row-of-2
		layoutInfo2, err2 := e.calculateRowLayout(pictures[1:3], "row-of-2", newAvailableHeight)
		if err2 != nil {
			fmt.Printf("Error (process3Split): Failed to calculate layout for Pics 2 & 3 on new page: %v. Aborting split.\n", err2)
			return 0
		}

		if layoutInfo2.TotalHeight > newAvailableHeight+1e-6 {
			fmt.Printf("Error (process3Split): Calculated height (%.2f) for Pics 2 & 3 exceeds available height (%.2f) on new page. Aborting split.\n", layoutInfo2.TotalHeight, newAvailableHeight)
			return 0
		}

		e.placePicturesInRow(pictures[1:3], layoutInfo2)
		// Return height used on the *last* page
		return layoutInfo2.TotalHeight

	} else {
		// --- Attempt 3: Pic 1 didn't fit either. Place all 3 on a new page ---
		fmt.Println("Debug (process3Split): Pic 1 also failed/didn't fit. Placing all 3 on new page.")
		e.newPage()
		newAvailableHeight := (e.marginTop + e.availableHeight) - e.currentY
		fmt.Printf("Debug (process3Split): Placing all 3 on new page (Page %d). Available H: %.2f\n", e.currentPage.Page, newAvailableHeight)

		// Retry calculation for all 3 on the new page
		layoutInfo3Retry, err3Retry := e.calculateThreePicturesLayout(pictures, newAvailableHeight)
		if err3Retry != nil {
			fmt.Printf("Error (process3Split): Failed to calculate layout for all 3 pics even on new page: %v\n", err3Retry)
			return 0
		}
		if layoutInfo3Retry.TotalHeight > newAvailableHeight+1e-6 {
			fmt.Printf("Error (process3Split): Calculated height (%.2f) for 3 pics exceeds available height (%.2f) on new page.\n", layoutInfo3Retry.TotalHeight, newAvailableHeight)
			return 0
		}

		e.placePicturesInTemplate(pictures, layoutInfo3Retry)
		// Return height used on the new page
		return layoutInfo3Retry.TotalHeight
	}
}
