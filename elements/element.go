package elements

import (
	"fmt"
	"strings"

	"github.com/hoonfeng/svg/types"
)

// BaseElement 是所有SVG元素的基础实现
type BaseElement struct {
	tag        string
	id         string
	attributes map[string]string
	children   []types.Element
	parent     types.Element
}

// NewBaseElement 创建一个新的基础元素
func NewBaseElement(tag string) *BaseElement {
	return &BaseElement{
		tag:        tag,
		attributes: make(map[string]string),
		children:   make([]types.Element, 0),
	}
}

// Tag 返回元素的标签名
func (e *BaseElement) Tag() string {
	return e.tag
}

// ID 返回元素的ID
func (e *BaseElement) ID() string {
	return e.id
}

// SetID 设置元素的ID
func (e *BaseElement) SetID(id string) {
	e.id = id
	e.attributes["id"] = id
}

// Attributes 返回元素的所有属性
func (e *BaseElement) Attributes() map[string]string {
	return e.attributes
}

// GetAttributes 是Attributes方法的别名，用于实现Element接口
func (e *BaseElement) GetAttributes() map[string]string {
	return e.Attributes()
}

// SetAttribute 设置元素的属性
func (e *BaseElement) SetAttribute(name, value string) {
	e.attributes[name] = value
}

// GetAttribute 获取元素的属性
func (e *BaseElement) GetAttribute(name string, defaultValue ...string) (string, bool) {
	value, ok := e.attributes[name]
	if !ok && len(defaultValue) > 0 {
		return defaultValue[0], true
	}
	return value, ok
}

// RemoveAttribute 移除元素的属性
func (e *BaseElement) RemoveAttribute(name string) {
	delete(e.attributes, name)
}

// Children 返回元素的子元素
func (e *BaseElement) Children() []types.Element {
	return e.children
}

// AppendChild 添加子元素
func (e *BaseElement) AppendChild(child types.Element) {
	e.children = append(e.children, child)
	child.SetParent(e)
}

// RemoveChild 移除子元素
func (e *BaseElement) RemoveChild(child types.Element) {
	for i, c := range e.children {
		if c == child {
			e.children = append(e.children[:i], e.children[i+1:]...)
			child.SetParent(nil)
			return
		}
	}
}

// Parent 返回元素的父元素
func (e *BaseElement) Parent() types.Element {
	return e.parent
}

// SetParent 设置元素的父元素
func (e *BaseElement) SetParent(parent types.Element) {
	e.parent = parent
}

// Clone 克隆元素
func (e *BaseElement) Clone() types.Element {
	clone := NewBaseElement(e.tag)
	clone.id = e.id

	// 复制属性
	for k, v := range e.attributes {
		clone.attributes[k] = v
	}

	// 复制子元素
	for _, child := range e.children {
		childClone := child.Clone()
		clone.AppendChild(childClone)
	}

	return clone
}

// ToXML 将元素转换为XML字符串
func (e *BaseElement) ToXML() string {
	var sb strings.Builder

	sb.WriteString("<")
	sb.WriteString(e.tag)

	// 添加属性
	for name, value := range e.attributes {
		sb.WriteString(fmt.Sprintf(` %s="%s"`, name, value))
	}

	if len(e.children) == 0 {
		sb.WriteString("/>")
	} else {
		sb.WriteString(">")

		// 添加子元素
		for _, child := range e.children {
			sb.WriteString(child.ToXML())
		}

		sb.WriteString("</")
		sb.WriteString(e.tag)
		sb.WriteString(">")
	}

	return sb.String()
}

// Circle 表示SVG圆形元素
type Circle struct {
	*BaseElement
}

// NewCircle 创建一个新的圆形元素
func NewCircle(cx, cy, r float64) *Circle {
	circle := &Circle{
		BaseElement: NewBaseElement("circle"),
	}
	circle.SetAttribute("cx", fmt.Sprintf("%f", cx))
	circle.SetAttribute("cy", fmt.Sprintf("%f", cy))
	circle.SetAttribute("r", fmt.Sprintf("%f", r))
	return circle
}

// Rect 表示SVG矩形元素
type Rect struct {
	*BaseElement
}

// NewRect 创建一个新的矩形元素
func NewRect(x, y, width, height float64) *Rect {
	rect := &Rect{
		BaseElement: NewBaseElement("rect"),
	}
	rect.SetAttribute("x", fmt.Sprintf("%f", x))
	rect.SetAttribute("y", fmt.Sprintf("%f", y))
	rect.SetAttribute("width", fmt.Sprintf("%f", width))
	rect.SetAttribute("height", fmt.Sprintf("%f", height))
	return rect
}

// Ellipse 表示SVG椭圆元素
type Ellipse struct {
	*BaseElement
}

// NewEllipse 创建一个新的椭圆元素
func NewEllipse(cx, cy, rx, ry float64) *Ellipse {
	ellipse := &Ellipse{
		BaseElement: NewBaseElement("ellipse"),
	}
	ellipse.SetAttribute("cx", fmt.Sprintf("%f", cx))
	ellipse.SetAttribute("cy", fmt.Sprintf("%f", cy))
	ellipse.SetAttribute("rx", fmt.Sprintf("%f", rx))
	ellipse.SetAttribute("ry", fmt.Sprintf("%f", ry))
	return ellipse
}

// Line 表示SVG线段元素
type Line struct {
	*BaseElement
}

// NewLine 创建一个新的线段元素
func NewLine(x1, y1, x2, y2 float64) *Line {
	line := &Line{
		BaseElement: NewBaseElement("line"),
	}
	line.SetAttribute("x1", fmt.Sprintf("%f", x1))
	line.SetAttribute("y1", fmt.Sprintf("%f", y1))
	line.SetAttribute("x2", fmt.Sprintf("%f", x2))
	line.SetAttribute("y2", fmt.Sprintf("%f", y2))
	return line
}

// Polyline 表示SVG折线元素
type Polyline struct {
	*BaseElement
}

// NewPolyline 创建一个新的折线元素
func NewPolyline(points []types.Point) *Polyline {
	polyline := &Polyline{
		BaseElement: NewBaseElement("polyline"),
	}

	pointsStr := ""
	for i, point := range points {
		if i > 0 {
			pointsStr += " "
		}
		pointsStr += fmt.Sprintf("%f,%f", point.X, point.Y)
	}

	polyline.SetAttribute("points", pointsStr)
	return polyline
}

// Polygon 表示SVG多边形元素
type Polygon struct {
	*BaseElement
}

// NewPolygon 创建一个新的多边形元素
func NewPolygon(points []types.Point) *Polygon {
	polygon := &Polygon{
		BaseElement: NewBaseElement("polygon"),
	}

	pointsStr := ""
	for i, point := range points {
		if i > 0 {
			pointsStr += " "
		}
		pointsStr += fmt.Sprintf("%f,%f", point.X, point.Y)
	}

	polygon.SetAttribute("points", pointsStr)
	return polygon
}

// Path 表示SVG路径元素
type Path struct {
	*BaseElement
}

// NewPath 创建一个新的路径元素
func NewPath(d string) *Path {
	path := &Path{
		BaseElement: NewBaseElement("path"),
	}
	path.SetAttribute("d", d)
	return path
}

// Text 表示SVG文本元素
type Text struct {
	*BaseElement
	content string
}

// NewText 创建一个新的文本元素
func NewText(x, y float64, content string) *Text {
	text := &Text{
		BaseElement: NewBaseElement("text"),
		content:     content,
	}
	text.SetAttribute("x", fmt.Sprintf("%f", x))
	text.SetAttribute("y", fmt.Sprintf("%f", y))
	// 设置默认字体属性
	text.SetAttribute("font-family", "sans-serif")
	text.SetAttribute("font-size", "16")
	text.SetAttribute("text-anchor", "start")
	text.SetAttribute("alignment-baseline", "alphabetic")
	return text
}

// SetContent 设置文本内容
func (t *Text) SetContent(content string) {
	t.content = content
}

// GetContent 获取文本内容
func (t *Text) GetContent() string {
	return t.content
}

// SetFontFamily 设置字体族
func (t *Text) SetFontFamily(fontFamily string) {
	t.SetAttribute("font-family", fontFamily)
}

// SetFontSize 设置字体大小
func (t *Text) SetFontSize(fontSize float64) {
	t.SetAttribute("font-size", fmt.Sprintf("%f", fontSize))
}

// SetFontWeight 设置字体粗细
func (t *Text) SetFontWeight(fontWeight string) {
	t.SetAttribute("font-weight", fontWeight)
}

// SetFontStyle 设置字体样式
func (t *Text) SetFontStyle(fontStyle string) {
	t.SetAttribute("font-style", fontStyle)
}

// SetTextAnchor 设置文本锚点
func (t *Text) SetTextAnchor(textAnchor string) {
	t.SetAttribute("text-anchor", textAnchor)
}

// SetAlignmentBaseline 设置基线对齐
func (t *Text) SetAlignmentBaseline(alignmentBaseline string) {
	t.SetAttribute("alignment-baseline", alignmentBaseline)
}

// SetFill 设置填充颜色
func (t *Text) SetFill(fill string) {
	t.SetAttribute("fill", fill)
}

// SetStroke 设置描边颜色
func (t *Text) SetStroke(stroke string) {
	t.SetAttribute("stroke", stroke)
}

// SetStrokeWidth 设置描边宽度
func (t *Text) SetStrokeWidth(strokeWidth float64) {
	t.SetAttribute("stroke-width", fmt.Sprintf("%f", strokeWidth))
}

// ToXML 重写ToXML方法以包含文本内容
func (t *Text) ToXML() string {
	var sb strings.Builder

	sb.WriteString("<")
	sb.WriteString(t.Tag())

	// 添加属性
	for name, value := range t.Attributes() {
		sb.WriteString(fmt.Sprintf(` %s="%s"`, name, value))
	}

	sb.WriteString(">")
	sb.WriteString(t.content)

	// 添加子元素
	for _, child := range t.Children() {
		sb.WriteString(child.ToXML())
	}

	sb.WriteString("</")
	sb.WriteString(t.Tag())
	sb.WriteString(">")

	return sb.String()
}

// Group 表示SVG组元素
type Group struct {
	*BaseElement
}

// NewGroup 创建一个新的组元素
func NewGroup() *Group {
	return &Group{
		BaseElement: NewBaseElement("g"),
	}
}

// SVG 表示嵌套的SVG元素
type SVG struct {
	*BaseElement
}

// NewSVG 创建一个新的SVG元素
func NewSVG(x, y, width, height float64) *SVG {
	svg := &SVG{
		BaseElement: NewBaseElement("svg"),
	}
	svg.SetAttribute("x", fmt.Sprintf("%f", x))
	svg.SetAttribute("y", fmt.Sprintf("%f", y))
	svg.SetAttribute("width", fmt.Sprintf("%f", width))
	svg.SetAttribute("height", fmt.Sprintf("%f", height))
	return svg
}

// Image 表示SVG图像元素
type Image struct {
	*BaseElement
}

// NewImage 创建一个新的图像元素
func NewImage(x, y, width, height float64, href string) *Image {
	image := &Image{
		BaseElement: NewBaseElement("image"),
	}
	image.SetAttribute("x", fmt.Sprintf("%f", x))
	image.SetAttribute("y", fmt.Sprintf("%f", y))
	image.SetAttribute("width", fmt.Sprintf("%f", width))
	image.SetAttribute("height", fmt.Sprintf("%f", height))
	image.SetAttribute("href", href)
	return image
}
