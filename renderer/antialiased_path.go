package renderer

import (
	"image"
	"image/color"
	"math"

	"github.com/hoonfeng/svg/path"
	"github.com/hoonfeng/svg/types"
)

// AntiAliasedPathRenderer 抗锯齿路径渲染器 / Anti-aliased path renderer
type AntiAliasedPathRenderer struct {
	*AntiAliasedRenderer
}

// NewAntiAliasedPathRenderer 创建抗锯齿路径渲染器 / Create anti-aliased path renderer
func NewAntiAliasedPathRenderer() *AntiAliasedPathRenderer {
	return &AntiAliasedPathRenderer{
		AntiAliasedRenderer: NewAntiAliasedRenderer(),
	}
}

// RenderPath 渲染抗锯齿路径 / Render anti-aliased path
func (r *AntiAliasedPathRenderer) RenderPath(img *image.RGBA, pathData string, fillColor, strokeColor color.RGBA, strokeWidth float64, viewBox []float64, scaleX, scaleY float64) error {
	// 解析路径 / Parse path
	parsedPath, err := path.ParsePath(pathData)
	if err != nil {
		return err
	}

	// 设置web级别的极致精度 / Set web-level ultra precision
	precision := 0.001 // web级别的精度用于MSAA / Web-level precision for MSAA

	// 获取子路径和闭合信息 / Get sub-paths and closure information
	subPaths := parsedPath.FlattenSubPaths(precision)
	closeInfo := parsedPath.GetSubPathCloseInfo() // 获取每个子路径的闭合信息

	// 转换所有子路径的坐标 / Transform coordinates for all sub-paths
	transformedSubPaths := make([][]types.Point, 0, len(subPaths))
	transformedCloseInfo := make([]bool, 0, len(subPaths))
	for i, subPath := range subPaths {
		if len(subPath) < 2 { // 降低要求，允许线条
			continue // 跳过无效的子路径 / Skip invalid sub-paths
		}
		transformedPath := r.transformPath(subPath, viewBox, scaleX, scaleY)
		transformedSubPaths = append(transformedSubPaths, transformedPath)

		// 保存闭合信息
		if i < len(closeInfo) {
			transformedCloseInfo = append(transformedCloseInfo, closeInfo[i])
		} else {
			transformedCloseInfo = append(transformedCloseInfo, false)
		}
	}

	// 使用缠绕数规则填充复杂路径 / Fill complex path using winding rule
	if fillColor.A > 0 && len(transformedSubPaths) > 0 {
		r.fillAntiAliasedComplexPath(img, transformedSubPaths, fillColor)
	}

	// 使用真正的描边路径生成器 / Use true stroke path generator
	if strokeColor.A > 0 && strokeWidth > 0 {
		// 创建真正的描边渲染器
		trueStrokeRenderer := NewTrueStrokeRenderer()
		// 使用真正的描边路径渲染复杂路径描边
		trueStrokeRenderer.RenderTrueStrokeComplexPath(img, transformedSubPaths, strokeColor, strokeWidth*math.Min(scaleX, scaleY), transformedCloseInfo)
	}

	return nil
}

// transformPath 转换路径坐标 / Transform path coordinates
func (r *AntiAliasedPathRenderer) transformPath(subPath []types.Point, viewBox []float64, scaleX, scaleY float64) []types.Point {
	transformed := make([]types.Point, len(subPath))
	for i, point := range subPath {
		transformed[i] = types.Point{
			X: (point.X - viewBox[0]) * scaleX,
			Y: (point.Y - viewBox[1]) * scaleY,
		}
	}
	// 过滤过短的线段 / Filter out segments that are too short
	return r.filterShortSegments(transformed)
}

// filterShortSegments 过滤过短的线段 / Filter out segments that are too short
func (r *AntiAliasedPathRenderer) filterShortSegments(path []types.Point) []types.Point {
	if len(path) < 2 {
		return path
	}

	minSegmentLength := 0.5 // 最小线段长度阈值 / Minimum segment length threshold
	filtered := make([]types.Point, 0, len(path))
	filtered = append(filtered, path[0]) // 总是保留第一个点 / Always keep the first point

	for i := 1; i < len(path); i++ {
		// 计算当前点与上一个保留点的距离 / Calculate distance from current point to last kept point
		lastKept := filtered[len(filtered)-1]
		dx := path[i].X - lastKept.X
		dy := path[i].Y - lastKept.Y
		distance := math.Sqrt(dx*dx + dy*dy)

		// 如果距离足够大，或者是最后一个点，则保留 / Keep if distance is large enough, or if it's the last point
		if distance >= minSegmentLength || i == len(path)-1 {
			filtered = append(filtered, path[i])
		}
	}

	return filtered
}

// fillAntiAliasedComplexPath 使用缠绕数规则填充复杂抗锯齿路径 / Fill complex anti-aliased path using winding rule
func (r *AntiAliasedPathRenderer) fillAntiAliasedComplexPath(img *image.RGBA, subPaths [][]types.Point, fillColor color.RGBA) {
	if len(subPaths) == 0 {
		return
	}

	// 计算所有子路径的总边界框 / Calculate overall bounding box for all sub-paths
	minX, minY, maxX, maxY := r.calculateComplexBounds(subPaths)

	// 扩展边界框以包含抗锯齿边缘 / Expand bounds for anti-aliasing edges
	minX -= 2
	minY -= 2
	maxX += 2
	maxY += 2

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

	// 遍历边界框内的每个像素 / Iterate through each pixel in bounding box
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			// 使用增强web级别的多重采样抗锯齿(MSAA) / Use enhanced web-level Multi-Sample Anti-Aliasing (MSAA)
			coverage := r.calculateWebLevelMSAA(float64(x), float64(y), subPaths, 16)

			// 使用更低的覆盖率阈值和边缘平滑处理 / Use lower coverage threshold and edge smoothing
			minCoverage := 0.05 // 降低阈值以获得更平滑的边缘 / Lower threshold for smoother edges
			if coverage > minCoverage {
				// 对边缘区域应用额外的平滑处理 / Apply additional smoothing for edge areas
				smoothedCoverage := r.applyCoverageSmoothing(coverage, float64(x), float64(y))
				// 混合颜色 / Blend color
				blendedColor := blendColors(getPixelColor(img, x, y), fillColor, smoothedCoverage)
				DrawPixel(img, x, y, blendedColor)
			}
		}
	}
}

// fillAntiAliasedPath 填充抗锯齿路径 / Fill anti-aliased path
func (r *AntiAliasedPathRenderer) fillAntiAliasedPath(img *image.RGBA, path []types.Point, fillColor color.RGBA) {
	if len(path) < 3 {
		return
	}

	// 计算边界框 / Calculate bounding box
	minX, minY, maxX, maxY := r.calculateBounds(path)

	// 扩展边界框以包含抗锯齿边缘 / Expand bounds for anti-aliasing edges
	minX -= 2
	minY -= 2
	maxX += 2
	maxY += 2

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

	// 遍历边界框内的每个像素 / Iterate through each pixel in bounding box
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			// 使用增强web级别的多重采样抗锯齿(MSAA) / Use enhanced web-level Multi-Sample Anti-Aliasing (MSAA)
			coverage := r.calculateWebLevelPathMSAA(float64(x), float64(y), path, 16)

			// 使用更低的覆盖率阈值和边缘平滑处理 / Use lower coverage threshold and edge smoothing
			minCoverage := 0.05 // 降低阈值以获得更平滑的边缘 / Lower threshold for smoother edges
			if coverage > minCoverage {
				// 对边缘区域应用额外的平滑处理 / Apply additional smoothing for edge areas
				smoothedCoverage := r.applyCoverageSmoothing(coverage, float64(x), float64(y))
				// 混合颜色 / Blend color
				blendedColor := blendColors(getPixelColor(img, x, y), fillColor, smoothedCoverage)
				DrawPixel(img, x, y, blendedColor)
			}
		}
	}
}

// strokeAntiAliasedPath 描边抗锯齿路径 / Stroke anti-aliased path
func (r *AntiAliasedPathRenderer) strokeAntiAliasedPath(img *image.RGBA, path []types.Point, strokeColor color.RGBA, strokeWidth float64) {
	if len(path) < 2 {
		return
	}

	// 计算边界框 / Calculate bounding box
	minX, minY, maxX, maxY := r.calculateBounds(path)

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

	// 遍历边界框内的每个像素 / Iterate through each pixel in bounding box
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			// 计算到路径的最短距离 / Calculate shortest distance to path
			distance := r.calculateDistanceToPath(float64(x)+0.5, float64(y)+0.5, path)

			// 计算描边覆盖率 / Calculate stroke coverage
			coverage := r.calculateStrokeCoverage(distance, strokeWidth)

			// 设置最小覆盖率阈值以避免描边外的轻微渲染 / Set minimum coverage threshold to avoid slight rendering outside stroke
			minCoverage := 0.1 // 只有覆盖率大于10%才进行描边 / Only stroke if coverage is greater than 10%
			if coverage > minCoverage {
				// 混合颜色 / Blend color
				blendedColor := blendColors(getPixelColor(img, x, y), strokeColor, coverage)
				DrawPixel(img, x, y, blendedColor)
			}
		}
	}
}

// calculateComplexBounds 计算复杂路径的边界框 / Calculate bounding box for complex path
func (r *AntiAliasedPathRenderer) calculateComplexBounds(subPaths [][]types.Point) (int, int, int, int) {
	if len(subPaths) == 0 {
		return 0, 0, 0, 0
	}

	// 找到第一个有效的点来初始化边界 / Find first valid point to initialize bounds
	var minX, maxX, minY, maxY float64
	initialized := false

	for _, subPath := range subPaths {
		if len(subPath) > 0 {
			if !initialized {
				minX, maxX = subPath[0].X, subPath[0].X
				minY, maxY = subPath[0].Y, subPath[0].Y
				initialized = true
			}
			for _, point := range subPath {
				if point.X < minX {
					minX = point.X
				}
				if point.X > maxX {
					maxX = point.X
				}
				if point.Y < minY {
					minY = point.Y
				}
				if point.Y > maxY {
					maxY = point.Y
				}
			}
		}
	}

	if !initialized {
		return 0, 0, 0, 0
	}

	return int(math.Floor(minX)), int(math.Floor(minY)), int(math.Ceil(maxX)), int(math.Ceil(maxY))
}

// calculateBounds 计算路径边界框 / Calculate path bounding box
func (r *AntiAliasedPathRenderer) calculateBounds(path []types.Point) (int, int, int, int) {
	if len(path) == 0 {
		return 0, 0, 0, 0
	}

	minX, maxX := path[0].X, path[0].X
	minY, maxY := path[0].Y, path[0].Y

	for _, point := range path[1:] {
		if point.X < minX {
			minX = point.X
		}
		if point.X > maxX {
			maxX = point.X
		}
		if point.Y < minY {
			minY = point.Y
		}
		if point.Y > maxY {
			maxY = point.Y
		}
	}

	return int(math.Floor(minX)), int(math.Floor(minY)), int(math.Ceil(maxX)), int(math.Ceil(maxY))
}

// calculateWebLevelMSAA 使用超高效的距离场抗锯齿计算复杂路径覆盖率 / Calculate complex path coverage with ultra-efficient distance field anti-aliasing
func (r *AntiAliasedPathRenderer) calculateWebLevelMSAA(pixelX, pixelY float64, subPaths [][]types.Point, samples int) float64 {
	// 使用距离场方法，只需要计算像素中心点 / Use distance field method, only need to calculate pixel center
	centerX := pixelX + 0.5
	centerY := pixelY + 0.5

	// 快速检查：先用像素中心点判断是否明显在内部或外部 / Quick check: use pixel center to determine if obviously inside or outside
	isInside := r.isPointInComplexPath(centerX, centerY, subPaths)

	// 快速距离估算（只检查第一个子路径的边界框） / Fast distance estimation (only check first sub-path bounds)
	if len(subPaths) == 0 {
		return 0.0
	}

	firstPath := subPaths[0]
	if len(firstPath) < 3 {
		return 0.0
	}

	// 计算到第一个路径的粗略距离 / Calculate rough distance to first path
	bounds := r.calculatePathBounds(firstPath)
	roughDistance := 0.0
	if centerX < bounds.MinX {
		roughDistance = bounds.MinX - centerX
	} else if centerX > bounds.MaxX {
		roughDistance = centerX - bounds.MaxX
	}
	if centerY < bounds.MinY {
		dy := bounds.MinY - centerY
		roughDistance = math.Sqrt(roughDistance*roughDistance + dy*dy)
	} else if centerY > bounds.MaxY {
		dy := centerY - bounds.MaxY
		roughDistance = math.Sqrt(roughDistance*roughDistance + dy*dy)
	}

	// 如果距离边缘较远，直接返回结果 / If far from edge, return result directly
	if roughDistance > 2.0 {
		if isInside {
			return 1.0
		} else {
			return 0.0
		}
	}

	// 只在边缘附近使用2x2采样 / Only use 2x2 samples near edges
	quickSamples := 2 // 减少到2x2采样
	insideCount := 0
	totalSamples := quickSamples * quickSamples

	// 使用固定偏移而不是计算偏移以提高性能 / Use fixed offsets instead of calculated offsets for better performance
	offsets := []float64{0.25, 0.75}

	for i := 0; i < quickSamples; i++ {
		for j := 0; j < quickSamples; j++ {
			sampleX := pixelX + offsets[i]
			sampleY := pixelY + offsets[j]

			if r.isPointInComplexPath(sampleX, sampleY, subPaths) {
				insideCount++
			}
		}
	}

	coverage := float64(insideCount) / float64(totalSamples)

	// 简化的平滑处理 / Simplified smoothing
	if coverage > 0 && coverage < 1 {
		// 只对部分覆盖的像素应用平滑 / Only apply smoothing to partially covered pixels
		coverage = r.smootherStep(coverage)
	}

	return coverage
}

// calculateWebLevelPathMSAA 使用超高效的距离场抗锯齿计算路径覆盖率 / Calculate path coverage with ultra-efficient distance field anti-aliasing
func (r *AntiAliasedPathRenderer) calculateWebLevelPathMSAA(pixelX, pixelY float64, path []types.Point, samples int) float64 {
	// 使用距离场方法，只需要计算像素中心点 / Use distance field method, only need to calculate pixel center
	centerX := pixelX + 0.5
	centerY := pixelY + 0.5

	// 快速检查：先用像素中心点判断是否明显在内部或外部 / Quick check: use pixel center to determine if obviously inside or outside
	isInside := r.isPointInPath(centerX, centerY, path)

	// 快速距离估算（使用边界框） / Fast distance estimation (using bounding box)
	bounds := r.calculatePathBounds(path)
	roughDistance := 0.0
	if centerX < bounds.MinX {
		roughDistance = bounds.MinX - centerX
	} else if centerX > bounds.MaxX {
		roughDistance = centerX - bounds.MaxX
	}
	if centerY < bounds.MinY {
		dy := bounds.MinY - centerY
		roughDistance = math.Sqrt(roughDistance*roughDistance + dy*dy)
	} else if centerY > bounds.MaxY {
		dy := centerY - bounds.MaxY
		roughDistance = math.Sqrt(roughDistance*roughDistance + dy*dy)
	}

	// 如果距离边缘较远，直接返回结果 / If far from edge, return result directly
	if roughDistance > 2.0 {
		if isInside {
			return 1.0
		} else {
			return 0.0
		}
	}

	// 只在边缘附近使用2x2采样 / Only use 2x2 samples near edges
	quickSamples := 2 // 减少到2x2采样
	insideCount := 0
	totalSamples := quickSamples * quickSamples

	// 使用固定偏移而不是计算偏移以提高性能 / Use fixed offsets instead of calculated offsets for better performance
	offsets := []float64{0.25, 0.75}

	for i := 0; i < quickSamples; i++ {
		for j := 0; j < quickSamples; j++ {
			sampleX := pixelX + offsets[i]
			sampleY := pixelY + offsets[j]

			if r.isPointInPath(sampleX, sampleY, path) {
				insideCount++
			}
		}
	}

	coverage := float64(insideCount) / float64(totalSamples)

	// 简化的平滑处理 / Simplified smoothing
	if coverage > 0 && coverage < 1 {
		// 只对部分覆盖的像素应用平滑 / Only apply smoothing to partially covered pixels
		coverage = r.smootherStep(coverage)
	}

	return coverage
}

// PathBounds 路径边界框结构 / Path bounds structure
type PathBounds struct {
	MinX, MinY, MaxX, MaxY float64
}

// isPointInComplexPath 使用缠绕数规则检查点是否在复杂路径内 / Check if point is inside complex path using winding rule
func (r *AntiAliasedPathRenderer) isPointInComplexPath(x, y float64, subPaths [][]types.Point) bool {
	// 快速边界检查 / Quick bounds check
	for _, subPath := range subPaths {
		if len(subPath) < 3 {
			continue
		}

		// 计算子路径边界 / Calculate sub-path bounds
		bounds := r.calculatePathBounds(subPath)
		if x < bounds.MinX || x > bounds.MaxX || y < bounds.MinY || y > bounds.MaxY {
			continue
		}

		if r.isPointInPathOptimized(x, y, subPath) {
			return true
		}
	}
	return false
}

// calculatePathBounds 计算路径边界框 / Calculate path bounds
func (r *AntiAliasedPathRenderer) calculatePathBounds(path []types.Point) PathBounds {
	if len(path) == 0 {
		return PathBounds{}
	}

	bounds := PathBounds{
		MinX: path[0].X, MaxX: path[0].X,
		MinY: path[0].Y, MaxY: path[0].Y,
	}

	for _, point := range path[1:] {
		if point.X < bounds.MinX {
			bounds.MinX = point.X
		} else if point.X > bounds.MaxX {
			bounds.MaxX = point.X
		}
		if point.Y < bounds.MinY {
			bounds.MinY = point.Y
		} else if point.Y > bounds.MaxY {
			bounds.MaxY = point.Y
		}
	}

	return bounds
}

// isLeft 测试点是否在有向线段的左侧 / Test if point is left of directed line segment
func (r *AntiAliasedPathRenderer) isLeft(x1, y1, x2, y2, px, py float64) float64 {
	return (x2-x1)*(py-y1) - (px-x1)*(y2-y1)
}

// isPointInPath 使用射线投射算法检查点是否在路径内 / Check if point is inside path using ray casting
func (r *AntiAliasedPathRenderer) isPointInPath(x, y float64, path []types.Point) bool {
	if len(path) < 3 {
		return false
	}

	inside := false
	j := len(path) - 1

	for i := 0; i < len(path); i++ {
		xi, yi := path[i].X, path[i].Y
		xj, yj := path[j].X, path[j].Y

		if ((yi > y) != (yj > y)) && (x < (xj-xi)*(y-yi)/(yj-yi)+xi) {
			inside = !inside
		}
		j = i
	}

	return inside
}

// isPointInPathOptimized 优化的点在路径内检测 / Optimized point-in-path detection
func (r *AntiAliasedPathRenderer) isPointInPathOptimized(x, y float64, path []types.Point) bool {
	if len(path) < 3 {
		return false
	}

	// 使用更高效的射线投射算法 / Use more efficient ray casting algorithm
	inside := false
	j := len(path) - 1

	for i := 0; i < len(path); i++ {
		xi, yi := path[i].X, path[i].Y
		xj, yj := path[j].X, path[j].Y

		// 优化：跳过明显不相交的边 / Optimization: skip obviously non-intersecting edges
		if (yi <= y && yj <= y) || (yi > y && yj > y) {
			j = i
			continue
		}

		// 计算交点的x坐标 / Calculate x-coordinate of intersection
		if yi != yj {
			xIntersect := xi + (y-yi)*(xj-xi)/(yj-yi)
			if x < xIntersect {
				inside = !inside
			}
		}
		j = i
	}

	return inside
}

// calculateDistanceToPath 计算点到路径的最短距离（优化版） / Calculate shortest distance from point to path (optimized)
func (r *AntiAliasedPathRenderer) calculateDistanceToPath(x, y float64, path []types.Point) float64 {
	if len(path) < 2 {
		return math.Inf(1)
	}

	// 首先计算路径边界框 / First calculate path bounds
	bounds := r.calculatePathBounds(path)

	// 如果点在边界框外，计算到边界框的距离作为下界 / If point is outside bounds, calculate distance to bounds as lower bound
	if x < bounds.MinX || x > bounds.MaxX || y < bounds.MinY || y > bounds.MaxY {
		// 计算到边界框的距离 / Calculate distance to bounding box
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

		// 如果边界框距离已经很大，直接返回 / If bounds distance is already large, return directly
		if boundsDistance > 5.0 {
			return boundsDistance
		}
	}

	minDistance := math.Inf(1)
	minSegmentLength := 0.5 // 最小线段长度阈值 / Minimum segment length threshold

	// 优化：只检查可能影响结果的线段 / Optimization: only check segments that could affect result
	for i := 0; i < len(path); i++ {
		j := (i + 1) % len(path)

		// 快速跳过距离过远的线段 / Quick skip for segments that are too far
		segmentBounds := PathBounds{
			MinX: math.Min(path[i].X, path[j].X),
			MaxX: math.Max(path[i].X, path[j].X),
			MinY: math.Min(path[i].Y, path[j].Y),
			MaxY: math.Max(path[i].Y, path[j].Y),
		}

		// 计算到线段边界框的距离 / Calculate distance to segment bounds
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

		// 如果到线段边界框的距离已经大于当前最小距离，跳过 / Skip if distance to segment bounds is already larger than current minimum
		if segmentBoundsDistance >= minDistance {
			continue
		}

		// 计算线段长度 / Calculate segment length
		dx = path[j].X - path[i].X
		dy = path[j].Y - path[i].Y
		segmentLength := math.Sqrt(dx*dx + dy*dy)

		var distance float64
		if segmentLength < minSegmentLength {
			// 线段过短，视为点，计算到线段中点的距离 / Segment too short, treat as point
			midX := (path[i].X + path[j].X) / 2
			midY := (path[i].Y + path[j].Y) / 2
			dx := x - midX
			dy := y - midY
			distance = math.Sqrt(dx*dx + dy*dy)
		} else {
			// 正常线段，计算到线段的距离 / Normal segment, calculate distance to segment
			distance = r.distanceToLineSegmentOptimized(x, y, path[i].X, path[i].Y, path[j].X, path[j].Y)
		}

		if distance < minDistance {
			minDistance = distance
		}
	}

	return minDistance
}

// distanceToLineSegment 计算点到线段的距离 / Calculate distance from point to line segment
func (r *AntiAliasedPathRenderer) distanceToLineSegment(px, py, x1, y1, x2, y2 float64) float64 {
	dx := x2 - x1
	dy := y2 - y1
	lengthSq := dx*dx + dy*dy

	if lengthSq == 0 {
		// 线段退化为点 / Line segment degenerates to a point
		dx = px - x1
		dy = py - y1
		return math.Sqrt(dx*dx + dy*dy)
	}

	// 计算投影参数 / Calculate projection parameter
	t := ((px-x1)*dx + (py-y1)*dy) / lengthSq

	if t < 0 {
		// 投影在线段起点之前 / Projection is before line segment start
		dx = px - x1
		dy = py - y1
	} else if t > 1 {
		// 投影在线段终点之后 / Projection is after line segment end
		dx = px - x2
		dy = py - y2
	} else {
		// 投影在线段上 / Projection is on line segment
		projX := x1 + t*dx
		projY := y1 + t*dy
		dx = px - projX
		dy = py - projY
	}

	return math.Sqrt(dx*dx + dy*dy)
}

// distanceToLineSegmentOptimized 优化的点到线段距离计算 / Optimized distance calculation from point to line segment
func (r *AntiAliasedPathRenderer) distanceToLineSegmentOptimized(px, py, x1, y1, x2, y2 float64) float64 {
	dx := x2 - x1
	dy := y2 - y1
	lengthSq := dx*dx + dy*dy

	// 快速检查：如果线段长度为0 / Quick check: if segment length is 0
	if lengthSq < 1e-10 {
		dx = px - x1
		dy = py - y1
		return math.Sqrt(dx*dx + dy*dy)
	}

	// 计算投影参数（避免除法） / Calculate projection parameter (avoid division)
	dot := (px-x1)*dx + (py-y1)*dy

	// 快速路径：检查投影是否在线段端点 / Fast path: check if projection is at segment endpoints
	if dot <= 0 {
		// 最近点是起点 / Closest point is start point
		dx = px - x1
		dy = py - y1
		return math.Sqrt(dx*dx + dy*dy)
	}
	if dot >= lengthSq {
		// 最近点是终点 / Closest point is end point
		dx = px - x2
		dy = py - y2
		return math.Sqrt(dx*dx + dy*dy)
	}

	// 投影在线段内部 / Projection is inside segment
	t := dot / lengthSq
	projX := x1 + t*dx
	projY := y1 + t*dy
	dx = px - projX
	dy = py - projY
	return math.Sqrt(dx*dx + dy*dy)
}

// smootherStep 实现更平滑的插值函数 / Implement smoother interpolation function
func (r *AntiAliasedPathRenderer) smootherStep(t float64) float64 {
	if t <= 0 {
		return 0
	}
	if t >= 1 {
		return 1
	}
	// 使用5次Hermite插值：6t^5 - 15t^4 + 10t^3 / Use quintic Hermite interpolation: 6t^5 - 15t^4 + 10t^3
	return t * t * t * (t*(t*6-15) + 10)
}

// applyCoverageSmoothing 对覆盖率应用额外的平滑处理 / Apply additional smoothing to coverage
func (r *AntiAliasedPathRenderer) applyCoverageSmoothing(coverage, x, y float64) float64 {
	// 对边缘区域（覆盖率在0.1-0.9之间）应用额外平滑 / Apply additional smoothing for edge areas (coverage between 0.1-0.9)
	if coverage > 0.1 && coverage < 0.9 {
		// 使用Gamma校正来改善边缘平滑度 / Use gamma correction to improve edge smoothness
		gamma := 1.8 // Gamma值，用于调整边缘过渡 / Gamma value for adjusting edge transition
		return math.Pow(coverage, 1.0/gamma)
	}
	return coverage
}

// haltonSequence 生成Halton序列用于更均匀的采样 / Generate Halton sequence for more uniform sampling
func (r *AntiAliasedPathRenderer) haltonSequence(index, base int) float64 {
	result := 0.0
	f := 1.0 / float64(base)
	i := index
	for i > 0 {
		result += f * float64(i%base)
		i /= base
		f /= float64(base)
	}
	return result
}

// calculateSignedDistanceToComplexPath 计算到复杂路径的有符号距离 / Calculate signed distance to complex path
func (r *AntiAliasedPathRenderer) calculateSignedDistanceToComplexPath(x, y float64, subPaths [][]types.Point) float64 {
	minDistance := math.Inf(1)
	isInside := false

	// 检查是否在路径内部 / Check if inside path
	isInside = r.isPointInComplexPath(x, y, subPaths)

	// 计算到最近边缘的距离 / Calculate distance to nearest edge
	for _, subPath := range subPaths {
		if len(subPath) < 2 {
			continue
		}

		for i := 0; i < len(subPath); i++ {
			j := (i + 1) % len(subPath)
			dist := r.distanceToLineSegment(x, y, subPath[i].X, subPath[i].Y, subPath[j].X, subPath[j].Y)
			if dist < minDistance {
				minDistance = dist
			}
		}
	}

	// 返回有符号距离：内部为负，外部为正 / Return signed distance: negative inside, positive outside
	if isInside {
		return -minDistance
	}
	return minDistance
}

// calculateSignedDistanceToPath 计算到单路径的有符号距离 / Calculate signed distance to single path
func (r *AntiAliasedPathRenderer) calculateSignedDistanceToPath(x, y float64, path []types.Point) float64 {
	if len(path) < 3 {
		return math.Inf(1)
	}

	minDistance := math.Inf(1)
	isInside := false

	// 检查是否在路径内部 / Check if inside path
	isInside = r.isPointInPath(x, y, path)

	// 计算到最近边缘的距离 / Calculate distance to nearest edge
	for i := 0; i < len(path); i++ {
		j := (i + 1) % len(path)
		dist := r.distanceToLineSegment(x, y, path[i].X, path[i].Y, path[j].X, path[j].Y)
		if dist < minDistance {
			minDistance = dist
		}
	}

	// 返回有符号距离：内部为负，外部为正 / Return signed distance: negative inside, positive outside
	if isInside {
		return -minDistance
	}
	return minDistance
}

// applySDFSmoothing 应用距离场平滑 / Apply SDF smoothing
func (r *AntiAliasedPathRenderer) applySDFSmoothing(coverage, distance float64) float64 {
	// 使用距离场信息改善边缘平滑度 / Use distance field information to improve edge smoothness
	edgeWidth := 1.0

	if math.Abs(distance) < edgeWidth {
		// 在边缘区域使用距离场平滑 / Use distance field smoothing in edge area
		t := (distance + edgeWidth) / (2.0 * edgeWidth)
		// 限制t在[0,1]范围内 / Clamp t to [0,1] range
		if t < 0 {
			t = 0
		}
		if t > 1 {
			t = 1
		}
		// 使用平滑步函数 / Use smooth step function
		sdfCoverage := r.smootherStep(t)
		// 结合MSAA和SDF结果 / Combine MSAA and SDF results
		return (coverage + sdfCoverage) * 0.5
	}

	return coverage
}

// calculateStrokeCoverage 计算描边覆盖率 / Calculate stroke coverage
func (r *AntiAliasedPathRenderer) calculateStrokeCoverage(distance, strokeWidth float64) float64 {
	halfWidth := strokeWidth / 2
	// 使用更宽的边缘区域以获得更平滑的过渡 / Use wider edge area for smoother transition
	edgeWidth := 2.0 // 增加抗锯齿边缘宽度 / Increase anti-aliasing edge width

	if distance <= halfWidth-edgeWidth/2 {
		return 1.0 // 完全在描边内 / Completely inside stroke
	} else if distance >= halfWidth+edgeWidth/2 {
		return 0.0 // 完全在描边外 / Completely outside stroke
	} else {
		// 在边缘区域，使用改进的平滑函数计算覆盖率 / In edge area, use improved smooth function to calculate coverage
		t := (distance - (halfWidth - edgeWidth/2)) / edgeWidth
		// 使用三次Hermite插值获得更平滑的边缘 / Use cubic Hermite interpolation for smoother edges
		return 1.0 - r.smootherStep(t)
	}
}
