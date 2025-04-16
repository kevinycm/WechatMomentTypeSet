package calculate

import "fmt"

// processLayoutForNinePictures handles the layout calculation and placement for exactly 9 pictures.
// It includes logic for handling split_required (to 3+6) and force_new_page signals.
func (e *ContinuousLayoutEngine) processLayoutForNinePictures(pictures []Picture, layoutAvailableHeight float64) float64 {
	numPics := 9
	if len(pictures) != numPics {
		fmt.Printf("Error: processLayoutForNinePictures called with %d pictures. Skipping.\n", len(pictures))
		return 0
	}

	layoutInfo, err := e.calculateNinePicturesLayout(pictures, layoutAvailableHeight)

	// --- Special handling for 9-pic ---
	if err != nil {
		switch err.Error() {
		case "force_new_page":
			fmt.Println("Info: Forcing new page for 9-picture layout due to wide/tall images not fitting.")
			e.newPage()
			newLayoutAvailableHeight := (e.marginTop + e.availableHeight) - e.currentY // Recalculate
			// Retry calculation on the new page
			fmt.Printf("Debug (Case 9 - Retry): Attempting calculation on new page. Available height: %.2f\n", newLayoutAvailableHeight)
			layoutInfo, err = e.calculateNinePicturesLayout(pictures, newLayoutAvailableHeight)
			if err != nil {
				// +++ Log Retry Failure +++
				fmt.Printf("Error (Case 9 - Retry Failed): Calculation on new page failed with error: %v. Skipping placement.\n", err)
				// +++ End Log +++
				// fmt.Printf("Error calculating 9-pic layout even on new page: %v. Skipping.\n", err) // Old log
				return 0
			}
			// +++ Log Retry Success +++
			fmt.Printf("Debug (Case 9 - Retry Succeeded): Calculation on new page successful. Layout Height: %.2f\n", layoutInfo.TotalHeight)
			// +++ End Log +++
			// If retry succeeds, err is nil...

		case "split_required":
			fmt.Println("Info: Splitting 9-picture layout across pages (3+6).")
			// --- Refactored Logic: Calculate first, then decide page break ---

			// 1. Try calculating layout for the first 3 pictures on the current page
			fmt.Printf("Debug (9-Split): Attempting calculation for first 3 (available: %.2f)\n", layoutAvailableHeight)
			layoutInfo3, err3 := e.calculateThreePicturesLayout(pictures[0:3], layoutAvailableHeight)
			var heightUsed1 float64 = 0

			// 2. Check if first group calculation succeeded AND fits
			if err3 == nil && layoutInfo3.TotalHeight <= layoutAvailableHeight {
				// Fits on current page
				fmt.Printf("Debug (9-Split): First 3 fit on current page (height: %.2f <= available: %.2f)\n", layoutInfo3.TotalHeight, layoutAvailableHeight)
				e.placePicturesInTemplate(pictures[0:3], layoutInfo3)
				heightUsed1 = layoutInfo3.TotalHeight
				e.currentY += heightUsed1 // Update Y on current page

				// Add spacing if possible
				spacingNeeded := e.elementSpacing
				spaceAfterFirstGroup := (e.marginTop + e.availableHeight) - e.currentY
				if spaceAfterFirstGroup >= spacingNeeded {
					fmt.Printf("Debug (9-Split): Adding spacing (%.2f) after first 3 on same page.\n", spacingNeeded)
					e.currentY += spacingNeeded
				} else {
					fmt.Printf("Debug (9-Split): Not enough space (%.2f) for spacing (%.2f).\n", spaceAfterFirstGroup, spacingNeeded)
				}
			} else {
				// Does not fit or calculation error - Force New Page
				if err3 != nil {
					fmt.Printf("Debug (9-Split): Initial calc for first 3 failed: %v. Forcing new page.\n", err3)
				} else {
					fmt.Printf("Debug (9-Split): First 3 don't fit (%.2f > %.2f). Forcing new page.\n", layoutInfo3.TotalHeight, layoutAvailableHeight)
				}
				e.newPage()                                                            // Go to new page
				layoutAvailableHeight = (e.marginTop + e.availableHeight) - e.currentY // Recalculate for new page

				// Retry calculating & placing first 3 on the new page
				fmt.Printf("Debug (9-Split): Retrying calculation & placement for first 3 on new page (available: %.2f)\n", layoutAvailableHeight)
				layoutInfo3New, err3New := e.calculateThreePicturesLayout(pictures[0:3], layoutAvailableHeight)
				if err3New != nil || layoutInfo3New.TotalHeight > layoutAvailableHeight {
					// Handle error: Failed even on new page
					if err3New != nil {
						fmt.Printf("Error: Failed to place first 3 (9-split) even on new page: %v\n", err3New)
					} else {
						fmt.Printf("Error: First 3 (9-split) too tall (%.2f) for new page (%.2f).\n", layoutInfo3New.TotalHeight, layoutAvailableHeight)
					}
					return 0
				}
				// Place first 3 successfully on new page
				e.placePicturesInTemplate(pictures[0:3], layoutInfo3New)
				heightUsed1 = layoutInfo3New.TotalHeight
				e.currentY += heightUsed1 // Update Y on new page

				// Add spacing if possible on new page
				spacingNeeded := e.elementSpacing
				spaceAfterFirstGroup := (e.marginTop + e.availableHeight) - e.currentY
				if spaceAfterFirstGroup >= spacingNeeded {
					fmt.Printf("Debug (9-Split): Adding spacing (%.2f) after first 3 on new page.\n", spacingNeeded)
					e.currentY += spacingNeeded
				} else {
					fmt.Printf("Debug (9-Split): Not enough space (%.2f) for spacing (%.2f) on new page.\n", spaceAfterFirstGroup, spacingNeeded)
				}
			}

			// --- Force New Page for the second group ---
			fmt.Printf("Debug (9-Split): Forcing new page after placing first 3 pictures (Y was %.2f).\n", e.currentY)
			e.newPage()

			// 3. Place the remaining 6 pictures on the new page
			currentPageNum := e.currentPage.Page
			layoutAvailableHeightStep2 := (e.marginTop + e.availableHeight) - e.currentY // Height available on the new page
			fmt.Printf("Debug (9-Split): Placing remaining 6 (pics %d-%d) on New Page %d. Available height: %.2f\n", pictures[3].Index, pictures[8].Index, currentPageNum, layoutAvailableHeightStep2)

			heightUsed2 := e.processLayoutForSixPictures(pictures[3:9], layoutAvailableHeightStep2)
			if heightUsed2 <= 1e-6 {
				fmt.Println("Error: Failed to place remaining 6 pictures during 9-pic split on new page. First part placed, but operation incomplete.")
				return 0 // Return 0 as the second part failed
			}
			// If the second group (6 pics) was placed successfully on the new page,
			// return the height it consumed on that new page.
			// processPictures will correctly add this height to the currentY of the new page.
			return heightUsed2

		default:
			fmt.Printf("Error calculating layout for 9 pictures: %v. Skipping placement.\n", err)
			return 0
		}
	}

	// If we reach here, err is nil, and layoutInfo is valid for placement.
	// +++ Log Before Placing +++
	fmt.Printf("Debug (Place - 9 Pics): Placing %d pictures. CurrentY before placement: %.2f. Layout TotalHeight: %.2f\n", len(pictures), e.currentY, layoutInfo.TotalHeight)
	// +++ End Log +++
	e.placePicturesInTemplate(pictures, layoutInfo)

	return layoutInfo.TotalHeight
}
