package types

// ElementBase 是所有SVG元素的基类
type ElementBase struct {
	id         string
	tag        string
	attributes map[string]string
	children   []Element
}

// ID 返回元素的ID
func (e *ElementBase) ID() string {
	return e.id
}

// GetAttributes 返回元素的属性
func (e *ElementBase) GetAttributes() map[string]string {
	return e.attributes
}

// Children 返回元素的子元素
func (e *ElementBase) Children() []Element {
	return e.children
}

// SetID 设置元素的ID
func (e *ElementBase) SetID(id string) {
	e.id = id
}

// SetAttributes 设置元素的属性
func (e *ElementBase) SetAttributes(name, value string) {
	e.attributes[name] = value
}

// AddChild 添加子元素
func (e *ElementBase) AddChild(child Element) {
	e.children = append(e.children, child)
}

// Tag 返回元素标签名
func (e *ElementBase) Tag() string {
	return e.tag
}

// SetTag 设置元素标签名
func (e *ElementBase) SetTag(tag string) {
	e.tag = tag
}
