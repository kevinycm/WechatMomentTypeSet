package calculate

import (
	"errors"
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
		entries:        entries,
		marginLeft:     142,
		marginRight:    142,
		marginTop:      189,
		marginBottom:   189,
		timeHeight:     100,
		fontSize:       66.67, // 对应72DPI的16px
		lineHeight:     100,   // 对应72DPI的24px
		entrySpacing:   150,   // 条目之间的间距
		elementSpacing: 30,    // 元素整体之间的间距
		imageSpacing:   15,    // 图片之间的间距
		minWideHeight:  600,   // Min height for Wide pics (AR >= 3)
		minTallHeight:  800,   // Min height for Tall pics (AR <= 1/3)

		// Added slices for 1-9 pictures (index 0 unused)
		minLandscapeHeights: []float64{600, 600, 400, 600, 600, 600, 600, 600, 600}, // 横图 1-9 张
		minPortraitHeights:  []float64{800, 800, 600, 800, 800, 800, 800, 800, 800}, // 竖图 1-9 张

		singleImageHeight: 3130, // 设置单张竖图的最大高度
		singleImageWidth:  2124, // 设置单张横图的最大度
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

	// --- Check for Ultra-Wide or Ultra-Tall Pictures ---
	hasUltraWideOrTall := false
	ultraThreshold := 4.0
	for _, pic := range pictures {
		if pic.Height > 0 && pic.Width > 0 {
			ar := float64(pic.Width) / float64(pic.Height)
			if ar >= ultraThreshold || ar <= (1.0/ultraThreshold) { // Check for wide OR tall
				hasUltraWideOrTall = true
				picType := "ultra-wide"
				if ar <= (1.0 / ultraThreshold) {
					picType = "ultra-tall"
				}
				fmt.Printf("Debug (ProcessPics): %s picture detected (Index %d, AR %.2f). Will use dynamic layout.\n", picType, pic.Index, ar)
				break
			}
		}
	}

	// --- Choose Layout Strategy ---
	if hasUltraWideOrTall {
		// +++ Use NEW Dynamic Row-by-Row Strategy +++
		fmt.Println("Debug (ProcessPics): Using dynamic row layout strategy.")

		currentIndex := 0
		for currentIndex < numPicsTotal {
			// 1. Calculate Required Spacing Before this row
			requiredSpacing := e.requiredSpacingBeforeElement()

			// 2. Estimate minimum height for the *next potential row* (Simple estimate)
			estimatedMinHeight := 100.0 // Baseline estimate
			if currentIndex < numPicsTotal {
				pic := pictures[currentIndex]
				ar := 1.0
				if pic.Height > 0 {
					ar = float64(pic.Width) / float64(pic.Height)
				}
				estimatedMinHeight = GetRequiredMinHeight(e, GetPictureType(ar), 1)
			}

			// 3. Centralized Pagination Check (Before placing the row)
			pageRemainingPhysicalHeight := (e.marginTop + e.availableHeight) - e.currentY
			needsNewPage := false
			if e.currentY > e.marginTop { // Only paginate if not at top
				if pageRemainingPhysicalHeight-requiredSpacing < estimatedMinHeight {
					needsNewPage = true
				}
			}

			if needsNewPage {
				e.newPage()
				requiredSpacing = 0
			}

			// 4. Apply Spacing
			e.currentY += requiredSpacing

			// 5. Calculate Available Height for Row Placement
			layoutAvailableHeight := (e.marginTop + e.availableHeight) - e.currentY
			if layoutAvailableHeight <= 1e-6 { // Use tolerance
				fmt.Printf("Warning: No available height left (%.2f) on page %d for picture row starting at index %d. Skipping remaining pictures.\n", layoutAvailableHeight, e.currentPage.Page, currentIndex)
				break
			}

			// 6. Determine which pictures form the next row
			picsInNextRow, rowConfigType := e.determineNextPictureRow(pictures[currentIndex:], layoutAvailableHeight)
			numPicturesConsumed := len(picsInNextRow)

			if numPicturesConsumed == 0 {
				fmt.Printf("Warning: Could not determine layout for picture index %d. Skipping remaining pictures.\n", currentIndex)
				break
			}

			// 7. Calculate the actual layout for this specific row
			rowLayoutInfo, err := e.calculateRowLayout(picsInNextRow, rowConfigType, layoutAvailableHeight)
			if err != nil {
				fmt.Printf("Error calculating row layout for %d pictures (type: %s) starting at index %d: %v. Skipping row.\n", numPicturesConsumed, rowConfigType, currentIndex, err)
				// Simple strategy: skip the problematic picture(s) and try next
				currentIndex += numPicturesConsumed
				continue // Try the next iteration
			}

			// 8. Double check if calculated height fits (should be handled by calculateRowLayout ideally)
			if rowLayoutInfo.TotalHeight > layoutAvailableHeight+1e-6 {
				fmt.Printf("Error: Calculated row height (%.2f) exceeds available height (%.2f) for %d pics (type: %s) starting at index %d. Skipping row.\n", rowLayoutInfo.TotalHeight, layoutAvailableHeight, numPicturesConsumed, rowConfigType, currentIndex)
				currentIndex += numPicturesConsumed
				continue
			}

			// Handle case where layout calculation yields zero height (shouldn't happen ideally)
			if rowLayoutInfo.TotalHeight <= 1e-6 {
				fmt.Printf("Warning: Calculated row layout for %d pics (type: %s) starting at %d resulted in zero height. Skipping row.\n", numPicturesConsumed, rowConfigType, currentIndex)
				currentIndex += numPicturesConsumed
				continue
			}

			// 9. Place the pictures for this row
			e.placePicturesInRow(picsInNextRow, rowLayoutInfo)

			// 10. Update Y Coordinate
			e.currentY += rowLayoutInfo.TotalHeight

			// 11. Advance Index
			currentIndex += numPicturesConsumed
		} // End for loop

	} else {
		// +++ Use OLD Standard Templated Strategy +++
		fmt.Println("Debug (ProcessPics): No ultra-wide/tall pictures. Using standard templated layout strategy.")
		processPicturesOldStrategy(e, pictures)
	}
}

// --- Placeholder Signatures for New Dynamic Layout Helper Functions ---

// determineNextPictureRow looks at the remaining pictures and decides how many (and which)
// should form the next row based on rules (ultra-wide/tall, 3, 2, 1) and available height.
// Returns the pictures for the row and a string indicating the configuration type.
func (e *ContinuousLayoutEngine) determineNextPictureRow(remainingPics []Picture, availableHeight float64) (picsForNextRow []Picture, rowConfigType string) {
	numRemaining := len(remainingPics)
	if numRemaining == 0 {
		return []Picture{}, ""
	}

	ultraThreshold := 4.0
	// Helper function to check if a picture is ultra-wide or ultra-tall
	isUltra := func(pic Picture) bool {
		if pic.Height > 0 && pic.Width > 0 {
			ar := float64(pic.Width) / float64(pic.Height)
			return ar >= ultraThreshold || ar <= (1.0/ultraThreshold)
		}
		return false // Treat invalid dimensions as not ultra
	}

	// --- Decision Logic ---

	// 1. Check the first picture
	pic1 := remainingPics[0]
	if isUltra(pic1) {
		fmt.Println("Debug (determineRow): First pic is ultra, forming row-of-1-ultra.")
		return remainingPics[0:1], "row-of-1-ultra"
	}

	// 2. If pic1 is normal, try forming a row of 3
	if numRemaining >= 3 {
		pic2 := remainingPics[1]
		pic3 := remainingPics[2]
		if !isUltra(pic2) && !isUltra(pic3) {
			// All three are normal, form a row of 3
			fmt.Println("Debug (determineRow): First 3 pics are normal, forming row-of-3.")
			return remainingPics[0:3], "row-of-3"
		} else {
			fmt.Println("Debug (determineRow): Cannot form row-of-3 (pic 2 or 3 is ultra).")
		}
	} // Implicitly falls through if less than 3 remain or condition not met

	// 3. If row of 3 not formed, try forming a row of 2
	if numRemaining >= 2 {
		pic2 := remainingPics[1]
		if !isUltra(pic2) {
			// Both pic1 and pic2 are normal, form a row of 2
			fmt.Println("Debug (determineRow): First 2 pics are normal, forming row-of-2.")
			return remainingPics[0:2], "row-of-2"
		} else {
			fmt.Println("Debug (determineRow): Cannot form row-of-2 (pic 2 is ultra).")
		}
	} // Implicitly falls through if less than 2 remain or condition not met

	// 4. Default to row of 1 (since pic1 is known to be normal here)
	fmt.Println("Debug (determineRow): Defaulting to row-of-1 (normal pic).")
	return remainingPics[0:1], "row-of-1"
}

// calculateRowLayout calculates the geometry for a specific row configuration.
// It *must* respect availableHeight and check minimum heights, returning the best effort layout.
func (e *ContinuousLayoutEngine) calculateRowLayout(picsInRow []Picture, rowConfigType string, availableHeight float64) (TemplateLayout, error) {
	fmt.Printf("Debug (calculateRowLayout): Calculating for type '%s', %d pics, availableHeight %.2f\n", rowConfigType, len(picsInRow), availableHeight)

	switch rowConfigType {
	case "row-of-1-ultra", "row-of-1":
		if len(picsInRow) != 1 {
			return TemplateLayout{}, fmt.Errorf("calculateRowLayout: expected 1 picture for row-of-1, got %d", len(picsInRow))
		}
		pic := picsInRow[0]

		// --- Calculate Single Picture Layout (Adapted from processSinglePictureLayoutAndPlace) ---
		aspectRatio := 1.0
		validAR := false
		if pic.Height > 0 && pic.Width > 0 {
			aspectRatio = float64(pic.Width) / float64(pic.Height)
			validAR = true
		} else {
			fmt.Printf("Warning (calculateRowLayout): Invalid dimensions for picture index %d. Using default AR=1.\n", pic.Index)
		}

		picType := GetPictureType(aspectRatio)
		finalWidth := 0.0
		finalHeight := 0.0

		// Use a small positive value for availableHeight if it's near zero to avoid division issues
		if availableHeight <= 1e-6 {
			fmt.Printf("Warning (calculateRowLayout): Available height is near zero (%.2f). Cannot calculate layout.\n", availableHeight)
			// Return zero-height layout, the caller should handle this.
			return TemplateLayout{TotalHeight: 0}, errors.New("available height is too small")
		}

		if !validAR {
			// Fallback for invalid AR: Use available width and a reasonable capped height
			fmt.Printf("Warning (calculateRowLayout): Using fallback dimensions for Pic %d due to invalid AR.\n", pic.Index)
			finalWidth = e.availableWidth
			// Estimate height based on min landscape, capped by available
			finalHeight = math.Min(GetRequiredMinHeight(e, "landscape", 1), availableHeight)
			if finalHeight < 1.0 {
				finalHeight = 1.0
			}
		} else {
			// Calculate dimensions based on picture type and available space
			switch picType {
			case "wide", "landscape":
				// Fit width first
				finalWidth = e.availableWidth
				finalHeight = finalWidth / aspectRatio
				// Scale down if height exceeds available
				if finalHeight > availableHeight {
					finalHeight = availableHeight
					finalWidth = finalHeight * aspectRatio
				}
			case "tall", "portrait", "square", "unknown": // Treat square/unknown like portrait/tall
				// Fit height first
				finalHeight = availableHeight
				finalWidth = finalHeight * aspectRatio
				// Scale down if width exceeds available
				if finalWidth > e.availableWidth {
					finalWidth = e.availableWidth
					finalHeight = finalWidth / aspectRatio
				}
			}
		}

		// Ensure positive dimensions after calculations
		if finalWidth < 1.0 {
			finalWidth = 1.0
		}
		if finalHeight < 1.0 {
			finalHeight = 1.0
		}

		// --- Check Minimum Height (Log warning, but don't error out for dynamic layout) ---
		requiredMinHeight := GetRequiredMinHeight(e, picType, 1) // Check against single pic requirement
		if finalHeight < requiredMinHeight {
			fmt.Printf("Warning (calculateRowLayout): Single picture (Index %d, Type %s) layout height %.2f does not meet minimum %.2f.\n", pic.Index, picType, finalHeight, requiredMinHeight)
		}

		// Return the layout
		return TemplateLayout{
			Positions:   [][]float64{{0, 0}}, // Position relative to row start
			Dimensions:  [][]float64{{finalWidth, finalHeight}},
			TotalHeight: finalHeight,
			TotalWidth:  finalWidth, // Actual width used by this picture
		}, nil

	case "row-of-2":
		if len(picsInRow) != 2 {
			return TemplateLayout{}, fmt.Errorf("calculateRowLayout: expected 2 pictures for row-of-2, got %d", len(picsInRow))
		}
		// Calculate initial uniform layout
		initialWidths, initialHeight, err := e.calculateUniformRowLayout(picsInRow)
		if err != nil {
			return TemplateLayout{}, fmt.Errorf("failed to calculate initial uniform layout for row-of-2: %w", err)
		}

		// Scale if necessary
		finalWidths := initialWidths
		finalHeight := initialHeight
		scale := 1.0
		if finalHeight > availableHeight {
			if finalHeight > 1e-6 {
				scale = availableHeight / finalHeight
				finalHeight *= scale
				for i := range finalWidths {
					finalWidths[i] *= scale
					if finalWidths[i] < 1.0 {
						finalWidths[i] = 1.0
					}
				}
			} else {
				return TemplateLayout{TotalHeight: 0}, errors.New("initial calculated height near zero for row-of-2")
			}
		}

		// Check minimum heights (log warnings)
		for _, pic := range picsInRow {
			ar := 1.0
			if pic.Height > 0 {
				ar = float64(pic.Width) / float64(pic.Height)
			}
			picType := GetPictureType(ar)
			requiredMinHeight := GetRequiredMinHeight(e, picType, 2) // Check against 2-pic requirement
			if finalHeight < requiredMinHeight {
				fmt.Printf("Warning (calculateRowLayout): Row-of-2 picture (Index %d, Type %s) layout height %.2f does not meet minimum %.2f.\n", pic.Index, picType, finalHeight, requiredMinHeight)
			}
		}

		// Construct layout
		positions := [][]float64{{0, 0}, {finalWidths[0] + e.imageSpacing, 0}}
		dimensions := [][]float64{{finalWidths[0], finalHeight}, {finalWidths[1], finalHeight}}
		totalWidth := finalWidths[0] + e.imageSpacing + finalWidths[1]

		return TemplateLayout{
			Positions:   positions,
			Dimensions:  dimensions,
			TotalHeight: finalHeight,
			TotalWidth:  totalWidth, // More accurate width
		}, nil

	case "row-of-3":
		if len(picsInRow) != 3 {
			return TemplateLayout{}, fmt.Errorf("calculateRowLayout: expected 3 pictures for row-of-3, got %d", len(picsInRow))
		}
		// Calculate initial uniform layout
		initialWidths, initialHeight, err := e.calculateUniformRowLayout(picsInRow)
		if err != nil {
			return TemplateLayout{}, fmt.Errorf("failed to calculate initial uniform layout for row-of-3: %w", err)
		}

		// Scale if necessary
		finalWidths := initialWidths
		finalHeight := initialHeight
		scale := 1.0
		if finalHeight > availableHeight {
			if finalHeight > 1e-6 {
				scale = availableHeight / finalHeight
				finalHeight *= scale
				for i := range finalWidths {
					finalWidths[i] *= scale
					if finalWidths[i] < 1.0 {
						finalWidths[i] = 1.0
					}
				}
			} else {
				return TemplateLayout{TotalHeight: 0}, errors.New("initial calculated height near zero for row-of-3")
			}
		}

		// Check minimum heights (log warnings)
		for _, pic := range picsInRow {
			ar := 1.0
			if pic.Height > 0 {
				ar = float64(pic.Width) / float64(pic.Height)
			}
			picType := GetPictureType(ar)
			requiredMinHeight := GetRequiredMinHeight(e, picType, 3) // Check against 3-pic requirement
			if finalHeight < requiredMinHeight {
				fmt.Printf("Warning (calculateRowLayout): Row-of-3 picture (Index %d, Type %s) layout height %.2f does not meet minimum %.2f.\n", pic.Index, picType, finalHeight, requiredMinHeight)
			}
		}

		// Construct layout
		pos1X := finalWidths[0] + e.imageSpacing
		pos2X := pos1X + finalWidths[1] + e.imageSpacing
		positions := [][]float64{{0, 0}, {pos1X, 0}, {pos2X, 0}}
		dimensions := [][]float64{{finalWidths[0], finalHeight}, {finalWidths[1], finalHeight}, {finalWidths[2], finalHeight}}
		totalWidth := finalWidths[0] + e.imageSpacing + finalWidths[1] + e.imageSpacing + finalWidths[2]

		return TemplateLayout{
			Positions:   positions,
			Dimensions:  dimensions,
			TotalHeight: finalHeight,
			TotalWidth:  totalWidth, // More accurate width
		}, nil
	}
	return TemplateLayout{}, fmt.Errorf("calculateRowLayout not implemented for type: %s (%d pics)", rowConfigType, len(picsInRow))
}

// placePicturesInRow places the pictures according to the row's calculated layout.
func (e *ContinuousLayoutEngine) placePicturesInRow(picsInRow []Picture, rowLayout TemplateLayout) {
	// TODO: Implement placement logic (adapt placePicturesInTemplate)
	fmt.Printf("Debug (placePicturesInRow): Placing %d pictures. Layout TotalHeight: %.2f\n", len(picsInRow), rowLayout.TotalHeight)

	if len(picsInRow) != len(rowLayout.Positions) || len(picsInRow) != len(rowLayout.Dimensions) {
		fmt.Println("Error: Mismatch between picture count and layout information in placePicturesInRow.")
		return
	}

	// Ensure entry exists
	if len(e.currentPage.Entries) == 0 {
		fmt.Println("Warning: placePicturesInRow called with no current entry. Creating one.")
		e.currentPage.Entries = append(e.currentPage.Entries, PageEntry{})
	}
	currentEntry := &e.currentPage.Entries[len(e.currentPage.Entries)-1]
	startY := e.currentY // Top Y coordinate for this row

	// Calculate Centering Offset for the row
	actualRowWidth := 0.0 // Calculate actual width from dimensions/positions if not provided reliably
	lastPicIndex := len(rowLayout.Positions) - 1
	if lastPicIndex >= 0 && len(rowLayout.Positions[lastPicIndex]) == 2 && len(rowLayout.Dimensions[lastPicIndex]) == 2 {
		actualRowWidth = rowLayout.Positions[lastPicIndex][0] + rowLayout.Dimensions[lastPicIndex][0]
	} else if rowLayout.TotalWidth > 0 {
		actualRowWidth = rowLayout.TotalWidth // Use if available
	} else {
		// Estimate width if not properly calculated (fallback)
		for _, dim := range rowLayout.Dimensions {
			if len(dim) == 2 {
				actualRowWidth += dim[0]
			}
		}
		actualRowWidth += float64(len(picsInRow)-1) * e.imageSpacing
	}

	offsetX := 0.0
	if actualRowWidth < e.availableWidth {
		offsetX = (e.availableWidth - actualRowWidth) / 2.0
	}
	if offsetX < 0 {
		offsetX = 0
	}

	for i, pic := range picsInRow {
		relativeX := rowLayout.Positions[i][0]
		relativeY := rowLayout.Positions[i][1] // Relative Y within the row (should usually be 0)
		width := rowLayout.Dimensions[i][0]
		height := rowLayout.Dimensions[i][1]

		// Apply centering offset and current Y
		absX0 := e.marginLeft + offsetX + relativeX
		absY0 := startY + relativeY
		absX1 := absX0 + width
		absY1 := absY0 + height

		area := [][]float64{{absX0, absY0}, {absX1, absY1}}

		// Ensure Pictures slice is initialized
		if currentEntry.Pictures == nil {
			currentEntry.Pictures = make([]Picture, 0, len(picsInRow))
		}
		currentEntry.Pictures = append(currentEntry.Pictures, Picture{
			Index:  pic.Index,
			Area:   area,
			URL:    pic.URL,
			Width:  int(math.Round(width)),
			Height: int(math.Round(height)),
		})
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

// calculateUniformRowLayout calculates the dimensions for a row of pictures aiming for a uniform height,
// such that the total width exactly fills the availableWidth.
// It returns the calculated widths for each picture and the calculated uniform height.
func (e *ContinuousLayoutEngine) calculateUniformRowLayout(picsInRow []Picture) (widths []float64, uniformHeight float64, err error) {
	numPics := len(picsInRow)
	if numPics == 0 {
		return nil, 0, errors.New("calculateUniformRowLayout: called with zero pictures")
	}

	AW := e.availableWidth
	spacing := e.imageSpacing
	totalSpacing := float64(numPics-1) * spacing

	// Calculate sum of aspect ratios (W/H)
	sumAR := 0.0
	ARs := make([]float64, numPics)
	for i, pic := range picsInRow {
		if pic.Height <= 0 || pic.Width <= 0 {
			// Handle invalid dimensions - default to AR=1?
			fmt.Printf("Warning (calculateUniformRowLayout): Invalid dimensions for pic index %d. Using AR=1.\n", pic.Index)
			ARs[i] = 1.0
		} else {
			ARs[i] = float64(pic.Width) / float64(pic.Height)
		}
		sumAR += ARs[i]
	}

	if sumAR <= 1e-6 { // Avoid division by zero if all ARs are tiny/invalid
		return nil, 0, errors.New("calculateUniformRowLayout: sum of aspect ratios is zero or negative")
	}

	// Calculate the uniform height H = (AvailableWidth - TotalSpacing) / SumOfAspectRatios
	availableWidthForPics := AW - totalSpacing
	if availableWidthForPics <= 0 {
		// Not enough space even without pictures
		return nil, 0, errors.New("calculateUniformRowLayout: available width is less than or equal to total spacing")
	}
	uniformHeight = availableWidthForPics / sumAR

	if uniformHeight <= 1e-6 {
		return nil, 0, errors.New("calculateUniformRowLayout: calculated uniform height is zero or negative")
	}

	// Calculate individual widths W_i = AR_i * H
	widths = make([]float64, numPics)
	calculatedTotalWidth := 0.0
	for i := 0; i < numPics; i++ {
		widths[i] = ARs[i] * uniformHeight
		calculatedTotalWidth += widths[i]
	}
	calculatedTotalWidth += totalSpacing

	// Optional: Adjust widths slightly due to potential floating point inaccuracies to ensure sum matches AW
	widthAdjustment := (AW - calculatedTotalWidth) / float64(numPics)
	for i := range widths {
		widths[i] += widthAdjustment
		if widths[i] < 1.0 {
			widths[i] = 1.0
		} // Ensure minimum width
	}

	return widths, uniformHeight, nil
}

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
		// --- UPDATED: Call new function with split logic ---
		fmt.Println("Debug (TemplateDispatch): Calling processThreePicturesWithSplitLogic for 3 pictures.")
		return e.processThreePicturesWithSplitLogic(pictures, layoutAvailableHeight)
	case 4:
		// --- UPDATED: Call new function with split logic ---
		fmt.Println("Debug (TemplateDispatch): Calling processFourPicturesWithSplitLogic for 4 pictures.")
		return e.processFourPicturesWithSplitLogic(pictures, layoutAvailableHeight)
	case 5:
		// --- UPDATED: Call new function with split logic ---
		fmt.Println("Debug (TemplateDispatch): Calling processFivePicturesWithSplitLogic for 5 pictures.")
		return e.processFivePicturesWithSplitLogic(pictures, layoutAvailableHeight)
	case 6:
		// --- UPDATED: Call new function with split logic ---
		fmt.Println("Debug (TemplateDispatch): Calling processSixPicturesWithSplitLogic for 6 pictures.")
		return e.processSixPicturesWithSplitLogic(pictures, layoutAvailableHeight)
	case 7:
		// --- CORRECTED: Call the function with split logic ---
		fmt.Println("Debug (TemplateDispatch): Calling processSevenPicturesWithSplitLogic for 7 pictures.")
		return e.processSevenPicturesWithSplitLogic(pictures, layoutAvailableHeight)
	case 8:
		// --- UPDATED: Call new function with split logic ---
		fmt.Println("Debug (TemplateDispatch): Calling processEightPicturesWithSplitLogic for 8 pictures.")
		return e.processEightPicturesWithSplitLogic(pictures, layoutAvailableHeight)
	case 9:
		// --- UPDATED: Call new function with split logic ---
		fmt.Println("Debug (TemplateDispatch): Calling processNinePicturesWithSplitLogic for 9 pictures.")
		return e.processNinePicturesWithSplitLogic(pictures, layoutAvailableHeight)
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

// processTwoPicGroup attempts to calculate and place a group of 2 pictures.
// It handles pagination if the group doesn't fit the initial availableHeight.
// It updates the engine's currentY state directly upon successful placement.
// Returns the height used by the group and any critical error.
func (e *ContinuousLayoutEngine) processTwoPicGroup(
	groupPics []Picture,
	groupNum int, // For logging (e.g., 1st, 2nd, 3rd group of 2)
	layoutAvailableHeight float64,
) (heightUsed float64, err error) {

	if len(groupPics) != 2 {
		return 0, fmt.Errorf("processTwoPicGroup: expected 2 pictures, got %d", len(groupPics))
	}
	var layoutInfo TemplateLayout
	var calcErr error
	const MINIMUM_VIABLE_ROW_HEIGHT = 50.0 // Minimum pixels high for a row to be considered valid

	fmt.Printf("Debug (Split-Group 2pic-%d): Attempting calculation. AvailableHeight: %.2f\n", groupNum, layoutAvailableHeight)

	// Use calculateRowLayout for 2 pictures
	layoutInfo, calcErr = e.calculateRowLayout(groupPics, "row-of-2", layoutAvailableHeight)

	// --- Check Fit and Place or Retry ---
	// ADDED: Check if calculated height is too small to be viable
	isTooSmall := layoutInfo.TotalHeight < MINIMUM_VIABLE_ROW_HEIGHT
	if calcErr == nil && !isTooSmall && layoutInfo.TotalHeight <= layoutAvailableHeight+1e-6 { // Add tolerance
		// Fits on current page and is viable
		fmt.Printf("Debug (Split-Group 2pic-%d): Group fits and is viable on current page (Page %d, height: %.2f <= available: %.2f)\n", groupNum, e.currentPage.Page, layoutInfo.TotalHeight, layoutAvailableHeight)
		e.placePicturesInRow(groupPics, layoutInfo)
		e.currentY += layoutInfo.TotalHeight // Update Y
		return layoutInfo.TotalHeight, nil
	} else {
		// Doesn't fit, calculation error, or too small - Force New Page
		reason := "Unknown reason"
		if calcErr != nil {
			reason = fmt.Sprintf("Initial calc failed: %v", calcErr)
		} else if isTooSmall {
			reason = fmt.Sprintf("Calculated height %.2f is less than minimum viable %.2f", layoutInfo.TotalHeight, MINIMUM_VIABLE_ROW_HEIGHT)
		} else { // Must be layoutInfo.TotalHeight > layoutAvailableHeight
			reason = fmt.Sprintf("Group doesn't fit (%.2f > %.2f)", layoutInfo.TotalHeight, layoutAvailableHeight)
		}
		fmt.Printf("Debug (Split-Group 2pic-%d): %s. Forcing new page.\n", groupNum, reason)

		e.newPage()
		newAvailableHeight := (e.marginTop + e.availableHeight) - e.currentY

		// Retry calculation on the new page
		fmt.Printf("Debug (Split-Group 2pic-%d): Retrying calculation & placement on new page (Page %d, available: %.2f)\n", groupNum, e.currentPage.Page, newAvailableHeight)
		layoutInfo, calcErr = e.calculateRowLayout(groupPics, "row-of-2", newAvailableHeight)

		// ADDED: Check if calculated height is too small on retry
		isTooSmallRetry := layoutInfo.TotalHeight < MINIMUM_VIABLE_ROW_HEIGHT
		if calcErr != nil || isTooSmallRetry || layoutInfo.TotalHeight > newAvailableHeight+1e-6 { // Add tolerance
			// Handle error: Failed even on new page
			errMsg := fmt.Sprintf("failed to place 2-pic group %d even on new page", groupNum)
			if calcErr != nil {
				errMsg += fmt.Sprintf(": %v", calcErr)
			} else if isTooSmallRetry {
				errMsg += fmt.Sprintf(" (height %.2f < min viable %.2f)", layoutInfo.TotalHeight, MINIMUM_VIABLE_ROW_HEIGHT)
			} else {
				errMsg += fmt.Sprintf(" (height %.2f > available %.2f)", layoutInfo.TotalHeight, newAvailableHeight)
			}
			fmt.Printf("Error: %s\n", errMsg)
			return 0, errors.New(errMsg)
		}
		// Place successfully on new page
		e.placePicturesInRow(groupPics, layoutInfo)
		e.currentY += layoutInfo.TotalHeight // Update Y on new page
		return layoutInfo.TotalHeight, nil
	}
}
