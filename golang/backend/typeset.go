package backend

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
)

// Component represents a basic component in the layout
type Component interface {
	GetID() string
	GetWidth() float64
	GetHeight() float64
	GetX() float64
	GetY() float64
	SetPosition(x, y float64)
	SetSize(width, height float64)
}

// PictureComponent represents a picture in the moment
type PictureComponent struct {
	ID             string  `json:"id"`
	X              float64 `json:"x"`
	Y              float64 `json:"y"`
	Width          float64 `json:"width"`
	Height         float64 `json:"height"`
	URL            string  `json:"url"`
	OriginalWidth  float64 `json:"original_width"`
	OriginalHeight float64 `json:"original_height"`
}

func NewPictureComponent(id string, width, height float64, url string) *PictureComponent {
	return &PictureComponent{
		ID:             id,
		OriginalWidth:  width,
		OriginalHeight: height,
		URL:            url,
	}
}

func (c *PictureComponent) GetID() string {
	return c.ID
}

func (c *PictureComponent) GetWidth() float64 {
	return c.Width
}

func (c *PictureComponent) GetHeight() float64 {
	return c.Height
}

func (c *PictureComponent) GetX() float64 {
	return c.X
}

func (c *PictureComponent) GetY() float64 {
	return c.Y
}

func (c *PictureComponent) SetPosition(x, y float64) {
	c.X = x
	c.Y = y
}

func (c *PictureComponent) SetSize(width, height float64) {
	c.Width = width
	c.Height = height
}

// Calculate aspect ratio preserving size
func (c *PictureComponent) CalculateSize(maxWidth, maxHeight float64) {
	ratio := math.Min(maxWidth/c.OriginalWidth, maxHeight/c.OriginalHeight)
	c.Width = c.OriginalWidth * ratio
	c.Height = c.OriginalHeight * ratio
}

// TextComponent represents a text block in the moment
type TextComponent struct {
	ID         string  `json:"id"`
	X          float64 `json:"x"`
	Y          float64 `json:"y"`
	Width      float64 `json:"width"`
	Height     float64 `json:"height"`
	Text       string  `json:"text"`
	LineHeight float64 `json:"line_height"`
	FontSize   float64 `json:"font_size"`
}

func NewTextComponent(id string, text string) *TextComponent {
	return &TextComponent{
		ID:         id,
		Text:       text,
		Width:      600, // Default width
		LineHeight: 24,  // Default line height
		FontSize:   16,  // Default font size
	}
}

func (c *TextComponent) GetID() string {
	return c.ID
}

func (c *TextComponent) GetWidth() float64 {
	return c.Width
}

func (c *TextComponent) GetHeight() float64 {
	return c.Height
}

func (c *TextComponent) GetX() float64 {
	return c.X
}

func (c *TextComponent) GetY() float64 {
	return c.Y
}

func (c *TextComponent) SetPosition(x, y float64) {
	c.X = x
	c.Y = y
}

func (c *TextComponent) SetSize(width, height float64) {
	c.Width = width
	c.Height = height
}

// Calculate height based on text content
func (c *TextComponent) CalculateHeight() {
	lines := strings.Count(c.Text, "\n") + 1
	c.Height = float64(lines) * c.LineHeight
}

// TimeComponent represents the time display in the moment
type TimeComponent struct {
	ID     string  `json:"id"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Time   string  `json:"time"`
}

func NewTimeComponent(id string, time string) *TimeComponent {
	return &TimeComponent{
		ID:     id,
		Time:   time,
		Width:  200,
		Height: 30,
	}
}

func (c *TimeComponent) GetID() string {
	return c.ID
}

func (c *TimeComponent) GetWidth() float64 {
	return c.Width
}

func (c *TimeComponent) GetHeight() float64 {
	return c.Height
}

func (c *TimeComponent) GetX() float64 {
	return c.X
}

func (c *TimeComponent) GetY() float64 {
	return c.Y
}

func (c *TimeComponent) SetPosition(x, y float64) {
	c.X = x
	c.Y = y
}

func (c *TimeComponent) SetSize(width, height float64) {
	c.Width = width
	c.Height = height
}

// LayoutManager handles the overall layout of the moment
type LayoutManager struct {
	Width        float64             `json:"width"`
	Height       float64             `json:"height"`
	Padding      float64             `json:"padding"`
	Spacing      float64             `json:"spacing"`
	Components   []interface{}       `json:"components"`
	TimeComp     *TimeComponent      `json:"time_component"`
	TextComp     *TextComponent      `json:"text_component"`
	PictureComps []*PictureComponent `json:"picture_components"`
}

func NewLayoutManager(width, height, padding, spacing float64) *LayoutManager {
	return &LayoutManager{
		Width:      width,
		Height:     height,
		Padding:    padding,
		Spacing:    spacing,
		Components: make([]interface{}, 0),
	}
}

func (lm *LayoutManager) AddTimeComponent(time string) {
	lm.TimeComp = NewTimeComponent("time", time)
	lm.Components = append(lm.Components, lm.TimeComp)
}

func (lm *LayoutManager) AddTextComponent(text string) {
	lm.TextComp = NewTextComponent("text", text)
	lm.TextComp.CalculateHeight()
	lm.Components = append(lm.Components, lm.TextComp)
}

func (lm *LayoutManager) AddPictureComponent(width, height float64, url string) {
	pic := NewPictureComponent(fmt.Sprintf("pic%d", len(lm.PictureComps)+1), width, height, url)
	lm.PictureComps = append(lm.PictureComps, pic)
	lm.Components = append(lm.Components, pic)
}

func (lm *LayoutManager) Layout() {
	// Calculate available width and height
	availableWidth := lm.Width - 2*lm.Padding
	availableHeight := lm.Height - 2*lm.Padding

	// Position time component
	if lm.TimeComp != nil {
		lm.TimeComp.SetPosition(lm.Padding, lm.Padding)
		availableHeight -= lm.TimeComp.GetHeight() + lm.Spacing
	}

	// Position text component
	if lm.TextComp != nil {
		textY := lm.Padding
		if lm.TimeComp != nil {
			textY += lm.TimeComp.GetHeight() + lm.Spacing
		}
		lm.TextComp.SetPosition(lm.Padding, textY)
		availableHeight -= lm.TextComp.GetHeight() + lm.Spacing
	}

	// Calculate grid layout for pictures
	if len(lm.PictureComps) > 0 {
		// Determine grid dimensions
		cols := int(math.Ceil(math.Sqrt(float64(len(lm.PictureComps)))))
		rows := int(math.Ceil(float64(len(lm.PictureComps)) / float64(cols)))

		// Calculate cell size
		cellWidth := (availableWidth - float64(cols-1)*lm.Spacing) / float64(cols)
		cellHeight := (availableHeight - float64(rows-1)*lm.Spacing) / float64(rows)

		// Position pictures
		picturesY := lm.Padding
		if lm.TimeComp != nil {
			picturesY += lm.TimeComp.GetHeight() + lm.Spacing
		}
		if lm.TextComp != nil {
			picturesY += lm.TextComp.GetHeight() + lm.Spacing
		}

		for i, pic := range lm.PictureComps {
			row := i / cols
			col := i % cols

			x := lm.Padding + float64(col)*(cellWidth+lm.Spacing)
			y := picturesY + float64(row)*(cellHeight+lm.Spacing)

			pic.CalculateSize(cellWidth, cellHeight)
			pic.SetPosition(x, y)
		}
	}
}

// ToJSON converts the layout to JSON
func (lm *LayoutManager) ToJSON() ([]byte, error) {
	return json.MarshalIndent(lm, "", "  ")
}

// LayoutEngine represents the layout engine
type LayoutEngine struct {
	entry           TestCase
	pages           []Page
	currentPage     *Page
	marginLeft      float64
	marginRight     float64
	marginTop       float64
	marginBottom    float64
	availableWidth  float64
	availableHeight float64
	timeHeight      float64
	lineHeight      float64
	currentY        float64
}

// Page represents a page in the layout
type Page struct {
	Page      int           `json:"page"`
	TimeArea  [][]float64   `json:"time_area"`
	Time      string        `json:"time"`
	TextAreas [][][]float64 `json:"text_areas"`
	Texts     []string      `json:"texts"`
	Pictures  []Picture     `json:"pictures"`
}

// Picture represents a picture in the layout
type Picture struct {
	Index  int         `json:"index"`
	Area   [][]float64 `json:"area"`
	URL    string      `json:"url"`
	Width  int         `json:"width"`
	Height int         `json:"height"`
}

// TestCase represents a test case
type TestCase struct {
	ID       int       `json:"id"`
	Time     string    `json:"time"`
	Text     string    `json:"text"`
	Pictures []Picture `json:"pictures"`
}

// NewLayoutEngine creates a new layout engine
func NewLayoutEngine(entry TestCase) *LayoutEngine {
	engine := &LayoutEngine{
		entry:        entry,
		marginLeft:   100,
		marginRight:  100,
		marginTop:    100,
		marginBottom: 100,
		timeHeight:   100,
		lineHeight:   40,
	}
	engine.availableWidth = 2480 - engine.marginLeft - engine.marginRight
	engine.availableHeight = 3508 - engine.marginTop - engine.marginBottom
	return engine
}

// ProcessEntry processes the entry and returns the layout result
func (e *LayoutEngine) ProcessEntry() ([]Page, error) {
	e.newPage()
	e.addTime()
	e.processText(e.entry.Text)
	e.processPictures(e.entry.Pictures)
	return e.pages, nil
}

func (e *LayoutEngine) newPage() {
	page := &Page{
		Page:      len(e.pages) + 1,
		TextAreas: make([][][]float64, 0),
		Texts:     make([]string, 0),
		Pictures:  make([]Picture, 0),
	}
	e.pages = append(e.pages, *page)
	e.currentPage = &e.pages[len(e.pages)-1]
	e.currentY = e.marginTop
}

func (e *LayoutEngine) addTime() {
	if e.currentPage.Page != 1 {
		return
	}
	timeStr := e.entry.Time
	x0 := e.marginLeft
	y0 := e.marginTop
	x1 := x0 + e.availableWidth
	y1 := y0 + e.timeHeight
	e.currentPage.TimeArea = [][]float64{{x0, y0}, {x1, y1}}
	e.currentPage.Time = timeStr
	e.currentY = y1
}

func (e *LayoutEngine) processText(text string) {
	if text == "" {
		return
	}
	lines := strings.Split(text, "\n")
	currentLine := 0
	for currentLine < len(lines) {
		var remainingHeight float64
		if e.currentPage.Page == 1 {
			remainingHeight = e.availableHeight - e.timeHeight
		} else {
			remainingHeight = e.availableHeight
		}
		availableLines := int(math.Min(float64(len(lines)-currentLine), remainingHeight/e.lineHeight))
		chunk := lines[currentLine : currentLine+availableLines]
		e.addTextChunk(chunk)
		currentLine += availableLines
		if currentLine < len(lines) {
			e.newPage()
		}
	}
}

func (e *LayoutEngine) addTextChunk(chunk []string) {
	startY := e.currentY
	textHeight := float64(len(chunk)) * e.lineHeight
	area := [][]float64{
		{e.marginLeft, startY},
		{e.marginLeft + e.availableWidth, startY + textHeight},
	}
	e.currentPage.TextAreas = append(e.currentPage.TextAreas, area)
	e.currentPage.Texts = append(e.currentPage.Texts, strings.Join(chunk, "\n"))
	e.currentY += textHeight
}

func (e *LayoutEngine) processPictures(pictures []Picture) {
	layout := e.getLayout(len(pictures))
	currentIdx := 0
	for _, row := range layout {
		endIdx := currentIdx + row
		if endIdx > len(pictures) {
			endIdx = len(pictures)
		}
		rowPics := pictures[currentIdx:endIdx]
		currentIdx += row
		e.processPictureRow(rowPics)
	}
}

func (e *LayoutEngine) getLayout(n int) []int {
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

func (e *LayoutEngine) processPictureRow(rowPics []Picture) {
	for {
		availableHeight := e.availableHeight - (e.currentY - e.marginTop)
		if availableHeight <= 0 {
			e.newPage()
			continue
		}

		totalWidth := 0.0
		commonHeight := availableHeight

		scaledWidths := make([]float64, len(rowPics))
		for i, pic := range rowPics {
			aspectRatio := float64(pic.Width) / float64(pic.Height)
			scaledWidth := commonHeight * aspectRatio
			scaledWidths[i] = scaledWidth
			totalWidth += scaledWidth
		}

		if len(rowPics) == 1 {
			if commonHeight <= availableHeight {
				e.placePictures(rowPics, scaledWidths, commonHeight)
				break
			} else {
				e.newPage()
				continue
			}
		} else {
			if totalWidth > e.availableWidth {
				scaleFactor := e.availableWidth / totalWidth
				commonHeight *= scaleFactor
				totalWidth = e.availableWidth
				for i := range scaledWidths {
					scaledWidths[i] *= scaleFactor
				}
			}

			if commonHeight <= availableHeight {
				if math.Abs(totalWidth-e.availableWidth) > 1 {
					e.newPage()
					continue
				}

				e.placePictures(rowPics, scaledWidths, commonHeight)
				break
			} else {
				e.newPage()
			}
		}
	}
}

func (e *LayoutEngine) placePictures(rowPics []Picture, scaledWidths []float64, commonHeight float64) {
	startY := e.currentY
	x := e.marginLeft
	for i, pic := range rowPics {
		area := [][]float64{
			{x, startY},
			{x + scaledWidths[i], startY + commonHeight},
		}
		e.currentPage.Pictures = append(e.currentPage.Pictures, Picture{
			Index: i + 1,
			Area:  area,
			URL:   pic.URL,
		})
		x += scaledWidths[i]
	}
	e.currentY += commonHeight
}

// LayoutArea 表示布局区域
type LayoutArea struct {
	Start [2]int `json:"start"`
	End   [2]int `json:"end"`
}

// LayoutPicture 表示布局后的图片信息
type LayoutPicture struct {
	URL  string     `json:"url"`
	Area LayoutArea `json:"area"`
}

// LayoutPage 表示布局后的页面
type LayoutPage struct {
	Time      string          `json:"time"`
	TimeArea  LayoutArea      `json:"time_area"`
	Texts     []string        `json:"texts"`
	TextAreas []LayoutArea    `json:"text_areas"`
	Pictures  []LayoutPicture `json:"pictures"`
}

// LayoutResult 表示布局结果
type LayoutResult struct {
	Pages []LayoutPage `json:"pages"`
}

// ProcessTestCase processes a test case and returns the layout result
func ProcessTestCase(testCase TestCase) (*LayoutManager, error) {
	// Create layout manager with page dimensions
	layoutManager := NewLayoutManager(800, 1200, 20, 10)

	// Add time component
	layoutManager.AddTimeComponent(testCase.Time)

	// Add text component
	layoutManager.AddTextComponent(testCase.Text)

	// Add picture components
	for _, pic := range testCase.Pictures {
		layoutManager.AddPictureComponent(float64(pic.Width), float64(pic.Height), pic.URL)
	}

	// Perform layout
	layoutManager.Layout()

	return layoutManager, nil
}
