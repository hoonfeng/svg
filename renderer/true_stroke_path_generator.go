package renderer

import (
	"image"
	"image/color"
	"math"

	"github.com/hoonfeng/svg/types"
)

// StrokeCapStyle 线帽样式 / Stroke cap style
type StrokeCapStyle int

const (
	CapButt   StrokeCapStyle = iota // 平头线帽 / Butt cap
	CapRound                        // 圆形线帽 / Round cap
	CapSquare                       // 方形线帽 / Square cap
)

// StrokeJoinStyle 线段连接样式 / Stroke join style
type StrokeJoinStyle int

const (
	JoinMiter StrokeJoinStyle = iota // 尖角连接 / Miter join
	JoinRound                        // 圆角连接 / Round join
	JoinBevel                        // 斜角连接 / Bevel join
)

// TrueStrokePathGenerator 真正的描边路径生成器 / True stroke path generator
type TrueStrokePathGenerator struct {
	CapStyle   StrokeCapStyle  // 线帽样式 / Cap style
	JoinStyle  StrokeJoinStyle // 连接样式 / Join style
	MiterLimit float64         // 尖角限制 / Miter limit
}

// NewTrueStrokePathGenerator 创建新的描边路径生成器 / Create new stroke path generator
func NewTrueStrokePathGenerator() *TrueStrokePathGenerator {
	return &TrueStrokePathGenerator{
		CapStyle:   CapRound,  // 默认圆形线帽 / Default round cap
		JoinStyle:  JoinRound, // 默认圆角连接 / Default round join
		MiterLimit: 4.0,       // 默认尖角限制 / Default miter limit
	}
}

// GenerateStrokePath 生成真正的描边路径 / Generate true stroke path
func (g *TrueStrokePathGenerator) GenerateStrokePath(path []types.Point, strokeWidth float64, closePath bool) []types.Point {
	if len(path) < 2 {
		return nil
	}

	halfWidth := strokeWidth / 2
	strokePath := make([]types.Point, 0)

	// 处理路径闭合 / Handle path closure
	processedPath := path
	if closePath && len(path) >= 3 {
		// 检查是否需要添加闭合线段 / Check if closure segment is needed
		firstPoint := path[0]
		lastPoint := path[len(path)-1]
		dx := firstPoint.X - lastPoint.X
		dy := firstPoint.Y - lastPoint.Y
		distance := math.Sqrt(dx*dx + dy*dy)
		if distance > 0.1 { // 如果起点和终点距离大于0.1像素，添加闭合线段
			processedPath = make([]types.Point, len(path)+1)
			copy(processedPath, path)
			processedPath[len(path)] = firstPoint
		}
	}

	// 生成左侧和右侧偏移路径 / Generate left and right offset paths
	leftPath := g.generateOffsetPath(processedPath, halfWidth, true)
	rightPath := g.generateOffsetPath(processedPath, halfWidth, false)

	// 构建完整的描边路径 / Build complete stroke path
	if len(leftPath) > 0 && len(rightPath) > 0 {
		// 添加左侧路径 / Add left path
		strokePath = append(strokePath, leftPath...)

		// 如果不是闭合路径，添加终点线帽 / Add end cap if not closed
		if !closePath {
			endCap := g.generateEndCap(processedPath[len(processedPath)-2], processedPath[len(processedPath)-1], halfWidth, false)
			strokePath = append(strokePath, endCap...)
		}

		// 添加右侧路径（反向）/ Add right path (reversed)
		for i := len(rightPath) - 1; i >= 0; i-- {
			strokePath = append(strokePath, rightPath[i])
		}

		// 如果不是闭合路径，添加起点线帽 / Add start cap if not closed
		if !closePath {
			startCap := g.generateEndCap(processedPath[1], processedPath[0], halfWidth, true)
			strokePath = append(strokePath, startCap...)
		}
	}

	return strokePath
}

// generateOffsetPath 生成偏移路径 / Generate offset path
func (g *TrueStrokePathGenerator) generateOffsetPath(path []types.Point, offset float64, isLeft bool) []types.Point {
	if len(path) < 2 {
		return nil
	}

	offsetPath := make([]types.Point, 0)

	for i := 0; i < len(path)-1; i++ {
		current := path[i]
		next := path[i+1]

		// 计算线段的法向量 / Calculate normal vector of segment
		dx := next.X - current.X
		dy := next.Y - current.Y
		length := math.Sqrt(dx*dx + dy*dy)

		if length < 1e-10 {
			continue // 跳过长度为0的线段 / Skip zero-length segments
		}

		// 归一化方向向量 / Normalize direction vector
		dx /= length
		dy /= length

		// 计算法向量（垂直向量）/ Calculate normal vector (perpendicular)
		normalX := -dy
		normalY := dx

		// 根据左右侧调整法向量方向 / Adjust normal direction for left/right
		if !isLeft {
			normalX = -normalX
			normalY = -normalY
		}

		// 计算偏移点 / Calculate offset points
		offsetStart := types.Point{
			X: current.X + normalX*offset,
			Y: current.Y + normalY*offset,
		}
		offsetEnd := types.Point{
			X: next.X + normalX*offset,
			Y: next.Y + normalY*offset,
		}

		// 处理线段连接 / Handle segment joins
		if i == 0 {
			// 第一个线段，直接添加起点 / First segment, add start point directly
			offsetPath = append(offsetPath, offsetStart)
		} else {
			// 处理与前一个线段的连接 / Handle join with previous segment
			joinPoints := g.generateJoin(path[i-1], current, next, offset, isLeft)
			offsetPath = append(offsetPath, joinPoints...)
		}

		// 添加线段终点 / Add segment end point
		if i == len(path)-2 {
			// 最后一个线段，添加终点 / Last segment, add end point
			offsetPath = append(offsetPath, offsetEnd)
		}
	}

	return offsetPath
}

// generateJoin 生成线段连接 / Generate segment join
func (g *TrueStrokePathGenerator) generateJoin(prev, current, next types.Point, offset float64, isLeft bool) []types.Point {
	joinPoints := make([]types.Point, 0)

	// 计算前一个线段的方向 / Calculate previous segment direction
	prevDx := current.X - prev.X
	prevDy := current.Y - prev.Y
	prevLength := math.Sqrt(prevDx*prevDx + prevDy*prevDy)

	// 计算下一个线段的方向 / Calculate next segment direction
	nextDx := next.X - current.X
	nextDy := next.Y - current.Y
	nextLength := math.Sqrt(nextDx*nextDx + nextDy*nextDy)

	if prevLength < 1e-10 || nextLength < 1e-10 {
		return joinPoints // 跳过长度为0的线段 / Skip zero-length segments
	}

	// 归一化方向向量 / Normalize direction vectors
	prevDx /= prevLength
	prevDy /= prevLength
	nextDx /= nextLength
	nextDy /= nextLength

	// 计算法向量 / Calculate normal vectors
	prevNormalX := -prevDy
	prevNormalY := prevDx
	nextNormalX := -nextDy
	nextNormalY := nextDx

	// 根据左右侧调整法向量方向 / Adjust normal direction for left/right
	if !isLeft {
		prevNormalX = -prevNormalX
		prevNormalY = -prevNormalY
		nextNormalX = -nextNormalX
		nextNormalY = -nextNormalY
	}

	// 计算偏移点 / Calculate offset points
	prevOffset := types.Point{
		X: current.X + prevNormalX*offset,
		Y: current.Y + prevNormalY*offset,
	}
	nextOffset := types.Point{
		X: current.X + nextNormalX*offset,
		Y: current.Y + nextNormalY*offset,
	}

	// 根据连接样式生成连接点 / Generate join points based on join style
	switch g.JoinStyle {
	case JoinMiter:
		// 尖角连接 / Miter join
		miterPoint := g.calculateMiterJoin(prevOffset, current, nextOffset, offset)
		if miterPoint != nil {
			joinPoints = append(joinPoints, *miterPoint)
		} else {
			// 尖角过长，回退到斜角连接 / Miter too long, fallback to bevel
			joinPoints = append(joinPoints, prevOffset, nextOffset)
		}
	case JoinRound:
		// 圆角连接 / Round join
		roundPoints := g.generateRoundJoin(prevOffset, current, nextOffset, offset)
		joinPoints = append(joinPoints, roundPoints...)
	case JoinBevel:
		// 斜角连接 / Bevel join
		joinPoints = append(joinPoints, prevOffset, nextOffset)
	}

	return joinPoints
}

// calculateMiterJoin 计算尖角连接 / Calculate miter join
func (g *TrueStrokePathGenerator) calculateMiterJoin(prevOffset, center, nextOffset types.Point, offset float64) *types.Point {
	// 计算两条偏移线的交点 / Calculate intersection of two offset lines
	// 使用线段交点公式 / Use line intersection formula

	// 前一条线的方向向量 / Previous line direction vector
	prevDx := center.X - prevOffset.X
	prevDy := center.Y - prevOffset.Y

	// 下一条线的方向向量 / Next line direction vector
	nextDx := center.X - nextOffset.X
	nextDy := center.Y - nextOffset.Y

	// 计算行列式 / Calculate determinant
	det := prevDx*nextDy - prevDy*nextDx

	if math.Abs(det) < 1e-10 {
		return nil // 线段平行，无交点 / Lines are parallel, no intersection
	}

	// 计算参数 / Calculate parameters
	dx := nextOffset.X - prevOffset.X
	dy := nextOffset.Y - prevOffset.Y
	t := (dx*nextDy - dy*nextDx) / det

	// 计算交点 / Calculate intersection point
	intersection := types.Point{
		X: prevOffset.X + t*prevDx,
		Y: prevOffset.Y + t*prevDy,
	}

	// 检查尖角长度限制 / Check miter length limit
	distance := math.Sqrt((intersection.X-center.X)*(intersection.X-center.X) + (intersection.Y-center.Y)*(intersection.Y-center.Y))
	if distance > offset*g.MiterLimit {
		return nil // 尖角过长 / Miter too long
	}

	return &intersection
}

// generateRoundJoin 生成圆角连接 / Generate round join
func (g *TrueStrokePathGenerator) generateRoundJoin(prevOffset, center, nextOffset types.Point, offset float64) []types.Point {
	roundPoints := make([]types.Point, 0)

	// 计算起始和结束角度 / Calculate start and end angles
	startAngle := math.Atan2(prevOffset.Y-center.Y, prevOffset.X-center.X)
	endAngle := math.Atan2(nextOffset.Y-center.Y, nextOffset.X-center.X)

	// 确保角度差在合理范围内 / Ensure angle difference is reasonable
	angleDiff := endAngle - startAngle
	if angleDiff > math.Pi {
		angleDiff -= 2 * math.Pi
	} else if angleDiff < -math.Pi {
		angleDiff += 2 * math.Pi
	}

	// 计算圆弧分段数 / Calculate arc segments
	segments := int(math.Ceil(math.Abs(angleDiff) / (math.Pi / 8))) // 每22.5度一个分段
	if segments < 2 {
		segments = 2
	}

	// 生成圆弧点 / Generate arc points
	roundPoints = append(roundPoints, prevOffset)
	for i := 1; i < segments; i++ {
		t := float64(i) / float64(segments)
		angle := startAngle + t*angleDiff
		point := types.Point{
			X: center.X + offset*math.Cos(angle),
			Y: center.Y + offset*math.Sin(angle),
		}
		roundPoints = append(roundPoints, point)
	}
	roundPoints = append(roundPoints, nextOffset)

	return roundPoints
}

// generateEndCap 生成线帽 / Generate end cap
func (g *TrueStrokePathGenerator) generateEndCap(prev, end types.Point, offset float64, isStart bool) []types.Point {
	capPoints := make([]types.Point, 0)

	// 计算线段方向 / Calculate segment direction
	dx := end.X - prev.X
	dy := end.Y - prev.Y
	length := math.Sqrt(dx*dx + dy*dy)

	if length < 1e-10 {
		return capPoints // 跳过长度为0的线段 / Skip zero-length segments
	}

	// 归一化方向向量 / Normalize direction vector
	dx /= length
	dy /= length

	// 计算法向量 / Calculate normal vector
	normalX := -dy
	normalY := dx

	// 根据起点/终点调整方向 / Adjust direction for start/end
	if isStart {
		dx = -dx
		dy = -dy
	}

	// 计算线帽的基础点 / Calculate cap base points
	leftPoint := types.Point{
		X: end.X + normalX*offset,
		Y: end.Y + normalY*offset,
	}
	rightPoint := types.Point{
		X: end.X - normalX*offset,
		Y: end.Y - normalY*offset,
	}

	// 根据线帽样式生成线帽 / Generate cap based on cap style
	switch g.CapStyle {
	case CapButt:
		// 平头线帽，直接连接 / Butt cap, direct connection
		capPoints = append(capPoints, leftPoint, rightPoint)
	case CapSquare:
		// 方形线帽 / Square cap
		extendedLeft := types.Point{
			X: leftPoint.X + dx*offset,
			Y: leftPoint.Y + dy*offset,
		}
		extendedRight := types.Point{
			X: rightPoint.X + dx*offset,
			Y: rightPoint.Y + dy*offset,
		}
		capPoints = append(capPoints, leftPoint, extendedLeft, extendedRight, rightPoint)
	case CapRound:
		// 圆形线帽 / Round cap
		roundCap := g.generateRoundCap(end, leftPoint, rightPoint, offset, dx, dy)
		capPoints = append(capPoints, roundCap...)
	}

	return capPoints
}

// generateRoundCap 生成圆形线帽 / Generate round cap
func (g *TrueStrokePathGenerator) generateRoundCap(center, leftPoint, rightPoint types.Point, offset, dx, dy float64) []types.Point {
	roundCap := make([]types.Point, 0)

	// 计算起始和结束角度 / Calculate start and end angles
	startAngle := math.Atan2(leftPoint.Y-center.Y, leftPoint.X-center.X)
	endAngle := math.Atan2(rightPoint.Y-center.Y, rightPoint.X-center.X)

	// 确保角度差为半圆 / Ensure angle difference is semicircle
	angleDiff := endAngle - startAngle
	if angleDiff > 0 {
		angleDiff -= 2 * math.Pi
	}
	if angleDiff > -math.Pi {
		angleDiff -= math.Pi
	}

	// 计算圆弧分段数 / Calculate arc segments
	segments := int(math.Ceil(math.Abs(angleDiff) / (math.Pi / 8))) // 每22.5度一个分段
	if segments < 4 {
		segments = 4
	}

	// 生成圆弧点 / Generate arc points
	roundCap = append(roundCap, leftPoint)
	for i := 1; i < segments; i++ {
		t := float64(i) / float64(segments)
		angle := startAngle + t*angleDiff
		point := types.Point{
			X: center.X + offset*math.Cos(angle),
			Y: center.Y + offset*math.Sin(angle),
		}
		roundCap = append(roundCap, point)
	}
	roundCap = append(roundCap, rightPoint)

	return roundCap
}

// TrueStrokeRenderer 真正的描边渲染器 / True stroke renderer
type TrueStrokeRenderer struct {
	*AntiAliasedPathRenderer
	PathGenerator *TrueStrokePathGenerator
}

// NewTrueStrokeRenderer 创建新的真正描边渲染器 / Create new true stroke renderer
func NewTrueStrokeRenderer() *TrueStrokeRenderer {
	return &TrueStrokeRenderer{
		AntiAliasedPathRenderer: NewAntiAliasedPathRenderer(),
		PathGenerator:           NewTrueStrokePathGenerator(),
	}
}

// RenderTrueStroke 渲染真正的描边 / Render true stroke
func (r *TrueStrokeRenderer) RenderTrueStroke(img *image.RGBA, path []types.Point, strokeColor color.RGBA, strokeWidth float64, closePath bool) {
	if len(path) < 2 {
		return
	}

	// 生成真正的描边路径 / Generate true stroke path
	strokePath := r.PathGenerator.GenerateStrokePath(path, strokeWidth, closePath)
	if len(strokePath) < 3 {
		return
	}

	// 使用专门的描边路径渲染方法 / Use specialized stroke path rendering method
	r.renderStrokePathDirect(img, strokePath, strokeColor)
}

// renderStrokePathDirect 直接渲染描边路径轮廓 / Directly render stroke path outline
func (r *TrueStrokeRenderer) renderStrokePathDirect(img *image.RGBA, strokePath []types.Point, strokeColor color.RGBA) {
	if len(strokePath) < 3 {
		return
	}

	// 计算路径边界 / Calculate path bounds
	bounds := r.calculatePathBounds(strokePath)
	minX := int(math.Floor(bounds.MinX - 2))
	maxX := int(math.Ceil(bounds.MaxX + 2))
	minY := int(math.Floor(bounds.MinY - 2))
	maxY := int(math.Ceil(bounds.MaxY + 2))

	// 确保边界在图像范围内 / Ensure bounds are within image
	if minX < 0 {
		minX = 0
	}
	if maxX >= img.Bounds().Dx() {
		maxX = img.Bounds().Dx() - 1
	}
	if minY < 0 {
		minY = 0
	}
	if maxY >= img.Bounds().Dy() {
		maxY = img.Bounds().Dy() - 1
	}

	// 遍历边界内的每个像素 / Iterate through each pixel in bounds
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			// 计算像素覆盖率 / Calculate pixel coverage
			coverage := r.calculateStrokePathCoverage(float64(x), float64(y), strokePath)
			if coverage > 0 {
				// 混合颜色 / Blend color
				r.blendPixel(img, x, y, strokeColor, coverage)
			}
		}
	}
}

// calculateStrokePathCoverage 计算像素对描边路径的覆盖率 / Calculate pixel coverage for stroke path
func (r *TrueStrokeRenderer) calculateStrokePathCoverage(pixelX, pixelY float64, strokePath []types.Point) float64 {
	// 使用4x4子像素采样 / Use 4x4 sub-pixel sampling
	samples := 4
	insideCount := 0
	totalSamples := samples * samples
	step := 1.0 / float64(samples)

	for i := 0; i < samples; i++ {
		for j := 0; j < samples; j++ {
			sampleX := pixelX + (float64(i)+0.5)*step
			sampleY := pixelY + (float64(j)+0.5)*step

			// 使用射线投射算法检查点是否在描边路径内 / Use ray casting to check if point is inside stroke path
			if r.isPointInStrokePath(sampleX, sampleY, strokePath) {
				insideCount++
			}
		}
	}

	return float64(insideCount) / float64(totalSamples)
}

// isPointInStrokePath 检查点是否在描边路径内 / Check if point is inside stroke path
func (r *TrueStrokeRenderer) isPointInStrokePath(x, y float64, strokePath []types.Point) bool {
	if len(strokePath) < 3 {
		return false
	}

	// 使用射线投射算法 / Use ray casting algorithm
	inside := false
	j := len(strokePath) - 1

	for i := 0; i < len(strokePath); i++ {
		xi, yi := strokePath[i].X, strokePath[i].Y
		xj, yj := strokePath[j].X, strokePath[j].Y

		if ((yi > y) != (yj > y)) && (x < (xj-xi)*(y-yi)/(yj-yi)+xi) {
			inside = !inside
		}
		j = i
	}

	return inside
}

// blendPixel 混合像素颜色 / Blend pixel color
func (r *TrueStrokeRenderer) blendPixel(img *image.RGBA, x, y int, colors color.RGBA, coverage float64) {
	if coverage <= 0 {
		return
	}
	if coverage > 1 {
		coverage = 1
	}

	// 获取当前像素颜色 / Get current pixel color
	currentColor := img.RGBAAt(x, y)

	// Alpha混合 / Alpha blending
	alpha := float64(colors.A) * coverage / 255.0
	invAlpha := 1.0 - alpha

	newR := uint8(float64(colors.R)*alpha + float64(currentColor.R)*invAlpha)
	newG := uint8(float64(colors.G)*alpha + float64(currentColor.G)*invAlpha)
	newB := uint8(float64(colors.B)*alpha + float64(currentColor.B)*invAlpha)
	newA := uint8(math.Min(255, float64(currentColor.A)+alpha*255))

	img.SetRGBA(x, y, color.RGBA{R: newR, G: newG, B: newB, A: newA})
}

// PathBounds 路径边界框结构 / Path bounds structure

// calculatePathBounds 计算路径边界框 / Calculate path bounds
func (r *TrueStrokeRenderer) calculatePathBounds(path []types.Point) PathBounds {
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

// RenderTrueStrokeComplexPath 渲染复杂路径的真正描边 / Render true stroke for complex path
func (r *TrueStrokeRenderer) RenderTrueStrokeComplexPath(img *image.RGBA, subPaths [][]types.Point, strokeColor color.RGBA, strokeWidth float64, closeSubPaths []bool) {
	if len(subPaths) == 0 {
		return
	}

	// 为每个子路径单独处理描边 / Process stroke for each sub-path separately
	for i, subPath := range subPaths {
		if len(subPath) < 2 {
			continue
		}

		// 确定是否闭合当前子路径 / Determine if current sub-path should be closed
		closePath := false
		if i < len(closeSubPaths) {
			closePath = closeSubPaths[i]
		}

		// 渲染当前子路径的真正描边 / Render true stroke for current sub-path
		r.RenderTrueStroke(img, subPath, strokeColor, strokeWidth, closePath)
	}
}
