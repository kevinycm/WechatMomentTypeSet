package calculate

import (
	"fmt"
	"math"
	"strings"
	"time"
	"unicode/utf8"
)

// NewContinuousLayoutEngine creates a new continuous layout engine
// UPDATED TO INITIALIZE NEW FIELDS
func NewContinuousLayoutEngine(entries []Entry) *ContinuousLayoutEngine {
	engine := &ContinuousLayoutEngine{
		entries:                          entries,
		marginLeft:                       142,
		marginRight:                      142,
		marginTop:                        189,
		marginBottom:                     189,
		timeHeight:                       100,
		fontSize:                         66.67, // 对应72DPI的16px
		lineHeight:                       100,   // 对应72DPI的24px
		entrySpacing:                     150,   // 条目之间的间距
		elementSpacing:                   30,    // 元素整体之间的间距
		imageSpacing:                     15,    // 图片之间的间距
		minWideHeight:                    600,   // Min height for Wide pics (AR >= 3)
		minTallHeight:                    800,   // Min height for Tall pics (AR <= 1/3)
		minLandscapeHeight:               600,   // Base Min height for Landscape (for < 5 pics)
		minPortraitHeight:                600,   // Base Min height for Portrait (for < 5 pics)
		minLandscapeHeightLargeGroup:     600,   // Min height Landscape (5-7 pics)
		minPortraitHeightLargeGroup:      800,   // Min height Portrait (5-7 pics)
		minLandscapeHeightVeryLargeGroup: 450,   // Added: Min height Landscape (>= 8 pics) - Lower value
		minPortraitHeightVeryLargeGroup:  900,   // Added: Min height Portrait (>= 8 pics) - Lower value
		singleImageHeight:                2808,  // 设置单张竖图的默认高度
		singleImageWidth:                 1695,  // 设置单张横图的默认宽度
		minImageHeight:                   800,   // 设置单张竖图的最小展示高度
		minImageWidth:                    1200,  // 设置单张横图的最小展示宽度
		bottomMargin:                     100,   // 底部边距
	}
	engine.availableWidth = 2480 - engine.marginLeft - engine.marginRight
	engine.availableHeight = 3508 - engine.marginTop - engine.marginBottom

	// 从第一个条目中获取年月信息
	if len(entries) > 0 {
		t, err := time.Parse("2006-01-02 15:04:05", entries[0].Time)
		if err == nil {
			engine.currentYearMonth = fmt.Sprintf("%d年%d月", t.Year(), t.Month())
		}
	}

	return engine
}

// ProcessEntries processes all entries and returns the layout result
func (e *ContinuousLayoutEngine) ProcessEntries() ([]ContinuousLayoutPage, error) {
	e.newPage() // Start with a fresh page

	for _, entry := range e.entries {
		// Let processEntry handle content placement and pagination internally
		e.processEntry(entry)
	}

	return e.pages, nil
}

func (e *ContinuousLayoutEngine) processEntry(entry Entry) {
	// Add entry spacing if this isn't the first element on the page
	// We need a more robust check than just e.currentY > e.marginTop
	// Check if the current page actually has content already placed
	// Simplified: If currentY was moved from the top margin, add spacing before the next entry.
	if e.currentY > e.marginTop {
		// Check if adding spacing would push us off the page
		if e.currentY+e.entrySpacing > e.availableHeight+e.marginTop {
			// If just the spacing doesn't fit, go to a new page
			e.newPage()
		} else {
			// Add spacing
			e.currentY += e.entrySpacing
		}
	}

	// 1. Process Time
	e.addTime(entry.Time)

	// 2. Process Text (handles its own internal pagination)
	if strings.TrimSpace(entry.Text) != "" {
		e.processText(entry.Text)
	}

	// 3. Process Pictures (handles its own row-by-row pagination)
	if len(entry.Pictures) > 0 {
		e.processPictures(entry.Pictures)
	}
}

// Modify addTime to handle potential page break *before* adding the time entry
func (e *ContinuousLayoutEngine) addTime(timeStr string) {
	// Check if space is available for the time block itself
	minTimeHeight := e.timeHeight
	// Calculate remaining space accurately at this point
	pageAvailableHeight := e.availableHeight - (e.currentY - e.marginTop)

	if pageAvailableHeight < minTimeHeight {
		// Not enough space even for the time block, create a new page
		e.newPage()
		// Reset currentY is handled by newPage
	}

	// Now place the time block (currentY is guaranteed to be valid)
	x0 := e.marginLeft
	y0 := e.currentY
	x1 := x0 + e.availableWidth
	y1 := y0 + e.timeHeight

	// Ensure the current page has an entry list initialized
	if e.currentPage.Entries == nil {
		e.currentPage.Entries = make([]PageEntry, 0)
	}

	// Create the entry *now*, right before placing the time
	entry := PageEntry{
		Time:      timeStr,
		TimeArea:  [][]float64{{x0, y0}, {x1, y1}},
		TextAreas: make([][][]float64, 0), // Initialize other fields
		Texts:     make([]string, 0),
		Pictures:  make([]Picture, 0),
	}

	// Parse and format time if possible
	layouts := []string{
		"2006-01-02 15:04:05",
		"2006年1月2日 15:04",
		"2006年01月02日 15:04",
	}
	var t time.Time
	var err error
	for _, layout := range layouts {
		t, err = time.Parse(layout, timeStr)
		if err == nil {
			break
		}
	}

	if err == nil {
		weekdayMap := map[time.Weekday]string{
			time.Sunday: "周日", time.Monday: "周一", time.Tuesday: "周二",
			time.Wednesday: "周三", time.Thursday: "周四", time.Friday: "周五",
			time.Saturday: "周六",
		}
		entry.DatePart = fmt.Sprintf("%d月%d日 %s", t.Month(), t.Day(), weekdayMap[t.Weekday()])
		entry.TimePart = fmt.Sprintf("%02d:%02d", t.Hour(), t.Minute())
	}

	// Add the newly created entry to the current page
	e.currentPage.Entries = append(e.currentPage.Entries, entry)

	// Update current Y position *after* placing time
	// Add elementSpacing only if text or pictures will follow
	// This spacing logic might be better handled *before* text/pictures are added.
	// Let's update Y simply by the time height for now.
	e.currentY = y1 // Position is now at the bottom of the time block
	// e.timeAreaBottom = y1 - e.marginTop // This seems less critical now
}

// Modify addTextChunk to correctly append to the *last* entry on the page
// And handle spacing *before* adding the text chunk
func (e *ContinuousLayoutEngine) addTextChunk(chunk []string) {
	if len(chunk) == 0 {
		return
	}

	// Ensure there is an entry to add to
	if len(e.currentPage.Entries) == 0 {
		// This case should ideally not happen if addTime always creates an entry
		// If it does, it means text is the very first element.
		entry := PageEntry{
			TextAreas: make([][][]float64, 0),
			Texts:     make([]string, 0),
			Pictures:  make([]Picture, 0),
		}
		e.currentPage.Entries = append(e.currentPage.Entries, entry)
	} // else, we assume addTime created the entry

	// Get the last entry on the current page
	currentEntryIndex := len(e.currentPage.Entries) - 1
	currentEntry := &e.currentPage.Entries[currentEntryIndex]

	// Add element spacing *before* the text if needed
	// Check if the entry currently only contains the TimeArea
	if len(currentEntry.TextAreas) == 0 && len(currentEntry.Pictures) == 0 && len(currentEntry.TimeArea) > 0 {
		// Check for space before adding spacing + text
		pageAvailableHeight := e.availableHeight - (e.currentY - e.marginTop)
		if pageAvailableHeight < e.elementSpacing+e.lineHeight {
			// Not enough space for spacing + 1 line of text, force new page
			e.newPage()
			// Create a new entry on the new page to hold the continued text
			// (or handle splitting the original entry across pages - more complex)
			// Simpler: Assume the text chunk continues in a new logical block.
			entry := PageEntry{
				// Copy relevant identifiers? No, treat as new block for now.
				TextAreas: make([][][]float64, 0),
				Texts:     make([]string, 0),
				Pictures:  make([]Picture, 0),
			}
			e.currentPage.Entries = append(e.currentPage.Entries, entry)
			currentEntryIndex = len(e.currentPage.Entries) - 1 // Update index
			currentEntry = &e.currentPage.Entries[currentEntryIndex]
			// No spacing needed at start of new page/entry
		} else {
			// Enough space for spacing + text line, add spacing now.
			e.currentY += e.elementSpacing
		}
	} // else, no spacing needed (continuation of text, or text is first element)

	startY := e.currentY
	textHeight := float64(len(chunk)) * e.lineHeight
	area := [][]float64{
		{e.marginLeft, startY},
		{e.marginLeft + e.availableWidth, startY + textHeight},
	}

	currentEntry.TextAreas = append(currentEntry.TextAreas, area)
	currentEntry.Texts = append(currentEntry.Texts, strings.Join(chunk, "\n"))

	// Update Y position after adding text chunk
	e.currentY += textHeight
}

// processText needs slight modification to handle page breaks correctly with new addTextChunk
func (e *ContinuousLayoutEngine) processText(text string) {
	if strings.TrimSpace(text) == "" {
		return
	}

	charsPerLine := int(e.availableWidth / e.fontSize)
	var lines []string
	for _, paragraph := range strings.Split(text, "\n") {
		charCount := utf8.RuneCountInString(paragraph)
		if charCount <= charsPerLine {
			lines = append(lines, paragraph)
		} else {
			runes := []rune(paragraph)
			for i := 0; i < len(runes); i += charsPerLine {
				end := i + charsPerLine
				if end > len(runes) {
					end = len(runes)
				}
				line := string(runes[i:end])
				lines = append(lines, line)
			}
		}
	}

	currentLine := 0
	for currentLine < len(lines) {
		// Calculate remaining space *before* attempting to add a chunk
		pageAvailableHeight := e.availableHeight - (e.currentY - e.marginTop)

		// Calculate required height for the *next* line (+ spacing if it's the first line after time)
		requiredHeight := e.lineHeight
		// Check if spacing needs to be added *before* this text block
		if currentLine == 0 && len(e.currentPage.Entries) > 0 {
			lastEntry := &e.currentPage.Entries[len(e.currentPage.Entries)-1]
			if len(lastEntry.TextAreas) == 0 && len(lastEntry.Pictures) == 0 && len(lastEntry.TimeArea) > 0 {
				requiredHeight += e.elementSpacing
			}
		}

		if pageAvailableHeight < requiredHeight {
			// Not enough space even for one more line (+ potentially spacing)
			e.newPage()
			// Create a new entry on the new page
			entry := PageEntry{
				TextAreas: make([][][]float64, 0),
				Texts:     make([]string, 0),
				Pictures:  make([]Picture, 0),
			}
			e.currentPage.Entries = append(e.currentPage.Entries, entry)
			// Recalculate available height for the loop check
			pageAvailableHeight = e.availableHeight // available on new page
			// Continue loop to place the line on the new page
		}

		// Determine how many lines *can* fit
		// Consider spacing needed before the *first* line of this chunk
		firstLineSpacing := 0.0
		if len(e.currentPage.Entries) > 0 {
			lastEntry := &e.currentPage.Entries[len(e.currentPage.Entries)-1]
			if len(lastEntry.TextAreas) == 0 && len(lastEntry.Pictures) == 0 && len(lastEntry.TimeArea) > 0 {
				firstLineSpacing = e.elementSpacing
			}
		}

		availableForText := pageAvailableHeight - firstLineSpacing
		if availableForText < 0 {
			availableForText = 0
		} // Handle edge case

		availableLines := int(math.Floor(availableForText / e.lineHeight))
		if availableLines <= 0 {
			// Should have been caught by the requiredHeight check, but as safety:
			if pageAvailableHeight >= requiredHeight { // If space exists but calculation weird
				availableLines = 1 // Place at least one line if possible
			} else {
				// This state implies a new page was needed but logic failed? Log error maybe.
				// Force a page break to avoid infinite loop
				e.newPage()
				entry := PageEntry{}
				e.currentPage.Entries = append(e.currentPage.Entries, entry)
				continue // Retry placement on new page
			}
		}

		numLinesToAdd := int(math.Min(float64(len(lines)-currentLine), float64(availableLines)))
		chunk := lines[currentLine : currentLine+numLinesToAdd]
		e.addTextChunk(chunk) // addTextChunk now handles spacing and page break check before adding
		currentLine += numLinesToAdd

		// No need to create new page here, addTextChunk handles breaks before adding
		// and the loop condition checks space at the start.
	}
}

func (e *ContinuousLayoutEngine) newPage() {
	page := &ContinuousLayoutPage{
		Page:    len(e.pages) + 1,
		Entries: make([]PageEntry, 0),
	}
	e.pages = append(e.pages, *page)
	e.currentPage = &e.pages[len(e.pages)-1]
	e.currentY = e.marginTop
	e.timeAreaBottom = 0
}

// processPictures handles layout and pagination for a block of pictures.
func (e *ContinuousLayoutEngine) processPictures(pictures []Picture) {
	numPicsTotal := len(pictures)
	if numPicsTotal == 0 {
		return
	}

	// --- Check for Ultra-Wide Pictures ---
	hasUltraWide := false
	ultraWideThreshold := 4.0
	for _, pic := range pictures {
		if pic.Height > 0 && pic.Width > 0 {
			ar := float64(pic.Width) / float64(pic.Height)
			if ar >= ultraWideThreshold {
				hasUltraWide = true
				fmt.Printf("Debug (ProcessPics): Ultra-wide picture detected (Index %d, AR %.2f). Will use dynamic layout.\n", pic.Index, ar)
				break
			}
		}
	}

	// --- Choose Layout Strategy ---
	if hasUltraWide {
		// +++ Use NEW Dynamic Row-by-Row Strategy +++
		fmt.Println("Debug (ProcessPics): Using dynamic row layout strategy.")
		// TODO: Implement the dynamic row layout logic here
		// (Loop through pictures, determine rows, calculate row layout, place row, update Y)
		fmt.Printf("Warning: Dynamic row layout for entry with ultra-wide pictures is not yet implemented. Skipping picture layout for this entry (starting index: %d).\n", pictures[0].Index) // Assuming pictures is not empty here
		// Placeholder: Fallback to old logic for now to avoid breaking everything
		// Remove this fallback once dynamic logic is implemented
		// processPicturesOldStrategy(e, pictures) // REMOVED Placeholder call

	} else {
		// +++ Use OLD Standard Templated Strategy +++
		fmt.Println("Debug (ProcessPics): No ultra-wide pictures. Using standard templated layout strategy.")
		processPicturesOldStrategy(e, pictures)
	}
}

// processPicturesOldStrategy contains the original logic using fixed templates
// Extracted to a separate function for clarity during refactoring
func processPicturesOldStrategy(e *ContinuousLayoutEngine, pictures []Picture) {
	numPics := len(pictures)
	if numPics == 0 { // Should be caught by caller, but double check
		return
	}

	// 1. Calculate Required Spacing Before Pictures
	requiredSpacing := e.requiredSpacingBeforeElement()

	// +++ Log Spacing Info +++
	fmt.Printf("Debug (OldStrategy): Before spacing. CurrentY: %.2f, RequiredSpacing: %.2f\n", e.currentY, requiredSpacing)
	// +++ End Log +++

	// 2. Centralized Pagination Check
	// Calculate remaining physical space on the current page
	pageRemainingPhysicalHeight := (e.marginTop + e.availableHeight) - e.currentY
	hasContentOnCurrentPage := e.currentY > e.marginTop // Check if anything was placed

	needsNewPage := false
	// --- Updated Pagination Logic ---
	if hasContentOnCurrentPage {
		minRequiredHeightForPics := 0.0
		switch numPics {
		case 1:
			pic := pictures[0]
			ar := 1.0
			if pic.Height > 0 && pic.Width > 0 {
				ar = float64(pic.Width) / float64(pic.Height)
			}
			picType := GetPictureType(ar)
			switch picType {
			case "wide":
				minRequiredHeightForPics = GetRequiredMinHeight(e, picType, numPics)
			case "tall":
				minRequiredHeightForPics = GetRequiredMinHeight(e, picType, numPics)
			case "landscape":
				minRequiredHeightForPics = GetRequiredMinHeight(e, picType, numPics)
			case "portrait":
				minRequiredHeightForPics = GetRequiredMinHeight(e, picType, numPics)
			default: // square or unknown
				minRequiredHeightForPics = GetRequiredMinHeight(e, "landscape", numPics) // Fallback for 1 pic
			}
		case 2:
			// For pagination check for 2 pics, use a generic minimum (e.g., landscape)
			// The actual layout calculation happens later.
			minRequiredHeightForPics = GetRequiredMinHeight(e, "landscape", numPics)
		default: // 3+ pics
			// For pagination check for 3+ pics, let's use the base landscape height as a generic minimum guess.
			// The actual detailed check happens in the layout functions.
			minRequiredHeightForPics = GetRequiredMinHeight(e, "landscape", numPics)
		}

		// Check if remaining space (after spacing) is less than the minimum needed
		if pageRemainingPhysicalHeight-requiredSpacing < minRequiredHeightForPics {
			needsNewPage = true
		}
	} // No pagination needed if current page is empty

	if needsNewPage {
		e.newPage()
		requiredSpacing = 0 // No spacing needed at top of new page
	}

	// 3. Apply Spacing (potentially on the new page)
	e.currentY += requiredSpacing
	// +++ Log Spacing Info +++
	fmt.Printf("Debug (OldStrategy): After spacing. CurrentY: %.2f\n", e.currentY)
	// +++ End Log +++

	// 4. Calculate Layout Available Height on the Target Page
	// This is the remaining space from currentY to the bottom margin boundary
	layoutAvailableHeight := (e.marginTop + e.availableHeight) - e.currentY
	if layoutAvailableHeight < 0 {
		layoutAvailableHeight = 0 // Avoid negative height if spacing pushed us over edge case
	}

	// +++ 添加日志：记录调用前的页面和Y坐标 +++
	pageBeforeLayout := e.currentPage.Page
	yBeforeLayout := e.currentY
	fmt.Printf("Debug (OldStrategy): Preparing to call layout function for %d pics. Page: %d, CurrentY: %.2f, AvailableH: %.2f\n", len(pictures), pageBeforeLayout, yBeforeLayout, layoutAvailableHeight)

	// 5. Call the appropriate layout processing function based on picture count
	var actualHeightUsed float64
	switch numPics {
	case 1:
		fmt.Println("Debug (OldStrategy): Calling handler for 1 picture.")
		actualHeightUsed = e.processSinglePictureLayoutAndPlace(pictures[0], layoutAvailableHeight)
	case 2:
		fmt.Println("Debug (OldStrategy): Calling handler for 2 pictures.")
		actualHeightUsed = e.processTwoPicturesLayoutAndPlace(pictures, layoutAvailableHeight)
	default: // 3 or more pictures
		fmt.Printf("Debug (OldStrategy): Calling processTemplatedLayoutAndPlace for %d pictures.\n", numPics)
		actualHeightUsed = e.processTemplatedLayoutAndPlace(pictures, layoutAvailableHeight) // Handles 3-9 and signals >9
	}
	// +++ 添加日志：记录调用后的页面、Y坐标和返回的高度 +++
	pageAfterLayout := e.currentPage.Page
	yAfterLayout := e.currentY
	fmt.Printf("Debug (OldStrategy): Returned from layout handler. Returned Height: %.2f. Page Before: %d, Page After: %d. Y Before: %.2f, Y After (engine state before update): %.2f\n", actualHeightUsed, pageBeforeLayout, pageAfterLayout, yBeforeLayout, yAfterLayout)

	// 6. Update Y Coordinate by Actual Placed Height
	// The layout function (single, two, or templated via specific functions)
	// calculates the layout and returns the total height it *should* occupy.
	// The placement happens relative to the currentY *before* this update.
	// We now update currentY by the height returned, assuming placement was successful.
	if actualHeightUsed > 0 { // Only update if a valid height was returned (not 0 or error codes like -2.0)
		e.currentY += actualHeightUsed
	} else if actualHeightUsed == -2.0 {
		// Handle split signal if necessary (though pagination should ideally prevent this call)
		fmt.Println("Debug (OldStrategy): Split signal received from layout function. CurrentY not updated.")
		// Potentially need logic here if a split during layout requires specific state changes.
	} else {
		// Handle other errors or zero height cases
		fmt.Printf("Debug (OldStrategy): Layout function returned non-positive height (%.2f). CurrentY not updated.\n", actualHeightUsed)
	}

	// +++ 添加日志：记录最终更新后的Y坐标 +++
	fmt.Printf("Debug (OldStrategy): Final CurrentY after potential update: %.2f\n", e.currentY)

	// Spacing *after* pictures is handled by the *next* element/entry's requiredSpacingBeforeElement check.
}

// calculateUniformRowHeightLayout calculates dimensions for a row aiming for uniform height,
// fitting within the availableWidth. It returns the calculated widths for each picture,
// (unused heights array), and the final uniform row height.
// --- Entire function calculateUniformRowHeightLayout removed from line 786 to 859 ---

// requiredSpacingBeforeElement checks if spacing is needed before the next element and returns the amount.
// It no longer modifies e.currentY directly.
func (e *ContinuousLayoutEngine) requiredSpacingBeforeElement() float64 {
	// Only add spacing if we are not at the very top of the page
	if e.currentY > e.marginTop {
		// Check if the current entry already has content.
		if len(e.currentPage.Entries) > 0 {
			lastEntry := &e.currentPage.Entries[len(e.currentPage.Entries)-1]
			// Check if the last entry has *any* content (time, text, or pictures)
			hasPreviousContent := len(lastEntry.TimeArea) > 0 || len(lastEntry.TextAreas) > 0 || len(lastEntry.Pictures) > 0

			// Only need spacing if there was previous content *and* we are not at the exact top margin
			if hasPreviousContent {
				return e.elementSpacing
			}
		}
	}
	return 0 // No spacing needed
}

// processTemplatedLayoutAndPlace handles layout for 3+ pictures using dedicated functions.
func (e *ContinuousLayoutEngine) processTemplatedLayoutAndPlace(pictures []Picture, layoutAvailableHeight float64) float64 {
	numPics := len(pictures)
	if numPics < 3 {
		fmt.Printf("Error: processTemplatedLayoutAndPlace called with %d pictures. Needs >= 3. Skipping.\n", numPics)
		return 0
	}

	switch numPics {
	case 3:
		// Call the specific function for 3 pictures
		return e.processLayoutForThreePictures(pictures, layoutAvailableHeight)
	case 4:
		// Call the specific function for 4 pictures
		return e.processLayoutForFourPictures(pictures, layoutAvailableHeight)
	case 5:
		// Call the specific function for 5 pictures
		return e.processLayoutForFivePictures(pictures, layoutAvailableHeight)
	case 6:
		// Call the specific function for 6 pictures
		return e.processLayoutForSixPictures(pictures, layoutAvailableHeight)
	case 7:
		// Call the specific function for 7 pictures
		return e.processLayoutForSevenPictures(pictures, layoutAvailableHeight)
	case 8:
		// Call the specific function for 8 pictures
		return e.processLayoutForEightPictures(pictures, layoutAvailableHeight)
	case 9:
		// Call the specific function for 9 pictures
		return e.processLayoutForNinePictures(pictures, layoutAvailableHeight)
	default:
		// For > 9 pictures, we currently don't have specific layouts. Signal split.
		fmt.Printf("Debug: No specific layout defined for %d pictures. Signaling split by returning error value -2.0.\n", numPics)
		// Treat this as a split signal, similar to how calculation functions might return split_required error.
		return -2.0 // Signal split_required
	}

	// Note: The common error handling and placement logic previously here is now integrated
	// into the specific processLayoutForXPictures functions or handled by the return values.
}

// placePicturesInTemplate places pictures based on the calculated template layout.
func (e *ContinuousLayoutEngine) placePicturesInTemplate(pictures []Picture, layout TemplateLayout) {
	if len(pictures) != len(layout.Positions) || len(pictures) != len(layout.Dimensions) {
		fmt.Println("Error: Mismatch between picture count and layout information in placePicturesInTemplate.")
		return
	}

	// Ensure entry exists
	if len(e.currentPage.Entries) == 0 {
		// This should ideally not happen if placement is always preceded by entry creation/selection.
		// If it does, create a temporary placeholder entry. This might indicate a logic flaw elsewhere.
		fmt.Println("Warning: placePicturesInTemplate called with no current entry on the page. Creating one.")
		e.currentPage.Entries = append(e.currentPage.Entries, PageEntry{})
	}
	currentEntry := &e.currentPage.Entries[len(e.currentPage.Entries)-1]
	startY := e.currentY // Top Y coordinate for the *entire* layout block

	// +++ Log StartY +++
	fmt.Printf("Debug (Place): Inside placePicturesInTemplate. startY (e.currentY): %.2f\n", startY)
	// +++ End Log +++

	// --- Calculate Actual Width and Centering Offset ---
	actualScaledWidth := 0.0
	for i := range layout.Positions {
		if len(layout.Positions[i]) == 2 && len(layout.Dimensions[i]) == 2 {
			rightEdge := layout.Positions[i][0] + layout.Dimensions[i][0]
			if rightEdge > actualScaledWidth {
				actualScaledWidth = rightEdge
			}
		}
	}

	offsetX := 0.0
	if actualScaledWidth < e.availableWidth {
		offsetX = (e.availableWidth - actualScaledWidth) / 2.0
	}
	if offsetX < 0 { // Safety check
		offsetX = 0
	}
	// --- End Calculation ---

	for i, pic := range pictures {
		relativeX := layout.Positions[i][0]
		relativeY := layout.Positions[i][1]
		width := layout.Dimensions[i][0]
		height := layout.Dimensions[i][1]

		// Apply centering offset
		absX0 := e.marginLeft + offsetX + relativeX
		absY0 := startY + relativeY
		absX1 := absX0 + width
		absY1 := absY0 + height

		area := [][]float64{
			{absX0, absY0},
			{absX1, absY1},
		}
		// Ensure Pictures slice is initialized if nil
		if currentEntry.Pictures == nil {
			currentEntry.Pictures = make([]Picture, 0, len(pictures))
		}
		currentEntry.Pictures = append(currentEntry.Pictures, Picture{
			Index:  pic.Index, // Preserve original index if available
			Area:   area,
			URL:    pic.URL,
			Width:  int(math.Round(width)),  // Store final rounded layout width
			Height: int(math.Round(height)), // Store final rounded layout height
		})
	}
	// currentY updated by caller (e.g., processTemplatedLayoutAndPlace)
}
