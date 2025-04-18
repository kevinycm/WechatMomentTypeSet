package waterfall

import (
	"fmt"
)

// processSixPicturesWithSplitLogic handles layout for 6 pictures based on the detailed 10-step rules.
// Returns height used ONLY if all 6 are placed together initially (Rule 1).
// Returns 0 in all other scenarios (splits, placements on new pages), caller relies on final e.currentY.
func (e *ContinuousLayoutEngine) processSixPicturesWithSplitLogic(pictures []Picture, layoutAvailableHeight float64) float64 {
	numPics := len(pictures)
	if numPics != 6 {
		fmt.Printf("Error (process6Split): Expected 6 pictures, got %d\n", numPics)
		return 0
	}

	const G1Start, G1End = 0, 2         // Group 1: Pics 0-1
	const G2Start, G2End = 2, 4         // Group 2: Pics 2-3
	const G3Start, G3End = 4, 6         // Group 3: Pics 4-5
	const G2FullStart, G2FullEnd = 2, 6 // Group 2 Full: Pics 2-5
	const tolerance = 1e-6

	// --- Rule 1: Try placing all 6 on the current page ---
	fmt.Printf("Debug (process6Split): Rule 1 - Attempting 6-pic layout on page %d (Avail H: %.2f).\n", e.currentPage.Page, layoutAvailableHeight)
	layoutInfo6, err6 := e.calculateSixPicturesLayout(pictures, layoutAvailableHeight)
	if err6 == nil && layoutInfo6.TotalHeight <= layoutAvailableHeight+tolerance {
		fmt.Printf("Debug (process6Split): Rule 1 - Success. Placing 6 pics (H: %.2f).\n", layoutInfo6.TotalHeight)
		e.placePicturesInTemplate(pictures, layoutInfo6)
		return layoutInfo6.TotalHeight // Return height used
	}
	fmt.Printf("Debug (process6Split): Rule 1 failed. Err: %v / Height: %.2f. Proceeding to Rule 2.\n", err6, layoutInfo6.TotalHeight)

	// --- Rule 2: Try placing Group 1 (0-1) on the current page ---
	fmt.Printf("Debug (process6Split): Rule 2 - Attempting G1 (0-1) on page %d (Avail H: %.2f).\n", e.currentPage.Page, layoutAvailableHeight)
	layoutInfoG1, errG1 := e.calculateTwoPicturesLayout(pictures[G1Start:G1End], layoutAvailableHeight)
	if errG1 == nil && layoutInfoG1.TotalHeight <= layoutAvailableHeight+tolerance {
		// Rule 2 Success Path: G1 fits on current page
		fmt.Printf("Debug (process6Split): Rule 2 - Success. Placing G1 (H: %.2f).\n", layoutInfoG1.TotalHeight)
		e.placePicturesInTemplate(pictures[G1Start:G1End], layoutInfoG1)
		e.currentY += layoutInfoG1.TotalHeight
		e.currentY += e.imageSpacing // Add spacing after G1
		currentAvailableHeightG1 := (e.marginTop + e.availableHeight) - e.currentY

		// Try placing Group 2 (2-3) on the same current page
		fmt.Printf("Debug (process6Split): Rule 2 - Attempting G2 (2-3) on same page %d (Avail H: %.2f).\n", e.currentPage.Page, currentAvailableHeightG1)
		layoutInfoG2, errG2 := e.calculateTwoPicturesLayout(pictures[G2Start:G2End], currentAvailableHeightG1)
		if errG2 == nil && layoutInfoG2.TotalHeight <= currentAvailableHeightG1+tolerance {
			// Rule 2 Success Path: G2 fits after G1 on current page
			fmt.Printf("Debug (process6Split): Rule 2 - Success. Placing G2 (H: %.2f).\n", layoutInfoG2.TotalHeight)
			e.placePicturesInTemplate(pictures[G2Start:G2End], layoutInfoG2)
			e.currentY += layoutInfoG2.TotalHeight
			e.currentY += e.imageSpacing // Add spacing after G2
			currentAvailableHeightG2 := (e.marginTop + e.availableHeight) - e.currentY

			// Try placing Group 3 (4-5) on the same current page
			fmt.Printf("Debug (process6Split): Rule 2 - Attempting G3 (4-5) on same page %d (Avail H: %.2f).\n", e.currentPage.Page, currentAvailableHeightG2)
			layoutInfoG3, errG3 := e.calculateTwoPicturesLayout(pictures[G3Start:G3End], currentAvailableHeightG2)
			if errG3 == nil && layoutInfoG3.TotalHeight <= currentAvailableHeightG2+tolerance {
				// Rule 2 Success Path: G3 fits after G2 on current page (2+2+2 success)
				fmt.Printf("Debug (process6Split): Rule 2 - Success. Placing G3 (H: %.2f). 2+2+2 on same page complete.\n", layoutInfoG3.TotalHeight)
				e.placePicturesInTemplate(pictures[G3Start:G3End], layoutInfoG3)
				e.currentY += layoutInfoG3.TotalHeight
				return 0
			} else {
				// Rule 2 Failed: G3 failed -> Rule 7: New page for G3
				fmt.Printf("Debug (process6Split): Rule 2 failed (G3 on same page). Err: %v / Height: %.2f. Proceeding to Rule 7 (New page for G3).\n", errG3, layoutInfoG3.TotalHeight)
				return e.placeG3OnNewPage(pictures[G3Start:G3End])
			}
		} else {
			// Rule 2 Failed: G2 failed -> Rule 5: New page, try G2-Full (2-5)
			fmt.Printf("Debug (process6Split): Rule 2 failed (G2 on same page). Err: %v / Height: %.2f. Proceeding to Rule 5 (New page, try G2-Full 2-5).\n", errG2, layoutInfoG2.TotalHeight)
			return e.processRule5Onwards(pictures[G2FullStart:G2FullEnd])
		}
	} else {
		// --- Rule 3: G1 failed on current page. New page, try 6-pic again. ---
		fmt.Printf("Debug (process6Split): Rule 2 failed (G1 on page %d). Err: %v / Height: %.2f. Proceeding to Rule 3 (New page).\n", e.currentPage.Page, errG1, layoutInfoG1.TotalHeight)
		e.newPage()
		e.currentY = e.marginTop
		newPageAvailableHeight1 := e.availableHeight

		fmt.Printf("Debug (process6Split): Rule 3 - Attempting 6-pic layout on new page %d (Avail H: %.2f).\n", e.currentPage.Page, newPageAvailableHeight1)
		layoutInfo6New, err6New := e.calculateSixPicturesLayout(pictures, newPageAvailableHeight1)
		if err6New == nil && layoutInfo6New.TotalHeight <= newPageAvailableHeight1+tolerance {
			// Rule 3 Success: 6 pics fit on the new page
			fmt.Printf("Debug (process6Split): Rule 3 - Success. Placing 6 pics (H: %.2f) on new page.\n", layoutInfo6New.TotalHeight)
			e.placePicturesInTemplate(pictures, layoutInfo6New)
			e.currentY += layoutInfo6New.TotalHeight
			return 0 // Return 0 as split across pages occurred
		} else {
			// --- Rule 4: 6-pic failed on new page. Try G1 (0-1) on new page. ---
			fmt.Printf("Debug (process6Split): Rule 3 failed (6-pic on new page). Err: %v / Height: %.2f. Proceeding to Rule 4.\n", err6New, layoutInfo6New.TotalHeight)
			// Note: We are still on the new page created.
			// Reset Y for placing G1 at the top of this new page.
			e.currentY = e.marginTop
			// Available height is still newPageAvailableHeight1.

			fmt.Printf("Debug (process6Split): Rule 4 - Attempting G1 (0-1) on new page %d (Avail H: %.2f).\n", e.currentPage.Page, newPageAvailableHeight1)
			layoutInfoG1New, errG1New := e.calculateTwoPicturesLayout(pictures[G1Start:G1End], newPageAvailableHeight1)
			if errG1New == nil && layoutInfoG1New.TotalHeight <= newPageAvailableHeight1+tolerance {
				// Rule 4 Success Path: G1 fits on new page
				fmt.Printf("Debug (process6Split): Rule 4 - Success. Placing G1 (H: %.2f) on new page.\n", layoutInfoG1New.TotalHeight)
				e.placePicturesInTemplate(pictures[G1Start:G1End], layoutInfoG1New)
				e.currentY += layoutInfoG1New.TotalHeight
				e.currentY += e.imageSpacing // Add spacing after G1
				newPageAvailableHeightG1 := (e.marginTop + e.availableHeight) - e.currentY

				// Try placing Group 2 (2-3) on the same new page
				fmt.Printf("Debug (process6Split): Rule 4 - Attempting G2 (2-3) on same new page %d (Avail H: %.2f).\n", e.currentPage.Page, newPageAvailableHeightG1)
				layoutInfoG2New, errG2New := e.calculateTwoPicturesLayout(pictures[G2Start:G2End], newPageAvailableHeightG1)
				if errG2New == nil && layoutInfoG2New.TotalHeight <= newPageAvailableHeightG1+tolerance {
					// Rule 4 Success Path: G2 fits after G1 on new page
					fmt.Printf("Debug (process6Split): Rule 4 - Success. Placing G2 (H: %.2f).\n", layoutInfoG2New.TotalHeight)
					e.placePicturesInTemplate(pictures[G2Start:G2End], layoutInfoG2New)
					e.currentY += layoutInfoG2New.TotalHeight
					e.currentY += e.imageSpacing // Add spacing after G2
					newPageAvailableHeightG2 := (e.marginTop + e.availableHeight) - e.currentY

					// Try placing Group 3 (4-5) on the same new page
					fmt.Printf("Debug (process6Split): Rule 4 - Attempting G3 (4-5) on same new page %d (Avail H: %.2f).\n", e.currentPage.Page, newPageAvailableHeightG2)
					layoutInfoG3New, errG3New := e.calculateTwoPicturesLayout(pictures[G3Start:G3End], newPageAvailableHeightG2)
					if errG3New == nil && layoutInfoG3New.TotalHeight <= newPageAvailableHeightG2+tolerance {
						// Rule 4 Success Path: G3 fits after G2 on new page (G1+G2+G3 success)
						fmt.Printf("Debug (process6Split): Rule 4 - Success. Placing G3 (H: %.2f). G1+G2+G3 on same new page complete.\n", layoutInfoG3New.TotalHeight)
						e.placePicturesInTemplate(pictures[G3Start:G3End], layoutInfoG3New)
						e.currentY += layoutInfoG3New.TotalHeight
						return 0
					} else {
						// Rule 4 Failed: G3 failed -> Rule 7: New page for G3
						fmt.Printf("Debug (process6Split): Rule 4 failed (G3 on same new page). Err: %v / Height: %.2f. Proceeding to Rule 7 (New page for G3).\n", errG3New, layoutInfoG3New.TotalHeight)
						return e.placeG3OnNewPage(pictures[G3Start:G3End])
					}
				} else {
					// Rule 4 Failed: G2 failed -> Rule 5: New page, try G2-Full (2-5)
					fmt.Printf("Debug (process6Split): Rule 4 failed (G2 on same new page). Err: %v / Height: %.2f. Proceeding to Rule 5 (New page, try G2-Full 2-5).\n", errG2New, layoutInfoG2New.TotalHeight)
					return e.processRule5Onwards(pictures[G2FullStart:G2FullEnd])
				}
			} else {
				// Rule 4 Failed Critically: G1 doesn't fit even on the new page.
				fmt.Printf("Warning (process6Split): Rule 4 - G1 (0-1) failed to place on new page %d (Avail H: %.2f). Err: %v / Height: %.2f. Proceeding to Rule 5 (New page, try G2-Full 2-5), G1 pics lost.\n", e.currentPage.Page, newPageAvailableHeight1, errG1New, layoutInfoG1New.TotalHeight)
				// Proceed to Rule 5, G1 is lost.
				return e.processRule5Onwards(pictures[G2FullStart:G2FullEnd])
			}
		}
	}
}

// processRule5Onwards handles Rule 5, 6, 7 logic: Try G2-Full(2-5) on new page,
// then G2(2-3)+G3(4-5) on another new page, then G3(4-5) on final new page.
func (e *ContinuousLayoutEngine) processRule5Onwards(picturesG2Full []Picture) float64 {
	const tolerance = 1e-6
	if len(picturesG2Full) != 4 {
		fmt.Printf("Error (process6Split Rule 5): Expected 4 pictures for G2-Full, got %d\n", len(picturesG2Full))
		return 0
	}

	// --- Rule 5: New page, try G2-Full (2-5) ---
	fmt.Printf("Debug (process6Split): Rule 5 - New page (Page %d) for G2-Full (2-5).\n", e.currentPage.Page+1)
	e.newPage()
	e.currentY = e.marginTop
	newPageAvailableHeight2 := e.availableHeight

	fmt.Printf("Debug (process6Split): Rule 5 - Attempting G2-Full (2-5) on new page %d (Avail H: %.2f).\n", e.currentPage.Page, newPageAvailableHeight2)
	layoutInfoG2Full, errG2Full := e.calculateFourPicturesLayout(picturesG2Full, newPageAvailableHeight2)
	if errG2Full == nil && layoutInfoG2Full.TotalHeight <= newPageAvailableHeight2+tolerance {
		// Rule 5 Success: G2-Full fits on its new page
		fmt.Printf("Debug (process6Split): Rule 5 - Success. Placing G2-Full (H: %.2f) on new page %d.\n", layoutInfoG2Full.TotalHeight, e.currentPage.Page)
		e.placePicturesInTemplate(picturesG2Full, layoutInfoG2Full)
		e.currentY += layoutInfoG2Full.TotalHeight
		return 0
	} else {
		// --- Rule 6: G2-Full failed on new page. Try G2 (2-3) on same new page. ---
		fmt.Printf("Debug (process6Split): Rule 5 failed (G2-Full on new page). Err: %v / Height: %.2f. Proceeding to Rule 6.\n", errG2Full, layoutInfoG2Full.TotalHeight)
		// Note: We are still on the page created for Rule 5.
		// Reset Y for placing G2 at the top of this page.
		e.currentY = e.marginTop
		// Available height is still newPageAvailableHeight2.

		fmt.Printf("Debug (process6Split): Rule 6 - Attempting G2 (2-3) on page %d (Avail H: %.2f).\n", e.currentPage.Page, newPageAvailableHeight2)
		layoutInfoG2, errG2 := e.calculateTwoPicturesLayout(picturesG2Full[0:2], newPageAvailableHeight2) // G2 pics are index 0,1 of picturesG2Full
		if errG2 == nil && layoutInfoG2.TotalHeight <= newPageAvailableHeight2+tolerance {
			// Rule 6 Success Path: G2 fits on this page
			fmt.Printf("Debug (process6Split): Rule 6 - Success. Placing G2 (H: %.2f).\n", layoutInfoG2.TotalHeight)
			e.placePicturesInTemplate(picturesG2Full[0:2], layoutInfoG2)
			e.currentY += layoutInfoG2.TotalHeight
			e.currentY += e.imageSpacing // Add spacing after G2
			currentAvailableHeightG2 := (e.marginTop + e.availableHeight) - e.currentY

			// Try placing Group 3 (4-5) on the same page
			fmt.Printf("Debug (process6Split): Rule 6 - Attempting G3 (4-5) on same page %d (Avail H: %.2f).\n", e.currentPage.Page, currentAvailableHeightG2)
			layoutInfoG3, errG3 := e.calculateTwoPicturesLayout(picturesG2Full[2:4], currentAvailableHeightG2) // G3 pics are index 2,3 of picturesG2Full
			if errG3 == nil && layoutInfoG3.TotalHeight <= currentAvailableHeightG2+tolerance {
				// Rule 6 Success Path: G3 fits after G2 (G2+G3 success)
				fmt.Printf("Debug (process6Split): Rule 6 - Success. Placing G3 (H: %.2f). G2+G3 on same page complete.\n", layoutInfoG3.TotalHeight)
				e.placePicturesInTemplate(picturesG2Full[2:4], layoutInfoG3)
				e.currentY += layoutInfoG3.TotalHeight
				return 0
			} else {
				// Rule 6 Failed: G3 failed -> Rule 7: New page for G3
				fmt.Printf("Debug (process6Split): Rule 6 failed (G3 on same page). Err: %v / Height: %.2f. Proceeding to Rule 7 (New page for G3).\n", errG3, layoutInfoG3.TotalHeight)
				return e.placeG3OnNewPage(picturesG2Full[2:4])
			}
		} else {
			// Rule 6 Failed Critically: G2 doesn't fit even on this page.
			fmt.Printf("Warning (process6Split): Rule 6 - G2 (2-3) failed to place on page %d (Avail H: %.2f). Err: %v / Height: %.2f. Proceeding to Rule 7 (New page for G3), G2 pics lost.\n", e.currentPage.Page, newPageAvailableHeight2, errG2, layoutInfoG2.TotalHeight)
			// Proceed to Rule 7 for G3, G2 is lost.
			return e.placeG3OnNewPage(picturesG2Full[2:4])
		}
	}
}

// placeG3OnNewPage handles Rule 7 logic: Create new page and place G3 (4-5).
func (e *ContinuousLayoutEngine) placeG3OnNewPage(picturesG3 []Picture) float64 {
	const tolerance = 1e-6
	if len(picturesG3) != 2 {
		fmt.Printf("Error (process6Split Rule 7): Expected 2 pictures for G3, got %d\n", len(picturesG3))
		return 0
	}

	// --- Rule 7: New page for G3 (4-5) ---
	fmt.Printf("Debug (process6Split): Rule 7 - New page (Page %d) for G3 (4-5).\n", e.currentPage.Page+1)
	e.newPage()
	e.currentY = e.marginTop
	newPageAvailableHeight3 := e.availableHeight

	fmt.Printf("Debug (process6Split): Rule 7 - Attempting G3 (4-5) on new page %d (Avail H: %.2f).\n", e.currentPage.Page, newPageAvailableHeight3)
	layoutInfoG3Final, errG3Final := e.calculateTwoPicturesLayout(picturesG3, newPageAvailableHeight3)
	if errG3Final == nil && layoutInfoG3Final.TotalHeight <= newPageAvailableHeight3+tolerance {
		fmt.Printf("Debug (process6Split): Rule 7 - Success. Placing G3 (H: %.2f) on new page %d.\n", layoutInfoG3Final.TotalHeight, e.currentPage.Page)
		e.placePicturesInTemplate(picturesG3, layoutInfoG3Final)
		e.currentY += layoutInfoG3Final.TotalHeight
	} else {
		// Rule 7 Failed: G3 failed even on its own dedicated page.
		fmt.Printf("Error (process6Split): Rule 7 - Critical failure. G3 (4-5) failed to place even on new page %d (Avail H: %.2f). Err: %v / Height: %.2f. Aborting placement of G3. Pics 4-5 lost.\n", e.currentPage.Page, newPageAvailableHeight3, errG3Final, layoutInfoG3Final.TotalHeight)
		// if errors.Is(errG3Final, ErrMinHeightConstraint) { ... }
	}
	return 0 // Always return 0 as split occurred.
}
