package waterfall

import (
	"fmt"
)

// processFourPicturesWithSplitLogic handles layout for 4 pictures based on the detailed 5-step rules.
// Returns height used ONLY if all 4 are placed together initially (Rule 1).
// Returns 0 in all other scenarios (splits, placements on new pages), caller relies on final e.currentY.
func (e *ContinuousLayoutEngine) processFourPicturesWithSplitLogic(pictures []Picture, layoutAvailableHeight float64) float64 {
	numPics := len(pictures)
	if numPics != 4 {
		fmt.Printf("Error (process4Split): Expected 4 pictures, got %d\n", numPics)
		return 0
	}

	const G1Start, G1End = 0, 2 // Group 1: Pics 0-1
	const G2Start, G2End = 2, 4 // Group 2: Pics 2-3
	const tolerance = 1e-6

	// --- Rule 1: Try placing all 4 on the current page ---
	fmt.Printf("Debug (process4Split): Rule 1 - Attempting 4-pic layout on page %d (Avail H: %.2f).\n", e.currentPage.Page, layoutAvailableHeight)
	layoutInfo4, err4 := e.calculateFourPicturesLayout(pictures, layoutAvailableHeight)
	if err4 == nil && layoutInfo4.TotalHeight <= layoutAvailableHeight+tolerance {
		fmt.Printf("Debug (process4Split): Rule 1 - Success. Placing 4 pics (H: %.2f).\n", layoutInfo4.TotalHeight)
		e.placePicturesInTemplate(pictures, layoutInfo4)
		return layoutInfo4.TotalHeight // Return height used
	}
	fmt.Printf("Debug (process4Split): Rule 1 failed. Err: %v / Height: %.2f. Proceeding to Rule 2.\n", err4, layoutInfo4.TotalHeight)

	// --- Rule 2: Try placing Group 1 (0-1) on the current page ---
	fmt.Printf("Debug (process4Split): Rule 2 - Attempting G1 (0-1) on page %d (Avail H: %.2f).\n", e.currentPage.Page, layoutAvailableHeight)
	layoutInfoG1, errG1 := e.calculateTwoPicturesLayout(pictures[G1Start:G1End], layoutAvailableHeight)
	if errG1 == nil && layoutInfoG1.TotalHeight <= layoutAvailableHeight+tolerance {
		// Rule 2 Success Path: G1 fits on current page
		fmt.Printf("Debug (process4Split): Rule 2 - Success. Placing G1 (H: %.2f).\n", layoutInfoG1.TotalHeight)
		e.placePicturesInTemplate(pictures[G1Start:G1End], layoutInfoG1)
		e.currentY += layoutInfoG1.TotalHeight
		e.currentY += e.imageSpacing // Add spacing after G1
		currentAvailableHeightG1 := (e.marginTop + e.availableHeight) - e.currentY

		// Try placing Group 2 (2-3) on the same current page
		fmt.Printf("Debug (process4Split): Rule 2 - Attempting G2 (2-3) on same page %d (Avail H: %.2f).\n", e.currentPage.Page, currentAvailableHeightG1)
		layoutInfoG2, errG2 := e.calculateTwoPicturesLayout(pictures[G2Start:G2End], currentAvailableHeightG1)
		if errG2 == nil && layoutInfoG2.TotalHeight <= currentAvailableHeightG1+tolerance {
			// Rule 2 Success Path: G2 fits after G1 on current page (2+2 success)
			fmt.Printf("Debug (process4Split): Rule 2 - Success. Placing G2 (H: %.2f). 2+2 on same page complete.\n", layoutInfoG2.TotalHeight)
			e.placePicturesInTemplate(pictures[G2Start:G2End], layoutInfoG2)
			e.currentY += layoutInfoG2.TotalHeight
			return 0
		} else {
			// Rule 2 Failed: G2 failed -> Rule 5: New page for G2
			fmt.Printf("Debug (process4Split): Rule 2 failed (G2 on same page). Err: %v / Height: %.2f. Proceeding to Rule 5 (New page for G2).\n", errG2, layoutInfoG2.TotalHeight)
			return e.placeG2TwoOnNewPage(pictures[G2Start:G2End])
		}
	} else {
		// --- Rule 3: G1 failed on current page. New page, try 4-pic again. ---
		fmt.Printf("Debug (process4Split): Rule 2 failed (G1 on page %d). Err: %v / Height: %.2f. Proceeding to Rule 3 (New page).\n", e.currentPage.Page, errG1, layoutInfoG1.TotalHeight)
		e.newPage()
		e.currentY = e.marginTop
		newPageAvailableHeight1 := e.availableHeight

		fmt.Printf("Debug (process4Split): Rule 3 - Attempting 4-pic layout on new page %d (Avail H: %.2f).\n", e.currentPage.Page, newPageAvailableHeight1)
		layoutInfo4New, err4New := e.calculateFourPicturesLayout(pictures, newPageAvailableHeight1)
		if err4New == nil && layoutInfo4New.TotalHeight <= newPageAvailableHeight1+tolerance {
			// Rule 3 Success: 4 pics fit on the new page
			fmt.Printf("Debug (process4Split): Rule 3 - Success. Placing 4 pics (H: %.2f) on new page.\n", layoutInfo4New.TotalHeight)
			e.placePicturesInTemplate(pictures, layoutInfo4New)
			e.currentY += layoutInfo4New.TotalHeight
			return 0 // Return 0 as split across pages occurred
		} else {
			// --- Rule 4: 4-pic failed on new page. Try G1 (0-1) on new page. ---
			fmt.Printf("Debug (process4Split): Rule 3 failed (4-pic on new page). Err: %v / Height: %.2f. Proceeding to Rule 4.\n", err4New, layoutInfo4New.TotalHeight)
			// Note: We are still on the new page created.
			// Reset Y for placing G1 at the top of this new page.
			e.currentY = e.marginTop
			// Available height is still newPageAvailableHeight1.

			fmt.Printf("Debug (process4Split): Rule 4 - Attempting G1 (0-1) on new page %d (Avail H: %.2f).\n", e.currentPage.Page, newPageAvailableHeight1)
			layoutInfoG1New, errG1New := e.calculateTwoPicturesLayout(pictures[G1Start:G1End], newPageAvailableHeight1)
			if errG1New == nil && layoutInfoG1New.TotalHeight <= newPageAvailableHeight1+tolerance {
				// Rule 4 Success Path: G1 fits on new page
				fmt.Printf("Debug (process4Split): Rule 4 - Success. Placing G1 (H: %.2f) on new page.\n", layoutInfoG1New.TotalHeight)
				e.placePicturesInTemplate(pictures[G1Start:G1End], layoutInfoG1New)
				e.currentY += layoutInfoG1New.TotalHeight
				e.currentY += e.imageSpacing // Add spacing after G1
				newPageAvailableHeightG1 := (e.marginTop + e.availableHeight) - e.currentY

				// Try placing Group 2 (2-3) on the same new page
				fmt.Printf("Debug (process4Split): Rule 4 - Attempting G2 (2-3) on same new page %d (Avail H: %.2f).\n", e.currentPage.Page, newPageAvailableHeightG1)
				layoutInfoG2New, errG2New := e.calculateTwoPicturesLayout(pictures[G2Start:G2End], newPageAvailableHeightG1)
				if errG2New == nil && layoutInfoG2New.TotalHeight <= newPageAvailableHeightG1+tolerance {
					// Rule 4 Success Path: G2 fits after G1 on new page (G1+G2 success)
					fmt.Printf("Debug (process4Split): Rule 4 - Success. Placing G2 (H: %.2f). G1+G2 on same new page complete.\n", layoutInfoG2New.TotalHeight)
					e.placePicturesInTemplate(pictures[G2Start:G2End], layoutInfoG2New)
					e.currentY += layoutInfoG2New.TotalHeight
					return 0
				} else {
					// Rule 4 Failed: G2 failed -> Rule 5: New page for G2
					fmt.Printf("Debug (process4Split): Rule 4 failed (G2 on same new page). Err: %v / Height: %.2f. Proceeding to Rule 5 (New page for G2).\n", errG2New, layoutInfoG2New.TotalHeight)
					return e.placeG2TwoOnNewPage(pictures[G2Start:G2End])
				}
			} else {
				// Rule 4 Failed Critically: G1 doesn't fit even on the new page.
				fmt.Printf("Warning (process4Split): Rule 4 - G1 (0-1) failed to place on new page %d (Avail H: %.2f). Err: %v / Height: %.2f. Proceeding to Rule 5 (New page for G2), G1 pics lost.\n", e.currentPage.Page, newPageAvailableHeight1, errG1New, layoutInfoG1New.TotalHeight)
				// Proceed to Rule 5, G1 is lost.
				return e.placeG2TwoOnNewPage(pictures[G2Start:G2End])
			}
		}
	}
}

// placeG2TwoOnNewPage handles Rule 5 logic: Create new page and place G2 (2-3).
func (e *ContinuousLayoutEngine) placeG2TwoOnNewPage(picturesG2 []Picture) float64 {
	const tolerance = 1e-6
	if len(picturesG2) != 2 {
		fmt.Printf("Error (process4Split Rule 5): Expected 2 pictures for G2, got %d\n", len(picturesG2))
		return 0
	}

	// --- Rule 5: New page for G2 (2-3) ---
	fmt.Printf("Debug (process4Split): Rule 5 - New page (Page %d) for G2 (2-3).\n", e.currentPage.Page+1)
	e.newPage()
	e.currentY = e.marginTop
	newPageAvailableHeight2 := e.availableHeight

	fmt.Printf("Debug (process4Split): Rule 5 - Attempting G2 (2-3) on new page %d (Avail H: %.2f).\n", e.currentPage.Page, newPageAvailableHeight2)
	layoutInfoG2Final, errG2Final := e.calculateTwoPicturesLayout(picturesG2, newPageAvailableHeight2)
	if errG2Final == nil && layoutInfoG2Final.TotalHeight <= newPageAvailableHeight2+tolerance {
		fmt.Printf("Debug (process4Split): Rule 5 - Success. Placing G2 (H: %.2f) on new page %d.\n", layoutInfoG2Final.TotalHeight, e.currentPage.Page)
		e.placePicturesInTemplate(picturesG2, layoutInfoG2Final)
		e.currentY += layoutInfoG2Final.TotalHeight
	} else {
		// Rule 5 Failed: G2 failed even on its own dedicated page.
		fmt.Printf("Error (process4Split): Rule 5 - Critical failure. G2 (2-3) failed to place even on new page %d (Avail H: %.2f). Err: %v / Height: %.2f. Aborting placement of G2. Pics 2-3 lost.\n", e.currentPage.Page, newPageAvailableHeight2, errG2Final, layoutInfoG2Final.TotalHeight)
		// if errors.Is(errG2Final, ErrMinHeightConstraint) { ... }
	}
	return 0 // Always return 0 as split occurred.
}
