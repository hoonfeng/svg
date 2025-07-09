package types

import (
	"testing"
)

// MockElement 用于测试的模拟元素
// MockElement is a mock element for testing
type MockElement struct {
	id         string
	attributes map[string]string
	children   []Element
	parent     Element
	tag        string
}

func NewMockElement(tag string) *MockElement {
	return &MockElement{
		tag:        tag,
		attributes: make(map[string]string),
		children:   make([]Element, 0),
	}
}

func (m *MockElement) ID() string                                              { return m.id }
func (m *MockElement) SetID(id string)                                         { m.id = id }
func (m *MockElement) GetAttributes() map[string]string                        { return m.attributes }
func (m *MockElement) SetAttribute(name, value string)                         { m.attributes[name] = value }
func (m *MockElement) GetAttribute(name string, defaultValue ...string) (string, bool) {
	value, ok := m.attributes[name]
	if !ok && len(defaultValue) > 0 {
		return defaultValue[0], true
	}
	return value, ok
}
func (m *MockElement) Children() []Element                                      { return m.children }
func (m *MockElement) AppendChild(child Element)                               { m.children = append(m.children, child) }
func (m *MockElement) ToXML() string                                            { return "<" + m.tag + "/>" }
func (m *MockElement) SetParent(parent Element)                                 { m.parent = parent }
func (m *MockElement) Clone() Element                                           { return NewMockElement(m.tag) }
func (m *MockElement) Tag() string                                              { return m.tag }

func TestNewDocument(t *testing.T) {
	doc := NewDocument(800, 600)
	if doc == nil {
		t.Fatal("Failed to create new document")
	}

	if doc.Width != "800" || doc.Height != "600" {
		t.Errorf("Document dimensions incorrect, got width=%s height=%s", doc.Width, doc.Height)
	}
}

func TestAddElement(t *testing.T) {
	doc := NewDocument(800, 600)
	circle := NewMockElement("circle")
	circle.SetID("test-circle")
	circle.SetAttribute("cx", "400")
	circle.SetAttribute("cy", "300")
	circle.SetAttribute("r", "100")

	doc.AppendElement(circle)
	if len(doc.Elements) != 1 {
		t.Error("Failed to add element to document")
	}
}

func TestFindElement(t *testing.T) {
	doc := NewDocument(800, 600)
	circle := NewMockElement("circle")
	circle.SetID("test-circle")
	circle.SetAttribute("cx", "400")
	circle.SetAttribute("cy", "300")
	circle.SetAttribute("r", "100")
	doc.AppendElement(circle)

	found := doc.FindElementByID("test-circle")
	if found == nil {
		t.Error("Failed to find element by ID")
	}
}

func TestToXML(t *testing.T) {
	doc := NewDocument(800, 600)
	circle := NewMockElement("circle")
	circle.SetAttribute("cx", "400")
	circle.SetAttribute("cy", "300")
	circle.SetAttribute("r", "100")
	doc.AppendElement(circle)

	xml := doc.ToXML()
	if len(xml) == 0 {
		t.Error("Generated XML is empty")
	}
}
