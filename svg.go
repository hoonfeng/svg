// Package svg provides a comprehensive SVG (Scalable Vector Graphics) library for Go
// svg包为Go提供了一个全面的SVG（可缩放矢量图形）库
package svg

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hoonfeng/svg/api"
	"github.com/hoonfeng/svg/io"
	"github.com/hoonfeng/svg/renderer"
	. "github.com/hoonfeng/svg/types"
)

// ============================================================================
// 核心SVG结构体 / Core SVG Structure
// ============================================================================

// SVG 统一的SVG操作接口 / Unified SVG operation interface
type SVG struct {
	doc     *Document         // SVG文档 / SVG document
	builder *api.SVGBuilder   // 构建器 / Builder
	gen     *api.SVGGenerator // 生成器 / Generator
	width   int               // 画布宽度 / Canvas width
	height  int               // 画布高度 / Canvas height
}

// ============================================================================
// 创建和初始化方法 / Creation and Initialization Methods
// ============================================================================

// New 创建新的SVG实例 / Create new SVG instance
func New(width, height int) *SVG {
	doc := NewDocument(width, height)
	builder := api.NewSVGBuilder(float64(width), float64(height))
	gen := api.NewSVGGenerator(float64(width), float64(height))
	return &SVG{
		doc:     doc,
		builder: builder,
		gen:     gen,
		width:   width,
		height:  height,
	}
}

// NewWithViewBox 创建带视图框的SVG实例 / Create SVG instance with viewBox
func NewWithViewBox(width, height int, viewX, viewY, viewWidth, viewHeight float64) *SVG {
	doc := NewDocument(width, height)
	doc.SetViewBox(viewX, viewY, viewWidth, viewHeight)
	builder := api.NewSVGBuilderWithViewBox(float64(width), float64(height), viewX, viewY, viewWidth, viewHeight)
	gen := api.NewSVGGenerator(float64(width), float64(height))
	return &SVG{
		doc:     doc,
		builder: builder,
		gen:     gen,
		width:   width,
		height:  height,
	}
}

// Load 从文件加载SVG / Load SVG from file
func Load(filename string) (*SVG, error) {
	doc, err := io.LoadSVG(filename)
	if err != nil {
		return nil, err
	}

	// 解析文档尺寸 / Parse document dimensions
	width, height := extractDimensions(doc)

	builder := api.NewSVGBuilder(width, height)
	gen := api.NewSVGGenerator(width, height)

	return &SVG{
		doc:     doc,
		builder: builder,
		gen:     gen,
		width:   int(width),
		height:  int(height),
	}, nil
}

// Parse 从字符串解析SVG / Parse SVG from string
func Parse(svgContent string) (*SVG, error) {
	doc, err := io.ParseSVG([]byte(svgContent))
	if err != nil {
		return nil, err
	}

	// 解析文档尺寸 / Parse document dimensions
	width, height := extractDimensions(doc)

	builder := api.NewSVGBuilder(width, height)
	gen := api.NewSVGGenerator(width, height)

	return &SVG{
		doc:     doc,
		builder: builder,
		gen:     gen,
		width:   int(width),
		height:  int(height),
	}, nil
}

// ============================================================================
// 基础绘图方法 / Basic Drawing Methods
// ============================================================================

// Rect 创建矩形 / Create rectangle
func (s *SVG) Rect(x, y, width, height float64) *RectElement {
	rectBuilder := s.builder.AddRect(x, y, width, height)
	return &RectElement{builder: rectBuilder, svg: s}
}

// Circle 创建圆形 / Create circle
func (s *SVG) Circle(cx, cy, r float64) *CircleElement {
	circleBuilder := s.builder.AddCircle(cx, cy, r)
	return &CircleElement{builder: circleBuilder, svg: s}
}

// Ellipse 创建椭圆 / Create ellipse
func (s *SVG) Ellipse(cx, cy, rx, ry float64) *EllipseElement {
	ellipseBuilder := s.builder.AddEllipse(cx, cy, rx, ry)
	return &EllipseElement{builder: ellipseBuilder, svg: s}
}

// Line 创建直线 / Create line
func (s *SVG) Line(x1, y1, x2, y2 float64) *LineElement {
	lineBuilder := s.builder.AddLine(x1, y1, x2, y2)
	return &LineElement{builder: lineBuilder, svg: s}
}

// Text 创建文本 / Create text
func (s *SVG) Text(x, y float64, text string) *TextElement {
	textBuilder := s.builder.AddText(x, y, text)
	return &TextElement{builder: textBuilder, svg: s}
}

// Path 创建路径 / Create path
func (s *SVG) Path(d string) *PathElement {
	pathBuilder := s.builder.AddPath(d)
	return &PathElement{builder: pathBuilder, svg: s}
}

// ============================================================================
// 高级绘图方法 / Advanced Drawing Methods
// ============================================================================

// Star 创建星形 / Create star
func (s *SVG) Star(cx, cy, outerRadius float64, points int) *ShapeElement {
	options := api.ShapeOptions{
		CenterX:     cx,
		CenterY:     cy,
		Size:        outerRadius,
		Points:      points,
		FillColor:   color.RGBA{255, 255, 0, 255}, // 默认黄色 / Default yellow
		StrokeColor: color.RGBA{0, 0, 0, 255},     // 默认黑色 / Default black
		StrokeWidth: 1,
	}
	s.gen.CreateShape("star", options)
	return &ShapeElement{svg: s, shapeType: "star", options: options}
}

// Heart 创建心形 / Create heart
func (s *SVG) Heart(cx, cy, size float64) *ShapeElement {
	options := api.ShapeOptions{
		CenterX:     cx,
		CenterY:     cy,
		Size:        size,
		FillColor:   color.RGBA{255, 0, 0, 255}, // 默认红色 / Default red
		StrokeColor: color.RGBA{0, 0, 0, 255},   // 默认黑色 / Default black
		StrokeWidth: 1,
	}
	s.gen.CreateShape("heart", options)
	return &ShapeElement{svg: s, shapeType: "heart", options: options}
}

// Polygon 创建多边形 / Create polygon
func (s *SVG) Polygon(cx, cy, radius float64, sides int) *ShapeElement {
	options := api.ShapeOptions{
		CenterX:     cx,
		CenterY:     cy,
		Size:        radius,
		Points:      sides,
		FillColor:   color.RGBA{0, 255, 0, 255}, // 默认绿色 / Default green
		StrokeColor: color.RGBA{0, 0, 0, 255},   // 默认黑色 / Default black
		StrokeWidth: 1,
	}
	s.gen.CreateShape("polygon", options)
	return &ShapeElement{svg: s, shapeType: "polygon", options: options}
}

// ============================================================================
// 图表方法 / Chart Methods
// ============================================================================

// BarChart 创建柱状图 / Create bar chart
func (s *SVG) BarChart(data []float64, x, y, width, height float64) *ChartElement {
	options := api.ChartOptions{
		Width:       width,
		Height:      height,
		FillColor:   color.RGBA{100, 150, 255, 255}, // 默认蓝色 / Default blue
		StrokeColor: color.RGBA{0, 0, 0, 255},       // 默认黑色 / Default black
	}
	s.gen.CreateChart("bar", data, options)
	return &ChartElement{svg: s, chartType: "bar", data: data, options: options}
}

// PieChart 创建饼图 / Create pie chart
func (s *SVG) PieChart(data []float64, cx, cy, radius float64) *ChartElement {
	options := api.ChartOptions{
		Width:       radius * 2,
		Height:      radius * 2,
		FillColor:   color.RGBA{255, 100, 100, 255}, // 默认粉色 / Default pink
		StrokeColor: color.RGBA{0, 0, 0, 255},       // 默认黑色 / Default black
	}
	s.gen.CreateChart("pie", data, options)
	return &ChartElement{svg: s, chartType: "pie", data: data, options: options}
}

// LineChart 创建折线图 / Create line chart
func (s *SVG) LineChart(data []float64, x, y, width, height float64) *ChartElement {
	options := api.ChartOptions{
		Width:       width,
		Height:      height,
		FillColor:   color.RGBA{0, 0, 0, 0},     // 透明填充 / Transparent fill
		StrokeColor: color.RGBA{255, 0, 0, 255}, // 默认红色 / Default red
	}
	s.gen.CreateChart("line", data, options)
	return &ChartElement{svg: s, chartType: "line", data: data, options: options}
}

// ============================================================================
// 图案和网格方法 / Pattern and Grid Methods
// ============================================================================

// Grid 创建网格 / Create grid
func (s *SVG) Grid(rows, cols int, cellWidth, cellHeight float64) *GridElement {
	options := api.GridOptions{
		LineColor: color.RGBA{128, 128, 128, 255}, // 默认灰色 / Default gray
		LineWidth: 1,
	}
	s.gen.CreateGrid(rows, cols, cellWidth, cellHeight, options)
	return &GridElement{svg: s, rows: rows, cols: cols, cellWidth: cellWidth, cellHeight: cellHeight, options: options}
}

// DotPattern 创建点图案 / Create dot pattern
func (s *SVG) DotPattern(spacing, radius float64) *PatternElement {
	options := api.PatternOptions{
		Spacing: spacing,
		Size:    radius,
		Color:   color.RGBA{0, 0, 0, 255}, // 默认黑色 / Default black
	}
	s.gen.CreatePattern("dots", options)
	return &PatternElement{svg: s, patternType: "dots", options: options}
}

// ============================================================================
// 样式和变换方法 / Style and Transform Methods
// ============================================================================

// Background 设置背景颜色 / Set background color
func (s *SVG) Background(bgColor color.Color) *SVG {
	s.builder.SetBackground(bgColor)
	return s
}

// Group 开始组 / Begin group
func (s *SVG) Group() *GroupElement {
	groupBuilder := s.builder.BeginGroup()
	return &GroupElement{builder: groupBuilder, svg: s}
}

// ============================================================================
// 高级输出方法 / Advanced Output Methods
// ============================================================================

// String 转换为SVG字符串 / Convert to SVG string
func (s *SVG) String() string {
	return s.doc.ToXML()
}

// Save 保存为SVG文件 / Save as SVG file
func (s *SVG) Save(filename string) error {
	return io.SaveSVG(s.doc, filename)
}

// Render 渲染为图像 / Render to image
func (s *SVG) Render(width, height int) (*image.RGBA, error) {
	if width <= 0 {
		width = s.width
	}
	if height <= 0 {
		height = s.height
	}
	return renderer.RenderDocument(s.doc, width, height)
}

// RenderToSize 渲染到指定尺寸 / Render to specified size
func (s *SVG) RenderToSize(width, height int) (*image.RGBA, error) {
	return renderer.RenderDocument(s.doc, width, height)
}

// SavePNG 保存为PNG文件 / Save as PNG file
func (s *SVG) SavePNG(filename string, width, height int) error {
	img, err := s.RenderToSize(width, height)
	if err != nil {
		return err
	}
	return SaveImageToPNG(img, filename)
}

// SaveJPEG 保存为JPEG文件 / Save as JPEG file
func (s *SVG) SaveJPEG(filename string, width, height int, quality int) error {
	img, err := s.RenderToSize(width, height)
	if err != nil {
		return err
	}
	return SaveImageToJPEG(img, filename, quality)
}

// SaveImage 保存为指定格式的图片文件 / Save as image file in specified format
func (s *SVG) SaveImage(filename string, width, height int, format string, quality ...int) error {
	img, err := s.RenderToSize(width, height)
	if err != nil {
		return err
	}

	ext := strings.ToLower(filepath.Ext(filename))
	if format != "" {
		ext = "." + strings.ToLower(format)
	}

	switch ext {
	case ".png":
		return SaveImageToPNG(img, filename)
	case ".jpg", ".jpeg":
		q := 90 // 默认质量 / Default quality
		if len(quality) > 0 {
			q = quality[0]
		}
		return SaveImageToJPEG(img, filename, q)
	default:
		return fmt.Errorf("unsupported image format: %s", ext)
	}
}

// GetImageData 获取图像数据 / Get image data
func (s *SVG) GetImageData(width, height int) (*image.RGBA, error) {
	return s.RenderToSize(width, height)
}

// GetPNGData 获取PNG格式的图像数据 / Get PNG format image data
func (s *SVG) GetPNGData(width, height int) ([]byte, error) {
	img, err := s.RenderToSize(width, height)
	if err != nil {
		return nil, err
	}
	return ImageToPNGBytes(img)
}

// GetJPEGData 获取JPEG格式的图像数据 / Get JPEG format image data
func (s *SVG) GetJPEGData(width, height int, quality int) ([]byte, error) {
	img, err := s.RenderToSize(width, height)
	if err != nil {
		return nil, err
	}
	return ImageToJPEGBytes(img, quality)
}

// GetDocument 获取文档 / Get document
func (s *SVG) GetDocument() *Document {
	return s.doc
}

// GetSize 获取画布尺寸 / Get canvas size
func (s *SVG) GetSize() (int, int) {
	return s.width, s.height
}

// SetSize 设置画布尺寸 / Set canvas size
func (s *SVG) SetSize(width, height int) *SVG {
	s.width = width
	s.height = height
	s.doc.Width = fmt.Sprintf("%d", width)
	s.doc.Height = fmt.Sprintf("%d", height)
	return s
}

// ============================================================================
// 元素类型绑定方法 / Element Type Binding Methods
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
