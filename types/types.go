package types

// Point 表示2D坐标点 / Represents a 2D coordinate point
type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Element SVG元素接口 / SVG element interface
type Element interface {
	// GetID 获取元素ID / Get element ID
	GetID() string

	// SetID 设置元素ID / Set element ID
	SetID(id string)

	// GetAttribute 获取属性值 / Get attribute value
	GetAttribute(name string) string

	// SetAttribute 设置属性值 / Set attribute value
	SetAttribute(name, value string)

	// GetChildren 获取子元素 / Get child elements
	GetChildren() []Element

	// AddChild 添加子元素 / Add child element
	AddChild(child Element)

	// RemoveChild 移除子元素 / Remove child element
	RemoveChild(child Element)

	// ToXML 转换为XML字符串 / Convert to XML string
	ToXML() string

	// GetParent 获取父元素 / Get parent element
	GetParent() Element

	// SetParent 设置父元素 / Set parent element
	SetParent(parent Element)

	// Clone 克隆元素 / Clone element
	Clone() Element

	// GetTagName 获取标签名 / Get tag name
	GetTagName() string
}