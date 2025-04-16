package calculate

// Picture represents a picture in the layout
type Picture struct {
	Index  int         `json:"index"`
	Area   [][]float64 `json:"area"`
	URL    string      `json:"url"`
	Width  int         `json:"width"`
	Height int         `json:"height"`
}

// Entry represents a single moment entry with time, text and pictures
type Entry struct {
	Time     string    `json:"time"`
	Text     string    `json:"text"`
	Pictures []Picture `json:"pictures"`
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

// ContinuousLayoutPage represents a single page in the continuous layout
type ContinuousLayoutPage struct {
	Page      int         `json:"page"`
	IsInsert  bool        `json:"is_insert"`  // 是否是插页
	YearMonth string      `json:"year_month"` // 年月信息，格式：2025年3月
	Entries   []PageEntry `json:"entries"`
}

// ContinuousLayoutEngine represents the continuous layout engine
type ContinuousLayoutEngine struct {
	entries                          []Entry
	pages                            []ContinuousLayoutPage
	currentPage                      *ContinuousLayoutPage
	marginLeft                       float64
	marginRight                      float64
	marginTop                        float64
	marginBottom                     float64
	availableWidth                   float64
	availableHeight                  float64
	timeHeight                       float64
	fontSize                         float64
	lineHeight                       float64
	currentY                         float64
	timeAreaBottom                   float64
	entrySpacing                     float64 // 条目之间的间距
	elementSpacing                   float64 // 元素整体之间的间距
	imageSpacing                     float64 // 图片之间的间距
	minWideHeight                    float64 // Min height for Wide pics (AR >= 3)
	minTallHeight                    float64 // Min height for Tall pics (AR <= 1/3)
	minLandscapeHeight               float64 // Base Min height for Landscape pics (1 < AR < 3) - Used for < 5 pics
	minPortraitHeight                float64 // Base Min height for Portrait pics (1/3 < AR < 1) - Used for < 5 pics
	minLandscapeHeightLargeGroup     float64 // Added: Min height for Landscape pics (>= 5 pics)
	minPortraitHeightLargeGroup      float64 // Added: Min height for Portrait pics (>= 5 pics)
	minLandscapeHeightVeryLargeGroup float64 // Added: Min height for Landscape pics (>= 8 pics)
	minPortraitHeightVeryLargeGroup  float64 // Added: Min height for Portrait pics (>= 8 pics)
	singleImageHeight                float64 // 单张竖图的默认高度
	singleImageWidth                 float64 // 单张横图的默认宽度
	minImageHeight                   float64 // 单张竖图的最小展示高度
	minImageWidth                    float64 // 单张横图的最小展示宽度
	currentYearMonth                 string
	bottomMargin                     float64
}

// TemplateLayout holds the calculated positions and dimensions for a template
type TemplateLayout struct {
	Positions   [][]float64 // Relative positions [x, y] for top-left corner of each pic within the layout block
	Dimensions  [][]float64 // Dimensions [width, height] for each pic
	TotalHeight float64     // Total height of the layout block (including internal spacing)
	TotalWidth  float64     // Total width (should generally match e.availableWidth)
}
