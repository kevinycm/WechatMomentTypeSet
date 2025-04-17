package calculate

import (
	"fmt"
)

// processEightPicturesWithSplitLogic handles layout for 8 pictures based on the 2+2+2+2 rule.
// UPDATED: Uses imageSpacing for inter-group spacing.
func (e *ContinuousLayoutEngine) processEightPicturesWithSplitLogic(pictures []Picture, layoutAvailableHeight float64) float64 {
	numPics := 8
	if len(pictures) != numPics {
		fmt.Printf("Error (process8Split): Incorrect number of pictures: %d\n", len(pictures))
		return 0
	}

	// --- Attempt 1: Try placing all 8 on the current page ---
	fmt.Println("Debug (process8Split): Attempting to place all 8 pictures initially.")
	layoutInfo8, err8 := e.calculateEightPicturesLayout(pictures, layoutAvailableHeight)

	if err8 == nil && layoutInfo8.TotalHeight <= layoutAvailableHeight+1e-6 { // Success and fits
		fmt.Println("Debug (process8Split): All 8 fit on current page directly.")
		e.placePicturesInTemplate(pictures, layoutInfo8)
		return layoutInfo8.TotalHeight
	}

	if err8 != nil {
		if err8.Error() == "force_new_page" {
			fmt.Println("Debug (process8Split): Initial calculation signaled force_new_page. Retrying on new page.")
			return e.placeAllEightOnNewPage(pictures) // Use helper for retry logic
		} else if err8.Error() == "split_required" {
			fmt.Println("Debug (process8Split): Initial calculation signaled split_required. Proceeding to 2+2+2+2.")
		} else {
			fmt.Printf("Debug (process8Split): Initial calculation failed (err: %v). Proceeding to split attempts.\n", err8)
		}
	} else if layoutInfo8.TotalHeight > layoutAvailableHeight { // Fits check failed
		fmt.Printf("Debug (process8Split): Initial calculation succeeded but doesn't fit (H: %.2f > Avail: %.2f). Proceeding to split attempts.\n", layoutInfo8.TotalHeight, layoutAvailableHeight)
	}

	// --- Attempt 2: Try 2 + 2 + 2 + 2 Split ---
	fmt.Println("Debug (process8Split): Attempting 2+2+2+2 split.")
	// Use processTwoPicGroup which handles its own retry/pagination logic
	// Capture state *after* each group call

	// --- Group 1 (Pics 0-1) ---
	heightUsed1, err2_1 := e.processTwoPicGroup(pictures[0:2], 1, layoutAvailableHeight)
	pageAfterG1 := e.currentPage.Page
	yAfterG1 := e.currentY

	if err2_1 != nil {
		fmt.Printf("Error (process8Split): Failed to place Group 1 (0-1): %v. Aborting.\n", err2_1)
		return 0
	}
	if heightUsed1 <= 1e-6 {
		fmt.Println("Warning (process8Split): Group 1 (0-1) was skipped.")
	}
	fmt.Printf("Debug (process8Split): Group 1 done. Height: %.2f. Ended on Page: %d, Y: %.2f\n", heightUsed1, pageAfterG1, yAfterG1)

	// --- Group 2 (Pics 2-3) ---
	availableHeightForG2 := 0.0
	yBeforeProcessingG2 := yAfterG1

	if pageAfterG1 == e.currentPage.Page && heightUsed1 > 1e-6 {
		requiredSpacingG2 := e.imageSpacing // CORRECTED
		currentAvailableHeightAfterG1 := (e.marginTop + e.availableHeight) - yAfterG1
		if currentAvailableHeightAfterG1 < requiredSpacingG2 {
			fmt.Printf("Debug (process8Split): Not enough space for image spacing before G2 (Avail %.2f < Spacing %.2f). Forcing new page.\n", currentAvailableHeightAfterG1, requiredSpacingG2)
			e.newPage()
			yBeforeProcessingG2 = e.currentY
			availableHeightForG2 = (e.marginTop + e.availableHeight) - yBeforeProcessingG2
		} else {
			fmt.Printf("Debug (process8Split): Adding image spacing %.2f before Group 2.\n", requiredSpacingG2)
			yBeforeProcessingG2 = yAfterG1 + requiredSpacingG2
			e.currentY = yBeforeProcessingG2
			availableHeightForG2 = currentAvailableHeightAfterG1 - requiredSpacingG2
		}
	} else {
		yBeforeProcessingG2 = e.currentY
		availableHeightForG2 = (e.marginTop + e.availableHeight) - yBeforeProcessingG2
	}

	heightUsed2, err2_2 := e.processTwoPicGroup(pictures[2:4], 2, availableHeightForG2)
	pageAfterG2 := e.currentPage.Page
	yAfterG2 := e.currentY

	if err2_2 != nil {
		fmt.Printf("Error (process8Split): Failed to place Group 2 (2-3): %v. Aborting remaining.\n", err2_2)
		return 0
	}
	if heightUsed2 <= 1e-6 {
		fmt.Println("Warning (process8Split): Group 2 (2-3) was skipped.")
	}
	fmt.Printf("Debug (process8Split): Group 2 done. Height: %.2f. Ended on Page: %d, Y: %.2f\n", heightUsed2, pageAfterG2, yAfterG2)

	// --- Group 3 (Pics 4-5) ---
	availableHeightForG3 := 0.0
	yBeforeProcessingG3 := yAfterG2

	if pageAfterG2 == e.currentPage.Page && heightUsed2 > 1e-6 {
		requiredSpacingG3 := e.imageSpacing // CORRECTED
		currentAvailableHeightAfterG2 := (e.marginTop + e.availableHeight) - yAfterG2
		if currentAvailableHeightAfterG2 < requiredSpacingG3 {
			fmt.Printf("Debug (process8Split): Not enough space for image spacing before G3 (Avail %.2f < Spacing %.2f). Forcing new page.\n", currentAvailableHeightAfterG2, requiredSpacingG3)
			e.newPage()
			yBeforeProcessingG3 = e.currentY
			availableHeightForG3 = (e.marginTop + e.availableHeight) - yBeforeProcessingG3
		} else {
			fmt.Printf("Debug (process8Split): Adding image spacing %.2f before Group 3.\n", requiredSpacingG3)
			yBeforeProcessingG3 = yAfterG2 + requiredSpacingG3
			e.currentY = yBeforeProcessingG3
			availableHeightForG3 = currentAvailableHeightAfterG2 - requiredSpacingG3
		}
	} else {
		yBeforeProcessingG3 = e.currentY
		availableHeightForG3 = (e.marginTop + e.availableHeight) - yBeforeProcessingG3
	}

	heightUsed3, err2_3 := e.processTwoPicGroup(pictures[4:6], 3, availableHeightForG3)
	pageAfterG3 := e.currentPage.Page
	yAfterG3 := e.currentY

	if err2_3 != nil {
		fmt.Printf("Error (process8Split): Failed to place Group 3 (4-5): %v. Aborting remaining.\n", err2_3)
		return 0
	}
	if heightUsed3 <= 1e-6 {
		fmt.Println("Warning (process8Split): Group 3 (4-5) was skipped.")
	}
	fmt.Printf("Debug (process8Split): Group 3 done. Height: %.2f. Ended on Page: %d, Y: %.2f\n", heightUsed3, pageAfterG3, yAfterG3)

	// --- Group 4 (Pics 6-7) ---
	availableHeightForG4 := 0.0
	yBeforeProcessingG4 := yAfterG3

	if pageAfterG3 == e.currentPage.Page && heightUsed3 > 1e-6 {
		requiredSpacingG4 := e.imageSpacing // CORRECTED
		currentAvailableHeightAfterG3 := (e.marginTop + e.availableHeight) - yAfterG3
		if currentAvailableHeightAfterG3 < requiredSpacingG4 {
			fmt.Printf("Debug (process8Split): Not enough space for image spacing before G4 (Avail %.2f < Spacing %.2f). Forcing new page.\n", currentAvailableHeightAfterG3, requiredSpacingG4)
			e.newPage()
			yBeforeProcessingG4 = e.currentY
			availableHeightForG4 = (e.marginTop + e.availableHeight) - yBeforeProcessingG4
		} else {
			fmt.Printf("Debug (process8Split): Adding image spacing %.2f before Group 4.\n", requiredSpacingG4)
			yBeforeProcessingG4 = yAfterG3 + requiredSpacingG4
			e.currentY = yBeforeProcessingG4
			availableHeightForG4 = currentAvailableHeightAfterG3 - requiredSpacingG4
		}
	} else {
		yBeforeProcessingG4 = e.currentY
		availableHeightForG4 = (e.marginTop + e.availableHeight) - yBeforeProcessingG4
	}

	heightUsed4, err2_4 := e.processTwoPicGroup(pictures[6:8], 4, availableHeightForG4)
	// pageAfterG4 := e.currentPage.Page
	// yAfterG4 := e.currentY

	if err2_4 != nil {
		fmt.Printf("Error (process8Split): Failed to place Group 4 (6-7): %v.\n", err2_4)
		return 0
	}
	if heightUsed4 <= 1e-6 {
		fmt.Println("Warning (process8Split): Group 4 (6-7) was skipped.")
	}
	// fmt.Printf("Debug (process8Split): Group 4 done. Height: %.2f. Ended on Page: %d, Y: %.2f\n", heightUsed4, pageAfterG4, yAfterG4)

	fmt.Println("Debug (process8Split): Completed 2+2+2+2 split processing.")
	return 0 // Return 0 indicates split path was taken
}

// fallbackSplitFourPlusFour needs similar spacing correction and linter fix
func (e *ContinuousLayoutEngine) fallbackSplitFourPlusFour(pictures []Picture, layoutAvailableHeight float64) float64 {
	fmt.Println("Debug (process8Split): Executing 4+4 fallback split.")

	// --- Group 1 (Pics 0-3) ---
	// LINTER FIX: Assign only height, assuming processFourPictures... returns float64
	heightUsed1 := e.processFourPicturesWithSplitLogic(pictures[0:4], layoutAvailableHeight)
	pageAfterG1 := e.currentPage.Page
	yAfterG1 := e.currentY

	// LINTER FIX: Check only height, not error
	if heightUsed1 <= 1e-6 {
		fmt.Println("Warning/Error (process8Split - 4+4): Failed to place or skipped first 4 pictures. Trying all 8 on new page.")
		return e.placeAllEightOnNewPage(pictures)
	}

	fmt.Printf("Debug (process8Split - 4+4): Group 1 (0-3) done. Height: %.2f. Ended on Page: %d, Y: %.2f\n", heightUsed1, pageAfterG1, yAfterG1)

	// --- Group 2 (Pics 4-7) ---
	availableHeightForG2 := 0.0
	yBeforeProcessingG2 := yAfterG1

	if pageAfterG1 == e.currentPage.Page && heightUsed1 > 1e-6 {
		requiredSpacingG2 := e.imageSpacing // CORRECTED
		currentAvailableHeightAfterG1 := (e.marginTop + e.availableHeight) - yAfterG1
		if currentAvailableHeightAfterG1 < requiredSpacingG2 {
			fmt.Printf("Debug (process8Split - 4+4): Not enough space for image spacing before G2 (Avail %.2f < Spacing %.2f). Forcing new page.\n", currentAvailableHeightAfterG1, requiredSpacingG2)
			e.newPage()
			yBeforeProcessingG2 = e.currentY
			availableHeightForG2 = (e.marginTop + e.availableHeight) - yBeforeProcessingG2
		} else {
			fmt.Printf("Debug (process8Split - 4+4): Adding image spacing %.2f before Group 2.\n", requiredSpacingG2)
			yBeforeProcessingG2 = yAfterG1 + requiredSpacingG2
			e.currentY = yBeforeProcessingG2
			availableHeightForG2 = currentAvailableHeightAfterG1 - requiredSpacingG2
		}
	} else {
		yBeforeProcessingG2 = e.currentY
		availableHeightForG2 = (e.marginTop + e.availableHeight) - yBeforeProcessingG2
	}

	// LINTER FIX: Assign only height
	heightUsed2 := e.processFourPicturesWithSplitLogic(pictures[4:8], availableHeightForG2)
	// pageAfterG2 := e.currentPage.Page
	// yAfterG2 := e.currentY

	// LINTER FIX: Check only height
	if heightUsed2 <= 1e-6 {
		fmt.Println("Warning/Error (process8Split - 4+4): Failed to place or skipped second 4 pictures.")
		// Decide on behavior: stop here returning 0? Or try to proceed?
		return 0 // Stop processing if second group fails/is skipped
	}
	// fmt.Printf("Debug (process8Split - 4+4): Group 2 (4-7) done. Height: %.2f. Ended on Page: %d, Y: %.2f\n", heightUsed2, pageAfterG2, yAfterG2)

	fmt.Println("Debug (process8Split - 4+4): Completed 4+4 split processing.")
	return 0 // Return 0 indicates split path was taken
}

// placeAllEightOnNewPage helper places all 8 on a new page.
func (e *ContinuousLayoutEngine) placeAllEightOnNewPage(pictures []Picture) float64 {
	fmt.Println("Debug (process8Split - Fallback): Attempting to place all 8 on a forced new page.")
	e.newPage()
	newAvailableHeight := e.availableHeight // Full height available
	layoutInfo8Retry, err8Retry := e.calculateEightPicturesLayout(pictures, newAvailableHeight)
	if err8Retry == nil && layoutInfo8Retry.TotalHeight <= newAvailableHeight+1e-6 {
		fmt.Println("Debug (process8Split - Fallback): Success placing all 8 on new page.")
		e.placePicturesInTemplate(pictures, layoutInfo8Retry)
		e.currentY += layoutInfo8Retry.TotalHeight // Update Y here!
		return layoutInfo8Retry.TotalHeight
	} else {
		fmt.Printf("Error (process8Split - Fallback): Failed to place all 8 even on new page (err: %v). Aborting.\n", err8Retry)
		return 0
	}
}

// --- Helper functions from process_layout_7_split.go are needed ---
// Need to ensure these are accessible, e.g., defined in a shared utility file or engine.go
// func removeLastNPicturesFromEntry(currentPage **ContinuousLayoutPage, n int) { ... }
// func getEntryPictureCount(currentPage **ContinuousLayoutPage) int { ... }
// func removePicturesAdded(currentPage **ContinuousLayoutPage, countBefore int) { ... }
// func GetPictureType(ar float64) string { ... }
// func GetRequiredMinHeight(e *ContinuousLayoutEngine, picType string, numPicsContext int) float64 { ... }
