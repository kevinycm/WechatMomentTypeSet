package calculate

import (
	"fmt"
)

// Removed Placeholder for GetPictureTypeAR - Assume defined elsewhere in package
// Removed Placeholder for GetRequiredMinHeight - Assume defined elsewhere in package

// processSevenPicturesWithSplitLogic handles layout for 7 pictures based on the 2+2+2+1 rule.
// UPDATED: Uses imageSpacing for inter-group spacing.
func (e *ContinuousLayoutEngine) processSevenPicturesWithSplitLogic(pictures []Picture, layoutAvailableHeight float64) float64 {
	numPics := 7
	if len(pictures) != numPics {
		fmt.Printf("Error (process7Split): Incorrect number of pictures: %d\n", len(pictures))
		return 0
	}

	// --- Attempt 1: Try placing all 7 on the current page ---
	fmt.Println("Debug (process7Split): Attempting to place all 7 pictures initially.")
	layoutInfo7, err7 := e.calculateSevenPicturesLayout(pictures, layoutAvailableHeight)

	if err7 == nil && layoutInfo7.TotalHeight <= layoutAvailableHeight+1e-6 { // Success and fits
		fmt.Println("Debug (process7Split): All 7 fit on current page directly.")
		e.placePicturesInTemplate(pictures, layoutInfo7)
		return layoutInfo7.TotalHeight
	}

	if err7 != nil {
		if err7.Error() == "force_new_page" {
			fmt.Println("Debug (process7Split): Initial calculation signaled force_new_page. Retrying on new page.")
			return e.placeAllSevenOnNewPage(pictures) // Use helper for retry logic
		}
		fmt.Printf("Debug (process7Split): Initial calculation failed or needs split (err: %v). Proceeding to split attempts.\n", err7)
	} else if layoutInfo7.TotalHeight > layoutAvailableHeight { // Fits check failed
		fmt.Printf("Debug (process7Split): Initial calculation succeeded but doesn't fit (H: %.2f > Avail: %.2f). Proceeding to split attempts.\n", layoutInfo7.TotalHeight, layoutAvailableHeight)
	}

	// --- Attempt 2: Try 2 + 2 + 2 + 1 Split ---
	fmt.Println("Debug (process7Split): Attempting 2+2+2+1 split.")
	// Use processTwoPicGroup which handles its own retry/pagination logic
	// Capture state *after* each group call

	// --- Group 1 (Pics 0-1) ---
	heightUsed1, err2_1 := e.processTwoPicGroup(pictures[0:2], 1, layoutAvailableHeight) // Group 1 (pics 0-1)
	pageAfterG1 := e.currentPage.Page
	yAfterG1 := e.currentY

	if err2_1 != nil {
		fmt.Printf("Error (process7Split): Failed to place Group 1 (0-1): %v. Falling back to 3+4.\n", err2_1)
		// Fallback to 3+4 logic if first group fails
		return e.fallbackSplitThreePlusFour(pictures, layoutAvailableHeight)
	}
	if heightUsed1 <= 1e-6 {
		fmt.Println("Warning (process7Split): Group 1 (0-1) was skipped. Falling back to 3+4.")
		// Fallback to 3+4 logic if first group is skipped
		return e.fallbackSplitThreePlusFour(pictures, layoutAvailableHeight)
	}
	fmt.Printf("Debug (process7Split): Group 1 done. Height: %.2f. Ended on Page: %d, Y: %.2f\n", heightUsed1, pageAfterG1, yAfterG1)

	// --- Group 2 (Pics 2-3) ---
	availableHeightForG2 := 0.0
	yBeforeProcessingG2 := yAfterG1

	if pageAfterG1 == e.currentPage.Page && heightUsed1 > 1e-6 {
		requiredSpacingG2 := e.imageSpacing // CORRECTED
		currentAvailableHeightAfterG1 := (e.marginTop + e.availableHeight) - yAfterG1
		if currentAvailableHeightAfterG1 < requiredSpacingG2 {
			fmt.Printf("Debug (process7Split): Not enough space for image spacing before G2 (Avail %.2f < Spacing %.2f). Forcing new page.\n", currentAvailableHeightAfterG1, requiredSpacingG2)
			e.newPage()
			yBeforeProcessingG2 = e.currentY
			availableHeightForG2 = (e.marginTop + e.availableHeight) - yBeforeProcessingG2
		} else {
			fmt.Printf("Debug (process7Split): Adding image spacing %.2f before Group 2.\n", requiredSpacingG2)
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
		fmt.Printf("Error (process7Split): Failed to place Group 2 (2-3): %v. Aborting remaining and falling back to placing 3-7 on new page.\n", err2_2)
		// Fallback: Place remaining 5 (indices 2-6) on a new page if G2 fails
		e.newPage()
		newAvailableHeight := (e.marginTop + e.availableHeight) - e.currentY
		return e.processFivePicturesWithSplitLogic(pictures[2:7], newAvailableHeight) // Place 3-7
	}
	if heightUsed2 <= 1e-6 {
		fmt.Println("Warning (process7Split): Group 2 (2-3) was skipped. Proceeding with Group 3.")
		// Continue even if skipped? Or fallback? Let's continue for now.
	}
	fmt.Printf("Debug (process7Split): Group 2 done. Height: %.2f. Ended on Page: %d, Y: %.2f\n", heightUsed2, pageAfterG2, yAfterG2)

	// --- Group 3 (Pics 4-5) ---
	availableHeightForG3 := 0.0
	yBeforeProcessingG3 := yAfterG2

	if pageAfterG2 == e.currentPage.Page && heightUsed2 > 1e-6 {
		requiredSpacingG3 := e.imageSpacing // CORRECTED
		currentAvailableHeightAfterG2 := (e.marginTop + e.availableHeight) - yAfterG2
		if currentAvailableHeightAfterG2 < requiredSpacingG3 {
			fmt.Printf("Debug (process7Split): Not enough space for image spacing before G3 (Avail %.2f < Spacing %.2f). Forcing new page.\n", currentAvailableHeightAfterG2, requiredSpacingG3)
			e.newPage()
			yBeforeProcessingG3 = e.currentY
			availableHeightForG3 = (e.marginTop + e.availableHeight) - yBeforeProcessingG3
		} else {
			fmt.Printf("Debug (process7Split): Adding image spacing %.2f before Group 3.\n", requiredSpacingG3)
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
		fmt.Printf("Error (process7Split): Failed to place Group 3 (4-5): %v. Aborting remaining and falling back to placing 5-7 on new page.\n", err2_3)
		// Fallback: Place remaining 3 (indices 4-6) on a new page if G3 fails
		e.newPage()
		newAvailableHeight := (e.marginTop + e.availableHeight) - e.currentY
		return e.processThreePicturesWithSplitLogic(pictures[4:7], newAvailableHeight) // Place 5-7
	}
	if heightUsed3 <= 1e-6 {
		fmt.Println("Warning (process7Split): Group 3 (4-5) was skipped. Proceeding with final picture.")
		// Continue even if skipped? Let's continue.
	}
	fmt.Printf("Debug (process7Split): Group 3 done. Height: %.2f. Ended on Page: %d, Y: %.2f\n", heightUsed3, pageAfterG3, yAfterG3)

	// --- Final Picture (Pic 6) ---
	availableHeightForG4 := 0.0
	yBeforeProcessingG4 := yAfterG3

	if pageAfterG3 == e.currentPage.Page && heightUsed3 > 1e-6 {
		requiredSpacingG4 := e.imageSpacing // CORRECTED (Spacing before the single pic)
		currentAvailableHeightAfterG3 := (e.marginTop + e.availableHeight) - yAfterG3
		if currentAvailableHeightAfterG3 < requiredSpacingG4 {
			fmt.Printf("Debug (process7Split): Not enough space for image spacing before Final Pic (Avail %.2f < Spacing %.2f). Forcing new page.\n", currentAvailableHeightAfterG3, requiredSpacingG4)
			e.newPage()
			yBeforeProcessingG4 = e.currentY
			availableHeightForG4 = (e.marginTop + e.availableHeight) - yBeforeProcessingG4
		} else {
			fmt.Printf("Debug (process7Split): Adding image spacing %.2f before Final Pic.\n", requiredSpacingG4)
			yBeforeProcessingG4 = yAfterG3 + requiredSpacingG4
			e.currentY = yBeforeProcessingG4
			availableHeightForG4 = currentAvailableHeightAfterG3 - requiredSpacingG4
		}
	} else {
		yBeforeProcessingG4 = e.currentY
		availableHeightForG4 = (e.marginTop + e.availableHeight) - yBeforeProcessingG4
	}

	// Assuming processSinglePictureLayoutAndPlace handles its own retry/pagination
	// and returns height used or 0 on failure/skip.
	lastPicHeight := e.processSinglePictureLayoutAndPlace(pictures[6], availableHeightForG4)
	// pageAfterG4 := e.currentPage.Page
	// yAfterG4 := e.currentY

	if lastPicHeight <= 1e-6 {
		fmt.Printf("Error/Warning (process7Split): Failed to place or skipped final picture (Index %d).\n", pictures[6].Index)
		// Decide behavior - return 0?
		return 0
	}
	// fmt.Printf("Debug (process7Split): Final Pic done. Height: %.2f. Ended on Page: %d, Y: %.2f\n", lastPicHeight, pageAfterG4, yAfterG4)

	fmt.Println("Debug (process7Split): Completed 2+2+2+1 split processing.")
	return 0 // Return 0 indicates split path was taken
}

// fallbackSplitThreePlusFour needs spacing correction between the groups
func (e *ContinuousLayoutEngine) fallbackSplitThreePlusFour(pictures []Picture, layoutAvailableHeight float64) float64 {
	fmt.Println("Debug (process7Split): Executing 3+4 fallback split.")

	// --- Group 1 (Pics 0-2) ---
	// Assume processThreePictures returns float64
	heightUsed1 := e.processThreePicturesWithSplitLogic(pictures[0:3], layoutAvailableHeight)
	pageAfterG1 := e.currentPage.Page
	yAfterG1 := e.currentY

	// Check if the first group failed or was skipped
	if heightUsed1 <= 1e-6 {
		fmt.Printf("Error/Warning (process7Split - 3+4): Failed/Skipped placement of first 3 pictures. Trying final fallback (all 7 on new page).\n")
		return e.placeAllSevenOnNewPage(pictures) // Try final fallback
	}

	fmt.Printf("Debug (process7Split - 3+4): Group 1 (0-2) done. Height: %.2f. Ended on Page: %d, Y: %.2f\n", heightUsed1, pageAfterG1, yAfterG1)

	// --- Group 2 (Pics 3-6) ---
	availableHeightForG2 := 0.0
	yBeforeProcessingG2 := yAfterG1

	if pageAfterG1 == e.currentPage.Page && heightUsed1 > 1e-6 {
		requiredSpacingG2 := e.imageSpacing // CORRECTED
		currentAvailableHeightAfterG1 := (e.marginTop + e.availableHeight) - yAfterG1
		if currentAvailableHeightAfterG1 < requiredSpacingG2 {
			fmt.Printf("Debug (process7Split - 3+4): Not enough space for image spacing before G2 (Avail %.2f < Spacing %.2f). Forcing new page.\n", currentAvailableHeightAfterG1, requiredSpacingG2)
			e.newPage()
			yBeforeProcessingG2 = e.currentY
			availableHeightForG2 = (e.marginTop + e.availableHeight) - yBeforeProcessingG2
		} else {
			fmt.Printf("Debug (process7Split - 3+4): Adding image spacing %.2f before Group 2.\n", requiredSpacingG2)
			yBeforeProcessingG2 = yAfterG1 + requiredSpacingG2
			e.currentY = yBeforeProcessingG2
			availableHeightForG2 = currentAvailableHeightAfterG1 - requiredSpacingG2
		}
	} else {
		yBeforeProcessingG2 = e.currentY
		availableHeightForG2 = (e.marginTop + e.availableHeight) - yBeforeProcessingG2
	}

	// Assume processFourPictures returns float64
	heightUsed2 := e.processFourPicturesWithSplitLogic(pictures[3:7], availableHeightForG2)
	// pageAfterG2 := e.currentPage.Page
	// yAfterG2 := e.currentY

	if heightUsed2 <= 1e-6 {
		fmt.Printf("Error/Warning (process7Split - 3+4): Failed/Skipped placement of second 4 pictures.\n")
		// Decide behavior: return 0 indicates failure of this split path
		return 0
	}
	// fmt.Printf("Debug (process7Split - 3+4): Group 2 (3-6) done. Height: %.2f. Ended on Page: %d, Y: %.2f\n", heightUsed2, pageAfterG2, yAfterG2)

	fmt.Println("Debug (process7Split - 3+4): Completed 3+4 split processing.")
	return 0 // Return 0 indicates split path was taken
}

// --- Helper functions assumed to exist or need implementation ---

// removeLastNPicturesFromEntry removes the last N pictures from the last entry on the current page.
func removeLastNPicturesFromEntry(currentPage **ContinuousLayoutPage, n int) {
	if currentPage == nil || *currentPage == nil || len((*currentPage).Entries) == 0 {
		fmt.Printf("Warning (removePictures): Current page or entries list is nil/empty.\n")
		return
	}
	lastEntry := &(*currentPage).Entries[len((*currentPage).Entries)-1]
	if len(lastEntry.Pictures) >= n {
		lastEntry.Pictures = lastEntry.Pictures[:len(lastEntry.Pictures)-n]
		fmt.Printf("Debug (removePictures): Removed last %d pictures from Page %d Entry %d.\n", n, (*currentPage).Page, len((*currentPage).Entries)-1)
	} else {
		fmt.Printf("Warning (removePictures): Tried to remove %d pictures, but entry only has %d.\n", n, len(lastEntry.Pictures))
	}
}

// getEntryPictureCount returns the number of pictures in the last entry of the current page.
func getEntryPictureCount(currentPage **ContinuousLayoutPage) int {
	if currentPage == nil || *currentPage == nil || len((*currentPage).Entries) == 0 {
		return 0
	}
	return len((*currentPage).Entries[len((*currentPage).Entries)-1].Pictures)
}

// removePicturesAdded removes pictures added since a certain count.
func removePicturesAdded(currentPage **ContinuousLayoutPage, countBefore int) {
	if currentPage == nil || *currentPage == nil || len((*currentPage).Entries) == 0 {
		// fmt.Printf("Warning (removePicturesAdded): Current page or entries list is nil/empty.\n")
		return
	}
	lastEntry := &(*currentPage).Entries[len((*currentPage).Entries)-1]
	countAfter := len(lastEntry.Pictures)
	numToRemove := countAfter - countBefore
	if numToRemove > 0 && countAfter >= numToRemove {
		lastEntry.Pictures = lastEntry.Pictures[:countAfter-numToRemove]
		fmt.Printf("Debug (removePicturesAdded): Removed last %d pictures added from Page %d Entry %d.\n", numToRemove, (*currentPage).Page, len((*currentPage).Entries)-1)
	} else if numToRemove > 0 {
		// This case indicates a logic error elsewhere or state corruption
		fmt.Printf("Warning (removePicturesAdded): Logic error? Tried to remove %d pics, but entry only has %d (started with %d).\n", numToRemove, countAfter, countBefore)
	}
}

// placeAllSevenOnNewPage helper remains the same
func (e *ContinuousLayoutEngine) placeAllSevenOnNewPage(pictures []Picture) float64 {
	fmt.Println("Debug (process7Split - Fallback): Attempting to place all 7 on a forced new page.")
	e.newPage()
	newAvailableHeight := e.availableHeight // Full height available
	layoutInfo7Retry, err7Retry := e.calculateSevenPicturesLayout(pictures, newAvailableHeight)
	if err7Retry == nil && layoutInfo7Retry.TotalHeight <= newAvailableHeight+1e-6 {
		fmt.Println("Debug (process7Split - Fallback): Success placing all 7 on new page.")
		e.placePicturesInTemplate(pictures, layoutInfo7Retry)
		e.currentY += layoutInfo7Retry.TotalHeight // Update Y here!
		return layoutInfo7Retry.TotalHeight
	} else {
		fmt.Printf("Error (process7Split - Fallback): Failed to place all 7 even on new page (err: %v). Aborting.\n", err7Retry)
		return 0
	}
}

// GetPictureType function (assuming it takes ar float64)
// Ensure definition exists elsewhere (e.g., engine.go)
// func GetPictureType(ar float64) string { ... }

// GetRequiredMinHeight function
// Ensure definition exists elsewhere (e.g., engine.go)
// func GetRequiredMinHeight(e *ContinuousLayoutEngine, picType string, numPicsContext int) float64 { ... }

// processSinglePictureLayoutAndPlace function
// Ensure definition exists elsewhere (e.g., engine.go)
// func (e *ContinuousLayoutEngine) processSinglePictureLayoutAndPlace(picture Picture, layoutAvailableHeight float64) float64 { ... }

// processThreePicturesWithSplitLogic function
// LINTER FIX: Assume returns float64
// Ensure definition exists elsewhere (e.g., process_layout_3_split.go)
// func (e *ContinuousLayoutEngine) processThreePicturesWithSplitLogic(pictures []Picture, layoutAvailableHeight float64) float64 { ... }

// processFivePicturesWithSplitLogic function
// Ensure definition exists elsewhere (e.g., process_layout_5_split.go)
// func (e *ContinuousLayoutEngine) processFivePicturesWithSplitLogic(pictures []Picture, layoutAvailableHeight float64) float64 { ... } // Note: Return error? Assume float64

// processFourPicturesWithSplitLogic function
// LINTER FIX: Assume returns float64
// Ensure definition exists elsewhere (e.g., process_layout_4_split.go)
// func (e *ContinuousLayoutEngine) processFourPicturesWithSplitLogic(pictures []Picture, layoutAvailableHeight float64) float64 { ... }
