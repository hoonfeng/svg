package parser

import (
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/hoonfeng/svg/elements"
	"github.com/hoonfeng/svg/types"
)

// XMLParser 是基于encoding/xml的SVG解析器
type XMLParser struct{}

// NewXMLParser 创建一个新的XML解析器
func NewXMLParser() *XMLParser {
	return &XMLParser{}
}

// Parse 解析SVG内容
func (p *XMLParser) Parse(r io.Reader) (*types.Document, error) {
	decoder := xml.NewDecoder(r)
	doc := &types.Document{
		Elements:   make([]types.Element, 0),
		Attributes: make(map[string]string),
	}

	// 解析XML
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("XML解析错误: %v", err)
		}

		switch se := token.(type) {
		case xml.StartElement:
			if se.Name.Local == "svg" {
				// 解析SVG根元素
				for _, attr := range se.Attr {
					switch attr.Name.Local {
					case "width":
						doc.Width = attr.Value
					case "height":
						doc.Height = attr.Value
					case "viewBox":
						doc.ViewBox = attr.Value
					default:
						doc.Attributes[attr.Name.Local] = attr.Value
					}
				}

				// 解析子元素
				if err := p.parseChildren(decoder, doc); err != nil {
					return nil, err
				}
			}
		}
	}

	return doc, nil
}

// ParseString 解析SVG字符串
func (p *XMLParser) ParseString(s string) (*types.Document, error) {
	return p.Parse(strings.NewReader(s))
}

// parseChildren 解析子元素
func (p *XMLParser) parseChildren(decoder *xml.Decoder, doc *types.Document) error {
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("XML解析错误: %v", err)
		}

		switch se := token.(type) {
		case xml.StartElement:
			element, err := p.parseElement(decoder, se)
			if err != nil {
				return err
			}
			doc.Elements = append(doc.Elements, element)
		case xml.EndElement:
			if se.Name.Local == "svg" {
				return nil
			}
		}
	}

	return nil
}

// parseElement 解析单个元素
func (p *XMLParser) parseElement(decoder *xml.Decoder, start xml.StartElement) (types.Element, error) {
	var element types.Element

	// 根据标签名创建对应的元素
	switch start.Name.Local {
	case "rect":
		element = p.parseRect(start)
	case "circle":
		element = p.parseCircle(start)
	case "ellipse":
		element = p.parseEllipse(start)
	case "line":
		element = p.parseLine(start)
	case "polyline":
		element = p.parsePolyline(start)
	case "polygon":
		element = p.parsePolygon(start)
	case "path":
		element = p.parsePath(start)
	case "text":
		element = p.parseText(start)
	case "g":
		element = p.parseGroup(start)
	case "svg":
		element = p.parseSVG(start)
	case "image":
		element = p.parseImage(start)
	default:
		// 对于未知元素，创建一个基本元素
		baseElement := elements.NewBaseElement(start.Name.Local)
		for _, attr := range start.Attr {
			baseElement.SetAttribute(attr.Name.Local, attr.Value)
		}
		element = baseElement
	}

	// 解析子元素
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("XML解析错误: %v", err)
		}

		switch se := token.(type) {
		case xml.StartElement:
			child, err := p.parseElement(decoder, se)
			if err != nil {
				return nil, err
			}
			element.AppendChild(child)
		case xml.EndElement:
			if se.Name.Local == start.Name.Local {
				return element, nil
			}
		case xml.CharData:
			// 处理文本内容
			if start.Name.Local == "text" {
				if textElement, ok := element.(*elements.Text); ok {
					textElement.SetContent(string(se))
				}
			}
		}
	}

	return element, nil
}

// parseRect 解析矩形元素
func (p *XMLParser) parseRect(start xml.StartElement) *elements.Rect {
	var x, y, width, height float64

	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "x":
			x, _ = strconv.ParseFloat(attr.Value, 64)
		case "y":
			y, _ = strconv.ParseFloat(attr.Value, 64)
		case "width":
			width, _ = strconv.ParseFloat(attr.Value, 64)
		case "height":
			height, _ = strconv.ParseFloat(attr.Value, 64)
		}
	}

	rect := elements.NewRect(x, y, width, height)

	// 设置其他属性
	for _, attr := range start.Attr {
		if attr.Name.Local != "x" && attr.Name.Local != "y" && attr.Name.Local != "width" && attr.Name.Local != "height" {
			rect.SetAttribute(attr.Name.Local, attr.Value)
		}
	}

	return rect
}

// parseCircle 解析圆形元素
func (p *XMLParser) parseCircle(start xml.StartElement) *elements.Circle {
	var cx, cy, r float64

	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "cx":
			cx, _ = strconv.ParseFloat(attr.Value, 64)
		case "cy":
			cy, _ = strconv.ParseFloat(attr.Value, 64)
		case "r":
			r, _ = strconv.ParseFloat(attr.Value, 64)
		}
	}

	circle := elements.NewCircle(cx, cy, r)

	// 设置其他属性
	for _, attr := range start.Attr {
		if attr.Name.Local != "cx" && attr.Name.Local != "cy" && attr.Name.Local != "r" {
			circle.SetAttribute(attr.Name.Local, attr.Value)
		}
	}

	return circle
}

// parseEllipse 解析椭圆元素
func (p *XMLParser) parseEllipse(start xml.StartElement) *elements.Ellipse {
	var cx, cy, rx, ry float64

	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "cx":
			cx, _ = strconv.ParseFloat(attr.Value, 64)
		case "cy":
			cy, _ = strconv.ParseFloat(attr.Value, 64)
		case "rx":
			rx, _ = strconv.ParseFloat(attr.Value, 64)
		case "ry":
			ry, _ = strconv.ParseFloat(attr.Value, 64)
		}
	}

	ellipse := elements.NewEllipse(cx, cy, rx, ry)

	// 设置其他属性
	for _, attr := range start.Attr {
		if attr.Name.Local != "cx" && attr.Name.Local != "cy" && attr.Name.Local != "rx" && attr.Name.Local != "ry" {
			ellipse.SetAttribute(attr.Name.Local, attr.Value)
		}
	}

	return ellipse
}

// parseLine 解析线段元素
func (p *XMLParser) parseLine(start xml.StartElement) *elements.Line {
	var x1, y1, x2, y2 float64

	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "x1":
			x1, _ = strconv.ParseFloat(attr.Value, 64)
		case "y1":
			y1, _ = strconv.ParseFloat(attr.Value, 64)
		case "x2":
			x2, _ = strconv.ParseFloat(attr.Value, 64)
		case "y2":
			y2, _ = strconv.ParseFloat(attr.Value, 64)
		}
	}

	line := elements.NewLine(x1, y1, x2, y2)

	// 设置其他属性
	for _, attr := range start.Attr {
		if attr.Name.Local != "x1" && attr.Name.Local != "y1" && attr.Name.Local != "x2" && attr.Name.Local != "y2" {
			line.SetAttribute(attr.Name.Local, attr.Value)
		}
	}

	return line
}

// parsePolyline 解析折线元素
func (p *XMLParser) parsePolyline(start xml.StartElement) *elements.Polyline {
	var points []types.Point

	for _, attr := range start.Attr {
		if attr.Name.Local == "points" {
			points = p.parsePoints(attr.Value)
			break
		}
	}

	polyline := elements.NewPolyline(points)

	// 设置其他属性
	for _, attr := range start.Attr {
		if attr.Name.Local != "points" {
			polyline.SetAttribute(attr.Name.Local, attr.Value)
		}
	}

	return polyline
}

// parsePolygon 解析多边形元素
func (p *XMLParser) parsePolygon(start xml.StartElement) *elements.Polygon {
	var points []types.Point

	for _, attr := range start.Attr {
		if attr.Name.Local == "points" {
			points = p.parsePoints(attr.Value)
			break
		}
	}

	polygon := elements.NewPolygon(points)

	// 设置其他属性
	for _, attr := range start.Attr {
		if attr.Name.Local != "points" {
			polygon.SetAttribute(attr.Name.Local, attr.Value)
		}
	}

	return polygon
}

// parsePath 解析路径元素
func (p *XMLParser) parsePath(start xml.StartElement) *elements.Path {
	var d string

	for _, attr := range start.Attr {
		if attr.Name.Local == "d" {
			d = attr.Value
			break
		}
	}

	path := elements.NewPath(d)

	// 设置其他属性
	for _, attr := range start.Attr {
		if attr.Name.Local != "d" {
			path.SetAttribute(attr.Name.Local, attr.Value)
		}
	}

	return path
}

// parseText 解析文本元素
func (p *XMLParser) parseText(start xml.StartElement) *elements.Text {
	var x, y float64

	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "x":
			x, _ = strconv.ParseFloat(attr.Value, 64)
		case "y":
			y, _ = strconv.ParseFloat(attr.Value, 64)
		}
	}

	text := elements.NewText(x, y, "")

	// 设置其他属性
	for _, attr := range start.Attr {
		if attr.Name.Local != "x" && attr.Name.Local != "y" {
			text.SetAttribute(attr.Name.Local, attr.Value)
		}
	}

	return text
}

// parseGroup 解析组元素
func (p *XMLParser) parseGroup(start xml.StartElement) *elements.Group {
	group := elements.NewGroup()

	// 设置属性
	for _, attr := range start.Attr {
		group.SetAttribute(attr.Name.Local, attr.Value)
	}

	return group
}

// parseSVG 解析嵌套的SVG元素
func (p *XMLParser) parseSVG(start xml.StartElement) *elements.SVG {
	var x, y, width, height float64

	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "x":
			x, _ = strconv.ParseFloat(attr.Value, 64)
		case "y":
			y, _ = strconv.ParseFloat(attr.Value, 64)
		case "width":
			width, _ = strconv.ParseFloat(attr.Value, 64)
		case "height":
			height, _ = strconv.ParseFloat(attr.Value, 64)
		}
	}

	svg := elements.NewSVG(x, y, width, height)

	// 设置其他属性
	for _, attr := range start.Attr {
		if attr.Name.Local != "x" && attr.Name.Local != "y" && attr.Name.Local != "width" && attr.Name.Local != "height" {
			svg.SetAttribute(attr.Name.Local, attr.Value)
		}
	}

	return svg
}

// parseImage 解析图像元素
func (p *XMLParser) parseImage(start xml.StartElement) *elements.Image {
	var x, y, width, height float64
	var href string

	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "x":
			x, _ = strconv.ParseFloat(attr.Value, 64)
		case "y":
			y, _ = strconv.ParseFloat(attr.Value, 64)
		case "width":
			width, _ = strconv.ParseFloat(attr.Value, 64)
		case "height":
			height, _ = strconv.ParseFloat(attr.Value, 64)
		case "href", "xlink:href":
			href = attr.Value
		}
	}

	image := elements.NewImage(x, y, width, height, href)

	// 设置其他属性
	for _, attr := range start.Attr {
		if attr.Name.Local != "x" && attr.Name.Local != "y" && attr.Name.Local != "width" && attr.Name.Local != "height" && attr.Name.Local != "href" && attr.Name.Local != "xlink:href" {
			image.SetAttribute(attr.Name.Local, attr.Value)
		}
	}

	return image
}

// parsePoints 解析点字符串
func (p *XMLParser) parsePoints(pointsStr string) []types.Point {
	points := make([]types.Point, 0)

	// 分割点字符串
	pairs := strings.Fields(strings.ReplaceAll(pointsStr, ",", " "))

	for i := 0; i < len(pairs); i += 2 {
		if i+1 < len(pairs) {
			x, _ := strconv.ParseFloat(pairs[i], 64)
			y, _ := strconv.ParseFloat(pairs[i+1], 64)
			points = append(points, types.Point{X: x, Y: y})
		}
	}

	return points
}
