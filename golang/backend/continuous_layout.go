package backend

import (
	"math"
	"strings"
	"unicode/utf8"
)

// ContinuousLayoutEngine represents the continuous layout engine
type ContinuousLayoutEngine struct {
	entries           []Entry
	pages             []Page
	currentPage       *Page
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
	singleImageHeight float64 // 单张图片的默认高度
}

// NewContinuousLayoutEngine creates a new continuous layout engine
func NewContinuousLayoutEngine(entries []Entry) *ContinuousLayoutEngine {
	engine := &ContinuousLayoutEngine{
		entries:           entries,
		marginLeft:        100,
		marginRight:       100,
		marginTop:         100,
		marginBottom:      160,
		timeHeight:        104,
		fontSize:          50,
		lineHeight:        75,
		entrySpacing:      150,  // 条目之间的间距
		elementSpacing:    30,   // 元素整体之间的间距
		imageSpacing:      20,   // 图片之间的间距
		singleImageHeight: 1695, // 设置单张图片的默认高度
	}
	engine.availableWidth = 2480 - engine.marginLeft - engine.marginRight
	engine.availableHeight = 3508 - engine.marginTop - engine.marginBottom
	return engine
}

// ProcessEntries processes all entries and returns the layout result
func (e *ContinuousLayoutEngine) ProcessEntries() ([]Page, error) {
	e.newPage()

	for i, entry := range e.entries {
		if i > 0 {
			// 如果不是第一个条目，添加间距
			e.currentY += e.entrySpacing
		}
		e.processEntry(entry)
	}

	return e.pages, nil
}

func (e *ContinuousLayoutEngine) processEntry(entry Entry) {
	// 1.1) 检查时间区域是否可以在当前页面展示
	availableHeight := e.availableHeight - (e.currentY - e.marginTop)
	minTimeHeight := e.timeHeight
	if e.currentY > e.marginTop {
		minTimeHeight += e.entrySpacing
	}

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
		minPictureHeight = e.calculateFirstRowPictureHeight(entry.Pictures) + e.elementSpacing
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
		// 单张图片使用默认高度
		return e.singleImageHeight
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
			availableHeight = e.availableHeight
		}

		availableLines := int(math.Floor(availableHeight / e.lineHeight))
		if availableLines <= 0 {
			e.newPage()
			continue
		}

		availableLines = int(math.Min(float64(len(lines)-currentLine), float64(availableLines)))
		chunk := lines[currentLine : currentLine+availableLines]
		e.addTextChunk(chunk)
		currentLine += availableLines

		if currentLine < len(lines) {
			e.newPage()
		}
	}
}

func (e *ContinuousLayoutEngine) newPage() {
	page := &Page{
		Page:      len(e.pages) + 1,
		TextAreas: make([][][]float64, 0),
		Texts:     make([]string, 0),
		Pictures:  make([]Picture, 0),
	}
	e.pages = append(e.pages, *page)
	e.currentPage = &e.pages[len(e.pages)-1]
	e.currentY = e.marginTop
	e.timeAreaBottom = 0
}

func (e *ContinuousLayoutEngine) addTime(timeStr string) {
	// 如果不是页面开始，添加条目间距
	if e.currentY > e.marginTop {
		e.currentY += e.entrySpacing // 使用条目间距
	}

	x0 := e.marginLeft
	y0 := e.currentY
	x1 := x0 + e.availableWidth
	y1 := y0 + e.timeHeight

	e.currentPage.TimeArea = [][]float64{{x0, y0}, {x1, y1}}
	e.currentPage.Time = timeStr

	// 更新当前位置和时间区域底部位置
	e.currentY = y1 + e.elementSpacing // 使用元素间距
	e.timeAreaBottom = e.currentY - e.marginTop
}

func (e *ContinuousLayoutEngine) addTextChunk(chunk []string) {
	startY := e.currentY
	textHeight := float64(len(chunk)) * e.lineHeight
	area := [][]float64{
		{e.marginLeft, startY},
		{e.marginLeft + e.availableWidth, startY + textHeight},
	}
	e.currentPage.TextAreas = append(e.currentPage.TextAreas, area)
	e.currentPage.Texts = append(e.currentPage.Texts, strings.Join(chunk, "\n"))
	e.currentY += textHeight + e.elementSpacing // 使用元素间距
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
	// 检查当前页面剩余空间
	availableHeight := e.availableHeight - (e.currentY - e.marginTop)

	// 如果剩余空间小于默认高度，创建新页面
	if availableHeight < e.singleImageHeight {
		e.newPage()
		availableHeight = e.availableHeight
	}

	// 计算图片的实际高度，使用可用高度和默认高度中的较小值
	actualHeight := math.Min(availableHeight, e.singleImageHeight)

	// 根据原始图片的宽高比计算对应的宽度
	aspectRatio := float64(pic.Width) / float64(pic.Height)
	width := actualHeight * aspectRatio

	// 确保宽度不超过可用宽度
	if width > e.availableWidth {
		width = e.availableWidth
		actualHeight = width / aspectRatio
	}

	// 计算水平居中的起始x坐标
	x := e.marginLeft + (e.availableWidth-width)/2

	// 创建图片区域
	area := [][]float64{
		{x, e.currentY},
		{x + width, e.currentY + actualHeight},
	}

	e.currentPage.Pictures = append(e.currentPage.Pictures, Picture{
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
	}

	// 布局这一行的图片
	startY := e.currentY
	x := e.marginLeft
	for _, pic := range rowPics {
		area := [][]float64{
			{x, startY},
			{x + width, startY + height},
		}
		e.currentPage.Pictures = append(e.currentPage.Pictures, Picture{
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
