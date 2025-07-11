package io

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/hoonfeng/svg/types"

	"github.com/hoonfeng/svg/elements"
)

// LoadSVG 从文件加载SVG文档
func LoadSVG(filename string) (*types.Document, error) {
	// 打开文件
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 读取文件内容
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// 解析XML
	return ParseSVG(data)
}

// xmlElement XML元素结构
type xmlElement struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
	Content string     `xml:",innerxml"`
}

// ParseSVG 从XML数据解析SVG文档
func ParseSVG(data []byte) (*types.Document, error) {
	// 定义XML结构
	type xmlSVG struct {
		XMLName  xml.Name     `xml:"svg"`
		Width    string       `xml:"width,attr"`
		Height   string       `xml:"height,attr"`
		ViewBox  string       `xml:"viewBox,attr"`
		Title    string       `xml:"title"`
		Desc     string       `xml:"desc"`
		Elements []xmlElement `xml:",any"`
	}

	// 解析XML
	var xmlDoc xmlSVG
	if err := xml.Unmarshal(data, &xmlDoc); err != nil {
		return nil, err
	}

	// 创建SVG文档 - 使用默认尺寸然后设置属性
	doc := types.NewDocument(800, 600)
	if xmlDoc.Width != "" {
		doc.Width = xmlDoc.Width
	}
	if xmlDoc.Height != "" {
		doc.Height = xmlDoc.Height
	}
	if xmlDoc.ViewBox != "" {
		doc.ViewBox = xmlDoc.ViewBox
	}
	if xmlDoc.Title != "" {
		doc.SetAttribute("title", xmlDoc.Title)
	}
	if xmlDoc.Desc != "" {
		doc.SetAttribute("desc", xmlDoc.Desc)
	}

	// 解析元素
	for _, xmlEl := range xmlDoc.Elements {
		element, err := parseElement(xmlEl)
		if err != nil {
			return nil, err
		}
		if element != nil {
			doc.AppendElement(element)
		}
	}

	return doc, nil
}

// parseElement 解析单个XML元素
func parseElement(xmlEl xmlElement) (types.Element, error) {
	switch xmlEl.XMLName.Local {
	case "rect":
		return parseRect(xmlEl.Attrs)
	case "circle":
		return parseCircle(xmlEl.Attrs)
	case "ellipse":
		return parseEllipse(xmlEl.Attrs)
	case "line":
		return parseLine(xmlEl.Attrs)
	case "polyline":
		return parsePolyline(xmlEl.Attrs)
	case "polygon":
		return parsePolygon(xmlEl.Attrs)
	case "path":
		return parsePath(xmlEl.Attrs)
	case "text":
		return parseText(xmlEl.Attrs, xmlEl.Content)
	case "g":
		return parseGroup(xmlEl)
	default:
		// 忽略不支持的元素
		return nil, nil
	}
}

// parseGroup 解析组元素及其子元素
func parseGroup(xmlEl xmlElement) (*elements.Group, error) {
	// 创建组元素
	group := elements.NewGroup()

	// 设置属性
	for _, attr := range xmlEl.Attrs {
		group.SetAttribute(attr.Name.Local, attr.Value)
	}

	// 解析子元素
	type xmlRoot struct {
		Elements []xmlElement `xml:",any"`
	}
	var root xmlRoot
	if err := xml.Unmarshal([]byte("<root>"+xmlEl.Content+"</root>"), &root); err != nil {
		return nil, err
	}

	// 递归解析子元素
	for _, childEl := range root.Elements {
		childElement, err := parseElement(childEl)
		if err != nil {
			return nil, err
		}
		if childElement != nil {
			group.AppendChild(childElement)
		}
	}

	return group, nil
}

// parseRect 解析矩形元素 / Parse rectangle element
func parseRect(attrs []xml.Attr) (*elements.Rect, error) {
	var x, y, width, height float64
	var err error

	for _, attr := range attrs {
		switch attr.Name.Local {
		case "x":
			if x, err = strconv.ParseFloat(attr.Value, 64); err != nil {
				return nil, fmt.Errorf("invalid x value '%s': %v", attr.Value, err)
			}
		case "y":
			if y, err = strconv.ParseFloat(attr.Value, 64); err != nil {
				return nil, fmt.Errorf("invalid y value '%s': %v", attr.Value, err)
			}
		case "width":
			if width, err = strconv.ParseFloat(attr.Value, 64); err != nil {
				return nil, fmt.Errorf("invalid width value '%s': %v", attr.Value, err)
			}
		case "height":
			if height, err = strconv.ParseFloat(attr.Value, 64); err != nil {
				return nil, fmt.Errorf("invalid height value '%s': %v", attr.Value, err)
			}
		}
	}

	rect := elements.NewRect(x, y, width, height)

	// 设置其他属性
	for _, attr := range attrs {
		if attr.Name.Local != "x" && attr.Name.Local != "y" && attr.Name.Local != "width" && attr.Name.Local != "height" {
			rect.SetAttribute(attr.Name.Local, attr.Value)
		}
	}

	return rect, nil
}

// parseCircle 解析圆形元素 / Parse circle element
func parseCircle(attrs []xml.Attr) (*elements.Circle, error) {
	var cx, cy, r float64
	var err error

	for _, attr := range attrs {
		switch attr.Name.Local {
		case "cx":
			if cx, err = strconv.ParseFloat(attr.Value, 64); err != nil {
				return nil, fmt.Errorf("invalid cx value '%s': %v", attr.Value, err)
			}
		case "cy":
			if cy, err = strconv.ParseFloat(attr.Value, 64); err != nil {
				return nil, fmt.Errorf("invalid cy value '%s': %v", attr.Value, err)
			}
		case "r":
			if r, err = strconv.ParseFloat(attr.Value, 64); err != nil {
				return nil, fmt.Errorf("invalid r value '%s': %v", attr.Value, err)
			}
		}
	}

	circle := elements.NewCircle(cx, cy, r)

	// 设置其他属性
	for _, attr := range attrs {
		if attr.Name.Local != "cx" && attr.Name.Local != "cy" && attr.Name.Local != "r" {
			circle.SetAttribute(attr.Name.Local, attr.Value)
		}
	}

	return circle, nil
}

// parseEllipse 解析椭圆元素 / Parse ellipse element
func parseEllipse(attrs []xml.Attr) (*elements.Ellipse, error) {
	var cx, cy, rx, ry float64
	var err error

	for _, attr := range attrs {
		switch attr.Name.Local {
		case "cx":
			if cx, err = strconv.ParseFloat(attr.Value, 64); err != nil {
				return nil, fmt.Errorf("invalid cx value '%s': %v", attr.Value, err)
			}
		case "cy":
			if cy, err = strconv.ParseFloat(attr.Value, 64); err != nil {
				return nil, fmt.Errorf("invalid cy value '%s': %v", attr.Value, err)
			}
		case "rx":
			if rx, err = strconv.ParseFloat(attr.Value, 64); err != nil {
				return nil, fmt.Errorf("invalid rx value '%s': %v", attr.Value, err)
			}
		case "ry":
			if ry, err = strconv.ParseFloat(attr.Value, 64); err != nil {
				return nil, fmt.Errorf("invalid ry value '%s': %v", attr.Value, err)
			}
		}
	}

	ellipse := elements.NewEllipse(cx, cy, rx, ry)

	// 设置其他属性
	for _, attr := range attrs {
		if attr.Name.Local != "cx" && attr.Name.Local != "cy" && attr.Name.Local != "rx" && attr.Name.Local != "ry" {
			ellipse.SetAttribute(attr.Name.Local, attr.Value)
		}
	}

	return ellipse, nil
}

// parseLine 解析线段元素 / Parse line element
func parseLine(attrs []xml.Attr) (*elements.Line, error) {
	var x1, y1, x2, y2 float64
	var err error

	for _, attr := range attrs {
		switch attr.Name.Local {
		case "x1":
			if x1, err = strconv.ParseFloat(attr.Value, 64); err != nil {
				return nil, fmt.Errorf("invalid x1 value '%s': %v", attr.Value, err)
			}
		case "y1":
			if y1, err = strconv.ParseFloat(attr.Value, 64); err != nil {
				return nil, fmt.Errorf("invalid y1 value '%s': %v", attr.Value, err)
			}
		case "x2":
			if x2, err = strconv.ParseFloat(attr.Value, 64); err != nil {
				return nil, fmt.Errorf("invalid x2 value '%s': %v", attr.Value, err)
			}
		case "y2":
			if y2, err = strconv.ParseFloat(attr.Value, 64); err != nil {
				return nil, fmt.Errorf("invalid y2 value '%s': %v", attr.Value, err)
			}
		}
	}

	line := elements.NewLine(x1, y1, x2, y2)

	// 设置其他属性
	for _, attr := range attrs {
		if attr.Name.Local != "x1" && attr.Name.Local != "y1" && attr.Name.Local != "x2" && attr.Name.Local != "y2" {
			line.SetAttribute(attr.Name.Local, attr.Value)
		}
	}

	return line, nil
}

// parsePolyline 解析折线元素
func parsePolyline(attrs []xml.Attr) (*elements.Polyline, error) {
	var points []types.Point

	for _, attr := range attrs {
		if attr.Name.Local == "points" {
			points = parsePoints(attr.Value)
			break
		}
	}

	polyline := elements.NewPolyline(points)

	// 设置其他属性
	for _, attr := range attrs {
		if attr.Name.Local != "points" {
			polyline.SetAttribute(attr.Name.Local, attr.Value)
		}
	}

	return polyline, nil
}

// parsePolygon 解析多边形元素
func parsePolygon(attrs []xml.Attr) (*elements.Polygon, error) {
	var points []types.Point

	for _, attr := range attrs {
		if attr.Name.Local == "points" {
			points = parsePoints(attr.Value)
			break
		}
	}

	polygon := elements.NewPolygon(points)

	// 设置其他属性
	for _, attr := range attrs {
		if attr.Name.Local != "points" {
			polygon.SetAttribute(attr.Name.Local, attr.Value)
		}
	}

	return polygon, nil
}

// parsePath 解析路径元素
func parsePath(attrs []xml.Attr) (*elements.Path, error) {
	var d string

	for _, attr := range attrs {
		if attr.Name.Local == "d" {
			d = attr.Value
			break
		}
	}

	path := elements.NewPath(d)

	// 设置其他属性
	for _, attr := range attrs {
		if attr.Name.Local != "d" {
			path.SetAttribute(attr.Name.Local, attr.Value)
		}
	}

	return path, nil
}

// parseText 解析文本元素 / Parse text element
func parseText(attrs []xml.Attr, content string) (*elements.Text, error) {
	var x, y float64
	var err error

	for _, attr := range attrs {
		switch attr.Name.Local {
		case "x":
			if x, err = strconv.ParseFloat(attr.Value, 64); err != nil {
				return nil, fmt.Errorf("invalid x value '%s': %v", attr.Value, err)
			}
		case "y":
			if y, err = strconv.ParseFloat(attr.Value, 64); err != nil {
				return nil, fmt.Errorf("invalid y value '%s': %v", attr.Value, err)
			}
		}
	}

	text := elements.NewText(x, y, content)

	// 设置其他属性
	for _, attr := range attrs {
		if attr.Name.Local != "x" && attr.Name.Local != "y" {
			text.SetAttribute(attr.Name.Local, attr.Value)
		}
	}

	return text, nil
}

// parsePoints 解析点坐标字符串 / Parse points coordinate string
func parsePoints(pointsStr string) []types.Point {
	points := make([]types.Point, 0)

	// 分割点字符串
	pairs := strings.Fields(strings.ReplaceAll(pointsStr, ",", " "))

	for i := 0; i < len(pairs); i += 2 {
		if i+1 < len(pairs) {
			x, err := strconv.ParseFloat(pairs[i], 64)
			if err != nil {
				continue // 跳过无效的坐标
			}
			y, err := strconv.ParseFloat(pairs[i+1], 64)
			if err != nil {
				continue // 跳过无效的坐标
			}
			points = append(points, types.Point{X: x, Y: y})
		}
	}

	return points
}
