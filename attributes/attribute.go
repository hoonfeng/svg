package attributes

import (
	"fmt"
	"image/color"
	"math"
	"strconv"
	"strings"
)

// ColorToHex 将Go的color.Color转换为SVG十六进制颜色字符串
func ColorToHex(c color.Color) string {
	r, g, b, a := c.RGBA()
	r, g, b, a = r>>8, g>>8, b>>8, a>>8

	if a == 255 {
		return fmt.Sprintf("#%02x%02x%02x", r, g, b)
	}

	// 如果有透明度，使用rgba()格式
	return fmt.Sprintf("rgba(%d,%d,%d,%.2f)", r, g, b, float64(a)/255.0)
}

// ParseColor 将颜色字符串解析为color.Color
func ParseColor(s string) (color.Color, error) {
	s = strings.TrimSpace(s)

	// 处理十六进制颜色
	if strings.HasPrefix(s, "#") {
		hex := s[1:]
		var r, g, b, a uint8

		switch len(hex) {
		case 3: // #RGB
			fmt.Sscanf(hex, "%1x%1x%1x", &r, &g, &b)
			r = r * 17
			g = g * 17
			b = b * 17
			a = 255
		case 6: // #RRGGBB
			fmt.Sscanf(hex, "%2x%2x%2x", &r, &g, &b)
			a = 255
		case 8: // #RRGGBBAA
			fmt.Sscanf(hex, "%2x%2x%2x%2x", &r, &g, &b, &a)
		default:
			return nil, fmt.Errorf("invalid hex color format: %s", s)
		}

		return color.RGBA{r, g, b, a}, nil
	}

	// 处理rgb()和rgba()格式
	if strings.HasPrefix(s, "rgb(") || strings.HasPrefix(s, "rgba(") {
		s = strings.TrimPrefix(s, "rgb(")
		s = strings.TrimPrefix(s, "rgba(")
		s = strings.TrimSuffix(s, ")")

		parts := strings.Split(s, ",")
		var r, g, b uint8
		var a float64 = 1.0

		if len(parts) < 3 || len(parts) > 4 {
			return nil, fmt.Errorf("invalid rgb/rgba format: %s", s)
		}

		fmt.Sscanf(parts[0], "%d", &r)
		fmt.Sscanf(parts[1], "%d", &g)
		fmt.Sscanf(parts[2], "%d", &b)

		if len(parts) == 4 {
			fmt.Sscanf(parts[3], "%f", &a)
		}

		return color.RGBA{r, g, b, uint8(a * 255)}, nil
	}

	// 处理命名颜色
	switch strings.ToLower(s) {
	case "black":
		return color.RGBA{0, 0, 0, 255}, nil
	case "white":
		return color.RGBA{255, 255, 255, 255}, nil
	case "red":
		return color.RGBA{255, 0, 0, 255}, nil
	case "green":
		return color.RGBA{0, 128, 0, 255}, nil
	case "blue":
		return color.RGBA{0, 0, 255, 255}, nil
	case "yellow":
		return color.RGBA{255, 255, 0, 255}, nil
	case "cyan", "aqua":
		return color.RGBA{0, 255, 255, 255}, nil
	case "magenta", "fuchsia":
		return color.RGBA{255, 0, 255, 255}, nil
	case "gray", "grey":
		return color.RGBA{128, 128, 128, 255}, nil
	case "silver":
		return color.RGBA{192, 192, 192, 255}, nil
	case "maroon":
		return color.RGBA{128, 0, 0, 255}, nil
	case "olive":
		return color.RGBA{128, 128, 0, 255}, nil
	case "navy":
		return color.RGBA{0, 0, 128, 255}, nil
	case "purple":
		return color.RGBA{128, 0, 128, 255}, nil
	case "teal":
		return color.RGBA{0, 128, 128, 255}, nil
	case "transparent":
		return color.RGBA{0, 0, 0, 0}, nil
	default:
		return nil, fmt.Errorf("unknown color name: %s", s)
	}
}

// Style 表示SVG样式
type Style struct {
	properties map[string]string
}

// NewStyle 创建一个新的样式
func NewStyle() *Style {
	return &Style{
		properties: make(map[string]string),
	}
}

// Set 设置样式属性
func (s *Style) Set(name, value string) {
	s.properties[name] = value
}

// Get 获取样式属性
func (s *Style) Get(name string) (string, bool) {
	value, ok := s.properties[name]
	return value, ok
}

// Remove 移除样式属性
func (s *Style) Remove(name string) {
	delete(s.properties, name)
}

// SetFill 设置填充颜色
func (s *Style) SetFill(c color.Color) {
	s.Set("fill", ColorToHex(c))
}

// SetStroke 设置描边颜色
func (s *Style) SetStroke(c color.Color) {
	s.Set("stroke", ColorToHex(c))
}

// SetStrokeWidth 设置描边宽度
func (s *Style) SetStrokeWidth(width float64) {
	s.Set("stroke-width", fmt.Sprintf("%f", width))
}

// SetOpacity 设置不透明度
func (s *Style) SetOpacity(opacity float64) {
	s.Set("opacity", fmt.Sprintf("%f", opacity))
}

// SetFontFamily 设置字体族
func (s *Style) SetFontFamily(family string) {
	s.Set("font-family", family)
}

// SetFontSize 设置字体大小
func (s *Style) SetFontSize(size float64) {
	s.Set("font-size", fmt.Sprintf("%fpx", size))
}

// SetFontWeight 设置字体粗细
func (s *Style) SetFontWeight(weight string) {
	s.Set("font-weight", weight)
}

// SetTextAnchor 设置文本锚点
func (s *Style) SetTextAnchor(anchor string) {
	s.Set("text-anchor", anchor)
}

// ToString 将样式转换为字符串
func (s *Style) ToString() string {
	var parts []string
	for name, value := range s.properties {
		parts = append(parts, fmt.Sprintf("%s: %s", name, value))
	}
	return strings.Join(parts, "; ")
}

// Matrix 表示2D变换矩阵
type Matrix struct {
	A, B, C, D, E, F float64
}

// Transform 表示SVG变换
type Transform struct {
	operations []string
	matrix     *Matrix // 缓存的矩阵表示
}

// NewTransform 创建一个新的变换
func NewTransform() *Transform {
	return &Transform{
		operations: make([]string, 0),
		matrix:     &Matrix{A: 1, B: 0, C: 0, D: 1, E: 0, F: 0}, // 单位矩阵
	}
}

// Translate 添加平移变换
func (t *Transform) Translate(x, y float64) *Transform {
	t.operations = append(t.operations, fmt.Sprintf("translate(%f,%f)", x, y))
	return t
}

// Scale 添加缩放变换
func (t *Transform) Scale(x, y float64) *Transform {
	t.operations = append(t.operations, fmt.Sprintf("scale(%f,%f)", x, y))
	return t
}

// Rotate 添加旋转变换
func (t *Transform) Rotate(angle float64) *Transform {
	t.operations = append(t.operations, fmt.Sprintf("rotate(%f)", angle))
	return t
}

// RotateAround 添加围绕指定点的旋转变换
func (t *Transform) RotateAround(angle, x, y float64) *Transform {
	t.operations = append(t.operations, fmt.Sprintf("rotate(%f,%f,%f)", angle, x, y))
	return t
}

// SkewX 添加X轴倾斜变换
func (t *Transform) SkewX(angle float64) *Transform {
	t.operations = append(t.operations, fmt.Sprintf("skewX(%f)", angle))
	return t
}

// SkewY 添加Y轴倾斜变换
func (t *Transform) SkewY(angle float64) *Transform {
	t.operations = append(t.operations, fmt.Sprintf("skewY(%f)", angle))
	return t
}

// Matrix 添加矩阵变换
func (t *Transform) Matrix(a, b, c, d, e, f float64) *Transform {
	t.operations = append(t.operations, fmt.Sprintf("matrix(%f,%f,%f,%f,%f,%f)", a, b, c, d, e, f))
	return t
}

// ToString 将变换转换为字符串
func (t *Transform) ToString() string {
	return strings.Join(t.operations, " ")
}

// GetMatrix 获取变换的矩阵表示
func (t *Transform) GetMatrix() *Matrix {
	if t.matrix == nil {
		t.matrix = &Matrix{A: 1, B: 0, C: 0, D: 1, E: 0, F: 0} // 单位矩阵
	}
	
	// 重新计算矩阵
	result := &Matrix{A: 1, B: 0, C: 0, D: 1, E: 0, F: 0} // 单位矩阵
	
	// 解析并应用所有变换操作
	for _, operation := range t.operations {
		opMatrix := parseTransformOperation(operation)
		if opMatrix != nil {
			result = multiplyMatrices(result, opMatrix)
		}
	}
	
	t.matrix = result
	return t.matrix
}

// SetMatrix 设置变换的矩阵
func (t *Transform) SetMatrix(m *Matrix) {
	t.matrix = m
}

// parseTransformOperation 解析单个变换操作并返回对应的矩阵
func parseTransformOperation(operation string) *Matrix {
	// 去除空格
	operation = strings.TrimSpace(operation)
	
	// 解析translate(x,y)
	if strings.HasPrefix(operation, "translate(") {
		params := extractParams(operation, "translate")
		if len(params) >= 2 {
			x, y := params[0], params[1]
			return &Matrix{A: 1, B: 0, C: 0, D: 1, E: x, F: y}
		}
	}
	
	// 解析scale(x,y)
	if strings.HasPrefix(operation, "scale(") {
		params := extractParams(operation, "scale")
		if len(params) >= 2 {
			sx, sy := params[0], params[1]
			return &Matrix{A: sx, B: 0, C: 0, D: sy, E: 0, F: 0}
		} else if len(params) == 1 {
			s := params[0]
			return &Matrix{A: s, B: 0, C: 0, D: s, E: 0, F: 0}
		}
	}
	
	// 解析rotate(angle) 或 rotate(angle,cx,cy)
	if strings.HasPrefix(operation, "rotate(") {
		params := extractParams(operation, "rotate")
		if len(params) >= 1 {
			angle := params[0] * math.Pi / 180 // 转换为弧度
			cos := math.Cos(angle)
			sin := math.Sin(angle)
			
			if len(params) >= 3 {
				// 围绕指定点旋转
				cx, cy := params[1], params[2]
				// 先平移到原点，旋转，再平移回去
				// T(cx,cy) * R(angle) * T(-cx,-cy)
				return &Matrix{
					A: cos, B: sin, C: -sin, D: cos,
					E: cx*(1-cos) + cy*sin, F: cy*(1-cos) - cx*sin,
				}
			} else {
				// 围绕原点旋转
				return &Matrix{A: cos, B: sin, C: -sin, D: cos, E: 0, F: 0}
			}
		}
	}
	
	// 解析skewX(angle)
	if strings.HasPrefix(operation, "skewX(") {
		params := extractParams(operation, "skewX")
		if len(params) >= 1 {
			angle := params[0] * math.Pi / 180 // 转换为弧度
			tan := math.Tan(angle)
			return &Matrix{A: 1, B: 0, C: tan, D: 1, E: 0, F: 0}
		}
	}
	
	// 解析skewY(angle)
	if strings.HasPrefix(operation, "skewY(") {
		params := extractParams(operation, "skewY")
		if len(params) >= 1 {
			angle := params[0] * math.Pi / 180 // 转换为弧度
			tan := math.Tan(angle)
			return &Matrix{A: 1, B: tan, C: 0, D: 1, E: 0, F: 0}
		}
	}
	
	// 解析matrix(a,b,c,d,e,f)
	if strings.HasPrefix(operation, "matrix(") {
		params := extractParams(operation, "matrix")
		if len(params) >= 6 {
			return &Matrix{
				A: params[0], B: params[1], C: params[2],
				D: params[3], E: params[4], F: params[5],
			}
		}
	}
	
	return nil
}

// extractParams 从变换操作字符串中提取参数
func extractParams(operation, funcName string) []float64 {
	// 移除函数名和括号
	start := len(funcName) + 1 // +1 for '('
	end := strings.LastIndex(operation, ")")
	if end == -1 {
		return nil
	}
	
	paramsStr := operation[start:end]
	// 替换逗号为空格
	paramsStr = strings.ReplaceAll(paramsStr, ",", " ")
	// 分割参数
	parts := strings.Fields(paramsStr)
	
	params := make([]float64, 0, len(parts))
	for _, part := range parts {
		if val, err := strconv.ParseFloat(part, 64); err == nil {
			params = append(params, val)
		}
	}
	
	return params
}

// multiplyMatrices 矩阵乘法
func multiplyMatrices(m1, m2 *Matrix) *Matrix {
	return &Matrix{
		A: m1.A*m2.A + m1.B*m2.C,
		B: m1.A*m2.B + m1.B*m2.D,
		C: m1.C*m2.A + m1.D*m2.C,
		D: m1.C*m2.B + m1.D*m2.D,
		E: m1.E*m2.A + m1.F*m2.C + m2.E,
		F: m1.E*m2.B + m1.F*m2.D + m2.F,
	}
}

// Gradient 表示SVG渐变的基础结构
type Gradient struct {
	ID       string
	Stops    []GradientStop
	GradType string // "linear" 或 "radial"
	Attrs    map[string]string
}

// GradientStop 表示渐变中的一个颜色停止点
type GradientStop struct {
	Offset  float64
	Color   color.Color
	Opacity float64
}

// NewLinearGradient 创建一个新的线性渐变
func NewLinearGradient(id string, x1, y1, x2, y2 float64) *Gradient {
	g := &Gradient{
		ID:       id,
		Stops:    make([]GradientStop, 0),
		GradType: "linear",
		Attrs:    make(map[string]string),
	}

	g.Attrs["x1"] = fmt.Sprintf("%f", x1)
	g.Attrs["y1"] = fmt.Sprintf("%f", y1)
	g.Attrs["x2"] = fmt.Sprintf("%f", x2)
	g.Attrs["y2"] = fmt.Sprintf("%f", y2)

	return g
}

// NewRadialGradient 创建一个新的径向渐变
func NewRadialGradient(id string, cx, cy, r, fx, fy float64) *Gradient {
	g := &Gradient{
		ID:       id,
		Stops:    make([]GradientStop, 0),
		GradType: "radial",
		Attrs:    make(map[string]string),
	}

	g.Attrs["cx"] = fmt.Sprintf("%f", cx)
	g.Attrs["cy"] = fmt.Sprintf("%f", cy)
	g.Attrs["r"] = fmt.Sprintf("%f", r)
	g.Attrs["fx"] = fmt.Sprintf("%f", fx)
	g.Attrs["fy"] = fmt.Sprintf("%f", fy)

	return g
}

// AddStop 添加一个渐变停止点
func (g *Gradient) AddStop(offset float64, c color.Color, opacity float64) {
	g.Stops = append(g.Stops, GradientStop{
		Offset:  offset,
		Color:   c,
		Opacity: opacity,
	})
}

// ToXML 将渐变转换为XML字符串
func (g *Gradient) ToXML() string {
	var sb strings.Builder

	if g.GradType == "linear" {
		sb.WriteString(fmt.Sprintf("<linearGradient id=\"%s\"", g.ID))
	} else {
		sb.WriteString(fmt.Sprintf("<radialGradient id=\"%s\"", g.ID))
	}

	// 添加属性
	for name, value := range g.Attrs {
		sb.WriteString(fmt.Sprintf(" %s=\"%s\"", name, value))
	}

	sb.WriteString(">")

	// 添加停止点
	for _, stop := range g.Stops {
		sb.WriteString(fmt.Sprintf("<stop offset=\"%f\" stop-color=\"%s\" stop-opacity=\"%f\" />",
			stop.Offset, ColorToHex(stop.Color), stop.Opacity))
	}

	if g.GradType == "linear" {
		sb.WriteString("</linearGradient>")
	} else {
		sb.WriteString("</radialGradient>")
	}

	return sb.String()
}

// Filter 表示SVG滤镜
type Filter struct {
	ID       string
	Elements []string
	Attrs    map[string]string
}

// NewFilter 创建一个新的滤镜
func NewFilter(id string) *Filter {
	return &Filter{
		ID:       id,
		Elements: make([]string, 0),
		Attrs:    make(map[string]string),
	}
}

// AddElement 添加一个滤镜元素
func (f *Filter) AddElement(element string) {
	f.Elements = append(f.Elements, element)
}

// SetAttribute 设置滤镜属性
func (f *Filter) SetAttribute(name, value string) {
	f.Attrs[name] = value
}

// ToXML 将滤镜转换为XML字符串
func (f *Filter) ToXML() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("<filter id=\"%s\"", f.ID))

	// 添加属性
	for name, value := range f.Attrs {
		sb.WriteString(fmt.Sprintf(" %s=\"%s\"", name, value))
	}

	sb.WriteString(">")

	// 添加滤镜元素
	for _, element := range f.Elements {
		sb.WriteString(element)
	}

	sb.WriteString("</filter>")

	return sb.String()
}
