package svg

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"strconv"
	"strings"

	"github.com/hoonfeng/svg/animation"
	"github.com/hoonfeng/svg/api"
	"github.com/hoonfeng/svg/io"
	"github.com/hoonfeng/svg/output"
	"github.com/hoonfeng/svg/parser"
	"github.com/hoonfeng/svg/renderer"
	"github.com/hoonfeng/svg/types"
)

// Document SVG文档结构 / SVG document structure
type Document struct {
	Width   string `xml:"width,attr"`
	Height  string `xml:"height,attr"`
	ViewBox string `xml:"viewBox,attr"`
	Xmlns   string `xml:"xmlns,attr"`
	Content string `xml:",innerxml"`
}

// SVG 主要结构体 / Main SVG structure
type SVG struct {
	gen      *api.Generator
	renderer *renderer.Renderer
	animator *animation.Animator
	io       *io.Manager
	output   *output.Manager
	parser   *parser.Parser
	doc      *Document
}

// New 创建新的SVG实例 / Create new SVG instance
func New(width, height float64) *SVG {
	return &SVG{
		gen:      api.NewGenerator(width, height),
		renderer: renderer.New(),
		animator: animation.New(),
		io:       io.New(),
		output:   output.New(),
		parser:   parser.New(),
		doc: &Document{
			Width:  fmt.Sprintf("%.0f", width),
			Height: fmt.Sprintf("%.0f", height),
			Xmlns:  "http://www.w3.org/2000/svg",
		},
	}
}

// NewWithViewBox 创建带视图框的SVG实例 / Create SVG instance with viewBox
func NewWithViewBox(width, height float64, viewBox string) *SVG {
	svg := New(width, height)
	svg.doc.ViewBox = viewBox
	return svg
}

// Load 从文件加载SVG / Load SVG from file
func Load(filename string) (*SVG, error) {
	svg := &SVG{
		renderer: renderer.New(),
		animator: animation.New(),
		io:       io.New(),
		output:   output.New(),
		parser:   parser.New(),
	}

	doc, err := svg.io.LoadFromFile(filename)
	if err != nil {
		return nil, err
	}

	svg.doc = doc
	svg.gen = api.NewFromDocument(doc)
	return svg, nil
}

// Parse 解析SVG字符串 / Parse SVG string
func Parse(svgContent string) (*SVG, error) {
	svg := &SVG{
		renderer: renderer.New(),
		animator: animation.New(),
		io:       io.New(),
		output:   output.New(),
		parser:   parser.New(),
	}

	doc, err := svg.parser.ParseString(svgContent)
	if err != nil {
		return nil, err
	}

	svg.doc = doc
	svg.gen = api.NewFromDocument(doc)
	return svg, nil
}

// ============================================================================
// 基础绘图方法 / Basic Drawing Methods
// ============================================================================

// Rect 绘制矩形 / Draw rectangle
func (s *SVG) Rect(x, y, width, height float64) *RectElement {
	builder := s.gen.Rect(x, y, width, height)
	return &RectElement{builder: builder, svg: s}
}

// Circle 绘制圆形 / Draw circle
func (s *SVG) Circle(cx, cy, r float64) *CircleElement {
	builder := s.gen.Circle(cx, cy, r)
	return &CircleElement{builder: builder, svg: s}
}

// Ellipse 绘制椭圆 / Draw ellipse
func (s *SVG) Ellipse(cx, cy, rx, ry float64) *EllipseElement {
	builder := s.gen.Ellipse(cx, cy, rx, ry)
	return &EllipseElement{builder: builder, svg: s}
}

// Line 绘制直线 / Draw line
func (s *SVG) Line(x1, y1, x2, y2 float64) *LineElement {
	builder := s.gen.Line(x1, y1, x2, y2)
	return &LineElement{builder: builder, svg: s}
}

// Text 绘制文本 / Draw text
func (s *SVG) Text(x, y float64, text string) *TextElement {
	builder := s.gen.Text(x, y, text)
	return &TextElement{builder: builder, svg: s}
}

// Path 绘制路径 / Draw path
func (s *SVG) Path(d string) *PathElement {
	builder := s.gen.Path(d)
	return &PathElement{builder: builder, svg: s}
}

// ============================================================================
// 高级绘图方法 / Advanced Drawing Methods
// ============================================================================

// Star 绘制星形 / Draw star
func (s *SVG) Star(cx, cy, outerRadius, innerRadius float64, points int) *ShapeElement {
	options := api.ShapeOptions{
		CenterX:     cx,
		CenterY:     cy,
		OuterRadius: outerRadius,
		InnerRadius: innerRadius,
		Points:      points,
	}
	return &ShapeElement{svg: s, shapeType: "star", options: options}
}

// Heart 绘制心形 / Draw heart
func (s *SVG) Heart(cx, cy, size float64) *ShapeElement {
	options := api.ShapeOptions{
		CenterX: cx,
		CenterY: cy,
		Size:    size,
	}
	return &ShapeElement{svg: s, shapeType: "heart", options: options}
}

// Polygon 绘制多边形 / Draw polygon
func (s *SVG) Polygon(points []types.Point) *ShapeElement {
	options := api.ShapeOptions{
		Points: len(points),
	}
	return &ShapeElement{svg: s, shapeType: "polygon", options: options}
}

// ============================================================================
// 图表方法 / Chart Methods
// ============================================================================

// BarChart 绘制柱状图 / Draw bar chart
func (s *SVG) BarChart(x, y, width, height float64, data []float64) *ChartElement {
	options := api.ChartOptions{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	}
	return &ChartElement{svg: s, chartType: "bar", data: data, options: options}
}

// PieChart 绘制饼图 / Draw pie chart
func (s *SVG) PieChart(cx, cy, radius float64, data []float64) *ChartElement {
	options := api.ChartOptions{
		CenterX: cx,
		CenterY: cy,
		Radius:  radius,
	}
	return &ChartElement{svg: s, chartType: "pie", data: data, options: options}
}

// LineChart 绘制折线图 / Draw line chart
func (s *SVG) LineChart(x, y, width, height float64, data []float64) *ChartElement {
	options := api.ChartOptions{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	}
	return &ChartElement{svg: s, chartType: "line", data: data, options: options}
}

// ============================================================================
// 图案和网格方法 / Pattern and Grid Methods
// ============================================================================

// Grid 绘制网格 / Draw grid
func (s *SVG) Grid(rows, cols int, cellWidth, cellHeight float64) *GridElement {
	options := api.GridOptions{}
	return &GridElement{
		svg:        s,
		rows:       rows,
		cols:       cols,
		cellWidth:  cellWidth,
		cellHeight: cellHeight,
		options:    options,
	}
}

// DotPattern 绘制点图案 / Draw dot pattern
func (s *SVG) DotPattern(spacing float64) *PatternElement {
	options := api.PatternOptions{
		Spacing: spacing,
	}
	return &PatternElement{svg: s, patternType: "dot", options: options}
}

// ============================================================================
// 样式和变换方法 / Style and Transform Methods
// ============================================================================

// Background 设置背景 / Set background
func (s *SVG) Background(color color.Color) *SVG {
	s.gen.SetBackground(color)
	return s
}

// Group 创建组 / Create group
func (s *SVG) Group() *GroupElement {
	builder := s.gen.Group()
	return &GroupElement{builder: builder, svg: s}
}

// ============================================================================
// 高级输出方法 / Advanced Output Methods
// ============================================================================

// String 转换为字符串 / Convert to string
func (s *SVG) String() string {
	return s.gen.String()
}

// Save 保存到文件 / Save to file
func (s *SVG) Save(filename string) error {
	return s.io.SaveToFile(filename, s.String())
}

// Render 渲染为图像 / Render as image
func (s *SVG) Render() (image.Image, error) {
	return s.renderer.RenderToImage(s.String())
}

// SavePNG 保存为PNG / Save as PNG
func (s *SVG) SavePNG(filename string) error {
	img, err := s.Render()
	if err != nil {
		return err
	}
	return SaveImageToPNG(img, filename)
}

// SaveJPEG 保存为JPEG / Save as JPEG
func (s *SVG) SaveJPEG(filename string, quality int) error {
	img, err := s.Render()
	if err != nil {
		return err
	}
	return SaveImageToJPEG(img, filename, quality)
}

// SaveImage 保存为指定格式的图像 / Save as image in specified format
func (s *SVG) SaveImage(filename, format string, quality int) error {
	img, err := s.Render()
	if err != nil {
		return err
	}

	switch strings.ToLower(format) {
	case "png":
		return SaveImageToPNG(img, filename)
	case "jpeg", "jpg":
		return SaveImageToJPEG(img, filename, quality)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

// GetImageData 获取图像数据 / Get image data
func (s *SVG) GetImageData() (image.Image, error) {
	return s.renderer.RenderToImage(s.String())
}

// GetPNGData 获取PNG数据 / Get PNG data
func (s *SVG) GetPNGData() ([]byte, error) {
	img, err := s.GetImageData()
	if err != nil {
		return nil, err
	}
	return ImageToPNGBytes(img)
}

// GetJPEGData 获取JPEG数据 / Get JPEG data
func (s *SVG) GetJPEGData(quality int) ([]byte, error) {
	img, err := s.GetImageData()
	if err != nil {
		return nil, err
	}
	return ImageToJPEGBytes(img, quality)
}

// ============================================================================
// 画布尺寸方法 / Canvas Size Methods
// ============================================================================

// GetWidth 获取画布宽度 / Get canvas width
func (s *SVG) GetWidth() float64 {
	width, _ := extractDimensions(s.doc)
	return width
}

// GetHeight 获取画布高度 / Get canvas height
func (s *SVG) GetHeight() float64 {
	_, height := extractDimensions(s.doc)
	return height
}

// SetSize 设置画布尺寸 / Set canvas size
func (s *SVG) SetSize(width, height float64) *SVG {
	s.doc.Width = fmt.Sprintf("%.0f", width)
	s.doc.Height = fmt.Sprintf("%.0f", height)
	s.gen.SetSize(width, height)
	return s
}

// ============================================================================
// 元素类型定义 / Element Type Definitions
// ============================================================================

// RectElement 矩形元素 / Rectangle element
type RectElement struct {
	builder *api.RectBuilder
	svg     *SVG
}

func (r *RectElement) Fill(c color.Color) *RectElement {
	r.builder.Fill(c)
	return r
}

func (r *RectElement) Stroke(c color.Color) *RectElement {
	r.builder.Stroke(c)
	return r
}

func (r *RectElement) StrokeWidth(width float64) *RectElement {
	r.builder.StrokeWidth(width)
	return r
}

func (r *RectElement) Rx(rx float64) *RectElement {
	r.builder.Rx(rx)
	return r
}

func (r *RectElement) Ry(ry float64) *RectElement {
	r.builder.Ry(ry)
	return r
}

func (r *RectElement) End() *SVG {
	r.builder.End()
	return r.svg
}

// CircleElement 圆形元素 / Circle element
type CircleElement struct {
	builder *api.CircleBuilder
	svg     *SVG
}

func (c *CircleElement) Fill(color color.Color) *CircleElement {
	c.builder.Fill(color)
	return c
}

func (c *CircleElement) Stroke(color color.Color) *CircleElement {
	c.builder.Stroke(color)
	return c
}

func (c *CircleElement) StrokeWidth(width float64) *CircleElement {
	c.builder.StrokeWidth(width)
	return c
}

func (c *CircleElement) End() *SVG {
	c.builder.End()
	return c.svg
}

// EllipseElement 椭圆元素 / Ellipse element
type EllipseElement struct {
	builder *api.EllipseBuilder
	svg     *SVG
}

func (e *EllipseElement) Fill(color color.Color) *EllipseElement {
	e.builder.Fill(color)
	return e
}

func (e *EllipseElement) Stroke(color color.Color) *EllipseElement {
	e.builder.Stroke(color)
	return e
}

func (e *EllipseElement) StrokeWidth(width float64) *EllipseElement {
	e.builder.StrokeWidth(width)
	return e
}

func (e *EllipseElement) End() *SVG {
	e.builder.End()
	return e.svg
}

// LineElement 直线元素 / Line element
type LineElement struct {
	builder *api.LineBuilder
	svg     *SVG
}

func (l *LineElement) Stroke(color color.Color) *LineElement {
	l.builder.Stroke(color)
	return l
}

func (l *LineElement) StrokeWidth(width float64) *LineElement {
	l.builder.StrokeWidth(width)
	return l
}

func (l *LineElement) End() *SVG {
	l.builder.End()
	return l.svg
}

// TextElement 文本元素 / Text element
type TextElement struct {
	builder *api.TextBuilder
	svg     *SVG
}

func (t *TextElement) Fill(color color.Color) *TextElement {
	t.builder.Fill(color)
	return t
}

func (t *TextElement) FontFamily(family string) *TextElement {
	t.builder.FontFamily(family)
	return t
}

func (t *TextElement) FontSize(size float64) *TextElement {
	t.builder.FontSize(size)
	return t
}

func (t *TextElement) FontWeight(weight string) *TextElement {
	t.builder.FontWeight(weight)
	return t
}

func (t *TextElement) End() *SVG {
	t.builder.End()
	return t.svg
}

// PathElement 路径元素 / Path element
type PathElement struct {
	builder *api.PathBuilder
	svg     *SVG
}

func (p *PathElement) Fill(color color.Color) *PathElement {
	p.builder.Fill(color)
	return p
}

func (p *PathElement) Stroke(color color.Color) *PathElement {
	p.builder.Stroke(color)
	return p
}

func (p *PathElement) StrokeWidth(width float64) *PathElement {
	p.builder.StrokeWidth(width)
	return p
}

func (p *PathElement) End() *SVG {
	p.builder.End()
	return p.svg
}

// ShapeElement 形状元素 / Shape element
type ShapeElement struct {
	svg       *SVG
	shapeType string
	options   api.ShapeOptions
}

func (s *ShapeElement) Fill(color color.Color) *ShapeElement {
	s.options.FillColor = color
	s.svg.gen.CreateShape(s.shapeType, s.options)
	return s
}

func (s *ShapeElement) Stroke(color color.Color) *ShapeElement {
	s.options.StrokeColor = color
	s.svg.gen.CreateShape(s.shapeType, s.options)
	return s
}

func (s *ShapeElement) StrokeWidth(width float64) *ShapeElement {
	s.options.StrokeWidth = width
	s.svg.gen.CreateShape(s.shapeType, s.options)
	return s
}

func (s *ShapeElement) End() *SVG {
	return s.svg
}

// ChartElement 图表元素 / Chart element
type ChartElement struct {
	svg       *SVG
	chartType string
	data      []float64
	options   api.ChartOptions
}

func (c *ChartElement) Fill(color color.Color) *ChartElement {
	c.options.FillColor = color
	c.svg.gen.CreateChart(c.chartType, c.data, c.options)
	return c
}

func (c *ChartElement) Stroke(color color.Color) *ChartElement {
	c.options.StrokeColor = color
	c.svg.gen.CreateChart(c.chartType, c.data, c.options)
	return c
}

func (c *ChartElement) End() *SVG {
	return c.svg
}

// GroupElement 组元素 / Group element
type GroupElement struct {
	builder *api.GroupBuilder
	svg     *SVG
}

func (g *GroupElement) Translate(x, y float64) *GroupElement {
	g.builder.Transform(fmt.Sprintf("translate(%.2f,%.2f)", x, y))
	return g
}

func (g *GroupElement) Scale(sx, sy float64) *GroupElement {
	g.builder.Transform(fmt.Sprintf("scale(%.2f,%.2f)", sx, sy))
	return g
}

func (g *GroupElement) Rotate(angle float64) *GroupElement {
	g.builder.Transform(fmt.Sprintf("rotate(%.2f)", angle))
	return g
}

func (g *GroupElement) Transform(transform string) *GroupElement {
	g.builder.Transform(transform)
	return g
}

func (g *GroupElement) End() *SVG {
	g.builder.End()
	return g.svg
}

// GridElement 网格元素 / Grid element
type GridElement struct {
	svg        *SVG
	rows       int
	cols       int
	cellWidth  float64
	cellHeight float64
	options    api.GridOptions
}

func (g *GridElement) LineColor(color color.Color) *GridElement {
	g.options.LineColor = color
	g.svg.gen.CreateGrid(g.rows, g.cols, g.cellWidth, g.cellHeight, g.options)
	return g
}

func (g *GridElement) LineWidth(width float64) *GridElement {
	g.options.LineWidth = width
	g.svg.gen.CreateGrid(g.rows, g.cols, g.cellWidth, g.cellHeight, g.options)
	return g
}

func (g *GridElement) End() *SVG {
	return g.svg
}

// PatternElement 图案元素 / Pattern element
type PatternElement struct {
	svg         *SVG
	patternType string
	options     api.PatternOptions
}

func (p *PatternElement) Color(color color.Color) *PatternElement {
	p.options.Color = color
	p.svg.gen.CreatePattern(p.patternType, p.options)
	return p
}

func (p *PatternElement) Spacing(spacing float64) *PatternElement {
	p.options.Spacing = spacing
	p.svg.gen.CreatePattern(p.patternType, p.options)
	return p
}

func (p *PatternElement) End() *SVG {
	return p.svg
}

// ============================================================================
// 辅助函数 / Helper Functions
// ============================================================================

// SaveImageToPNG 保存图像为PNG文件 / Save image as PNG file
func SaveImageToPNG(img image.Image, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return png.Encode(file, img)
}

// SaveImageToJPEG 保存图像为JPEG文件 / Save image as JPEG file
func SaveImageToJPEG(img image.Image, filename string, quality int) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return jpeg.Encode(file, img, &jpeg.Options{Quality: quality})
}

// ImageToPNGBytes 将图像转换为PNG字节数据 / Convert image to PNG bytes
func ImageToPNGBytes(img image.Image) ([]byte, error) {
	var buf strings.Builder
	err := png.Encode(&buf, img)
	if err != nil {
		return nil, err
	}
	return []byte(buf.String()), nil
}

// ImageToJPEGBytes 将图像转换为JPEG字节数据 / Convert image to JPEG bytes
func ImageToJPEGBytes(img image.Image, quality int) ([]byte, error) {
	var buf strings.Builder
	err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality})
	if err != nil {
		return nil, err
	}
	return []byte(buf.String()), nil
}

// extractDimensions 提取文档尺寸 / Extract document dimensions
func extractDimensions(doc *Document) (width, height float64) {
	width, height = 100, 100 // 默认值 / Default values

	if doc.Width != "" {
		if !strings.Contains(doc.Width, "%") {
			if val, err := parseFloat(doc.Width, 0); err == nil {
				width = val
			}
		}
	}

	if doc.Height != "" {
		if !strings.Contains(doc.Height, "%") {
			if val, err := parseFloat(doc.Height, 0); err == nil {
				height = val
			}
		}
	}

	return width, height
}

// parseFloat 解析浮点数 / Parse float
func parseFloat(s string, defaultValue float64) (float64, error) {
	s = strings.TrimSuffix(s, "px")
	s = strings.TrimSuffix(s, "pt")
	s = strings.TrimSuffix(s, "em")
	s = strings.TrimSuffix(s, "rem")

	if s == "" {
		return defaultValue, nil
	}

	return strconv.ParseFloat(s, 64)
}