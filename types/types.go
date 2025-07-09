package types

// Point 表示二维坐标点
type Point struct {
	X float64
	Y float64
}

// Element 表示SVG元素的接口
// Element interface represents an SVG element
type Element interface {
	// ID 返回元素的ID
	// ID returns the element's ID
	ID() string

	// SetID 设置元素的ID
	// SetID sets the element's ID
	SetID(id string)

	// GetAttributes 返回元素的属性
	// GetAttributes returns the element's attributes
	GetAttributes() map[string]string

	// SetAttribute 设置元素的属性
	// SetAttribute sets an attribute of the element
	SetAttribute(name, value string)

	// GetAttribute 获取元素的属性
	// GetAttribute gets an attribute of the element
	GetAttribute(name string, defaultValue ...string) (string, bool)

	// Children 返回元素的子元素
	// Children returns the element's children
	Children() []Element

	// AppendChild 添加子元素
	// AppendChild adds a child element
	AppendChild(child Element)

	// ToXML 将元素转换为XML字符串
	// ToXML converts the element to XML string
	ToXML() string

	// SetParent 设置父元素
	// SetParent sets the parent element
	SetParent(parent Element)

	// Clone 克隆元素
	// Clone clones the element
	Clone() Element

	// Tag 返回元素标签名
	// Tag returns the element's tag name
	Tag() string
}
