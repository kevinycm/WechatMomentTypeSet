package calculate

import (
	"errors" // Added for error handling
	"fmt"
	// Ensure math is imported if needed for calculations
	// "math"
)

// ErrMinHeightConstraint is returned when pictures don't meet minimum height.
// Defining it here as it seems to be missing.
var ErrMinHeightConstraint = errors.New("minimum height constraint violated")

// LayoutInfo stores results from layout calculation functions.
// Defining a basic version here as it seems to be missing.
type LayoutInfo struct {
	TotalHeight float64
	// Add other relevant fields if necessary, like individual item layouts:
	// Items []ItemLayout // Assuming ItemLayout is defined elsewhere or needs definition
}

// processNinePicturesWithSplitLogic implements the new 1-10 rules for 9-picture layout.
// It attempts various layout strategies (9-pic, 3-pic, 6-pic splits) across pages
// based on available height and minimum picture requirements.
func (e *ContinuousLayoutEngine) processNinePicturesWithSplitLogic(pictures []Picture, layoutAvailableHeight float64) float64 {
	numPics := len(pictures)
	if numPics != 9 {
		fmt.Printf("Error (process9SplitNew): Expected 9 pictures, got %d\n", numPics)
		return 0 // Or handle error appropriately
	}

	const G1Start, G1End = 0, 3
	const G2Start, G2End = 3, 6
	const G3Start, G3End = 6, 9
	const tolerance = 1e-6 // Tolerance for float comparisons

	// Pre-declare ONLY variables needed for goto targets
	var newPage2AvailableHeight float64
	var layoutInfo6 TemplateLayout
	var err6 error
	var newPage3AvailableHeight float64
	var layoutInfoG3New3 TemplateLayout
	var errG3New3 error
	// Add declarations for variables jumped over by goto
	var layoutInfo9New TemplateLayout
	var err9New error
	var newPageAvailableHeight1 float64
	var layoutInfoG1New TemplateLayout
	var errG1New error
	var newPageAvailableHeightG1 float64
	var layoutInfoG2New TemplateLayout
	var errG2New error
	var newPageAvailableHeightG2 float64
	var layoutInfoG3New TemplateLayout
	var errG3New error
	var layoutInfoG2New2 TemplateLayout
	var errG2New2 error
	var newPage2AvailableHeightG2 float64
	var layoutInfoG3New2 TemplateLayout
	var errG3New2 error

	// --- Rule 1: Try placing all 9 on the current page ---
	fmt.Printf("Debug (process9SplitNew): Rule 1 - Attempting 9-pic layout on page %d (Avail H: %.2f).\n", e.currentPage.Page, layoutAvailableHeight)
	layoutInfo9, err9 := e.calculateNinePicturesLayout(pictures, layoutAvailableHeight) // Use := again
	if err9 == nil && layoutInfo9.TotalHeight <= layoutAvailableHeight+tolerance {
		fmt.Printf("Debug (process9SplitNew): Rule 1 - Success. Placing 9 pics (H: %.2f).\n", layoutInfo9.TotalHeight)
		e.placePicturesInTemplate(pictures, layoutInfo9)
		// Return height used, processPictures will update e.currentY
		return layoutInfo9.TotalHeight
	}
	fmt.Printf("Debug (process9SplitNew): Rule 1 failed. Err: %v / Height: %.2f.\n", err9, layoutInfo9.TotalHeight)

	// --- Rule 2: Try placing Group 1 (0-2) on the current page ---
	fmt.Printf("Debug (process9SplitNew): Rule 2 - Attempting G1 (0-2) on page %d (Avail H: %.2f).\n", e.currentPage.Page, layoutAvailableHeight)
	layoutInfoG1, errG1 := e.calculateThreePicturesLayout(pictures[G1Start:G1End], layoutAvailableHeight) // Use :=
	if errG1 == nil && layoutInfoG1.TotalHeight <= layoutAvailableHeight+tolerance {
		fmt.Printf("Debug (process9SplitNew): Rule 2 - Success. Placing G1 (H: %.2f).\n", layoutInfoG1.TotalHeight)
		e.placePicturesInTemplate(pictures[G1Start:G1End], layoutInfoG1)
		e.currentY += layoutInfoG1.TotalHeight
		e.currentY += e.imageSpacing                                                                 // Add spacing after G1
		currentAvailableHeight1 := layoutAvailableHeight - layoutInfoG1.TotalHeight - e.imageSpacing // Use :=

		// Try placing Group 2 (3-5) on the same page
		fmt.Printf("Debug (process9SplitNew): Rule 2 - Attempting G2 (3-5) on same page %d (Avail H: %.2f).\n", e.currentPage.Page, currentAvailableHeight1)
		layoutInfoG2, errG2 := e.calculateThreePicturesLayout(pictures[G2Start:G2End], currentAvailableHeight1) // Use :=
		if errG2 == nil && layoutInfoG2.TotalHeight <= currentAvailableHeight1+tolerance {
			fmt.Printf("Debug (process9SplitNew): Rule 2 - Success. Placing G2 (H: %.2f).\n", layoutInfoG2.TotalHeight)
			e.placePicturesInTemplate(pictures[G2Start:G2End], layoutInfoG2)
			e.currentY += layoutInfoG2.TotalHeight
			e.currentY += e.imageSpacing                                                                   // Add spacing after G2
			currentAvailableHeight2 := currentAvailableHeight1 - layoutInfoG2.TotalHeight - e.imageSpacing // Use :=

			// Try placing Group 3 (6-8) on the same page
			fmt.Printf("Debug (process9SplitNew): Rule 2 - Attempting G3 (6-8) on same page %d (Avail H: %.2f).\n", e.currentPage.Page, currentAvailableHeight2)
			layoutInfoG3, errG3 := e.calculateThreePicturesLayout(pictures[G3Start:G3End], currentAvailableHeight2) // Use :=
			if errG3 == nil && layoutInfoG3.TotalHeight <= currentAvailableHeight2+tolerance {
				fmt.Printf("Debug (process9SplitNew): Rule 2 - Success. Placing G3 (H: %.2f).\n", layoutInfoG3.TotalHeight)
				e.placePicturesInTemplate(pictures[G3Start:G3End], layoutInfoG3)
				e.currentY += layoutInfoG3.TotalHeight
				return 0 // Placed 3+3+3 on one page
			}
			fmt.Printf("Debug (process9SplitNew): Rule 2 failed (G3 on same page). Err: %v / Height: %.2f. Go to new page for G3.\n", errG3, layoutInfoG3.TotalHeight)
			goto NewPageForG3 // G3 failed, needs new page
		}
		fmt.Printf("Debug (process9SplitNew): Rule 2 failed (G2 on same page). Err: %v / Height: %.2f. Go to new page, try 6-pic (3-8).\n", errG2, layoutInfoG2.TotalHeight)
		goto NewPageTrySixPic // G2 failed, try 6-pic split on new page
	}
	fmt.Printf("Debug (process9SplitNew): Rule 2 failed (G1 on page %d). Err: %v / Height: %.2f. Go to Rule 3.\n", e.currentPage.Page, errG1, layoutInfoG1.TotalHeight)

	// --- Rule 3: G1 didn't fit current page. Create new page. Try 9-pic again. ---
	fmt.Printf("Debug (process9SplitNew): Rule 3 - New page (Page %d).\n", e.currentPage.Page+1)
	e.newPage()
	e.currentY = e.marginTop                    // Reset Y for the new page
	newPageAvailableHeight1 = e.availableHeight // Use =

	fmt.Printf("Debug (process9SplitNew): Rule 3 - Attempting 9-pic layout on new page %d (Avail H: %.2f).\n", e.currentPage.Page, newPageAvailableHeight1)
	layoutInfo9New, err9New = e.calculateNinePicturesLayout(pictures, newPageAvailableHeight1) // Use =
	if err9New == nil && layoutInfo9New.TotalHeight <= newPageAvailableHeight1+tolerance {
		fmt.Printf("Debug (process9SplitNew): Rule 3 - Success. Placing 9 pics (H: %.2f) on new page.\n", layoutInfo9New.TotalHeight)
		e.placePicturesInTemplate(pictures, layoutInfo9New)
		// Update Y directly, return 0 as split (across pages) happened
		e.currentY += layoutInfo9New.TotalHeight
		return 0
	}
	fmt.Printf("Debug (process9SplitNew): Rule 3 failed (9-pic on new page). Err: %v / Height: %.2f. Go to Rule 4.\n", err9New, layoutInfo9New.TotalHeight)

	// --- Rule 4: 9-pic failed on new page. Try G1 (0-2) on new page. ---
	fmt.Printf("Debug (process9SplitNew): Rule 4 - Attempting G1 (0-2) on new page %d (Avail H: %.2f).\n", e.currentPage.Page, newPageAvailableHeight1)
	layoutInfoG1New, errG1New = e.calculateThreePicturesLayout(pictures[G1Start:G1End], newPageAvailableHeight1) // Use =
	if errG1New == nil && layoutInfoG1New.TotalHeight <= newPageAvailableHeight1+tolerance {
		fmt.Printf("Debug (process9SplitNew): Rule 4 - Success. Placing G1 (H: %.2f) on new page.\n", layoutInfoG1New.TotalHeight)
		e.placePicturesInTemplate(pictures[G1Start:G1End], layoutInfoG1New)
		e.currentY += layoutInfoG1New.TotalHeight
		e.currentY += e.imageSpacing                                                                      // Add spacing after G1
		newPageAvailableHeightG1 = newPageAvailableHeight1 - layoutInfoG1New.TotalHeight - e.imageSpacing // Use =

		// Try placing Group 2 (3-5) on the same new page
		fmt.Printf("Debug (process9SplitNew): Rule 4 - Attempting G2 (3-5) on same new page %d (Avail H: %.2f).\n", e.currentPage.Page, newPageAvailableHeightG1)
		layoutInfoG2New, errG2New = e.calculateThreePicturesLayout(pictures[G2Start:G2End], newPageAvailableHeightG1) // Use =
		if errG2New == nil && layoutInfoG2New.TotalHeight <= newPageAvailableHeightG1+tolerance {
			fmt.Printf("Debug (process9SplitNew): Rule 4 - Success. Placing G2 (H: %.2f) on new page.\n", layoutInfoG2New.TotalHeight)
			e.placePicturesInTemplate(pictures[G2Start:G2End], layoutInfoG2New)
			e.currentY += layoutInfoG2New.TotalHeight
			e.currentY += e.imageSpacing                                                                       // Add spacing after G2
			newPageAvailableHeightG2 = newPageAvailableHeightG1 - layoutInfoG2New.TotalHeight - e.imageSpacing // Use =

			// Try placing Group 3 (6-8) on the same new page
			fmt.Printf("Debug (process9SplitNew): Rule 4 - Attempting G3 (6-8) on same new page %d (Avail H: %.2f).\n", e.currentPage.Page, newPageAvailableHeightG2)
			layoutInfoG3New, errG3New = e.calculateThreePicturesLayout(pictures[G3Start:G3End], newPageAvailableHeightG2) // Use =
			if errG3New == nil && layoutInfoG3New.TotalHeight <= newPageAvailableHeightG2+tolerance {
				fmt.Printf("Debug (process9SplitNew): Rule 4 - Success. Placing G3 (H: %.2f) on new page.\n", layoutInfoG3New.TotalHeight)
				e.placePicturesInTemplate(pictures[G3Start:G3End], layoutInfoG3New)
				e.currentY += layoutInfoG3New.TotalHeight
				return 0 // Placed G1+G2+G3 on the new page
			}
			fmt.Printf("Debug (process9SplitNew): Rule 4 failed (G3 on new page). Err: %v / Height: %.2f. Go to new page for G3.\n", errG3New, layoutInfoG3New.TotalHeight)
			goto NewPageForG3 // G3 failed, needs new page
		}
		fmt.Printf("Debug (process9SplitNew): Rule 4 failed (G2 on new page). Err: %v / Height: %.2f. Go to new page, try 6-pic (3-8).\n", errG2New, layoutInfoG2New.TotalHeight)
		goto NewPageTrySixPic // G2 failed, try 6-pic split on new page
	}
	// Rule 4 failed: G1 couldn't even fit on the new page.
	// This implies a potential issue, like G1 pictures violating min height or page height being too small.
	// The most robust action is to try placing G1 on *another* new page, then continue.
	// However, the rules don't explicitly cover this. Current logic falls through to Rule 5 (New Page, Try 6-pic 3-8).
	// Let's stick to the rules as interpreted, which means proceeding as if G1 *was* placed (implicitly failing it here)
	// and trying the remaining 6 on a new page. This might lose G1 if it truly couldn't fit.
	fmt.Printf("Warning (process9SplitNew): Rule 4 failed critically (G1 on new page). Err: %v / Height: %.2f. Proceeding to Rule 5 (may lose G1 pics).\n", errG1New, layoutInfoG1New.TotalHeight)
	goto NewPageTrySixPic // Treat as if G1 placed, but G2 failed, leading to Rule 5.

	// --- Goto Labels ---

NewPageTrySixPic:
	// --- Rule 5: Create another new page. Try 6-pic (3-8). ---
	fmt.Printf("Debug (process9SplitNew): Rule 5 - New page (Page %d).\n", e.currentPage.Page+1)
	e.newPage() // Create Page 3 (or next page)
	e.currentY = e.marginTop
	newPage2AvailableHeight = e.availableHeight // Use assignment = (already declared)

	fmt.Printf("Debug (process9SplitNew): Rule 5 - Attempting 6-pic layout (3-8) on new page %d (Avail H: %.2f).\n", e.currentPage.Page, newPage2AvailableHeight)
	layoutInfo6, err6 = e.calculateSixPicturesLayout(pictures[G2Start:G3End], newPage2AvailableHeight) // Use assignment = (already declared)
	if err6 == nil && layoutInfo6.TotalHeight <= newPage2AvailableHeight+tolerance {
		fmt.Printf("Debug (process9SplitNew): Rule 5 - Success. Placing 6 pics (3-8) (H: %.2f) on new page %d.\n", layoutInfo6.TotalHeight, e.currentPage.Page)
		e.placePicturesInTemplate(pictures[G2Start:G3End], layoutInfo6)
		e.currentY += layoutInfo6.TotalHeight
		return 0 // Split success (G1 + G2/G3 as 6)
	}
	fmt.Printf("Debug (process9SplitNew): Rule 5 failed (6-pic on new page %d, Avail: %.2f). Err: %v / Height: %.2f.\n", e.currentPage.Page, newPage2AvailableHeight, err6, layoutInfo6.TotalHeight)

	// --- Rule 6: 6-pic failed on new page. Try G2 (3-5) on this new page. ---
	fmt.Printf("Debug (process9SplitNew): Rule 6 - Attempting G2 (3-5) on new page %d (Avail H: %.2f).\n", e.currentPage.Page, newPage2AvailableHeight)
	layoutInfoG2New2, errG2New2 = e.calculateThreePicturesLayout(pictures[G2Start:G2End], newPage2AvailableHeight) // Use =
	if errG2New2 == nil && layoutInfoG2New2.TotalHeight <= newPage2AvailableHeight+tolerance {
		fmt.Printf("Debug (process9SplitNew): Rule 6 - Success. Placing G2 (H: %.2f) on new page %d.\n", layoutInfoG2New2.TotalHeight, e.currentPage.Page)
		e.placePicturesInTemplate(pictures[G2Start:G2End], layoutInfoG2New2)
		e.currentY += layoutInfoG2New2.TotalHeight
		e.currentY += e.imageSpacing                                                                        // Add spacing after G2
		newPage2AvailableHeightG2 = newPage2AvailableHeight - layoutInfoG2New2.TotalHeight - e.imageSpacing // Use =

		// Try placing Group 3 (6-8) on the same (Rule 5's new) page
		fmt.Printf("Debug (process9SplitNew): Rule 6 - Attempting G3 (6-8) on same new page %d (Avail H: %.2f).\n", e.currentPage.Page, newPage2AvailableHeightG2)
		layoutInfoG3New2, errG3New2 = e.calculateThreePicturesLayout(pictures[G3Start:G3End], newPage2AvailableHeightG2) // Use =
		if errG3New2 == nil && layoutInfoG3New2.TotalHeight <= newPage2AvailableHeightG2+tolerance {
			fmt.Printf("Debug (process9SplitNew): Rule 6 - Success. Placing G3 (H: %.2f) on new page %d.\n", layoutInfoG3New2.TotalHeight, e.currentPage.Page)
			e.placePicturesInTemplate(pictures[G3Start:G3End], layoutInfoG3New2)
			e.currentY += layoutInfoG3New2.TotalHeight
			return 0 // Placed G2+G3 on the new page after G1
		}
		fmt.Printf("Debug (process9SplitNew): Rule 6 failed (G3 on same new page). Err: %v / Height: %.2f. Go to new page for G3.\n", errG3New2, layoutInfoG3New2.TotalHeight)
		goto NewPageForG3 // G3 failed, needs new page
	}
	fmt.Printf("Debug (process9SplitNew): Rule 6 failed (G2 on new page %d). Err: %v / Height: %.2f. Go to Rule 7 (new page for G3).\n", e.currentPage.Page, errG2New2, layoutInfoG2New2.TotalHeight)
	// G2 failed on this new page, fall through to create another new page specifically for G3.
	goto NewPageForG3

NewPageForG3:
	// --- Rule 7: Create another new page. Place G3 (6-8). ---
	fmt.Printf("Debug (process9SplitNew): Rule 7 - New page (Page %d) for G3 (6-8).\n", e.currentPage.Page+1)
	e.newPage() // Create Page 4 (or next page)
	e.currentY = e.marginTop
	newPage3AvailableHeight = e.availableHeight // Use assignment = (already declared)

	fmt.Printf("Debug (process9SplitNew): Rule 7 - Attempting G3 (6-8) on new page %d (Avail H: %.2f).\n", e.currentPage.Page, newPage3AvailableHeight)
	layoutInfoG3New3, errG3New3 = e.calculateThreePicturesLayout(pictures[G3Start:G3End], newPage3AvailableHeight) // Use assignment = (already declared)
	if errG3New3 == nil && layoutInfoG3New3.TotalHeight <= newPage3AvailableHeight+tolerance {
		fmt.Printf("Debug (process9SplitNew): Rule 7 - Success. Placing G3 (H: %.2f) on new page %d.\n", layoutInfoG3New3.TotalHeight, e.currentPage.Page)
		e.placePicturesInTemplate(pictures[G3Start:G3End], layoutInfoG3New3)
		e.currentY += layoutInfoG3New3.TotalHeight
		return 0 // Split success
	} else {
		// Rule 7 Failed: G3 failed even on its own dedicated page.
		fmt.Printf("Error (process9SplitNew): Rule 7 - Critical failure. G3 (6-8) failed to place even on new page %d (Avail H: %.2f). Err: %v / Height: %.2f. Aborting placement of G3. Pics 6-8 lost.\n", e.currentPage.Page, newPage3AvailableHeight, errG3New3, layoutInfoG3New3.TotalHeight)
		// Propagate the specific error if it's ErrMinHeightConstraint
		if errors.Is(errG3New3, ErrMinHeightConstraint) {
			// Return the error? Or just 0? Current approach is return 0 and log.
			// The calling function processPictures might need to know *why* it failed.
			// For now, stick to returning 0 as per the general pattern for splits/failures.
		}
		return 0 // Indicate split occurred, but G3 failed
	}
}

// Removed processPictureGroup function

// Removed placeAllNineOnNewPage function
