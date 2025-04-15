package backend

import (
	"fmt"
	"math"
	"strings"
	"time"
	"unicode/utf8"
)

// Layout represents a picture layout configuration
type Layout struct {
}

// ContinuousLayoutPage represents a single page in the continuous layout
type ContinuousLayoutPage struct {
	Page      int         `json:"page"`
	IsInsert  bool        `json:"is_insert"`  // 是否是插页
	YearMonth string      `json:"year_month"` // 年月信息，格式：2025年3月
	Entries   []PageEntry `json:"entries"`
}

// PageEntry represents a single entry's layout information on a page
type PageEntry struct {
	Time      string        `json:"time"`      // 格式：2025年3月30日 17:50
	DatePart  string        `json:"date_part"` // 格式：5月23日 周一
	TimePart  string        `json:"time_part"` // 格式：08:28
	TimeArea  [][]float64   `json:"time_area"`
	TextAreas [][][]float64 `json:"text_areas"`
	Texts     []string      `json:"texts"`
	Pictures  []Picture     `json:"pictures"`
}

// ContinuousLayoutEngine represents the continuous layout engine
type ContinuousLayoutEngine struct {
	entries            []Entry
	pages              []ContinuousLayoutPage
	currentPage        *ContinuousLayoutPage
	marginLeft         float64
	marginRight        float64
	marginTop          float64
	marginBottom       float64
	availableWidth     float64
	availableHeight    float64
	timeHeight         float64
	fontSize           float64
	lineHeight         float64
	currentY           float64
	timeAreaBottom     float64
	entrySpacing       float64 // 条目之间的间距
	elementSpacing     float64 // 元素整体之间的间距
	imageSpacing       float64 // 图片之间的间距
	minWideHeight      float64 // Min height for Wide pics (AR >= 3)
	minTallHeight      float64 // Min height for Tall pics (AR <= 1/3)
	minLandscapeHeight float64 // Min height for Landscape pics (1 < AR < 3)
	minPortraitHeight  float64 // Min height for Portrait pics (1/3 < AR < 1)
	singleImageHeight  float64 // 单张竖图的默认高度
	singleImageWidth   float64 // 单张横图的默认宽度
	minImageHeight     float64 // 单张竖图的最小展示高度
	minImageWidth      float64 // 单张横图的最小展示宽度
	currentYearMonth   string
	bottomMargin       float64
}

// Helper function to get picture type based on Aspect Ratio (AR)
func getPictureType(aspectRatio float64) string {
	if aspectRatio >= 3.0 {
		return "wide"
	} else if aspectRatio <= 1.0/3.0 {
		return "tall"
	} else if aspectRatio > 1.0 && aspectRatio < 3.0 {
		return "landscape"
	} else if aspectRatio > 1.0/3.0 && aspectRatio < 1.0 {
		return "portrait"
	} else if aspectRatio == 1.0 {
		return "square" // Added for completeness, rules don't specify minimums for square
	} else {
		return "unknown" // Should not happen with valid positive dimensions
	}
}

// NewContinuousLayoutEngine creates a new continuous layout engine
func NewContinuousLayoutEngine(entries []Entry) *ContinuousLayoutEngine {
	engine := &ContinuousLayoutEngine{
		entries:            entries,
		marginLeft:         142,
		marginRight:        142,
		marginTop:          189,
		marginBottom:       189,
		timeHeight:         100,
		fontSize:           66.67, // 对应72DPI的16px
		lineHeight:         100,   // 对应72DPI的24px
		entrySpacing:       150,   // 条目之间的间距
		elementSpacing:     30,    // 元素整体之间的间距
		imageSpacing:       15,    // 图片之间的间距
		minWideHeight:      400,   // Min height for Wide pics (AR >= 3)
		minTallHeight:      600,   // Min height for Tall pics (AR <= 1/3)
		minLandscapeHeight: 400,   // Min height for Landscape pics (1 < AR < 3)
		minPortraitHeight:  600,   // Min height for Portrait pics (1/3 < AR < 1)
		singleImageHeight:  2808,  // 设置单张竖图的默认高度
		singleImageWidth:   1695,  // 设置单张横图的默认宽度
		minImageHeight:     800,   // 设置单张竖图的最小展示高度
		minImageWidth:      1200,  // 设置单张横图的最小展示宽度
		bottomMargin:       100,   // 底部边距
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

// calculateEntryTotalHeight is potentially no longer needed for top-level pagination
// but might be useful elsewhere. Keep it for now, but remove its usage in ProcessEntries.
/*
func (e *ContinuousLayoutEngine) calculateEntryTotalHeight(entry Entry) float64 {
	// ... (previous implementation)
}
*/

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
	numPics := len(pictures)
	if numPics == 0 {
		return
	}

	// 1. Calculate Required Spacing Before Pictures
	requiredSpacing := e.requiredSpacingBeforeElement()

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
			picType := getPictureType(ar)
			switch picType {
			case "wide":
				minRequiredHeightForPics = e.minWideHeight
			case "tall":
				minRequiredHeightForPics = e.minTallHeight
			case "landscape":
				minRequiredHeightForPics = e.minLandscapeHeight
			case "portrait":
				minRequiredHeightForPics = e.minPortraitHeight
			default: // square or unknown
				minRequiredHeightForPics = e.minLandscapeHeight // Use landscape as fallback
			}
		case 2:
			// For 2 pics, use the generic min height for estimation/pagination decision.
			// The layout function itself will handle detailed placement later.
			minRequiredHeightForPics = e.minImageHeight // Or maybe min(minLandscape, minPortrait)? Using generic for now.
		default: // 3+ pics
			// Use template estimation or a generic minimum for pagination decision.
			minRequiredHeightForPics = e.minImageHeight // Using generic for now. Detailed height comes from template.
			// Ideally, we'd use estimatePicturesHeight here, but that involves calculations.
			// Let's stick to a simple minimum check for the pagination decision itself.
			// estimatedHeight := e.estimatePicturesHeight(pictures) // Can be complex
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

	// 4. Calculate Layout Available Height on the Target Page
	// This is the remaining space from currentY to the bottom margin boundary
	layoutAvailableHeight := (e.marginTop + e.availableHeight) - e.currentY
	if layoutAvailableHeight < 0 {
		layoutAvailableHeight = 0 // Avoid negative height if spacing pushed us over edge case
	}

	// 5. Call Layout & Placement Function
	var actualHeightUsed float64 = 0
	switch len(pictures) {
	case 1:
		actualHeightUsed = e.processSinglePictureLayoutAndPlace(pictures[0], layoutAvailableHeight)
	case 2:
		actualHeightUsed = e.processTwoPicturesLayoutAndPlace(pictures, layoutAvailableHeight)
	default:
		actualHeightUsed = e.processTemplatedLayoutAndPlace(pictures, layoutAvailableHeight)
	}

	// 6. Update Y Coordinate by Actual Placed Height
	// Y was already incremented by requiredSpacing before layout call
	e.currentY += actualHeightUsed
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

// processTemplatedLayout handles layout for 3+ pictures using predefined templates.
// --- Actual implementation for the refactored function ---
func (e *ContinuousLayoutEngine) processTemplatedLayoutAndPlace(pictures []Picture, layoutAvailableHeight float64) float64 {
	numPics := len(pictures)
	if numPics < 3 {
		fmt.Printf("Error: processTemplatedLayoutAndPlace called with %d pictures. Skipping.\n", numPics)
		return 0
	}

	// --- Select Template and Calculate Layout ---
	var layoutInfo TemplateLayout // Define a struct to hold calculated layout details
	var err error

	switch numPics {
	case 3:
		layoutInfo, err = e.calculateThreePicturesLayout(pictures, layoutAvailableHeight)
		// No split/force logic needed for 3 pics as it's the base case for splits
	case 4:
		layoutInfo, err = e.calculateFourPicturesLayout(pictures, layoutAvailableHeight)
		// --- Special handling for 4-pic ---
		if err != nil {
			switch err.Error() {
			case "force_new_page":
				fmt.Println("Info: Forcing new page for 4-picture layout due to wide/tall images not fitting.")
				e.newPage()
				layoutAvailableHeight = (e.marginTop + e.availableHeight) - e.currentY // Recalculate for new page
				// Retry calculation on the new page
				layoutInfo, err = e.calculateFourPicturesLayout(pictures, layoutAvailableHeight)
				if err != nil {
					// If it still fails (e.g., split required even on a full page, or other calc error)
					fmt.Printf("Error calculating 4-pic layout even on new page: %v. Skipping.\n", err)
					// We can't proceed with placement if the calculation failed on the new page
					return 0
				}
				// If calculation on new page succeeded, err is now nil,
				// and flow will continue to the common placement logic below the switch.

			case "split_required":
				fmt.Println("Info: Splitting 4-picture layout across pages.")
				// Check if there's enough space for even the first two using 2-pic logic's minimums
				minHeightForTwo := e.minImageHeight // Use the generic min height for 2-pic check
				if layoutAvailableHeight < minHeightForTwo {
					fmt.Printf("Error: Not enough space (%.2f) for even the first two pictures (min %.2f) during split. Skipping.\n", layoutAvailableHeight, minHeightForTwo)
					return 0 // Cannot place even the first part
				}

				// Place first two pictures
				heightUsed1 := e.processTwoPicturesLayoutAndPlace(pictures[0:2], layoutAvailableHeight)
				if heightUsed1 <= 1e-6 { // Check if first placement failed
					fmt.Println("Error: Failed to place first two pictures during split. Skipping rest.")
					return 0
				}
				e.currentY += heightUsed1 // Update Y after first successful placement

				// Start a new page for the next two
				e.newPage()
				newLayoutAvailableHeight := (e.marginTop + e.availableHeight) - e.currentY

				// Place last two pictures on the new page
				heightUsed2 := e.processTwoPicturesLayoutAndPlace(pictures[2:4], newLayoutAvailableHeight)
				if heightUsed2 <= 1e-6 { // Check if second placement failed
					fmt.Println("Error: Failed to place last two pictures during split. First part placed, but operation incomplete.")
					// Return 0 height used for *this call* because the second part failed.
					return 0
				}
				// If second placement succeeds, return the height used by the second part.
				// The caller (processPictures) will add this heightUsed2 to the new page's currentY.
				return heightUsed2 // Return immediately after successful split placement

			default:
				// Handle other calculation errors (not split_required or force_new_page)
				fmt.Printf("Error calculating layout for 4 pictures: %v. Skipping placement.\n", err)
				return 0 // Return immediately as we cannot place
			}
		}
		// If err was nil initially, or after force_new_page recalculation, flow continues...
	case 5:
		layoutInfo, err = e.calculateFivePicturesLayout(pictures, layoutAvailableHeight)
		// --- Special handling for 5-pic ---
		if err != nil {
			switch err.Error() {
			case "force_new_page":
				fmt.Println("Info: Forcing new page for 5-picture layout due to wide/tall images not fitting.")
				e.newPage()
				layoutAvailableHeight = (e.marginTop + e.availableHeight) - e.currentY // Recalculate for new page
				// Retry calculation on the new page
				layoutInfo, err = e.calculateFivePicturesLayout(pictures, layoutAvailableHeight)
				if err != nil {
					fmt.Printf("Error calculating 5-pic layout even on new page: %v. Skipping.\n", err)
					return 0
				}
				// If retry succeeds, err is nil, flow continues below

			case "split_required":
				fmt.Println("Info: Splitting 5-picture layout across pages (3+2).")
				// Estimate minimum height for the first group (3 pictures)
				// Using minImageHeight as a basic check, more accurate would require partial calculation.
				minHeightForThree := e.minImageHeight
				if layoutAvailableHeight < minHeightForThree {
					fmt.Printf("Error: Not enough space (%.2f) for even the first three pictures (min %.2f est.) during 5-pic split. Skipping.\n", layoutAvailableHeight, minHeightForThree)
					return 0
				}
				// Place first three pictures - Recursive call to handle 3-pic case
				heightUsed1 := e.processTemplatedLayoutAndPlace(pictures[0:3], layoutAvailableHeight)
				if heightUsed1 <= 1e-6 {
					fmt.Println("Error: Failed to place first three pictures during 5-pic split. Skipping rest.")
					return 0
				}
				e.currentY += heightUsed1
				// New page for the next two
				e.newPage()
				newLayoutAvailableHeight := (e.marginTop + e.availableHeight) - e.currentY
				// Place last two pictures
				heightUsed2 := e.processTwoPicturesLayoutAndPlace(pictures[3:5], newLayoutAvailableHeight)
				if heightUsed2 <= 1e-6 {
					fmt.Println("Error: Failed to place last two pictures during 5-pic split. First part placed, but operation incomplete.")
					return 0
				}
				return heightUsed2 // Return height used on the *new* page

			default:
				fmt.Printf("Error calculating layout for 5 pictures: %v. Skipping placement.\n", err)
				return 0
			}
		}
		// If err was nil initially, or after force_new_page recalculation, flow continues...
	case 6:
		layoutInfo, err = e.calculateSixPicturesLayout(pictures)
		// TODO: Add split/force handling for 6 pictures (e.g., 3+3 or 4+2)
	case 7:
		layoutInfo, err = e.calculateSevenPicturesLayout(pictures)
		// TODO: Add split/force handling for 7 pictures (e.g., 3+4 or 4+3)
	case 8:
		layoutInfo, err = e.calculateEightPicturesLayout(pictures)
		// TODO: Add split/force handling for 8 pictures (e.g., 4+4)
	case 9:
		layoutInfo, err = e.calculateNinePicturesLayout(pictures)
		// TODO: Add split/force handling for 9 pictures (e.g., 3+3+3 or 4+5)
	default:
		err = fmt.Errorf("template layout not implemented for %d pictures", numPics)
	}

	// --- Common Error Handling & Placement for Non-Split/Non-Force Cases ---
	if err != nil {
		// This handles errors from cases 3, 6, 7, 8, 9 or default case,
		// or calculation errors not caught by specific split/force handlers,
		// or errors from the second calculation attempt after force_new_page.
		fmt.Printf("Error calculating layout for %d pictures (final check): %v. Skipping placement.\n", numPics, err)
		return 0
	}

	// If we reach here, err is nil, and layoutInfo is valid for placement on the current page.
	e.placePicturesInTemplate(pictures, layoutInfo)
	return layoutInfo.TotalHeight
}

// TemplateLayout holds the calculated positions and dimensions for a template
type TemplateLayout struct {
	Positions   [][]float64 // Relative positions [x, y] for top-left corner of each pic within the layout block
	Dimensions  [][]float64 // Dimensions [width, height] for each pic
	TotalHeight float64     // Total height of the layout block (including internal spacing)
	TotalWidth  float64     // Total width (should generally match e.availableWidth)
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

// --- Helper function for 1 Top, 2 Bottom Stacked Template ---
func (e *ContinuousLayoutEngine) calculateLayout_1T2B(ARs []float64, types []string, AW, spacing float64) (TemplateLayout, bool, error) {
	layout := TemplateLayout{
		Positions:  make([][]float64, 3),
		Dimensions: make([][]float64, 3),
	}
	AR0, AR1, AR2 := ARs[0], ARs[1], ARs[2]

	// Pic 0 takes full width
	W0 := AW
	if W0 <= 1e-6 {
		return layout, false, fmt.Errorf("1T2B available width is zero")
	}
	H0 := 0.0
	if AR0 > 1e-6 {
		H0 = W0 / AR0
	}
	if H0 <= 1e-6 {
		return layout, false, fmt.Errorf("1T2B calculated zero height for top picture")
	}

	// Calculate bottom row height based on fitting pics 1 & 2 in available width
	bottomRowAvailableWidth := AW - spacing
	bottomTotalARSum := AR1 + AR2
	H_bottom := 0.0
	if bottomRowAvailableWidth > 1e-6 && bottomTotalARSum > 1e-6 {
		H_bottom = bottomRowAvailableWidth / bottomTotalARSum
	} else {
		return layout, false, fmt.Errorf("cannot calculate 1T2B bottom row height")
	}
	if H_bottom <= 1e-6 {
		return layout, false, fmt.Errorf("1T2B calculated zero height for bottom row")
	}

	W1 := H_bottom * AR1
	W2 := H_bottom * AR2

	layout.TotalHeight = H0 + spacing + H_bottom
	layout.TotalWidth = AW

	layout.Positions[0] = []float64{0, 0} // Pic 0 Top
	layout.Dimensions[0] = []float64{W0, H0}

	layout.Positions[1] = []float64{0, H0 + spacing} // Pic 1 Bottom Left
	layout.Dimensions[1] = []float64{W1, H_bottom}

	layout.Positions[2] = []float64{W1 + spacing, H0 + spacing} // Pic 2 Bottom Right
	layout.Dimensions[2] = []float64{W2, H_bottom}

	// Type-specific minimum height check
	meetsMin := true
	// heights := []float64{H0, H_bottom, H_bottom} // Heights for pics 0, 1, 2 respectively - Not strictly needed for check logic
	requiredMinHeights := make([]float64, 3)
	for i, picType := range types {
		switch picType {
		case "wide":
			requiredMinHeights[i] = e.minWideHeight
		case "tall":
			requiredMinHeights[i] = e.minTallHeight
		case "landscape":
			requiredMinHeights[i] = e.minLandscapeHeight
		case "portrait":
			requiredMinHeights[i] = e.minPortraitHeight
		default:
			requiredMinHeights[i] = e.minLandscapeHeight // Fallback
		}
		// Use correct height for check (H0 for pic 0, H_bottom for pics 1 & 2)
		checkHeight := H0
		if i > 0 {
			checkHeight = H_bottom
		}
		if checkHeight < requiredMinHeights[i] {
			meetsMin = false
			break
		}
	}
	// No explicit min width check based on rules

	return layout, meetsMin, nil
}

// --- Helper function for 3 in a Row Template ---
func (e *ContinuousLayoutEngine) calculateLayout_3Row(ARs []float64, types []string, AW, spacing float64) (TemplateLayout, bool, error) {
	layout := TemplateLayout{
		Positions:  make([][]float64, 3),
		Dimensions: make([][]float64, 3),
	}
	AR0, AR1, AR2 := ARs[0], ARs[1], ARs[2]

	rowAvailableWidth := AW - 2*spacing
	totalARSum := AR0 + AR1 + AR2
	H := 0.0
	if rowAvailableWidth > 1e-6 && totalARSum > 1e-6 {
		H = rowAvailableWidth / totalARSum
	} else {
		return layout, false, fmt.Errorf("cannot calculate 3-in-a-row layout (zero width or AR sum)")
	}
	if H <= 1e-6 {
		return layout, false, fmt.Errorf("3-in-a-row calculated zero height")
	}

	W0 := H * AR0
	W1 := H * AR1
	W2 := H * AR2
	widths := []float64{W0, W1, W2}

	layout.TotalHeight = H
	layout.TotalWidth = AW
	currentX := 0.0
	for i := 0; i < 3; i++ {
		layout.Positions[i] = []float64{currentX, 0}
		layout.Dimensions[i] = []float64{widths[i], H}
		if i < 2 {
			currentX += widths[i] + spacing
		}
	}

	// Type-specific minimum height check
	meetsMin := true
	requiredMinHeights := make([]float64, 3)
	for i, picType := range types {
		switch picType {
		case "wide":
			requiredMinHeights[i] = e.minWideHeight
		case "tall":
			requiredMinHeights[i] = e.minTallHeight
		case "landscape":
			requiredMinHeights[i] = e.minLandscapeHeight
		case "portrait":
			requiredMinHeights[i] = e.minPortraitHeight
		default:
			requiredMinHeights[i] = e.minLandscapeHeight // Fallback
		}
		// All pictures have the same height H in this layout
		if H < requiredMinHeights[i] {
			meetsMin = false
			break
		}
	}
	// No explicit min width check based on rules

	return layout, meetsMin, nil
}

// --- Helper function for 2 Left Stacked, 1 Right Template --- (Mirror of 1L2R)
func (e *ContinuousLayoutEngine) calculateLayout_2L1R(ARs []float64, types []string, AW, spacing float64) (TemplateLayout, bool, error) {
	layout := TemplateLayout{
		Positions:  make([][]float64, 3),
		Dimensions: make([][]float64, 3),
	}
	AR0, AR1, AR2 := ARs[0], ARs[1], ARs[2]

	// Calculate WL based on geometry (similar to WR in 1L2R)
	denominator := 1.0
	if AR0 > 1e-6 {
		denominator += AR2 / AR0
	}
	if AR1 > 1e-6 {
		denominator += AR2 / AR1
	}

	WL := 0.0
	if denominator > 1e-6 {
		WL = (AW - spacing*(1.0+AR2)) / denominator
	}

	if WL <= 1e-6 || WL > AW-spacing+1e-6 {
		return layout, false, fmt.Errorf("2L1R geometry infeasible (WL=%.2f)", WL)
	}

	W2 := AW - spacing - WL
	if W2 <= 1e-6 {
		return layout, false, fmt.Errorf("2L1R geometry infeasible (W2=%.2f)", W2)
	}

	H0 := 0.0
	if AR0 > 1e-6 {
		H0 = WL / AR0
	}
	H1 := 0.0
	if AR1 > 1e-6 {
		H1 = WL / AR1
	}
	H2 := 0.0
	if AR2 > 1e-6 {
		H2 = W2 / AR2
	}

	if H0 <= 1e-6 || H1 <= 1e-6 || H2 <= 1e-6 {
		return layout, false, fmt.Errorf("2L1R calculated zero height")
	}

	layout.TotalHeight = H2
	layout.TotalWidth = AW
	layout.Positions[0] = []float64{0, 0}
	layout.Dimensions[0] = []float64{WL, H0}
	layout.Positions[1] = []float64{0, H0 + spacing}
	layout.Dimensions[1] = []float64{WL, H1}
	layout.Positions[2] = []float64{WL + spacing, 0}
	layout.Dimensions[2] = []float64{W2, H2}

	// Type-specific minimum height check
	meetsMin := true
	heights := []float64{H0, H1, H2} // Heights for pics 0, 1, 2
	requiredMinHeights := make([]float64, 3)
	for i, picType := range types {
		switch picType {
		case "wide":
			requiredMinHeights[i] = e.minWideHeight
		case "tall":
			requiredMinHeights[i] = e.minTallHeight
		case "landscape":
			requiredMinHeights[i] = e.minLandscapeHeight
		case "portrait":
			requiredMinHeights[i] = e.minPortraitHeight
		default:
			requiredMinHeights[i] = e.minLandscapeHeight // Fallback
		}
		if heights[i] < requiredMinHeights[i] {
			meetsMin = false
			break
		}
	}
	// No explicit min width check based on rules

	return layout, meetsMin, nil
}

// --- Helper function for 2 Top, 1 Bottom Full Width Template ---
func (e *ContinuousLayoutEngine) calculateLayout_2T1B(ARs []float64, types []string, AW, spacing float64) (TemplateLayout, bool, error) {
	layout := TemplateLayout{
		Positions:  make([][]float64, 3),
		Dimensions: make([][]float64, 3),
	}
	AR0, AR1, AR2 := ARs[0], ARs[1], ARs[2]

	// Calculate top row height (H_top) based on fitting pics 0 & 1 in available width
	topRowAvailableWidth := AW - spacing
	topTotalARSum := AR0 + AR1
	H_top := 0.0
	if topRowAvailableWidth > 1e-6 && topTotalARSum > 1e-6 {
		H_top = topRowAvailableWidth / topTotalARSum
	} else {
		return layout, false, fmt.Errorf("cannot calculate 2T1B top row height")
	}
	if H_top <= 1e-6 {
		return layout, false, fmt.Errorf("2T1B calculated zero height for top row")
	}

	W0 := H_top * AR0
	W1 := H_top * AR1

	// Pic 2 takes full width at bottom
	W2 := AW
	H2 := 0.0
	if AR2 > 1e-6 {
		H2 = W2 / AR2
	}
	if H2 <= 1e-6 {
		return layout, false, fmt.Errorf("2T1B calculated zero height for bottom picture")
	}

	layout.TotalHeight = H_top + spacing + H2
	layout.TotalWidth = AW

	layout.Positions[0] = []float64{0, 0} // Pic 0 Top Left
	layout.Dimensions[0] = []float64{W0, H_top}
	layout.Positions[1] = []float64{W0 + spacing, 0} // Pic 1 Top Right
	layout.Dimensions[1] = []float64{W1, H_top}
	layout.Positions[2] = []float64{0, H_top + spacing} // Pic 2 Bottom Full
	layout.Dimensions[2] = []float64{W2, H2}

	// Type-specific minimum height check
	meetsMin := true
	// heights := []float64{H_top, H_top, H2} // Heights for pics 0, 1, 2 - Not strictly needed
	requiredMinHeights := make([]float64, 3)
	for i, picType := range types {
		switch picType {
		case "wide":
			requiredMinHeights[i] = e.minWideHeight
		case "tall":
			requiredMinHeights[i] = e.minTallHeight
		case "landscape":
			requiredMinHeights[i] = e.minLandscapeHeight
		case "portrait":
			requiredMinHeights[i] = e.minPortraitHeight
		default:
			requiredMinHeights[i] = e.minLandscapeHeight // Fallback
		}
		// Use correct height (H_top for 0, 1; H2 for 2)
		checkHeight := H_top
		if i == 2 {
			checkHeight = H2
		}
		if checkHeight < requiredMinHeights[i] {
			meetsMin = false
			break
		}
	}
	// No explicit min width check based on rules

	return layout, meetsMin, nil
}

// --- Helper function for 3 Top, 4 Bottom Template ---
func (e *ContinuousLayoutEngine) calculateLayout_3T4B(pictures []Picture, ARs []float64, AW, spacing, minH, minW float64) (TemplateLayout, bool, error) {
	layout := TemplateLayout{
		Positions:  make([][]float64, 7),
		Dimensions: make([][]float64, 7),
	}

	// 使用pictures参数获取图片分组，因为calculateUniformRowHeightLayout需要Picture对象
	row1Pics := pictures[0:3]
	row2Pics := pictures[3:7]

	widths1, _, height1 := e.calculateUniformRowHeightLayout(row1Pics, AW)
	widths2, _, height2 := e.calculateUniformRowHeightLayout(row2Pics, AW)

	if height1 <= 1e-6 || height2 <= 1e-6 {
		return layout, false, fmt.Errorf("failed to calculate row layouts for 3T4B")
	}
	W0, W1, W2 := widths1[0], widths1[1], widths1[2]
	W3, W4, W5, W6 := widths2[0], widths2[1], widths2[2], widths2[3]

	// --- Check Minimums ---*
	meetsMin := true
	if height1 < minH || height2 < minH {
		meetsMin = false
	}
	if W0 < minW || W1 < minW || W2 < minW || W3 < minW || W4 < minW || W5 < minW || W6 < minW {
		meetsMin = false
	}

	// --- Populate Layout Struct ---*
	layout.TotalHeight = height1 + spacing + height2
	layout.TotalWidth = AW
	// Row 1
	currentX := 0.0
	layout.Positions[0] = []float64{currentX, 0}
	layout.Dimensions[0] = []float64{W0, height1}
	currentX += W0 + spacing
	layout.Positions[1] = []float64{currentX, 0}
	layout.Dimensions[1] = []float64{W1, height1}
	currentX += W1 + spacing
	layout.Positions[2] = []float64{currentX, 0}
	layout.Dimensions[2] = []float64{W2, height1}
	// Row 2
	currentX = 0.0
	layout.Positions[3] = []float64{currentX, height1 + spacing}
	layout.Dimensions[3] = []float64{W3, height2}
	currentX += W3 + spacing
	layout.Positions[4] = []float64{currentX, height1 + spacing}
	layout.Dimensions[4] = []float64{W4, height2}
	currentX += W4 + spacing
	layout.Positions[5] = []float64{currentX, height1 + spacing}
	layout.Dimensions[5] = []float64{W5, height2}
	currentX += W5 + spacing
	layout.Positions[6] = []float64{currentX, height1 + spacing}
	layout.Dimensions[6] = []float64{W6, height2}

	return layout, meetsMin, nil
}

// --- Helper function for 4 Top, 3 Bottom Template ---
func (e *ContinuousLayoutEngine) calculateLayout_4T3B(pictures []Picture, ARs []float64, AW, spacing, minH, minW float64) (TemplateLayout, bool, error) {
	// Added definition based on expected structure
	layout := TemplateLayout{
		Positions:  make([][]float64, 7),
		Dimensions: make([][]float64, 7),
	}

	// Need pictures slice for calculateUniformRowHeightLayout
	if len(pictures) != 7 {
		return layout, false, fmt.Errorf("calculateLayout_4T3B requires 7 pictures")
	}
	row1Pics := pictures[0:4]
	row2Pics := pictures[4:7]

	widths1, _, height1 := e.calculateUniformRowHeightLayout(row1Pics, AW)
	widths2, _, height2 := e.calculateUniformRowHeightLayout(row2Pics, AW)

	if height1 <= 1e-6 || height2 <= 1e-6 {
		return layout, false, fmt.Errorf("failed to calculate row layouts for 4T3B")
	}
	W0, W1, W2, W3 := widths1[0], widths1[1], widths1[2], widths1[3]
	W4, W5, W6 := widths2[0], widths2[1], widths2[2]

	// --- Check Minimums ---
	meetsMin := true
	if height1 < minH || height2 < minH {
		meetsMin = false
	}
	if W0 < minW || W1 < minW || W2 < minW || W3 < minW || W4 < minW || W5 < minW || W6 < minW {
		meetsMin = false
	}

	// --- Populate Layout Struct ---
	layout.TotalHeight = height1 + spacing + height2
	layout.TotalWidth = AW
	// Row 1
	currentX := 0.0
	layout.Positions[0] = []float64{currentX, 0}
	layout.Dimensions[0] = []float64{W0, height1}
	currentX += W0 + spacing
	layout.Positions[1] = []float64{currentX, 0}
	layout.Dimensions[1] = []float64{W1, height1}
	currentX += W1 + spacing
	layout.Positions[2] = []float64{currentX, 0}
	layout.Dimensions[2] = []float64{W2, height1}
	currentX += W2 + spacing
	layout.Positions[3] = []float64{currentX, 0}
	layout.Dimensions[3] = []float64{W3, height1}
	// Row 2
	currentX = 0.0
	layout.Positions[4] = []float64{currentX, height1 + spacing}
	layout.Dimensions[4] = []float64{W4, height2}
	currentX += W4 + spacing
	layout.Positions[5] = []float64{currentX, height1 + spacing}
	layout.Dimensions[5] = []float64{W5, height2}
	currentX += W5 + spacing
	layout.Positions[6] = []float64{currentX, height1 + spacing}
	layout.Dimensions[6] = []float64{W6, height2}

	return layout, meetsMin, nil
}

// --- Helper function for 1 Left, 2 Right Stacked Template ---
func (e *ContinuousLayoutEngine) calculateLayout_1L2R(ARs []float64, types []string, AW, spacing float64) (TemplateLayout, bool, error) {
	// Added definition based on layout_3_pics.go version
	layout := TemplateLayout{
		Positions:  make([][]float64, 3),
		Dimensions: make([][]float64, 3),
	}
	AR0 := ARs[0]
	denominator := 1.0
	if ARs[1] > 1e-6 {
		denominator += AR0 / ARs[1]
	}
	if ARs[2] > 1e-6 {
		denominator += AR0 / ARs[2]
	}
	WR := 0.0
	if denominator > 1e-6 {
		WR = (AW - spacing*(1.0+AR0)) / denominator
	}
	if WR <= 1e-6 || WR > AW-spacing+1e-6 {
		return layout, false, fmt.Errorf("1L2R geometry infeasible (WR=%.2f)", WR)
	}
	W0 := AW - spacing - WR
	if W0 <= 1e-6 {
		return layout, false, fmt.Errorf("1L2R geometry infeasible (W0=%.2f)", W0)
	}
	H0 := 0.0
	if AR0 > 1e-6 {
		H0 = W0 / AR0
	}
	H1 := 0.0
	if ARs[1] > 1e-6 {
		H1 = WR / ARs[1]
	}
	H2 := 0.0
	if ARs[2] > 1e-6 {
		H2 = WR / ARs[2]
	}
	if H0 <= 1e-6 || H1 <= 1e-6 || H2 <= 1e-6 {
		return layout, false, fmt.Errorf("1L2R calculated zero height")
	}

	layout.TotalHeight = H0
	layout.TotalWidth = AW
	layout.Positions[0] = []float64{0, 0}
	layout.Dimensions[0] = []float64{W0, H0}
	layout.Positions[1] = []float64{W0 + spacing, 0}
	layout.Dimensions[1] = []float64{WR, H1}
	layout.Positions[2] = []float64{W0 + spacing, H1 + spacing}
	layout.Dimensions[2] = []float64{WR, H2}

	// Type-specific minimum height check
	meetsMin := true
	heights := []float64{H0, H1, H2} // Heights for pics 0, 1, 2
	requiredMinHeights := make([]float64, 3)
	for i, picType := range types {
		switch picType {
		case "wide":
			requiredMinHeights[i] = e.minWideHeight
		case "tall":
			requiredMinHeights[i] = e.minTallHeight
		case "landscape":
			requiredMinHeights[i] = e.minLandscapeHeight
		case "portrait":
			requiredMinHeights[i] = e.minPortraitHeight
		default:
			requiredMinHeights[i] = e.minLandscapeHeight // Fallback
		}
		if heights[i] < requiredMinHeights[i] {
			meetsMin = false
			break
		}
	}
	// No explicit min width check based on rules

	return layout, meetsMin, nil
}

// --- Helper function for 3 Columns (Vertical Stack) ---
func (e *ContinuousLayoutEngine) calculateLayout_3Col(ARs []float64, types []string, AW, spacing float64) (TemplateLayout, bool, error) {
	// Added definition based on layout_3_pics.go version
	layout := TemplateLayout{Positions: make([][]float64, 3), Dimensions: make([][]float64, 3)}
	AR0, AR1, AR2 := ARs[0], ARs[1], ARs[2]
	W0, W1, W2 := AW, AW, AW // Full width for each
	H0, H1, H2 := 0.0, 0.0, 0.0
	if AR0 > 1e-6 {
		H0 = W0 / AR0
	} else {
		return layout, false, fmt.Errorf("3Col calculated zero height for pic 0")
	}
	if AR1 > 1e-6 {
		H1 = W1 / AR1
	} else {
		return layout, false, fmt.Errorf("3Col calculated zero height for pic 1")
	}
	if AR2 > 1e-6 {
		H2 = W2 / AR2
	} else {
		return layout, false, fmt.Errorf("3Col calculated zero height for pic 2")
	}
	if H0 <= 1e-6 || H1 <= 1e-6 || H2 <= 1e-6 {
		return layout, false, fmt.Errorf("3Col calculated zero height")
	}
	heights := []float64{H0, H1, H2}
	layout.TotalHeight = H0 + H1 + H2 + 2*spacing
	layout.TotalWidth = AW
	currentY := 0.0
	for i := 0; i < 3; i++ {
		layout.Positions[i] = []float64{0, currentY}
		layout.Dimensions[i] = []float64{AW, heights[i]}
		if i < 2 {
			currentY += heights[i] + spacing
		}
	}

	// Type-specific minimum height check
	meetsMin := true
	requiredMinHeights := make([]float64, 3)
	for i, picType := range types {
		switch picType {
		case "wide":
			requiredMinHeights[i] = e.minWideHeight
		case "tall":
			requiredMinHeights[i] = e.minTallHeight
		case "landscape":
			requiredMinHeights[i] = e.minLandscapeHeight
		case "portrait":
			requiredMinHeights[i] = e.minPortraitHeight
		default:
			requiredMinHeights[i] = e.minLandscapeHeight // Fallback
		}
		if heights[i] < requiredMinHeights[i] {
			meetsMin = false
			break
		}
	}
	// No explicit min width check based on rules

	return layout, meetsMin, nil
}
