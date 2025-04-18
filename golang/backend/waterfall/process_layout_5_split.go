package waterfall

import (
	"fmt"
)

// processFivePicturesWithSplitLogic handles layout for 5 pictures according to the specified rules.
// Returns height used ONLY if all 5 are placed together (Rule 1 or Rule 3).
// Returns 0 in all split scenarios (Rule 2, 4, 5), caller relies on final e.currentY.
func (e *ContinuousLayoutEngine) processFivePicturesWithSplitLogic(pictures []Picture, layoutAvailableHeight float64) float64 {
	numPics := len(pictures)
	if numPics != 5 {
		fmt.Printf("Error (process5Split): Expected 5 pictures, got %d\\n", numPics)
		return 0
	}

	const G1Start, G1End = 0, 2 // Group 1: Pics 0-1
	const G2Start, G2End = 2, 5 // Group 2: Pics 2-4
	const tolerance = 1e-6

	// --- Rule 1: Try placing all 5 on the current page ---
	fmt.Printf("Debug (process5Split): Rule 1 - Attempting 5-pic layout on page %d (Avail H: %.2f).\\n", e.currentPage.Page, layoutAvailableHeight)
	layoutInfo5, err5 := e.calculateFivePicturesLayout(pictures, layoutAvailableHeight)
	if err5 == nil && layoutInfo5.TotalHeight <= layoutAvailableHeight+tolerance {
		fmt.Printf("Debug (process5Split): Rule 1 - Success. Placing 5 pics (H: %.2f).\\n", layoutInfo5.TotalHeight)
		e.placePicturesInTemplate(pictures, layoutInfo5)
		return layoutInfo5.TotalHeight // Return height used
	}
	fmt.Printf("Debug (process5Split): Rule 1 failed. Err: %v / Height: %.2f. Proceeding to Rule 2.\\n", err5, layoutInfo5.TotalHeight)

	// --- Rule 2: Try placing Group 1 (0-1) on the current page ---
	fmt.Printf("Debug (process5Split): Rule 2 - Attempting G1 (0-1) on page %d (Avail H: %.2f).\\n", e.currentPage.Page, layoutAvailableHeight)
	layoutInfoG1, errG1 := e.calculateTwoPicturesLayout(pictures[G1Start:G1End], layoutAvailableHeight)
	if errG1 == nil && layoutInfoG1.TotalHeight <= layoutAvailableHeight+tolerance {
		// Rule 2 Success Path: G1 fits on current page
		fmt.Printf("Debug (process5Split): Rule 2 - Success. Placing G1 (H: %.2f).\\n", layoutInfoG1.TotalHeight)
		e.placePicturesInTemplate(pictures[G1Start:G1End], layoutInfoG1)
		e.currentY += layoutInfoG1.TotalHeight
		e.currentY += e.imageSpacing // Add spacing after G1
		currentAvailableHeight := (e.marginTop + e.availableHeight) - e.currentY

		// Try placing Group 2 (2-4) on the same current page
		fmt.Printf("Debug (process5Split): Rule 2 - Attempting G2 (2-4) on same page %d (Avail H: %.2f).\\n", e.currentPage.Page, currentAvailableHeight)
		layoutInfoG2, errG2 := e.calculateThreePicturesLayout(pictures[G2Start:G2End], currentAvailableHeight)
		if errG2 == nil && layoutInfoG2.TotalHeight <= currentAvailableHeight+tolerance {
			// G2 fits on the same page
			fmt.Printf("Debug (process5Split): Rule 2 - Success. Placing G2 (H: %.2f) on same page.\\n", layoutInfoG2.TotalHeight)
			e.placePicturesInTemplate(pictures[G2Start:G2End], layoutInfoG2)
			e.currentY += layoutInfoG2.TotalHeight
			return 0 // Split 2+3 on same page complete
		} else {
			// G2 doesn't fit on the same page -> Rule 5: New page for G2
			fmt.Printf("Debug (process5Split): Rule 2 failed (G2 on same page). Err: %v / Height: %.2f. Proceeding to Rule 5 (New page for G2).\\n", errG2, layoutInfoG2.TotalHeight)
			return e.placeG2OnNewPage(pictures[G2Start:G2End]) // Call helper for Rule 5
		}
	}

	// --- Rule 3: G1 failed on current page (Rule 2 failed). New page, try 5-pic again. ---
	fmt.Printf("Debug (process5Split): Rule 2 failed (G1 on page %d). Err: %v / Height: %.2f. Proceeding to Rule 3 (New page).\\n", e.currentPage.Page, errG1, layoutInfoG1.TotalHeight)
	e.newPage()
	e.currentY = e.marginTop
	newPageAvailableHeight1 := e.availableHeight

	fmt.Printf("Debug (process5Split): Rule 3 - Attempting 5-pic layout on new page %d (Avail H: %.2f).\\n", e.currentPage.Page, newPageAvailableHeight1)
	layoutInfo5New, err5New := e.calculateFivePicturesLayout(pictures, newPageAvailableHeight1)
	if err5New == nil && layoutInfo5New.TotalHeight <= newPageAvailableHeight1+tolerance {
		// Rule 3 Success: 5 pics fit on the new page
		fmt.Printf("Debug (process5Split): Rule 3 - Success. Placing 5 pics (H: %.2f) on new page.\\n", layoutInfo5New.TotalHeight)
		e.placePicturesInTemplate(pictures, layoutInfo5New)
		e.currentY += layoutInfo5New.TotalHeight
		return 0 // Return 0 as split across pages occurred
	}

	// --- Rule 4: 5-pic failed on new page (Rule 3 failed). Try G1 (0-1) on new page. ---
	fmt.Printf("Debug (process5Split): Rule 3 failed (5-pic on new page). Err: %v / Height: %.2f. Proceeding to Rule 4.\\n", err5New, layoutInfo5New.TotalHeight)
	// Note: We are still on the new page created in Rule 3.
	// Reset Y for placing G1 at the top of this new page.
	e.currentY = e.marginTop
	// Available height is still newPageAvailableHeight1 for this attempt.

	fmt.Printf("Debug (process5Split): Rule 4 - Attempting G1 (0-1) on new page %d (Avail H: %.2f).\\n", e.currentPage.Page, newPageAvailableHeight1)
	layoutInfoG1New, errG1New := e.calculateTwoPicturesLayout(pictures[G1Start:G1End], newPageAvailableHeight1)
	if errG1New == nil && layoutInfoG1New.TotalHeight <= newPageAvailableHeight1+tolerance {
		// Rule 4 Success Path: G1 fits on new page
		fmt.Printf("Debug (process5Split): Rule 4 - Success. Placing G1 (H: %.2f) on new page.\\n", layoutInfoG1New.TotalHeight)
		e.placePicturesInTemplate(pictures[G1Start:G1End], layoutInfoG1New)
		e.currentY += layoutInfoG1New.TotalHeight
		e.currentY += e.imageSpacing // Add spacing after G1
		newPageAvailableHeightG1 := (e.marginTop + e.availableHeight) - e.currentY

		// Try placing Group 2 (2-4) on the same new page
		fmt.Printf("Debug (process5Split): Rule 4 - Attempting G2 (2-4) on same new page %d (Avail H: %.2f).\\n", e.currentPage.Page, newPageAvailableHeightG1)
		layoutInfoG2New, errG2New := e.calculateThreePicturesLayout(pictures[G2Start:G2End], newPageAvailableHeightG1)
		if errG2New == nil && layoutInfoG2New.TotalHeight <= newPageAvailableHeightG1+tolerance {
			// G2 fits on the same new page
			fmt.Printf("Debug (process5Split): Rule 4 - Success. Placing G2 (H: %.2f) on same new page.\\n", layoutInfoG2New.TotalHeight)
			e.placePicturesInTemplate(pictures[G2Start:G2End], layoutInfoG2New)
			e.currentY += layoutInfoG2New.TotalHeight
			return 0 // Split 2+3 on new page complete
		} else {
			// G2 doesn't fit on the new page after G1 -> Rule 5: New page for G2
			fmt.Printf("Debug (process5Split): Rule 4 failed (G2 on same new page). Err: %v / Height: %.2f. Proceeding to Rule 5 (New page for G2).\\n", errG2New, layoutInfoG2New.TotalHeight)
			return e.placeG2OnNewPage(pictures[G2Start:G2End]) // Call helper for Rule 5
		}
	} else {
		// Rule 4 Failed Critically: G1 doesn't fit even on the new page.
		fmt.Printf("Warning (process5Split): Rule 4 - G1 (0-1) failed to place on new page %d (Avail H: %.2f). Err: %v / Height: %.2f. Proceeding to Rule 5 (New page for G2), G1 pics lost.\\n", e.currentPage.Page, newPageAvailableHeight1, errG1New, layoutInfoG1New.TotalHeight)
		// Proceed to place G2 on yet another new page, G1 is lost.
		return e.placeG2OnNewPage(pictures[G2Start:G2End]) // Call helper for Rule 5
	}
}

// placeG2OnNewPage is a helper function implementing Rule 5 logic.
// It creates a new page and attempts to place Group 2 (pictures 2-4).
// Returns 0, indicating a split occurred.
func (e *ContinuousLayoutEngine) placeG2OnNewPage(picturesG2 []Picture) float64 {
	if len(picturesG2) != 3 {
		fmt.Printf("Error (placeG2OnNewPage): Expected 3 pictures for G2, got %d\\n", len(picturesG2))
		return 0 // Should not happen if called correctly
	}
	const tolerance = 1e-6 // Define tolerance locally
	fmt.Printf("Debug (process5Split): Rule 5 - New page (Page %d) for G2 (2-4).\\n", e.currentPage.Page+1)
	e.newPage() // Create Page 3 (or next)
	e.currentY = e.marginTop
	newPageAvailableHeight2 := e.availableHeight

	fmt.Printf("Debug (process5Split): Rule 5 - Attempting G2 (2-4) on new page %d (Avail H: %.2f).\\n", e.currentPage.Page, newPageAvailableHeight2)
	layoutInfoG2Final, errG2Final := e.calculateThreePicturesLayout(picturesG2, newPageAvailableHeight2)
	if errG2Final == nil && layoutInfoG2Final.TotalHeight <= newPageAvailableHeight2+tolerance {
		fmt.Printf("Debug (process5Split): Rule 5 - Success. Placing G2 (H: %.2f) on new page %d.\\n", layoutInfoG2Final.TotalHeight, e.currentPage.Page)
		e.placePicturesInTemplate(picturesG2, layoutInfoG2Final)
		e.currentY += layoutInfoG2Final.TotalHeight
		// return 0 // Split success - Already returns 0 implicitly by function end
	} else {
		// Rule 5 Failed: G2 failed even on its own dedicated page.
		fmt.Printf("Error (process5Split): Rule 5 - Critical failure. G2 (2-4) failed to place even on new page %d (Avail H: %.2f). Err: %v / Height: %.2f. Aborting placement of G2. Pics 2-4 lost.\\n", e.currentPage.Page, newPageAvailableHeight2, errG2Final, layoutInfoG2Final.TotalHeight)
		// if errors.Is(errG2Final, ErrMinHeightConstraint) { ... }
		// return 0 // Indicate split occurred, but G2 failed - Already returns 0 implicitly
	}
	return 0 // Always return 0 from this helper as a split occurred.
}
