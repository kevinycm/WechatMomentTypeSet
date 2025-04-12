package backend

import (
	"math"
	"strings"
	"unicode/utf8"
)

// ContinuousLayoutPage represents a single page in the continuous layout
type ContinuousLayoutPage struct {
	Page      int         `json:"page"`
	IsInsert  bool        `json:"is_insert"`  // 是否是插页
	YearMonth string      `json:"year_month"` // 年月信息，格式：2025年3月
	Entries   []PageEntry `json:"entries"`
}

// PageEntry represents a single entry's layout information on a page
type PageEntry struct {
	Time      string        `json:"time"` // 格式：2025年3月30日 17:50
	TimeArea  [][]float64   `json:"time_area"`
	TextAreas [][][]float64 `json:"text_areas"`
	Texts     []string      `json:"texts"`
	Pictures  []Picture     `json:"pictures"`
}

// ContinuousLayoutEngine represents the continuous layout engine
type ContinuousLayoutEngine struct {
	entries           []Entry
	pages             []ContinuousLayoutPage
	currentPage       *ContinuousLayoutPage
	marginLeft        float64
	marginRight       float64
	marginTop         float64
	marginBottom      float64
	availableWidth    float64
	availableHeight   float64
	timeHeight        float64
	fontSize          float64
	lineHeight        float64
	currentY          float64
	timeAreaBottom    float64
	entrySpacing      float64 // 条目之间的间距
	elementSpacing    float64 // 元素之间的间距
	imageSpacing      float64 // 图片之间的间距
	singleImageHeight float64 // 单张竖图的默认高度
	singleImageWidth  float64 // 单张横图的默认宽度
	minImageHeight    float64 // 单张竖图的最小展示高度
	minImageWidth     float64 // 单张横图的最小展示宽度
}

// NewContinuousLayoutEngine creates a new continuous layout engine
func NewContinuousLayoutEngine(entries []Entry) *ContinuousLayoutEngine {
	engine := &ContinuousLayoutEngine{
		entries:           entries,
		marginLeft:        160,
		marginRight:       160,
		marginTop:         160,
		marginBottom:      160,
		timeHeight:        104,
		fontSize:          50,
		lineHeight:        75,
		entrySpacing:      150,  // 条目之间的间距
		elementSpacing:    30,   // 元素整体之间的间距
		imageSpacing:      10,   // 图片之间的间距
		singleImageHeight: 2260, // 设置单张竖图的默认高度
		singleImageWidth:  1695, // 设置单张横图的默认宽度
		minImageHeight:    800,  // 设置单张竖图的最小展示高度
		minImageWidth:     1200, // 设置单张横图的最小展示宽度
	}
	engine.availableWidth = 2480 - engine.marginLeft - engine.marginRight
	engine.availableHeight = 3508 - engine.marginTop - engine.marginBottom
	return engine
}

// ProcessEntries processes all entries and returns the layout result
func (e *ContinuousLayoutEngine) ProcessEntries() ([]ContinuousLayoutPage, error) {
	e.newPage()

	for i, entry := range e.entries {
		// 计算当前条目的总高度
		totalHeight := e.calculateEntryTotalHeight(entry)

		// 检查当前页面剩余空间
		availableHeight := e.availableHeight - (e.currentY - e.marginTop)

		// 如果是第一个条目或剩余空间不足以完整放下当前条目，创建新页面
		if i == 0 || availableHeight < totalHeight {
			if i > 0 {
				e.newPage()
			}
			e.processEntry(entry)
		} else {
			// 有足够空间，直接在当前页面处理
			e.processEntry(entry)
		}
	}

	return e.pages, nil
}

// calculateEntryTotalHeight 计算一个条目的总高度
func (e *ContinuousLayoutEngine) calculateEntryTotalHeight(entry Entry) float64 {
	var totalHeight float64

	// 0. 如果不是页面开始位置，需要加上条目间距
	if e.currentY > e.marginTop {
		totalHeight += e.entrySpacing
	}

	// 1. 计算时间区域高度
	totalHeight += e.timeHeight + e.elementSpacing

	// 2. 计算文本区域高度
	if strings.TrimSpace(entry.Text) != "" {
		charsPerLine := int(e.availableWidth / e.fontSize)
		var lines []string
		for _, paragraph := range strings.Split(entry.Text, "\n") {
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
					lines = append(lines, string(runes[i:end]))
				}
			}
		}
		totalHeight += float64(len(lines))*e.lineHeight + e.elementSpacing
	}

	// 3. 计算图片区域高度
	if len(entry.Pictures) > 0 {
		if len(entry.Pictures) == 1 {
			// 单张图片
			pic := entry.Pictures[0]
			aspectRatio := float64(pic.Width) / float64(pic.Height)
			if aspectRatio > 1 {
				// 横图
				width := math.Min(e.singleImageWidth, e.availableWidth)
				if width >= e.minImageWidth {
					totalHeight += width/aspectRatio + e.elementSpacing
				}
			} else {
				// 竖图
				totalHeight += math.Min(e.singleImageHeight, e.availableHeight) + e.elementSpacing
			}
		} else {
			// 多张图片
			layout := e.getLayout(len(entry.Pictures))
			currentIndex := 0
			for _, rowCount := range layout {
				// 计算这一行的高度
				maxRowHeight := 0.0
				for i := 0; i < rowCount && currentIndex+i < len(entry.Pictures); i++ {
					pic := entry.Pictures[currentIndex+i]
					width := (e.availableWidth - e.imageSpacing*float64(rowCount-1)) / float64(rowCount)
					height := width / (float64(pic.Width) / float64(pic.Height))
					if height > maxRowHeight {
						maxRowHeight = height
					}
				}
				totalHeight += maxRowHeight + e.imageSpacing
				currentIndex += rowCount
			}
		}
	}

	return totalHeight
}

func (e *ContinuousLayoutEngine) processEntry(entry Entry) {
	// 如果不是页面开始位置，添加条目间距
	if e.currentY > e.marginTop {
		e.currentY += e.entrySpacing
	}

	// 1.1) 检查时间区域是否可以在当前页面展示
	availableHeight := e.availableHeight - (e.currentY - e.marginTop)
	minTimeHeight := e.timeHeight

	// 如果时间区域无法展示，创建新页面
	if availableHeight < minTimeHeight {
		e.newPage()
		availableHeight = e.availableHeight
	}

	// 1.2) 时间区域可以展示，需要判断文本或图片是否可以和时间一起展示
	hasText := strings.TrimSpace(entry.Text) != ""
	hasPictures := len(entry.Pictures) > 0

	// 计算至少一行文本需要的高度
	minTextHeight := e.lineHeight + e.elementSpacing
	// 计算至少一行图片需要的高度
	var minPictureHeight float64
	if hasPictures {
		if len(entry.Pictures) == 1 {
			// 单张图片的情况
			pic := entry.Pictures[0]
			aspectRatio := float64(pic.Width) / float64(pic.Height)
			if aspectRatio > 1 {
				// 横图
				width := math.Min(e.singleImageWidth, e.availableWidth)
				if width >= e.minImageWidth {
					height := width/aspectRatio + e.elementSpacing
					// 检查是否有足够的可展示空间（不包含底部边距）
					if height <= availableHeight {
						minPictureHeight = height
					} else {
						minPictureHeight = e.availableHeight + 1
					}
				} else {
					minPictureHeight = e.availableHeight + 1
				}
			} else {
				// 竖图
				if availableHeight >= e.minImageHeight {
					minPictureHeight = math.Min(e.singleImageHeight, availableHeight) + e.elementSpacing
				} else {
					minPictureHeight = e.availableHeight + 1
				}
			}
		} else {
			minPictureHeight = e.calculateFirstRowPictureHeight(entry.Pictures) + e.elementSpacing
		}
	}

	// 检查是否可以和时间一起展示
	canShowWithText := hasText && (availableHeight >= minTimeHeight+minTextHeight)
	canShowWithPictures := hasPictures && (availableHeight >= minTimeHeight+minPictureHeight)

	// 1.2.1) 或 1.2.2) 如果可以和时间一起展示，则在当前页面展示
	if canShowWithText || canShowWithPictures {
		e.addTime(entry.Time)

		if hasText {
			e.processText(entry.Text)
		}

		if hasPictures {
			e.processPictures(entry.Pictures)
		}
	} else {
		// 1.2.3) 或 1.2.4) 如果不能一起展示，创建新页面
		e.newPage()
		e.addTime(entry.Time)

		if hasText {
			e.processText(entry.Text)
		}

		if hasPictures {
			e.processPictures(entry.Pictures)
		}
	}
}

func (e *ContinuousLayoutEngine) calculateFirstRowPictureHeight(pictures []Picture) float64 {
	if len(pictures) == 0 {
		return 0
	}

	// 获取第一行图片数量
	layout := e.getLayout(len(pictures))
	firstRowCount := layout[0]

	if firstRowCount == 1 {
		// 单张图片处理
		pic := pictures[0]
		aspectRatio := float64(pic.Width) / float64(pic.Height)
		availableHeight := e.availableHeight - (e.currentY - e.marginTop)

		if aspectRatio > 1 {
			// 横图：使用默认宽度计算高度
			width := math.Min(e.singleImageWidth, e.availableWidth)
			if width < e.minImageWidth {
				return e.availableHeight + 1
			}
			height := width / aspectRatio
			// 检查是否有足够的可展示空间（不包含底部边距）
			if height > availableHeight {
				return e.availableHeight + 1
			}
			return height
		} else {
			// 竖图：使用默认高度
			if availableHeight < e.minImageHeight {
				return e.availableHeight + 1
			}
			return math.Min(e.singleImageHeight, availableHeight)
		}
	}

	// 多张图片计算第一行需要的高度
	totalWidth := e.availableWidth
	spacing := e.imageSpacing * float64(firstRowCount-1)
	availableWidth := totalWidth - spacing
	width := availableWidth / float64(firstRowCount)

	// 使用第一张图片的宽高比来计算高度
	aspectRatio := float64(pictures[0].Width) / float64(pictures[0].Height)
	return width / aspectRatio
}

func (e *ContinuousLayoutEngine) processText(text string) {
	if strings.TrimSpace(text) == "" {
		return
	}

	// 如果不是紧跟在时间区域后面，添加元素间距
	if e.currentY > e.marginTop+e.timeHeight {
		e.currentY += e.elementSpacing
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
		availableHeight := e.availableHeight - (e.currentY - e.marginTop)

		if availableHeight <= 0 {
			e.newPage()
			// 创建新的条目
			entry := PageEntry{
				TextAreas: make([][][]float64, 0),
				Texts:     make([]string, 0),
				Pictures:  make([]Picture, 0),
			}
			e.currentPage.Entries = append(e.currentPage.Entries, entry)
			availableHeight = e.availableHeight
		}

		availableLines := int(math.Floor(availableHeight / e.lineHeight))
		if availableLines <= 0 {
			e.newPage()
			// 创建新的条目
			entry := PageEntry{
				TextAreas: make([][][]float64, 0),
				Texts:     make([]string, 0),
				Pictures:  make([]Picture, 0),
			}
			e.currentPage.Entries = append(e.currentPage.Entries, entry)
			continue
		}

		availableLines = int(math.Min(float64(len(lines)-currentLine), float64(availableLines)))
		chunk := lines[currentLine : currentLine+availableLines]
		e.addTextChunk(chunk)
		currentLine += availableLines

		if currentLine < len(lines) {
			e.newPage()
			// 创建新的条目
			entry := PageEntry{
				TextAreas: make([][][]float64, 0),
				Texts:     make([]string, 0),
				Pictures:  make([]Picture, 0),
			}
			e.currentPage.Entries = append(e.currentPage.Entries, entry)
		}
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

func (e *ContinuousLayoutEngine) addTime(timeStr string) {
	x0 := e.marginLeft
	y0 := e.currentY
	x1 := x0 + e.availableWidth
	y1 := y0 + e.timeHeight

	// 创建新的 PageEntry
	entry := PageEntry{
		Time:      timeStr,
		TimeArea:  [][]float64{{x0, y0}, {x1, y1}},
		TextAreas: make([][][]float64, 0),
		Texts:     make([]string, 0),
		Pictures:  make([]Picture, 0),
	}
	e.currentPage.Entries = append(e.currentPage.Entries, entry)

	// 更新当前位置和时间区域底部位置
	e.currentY = y1 + e.elementSpacing
	e.timeAreaBottom = y1 - e.marginTop
}

func (e *ContinuousLayoutEngine) addTextChunk(chunk []string) {
	if len(chunk) == 0 {
		return
	}

	// 检查是否有条目，如果没有则创建一个空条目
	if len(e.currentPage.Entries) == 0 {
		e.newPage()
		entry := PageEntry{
			TextAreas: make([][][]float64, 0),
			Texts:     make([]string, 0),
			Pictures:  make([]Picture, 0),
		}
		e.currentPage.Entries = append(e.currentPage.Entries, entry)
	}

	startY := e.currentY
	textHeight := float64(len(chunk)) * e.lineHeight
	area := [][]float64{
		{e.marginLeft, startY},
		{e.marginLeft + e.availableWidth, startY + textHeight},
	}

	// 获取当前页面的最后一个条目
	currentEntry := &e.currentPage.Entries[len(e.currentPage.Entries)-1]
	currentEntry.TextAreas = append(currentEntry.TextAreas, area)
	currentEntry.Texts = append(currentEntry.Texts, strings.Join(chunk, "\n"))

	e.currentY += textHeight + e.elementSpacing
}

func (e *ContinuousLayoutEngine) processPictures(pictures []Picture) {
	if len(pictures) == 0 {
		return
	}

	// 单张图片的特殊处理
	if len(pictures) == 1 {
		e.processSinglePicture(pictures[0])
		return
	}

	// 多张图片的处理保持不变
	layout := e.getLayout(len(pictures))
	currentIdx := 0
	for _, row := range layout {
		endIdx := currentIdx + row
		if endIdx > len(pictures) {
			endIdx = len(pictures)
		}
		rowPics := pictures[currentIdx:endIdx]
		e.processPictureRow(rowPics)
		currentIdx += row
	}
}

func (e *ContinuousLayoutEngine) processSinglePicture(pic Picture) {
	// 检查是否有条目，如果没有则创建一个空条目
	if len(e.currentPage.Entries) == 0 {
		e.newPage()
		entry := PageEntry{
			TextAreas: make([][][]float64, 0),
			Texts:     make([]string, 0),
			Pictures:  make([]Picture, 0),
		}
		e.currentPage.Entries = append(e.currentPage.Entries, entry)
	}

	// 检查当前页面剩余空间（不包含底部边距）
	availableHeight := e.availableHeight - (e.currentY - e.marginTop)

	// 计算图片的宽高比
	aspectRatio := float64(pic.Width) / float64(pic.Height)

	var actualWidth, actualHeight float64

	if aspectRatio > 1 {
		// 横图：基于默认宽度计算
		actualWidth = math.Min(e.singleImageWidth, e.availableWidth)
		actualHeight = actualWidth / aspectRatio

		// 如果计算出的宽度小于最小展示宽度或高度超过可用空间，创建新页面
		if actualWidth < e.minImageWidth || actualHeight > availableHeight {
			e.newPage()
			// 创建新的条目
			entry := PageEntry{
				TextAreas: make([][][]float64, 0),
				Texts:     make([]string, 0),
				Pictures:  make([]Picture, 0),
			}
			e.currentPage.Entries = append(e.currentPage.Entries, entry)
			actualWidth = math.Min(e.singleImageWidth, e.availableWidth)
			actualHeight = actualWidth / aspectRatio
		}
	} else {
		// 竖图：基于默认高度计算
		actualHeight = math.Min(e.singleImageHeight, availableHeight)
		actualWidth = actualHeight * aspectRatio

		// 如果计算出的高度小于最小展示高度，创建新页面
		if actualHeight < e.minImageHeight {
			e.newPage()
			// 创建新的条目
			entry := PageEntry{
				TextAreas: make([][][]float64, 0),
				Texts:     make([]string, 0),
				Pictures:  make([]Picture, 0),
			}
			e.currentPage.Entries = append(e.currentPage.Entries, entry)
			actualHeight = math.Min(e.singleImageHeight, e.availableHeight)
			actualWidth = actualHeight * aspectRatio
		}
	}

	// 确保宽度不超过可用宽度
	if actualWidth > e.availableWidth {
		actualWidth = e.availableWidth
		actualHeight = actualWidth / aspectRatio
	}

	// 计算水平居中的起始x坐标
	x := e.marginLeft + (e.availableWidth-actualWidth)/2

	// 创建图片区域
	area := [][]float64{
		{x, e.currentY},
		{x + actualWidth, e.currentY + actualHeight},
	}

	// 获取当前页面的最后一个条目
	currentEntry := &e.currentPage.Entries[len(e.currentPage.Entries)-1]
	currentEntry.Pictures = append(currentEntry.Pictures, Picture{
		Index: pic.Index,
		Area:  area,
		URL:   pic.URL,
	})

	// 更新当前Y坐标
	e.currentY += actualHeight + e.elementSpacing
}

func (e *ContinuousLayoutEngine) processPictureRow(rowPics []Picture) {
	if len(rowPics) == 0 {
		return
	}

	// 检查是否有条目，如果没有则创建一个空条目
	if len(e.currentPage.Entries) == 0 {
		e.newPage()
		entry := PageEntry{
			TextAreas: make([][][]float64, 0),
			Texts:     make([]string, 0),
			Pictures:  make([]Picture, 0),
		}
		e.currentPage.Entries = append(e.currentPage.Entries, entry)
	}

	// 计算这一行的理想高度和宽度
	totalWidth := e.availableWidth
	spacing := e.imageSpacing * float64(len(rowPics)-1)
	availableWidth := totalWidth - spacing
	width := availableWidth / float64(len(rowPics))
	aspectRatio := float64(rowPics[0].Width) / float64(rowPics[0].Height)
	height := width / aspectRatio

	// 检查当前页面剩余空间
	availableHeight := e.availableHeight - (e.currentY - e.marginTop)
	if availableHeight < height {
		e.newPage()
		// 创建新的条目
		entry := PageEntry{
			TextAreas: make([][][]float64, 0),
			Texts:     make([]string, 0),
			Pictures:  make([]Picture, 0),
		}
		e.currentPage.Entries = append(e.currentPage.Entries, entry)
	}

	// 布局这一行的图片
	startY := e.currentY
	x := e.marginLeft
	for _, pic := range rowPics {
		area := [][]float64{
			{x, startY},
			{x + width, startY + height},
		}
		e.currentPage.Entries[len(e.currentPage.Entries)-1].Pictures = append(e.currentPage.Entries[len(e.currentPage.Entries)-1].Pictures, Picture{
			Index: pic.Index,
			Area:  area,
			URL:   pic.URL,
		})
		x += width + e.imageSpacing
	}

	// 更新当前Y坐标，如果不是最后一行，添加行间距
	e.currentY = startY + height
	if len(rowPics) > 0 {
		e.currentY += e.imageSpacing
	}
}

func (e *ContinuousLayoutEngine) getLayout(n int) []int {
	layoutRules := map[int][]int{
		1: {1},
		2: {2},
		3: {3},
		4: {2, 2},
		5: {2, 3},
		6: {3, 3},
		7: {2, 2, 3},
		8: {3, 3, 2},
		9: {3, 3, 3},
	}
	return layoutRules[n]
}
