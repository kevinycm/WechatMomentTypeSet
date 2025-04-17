package calculate

import (
	"fmt"
)

// processEightPicturesWithSplitLogic handles layout for 8 pictures based on the detailed 10-step rules.
// Returns height used ONLY if all 8 are placed together initially (Rule 1).
// Returns 0 in all other scenarios (splits, placements on new pages), caller relies on final e.currentY.
func (e *ContinuousLayoutEngine) processEightPicturesWithSplitLogic(pictures []Picture, layoutAvailableHeight float64) float64 {
	numPics := len(pictures)
	if numPics != 8 {
		fmt.Printf("Error (process8Split): Expected 8 pictures, got %d\\n", numPics)
		return 0
	}

	const G1Start, G1End = 0, 2         // Group 1: Pics 0-1
	const G2Start, G2End = 2, 5         // Group 2: Pics 2-4 (3 pictures)
	const G3Start, G3End = 5, 8         // Group 3: Pics 5-7 (3 pictures)
	const G2FullStart, G2FullEnd = 2, 8 // Group 2 Full: Pics 2-7 (6 pictures)
	const tolerance = 1e-6

	// --- Rule 1: Try placing all 8 on the current page ---
	fmt.Printf("Debug (process8Split): Rule 1 - Attempting 8-pic layout on page %d (Avail H: %.2f).\\n", e.currentPage.Page, layoutAvailableHeight)
	layoutInfo8, err8 := e.calculateEightPicturesLayout(pictures, layoutAvailableHeight)
	if err8 == nil && layoutInfo8.TotalHeight <= layoutAvailableHeight+tolerance {
		fmt.Printf("Debug (process8Split): Rule 1 - Success. Placing 8 pics (H: %.2f).\\n", layoutInfo8.TotalHeight)
		e.placePicturesInTemplate(pictures, layoutInfo8)
		return layoutInfo8.TotalHeight // Return height used
	}
	// Check for force_new_page error from calculateEightPicturesLayout (specific rule for 8 pics)
	if err8 != nil && err8.Error() == "force_new_page" {
		fmt.Println("Debug (process8Split): Rule 1 calculation signaled force_new_page. Placing all 8 on new page.")
		return e.placeAllEightOnNewPage(pictures) // Use helper for retry logic
	}
	fmt.Printf("Debug (process8Split): Rule 1 failed. Err: %v / Height: %.2f. Proceeding to Rule 2.\\n", err8, layoutInfo8.TotalHeight)

	// --- Rule 2: Try placing Group 1 (0-1) on the current page ---
	fmt.Printf("Debug (process8Split): Rule 2 - Attempting G1 (0-1) on page %d (Avail H: %.2f).\\n", e.currentPage.Page, layoutAvailableHeight)
	layoutInfoG1, errG1 := e.calculateTwoPicturesLayout(pictures[G1Start:G1End], layoutAvailableHeight)
	if errG1 == nil && layoutInfoG1.TotalHeight <= layoutAvailableHeight+tolerance {
		// Rule 2 Success Path: G1 fits on current page
		fmt.Printf("Debug (process8Split): Rule 2 - Success. Placing G1 (H: %.2f).\\n", layoutInfoG1.TotalHeight)
		e.placePicturesInTemplate(pictures[G1Start:G1End], layoutInfoG1)
		e.currentY += layoutInfoG1.TotalHeight
		e.currentY += e.imageSpacing // Add spacing after G1
		currentAvailableHeightG1 := (e.marginTop + e.availableHeight) - e.currentY

		// Try placing Group 2 (2-4) on the same current page
		fmt.Printf("Debug (process8Split): Rule 2 - Attempting G2 (2-4) on same page %d (Avail H: %.2f).\\n", e.currentPage.Page, currentAvailableHeightG1)
		layoutInfoG2, errG2 := e.calculateThreePicturesLayout(pictures[G2Start:G2End], currentAvailableHeightG1)
		if errG2 == nil && layoutInfoG2.TotalHeight <= currentAvailableHeightG1+tolerance {
			// Rule 2 Success Path: G2 fits after G1 on current page
			fmt.Printf("Debug (process8Split): Rule 2 - Success. Placing G2 (H: %.2f).\\n", layoutInfoG2.TotalHeight)
			e.placePicturesInTemplate(pictures[G2Start:G2End], layoutInfoG2)
			e.currentY += layoutInfoG2.TotalHeight
			e.currentY += e.imageSpacing // Add spacing after G2
			currentAvailableHeightG2 := (e.marginTop + e.availableHeight) - e.currentY

			// Try placing Group 3 (5-7) on the same current page
			fmt.Printf("Debug (process8Split): Rule 2 - Attempting G3 (5-7) on same page %d (Avail H: %.2f).\\n", e.currentPage.Page, currentAvailableHeightG2)
			layoutInfoG3, errG3 := e.calculateThreePicturesLayout(pictures[G3Start:G3End], currentAvailableHeightG2)
			if errG3 == nil && layoutInfoG3.TotalHeight <= currentAvailableHeightG2+tolerance {
				// Rule 2 Success Path: G3 fits after G2 on current page (2+3+3 success)
				fmt.Printf("Debug (process8Split): Rule 2 - Success. Placing G3 (H: %.2f). 2+3+3 on same page complete.\\n", layoutInfoG3.TotalHeight)
				e.placePicturesInTemplate(pictures[G3Start:G3End], layoutInfoG3)
				e.currentY += layoutInfoG3.TotalHeight
				return 0
			} else {
				// Rule 2 Failed: G3 failed -> Rule 7: New page for G3
				fmt.Printf("Debug (process8Split): Rule 2 failed (G3 on same page). Err: %v / Height: %.2f. Proceeding to Rule 7 (New page for G3).\\n", errG3, layoutInfoG3.TotalHeight)
				return e.placeG3ThreeOnNewPage(pictures[G3Start:G3End])
			}
		} else {
			// Rule 2 Failed: G2 failed -> Rule 5: New page, try G2-Full (2-7)
			fmt.Printf("Debug (process8Split): Rule 2 failed (G2 on same page). Err: %v / Height: %.2f. Proceeding to Rule 5 (New page, try G2-Full 2-7).\\n", errG2, layoutInfoG2.TotalHeight)
			return e.processRule5OnwardsFor8Pics(pictures[G2FullStart:G2FullEnd])
		}
	} else {
		// --- Rule 3: G1 failed on current page. New page, try 8-pic again. ---
		fmt.Printf("Debug (process8Split): Rule 2 failed (G1 on page %d). Err: %v / Height: %.2f. Proceeding to Rule 3 (New page).\\n", e.currentPage.Page, errG1, layoutInfoG1.TotalHeight)
		e.newPage()
		e.currentY = e.marginTop
		newPageAvailableHeight1 := e.availableHeight

		fmt.Printf("Debug (process8Split): Rule 3 - Attempting 8-pic layout on new page %d (Avail H: %.2f).\\n", e.currentPage.Page, newPageAvailableHeight1)
		layoutInfo8New, err8New := e.calculateEightPicturesLayout(pictures, newPageAvailableHeight1)
		if err8New == nil && layoutInfo8New.TotalHeight <= newPageAvailableHeight1+tolerance {
			// Rule 3 Success: 8 pics fit on the new page
			fmt.Printf("Debug (process8Split): Rule 3 - Success. Placing 8 pics (H: %.2f) on new page.\\n", layoutInfo8New.TotalHeight)
			e.placePicturesInTemplate(pictures, layoutInfo8New)
			e.currentY += layoutInfo8New.TotalHeight
			return 0 // Return 0 as split across pages occurred
		} else {
			// --- Rule 4: 8-pic failed on new page. Try G1 (0-1) on new page. ---
			fmt.Printf("Debug (process8Split): Rule 3 failed (8-pic on new page). Err: %v / Height: %.2f. Proceeding to Rule 4.\\n", err8New, layoutInfo8New.TotalHeight)
			// Note: We are still on the new page created.
			// Reset Y for placing G1 at the top of this new page.
			e.currentY = e.marginTop
			// Available height is still newPageAvailableHeight1.

			fmt.Printf("Debug (process8Split): Rule 4 - Attempting G1 (0-1) on new page %d (Avail H: %.2f).\\n", e.currentPage.Page, newPageAvailableHeight1)
			layoutInfoG1New, errG1New := e.calculateTwoPicturesLayout(pictures[G1Start:G1End], newPageAvailableHeight1)
			if errG1New == nil && layoutInfoG1New.TotalHeight <= newPageAvailableHeight1+tolerance {
				// Rule 4 Success Path: G1 fits on new page
				fmt.Printf("Debug (process8Split): Rule 4 - Success. Placing G1 (H: %.2f) on new page.\\n", layoutInfoG1New.TotalHeight)
				e.placePicturesInTemplate(pictures[G1Start:G1End], layoutInfoG1New)
				e.currentY += layoutInfoG1New.TotalHeight
				e.currentY += e.imageSpacing // Add spacing after G1
				newPageAvailableHeightG1 := (e.marginTop + e.availableHeight) - e.currentY

				// Try placing Group 2 (2-4) on the same new page
				fmt.Printf("Debug (process8Split): Rule 4 - Attempting G2 (2-4) on same new page %d (Avail H: %.2f).\\n", e.currentPage.Page, newPageAvailableHeightG1)
				layoutInfoG2New, errG2New := e.calculateThreePicturesLayout(pictures[G2Start:G2End], newPageAvailableHeightG1)
				if errG2New == nil && layoutInfoG2New.TotalHeight <= newPageAvailableHeightG1+tolerance {
					// Rule 4 Success Path: G2 fits after G1 on new page
					fmt.Printf("Debug (process8Split): Rule 4 - Success. Placing G2 (H: %.2f).\\n", layoutInfoG2New.TotalHeight)
					e.placePicturesInTemplate(pictures[G2Start:G2End], layoutInfoG2New)
					e.currentY += layoutInfoG2New.TotalHeight
					e.currentY += e.imageSpacing // Add spacing after G2
					newPageAvailableHeightG2 := (e.marginTop + e.availableHeight) - e.currentY

					// Try placing Group 3 (5-7) on the same new page
					fmt.Printf("Debug (process8Split): Rule 4 - Attempting G3 (5-7) on same new page %d (Avail H: %.2f).\\n", e.currentPage.Page, newPageAvailableHeightG2)
					layoutInfoG3New, errG3New := e.calculateThreePicturesLayout(pictures[G3Start:G3End], newPageAvailableHeightG2)
					if errG3New == nil && layoutInfoG3New.TotalHeight <= newPageAvailableHeightG2+tolerance {
						// Rule 4 Success Path: G3 fits after G2 on new page (G1+G2+G3 success)
						fmt.Printf("Debug (process8Split): Rule 4 - Success. Placing G3 (H: %.2f). G1+G2+G3 on same new page complete.\\n", layoutInfoG3New.TotalHeight)
						e.placePicturesInTemplate(pictures[G3Start:G3End], layoutInfoG3New)
						e.currentY += layoutInfoG3New.TotalHeight
						return 0
					} else {
						// Rule 4 Failed: G3 failed -> Rule 7: New page for G3
						fmt.Printf("Debug (process8Split): Rule 4 failed (G3 on same new page). Err: %v / Height: %.2f. Proceeding to Rule 7 (New page for G3).\\n", errG3New, layoutInfoG3New.TotalHeight)
						return e.placeG3ThreeOnNewPage(pictures[G3Start:G3End])
					}
				} else {
					// Rule 4 Failed: G2 failed -> Rule 5: New page, try G2-Full (2-7)
					fmt.Printf("Debug (process8Split): Rule 4 failed (G2 on same new page). Err: %v / Height: %.2f. Proceeding to Rule 5 (New page, try G2-Full 2-7).\\n", errG2New, layoutInfoG2New.TotalHeight)
					return e.processRule5OnwardsFor8Pics(pictures[G2FullStart:G2FullEnd])
				}
			} else {
				// Rule 4 Failed Critically: G1 doesn't fit even on the new page.
				fmt.Printf("Warning (process8Split): Rule 4 - G1 (0-1) failed to place on new page %d (Avail H: %.2f). Err: %v / Height: %.2f. Proceeding to Rule 5 (New page, try G2-Full 2-7), G1 pics lost.\\n", e.currentPage.Page, newPageAvailableHeight1, errG1New, layoutInfoG1New.TotalHeight)
				// Proceed to Rule 5, G1 is lost.
				return e.processRule5OnwardsFor8Pics(pictures[G2FullStart:G2FullEnd])
			}
		}
	}
}

// processRule5OnwardsFor8Pics handles Rule 5, 6, 7 logic for 8 pictures.
func (e *ContinuousLayoutEngine) processRule5OnwardsFor8Pics(picturesG2Full []Picture) float64 {
	const tolerance = 1e-6
	if len(picturesG2Full) != 6 {
		fmt.Printf("Error (process8Split Rule 5): Expected 6 pictures for G2-Full, got %d\\n", len(picturesG2Full))
		return 0
	}

	// --- Rule 5: New page, try G2-Full (2-7) ---
	fmt.Printf("Debug (process8Split): Rule 5 - New page (Page %d) for G2-Full (2-7).\\n", e.currentPage.Page+1)
	e.newPage()
	e.currentY = e.marginTop
	newPageAvailableHeight2 := e.availableHeight

	fmt.Printf("Debug (process8Split): Rule 5 - Attempting G2-Full (2-7) on new page %d (Avail H: %.2f).\\n", e.currentPage.Page, newPageAvailableHeight2)
	// Use the appropriate calculate function for 6 pictures
	layoutInfoG2Full, errG2Full := e.calculateSixPicturesLayout(picturesG2Full, newPageAvailableHeight2)
	if errG2Full == nil && layoutInfoG2Full.TotalHeight <= newPageAvailableHeight2+tolerance {
		// Rule 5 Success: G2-Full fits on its new page
		fmt.Printf("Debug (process8Split): Rule 5 - Success. Placing G2-Full (H: %.2f) on new page %d.\\n", layoutInfoG2Full.TotalHeight, e.currentPage.Page)
		e.placePicturesInTemplate(picturesG2Full, layoutInfoG2Full)
		e.currentY += layoutInfoG2Full.TotalHeight
		return 0
	} else {
		// --- Rule 6: G2-Full failed on new page. Try G2 (2-4) on same new page. ---
		fmt.Printf("Debug (process8Split): Rule 5 failed (G2-Full on new page). Err: %v / Height: %.2f. Proceeding to Rule 6.\\n", errG2Full, layoutInfoG2Full.TotalHeight)
		// Note: We are still on the page created for Rule 5.
		// Reset Y for placing G2 at the top of this page.
		e.currentY = e.marginTop
		// Available height is still newPageAvailableHeight2.

		fmt.Printf("Debug (process8Split): Rule 6 - Attempting G2 (2-4) on page %d (Avail H: %.2f).\\n", e.currentPage.Page, newPageAvailableHeight2)
		layoutInfoG2, errG2 := e.calculateThreePicturesLayout(picturesG2Full[0:3], newPageAvailableHeight2) // G2 pics are index 0,1,2 of picturesG2Full
		if errG2 == nil && layoutInfoG2.TotalHeight <= newPageAvailableHeight2+tolerance {
			// Rule 6 Success Path: G2 fits on this page
			fmt.Printf("Debug (process8Split): Rule 6 - Success. Placing G2 (H: %.2f).\\n", layoutInfoG2.TotalHeight)
			e.placePicturesInTemplate(picturesG2Full[0:3], layoutInfoG2)
			e.currentY += layoutInfoG2.TotalHeight
			e.currentY += e.imageSpacing // Add spacing after G2
			currentAvailableHeightG2 := (e.marginTop + e.availableHeight) - e.currentY

			// Try placing Group 3 (5-7) on the same page
			fmt.Printf("Debug (process8Split): Rule 6 - Attempting G3 (5-7) on same page %d (Avail H: %.2f).\\n", e.currentPage.Page, currentAvailableHeightG2)
			layoutInfoG3, errG3 := e.calculateThreePicturesLayout(picturesG2Full[3:6], currentAvailableHeightG2) // G3 pics are index 3,4,5 of picturesG2Full
			if errG3 == nil && layoutInfoG3.TotalHeight <= currentAvailableHeightG2+tolerance {
				// Rule 6 Success Path: G3 fits after G2 (G2+G3 success)
				fmt.Printf("Debug (process8Split): Rule 6 - Success. Placing G3 (H: %.2f). G2+G3 on same page complete.\\n", layoutInfoG3.TotalHeight)
				e.placePicturesInTemplate(picturesG2Full[3:6], layoutInfoG3)
				e.currentY += layoutInfoG3.TotalHeight
				return 0
			} else {
				// Rule 6 Failed: G3 failed -> Rule 7: New page for G3
				fmt.Printf("Debug (process8Split): Rule 6 failed (G3 on same page). Err: %v / Height: %.2f. Proceeding to Rule 7 (New page for G3).\\n", errG3, layoutInfoG3.TotalHeight)
				return e.placeG3ThreeOnNewPage(picturesG2Full[3:6])
			}
		} else {
			// Rule 6 Failed Critically: G2 doesn't fit even on this page.
			fmt.Printf("Warning (process8Split): Rule 6 - G2 (2-4) failed to place on page %d (Avail H: %.2f). Err: %v / Height: %.2f. Proceeding to Rule 7 (New page for G3), G2 pics lost.\\n", e.currentPage.Page, newPageAvailableHeight2, errG2, layoutInfoG2.TotalHeight)
			// Proceed to Rule 7 for G3, G2 is lost.
			return e.placeG3ThreeOnNewPage(picturesG2Full[3:6])
		}
	}
}

// placeG3ThreeOnNewPage handles Rule 7 logic for 8 pictures: Create new page and place G3 (5-7).
func (e *ContinuousLayoutEngine) placeG3ThreeOnNewPage(picturesG3 []Picture) float64 {
	const tolerance = 1e-6
	if len(picturesG3) != 3 {
		fmt.Printf("Error (process8Split Rule 7): Expected 3 pictures for G3, got %d\\n", len(picturesG3))
		return 0
	}

	// --- Rule 7: New page for G3 (5-7) ---
	fmt.Printf("Debug (process8Split): Rule 7 - New page (Page %d) for G3 (5-7).\\n", e.currentPage.Page+1)
	e.newPage()
	e.currentY = e.marginTop
	newPageAvailableHeight3 := e.availableHeight

	fmt.Printf("Debug (process8Split): Rule 7 - Attempting G3 (5-7) on new page %d (Avail H: %.2f).\\n", e.currentPage.Page, newPageAvailableHeight3)
	layoutInfoG3Final, errG3Final := e.calculateThreePicturesLayout(picturesG3, newPageAvailableHeight3)
	if errG3Final == nil && layoutInfoG3Final.TotalHeight <= newPageAvailableHeight3+tolerance {
		fmt.Printf("Debug (process8Split): Rule 7 - Success. Placing G3 (H: %.2f) on new page %d.\\n", layoutInfoG3Final.TotalHeight, e.currentPage.Page)
		e.placePicturesInTemplate(picturesG3, layoutInfoG3Final)
		e.currentY += layoutInfoG3Final.TotalHeight
	} else {
		// Rule 7 Failed: G3 failed even on its own dedicated page.
		fmt.Printf("Error (process8Split): Rule 7 - Critical failure. G3 (5-7) failed to place even on new page %d (Avail H: %.2f). Err: %v / Height: %.2f. Aborting placement of G3. Pics 5-7 lost.\\n", e.currentPage.Page, newPageAvailableHeight3, errG3Final, layoutInfoG3Final.TotalHeight)
		// if errors.Is(errG3Final, ErrMinHeightConstraint) { ... }
	}
	return 0 // Always return 0 as split occurred.
}

// placeAllEightOnNewPage helper (used in Rule 1 fallback)
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
		fmt.Printf("Error (process8Split - Fallback): Failed to place all 8 even on new page (err: %v). Aborting.\\n", err8Retry)
		return 0
	}
}
