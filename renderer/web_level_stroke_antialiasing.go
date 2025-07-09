package renderer

import (
	"image"
	"image/color"
	"math"

	"github.com/hoonfeng/svg/types"
)

// WebLevelStrokeRenderer Web级别描边抗锯齿渲染器
type WebLevelStrokeRenderer struct {
	*AntiAliasedPathRenderer
}

// NewWebLevelStrokeRenderer 创建Web级别描边抗锯齿渲染器
func NewWebLevelStrokeRenderer() *WebLevelStrokeRenderer {
	return &WebLevelStrokeRenderer{
		AntiAliasedPathRenderer: NewAntiAliasedPathRenderer(),
	}
}

// StrokeWebLevelAntiAliased 使用Web级别超高效抗锯齿算法进行描边
func (r *WebLevelStrokeRenderer) StrokeWebLevelAntiAliased(img *image.RGBA, path []types.Point, strokeColor color.RGBA, strokeWidth float64, closePath bool) {
	if len(path) < 2 {
		return
	}

	// 根据closePath参数决定是否闭合路径
	processedPath := path
	if closePath && len(path) >= 3 {
		// 只有在明确要求闭合且路径有足够点时才闭合
		firstPoint := path[0]
		lastPoint := path[len(path)-1]
		// 检查是否已经闭合（避免重复闭合）
		dx := firstPoint.X - lastPoint.X
		dy := firstPoint.Y - lastPoint.Y
		distance := math.Sqrt(dx*dx + dy*dy)
		if distance > 1.0 { // 如果起点和终点距离大于1像素，才添加闭合线段
			processedPath = make([]types.Point, len(path)+1)
			copy(processedPath, path)
			processedPath[len(path)] = firstPoint
		}
	}

	// 计算边界框 / Calculate bounding box
	minX, minY, maxX, maxY := r.calculateBounds(processedPath)

	// 扩展边界框以包含描边宽度和抗锯齿边缘 / Expand bounds for stroke width and anti-aliasing
	expansion := int(math.Ceil(strokeWidth/2)) + 2
	minX -= expansion
	minY -= expansion
	maxX += expansion
	maxY += expansion

	// 限制边界框在图像范围内 / Clamp bounds to image dimensions
	imgBounds := img.Bounds()
	if minX < imgBounds.Min.X {
		minX = imgBounds.Min.X
	}
	if minY < imgBounds.Min.Y {
		minY = imgBounds.Min.Y
	}
	if maxX >= imgBounds.Max.X {
		maxX = imgBounds.Max.X - 1
	}
	if maxY >= imgBounds.Max.Y {
		maxY = imgBounds.Max.Y - 1
	}

	// 使用Web级别超高效距离场抗锯齿算法
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			// 计算到路径的最短距离（使用优化算法）
			distance := r.calculateWebLevelDistanceToPath(float64(x)+0.5, float64(y)+0.5, processedPath)

			// 使用Web级别描边覆盖率计算
			coverage := r.calculateWebLevelStrokeCoverage(distance, strokeWidth)

			// 设置更严格的覆盖率阈值以避免描边外的轻微渲染
			minCoverage := 0.05 // 只有覆盖率大于5%才进行描边
			if coverage > minCoverage {
				// 混合颜色
				blendedColor := blendColors(getPixelColor(img, x, y), strokeColor, coverage)
				DrawPixel(img, x, y, blendedColor)
			}
		}
	}
}

// calculateWebLevelDistanceToPath Web级别超高效距离计算
func (r *WebLevelStrokeRenderer) calculateWebLevelDistanceToPath(x, y float64, path []types.Point) float64 {
	if len(path) < 2 {
		return math.Inf(1)
	}

	// 使用边界框快速预筛选
	bounds := r.calculatePathBounds(path)

	// 计算到边界框的距离作为下界
	dx := 0.0
	dy := 0.0
	if x < bounds.MinX {
		dx = bounds.MinX - x
	} else if x > bounds.MaxX {
		dx = x - bounds.MaxX
	}
	if y < bounds.MinY {
		dy = bounds.MinY - y
	} else if y > bounds.MaxY {
		dy = y - bounds.MaxY
	}
	boundsDistance := math.Sqrt(dx*dx + dy*dy)

	// 如果边界框距离已经很大，直接返回（超高效优化）
	if boundsDistance > 10.0 {
		return boundsDistance
	}

	minDistance := math.Inf(1)
	minSegmentLength := 0.3 // 更严格的最小线段长度阈值

	// 只检查可能影响结果的线段（Web级别优化）
	for i := 0; i < len(path)-1; i++ {
		// 快速跳过距离过远的线段
		segmentBounds := PathBounds{
			MinX: math.Min(path[i].X, path[i+1].X),
			MaxX: math.Max(path[i].X, path[i+1].X),
			MinY: math.Min(path[i].Y, path[i+1].Y),
			MaxY: math.Max(path[i].Y, path[i+1].Y),
		}

		// 计算到线段边界框的距离
		dx := 0.0
		dy := 0.0
		if x < segmentBounds.MinX {
			dx = segmentBounds.MinX - x
		} else if x > segmentBounds.MaxX {
			dx = x - segmentBounds.MaxX
		}
		if y < segmentBounds.MinY {
			dy = segmentBounds.MinY - y
		} else if y > segmentBounds.MaxY {
			dy = y - segmentBounds.MaxY
		}
		segmentBoundsDistance := math.Sqrt(dx*dx + dy*dy)

		// 如果到线段边界框的距离已经大于当前最小距离，跳过
		if segmentBoundsDistance >= minDistance {
			continue
		}

		// 计算线段长度
		dx = path[i+1].X - path[i].X
		dy = path[i+1].Y - path[i].Y
		segmentLength := math.Sqrt(dx*dx + dy*dy)

		var distance float64
		if segmentLength < minSegmentLength {
			// 线段过短，视为点
			midX := (path[i].X + path[i+1].X) / 2
			midY := (path[i].Y + path[i+1].Y) / 2
			dx := x - midX
			dy := y - midY
			distance = math.Sqrt(dx*dx + dy*dy)
		} else {
			// 使用超高效距离计算
			distance = r.distanceToLineSegmentUltraFast(x, y, path[i].X, path[i].Y, path[i+1].X, path[i+1].Y)
		}

		if distance < minDistance {
			minDistance = distance
		}
	}

	return minDistance
}

// distanceToLineSegmentUltraFast 超高效点到线段距离计算
func (r *WebLevelStrokeRenderer) distanceToLineSegmentUltraFast(px, py, x1, y1, x2, y2 float64) float64 {
	dx := x2 - x1
	dy := y2 - y1
	lengthSq := dx*dx + dy*dy

	// 超快检查：如果线段长度为0
	if lengthSq < 1e-12 {
		dx = px - x1
		dy = py - y1
		return math.Sqrt(dx*dx + dy*dy)
	}

	// 计算投影参数（避免除法，使用更高效的方法）
	dot := (px-x1)*dx + (py-y1)*dy

	// 超快路径：检查投影是否在线段端点
	if dot <= 0 {
		// 最近点是起点
		dx = px - x1
		dy = py - y1
		return math.Sqrt(dx*dx + dy*dy)
	}
	if dot >= lengthSq {
		// 最近点是终点
		dx = px - x2
		dy = py - y2
		return math.Sqrt(dx*dx + dy*dy)
	}

	// 投影在线段内部（使用高精度计算）
	t := dot / lengthSq
	projX := x1 + t*dx
	projY := y1 + t*dy
	dx = px - projX
	dy = py - projY
	return math.Sqrt(dx*dx + dy*dy)
}

// calculateWebLevelStrokeCoverage Web级别描边覆盖率计算
func (r *WebLevelStrokeRenderer) calculateWebLevelStrokeCoverage(distance, strokeWidth float64) float64 {
	halfWidth := strokeWidth / 2
	// 使用更精确的边缘区域以获得Web级别的平滑过渡
	edgeWidth := 1.5 // Web级别抗锯齿边缘宽度

	if distance <= halfWidth-edgeWidth/2 {
		return 1.0 // 完全在描边内
	} else if distance >= halfWidth+edgeWidth/2 {
		return 0.0 // 完全在描边外
	} else {
		// 在边缘区域，使用Web级别的超平滑函数计算覆盖率
		t := (distance - (halfWidth - edgeWidth/2)) / edgeWidth
		// 使用Web级别的超平滑Hermite插值
		return 1.0 - r.webLevelSmoothStep(t)
	}
}

// webLevelSmoothStep Web级别超平滑插值函数
func (r *WebLevelStrokeRenderer) webLevelSmoothStep(t float64) float64 {
	if t <= 0 {
		return 0
	}
	if t >= 1 {
		return 1
	}
	// 使用7次Hermite插值获得Web级别的超平滑边缘：-20t^7 + 70t^6 - 84t^5 + 35t^4
	t2 := t * t
	t3 := t2 * t
	t4 := t3 * t
	t5 := t4 * t
	t6 := t5 * t
	t7 := t6 * t
	return -20*t7 + 70*t6 - 84*t5 + 35*t4
}

// StrokeComplexPathWebLevel 使用Web级别抗锯齿渲染复杂路径描边
func (r *WebLevelStrokeRenderer) StrokeComplexPathWebLevel(img *image.RGBA, subPaths [][]types.Point, strokeColor color.RGBA, strokeWidth float64, closeSubPaths []bool) {
	if len(subPaths) == 0 {
		return
	}

	// 为每个子路径单独处理描边
	for i, subPath := range subPaths {
		if len(subPath) < 2 {
			continue
		}

		// 确定是否闭合当前子路径
		closePath := false
		if i < len(closeSubPaths) {
			closePath = closeSubPaths[i]
		}

		// 使用Web级别抗锯齿渲染当前子路径
		r.StrokeWebLevelAntiAliased(img, subPath, strokeColor, strokeWidth, closePath)
	}
}
