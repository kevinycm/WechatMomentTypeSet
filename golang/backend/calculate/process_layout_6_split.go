package calculate

import (
	"fmt"
)

// processSixPicturesWithSplitLogic handles layout for 6 pictures,
// attempting 2+2+2 split if all 6 don't fit initially.
// UPDATED: Uses imageSpacing for inter-group spacing.
func (e *ContinuousLayoutEngine) processSixPicturesWithSplitLogic(pictures []Picture, layoutAvailableHeight float64) float64 {
	numPics := 6
	if len(pictures) != numPics {
		fmt.Printf("Error (process6Split): Incorrect number of pictures: %d\n", len(pictures))
		return 0
	}

	// --- Attempt 1: Try placing all 6 on the current page ---
	fmt.Println("Debug (process6Split): Attempting to place all 6 pictures initially.")
	layoutInfo6, err6 := e.calculateSixPicturesLayout(pictures, layoutAvailableHeight)

	if err6 == nil && layoutInfo6.TotalHeight <= layoutAvailableHeight+1e-6 { // Success and fits
		fmt.Println("Debug (process6Split): All 6 fit on current page.")
		e.placePicturesInTemplate(pictures, layoutInfo6)
		return layoutInfo6.TotalHeight
	}

	// --- Attempt 2: Try splitting 2+2+2 ---
	fmt.Printf("Debug (process6Split): All 6 failed/didn't fit (err: %v, H: %.2f > Avail: %.2f). Attempting 2+2+2 split.\n", err6, layoutInfo6.TotalHeight, layoutAvailableHeight)

	// Use processTwoPicGroup which handles its own retry/pagination logic
	// Capture state *after* each group call

	// --- Group 1 (Pics 0-1) ---
	fmt.Println("Debug (process6Split): Processing Group 1 (Pics 0-1)")
	heightUsed1, errG1 := e.processTwoPicGroup(pictures[0:2], 1, layoutAvailableHeight)
	pageAfterG1 := e.currentPage.Page
	yAfterG1 := e.currentY

	if errG1 != nil {
		fmt.Printf("Error (process6Split): Failed processing group 1: %v\n", errG1)
		// Try placing all 6 on a new page as a final fallback
		return e.placeAllSixOnNewPage(pictures)
	}
	if heightUsed1 <= 1e-6 {
		fmt.Println("Warning (process6Split): Group 1 (0-1) was skipped. Aborting 2+2+2 split.")
		// Fallback if first group skipped
		return e.placeAllSixOnNewPage(pictures)
	}
	fmt.Printf("Debug (process6Split): Group 1 done. Height: %.2f. Ended on Page: %d, Y: %.2f\n", heightUsed1, pageAfterG1, yAfterG1)

	// --- Group 2 (Pics 2-3) ---
	fmt.Println("Debug (process6Split): Processing Group 2 (Pics 2-3)")
	availableHeightForG2 := 0.0
	yBeforeProcessingG2 := yAfterG1

	if pageAfterG1 == e.currentPage.Page && heightUsed1 > 1e-6 {
		requiredSpacingG2 := e.imageSpacing // CORRECTED
		currentAvailableHeightAfterG1 := (e.marginTop + e.availableHeight) - yAfterG1
		if currentAvailableHeightAfterG1 < requiredSpacingG2 {
			fmt.Printf("Debug (process6Split): Not enough space for image spacing before G2 (Avail %.2f < Spacing %.2f). Forcing new page.\n", currentAvailableHeightAfterG1, requiredSpacingG2)
			e.newPage()
			yBeforeProcessingG2 = e.currentY
			availableHeightForG2 = (e.marginTop + e.availableHeight) - yBeforeProcessingG2
		} else {
			fmt.Printf("Debug (process6Split): Adding image spacing %.2f before Group 2.\n", requiredSpacingG2)
			yBeforeProcessingG2 = yAfterG1 + requiredSpacingG2
			e.currentY = yBeforeProcessingG2
			availableHeightForG2 = currentAvailableHeightAfterG1 - requiredSpacingG2
		}
	} else {
		yBeforeProcessingG2 = e.currentY
		availableHeightForG2 = (e.marginTop + e.availableHeight) - yBeforeProcessingG2
	}

	heightUsed2, errG2 := e.processTwoPicGroup(pictures[2:4], 2, availableHeightForG2)
	pageAfterG2 := e.currentPage.Page
	yAfterG2 := e.currentY

	if errG2 != nil {
		fmt.Printf("Error (process6Split): Failed processing group 2: %v\n", errG2)
		// If group 2 fails, we can't reliably continue the split.
		// Maybe try placing all 6 on new page? Or just fail?
		// For now, return 0, indicating split failed. Need strategy for partial placement.
		return 0
	}
	if heightUsed2 <= 1e-6 {
		fmt.Println("Warning (process6Split): Group 2 (2-3) was skipped. Proceeding with Group 3.")
		// Continue?
	}
	fmt.Printf("Debug (process6Split): Group 2 done. Height: %.2f. Ended on Page: %d, Y: %.2f\n", heightUsed2, pageAfterG2, yAfterG2)

	// --- Group 3 (Pics 4-5) ---
	fmt.Println("Debug (process6Split): Processing Group 3 (Pics 4-5)")
	availableHeightForG3 := 0.0
	yBeforeProcessingG3 := yAfterG2

	if pageAfterG2 == e.currentPage.Page && heightUsed2 > 1e-6 {
		requiredSpacingG3 := e.imageSpacing // CORRECTED
		currentAvailableHeightAfterG2 := (e.marginTop + e.availableHeight) - yAfterG2
		if currentAvailableHeightAfterG2 < requiredSpacingG3 {
			fmt.Printf("Debug (process6Split): Not enough space for image spacing before G3 (Avail %.2f < Spacing %.2f). Forcing new page.\n", currentAvailableHeightAfterG2, requiredSpacingG3)
			e.newPage()
			yBeforeProcessingG3 = e.currentY
			availableHeightForG3 = (e.marginTop + e.availableHeight) - yBeforeProcessingG3
		} else {
			fmt.Printf("Debug (process6Split): Adding image spacing %.2f before Group 3.\n", requiredSpacingG3)
			yBeforeProcessingG3 = yAfterG2 + requiredSpacingG3
			e.currentY = yBeforeProcessingG3
			availableHeightForG3 = currentAvailableHeightAfterG2 - requiredSpacingG3
		}
	} else {
		yBeforeProcessingG3 = e.currentY
		availableHeightForG3 = (e.marginTop + e.availableHeight) - yBeforeProcessingG3
	}

	heightUsed3, errG3 := e.processTwoPicGroup(pictures[4:6], 3, availableHeightForG3)
	// pageAfterG3 := e.currentPage.Page
	// yAfterG3 := e.currentY

	if errG3 != nil {
		fmt.Printf("Error (process6Split): Failed processing group 3: %v\n", errG3)
		return 0 // Fail the overall process
	}
	if heightUsed3 <= 1e-6 {
		fmt.Println("Warning (process6Split): Group 3 (4-5) was skipped.")
	}
	// fmt.Printf("Debug (process6Split): Group 3 done. Height: %.2f. Ended on Page: %d, Y: %.2f\n", heightUsed3, pageAfterG3, yAfterG3)

	fmt.Println("Debug (process6Split): Completed 2+2+2 split processing.")
	// Return 0 to indicate split path was taken. Caller relies on final e.currentY.
	return 0
}

// placeAllSixOnNewPage is a helper for the fallback case in 2+2+2 split
func (e *ContinuousLayoutEngine) placeAllSixOnNewPage(pictures []Picture) float64 {
	fmt.Println("Debug (process6Split): Fallback: Placing all 6 on new page.")
	e.newPage()
	newAvailableHeight := (e.marginTop + e.availableHeight) - e.currentY
	fmt.Printf("Debug (process6Split Fallback): Placing all 6 on new page (Page %d). Available H: %.2f\n", e.currentPage.Page, newAvailableHeight)

	layoutInfo6Retry, err6Retry := e.calculateSixPicturesLayout(pictures, newAvailableHeight)
	if err6Retry != nil {
		fmt.Printf("Error (process6Split Fallback): Failed to calculate layout for all 6 pics even on new page: %v\n", err6Retry)
		return 0
	}
	if layoutInfo6Retry.TotalHeight > newAvailableHeight+1e-6 {
		fmt.Printf("Error (process6Split Fallback): Calculated height (%.2f) for 6 pics exceeds available height (%.2f) on new page.\n", layoutInfo6Retry.TotalHeight, newAvailableHeight)
		return 0
	}

	e.placePicturesInTemplate(pictures, layoutInfo6Retry)
	// We need to update Y coordinate here since the caller might rely on it
	e.currentY += layoutInfo6Retry.TotalHeight
	return layoutInfo6Retry.TotalHeight // Return height used in this fallback case
}
