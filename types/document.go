package types

import (
	"fmt"
	"io"
	"strings"
)

// Document 表示SVG文档
type Document struct {
	Width      string
	Height     string
	ViewBox    string
	Elements   []Element
	Attributes map[string]string
	Defs       []Element // 定义区域中的元素
}

// NewDocument 创建一个新的SVG文档
func NewDocument(width, height int) *Document {
	return &Document{
		Width:      fmt.Sprintf("%d", width),
		Height:     fmt.Sprintf("%d", height),
		Elements:   make([]Element, 0),
		Attributes: make(map[string]string),
		Defs:       make([]Element, 0),
	}
}

// SetViewBox 设置视图框
func (d *Document) SetViewBox(minX, minY, width, height float64) {
	d.ViewBox = fmt.Sprintf("%f %f %f %f", minX, minY, width, height)
}

// SetAttribute 设置文档属性
func (d *Document) SetAttribute(name, value string) {
	d.Attributes[name] = value
}

// GetAttribute 获取文档属性
func (d *Document) GetAttribute(name string) (string, bool) {
	value, ok := d.Attributes[name]
	return value, ok
}

// AppendElement 添加元素到文档
func (d *Document) AppendElement(element Element) {
	d.Elements = append(d.Elements, element)
}

// AddDef 添加元素到定义区域
func (d *Document) AddDef(element Element) {
	d.Defs = append(d.Defs, element)
}

// WriteTo 将SVG文档写入io.Writer
func (d *Document) WriteTo(w io.Writer) error {
	// 写入XML声明和DOCTYPE
	if _, err := io.WriteString(w, "<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"no\"?>\n"); err != nil {
		return err
	}
	if _, err := io.WriteString(w, "<!DOCTYPE svg PUBLIC \"-//W3C//DTD SVG 1.1//EN\" \"http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd\">\n"); err != nil {
		return err
	}

	// 写入SVG根元素开始标签
	if _, err := io.WriteString(w, "<svg"); err != nil {
		return err
	}

	// 写入命名空间
	if _, err := io.WriteString(w, " xmlns=\"http://www.w3.org/2000/svg\""); err != nil {
		return err
	}
	if _, err := io.WriteString(w, " xmlns:xlink=\"http://www.w3.org/1999/xlink\""); err != nil {
		return err
	}

	// 写入宽度和高度
	if _, err := io.WriteString(w, fmt.Sprintf(" width=\"%s\" height=\"%s\"", d.Width, d.Height)); err != nil {
		return err
	}

	// 写入视图框
	if d.ViewBox != "" {
		if _, err := io.WriteString(w, fmt.Sprintf(" viewBox=\"%s\"", d.ViewBox)); err != nil {
			return err
		}
	}

	// 写入其他属性
	for name, value := range d.Attributes {
		if _, err := io.WriteString(w, fmt.Sprintf(" %s=\"%s\"", name, value)); err != nil {
			return err
		}
	}

	// 结束开始标签
	if _, err := io.WriteString(w, ">\n"); err != nil {
		return err
	}

	// 写入定义区域
	if len(d.Defs) > 0 {
		if _, err := io.WriteString(w, "<defs>\n"); err != nil {
			return err
		}

		for _, def := range d.Defs {
			if _, err := io.WriteString(w, def.ToXML()); err != nil {
				return err
			}
			if _, err := io.WriteString(w, "\n"); err != nil {
				return err
			}
		}

		if _, err := io.WriteString(w, "</defs>\n"); err != nil {
			return err
		}
	}

	// 写入所有元素
	for _, element := range d.Elements {
		if _, err := io.WriteString(w, element.ToXML()); err != nil {
			return err
		}
		if _, err := io.WriteString(w, "\n"); err != nil {
			return err
		}
	}

	// 写入SVG结束标签
	if _, err := io.WriteString(w, "</svg>\n"); err != nil {
		return err
	}

	return nil
}

// ToXML 将SVG文档转换为XML字符串
func (d *Document) ToXML() string {
	var sb strings.Builder

	// XML声明和DOCTYPE
	sb.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"no\"?>\n")
	sb.WriteString("<!DOCTYPE svg PUBLIC \"-//W3C//DTD SVG 1.1//EN\" \"http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd\">\n")

	// SVG根元素开始标签
	sb.WriteString("<svg")

	// 命名空间
	sb.WriteString(" xmlns=\"http://www.w3.org/2000/svg\"")
	sb.WriteString(" xmlns:xlink=\"http://www.w3.org/1999/xlink\"")

	// 宽度和高度
	sb.WriteString(fmt.Sprintf(" width=\"%s\" height=\"%s\"", d.Width, d.Height))

	// 视图框
	if d.ViewBox != "" {
		sb.WriteString(fmt.Sprintf(" viewBox=\"%s\"", d.ViewBox))
	}

	// 其他属性
	for name, value := range d.Attributes {
		sb.WriteString(fmt.Sprintf(" %s=\"%s\"", name, value))
	}

	// 结束开始标签
	sb.WriteString(">\n")

	// 定义区域
	if len(d.Defs) > 0 {
		sb.WriteString("<defs>\n")

		for _, def := range d.Defs {
			sb.WriteString(def.ToXML())
			sb.WriteString("\n")
		}

		sb.WriteString("</defs>\n")
	}

	// 所有元素
	for _, element := range d.Elements {
		sb.WriteString(element.ToXML())
		sb.WriteString("\n")
	}

	// SVG结束标签
	sb.WriteString("</svg>\n")

	return sb.String()
}

// FindElementByID 通过ID查找元素
func (d *Document) FindElementByID(id string) Element {
	return findElementByID(d.Elements, id)
}

// findElementByID 递归查找具有指定ID的元素
func findElementByID(elements []Element, id string) Element {
	for _, element := range elements {
		if element.ID() == id {
			return element
		}

		// 递归查找子元素
		if found := findElementByID(element.Children(), id); found != nil {
			return found
		}
	}

	return nil
}
