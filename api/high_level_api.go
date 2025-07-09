// Package api provides high-level APIs for SVG creation and manipulation
// api包为SVG创建和操作提供高级API
package api

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/hoonfeng/svg/elements"
	"github.com/hoonfeng/svg/types"
)

// SVGBuilder 高级SVG构建器 / High-level SVG builder
type SVGBuilder struct {
	doc          *types.Document
	currentGroup *elements.Group
	groupStack   []*elements.Group
}

// NewSVGBuilder 创建新的SVG构建器 / Create new SVG builder
func NewSVGBuilder(width, height float64) *SVGBuilder {
	doc := types.NewDocument(int(width), int(height))
	doc.SetViewBox(0, 0, width, height)

	return &SVGBuilder{
		doc:        doc,
		groupStack: make([]*elements.Group, 0),
	}
}

// NewSVGBuilderWithViewBox 创建带视口的SVG构建器 / Create SVG builder with viewBox
func NewSVGBuilderWithViewBox(width, height, viewX, viewY, viewWidth, viewHeight float64) *SVGBuilder {
	doc := types.NewDocument(int(width), int(height))
	doc.SetViewBox(viewX, viewY, viewWidth, viewHeight)

	return &SVGBuilder{
		doc:        doc,
		groupStack: make([]*elements.Group, 0),
	}
}

// GetDocument 获取构建的文档 / Get built document
func (b *SVGBuilder) GetDocument() *types.Document {
	return b.doc
}

// SetBackground 设置背景颜色 / Set background color
func (b *SVGBuilder) SetBackground(bgColor color.Color) *SVGBuilder {
	// 创建背景矩形 / Create background rectangle
	width, height := b.getDocumentSize()
	bg := elements.NewRect(0, 0, width, height)
	bg.SetAttribute("fill", colorToString(bgColor))
	bg.SetAttribute("stroke", "none")

	// 插入到文档开头 / Insert at document beginning
	elements := b.doc.Elements
	b.doc.Elements = make([]types.Element, 0, len(elements)+1)
	b.doc.Elements = append(b.doc.Elements, bg)
	b.doc.Elements = append(b.doc.Elements, elements...)

	return b
}

// AddRect 添加矩形 / Add rectangle
func (b *SVGBuilder) AddRect(x, y, width, height float64) *RectBuilder {
	rect := elements.NewRect(x, y, width, height)
	b.addElement(rect)
	return &RectBuilder{rect: rect, builder: b}
}

// AddCircle 添加圆形 / Add circle
func (b *SVGBuilder) AddCircle(cx, cy, r float64) *CircleBuilder {
	circle := elements.NewCircle(cx, cy, r)
	b.addElement(circle)
	return &CircleBuilder{circle: circle, builder: b}
}

// AddEllipse 添加椭圆 / Add ellipse
func (b *SVGBuilder) AddEllipse(cx, cy, rx, ry float64) *EllipseBuilder {
	ellipse := elements.NewEllipse(cx, cy, rx, ry)
	b.addElement(ellipse)
	return &EllipseBuilder{ellipse: ellipse, builder: b}
}

// AddLine 添加直线 / Add line
func (b *SVGBuilder) AddLine(x1, y1, x2, y2 float64) *LineBuilder {
	line := elements.NewLine(x1, y1, x2, y2)
	b.addElement(line)
	return &LineBuilder{line: line, builder: b}
}

// AddText 添加文本 / Add text
func (b *SVGBuilder) AddText(x, y float64, text string) *TextBuilder {
	textElement := elements.NewText(x, y, text)
	b.addElement(textElement)
	return &TextBuilder{text: textElement, builder: b}
}

// AddPath 添加路径 / Add path
func (b *SVGBuilder) AddPath(pathData string) *PathBuilder {
	path := elements.NewPath(pathData)
	b.addElement(path)
	return &PathBuilder{path: path, builder: b}
}

// BeginGroup 开始组 / Begin group
func (b *SVGBuilder) BeginGroup() *GroupBuilder {
	group := elements.NewGroup()
	b.addElement(group)

	// 推入组栈 / Push to group stack
	b.groupStack = append(b.groupStack, b.currentGroup)
	b.currentGroup = group

	return &GroupBuilder{group: group, builder: b}
}

// EndGroup 结束组 / End group
func (b *SVGBuilder) EndGroup() *SVGBuilder {
	if len(b.groupStack) > 0 {
		// 弹出组栈 / Pop from group stack
		b.currentGroup = b.groupStack[len(b.groupStack)-1]
		b.groupStack = b.groupStack[:len(b.groupStack)-1]
	}
	return b
}

// addElement 添加元素到当前容器 / Add element to current container
func (b *SVGBuilder) addElement(element types.Element) {
	if b.currentGroup != nil {
		b.currentGroup.AppendChild(element)
	} else {
		b.doc.AppendElement(element)
	}
}

// getDocumentSize 获取文档尺寸 / Get document size
func (b *SVGBuilder) getDocumentSize() (float64, float64) {
	width := 100.0
	height := 100.0

	// 尝试解析Width字符串
	if b.doc.Width != "" {
		if w, err := strconv.ParseFloat(b.doc.Width, 64); err == nil {
			width = w
		}
	}

	// 尝试解析Height字符串
	if b.doc.Height != "" {
		if h, err := strconv.ParseFloat(b.doc.Height, 64); err == nil {
			height = h
		}
	}

	return width, height
}

// colorToString 将颜色转换为字符串 / Convert color to string
func colorToString(c color.Color) string {
	r, g, b, a := c.RGBA()
	if a == 0xFFFF {
		return fmt.Sprintf("rgb(%d,%d,%d)", r>>8, g>>8, b>>8)
	} else {
		return fmt.Sprintf("rgba(%d,%d,%d,%.2f)", r>>8, g>>8, b>>8, float64(a)/0xFFFF)
	}
}

// Element builders for fluent API
// 元素构建器，用于流式API

// RectBuilder 矩形构建器 / Rectangle builder
type RectBuilder struct {
	rect    *elements.Rect
	builder *SVGBuilder
}

// Fill 设置填充颜色 / Set fill color
func (rb *RectBuilder) Fill(color color.Color) *RectBuilder {
	rb.rect.SetAttribute("fill", colorToString(color))
	return rb
}

// Stroke 设置描边颜色 / Set stroke color
func (rb *RectBuilder) Stroke(color color.Color) *RectBuilder {
	rb.rect.SetAttribute("stroke", colorToString(color))
	return rb
}

// StrokeWidth 设置描边宽度 / Set stroke width
func (rb *RectBuilder) StrokeWidth(width float64) *RectBuilder {
	rb.rect.SetAttribute("stroke-width", fmt.Sprintf("%.2f", width))
	return rb
}

// Rx 设置圆角半径X / Set border radius X
func (rb *RectBuilder) Rx(rx float64) *RectBuilder {
	rb.rect.SetAttribute("rx", fmt.Sprintf("%.2f", rx))
	return rb
}

// Ry 设置圆角半径Y / Set border radius Y
func (rb *RectBuilder) Ry(ry float64) *RectBuilder {
	rb.rect.SetAttribute("ry", fmt.Sprintf("%.2f", ry))
	return rb
}

// End 结束矩形构建 / End rectangle building
func (rb *RectBuilder) End() *SVGBuilder {
	return rb.builder
}

// CircleBuilder 圆形构建器 / Circle builder
type CircleBuilder struct {
	circle  *elements.Circle
	builder *SVGBuilder
}

// Fill 设置填充颜色 / Set fill color
func (cb *CircleBuilder) Fill(color color.Color) *CircleBuilder {
	cb.circle.SetAttribute("fill", colorToString(color))
	return cb
}

// Stroke 设置描边颜色 / Set stroke color
func (cb *CircleBuilder) Stroke(color color.Color) *CircleBuilder {
	cb.circle.SetAttribute("stroke", colorToString(color))
	return cb
}

// StrokeWidth 设置描边宽度 / Set stroke width
func (cb *CircleBuilder) StrokeWidth(width float64) *CircleBuilder {
	cb.circle.SetAttribute("stroke-width", fmt.Sprintf("%.2f", width))
	return cb
}

// End 结束圆形构建 / End circle building
func (cb *CircleBuilder) End() *SVGBuilder {
	return cb.builder
}

// EllipseBuilder 椭圆构建器 / Ellipse builder
type EllipseBuilder struct {
	ellipse *elements.Ellipse
	builder *SVGBuilder
}

// Fill 设置填充颜色 / Set fill color
func (eb *EllipseBuilder) Fill(color color.Color) *EllipseBuilder {
	eb.ellipse.SetAttribute("fill", colorToString(color))
	return eb
}

// Stroke 设置描边颜色 / Set stroke color
func (eb *EllipseBuilder) Stroke(color color.Color) *EllipseBuilder {
	eb.ellipse.SetAttribute("stroke", colorToString(color))
	return eb
}

// StrokeWidth 设置描边宽度 / Set stroke width
func (eb *EllipseBuilder) StrokeWidth(width float64) *EllipseBuilder {
	eb.ellipse.SetAttribute("stroke-width", fmt.Sprintf("%.2f", width))
	return eb
}

// End 结束椭圆构建 / End ellipse building
func (eb *EllipseBuilder) End() *SVGBuilder {
	return eb.builder
}

// LineBuilder 直线构建器 / Line builder
type LineBuilder struct {
	line    *elements.Line
	builder *SVGBuilder
}

// Stroke 设置描边颜色 / Set stroke color
func (lb *LineBuilder) Stroke(color color.Color) *LineBuilder {
	lb.line.SetAttribute("stroke", colorToString(color))
	return lb
}

// StrokeWidth 设置描边宽度 / Set stroke width
func (lb *LineBuilder) StrokeWidth(width float64) *LineBuilder {
	lb.line.SetAttribute("stroke-width", fmt.Sprintf("%.2f", width))
	return lb
}

// End 结束直线构建 / End line building
func (lb *LineBuilder) End() *SVGBuilder {
	return lb.builder
}

// TextBuilder 文本构建器 / Text builder
type TextBuilder struct {
	text    *elements.Text
	builder *SVGBuilder
}

// Fill 设置填充颜色 / Set fill color
func (tb *TextBuilder) Fill(color color.Color) *TextBuilder {
	tb.text.SetAttribute("fill", colorToString(color))
	return tb
}

// FontFamily 设置字体族 / Set font family
func (tb *TextBuilder) FontFamily(family string) *TextBuilder {
	tb.text.SetAttribute("font-family", family)
	return tb
}

// FontSize 设置字体大小 / Set font size
func (tb *TextBuilder) FontSize(size float64) *TextBuilder {
	tb.text.SetAttribute("font-size", fmt.Sprintf("%.2f", size))
	return tb
}

// FontWeight 设置字体粗细 / Set font weight
func (tb *TextBuilder) FontWeight(weight string) *TextBuilder {
	tb.text.SetAttribute("font-weight", weight)
	return tb
}

// TextAnchor 设置文本锚点 / Set text anchor
func (tb *TextBuilder) TextAnchor(anchor string) *TextBuilder {
	tb.text.SetAttribute("text-anchor", anchor)
	return tb
}

// End 结束文本构建 / End text building
func (tb *TextBuilder) End() *SVGBuilder {
	return tb.builder
}

// PathBuilder 路径构建器 / Path builder
type PathBuilder struct {
	path    *elements.Path
	builder *SVGBuilder
}

// Fill 设置填充颜色 / Set fill color
func (pb *PathBuilder) Fill(color color.Color) *PathBuilder {
	pb.path.SetAttribute("fill", colorToString(color))
	return pb
}

// Stroke 设置描边颜色 / Set stroke color
func (pb *PathBuilder) Stroke(color color.Color) *PathBuilder {
	pb.path.SetAttribute("stroke", colorToString(color))
	return pb
}

// StrokeWidth 设置描边宽度 / Set stroke width
func (pb *PathBuilder) StrokeWidth(width float64) *PathBuilder {
	pb.path.SetAttribute("stroke-width", fmt.Sprintf("%.2f", width))
	return pb
}

// End 结束路径构建 / End path building
func (pb *PathBuilder) End() *SVGBuilder {
	return pb.builder
}

// GroupBuilder 组构建器 / Group builder
type GroupBuilder struct {
	group   *elements.Group
	builder *SVGBuilder
}

// Transform 设置变换 / Set transform
func (gb *GroupBuilder) Transform(transform string) *GroupBuilder {
	gb.group.SetAttribute("transform", transform)
	return gb
}

// Translate 设置平移变换 / Set translate transform
func (gb *GroupBuilder) Translate(x, y float64) *GroupBuilder {
	transform := fmt.Sprintf("translate(%.2f,%.2f)", x, y)
	gb.group.SetAttribute("transform", transform)
	return gb
}

// Scale 设置缩放变换 / Set scale transform
func (gb *GroupBuilder) Scale(sx, sy float64) *GroupBuilder {
	transform := fmt.Sprintf("scale(%.2f,%.2f)", sx, sy)
	gb.group.SetAttribute("transform", transform)
	return gb
}

// Rotate 设置旋转变换 / Set rotate transform
func (gb *GroupBuilder) Rotate(angle, cx, cy float64) *GroupBuilder {
	transform := fmt.Sprintf("rotate(%.2f,%.2f,%.2f)", angle, cx, cy)
	gb.group.SetAttribute("transform", transform)
	return gb
}

// End 结束组构建 / End group building
func (gb *GroupBuilder) End() *SVGBuilder {
	return gb.builder.EndGroup()
}
