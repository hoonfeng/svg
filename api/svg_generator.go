// Package api provides SVG generation APIs
// api包提供SVG生成API
package api

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hoonfeng/svg/types"
)

// SVGGenerator SVG生成器 / SVG generator
type SVGGenerator struct {
	builder *SVGBuilder
}

// NewSVGGenerator 创建新的SVG生成器 / Create new SVG generator
func NewSVGGenerator(width, height float64) *SVGGenerator {
	return &SVGGenerator{
		builder: NewSVGBuilder(width, height),
	}
}

// GetBuilder 获取构建器 / Get builder
func (g *SVGGenerator) GetBuilder() *SVGBuilder {
	return g.builder
}

// GetDocument 获取文档 / Get document
func (g *SVGGenerator) GetDocument() *types.Document {
	return g.builder.GetDocument()
}

// CreateChart 创建图表 / Create chart
func (g *SVGGenerator) CreateChart(chartType string, data []float64, options ChartOptions) *SVGGenerator {
	switch chartType {
	case "bar":
		g.createBarChart(data, options)
	case "line":
		g.createLineChart(data, options)
	case "pie":
		g.createPieChart(data, options)
	default:
		g.createBarChart(data, options)
	}
	return g
}

// CreateGrid 创建网格 / Create grid
func (g *SVGGenerator) CreateGrid(rows, cols int, cellWidth, cellHeight float64, options GridOptions) *SVGGenerator {
	width := float64(cols) * cellWidth
	height := float64(rows) * cellHeight

	// 创建网格组 / Create grid group
	gridGroup := g.builder.BeginGroup()

	// 绘制垂直线 / Draw vertical lines
	for i := 0; i <= cols; i++ {
		x := float64(i) * cellWidth
		g.builder.AddLine(x, 0, x, height).
			Stroke(options.LineColor).
			StrokeWidth(options.LineWidth).
			End()
	}

	// 绘制水平线 / Draw horizontal lines
	for i := 0; i <= rows; i++ {
		y := float64(i) * cellHeight
		g.builder.AddLine(0, y, width, y).
			Stroke(options.LineColor).
			StrokeWidth(options.LineWidth).
			End()
	}

	gridGroup.End()
	return g
}

// CreatePattern 创建图案 / Create pattern
func (g *SVGGenerator) CreatePattern(patternType string, options PatternOptions) *SVGGenerator {
	switch patternType {
	case "dots":
		g.createDotPattern(options)
	case "stripes":
		g.createStripePattern(options)
	case "checkerboard":
		g.createCheckerboardPattern(options)
	default:
		g.createDotPattern(options)
	}
	return g
}

// CreateShape 创建复杂形状 / Create complex shape
func (g *SVGGenerator) CreateShape(shapeType string, options ShapeOptions) *SVGGenerator {
	switch shapeType {
	case "star":
		g.createStar(options)
	case "polygon":
		g.createPolygon(options)
	case "arrow":
		g.createArrow(options)
	case "heart":
		g.createHeart(options)
	default:
		g.createStar(options)
	}
	return g
}

// createBarChart 创建柱状图 / Create bar chart
func (g *SVGGenerator) createBarChart(data []float64, options ChartOptions) {
	if len(data) == 0 {
		return
	}

	// 计算最大值 / Calculate max value
	maxValue := data[0]
	for _, v := range data {
		if v > maxValue {
			maxValue = v
		}
	}

	// 计算柱子宽度 / Calculate bar width
	barWidth := options.Width / float64(len(data)) * 0.8
	barSpacing := options.Width / float64(len(data)) * 0.2

	// 绘制柱子 / Draw bars
	for i, value := range data {
		barHeight := (value / maxValue) * options.Height
		x := float64(i)*(barWidth+barSpacing) + barSpacing/2
		y := options.Height - barHeight

		g.builder.AddRect(x, y, barWidth, barHeight).
			Fill(options.FillColor).
			Stroke(options.StrokeColor).
			StrokeWidth(1).
			End()
	}
}

// createLineChart 创建折线图 / Create line chart
func (g *SVGGenerator) createLineChart(data []float64, options ChartOptions) {
	if len(data) < 2 {
		return
	}

	// 计算最大值和最小值 / Calculate max and min values
	maxValue, minValue := data[0], data[0]
	for _, v := range data {
		if v > maxValue {
			maxValue = v
		}
		if v < minValue {
			minValue = v
		}
	}

	valueRange := maxValue - minValue
	if valueRange == 0 {
		valueRange = 1
	}

	// 构建路径数据 / Build path data
	pathData := ""
	for i, value := range data {
		x := float64(i) * options.Width / float64(len(data)-1)
		y := options.Height - ((value-minValue)/valueRange)*options.Height

		if i == 0 {
			pathData += fmt.Sprintf("M %.2f %.2f", x, y)
		} else {
			pathData += fmt.Sprintf(" L %.2f %.2f", x, y)
		}
	}

	g.builder.AddPath(pathData).
		Fill(color.RGBA{0, 0, 0, 0}). // 透明填充 / Transparent fill
		Stroke(options.StrokeColor).
		StrokeWidth(2).
		End()
}

// createPieChart 创建饼图 / Create pie chart
func (g *SVGGenerator) createPieChart(data []float64, options ChartOptions) {
	if len(data) == 0 {
		return
	}

	// 计算总和 / Calculate total
	total := 0.0
	for _, v := range data {
		total += v
	}

	if total == 0 {
		return
	}

	// 计算中心和半径 / Calculate center and radius
	cx := options.Width / 2
	cy := options.Height / 2
	radius := math.Min(options.Width, options.Height) / 2 * 0.8

	// 绘制扇形 / Draw sectors
	startAngle := 0.0
	colors := generateColors(len(data))

	for i, value := range data {
		angle := (value / total) * 2 * math.Pi
		endAngle := startAngle + angle

		// 计算路径 / Calculate path
		pathData := g.createArcPath(cx, cy, radius, startAngle, endAngle)

		g.builder.AddPath(pathData).
			Fill(colors[i]).
			Stroke(options.StrokeColor).
			StrokeWidth(1).
			End()

		startAngle = endAngle
	}
}

// createDotPattern 创建点图案 / Create dot pattern
func (g *SVGGenerator) createDotPattern(options PatternOptions) {
	for y := options.Spacing; y < options.Height; y += options.Spacing {
		for x := options.Spacing; x < options.Width; x += options.Spacing {
			g.builder.AddCircle(x, y, options.Size/2).
				Fill(options.Color).
				End()
		}
	}
}

// createStripePattern 创建条纹图案 / Create stripe pattern
func (g *SVGGenerator) createStripePattern(options PatternOptions) {
	for x := 0.0; x < options.Width; x += options.Spacing {
		g.builder.AddRect(x, 0, options.Size, options.Height).
			Fill(options.Color).
			End()
	}
}

// createCheckerboardPattern 创建棋盘图案 / Create checkerboard pattern
func (g *SVGGenerator) createCheckerboardPattern(options PatternOptions) {
	rows := int(options.Height / options.Spacing)
	cols := int(options.Width / options.Spacing)

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			if (row+col)%2 == 0 {
				x := float64(col) * options.Spacing
				y := float64(row) * options.Spacing
				g.builder.AddRect(x, y, options.Spacing, options.Spacing).
					Fill(options.Color).
					End()
			}
		}
	}
}

// createStar 创建星形 / Create star
func (g *SVGGenerator) createStar(options ShapeOptions) {
	points := 5
	if options.Points > 0 {
		points = options.Points
	}

	cx := options.CenterX
	cy := options.CenterY
	outerRadius := options.Size
	innerRadius := outerRadius * 0.4

	pathData := ""
	for i := 0; i < points*2; i++ {
		angle := float64(i) * math.Pi / float64(points)
		var radius float64
		if i%2 == 0 {
			radius = outerRadius
		} else {
			radius = innerRadius
		}

		x := cx + radius*math.Cos(angle-math.Pi/2)
		y := cy + radius*math.Sin(angle-math.Pi/2)

		if i == 0 {
			pathData += fmt.Sprintf("M %.2f %.2f", x, y)
		} else {
			pathData += fmt.Sprintf(" L %.2f %.2f", x, y)
		}
	}
	pathData += " Z"

	g.builder.AddPath(pathData).
		Fill(options.FillColor).
		Stroke(options.StrokeColor).
		StrokeWidth(options.StrokeWidth).
		End()
}

// createPolygon 创建多边形 / Create polygon
func (g *SVGGenerator) createPolygon(options ShapeOptions) {
	points := 6
	if options.Points > 0 {
		points = options.Points
	}

	cx := options.CenterX
	cy := options.CenterY
	radius := options.Size

	pathData := ""
	for i := 0; i < points; i++ {
		angle := float64(i) * 2 * math.Pi / float64(points)
		x := cx + radius*math.Cos(angle-math.Pi/2)
		y := cy + radius*math.Sin(angle-math.Pi/2)

		if i == 0 {
			pathData += fmt.Sprintf("M %.2f %.2f", x, y)
		} else {
			pathData += fmt.Sprintf(" L %.2f %.2f", x, y)
		}
	}
	pathData += " Z"

	g.builder.AddPath(pathData).
		Fill(options.FillColor).
		Stroke(options.StrokeColor).
		StrokeWidth(options.StrokeWidth).
		End()
}

// createArrow 创建箭头 / Create arrow
func (g *SVGGenerator) createArrow(options ShapeOptions) {
	length := options.Size
	width := length * 0.3
	headLength := length * 0.3
	headWidth := width * 2

	cx := options.CenterX
	cy := options.CenterY

	// 箭头主体 / Arrow body
	body := fmt.Sprintf("M %.2f %.2f L %.2f %.2f L %.2f %.2f L %.2f %.2f Z",
		cx-length/2, cy-width/2,
		cx+length/2-headLength, cy-width/2,
		cx+length/2-headLength, cy+width/2,
		cx-length/2, cy+width/2)

	// 箭头头部 / Arrow head
	head := fmt.Sprintf("M %.2f %.2f L %.2f %.2f L %.2f %.2f Z",
		cx+length/2-headLength, cy-headWidth/2,
		cx+length/2, cy,
		cx+length/2-headLength, cy+headWidth/2)

	pathData := body + " " + head

	g.builder.AddPath(pathData).
		Fill(options.FillColor).
		Stroke(options.StrokeColor).
		StrokeWidth(options.StrokeWidth).
		End()
}

// createHeart 创建心形 / Create heart
func (g *SVGGenerator) createHeart(options ShapeOptions) {
	size := options.Size
	cx := options.CenterX
	cy := options.CenterY

	// 心形路径 / Heart path
	pathData := fmt.Sprintf("M %.2f %.2f C %.2f %.2f %.2f %.2f %.2f %.2f C %.2f %.2f %.2f %.2f %.2f %.2f Z",
		cx, cy+size*0.3,
		cx-size*0.7, cy-size*0.3, cx-size*0.3, cy-size*0.3, cx, cy,
		cx+size*0.3, cy-size*0.3, cx+size*0.7, cy-size*0.3, cx, cy+size*0.3)

	g.builder.AddPath(pathData).
		Fill(options.FillColor).
		Stroke(options.StrokeColor).
		StrokeWidth(options.StrokeWidth).
		End()
}

// createArcPath 创建弧形路径 / Create arc path
func (g *SVGGenerator) createArcPath(cx, cy, radius, startAngle, endAngle float64) string {
	x1 := cx + radius*math.Cos(startAngle)
	y1 := cy + radius*math.Sin(startAngle)
	x2 := cx + radius*math.Cos(endAngle)
	y2 := cy + radius*math.Sin(endAngle)

	largeArc := 0
	if endAngle-startAngle > math.Pi {
		largeArc = 1
	}

	return fmt.Sprintf("M %.2f %.2f L %.2f %.2f A %.2f %.2f 0 %d 1 %.2f %.2f Z",
		cx, cy, x1, y1, radius, radius, largeArc, x2, y2)
}

// generateColors 生成颜色数组 / Generate color array
func generateColors(count int) []color.Color {
	colors := make([]color.Color, count)
	for i := 0; i < count; i++ {
		hue := float64(i) * 360.0 / float64(count)
		colors[i] = hslToRGB(hue, 70, 50)
	}
	return colors
}

// hslToRGB 将HSL转换为RGB / Convert HSL to RGB
func hslToRGB(h, s, l float64) color.Color {
	h = h / 360.0
	s = s / 100.0
	l = l / 100.0

	var r, g, b float64

	if s == 0 {
		r = l
		g = l
		b = l
	} else {
		hue2rgb := func(p, q, t float64) float64 {
			if t < 0 {
				t += 1
			}
			if t > 1 {
				t -= 1
			}
			if t < 1.0/6.0 {
				return p + (q-p)*6*t
			}
			if t < 1.0/2.0 {
				return q
			}
			if t < 2.0/3.0 {
				return p + (q-p)*(2.0/3.0-t)*6
			}
			return p
		}

		var q float64
		if l < 0.5 {
			q = l * (1 + s)
		} else {
			q = l + s - l*s
		}
		p := 2*l - q

		r = hue2rgb(p, q, h+1.0/3.0)
		g = hue2rgb(p, q, h)
		b = hue2rgb(p, q, h-1.0/3.0)
	}

	return color.RGBA{
		R: uint8(r * 255),
		G: uint8(g * 255),
		B: uint8(b * 255),
		A: 255,
	}
}

// Options structures / 选项结构体

// ChartOptions 图表选项 / Chart options
type ChartOptions struct {
	Width       float64
	Height      float64
	FillColor   color.Color
	StrokeColor color.Color
}

// GridOptions 网格选项 / Grid options
type GridOptions struct {
	LineColor color.Color
	LineWidth float64
}

// PatternOptions 图案选项 / Pattern options
type PatternOptions struct {
	Width   float64
	Height  float64
	Spacing float64
	Size    float64
	Color   color.Color
}

// ShapeOptions 形状选项 / Shape options
type ShapeOptions struct {
	CenterX     float64
	CenterY     float64
	Size        float64
	Points      int
	FillColor   color.Color
	StrokeColor color.Color
	StrokeWidth float64
}
