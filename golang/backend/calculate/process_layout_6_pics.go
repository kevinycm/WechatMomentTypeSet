package calculate

import "fmt"

// processLayoutForSixPictures handles the layout calculation and placement for exactly 6 pictures.
// It includes logic for handling split_required (to 3+3) and force_new_page signals.
func (e *ContinuousLayoutEngine) processLayoutForSixPictures(pictures []Picture, layoutAvailableHeight float64) float64 {
	numPics := 6
	if len(pictures) != numPics {
		fmt.Printf("Error: processLayoutForSixPictures called with %d pictures. Skipping.\n", len(pictures))
		return 0
	}

	// +++ Add Debug Logs for Case 6 +++
	fmt.Printf("Debug (Case 6 Entry): Entering with available height: %.2f\n", layoutAvailableHeight)
	layoutInfo, err := e.calculateSixPicturesLayout(pictures, layoutAvailableHeight)
	if err != nil {
		fmt.Printf("Debug (Case 6 Entry): calculateSixPicturesLayout returned error: %v\n", err)
	} else {
		fmt.Printf("Debug (Case 6 Entry): calculateSixPicturesLayout succeeded. TotalHeight: %.2f\n", layoutInfo.TotalHeight)
	}
	// +++ End Debug Logs +++

	// --- Special handling for 6-pic ---
	if err != nil {
		switch err.Error() {
		case "force_new_page":
			fmt.Println("Info: Forcing new page for 6-picture layout due to wide/tall images not fitting.")
			e.newPage()
			newLayoutAvailableHeight := (e.marginTop + e.availableHeight) - e.currentY // Recalculate for new page
			// Retry calculation on the new page
			fmt.Printf("Debug (Case 6 - Retry): Attempting calculation on new page. Available height: %.2f\n", newLayoutAvailableHeight)
			layoutInfo, err = e.calculateSixPicturesLayout(pictures, newLayoutAvailableHeight)
			if err != nil {
				fmt.Printf("Debug (Case 6 - Retry): Calculation failed even on new page: %v. Skipping.\n", err)
				return 0
			} else {
				fmt.Printf("Debug (Case 6 - Retry): Calculation succeeded on new page. TotalHeight: %.2f\n", layoutInfo.TotalHeight)
			}
			// If retry succeeds, err is nil, flow continues below

		case "split_required":
			fmt.Println("Info: Splitting 6-picture layout across pages (3+3).")
			// --- Refactored Logic for 6-Pic Split: Calculate first ---
			fmt.Printf("Debug (Case 6 - Split 3+3): Entering refactored logic. Initial available height: %.2f\n", layoutAvailableHeight)

			// 1. Try calculating layout for the first 3 pictures
			layoutInfo3, err3 := e.calculateThreePicturesLayout(pictures[0:3], layoutAvailableHeight)
			var heightUsed1 float64 = 0
			var spacingAdded float64 = 0

			// 2. Check if first group calculation succeeded AND fits in current available height
			if err3 == nil && layoutInfo3.TotalHeight <= layoutAvailableHeight {
				// Fits on current page/space
				fmt.Printf("Debug (Case 6 - Split 3+3): First 3 fit (height: %.2f <= available: %.2f)\n", layoutInfo3.TotalHeight, layoutAvailableHeight)
				e.placePicturesInTemplate(pictures[0:3], layoutInfo3)
				heightUsed1 = layoutInfo3.TotalHeight
				e.currentY += heightUsed1

				// Add spacing if possible
				spacingNeeded := e.elementSpacing
				spaceAfterFirstGroup := (e.marginTop + e.availableHeight) - e.currentY
				if spaceAfterFirstGroup >= spacingNeeded {
					fmt.Printf("Debug (Case 6 - Split 3+3): Adding spacing (%.2f) after first 3.\n", spacingNeeded)
					e.currentY += spacingNeeded
					spacingAdded = spacingNeeded
				} else {
					fmt.Printf("Debug (Case 6 - Split 3+3): Not enough space (%.2f) for spacing (%.2f).\n", spaceAfterFirstGroup, spacingNeeded)
				}
			} else {
				// Does not fit or calculation error - Force New Page
				if err3 != nil {
					fmt.Printf("Debug (Case 6 - Split 3+3): Initial calc for first 3 failed: %v. Forcing new page.\n", err3)
				} else {
					fmt.Printf("Debug (Case 6 - Split 3+3): First 3 don't fit (%.2f > %.2f). Forcing new page.\n", layoutInfo3.TotalHeight, layoutAvailableHeight)
				}
				e.newPage()
				layoutAvailableHeight = (e.marginTop + e.availableHeight) - e.currentY // Recalculate for new page

				// Retry calculating & placing first 3 on the new page
				fmt.Printf("Debug (Case 6 - Split 3+3): Retrying calc & place for first 3 on new page (available: %.2f)\n", layoutAvailableHeight)
				layoutInfo3New, err3New := e.calculateThreePicturesLayout(pictures[0:3], layoutAvailableHeight)
				if err3New != nil || layoutInfo3New.TotalHeight > layoutAvailableHeight {
					if err3New != nil {
						fmt.Printf("Error: Failed to place first 3 pictures during 6-pic split even on new page: %v\n", err3New)
					} else {
						fmt.Printf("Error: First 3 pictures too tall (%.2f) for new page (%.2f) during 6-pic split.\n", layoutInfo3New.TotalHeight, layoutAvailableHeight)
					}
					return 0 // Cannot place even the first group
				}

				// Place first 3 successfully on new page
				e.placePicturesInTemplate(pictures[0:3], layoutInfo3New)
				heightUsed1 = layoutInfo3New.TotalHeight
				e.currentY += heightUsed1

				// Add spacing if possible on new page
				spacingNeeded := e.elementSpacing
				spaceAfterFirstGroup := (e.marginTop + e.availableHeight) - e.currentY
				if spaceAfterFirstGroup >= spacingNeeded {
					fmt.Printf("Debug (Case 6 - Split 3+3): Adding spacing (%.2f) after first 3 on new page.\n", spacingNeeded)
					e.currentY += spacingNeeded
					spacingAdded = spacingNeeded
				} else {
					fmt.Printf("Debug (Case 6 - Split 3+3): Not enough space (%.2f) for spacing (%.2f) on new page.\n", spaceAfterFirstGroup, spacingNeeded)
				}
			}

			// 3. Place the remaining 3 pictures
			currentPageBeforeSecondGroup := e.currentPage.Page
			remainingHeightOnPage := (e.marginTop + e.availableHeight) - e.currentY
			if remainingHeightOnPage < 0 {
				remainingHeightOnPage = 0
			}
			fmt.Printf("Debug (Case 6 - Split 3+3): Placing remaining 3 (remaining H: %.2f)\n", remainingHeightOnPage)

			// Recursively call to handle the last 3 pictures - Use the specific 3-pic function
			heightUsed2 := e.processLayoutForThreePictures(pictures[3:6], remainingHeightOnPage)

			// Check if the second group placement failed INITIALLY
			if heightUsed2 <= 1e-6 {
				// --- Initial placement failed, try on new page ---
				fmt.Printf("Info (Case 6 - Split 3+3): Second group (pics %d-%d) failed placement in remaining space (%.2f). Attempting on new page.\n", pictures[3].Index, pictures[5].Index, remainingHeightOnPage)
				e.newPage()
				newLayoutAvailableHeight := (e.marginTop + e.availableHeight) - e.currentY
				fmt.Printf("Debug (Case 6 - Split 3+3): Retrying placement for second 3 on new page (available: %.2f)\n", newLayoutAvailableHeight)
				retryHeightUsed2 := e.processLayoutForThreePictures(pictures[3:6], newLayoutAvailableHeight)

				// Check if RETRY failed
				if retryHeightUsed2 <= 1e-6 {
					fmt.Println("Error: Failed to place second group of 3 pictures even on a new page during 6-pic split. Aborting split.")
					return 0 // Overall split failed
				} else {
					// RETRY SUCCEEDED on new page.
					// Return the height used by the FIRST group on the PREVIOUS page.
					fmt.Printf("Debug (Case 6 - Split 3+3): Second group placed on new page (height: %.2f). Returning height of first group only (%.2f + %.2f).\n", retryHeightUsed2, heightUsed1, spacingAdded)
					return heightUsed1 + spacingAdded
				}
			} else {
				// --- Initial placement succeeded ---
				// Second group placed successfully INITIALLY after first group (+ spacing) on the same page span
				lastPageAfterSecondGroup := e.currentPage.Page
				if lastPageAfterSecondGroup == currentPageBeforeSecondGroup {
					fmt.Printf("Debug (Case 6 - Split 3+3): Success, both groups on same page span (Page %d). Returning total height: %.2f (H1=%.2f + S=%.2f + H2=%.2f)\n", currentPageBeforeSecondGroup, heightUsed1+spacingAdded+heightUsed2, heightUsed1, spacingAdded, heightUsed2)
					return heightUsed1 + spacingAdded + heightUsed2 // Return combined height
				} else {
					// This case should technically not happen if processLayoutForThreePictures doesn't cause page breaks itself,
					// but as a safeguard, return the height used on the last page.
					fmt.Printf("Debug (Case 6 - Split 3+3): Success, second group possibly on new page (%d). Returning H2: %.2f\n", lastPageAfterSecondGroup, heightUsed2)
					return heightUsed2
				}
			}
			// --- End Refactored Logic for 6-Pic Split ---

		default:
			fmt.Printf("Error calculating layout for 6 pictures (unhandled error type): %v. Skipping placement.\n", err)
			return 0
		}
	}

	// If we reach here, err is nil, and layoutInfo is valid for placement.
	// +++ Log Before Placing +++
	fmt.Printf("Debug (Place - 6 Pics): Placing %d pictures. CurrentY before placement: %.2f. Layout TotalHeight: %.2f\n", len(pictures), e.currentY, layoutInfo.TotalHeight)
	// +++ End Log +++
	e.placePicturesInTemplate(pictures, layoutInfo)

	return layoutInfo.TotalHeight
}
