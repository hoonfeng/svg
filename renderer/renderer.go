package renderer

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"strconv"
	"strings"

	"github.com/hoonfeng/svg/font"
	"github.com/hoonfeng/svg/types"
)

// ImageRenderer 表示SVG到图像的渲染器
type ImageRenderer struct {
}

// NewImageRenderer 创建新的图像渲染器
func NewImageRenderer() *ImageRenderer {
	return &ImageRenderer{}
}

// NewImage 创建新的图像（为了兼容测试）
func NewImage(width, height int) *image.RGBA {
	return CreateImage(width, height, color.RGBA{0, 0, 0, 0})
}

// RenderElement 渲染单个元素（为了兼容测试）
func RenderElement(img *image.RGBA, element types.Element) error {
	renderer := NewImageRenderer()
	viewBox := []float64{0, 0, float64(img.Bounds().Dx()), float64(img.Bounds().Dy())}
	return renderer.renderElement(img, element, viewBox, 1.0, 1.0)
}

// RenderDocument 渲染整个文档（为了兼容测试）
func RenderDocument(doc *types.Document, width, height int) (*image.RGBA, error) {
	renderer := NewImageRenderer()
	return renderer.Render(doc, width, height)
}

// Render 将SVG文档渲染为图像
func (r *ImageRenderer) Render(doc *types.Document, width, height int) (*image.RGBA, error) {
	// 创建图像，使用透明背景 / Create image with transparent background
	img := CreateImage(width, height, color.RGBA{0, 0, 0, 0})

	// 解析视口
	viewBox := parseViewBox(doc.ViewBox)

	// 计算缩放比例
	scaleX := float64(width) / (viewBox[2] - viewBox[0])
	scaleY := float64(height) / (viewBox[3] - viewBox[1])

	// 渲染元素
	for _, element := range doc.Elements {
		err := r.renderElement(img, element, viewBox, scaleX, scaleY)
		if err != nil {
			return nil, err
		}
	}

	return img, nil
}

// renderElement 渲染单个SVG元素
func (r *ImageRenderer) renderElement(img *image.RGBA, element types.Element, viewBox []float64, scaleX, scaleY float64) error {
	switch element.Tag() {
	case "rect":
		return r.renderRect(img, element, viewBox, scaleX, scaleY)
	case "circle":
		return r.renderCircle(img, element, viewBox, scaleX, scaleY)
	case "ellipse":
		return r.renderEllipse(img, element, viewBox, scaleX, scaleY)
	case "line":
		return r.renderLine(img, element, viewBox, scaleX, scaleY)
	case "polyline":
		return r.renderPolyline(img, element, viewBox, scaleX, scaleY)
	case "polygon":
		return r.renderPolygon(img, element, viewBox, scaleX, scaleY)
	case "path":
		return r.renderPath(img, element, viewBox, scaleX, scaleY)
	case "text":
		return r.renderText(img, element, viewBox, scaleX, scaleY)
	case "g":
		// 组元素的渲染需要解析内容中的子元素
		// 简化实现，暂不支持组元素
		return nil
	default:
		return fmt.Errorf("不支持的元素类型: %s", element.Tag())
	}
}

// renderRect 渲染矩形元素
func (r *ImageRenderer) renderRect(img *image.RGBA, element types.Element, viewBox []float64, scaleX, scaleY float64) error {
	attrs := element.GetAttributes()

	// 解析属性
	x, _ := parseFloat(attrs["x"], 0)
	y, _ := parseFloat(attrs["y"], 0)
	width, _ := parseFloat(attrs["width"], 0)
	height, _ := parseFloat(attrs["height"], 0)

	// 转换坐标
	x1 := int((x - viewBox[0]) * scaleX)
	y1 := int((y - viewBox[1]) * scaleY)
	w := int(width * scaleX)
	h := int(height * scaleY)

	// 解析颜色
	fillColor := parseColor(attrs["fill"], color.RGBA{0, 0, 0, 0})
	strokeColor := parseColor(attrs["stroke"], color.RGBA{0, 0, 0, 255})

	// 判断是填充还是描边 / Determine if fill or stroke
	hasFill := attrs["fill"] != "none" && attrs["fill"] != ""
	hasStroke := attrs["stroke"] != "none" && attrs["stroke"] != ""

	// 绘制矩形
	if hasFill && fillColor != (color.RGBA{0, 0, 0, 0}) {
		DrawRect(img, x1, y1, w, h, fillColor, true)
	}

	if hasStroke && strokeColor != (color.RGBA{0, 0, 0, 0}) {
		DrawRect(img, x1, y1, w, h, strokeColor, false)
	}

	// 如果既没有填充也没有描边，默认使用填充 / Default to fill if neither fill nor stroke
	if !hasFill && !hasStroke {
		DrawRect(img, x1, y1, w, h, color.RGBA{0, 0, 0, 255}, true)
	}

	return nil
}

// renderCircle 渲染圆形元素
func (r *ImageRenderer) renderCircle(img *image.RGBA, element types.Element, viewBox []float64, scaleX, scaleY float64) error {
	attrs := element.GetAttributes()

	// 解析属性
	cx, _ := parseFloat(attrs["cx"], 0)
	cy, _ := parseFloat(attrs["cy"], 0)
	radius, _ := parseFloat(attrs["r"], 0)

	// 转换坐标
	centerX := int((cx - viewBox[0]) * scaleX)
	centerY := int((cy - viewBox[1]) * scaleY)
	circleRadius := int(radius * ((scaleX + scaleY) / 2))

	// 解析颜色
	fillColor := parseColor(attrs["fill"], color.RGBA{0, 0, 0, 0})
	strokeColor := parseColor(attrs["stroke"], color.RGBA{0, 0, 0, 255})

	// 判断是填充还是描边 / Determine if fill or stroke
	hasFill := attrs["fill"] != "none" && attrs["fill"] != ""
	hasStroke := attrs["stroke"] != "none" && attrs["stroke"] != ""

	// 绘制圆形
	if hasFill && fillColor != (color.RGBA{0, 0, 0, 0}) {
		DrawCircle(img, centerX, centerY, circleRadius, fillColor, true)
	}

	if hasStroke && strokeColor != (color.RGBA{0, 0, 0, 0}) {
		DrawCircle(img, centerX, centerY, circleRadius, strokeColor, false)
	}

	// 如果既没有填充也没有描边，默认使用填充 / Default to fill if neither fill nor stroke
	if !hasFill && !hasStroke {
		DrawCircle(img, centerX, centerY, circleRadius, color.RGBA{0, 0, 0, 255}, true)
	}

	return nil
}

// renderEllipse 渲染椭圆元素
func (r *ImageRenderer) renderEllipse(img *image.RGBA, element types.Element, viewBox []float64, scaleX, scaleY float64) error {
	attrs := element.GetAttributes()

	// 解析属性
	cx, _ := parseFloat(attrs["cx"], 0)
	cy, _ := parseFloat(attrs["cy"], 0)
	rx, _ := parseFloat(attrs["rx"], 0)
	ry, _ := parseFloat(attrs["ry"], 0)

	// 转换坐标
	centerX := int((cx - viewBox[0]) * scaleX)
	centerY := int((cy - viewBox[1]) * scaleY)
	radiusX := int(rx * scaleX)
	radiusY := int(ry * scaleY)

	// 解析颜色
	fillColor := parseColor(attrs["fill"], color.RGBA{0, 0, 0, 0})
	strokeColor := parseColor(attrs["stroke"], color.RGBA{0, 0, 0, 255})

	// 判断是填充还是描边 / Determine if fill or stroke
	hasFill := attrs["fill"] != "none" && attrs["fill"] != ""
	hasStroke := attrs["stroke"] != "none" && attrs["stroke"] != ""

	// 绘制椭圆
	if hasFill && fillColor != (color.RGBA{0, 0, 0, 0}) {
		DrawEllipse(img, centerX, centerY, radiusX, radiusY, fillColor, true)
	}

	if hasStroke && strokeColor != (color.RGBA{0, 0, 0, 0}) {
		DrawEllipse(img, centerX, centerY, radiusX, radiusY, strokeColor, false)
	}

	// 如果既没有填充也没有描边，默认使用填充 / Default to fill if neither fill nor stroke
	if !hasFill && !hasStroke {
		DrawEllipse(img, centerX, centerY, radiusX, radiusY, color.RGBA{0, 0, 0, 255}, true)
	}

	return nil
}

// renderLine 渲染线段元素
func (r *ImageRenderer) renderLine(img *image.RGBA, element types.Element, viewBox []float64, scaleX, scaleY float64) error {
	attrs := element.GetAttributes()

	// 解析属性
	x1, _ := parseFloat(attrs["x1"], 0)
	y1, _ := parseFloat(attrs["y1"], 0)
	x2, _ := parseFloat(attrs["x2"], 0)
	y2, _ := parseFloat(attrs["y2"], 0)

	// 转换坐标
	px1 := int((x1 - viewBox[0]) * scaleX)
	py1 := int((y1 - viewBox[1]) * scaleY)
	px2 := int((x2 - viewBox[0]) * scaleX)
	py2 := int((y2 - viewBox[1]) * scaleY)

	// 解析颜色
	strokeColor := parseColor(attrs["stroke"], color.RGBA{0, 0, 0, 255})

	// 绘制线段
	DrawLine(img, px1, py1, px2, py2, strokeColor)

	return nil
}

// renderPolyline 渲染折线元素
func (r *ImageRenderer) renderPolyline(img *image.RGBA, element types.Element, viewBox []float64, scaleX, scaleY float64) error {
	attrs := element.GetAttributes()

	// 解析属性
	pointsStr := attrs["points"]
	points := parsePoints(pointsStr)

	// 解析颜色
	strokeColor := parseColor(attrs["stroke"], color.RGBA{0, 0, 0, 255})

	// 绘制折线
	for i := 1; i < len(points); i++ {
		x1 := int((points[i-1].X - viewBox[0]) * scaleX)
		y1 := int((points[i-1].Y - viewBox[1]) * scaleY)
		x2 := int((points[i].X - viewBox[0]) * scaleX)
		y2 := int((points[i].Y - viewBox[1]) * scaleY)

		DrawLine(img, x1, y1, x2, y2, strokeColor)
	}

	return nil
}

// renderPolygon 渲染多边形元素
func (r *ImageRenderer) renderPolygon(img *image.RGBA, element types.Element, viewBox []float64, scaleX, scaleY float64) error {
	attrs := element.GetAttributes()

	// 解析属性
	pointsStr := attrs["points"]
	points := parsePoints(pointsStr)

	// 解析颜色
	strokeColor := parseColor(attrs["stroke"], color.RGBA{0, 0, 0, 255})

	// 绘制多边形
	for i := 1; i < len(points); i++ {
		x1 := int((points[i-1].X - viewBox[0]) * scaleX)
		y1 := int((points[i-1].Y - viewBox[1]) * scaleY)
		x2 := int((points[i].X - viewBox[0]) * scaleX)
		y2 := int((points[i].Y - viewBox[1]) * scaleY)

		DrawLine(img, x1, y1, x2, y2, strokeColor)
	}

	// 闭合多边形
	if len(points) > 1 {
		x1 := int((points[len(points)-1].X - viewBox[0]) * scaleX)
		y1 := int((points[len(points)-1].Y - viewBox[1]) * scaleY)
		x2 := int((points[0].X - viewBox[0]) * scaleX)
		y2 := int((points[0].Y - viewBox[1]) * scaleY)

		DrawLine(img, x1, y1, x2, y2, strokeColor)
	}

	return nil
}

// renderPath 渲染路径元素（使用抗锯齿） / Render path element (with anti-aliasing)
func (r *ImageRenderer) renderPath(img *image.RGBA, element types.Element, viewBox []float64, scaleX, scaleY float64) error {
	attrs := element.GetAttributes()

	// 获取路径数据 / Get path data
	pathData, exists := attrs["d"]
	if !exists {
		return fmt.Errorf("路径元素缺少'd'属性")
	}

	// 获取样式 / Get styles
	fillColor := r.getFillColor(attrs)
	strokeColor := r.getStrokeColor(attrs)
	strokeWidth := r.getStrokeWidth(attrs)

	// 创建抗锯齿路径渲染器 / Create anti-aliased path renderer
	aaPathRenderer := NewAntiAliasedPathRenderer()

	// 使用抗锯齿路径渲染器渲染路径 / Render path using anti-aliased path renderer
	return aaPathRenderer.RenderPath(img, pathData, fillColor, strokeColor, strokeWidth, viewBox, scaleX, scaleY)
}

// renderText 渲染文本元素
func (r *ImageRenderer) renderText(img *image.RGBA, element types.Element, viewBox []float64, scaleX, scaleY float64) error {
	attrs := element.GetAttributes()

	// 解析位置属性
	x, _ := parseFloat(attrs["x"], 0)
	y, _ := parseFloat(attrs["y"], 0)

	// 转换坐标
	renderX := (x - viewBox[0]) * scaleX
	renderY := (y - viewBox[1]) * scaleY

	// 获取文本内容
	var textContent string
	if textElement, ok := element.(interface{ GetContent() string }); ok {
		textContent = textElement.GetContent()
	} else {
		return fmt.Errorf("无法获取文本内容")
	}

	if textContent == "" {
		return nil // 空文本不需要渲染
	}

	// 创建文本样式
	style := r.createTextStyleFromAttributes(attrs, scaleX, scaleY)

	// 使用SVG文本渲染器渲染文本
	textRenderer := font.DefaultTextRenderer
	return textRenderer.RenderText(img, textContent, renderX, renderY, style)
}

// createTextStyleFromAttributes 从SVG属性创建文本样式
func (r *ImageRenderer) createTextStyleFromAttributes(attrs map[string]string, scaleX, scaleY float64) *font.TextStyle {
	style := font.NewTextStyle()

	// 解析字体族
	if fontFamily, ok := attrs["font-family"]; ok {
		style.FontFamily = fontFamily
	}

	// 解析字体大小
	if fontSizeStr, ok := attrs["font-size"]; ok {
		if fontSize, err := parseFloat(fontSizeStr, 16); err == nil {
			// 应用缩放
			style.FontSize = fontSize * ((scaleX + scaleY) / 2)
		}
	}

	// 解析字体粗细
	if fontWeight, ok := attrs["font-weight"]; ok {
		style.FontWeight = parseFontWeight(fontWeight)
	}

	// 解析字体样式
	if fontStyle, ok := attrs["font-style"]; ok {
		style.FontStyle = parseFontStyle(fontStyle)
	}

	// 解析文本锚点
	if textAnchor, ok := attrs["text-anchor"]; ok {
		switch textAnchor {
		case "start":
			style.TextAnchor = font.TextAnchorStart
		case "middle":
			style.TextAnchor = font.TextAnchorMiddle
		case "end":
			style.TextAnchor = font.TextAnchorEnd
		}
	}

	// 解析基线对齐
	if alignmentBaseline, ok := attrs["alignment-baseline"]; ok {
		switch alignmentBaseline {
		case "alphabetic":
			style.AlignmentBaseline = font.AlignmentBaselineAlphabetic
		case "middle":
			style.AlignmentBaseline = font.AlignmentBaselineMiddle
		case "hanging":
			style.AlignmentBaseline = font.AlignmentBaselineHanging
		case "top":
			style.AlignmentBaseline = font.AlignmentBaselineTop
		case "bottom":
			style.AlignmentBaseline = font.AlignmentBaselineBottom
		}
	}

	// 解析填充颜色
	if fill, ok := attrs["fill"]; ok {
		fillColor := parseColor(fill, color.RGBA{0, 0, 0, 255})
		style.Fill = &image.Uniform{C: fillColor}
	}

	// 解析描边颜色
	if stroke, ok := attrs["stroke"]; ok && stroke != "none" {
		strokeColor := parseColor(stroke, color.RGBA{0, 0, 0, 255})
		style.Stroke = &image.Uniform{C: strokeColor}
	}

	// 解析描边宽度
	if strokeWidthStr, ok := attrs["stroke-width"]; ok {
		if strokeWidth, err := parseFloat(strokeWidthStr, 0); err == nil {
			// 应用缩放
			style.StrokeWidth = strokeWidth * ((scaleX + scaleY) / 2)
		}
	}

	return style
}

// parseViewBox 解析视口
func parseViewBox(viewBox string) []float64 {
	// 如果viewBox为空，返回默认值
	if viewBox == "" {
		return []float64{0, 0, 800, 600} // 默认viewBox
	}

	parts := strings.Fields(viewBox)
	result := make([]float64, 4)

	for i := 0; i < len(parts) && i < 4; i++ {
		value, _ := strconv.ParseFloat(parts[i], 64)
		result[i] = value
	}

	// 确保宽度和高度不为零
	if result[2] == 0 {
		result[2] = 800 // 默认宽度
	}
	if result[3] == 0 {
		result[3] = 600 // 默认高度
	}

	return result
}

// parseFloat 解析浮点数
func parseFloat(s string, defaultValue float64) (float64, error) {
	if s == "" {
		return defaultValue, nil
	}

	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return defaultValue, err
	}

	return value, nil
}

// parseColor 解析颜色
func parseColor(s string, defaultColor color.RGBA) color.RGBA {
	if s == "" || s == "none" {
		return defaultColor
	}

	// 处理十六进制颜色
	if strings.HasPrefix(s, "#") {
		hex := s[1:]

		// 处理简写形式 (#RGB)
		if len(hex) == 3 {
			r := parseHex(hex[0:1] + hex[0:1])
			g := parseHex(hex[1:2] + hex[1:2])
			b := parseHex(hex[2:3] + hex[2:3])
			return color.RGBA{r, g, b, 255}
		}

		// 处理完整形式 (#RRGGBB)
		if len(hex) == 6 {
			r := parseHex(hex[0:2])
			g := parseHex(hex[2:4])
			b := parseHex(hex[4:6])
			return color.RGBA{r, g, b, 255}
		}
	}

	// 处理命名颜色
	switch s {
	case "black":
		return color.RGBA{0, 0, 0, 255}
	case "white":
		return color.RGBA{255, 255, 255, 255}
	case "red":
		return color.RGBA{255, 0, 0, 255}
	case "green":
		return color.RGBA{0, 128, 0, 255}
	case "blue":
		return color.RGBA{0, 0, 255, 255}
	case "yellow":
		return color.RGBA{255, 255, 0, 255}
	case "cyan":
		return color.RGBA{0, 255, 255, 255}
	case "magenta":
		return color.RGBA{255, 0, 255, 255}
	case "gray":
		return color.RGBA{128, 128, 128, 255}
	}

	return defaultColor
}

// parseHex 解析十六进制颜色值
func parseHex(s string) uint8 {
	value, _ := strconv.ParseUint(s, 16, 8)
	return uint8(value)
}

// parseFontWeight 解析字体粗细 / Parse font weight
func parseFontWeight(s string) font.FontWeight {
	switch s {
	case "100":
		return font.FontWeight100
	case "200":
		return font.FontWeight200
	case "300", "light":
		return font.FontWeight300
	case "400", "normal":
		return font.FontWeight400
	case "500", "medium":
		return font.FontWeight500
	case "600", "semibold":
		return font.FontWeight600
	case "700", "bold":
		return font.FontWeight700
	case "800":
		return font.FontWeight800
	case "900", "black":
		return font.FontWeight900
	case "lighter":
		return font.FontWeightLighter
	case "bolder":
		return font.FontWeightBolder
	default:
		return font.FontWeightNormal
	}
}

// parseFontStyle 解析字体样式 / Parse font style
func parseFontStyle(s string) font.FontStyle {
	switch s {
	case "italic":
		return font.FontStyleItalic
	case "oblique":
		return font.FontStyleOblique
	case "normal":
		return font.FontStyleNormal
	default:
		return font.FontStyleNormal
	}
}

// parsePoints 解析点列表
func parsePoints(s string) []types.Point {
	parts := strings.Fields(strings.Replace(s, ",", " ", -1))
	points := []types.Point{}

	for i := 0; i < len(parts)-1; i += 2 {
		x, errX := strconv.ParseFloat(parts[i], 64)
		y, errY := strconv.ParseFloat(parts[i+1], 64)

		if errX == nil && errY == nil {
			points = append(points, types.Point{X: x, Y: y})
		}
	}

	return points
}

// validateAndFixPath 验证并修复路径，确保路径正确闭合
func (r *ImageRenderer) validateAndFixPath(points []types.Point) []types.Point {
	if len(points) < 3 {
		return points
	}

	// 移除重复的相邻点
	cleanPoints := make([]types.Point, 0, len(points))
	cleanPoints = append(cleanPoints, points[0])

	for i := 1; i < len(points); i++ {
		last := cleanPoints[len(cleanPoints)-1]
		current := points[i]

		// 检查是否为重复点（距离小于0.1像素）
		dx := current.X - last.X
		dy := current.Y - last.Y
		distSq := dx*dx + dy*dy

		if distSq > 0.01 { // 0.1像素的平方
			cleanPoints = append(cleanPoints, current)
		}
	}

	// 确保路径闭合
	if len(cleanPoints) >= 3 {
		first := cleanPoints[0]
		last := cleanPoints[len(cleanPoints)-1]

		dx := last.X - first.X
		dy := last.Y - first.Y
		distSq := dx*dx + dy*dy

		// 如果路径未闭合（距离大于1像素），添加闭合点
		if distSq > 1.0 {
			cleanPoints = append(cleanPoints, first)
		}
	}

	return cleanPoints
}

// EdgeInfo 存储边的信息用于缠绕规则计算 / Edge information for winding rule calculation
type EdgeInfo struct {
	X1, Y1, X2, Y2 float64 // 边的起点和终点 / Edge start and end points
	Direction      int     // 边的方向：1表示向上，-1表示向下 / Edge direction: 1 for up, -1 for down
}

// IntersectionInfo 交点信息 / Intersection information
type IntersectionInfo struct {
	X         float64 // 交点的x坐标 / X coordinate of intersection
	Direction int     // 边的方向 / Edge direction
}

// FillPath 公开的填充路径方法 / Public fill path method
func (r *ImageRenderer) FillPath(img *image.RGBA, points []types.Point, fillColor color.RGBA) {
	r.fillPathWithWindingRule(img, points, fillColor)
}

// FillSubPathsWithWindingRule 公开的填充多个子路径方法 / Public fill multiple sub-paths method
func (r *ImageRenderer) FillSubPathsWithWindingRule(img *image.RGBA, subPaths [][]types.Point, fillColor color.RGBA) {
	r.fillSubPathsWithWindingRule(img, subPaths, fillColor)
}

// fillPath 填充路径 / Fill path using high-precision scanline algorithm with anti-aliasing
func (r *ImageRenderer) fillPath(img *image.RGBA, points []types.Point, fillColor color.RGBA) {
	r.fillPathWithWindingRule(img, points, fillColor)
}

// fillPathWithWindingRule 使用非零缠绕规则填充路径 / Fill path using nonzero winding rule
func (r *ImageRenderer) fillPathWithWindingRule(img *image.RGBA, points []types.Point, fillColor color.RGBA) {
	if len(points) < 3 {
		return
	}

	// 验证并修复路径
	points = r.validateAndFixPath(points)
	if len(points) < 3 {
		return
	}

	// 构建边表
	edges := r.buildEdgeTable(points)
	if len(edges) == 0 {
		return
	}

	// 找到边界框
	_, _, minY, maxY := r.getBoundingBox(points)

	// 扫描线填充
	startY := int(math.Floor(minY))
	endY := int(math.Ceil(maxY))

	for y := startY; y <= endY; y++ {
		// 找到与当前扫描线相交的边
		intersections := r.findIntersections(edges, float64(y))
		if len(intersections) == 0 {
			continue
		}

		// 按x坐标排序交点
		r.sortIntersections(intersections)

		// 使用非零缠绕规则填充
		r.fillScanlineWithWinding(img, intersections, y, fillColor)
	}
}

// fillSubPathsWithWindingRule 使用缠绕规则填充多个子路径
// 每个子路径独立处理，避免跨子路径的连接线问题
func (r *ImageRenderer) fillSubPathsWithWindingRule(img *image.RGBA, subPaths [][]types.Point, fillColor color.RGBA) {
	if len(subPaths) == 0 {
		return
	}

	// 收集所有子路径的边（每个子路径独立构建边表）
	allEdges := []EdgeInfo{}
	allPoints := []types.Point{} // 只用于边界框计算，不用于边构建

	for _, subPath := range subPaths {
		if len(subPath) < 3 {
			continue
		}

		// 验证并修复子路径
		subPath = r.validateAndFixPath(subPath)
		if len(subPath) < 3 {
			continue
		}

		// 为每个子路径独立构建边表，避免跨子路径连接
		subEdges := r.buildEdgeTable(subPath)
		allEdges = append(allEdges, subEdges...)

		// 只收集点用于边界框计算
		allPoints = append(allPoints, subPath...)
	}

	if len(allEdges) == 0 || len(allPoints) == 0 {
		return
	}

	// 找到所有点的边界框
	_, _, minY, maxY := r.getBoundingBox(allPoints)

	// 扫描线填充
	startY := int(math.Floor(minY))
	endY := int(math.Ceil(maxY))

	for y := startY; y <= endY; y++ {
		// 找到与当前扫描线相交的边
		intersections := r.findIntersections(allEdges, float64(y))
		if len(intersections) == 0 {
			continue
		}

		// 按x坐标排序交点
		r.sortIntersections(intersections)

		// 使用非零缠绕规则填充
		r.fillScanlineWithWinding(img, intersections, y, fillColor)
	}
}

// buildEdgeTable 构建边表 / Build edge table
func (r *ImageRenderer) buildEdgeTable(points []types.Point) []EdgeInfo {
	edges := make([]EdgeInfo, 0, len(points))

	for i := 0; i < len(points); i++ {
		j := (i + 1) % len(points)
		p1, p2 := points[i], points[j]

		// 跳过水平边
		if math.Abs(p1.Y-p2.Y) < 1e-10 {
			continue
		}

		// 确定边的方向
		direction := 1
		if p1.Y > p2.Y {
			direction = -1
			// 交换点，确保p1在下方，p2在上方
			p1, p2 = p2, p1
		}

		edges = append(edges, EdgeInfo{
			X1:        p1.X,
			Y1:        p1.Y,
			X2:        p2.X,
			Y2:        p2.Y,
			Direction: direction,
		})
	}

	return edges
}

// getBoundingBox 获取边界框 / Get bounding box
func (r *ImageRenderer) getBoundingBox(points []types.Point) (minX, maxX, minY, maxY float64) {
	minX, maxX = points[0].X, points[0].X
	minY, maxY = points[0].Y, points[0].Y

	for _, p := range points {
		if p.X < minX {
			minX = p.X
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}

	return
}

// findIntersections 找到与扫描线的交点 / Find intersections with scanline
func (r *ImageRenderer) findIntersections(edges []EdgeInfo, y float64) []IntersectionInfo {
	intersections := make([]IntersectionInfo, 0)

	for _, edge := range edges {
		// 检查边是否与扫描线相交
		if y >= edge.Y1 && y < edge.Y2 {
			// 计算交点的x坐标
			x := edge.X1 + (y-edge.Y1)*(edge.X2-edge.X1)/(edge.Y2-edge.Y1)
			intersections = append(intersections, IntersectionInfo{
				X:         x,
				Direction: edge.Direction,
			})
		}
	}

	return intersections
}

// sortIntersections 按x坐标排序交点 / Sort intersections by x coordinate
func (r *ImageRenderer) sortIntersections(intersections []IntersectionInfo) {
	for i := 0; i < len(intersections)-1; i++ {
		for j := i + 1; j < len(intersections); j++ {
			if intersections[i].X > intersections[j].X {
				intersections[i], intersections[j] = intersections[j], intersections[i]
			}
		}
	}
}

// fillScanlineWithWinding 使用缠绕数规则填充扫描线 / Fill scanline using winding number rule
func (r *ImageRenderer) fillScanlineWithWinding(img *image.RGBA, intersections []IntersectionInfo, y int, fillColor color.RGBA) {
	if len(intersections) == 0 {
		return
	}

	windingNumber := 0
	lastX := int(math.Floor(intersections[0].X))

	for i, intersection := range intersections {
		currentX := int(math.Floor(intersection.X))

		// 如果缠绕数非零，填充从lastX到currentX的像素
		if windingNumber != 0 {
			for x := lastX; x < currentX; x++ {
				DrawPixel(img, x, y, fillColor)
			}
		}

		// 更新缠绕数
		windingNumber += intersection.Direction

		// 更新lastX为下一个区间的起点
		if i < len(intersections)-1 {
			lastX = currentX
		}
	}
}

// fillPathLegacy 原始填充路径方法（保留作为备用）/ Legacy fill path method (kept as backup)
func (r *ImageRenderer) fillPathLegacy(img *image.RGBA, points []types.Point, fillColor color.RGBA) {
	if len(points) < 3 {
		return // 至少需要3个点才能形成一个可填充的区域 / At least 3 points needed for fillable area
	}

	// 验证并修复路径
	points = r.validateAndFixPath(points)
	if len(points) < 3 {
		return
	}

	// 高精度扫描线填充算法 / High-precision scanline fill algorithm
	// 找到边界框 / Find bounding box
	minX, maxX := points[0].X, points[0].X
	minY, maxY := points[0].Y, points[0].Y

	for _, p := range points {
		if p.X < minX {
			minX = p.X
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}

	// 扫描线填充算法 / Scanline fill algorithm
	startY := math.Floor(minY)
	endY := math.Ceil(maxY)

	for y := startY; y <= endY; y++ {
		intersections := []float64{}

		// 找到与扫描线的交点 / Find intersections with scanline
		for i := 0; i < len(points); i++ {
			j := (i + 1) % len(points)
			p1, p2 := points[i], points[j]

			// 跳过水平边和重复点 / Skip horizontal edges and duplicate points
			if math.Abs(p1.Y-p2.Y) < 1e-10 {
				continue
			}

			// 使用更严格的边界检查，避免端点重复计算
			// 包含下边界和上边界，确保完整填充
			minEdgeY := math.Min(p1.Y, p2.Y)
			maxEdgeY := math.Max(p1.Y, p2.Y)

			if float64(y) >= minEdgeY && float64(y) <= maxEdgeY {
				// 计算交点的x坐标 / Calculate intersection x coordinate
				x := p1.X + (float64(y)-p1.Y)*(p2.X-p1.X)/(p2.Y-p1.Y)
				intersections = append(intersections, x)
			}
		}

		// 对交点进行排序 / Sort intersections
		for i := 0; i < len(intersections)-1; i++ {
			for j := i + 1; j < len(intersections); j++ {
				if intersections[i] > intersections[j] {
					intersections[i], intersections[j] = intersections[j], intersections[i]
				}
			}
		}

		// 填充交点之间的像素 / Fill pixels between intersections
		for i := 0; i < len(intersections)-1; i += 2 {
			startX := int(math.Round(intersections[i]))
			endX := int(math.Round(intersections[i+1]))
			for x := startX; x <= endX; x++ {
				DrawPixel(img, x, int(y), fillColor)
			}
		}
	}
}

// drawAntiAliasedPixel 绘制抗锯齿像素 / Draw anti-aliased pixel
func (r *ImageRenderer) drawAntiAliasedPixel(img *image.RGBA, x, y int, fillColor color.RGBA, alpha float64) {
	if x < 0 || y < 0 || x >= img.Bounds().Dx() || y >= img.Bounds().Dy() {
		return
	}

	// 获取当前像素颜色 / Get current pixel color
	currentColor := img.RGBAAt(x, y)

	// Alpha混合 / Alpha blending
	newR := uint8(float64(currentColor.R)*(1.0-alpha) + float64(fillColor.R)*alpha)
	newG := uint8(float64(currentColor.G)*(1.0-alpha) + float64(fillColor.G)*alpha)
	newB := uint8(float64(currentColor.B)*(1.0-alpha) + float64(fillColor.B)*alpha)
	newA := uint8(math.Max(float64(currentColor.A), float64(fillColor.A)*alpha))

	img.SetRGBA(x, y, color.RGBA{newR, newG, newB, newA})
}

// strokePath 描边路径
func (r *ImageRenderer) strokePath(img *image.RGBA, points []types.Point, strokeColor color.RGBA, strokeWidth float64) {
	if len(points) < 2 {
		return // 至少需要2个点才能绘制线条
	}

	// 如果线宽小于等于1，使用简单的单像素线条
	if strokeWidth <= 1.0 {
		for i := 1; i < len(points); i++ {
			x1 := int(points[i-1].X)
			y1 := int(points[i-1].Y)
			x2 := int(points[i].X)
			y2 := int(points[i].Y)
			DrawLine(img, x1, y1, x2, y2, strokeColor)
		}
		return
	}

	// 对于较粗的线条，使用多条平行线来模拟
	halfWidth := strokeWidth / 2.0
	for i := 1; i < len(points); i++ {
		x1, y1 := points[i-1].X, points[i-1].Y
		x2, y2 := points[i].X, points[i].Y

		// 计算线段的方向向量
		dx := x2 - x1
		dy := y2 - y1
		length := math.Sqrt(dx*dx + dy*dy)

		if length == 0 {
			continue
		}

		// 计算垂直于线段的单位向量
		perpX := -dy / length
		perpY := dx / length

		// 绘制多条平行线来模拟粗线
		steps := int(math.Ceil(strokeWidth))
		for step := -steps / 2; step <= steps/2; step++ {
			offset := float64(step) * (strokeWidth / float64(steps))
			if math.Abs(offset) > halfWidth {
				continue
			}

			startX := int(x1 + perpX*offset)
			startY := int(y1 + perpY*offset)
			endX := int(x2 + perpX*offset)
			endY := int(y2 + perpY*offset)

			DrawLine(img, startX, startY, endX, endY, strokeColor)
		}
	}
}

// getFillColor 获取填充颜色 / Get fill color
func (r *ImageRenderer) getFillColor(attrs map[string]string) color.RGBA {
	fillAttr := attrs["fill"]
	if fillAttr == "none" {
		return color.RGBA{0, 0, 0, 0} // 透明 / Transparent
	}
	if fillAttr == "" {
		// SVG标准：如果没有设置fill属性，默认为黑色 / SVG standard: default to black if no fill attribute
		return color.RGBA{0, 0, 0, 255} // 默认黑色 / Default black
	}
	return parseColor(fillAttr, color.RGBA{0, 0, 0, 255})
}

// getStrokeColor 获取描边颜色
func (r *ImageRenderer) getStrokeColor(attrs map[string]string) color.RGBA {
	strokeAttr := attrs["stroke"]
	if strokeAttr == "none" || strokeAttr == "" {
		return color.RGBA{0, 0, 0, 0} // 透明
	}
	return parseColor(strokeAttr, color.RGBA{0, 0, 0, 255})
}

// getStrokeWidth 获取描边宽度
func (r *ImageRenderer) getStrokeWidth(attrs map[string]string) float64 {
	strokeWidth, _ := parseFloat(attrs["stroke-width"], 1)
	return strokeWidth
}
