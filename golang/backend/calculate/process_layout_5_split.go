package calculate

import (
	"fmt"
)

// processFivePicturesWithSplitLogic handles layout for 5 pictures,
// attempting 2+3 split if all 5 don't fit initially.
// UPDATED: Uses imageSpacing for inter-group spacing.
func (e *ContinuousLayoutEngine) processFivePicturesWithSplitLogic(pictures []Picture, layoutAvailableHeight float64) float64 {
	numPics := 5
	if len(pictures) != numPics {
		fmt.Printf("Error (process5Split): Incorrect number of pictures: %d\n", len(pictures))
		return 0
	}

	// --- Attempt 1: Try placing all 5 on the current page ---
	fmt.Println("Debug (process5Split): Attempting to place all 5 pictures initially.")
	layoutInfo5, err5 := e.calculateFivePicturesLayout(pictures, layoutAvailableHeight)

	if err5 == nil && layoutInfo5.TotalHeight <= layoutAvailableHeight+1e-6 { // Success and fits
		fmt.Println("Debug (process5Split): All 5 fit on current page.")
		e.placePicturesInTemplate(pictures, layoutInfo5)
		return layoutInfo5.TotalHeight
	}

	// --- Attempt 2: Try placing Pics 1-2 on current page, Pics 3-5 on next (or same if space) ---
	fmt.Printf("Debug (process5Split): All 5 failed/didn't fit (err: %v, H: %.2f > Avail: %.2f). Attempting 2+3 split.\n", err5, layoutInfo5.TotalHeight, layoutAvailableHeight)

	// Try calculating layout for just the first two pictures
	// Use processTwoPicGroup to handle potential retry/pagination for the first group
	fmt.Println("Debug (process5Split): Processing Group 1 (Pics 0-1).")
	heightUsed1, errG1 := e.processTwoPicGroup(pictures[0:2], 1, layoutAvailableHeight)
	pageAfterG1 := e.currentPage.Page
	yAfterG1 := e.currentY

	if errG1 != nil {
		fmt.Printf("Error (process5Split): Failed processing group 1 (Pics 0-1): %v. Placing all 5 on new page as fallback.\n", errG1)
		// Fallback if first group fails
		return e.placeAllFiveOnNewPage(pictures)
	}
	if heightUsed1 <= 1e-6 {
		fmt.Println("Warning (process5Split): Group 1 (0-1) was skipped. Placing all 5 on new page as fallback.")
		// Fallback if first group skipped
		return e.placeAllFiveOnNewPage(pictures)
	}
	fmt.Printf("Debug (process5Split): Group 1 done. Height: %.2f. Ended on Page: %d, Y: %.2f\n", heightUsed1, pageAfterG1, yAfterG1)

	// --- Group 2 (Pics 2-4) --- Now try pics 3-5
	fmt.Println("Debug (process5Split): Processing Group 2 (Pics 2-4).")
	availableHeightForG2 := 0.0
	yBeforeProcessingG2 := yAfterG1

	if pageAfterG1 == e.currentPage.Page && heightUsed1 > 1e-6 {
		requiredSpacingG2 := e.imageSpacing // CORRECTED
		currentAvailableHeightAfterG1 := (e.marginTop + e.availableHeight) - yAfterG1
		if currentAvailableHeightAfterG1 < requiredSpacingG2 {
			fmt.Printf("Debug (process5Split): Not enough space for image spacing before G2 (Avail %.2f < Spacing %.2f). Forcing new page.\n", currentAvailableHeightAfterG1, requiredSpacingG2)
			e.newPage()
			yBeforeProcessingG2 = e.currentY
			availableHeightForG2 = (e.marginTop + e.availableHeight) - yBeforeProcessingG2
		} else {
			fmt.Printf("Debug (process5Split): Adding image spacing %.2f before Group 2.\n", requiredSpacingG2)
			yBeforeProcessingG2 = yAfterG1 + requiredSpacingG2
			e.currentY = yBeforeProcessingG2
			availableHeightForG2 = currentAvailableHeightAfterG1 - requiredSpacingG2
		}
	} else {
		yBeforeProcessingG2 = e.currentY
		availableHeightForG2 = (e.marginTop + e.availableHeight) - yBeforeProcessingG2
	}

	// Use processPictureGroup for the remaining 3, as it handles the retry logic internally
	heightUsed2, errG2 := e.processPictureGroup(pictures[2:5], 2, availableHeightForG2) // Group 2 (Pics 2-4)
	// pageAfterG2 := e.currentPage.Page // Capture if needed
	// yAfterG2 := e.currentY

	if errG2 != nil {
		fmt.Printf("Error (process5Split): Failed processing group 2 (Pics 2-4): %v\n", errG2)
		// If G2 fails, the split is incomplete. Return 0. Need strategy for partial placement.
		return 0
	}
	if heightUsed2 <= 1e-6 {
		fmt.Println("Warning (process5Split): Group 2 (Pics 2-4) was skipped.")
	}
	// fmt.Printf("Debug (process5Split): Group 2 done. Height: %.2f. Ended on Page: %d, Y: %.2f\n", heightUsed2, pageAfterG2, yAfterG2)

	fmt.Println("Debug (process5Split): Completed 2+3 split processing.")
	// Return 0 to indicate split path was taken. Caller relies on final e.currentY.
	return 0
}

// placeAllFiveOnNewPage helper (merged from Attempt 3 logic)
func (e *ContinuousLayoutEngine) placeAllFiveOnNewPage(pictures []Picture) float64 {
	fmt.Println("Debug (process5Split): Fallback: Placing all 5 on new page.")
	e.newPage()
	newAvailableHeight := (e.marginTop + e.availableHeight) - e.currentY
	fmt.Printf("Debug (process5Split Fallback): Placing all 5 on new page (Page %d). Available H: %.2f\n", e.currentPage.Page, newAvailableHeight)

	// Retry calculation for all 5 on the new page
	layoutInfo5Retry, err5Retry := e.calculateFivePicturesLayout(pictures, newAvailableHeight)
	if err5Retry != nil {
		fmt.Printf("Error (process5Split Fallback): Failed to calculate layout for all 5 pics even on new page: %v\n", err5Retry)
		return 0
	}
	if layoutInfo5Retry.TotalHeight > newAvailableHeight+1e-6 {
		fmt.Printf("Error (process5Split Fallback): Calculated height (%.2f) for 5 pics exceeds available height (%.2f) on new page.\n", layoutInfo5Retry.TotalHeight, newAvailableHeight)
		return 0
	}

	e.placePicturesInTemplate(pictures, layoutInfo5Retry)
	// Update Y and return height for this fallback placement
	e.currentY += layoutInfo5Retry.TotalHeight
	return layoutInfo5Retry.TotalHeight
}
