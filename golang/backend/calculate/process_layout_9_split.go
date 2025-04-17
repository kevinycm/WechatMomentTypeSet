package calculate

import (
	"errors"
	"fmt"
	// Ensure math is imported if needed for calculations
	// Added for error handling
)

// processPictureGroup attempts to calculate and place a group of pictures (typically 3).
// It handles pagination if the group doesn't fit the initial availableHeight.
// It updates the engine's currentY state directly upon successful placement.
// Returns the height used by the group and any critical error.
func (e *ContinuousLayoutEngine) processPictureGroup(
	groupPics []Picture,
	groupNum int, // For logging (e.g., 1st, 2nd, 3rd group of 3)
	layoutAvailableHeight float64,
) (heightUsed float64, err error) {

	if len(groupPics) == 0 {
		return 0, errors.New("processPictureGroup: called with empty picture group")
	}
	numPicsInGroup := len(groupPics)
	var layoutInfo TemplateLayout
	var calcErr error
	var initialPlacementErr error // Store error from initial placement attempt

	fmt.Printf("Debug (Split-Group %d): Attempting calculation. AvailableHeight: %.2f\n", groupNum, layoutAvailableHeight)

	// Choose calculation function based on group size (primarily 3 for 9-split)
	switch numPicsInGroup {
	case 3:
		// Use calculateLayout method that checks min height and returns error
		layoutInfo, calcErr = e.calculateThreePicturesLayout(groupPics, layoutAvailableHeight)
	// Add cases for 6, 4, 2 if this helper is generalized later
	default:
		return 0, fmt.Errorf("processPictureGroup: unsupported group size %d", numPicsInGroup)
	}

	// --- Check Fit and Place or Retry ---
	if calcErr == nil && layoutInfo.TotalHeight <= layoutAvailableHeight+1e-6 { // Add tolerance
		// Fits on current page
		fmt.Printf("Debug (Split-Group %d): Group fits on current page (Page %d, height: %.2f <= available: %.2f)\n", groupNum, e.currentPage.Page, layoutInfo.TotalHeight, layoutAvailableHeight)
		e.placePicturesInTemplate(groupPics, layoutInfo)
		e.currentY += layoutInfo.TotalHeight // Update Y
		return layoutInfo.TotalHeight, nil
	} else {
		// Doesn't fit or calculation error - Store the error and force New Page
		if calcErr != nil {
			initialPlacementErr = calcErr // Store calculation error (e.g., min height)
			fmt.Printf("Debug (Split-Group %d): Initial calc failed: %v. Forcing new page.\n", groupNum, calcErr)
		} else {
			// Only reason left is height doesn't fit
			initialPlacementErr = fmt.Errorf("group height %.2f exceeds available %.2f", layoutInfo.TotalHeight, layoutAvailableHeight)
			fmt.Printf("Debug (Split-Group %d): Group doesn't fit (%.2f > %.2f). Forcing new page.\n", groupNum, layoutInfo.TotalHeight, layoutAvailableHeight)
		}
		e.newPage()
		newAvailableHeight := (e.marginTop + e.availableHeight) - e.currentY

		// Retry calculation on the new page
		fmt.Printf("Debug (Split-Group %d): Retrying calculation & placement on new page (Page %d, available: %.2f)\n", groupNum, e.currentPage.Page, newAvailableHeight)
		switch numPicsInGroup {
		case 3:
			layoutInfo, calcErr = e.calculateThreePicturesLayout(groupPics, newAvailableHeight)
		default: // Should not happen based on initial check
			return 0, fmt.Errorf("processPictureGroup: retry unsupported group size %d", numPicsInGroup)
		}

		// --- Error Handling for Retry ---
		if calcErr != nil { // Check if calculation itself failed on retry
			// Check if the error is the specific minimum height constraint failure
			if calcErr.Error() == "no layout satisfied minimum height requirements for 3 pictures" {
				fmt.Printf("Warning (Split-Group %d): Retry failed - layout constraints (min height) not met even on new page. Skipping group. Initial error: %v\n", groupNum, initialPlacementErr)
				// Return 0 height and nil error - indicates failure but prevents higher-level errors if desired
				return 0, fmt.Errorf("skipped group %d due to min height failure on new page (initial error: %w)", groupNum, initialPlacementErr)
			} else {
				// For any other calculation error on retry, report it
				errMsg := fmt.Sprintf("failed to place group %d even on new page: %v (initial error: %v)", groupNum, calcErr, initialPlacementErr)
				fmt.Printf("Error: %s\n", errMsg)
				return 0, errors.New(errMsg) // Propagate other errors
			}
		} else if layoutInfo.TotalHeight > newAvailableHeight+1e-6 { // Check if it fits height-wise on retry
			// Calculation succeeded, but still doesn't fit the new page height
			errMsg := fmt.Sprintf("failed to place group %d even on new page (height %.2f > available %.2f) (initial error: %v)", groupNum, layoutInfo.TotalHeight, newAvailableHeight, initialPlacementErr)
			fmt.Printf("Error: %s\n", errMsg)
			return 0, errors.New(errMsg) // Propagate height overflow error
		}
		// --- End Error Handling ---

		// If we reach here, retry was successful (calc OK, fits height)
		fmt.Printf("Debug (Split-Group %d): Placed group successfully on new page (Page %d).\n", groupNum, e.currentPage.Page)
		// Place successfully on new page
		e.placePicturesInTemplate(groupPics, layoutInfo)
		e.currentY += layoutInfo.TotalHeight // Update Y on new page
		return layoutInfo.TotalHeight, nil
	}
}

// processNinePicturesWithSplitLogic implements the new 1-5 rules.
// It attempts 9-pic layout first, then falls back to 3+3+3 with proactive page breaks.
// Relies on processPictureGroup to handle its own placement and internal page breaks if needed.
// Returns the *initial* height calculated for the first successful placement (either 9-pic or first group)
// or 0 if splitting occurs or fails completely. The caller should rely on e.currentY for final position.
func (e *ContinuousLayoutEngine) processNinePicturesWithSplitLogic(pictures []Picture, layoutAvailableHeight float64) float64 {
	numPics := 9
	if len(pictures) != numPics {
		fmt.Printf("Error (process9Split): Incorrect number of pictures: %d\n", len(pictures))
		return 0
	}

	// --- Rule 1: Attempt placing all 9 on the current page ---
	fmt.Println("Debug (process9Split): Rule 1 - Attempting to place all 9 pictures initially.")
	layoutInfo9, err9 := e.calculateNinePicturesLayout(pictures, layoutAvailableHeight)

	if err9 == nil && layoutInfo9.TotalHeight <= layoutAvailableHeight+1e-6 {
		fmt.Println("Debug (process9Split): Rule 1 - Success! All 9 fit on current page directly.")
		e.placePicturesInTemplate(pictures, layoutInfo9)
		return layoutInfo9.TotalHeight
	} else {
		if err9 != nil {
			fmt.Printf("Debug (process9Split): Rule 1 failed (calculation error: %v). Proceeding to split.\n", err9)
		} else {
			fmt.Printf("Debug (process9Split): Rule 1 failed (height %.2f > available %.2f). Proceeding to split.\n", layoutInfo9.TotalHeight, layoutAvailableHeight)
		}
	}

	// --- Rule 1 Failed: Proceed to 3+3+3 Split ---
	fmt.Println("Debug (process9Split): Rule 2 - Attempting 3+3+3 split.")

	// --- Group 1 (Pics 0-2) ---
	fmt.Println("Debug (process9Split): Processing Group 1 (0-2).")
	// yBeforeG1 := e.currentY
	heightG1, errG1 := e.processPictureGroup(pictures[0:3], 1, layoutAvailableHeight)
	pageAfterG1 := e.currentPage.Page // Capture state AFTER processing G1
	yAfterG1 := e.currentY

	if errG1 != nil {
		fmt.Printf("Error (process9Split): Failed to place Group 1 (0-2): %v. Aborting.\n", errG1)
		// Even if G1 failed, Y might have changed if newPage was called internally. Reset?
		// Let's assume processPictureGroup handles its state consistently on failure.
		// For safety, maybe reset Y to before G1? Depends on desired behavior.
		// e.currentY = yBeforeG1 // Optional: Reset Y if G1 fails critically
		return 0
	}
	if heightG1 <= 1e-6 {
		fmt.Println("Warning (process9Split): Group 1 (0-2) was skipped.")
	}
	fmt.Printf("Debug (process9Split): Group 1 done. Height: %.2f. Ended on Page: %d, Y: %.2f\n", heightG1, pageAfterG1, yAfterG1)

	// --- Group 2 (Pics 3-5) ---
	fmt.Println("Debug (process9Split): Processing Group 2 (3-5).")
	// Check if spacing is needed *before* Group 2
	availableHeightForG2 := 0.0
	yBeforeProcessingG2 := yAfterG1 // Start from where G1 ended

	// Check if G2 needs to start on a new page due to spacing vs available height *after* G1
	if pageAfterG1 == e.currentPage.Page && heightG1 > 1e-6 { // Only add spacing if G1 was actually placed on the page G2 might start on.
		requiredSpacingG2 := e.imageSpacing
		currentAvailableHeightAfterG1 := (e.marginTop + e.availableHeight) - yAfterG1
		if currentAvailableHeightAfterG1 < requiredSpacingG2 {
			fmt.Printf("Debug (process9Split): Not enough space for image spacing before G2 (Avail %.2f < Spacing %.2f). Forcing new page.\n", currentAvailableHeightAfterG1, requiredSpacingG2)
			e.newPage()
			yBeforeProcessingG2 = e.currentY // Y is now marginTop on new page
			availableHeightForG2 = (e.marginTop + e.availableHeight) - yBeforeProcessingG2
		} else {
			fmt.Printf("Debug (process9Split): Adding image spacing %.2f before Group 2.\n", requiredSpacingG2)
			yBeforeProcessingG2 = yAfterG1 + requiredSpacingG2 // Calculate where G2 *should* start
			e.currentY = yBeforeProcessingG2                   // Update engine state *before* calling processPictureGroup
			availableHeightForG2 = currentAvailableHeightAfterG1 - requiredSpacingG2
		}
	} else {
		// G2 starts on whatever page G1 left us on (potentially a new page)
		// Y should already be correctly set by G1's processing or the newPage call above.
		yBeforeProcessingG2 = e.currentY
		availableHeightForG2 = (e.marginTop + e.availableHeight) - yBeforeProcessingG2
	}

	// Process Group 2
	heightG2, errG2 := e.processPictureGroup(pictures[3:6], 2, availableHeightForG2)
	pageAfterG2 := e.currentPage.Page // Capture state AFTER processing G2
	yAfterG2 := e.currentY

	if errG2 != nil {
		fmt.Printf("Error (process9Split): Failed to place Group 2 (3-5): %v. Aborting remaining.\n", errG2)
		// Reset Y to before G2 attempt?
		// e.currentY = yBeforeProcessingG2 // Optional reset
		return 0
	}
	if heightG2 <= 1e-6 {
		fmt.Println("Warning (process9Split): Group 2 (3-5) was skipped.")
	}
	fmt.Printf("Debug (process9Split): Group 2 done. Height: %.2f. Ended on Page: %d, Y: %.2f\n", heightG2, pageAfterG2, yAfterG2)

	// --- Group 3 (Pics 6-8) ---
	fmt.Println("Debug (process9Split): Processing Group 3 (6-8).")
	// Check if spacing is needed *before* Group 3
	availableHeightForG3 := 0.0
	yBeforeProcessingG3 := yAfterG2 // Start from where G2 ended

	// Check if G3 needs to start on a new page due to spacing vs available height *after* G2
	if pageAfterG2 == e.currentPage.Page && heightG2 > 1e-6 { // Only add spacing if G2 was actually placed on the page G3 might start on.
		requiredSpacingG3 := e.imageSpacing
		currentAvailableHeightAfterG2 := (e.marginTop + e.availableHeight) - yAfterG2
		if currentAvailableHeightAfterG2 < requiredSpacingG3 {
			fmt.Printf("Debug (process9Split): Not enough space for image spacing before G3 (Avail %.2f < Spacing %.2f). Forcing new page.\n", currentAvailableHeightAfterG2, requiredSpacingG3)
			e.newPage()
			yBeforeProcessingG3 = e.currentY // Y is now marginTop on new page
			availableHeightForG3 = (e.marginTop + e.availableHeight) - yBeforeProcessingG3
		} else {
			fmt.Printf("Debug (process9Split): Adding image spacing %.2f before Group 3.\n", requiredSpacingG3)
			yBeforeProcessingG3 = yAfterG2 + requiredSpacingG3 // Calculate where G3 *should* start
			e.currentY = yBeforeProcessingG3                   // Update engine state *before* calling processPictureGroup
			availableHeightForG3 = currentAvailableHeightAfterG2 - requiredSpacingG3
		}
	} else {
		// G3 starts on whatever page G2 left us on (potentially a new page)
		// Y should already be correctly set by G2's processing or the newPage call above.
		yBeforeProcessingG3 = e.currentY
		availableHeightForG3 = (e.marginTop + e.availableHeight) - yBeforeProcessingG3
	}

	// Process Group 3
	heightG3, errG3 := e.processPictureGroup(pictures[6:9], 3, availableHeightForG3)
	// pageAfterG3 := e.currentPage.Page // Capture final state if needed
	// yAfterG3 := e.currentY

	if errG3 != nil {
		fmt.Printf("Error (process9Split): Failed to place Group 3 (6-8): %v.\n", errG3)
		// Reset Y to before G3 attempt?
		// e.currentY = yBeforeProcessingG3 // Optional reset
		return 0
	}
	if heightG3 <= 1e-6 {
		fmt.Println("Warning (process9Split): Group 3 (6-8) was skipped.")
	}
	// Final Y position is now managed by processPictureGroup calls.
	// fmt.Printf("Debug (process9Split): Group 3 done. Height: %.2f. Ended on Page: %d, Y: %.2f\n", heightG3, pageAfterG3, yAfterG3)

	fmt.Println("Debug (process9Split): Completed 3+3+3 split processing.")
	// Return 0 to indicate split path was taken. Caller relies on final e.currentY.
	return 0
}

// placeAllNineOnNewPage is a helper function specifically for retrying the 9-picture layout on a new page.
func (e *ContinuousLayoutEngine) placeAllNineOnNewPage(pictures []Picture) float64 {
	e.newPage()
	newAvailableHeight := (e.marginTop + e.availableHeight) - e.currentY
	fmt.Printf("Debug (placeAllNineOnNewPage): Retrying 9-pic layout on new page (Page %d, Available H: %.2f)\n", e.currentPage.Page, newAvailableHeight)
	layoutInfo, err := e.calculateNinePicturesLayout(pictures, newAvailableHeight)

	if err == nil && layoutInfo.TotalHeight <= newAvailableHeight+1e-6 {
		fmt.Println("Debug (placeAllNineOnNewPage): Successfully placed all 9 on new page.")
		e.placePicturesInTemplate(pictures, layoutInfo)
		return layoutInfo.TotalHeight
	}

	// Failure even on a new page
	errMsg := "unknown error"
	if err != nil {
		errMsg = err.Error()
	} else if layoutInfo.TotalHeight > newAvailableHeight {
		errMsg = fmt.Sprintf("height %.2f exceeds available %.2f", layoutInfo.TotalHeight, newAvailableHeight)
	}
	fmt.Printf("Error (placeAllNineOnNewPage): Failed to place all 9 pictures even on a new page: %s. Returning 0.\n", errMsg)
	// Decide how to handle this critical failure. Returning 0 might lead to lost pictures.
	// Consider returning an error or a special value if the caller needs to know.
	// For now, returning 0 as per original structure's likely expectation.
	return 0
}

// --- Helper function dependencies (ensure these are available) ---
// func (e *ContinuousLayoutEngine) processPictureGroup(...) -> Needs access
// func (e *ContinuousLayoutEngine) calculateNinePicturesLayout(...) -> Needs access
// func (e *ContinuousLayoutEngine) processSixPicturesWithSplitLogic(...) -> Needs access
// func (e *ContinuousLayoutEngine) processThreePicturesWithSplitLogic(...) -> Needs access
// func removePicturesAdded(...) -> Needs definition or import
// func getEntryPictureCount(...) -> Needs definition or import
// func GetPictureType(...) -> Needs definition or import
// func GetRequiredMinHeight(...) -> Needs definition or import
// func (e *ContinuousLayoutEngine) placePicturesInTemplate(...) -> Needs access
// func (e *ContinuousLayoutEngine) newPage() -> Needs access
// func (e *ContinuousLayoutEngine) requiredSpacingBeforeElement() -> Needs access
