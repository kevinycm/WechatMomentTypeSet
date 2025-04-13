package backend

const (
	minImageHeight   = 800   // 最小展示高度
	minImageWidth    = 1200  // 最小展示宽度
	MinPictureHeight = 200.0 // 最小展示高度
	MinPictureWidth  = 150.0 // 最小展示宽度
)

// getLayout 根据图片数量返回布局规则
func getLayout(n int) []int {
	switch n {
	case 1:
		return []int{1}
	case 2:
		return []int{2}
	case 3:
		return []int{3}
	case 4:
		return []int{2, 2}
	case 5:
		return []int{2, 3}
	case 6:
		return []int{3, 3}
	case 7:
		return []int{2, 2, 3}
	case 8:
		return []int{3, 3, 2}
	case 9:
		return []int{3, 3, 3}
	default:
		return []int{1}
	}
}

// calculateSinglePictureLayout 计算单张图片的布局
// 参数：
//   - availableWidth: 可用宽度
//   - availableHeight: 可用高度
//   - imageWidth: 图片原始宽度
//   - imageHeight: 图片原始高度
//
// 返回：
//   - width: 计算后的宽度
//   - height: 计算后的高度
//   - needNewPage: 是否需要新开一页
func calculateSinglePictureLayout(availableWidth, availableHeight, originalWidth, originalHeight float64) (float64, float64, bool) {
	if originalWidth == 0 || originalHeight == 0 {
		return 0, 0, false
	}

	aspectRatio := originalWidth / originalHeight
	isLongVerticalImage := aspectRatio < 1.0/1.5 // 长宽比小于1/1.5的竖图

	// 以可用区域高度作为基准计算尺寸
	height := availableHeight
	width := height * aspectRatio

	// 如果宽度超过可用宽度，以宽度为基准重新计算
	if width > availableWidth {
		width = availableWidth
		height = width / aspectRatio
	}

	// 检查是否需要新开一页
	needNewPage := false

	// 对于非长竖图，检查最小展示尺寸
	if !isLongVerticalImage {
		if width < MinPictureWidth || height < MinPictureHeight {
			needNewPage = true
		}
	}

	return width, height, needNewPage
}

// calculateSinglePicturePosition 计算单张图片的位置，使其水平居中显示
// 参数：
//   - availableWidth: 可用宽度
//   - availableHeight: 可用高度
//   - imageWidth: 图片宽度
//   - imageHeight: 图片高度
//
// 返回：
//   - x: 图片左上角x坐标
func calculateSinglePicturePosition(availableWidth, availableHeight, imageWidth, imageHeight float64) float64 {
	// 计算水平居中位置
	return (availableWidth - imageWidth) / 2
}

// calculatePictureRowLayout 计算一行图片的布局
// 参数：
//   - availableWidth: 可用宽度
//   - availableHeight: 可用高度
//   - pictures: 图片数组
//   - spacing: 图片之间的间距
//
// 返回：
//   - widths: 每张图片的宽度
//   - heights: 每张图片的高度
//   - needNewPage: 是否需要新开一页
func calculatePictureRowLayout(availableWidth, availableHeight float64, pictures []Picture, spacing float64) (widths, heights []float64, needNewPage bool) {
	// 计算每张图片的宽高比
	aspectRatios := make([]float64, len(pictures))
	for i, pic := range pictures {
		aspectRatios[i] = float64(pic.Width) / float64(pic.Height)
	}

	// 计算每张图片的宽度
	totalSpacing := spacing * float64(len(pictures)-1)
	availableWidthForPictures := availableWidth - totalSpacing
	widthPerPicture := availableWidthForPictures / float64(len(pictures))

	// 计算每张图片的高度
	widths = make([]float64, len(pictures))
	heights = make([]float64, len(pictures))
	maxHeight := 0.0

	for i, aspectRatio := range aspectRatios {
		widths[i] = widthPerPicture
		heights[i] = widthPerPicture / aspectRatio
		if heights[i] > maxHeight {
			maxHeight = heights[i]
		}
	}

	// 检查是否满足最小尺寸要求
	if maxHeight < minImageHeight {
		// 尝试以最小高度为基准计算
		maxHeight = minImageHeight
		for i, aspectRatio := range aspectRatios {
			heights[i] = maxHeight
			widths[i] = maxHeight * aspectRatio
		}

		// 如果计算出的总宽度超过可用宽度，需要新开一页
		totalWidth := 0.0
		for _, width := range widths {
			totalWidth += width
		}
		if totalWidth+totalSpacing > availableWidth {
			needNewPage = true
		}
	}

	return widths, heights, needNewPage
}

// calculatePictureRowPosition 计算一行图片的位置
// 参数：
//   - availableWidth: 可用宽度
//   - availableHeight: 可用高度
//   - widths: 每张图片的宽度
//   - heights: 每张图片的高度
//   - spacing: 图片之间的间距
//
// 返回：
//   - positions: 每张图片的位置 [x, y]
func calculatePictureRowPosition(availableWidth, availableHeight float64, widths, heights []float64, spacing float64) [][]float64 {
	positions := make([][]float64, len(widths))
	totalWidth := 0.0
	for _, width := range widths {
		totalWidth += width
	}
	totalWidth += spacing * float64(len(widths)-1)

	// 计算起始x坐标，使整行图片居中
	startX := (availableWidth - totalWidth) / 2
	currentX := startX

	// 计算每张图片的位置
	for i := range widths {
		positions[i] = []float64{currentX, 0}
		currentX += widths[i] + spacing
	}

	return positions
}
