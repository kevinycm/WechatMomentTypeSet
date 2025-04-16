package calculate

import "fmt"

// processLayoutForEightPictures handles the layout calculation and placement for exactly 8 pictures.
// It includes logic for handling split_required (to 4+4) and force_new_page signals.
func (e *ContinuousLayoutEngine) processLayoutForEightPictures(pictures []Picture, layoutAvailableHeight float64) float64 {
	numPics := 8
	if len(pictures) != numPics {
		fmt.Printf("Error: processLayoutForEightPictures called with %d pictures. Skipping.\n", len(pictures))
		return 0
	}

	layoutInfo, err := e.calculateEightPicturesLayout(pictures, layoutAvailableHeight)

	// --- Special handling for 8-pic ---
	if err != nil {
		switch err.Error() {
		case "force_new_page":
			fmt.Println("Info: Forcing new page for 8-picture layout due to wide/tall images not fitting.")
			e.newPage()
			newLayoutAvailableHeight := (e.marginTop + e.availableHeight) - e.currentY // Recalculate
			// Retry calculation on the new page
			fmt.Printf("Debug (Case 8 - Retry): Attempting calculation on new page. Available height: %.2f\n", newLayoutAvailableHeight)
			layoutInfo, err = e.calculateEightPicturesLayout(pictures, newLayoutAvailableHeight)
			if err != nil {
				fmt.Printf("Error calculating 8-pic layout even on new page: %v. Skipping.\n", err)
				return 0
			}
			// If retry succeeds, err is nil, flow continues below

		case "split_required":
			fmt.Println("Info: Splitting 8-picture layout into 4+4.")
			// --- Refactored Logic for 8-Pic Split: Calculate first ---
			fmt.Printf("Debug (Case 8 - Split 4+4): Entering refactored logic. Initial available height: %.2f\n", layoutAvailableHeight)

			// 1. Try calculating layout for the first 4 pictures
			layoutInfo4, err4 := e.calculateFourPicturesLayout(pictures[0:4], layoutAvailableHeight)
			var heightUsed1 float64 = 0
			var spacingAdded float64 = 0

			// 2. Check if first group calculation succeeded AND fits
			if err4 == nil && layoutInfo4.TotalHeight <= layoutAvailableHeight {
				// Fits on current page/space
				fmt.Printf("Debug (Case 8 - Split 4+4): First 4 fit (height: %.2f <= available: %.2f)\n", layoutInfo4.TotalHeight, layoutAvailableHeight)
				e.placePicturesInTemplate(pictures[0:4], layoutInfo4)
				heightUsed1 = layoutInfo4.TotalHeight
				e.currentY += heightUsed1

				// Add spacing if possible
				spacingNeeded := e.elementSpacing
				spaceAfterFirstGroup := (e.marginTop + e.availableHeight) - e.currentY
				if spaceAfterFirstGroup >= spacingNeeded {
					fmt.Printf("Debug (Case 8 - Split 4+4): Adding spacing (%.2f) after first 4.\n", spacingNeeded)
					e.currentY += spacingNeeded
					spacingAdded = spacingNeeded
				} else {
					fmt.Printf("Debug (Case 8 - Split 4+4): Not enough space (%.2f) for spacing (%.2f).\n", spaceAfterFirstGroup, spacingNeeded)
				}
			} else {
				// Does not fit or calculation error - Force New Page
				if err4 != nil {
					fmt.Printf("Debug (Case 8 - Split 4+4): Initial calc for first 4 failed: %v. Forcing new page.\n", err4)
				} else {
					fmt.Printf("Debug (Case 8 - Split 4+4): First 4 don't fit (%.2f > %.2f). Forcing new page.\n", layoutInfo4.TotalHeight, layoutAvailableHeight)
				}
				e.newPage()
				layoutAvailableHeight = (e.marginTop + e.availableHeight) - e.currentY // Recalculate for new page

				// Retry calculating & placing first 4 on the new page
				fmt.Printf("Debug (Case 8 - Split 4+4): Retrying calc & place for first 4 on new page (available: %.2f)\n", layoutAvailableHeight)
				layoutInfo4New, err4New := e.calculateFourPicturesLayout(pictures[0:4], layoutAvailableHeight)
				if err4New != nil || layoutInfo4New.TotalHeight > layoutAvailableHeight {
					if err4New != nil {
						fmt.Printf("Error: Failed to place first 4 pictures (8-split) even on new page: %v\n", err4New)
					} else {
						fmt.Printf("Error: First 4 pictures (8-split) too tall (%.2f) for new page (%.2f).\n", layoutInfo4New.TotalHeight, layoutAvailableHeight)
					}
					return 0 // Cannot place even the first group
				}
				// Place first 4 successfully on new page
				e.placePicturesInTemplate(pictures[0:4], layoutInfo4New)
				heightUsed1 = layoutInfo4New.TotalHeight
				e.currentY += heightUsed1

				// Add spacing if possible on new page
				spacingNeeded := e.elementSpacing
				spaceAfterFirstGroup := (e.marginTop + e.availableHeight) - e.currentY
				if spaceAfterFirstGroup >= spacingNeeded {
					fmt.Printf("Debug (Case 8 - Split 4+4): Adding spacing (%.2f) after first 4 on new page.\n", spacingNeeded)
					e.currentY += spacingNeeded
					spacingAdded = spacingNeeded
				} else {
					fmt.Printf("Debug (Case 8 - Split 4+4): Not enough space (%.2f) for spacing (%.2f) on new page.\n", spaceAfterFirstGroup, spacingNeeded)
				}
			}

			// 3. Place the remaining 4 pictures
			currentPageBeforeSecondGroup := e.currentPage.Page
			remainingHeightOnPage := (e.marginTop + e.availableHeight) - e.currentY
			if remainingHeightOnPage < 0 {
				remainingHeightOnPage = 0
			}
			fmt.Printf("Debug (Case 8 - Split 4+4): Placing remaining 4 (remaining H: %.2f)\n", remainingHeightOnPage)

			// Recursively call to handle the last 4 pictures - Use the specific 4-pic function
			heightUsed2 := e.processLayoutForFourPictures(pictures[4:8], remainingHeightOnPage)

			// Check if the second group placement failed INITIALLY
			if heightUsed2 <= 1e-6 {
				// --- Initial placement failed, try on new page ---
				fmt.Printf("Info (Case 8 - Split 4+4): Second group (pics %d-%d) failed placement in remaining space (%.2f). Attempting on new page.\n", pictures[4].Index, pictures[7].Index, remainingHeightOnPage)
				e.newPage()
				newLayoutAvailableHeight := (e.marginTop + e.availableHeight) - e.currentY
				fmt.Printf("Debug (Case 8 - Split 4+4): Retrying placement for second 4 on new page (available: %.2f)\n", newLayoutAvailableHeight)
				retryHeightUsed2 := e.processLayoutForFourPictures(pictures[4:8], newLayoutAvailableHeight)

				// Check if RETRY failed
				if retryHeightUsed2 <= 1e-6 {
					fmt.Println("Error: Failed to place second group of 4 pictures even on a new page during 8-pic split. Aborting split.")
					return 0 // Overall split failed
				} else {
					// RETRY SUCCEEDED on new page.
					// Return the height used by the FIRST group on the PREVIOUS page.
					fmt.Printf("Debug (Case 8 - Split 4+4): Second group placed on new page (height: %.2f). Returning height of first group only (%.2f + %.2f).\n", retryHeightUsed2, heightUsed1, spacingAdded)
					return heightUsed1 + spacingAdded
				}
			} else {
				// --- Initial placement succeeded ---
				// Second group placed successfully INITIALLY after first group (+ spacing) on the same page span
				lastPageAfterSecondGroup := e.currentPage.Page
				if lastPageAfterSecondGroup == currentPageBeforeSecondGroup {
					fmt.Printf("Debug (Case 8 - Split 4+4): Success, both groups on same page span (Page %d). Returning total height: %.2f (H1=%.2f + S=%.2f + H2=%.2f)\n", currentPageBeforeSecondGroup, heightUsed1+spacingAdded+heightUsed2, heightUsed1, spacingAdded, heightUsed2)
					return heightUsed1 + spacingAdded + heightUsed2 // Return combined height
				} else {
					// This case implies the 4-pic handler internally caused a page break.
					fmt.Printf("Debug (Case 8 - Split 4+4): Success, second group possibly on new page (%d). Returning H2: %.2f\n", lastPageAfterSecondGroup, heightUsed2)
					return heightUsed2
				}
			}
			// --- End Refactored Logic for 8-Pic Split ---

		default:
			fmt.Printf("Error calculating layout for 8 pictures: %v. Skipping placement.\n", err)
			return 0
		}
	}

	// If we reach here, err is nil, and layoutInfo is valid for placement.
	// +++ Log Before Placing +++
	fmt.Printf("Debug (Place - 8 Pics): Placing %d pictures. CurrentY before placement: %.2f. Layout TotalHeight: %.2f\n", len(pictures), e.currentY, layoutInfo.TotalHeight)
	// +++ End Log +++
	e.placePicturesInTemplate(pictures, layoutInfo)

	return layoutInfo.TotalHeight
}
