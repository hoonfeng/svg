# SVG库API参考文档 / API Reference

## 📖 概述 / Overview

本文档提供了SVG库的完整API参考，包括所有公共接口、方法、参数和返回值的详细说明。

This document provides a complete API reference for the SVG library, including detailed descriptions of all public interfaces, methods, parameters, and return values.

## 🏗️ 核心结构 / Core Structures

### SVG 结构体 / SVG Struct

主要的SVG画布结构体，用于创建和管理SVG文档。

The main SVG canvas structure for creating and managing SVG documents.

```go
type SVG struct {
    Width      int                // 画布宽度 / Canvas width
    Height     int                // 画布高度 / Canvas height
    Elements   []Element          // 元素列表 / Element list
    Background color.Color        // 背景颜色 / Background color
    ViewBox    string            // 视图框 / ViewBox
    // 私有字段... / Private fields...
}
```

#### 构造函数 / Constructor

```go
// New 创建新的SVG画布
// New creates a new SVG canvas
func New(width, height int) *SVG
```

**参数 / Parameters:**
- `width int`: 画布宽度（像素）/ Canvas width in pixels
- `height int`: 画布高度（像素）/ Canvas height in pixels

**返回值 / Returns:**
- `*SVG`: SVG画布实例 / SVG canvas instance

**示例 / Example:**
```go
canvas := svg.New(800, 600)
```

### Element 接口 / Element Interface

所有SVG元素的基础接口。

Base interface for all SVG elements.

```go
type Element interface {
    // ToSVG 将元素转换为SVG字符串
    // ToSVG converts element to SVG string
    ToSVG() string
    
    // GetBounds 获取元素边界
    // GetBounds gets element bounds
    GetBounds() (x, y, width, height float64)
    
    // Clone 克隆元素
    // Clone clones the element
    Clone() Element
}
```

## 🎨 基础图形API / Basic Shapes API

### 矩形 / Rectangle

```go
// Rect 创建矩形元素
// Rect creates a rectangle element
func (s *SVG) Rect(x, y, width, height float64) *RectElement
```

**参数 / Parameters:**
- `x float64`: X坐标 / X coordinate
- `y float64`: Y坐标 / Y coordinate
- `width float64`: 宽度 / Width
- `height float64`: 高度 / Height

**返回值 / Returns:**
- `*RectElement`: 矩形元素，支持链式调用 / Rectangle element with method chaining

**链式方法 / Chaining Methods:**
```go
type RectElement struct {
    // 继承基础样式方法 / Inherits base style methods
}

// Fill 设置填充颜色
// Fill sets fill color
func (r *RectElement) Fill(color interface{}) *RectElement

// Stroke 设置描边颜色
// Stroke sets stroke color
func (r *RectElement) Stroke(color interface{}) *RectElement

// StrokeWidth 设置描边宽度
// StrokeWidth sets stroke width
func (r *RectElement) StrokeWidth(width float64) *RectElement

// Rx 设置X轴圆角半径
// Rx sets X-axis border radius
func (r *RectElement) Rx(radius float64) *RectElement

// Ry 设置Y轴圆角半径
// Ry sets Y-axis border radius
func (r *RectElement) Ry(radius float64) *RectElement

// Transform 设置变换
// Transform sets transformation
func (r *RectElement) Transform(transform string) *RectElement
```

**示例 / Example:**
```go
canvas.Rect(10, 10, 100, 50).
    Fill(color.RGBA{255, 0, 0, 255}).
    Stroke(color.RGBA{0, 0, 0, 255}).
    StrokeWidth(2).
    Rx(5).
    Ry(5)
```

### 圆形 / Circle

```go
// Circle 创建圆形元素
// Circle creates a circle element
func (s *SVG) Circle(cx, cy, r float64) *CircleElement
```

**参数 / Parameters:**
- `cx float64`: 圆心X坐标 / Center X coordinate
- `cy float64`: 圆心Y坐标 / Center Y coordinate
- `r float64`: 半径 / Radius

**返回值 / Returns:**
- `*CircleElement`: 圆形元素，支持链式调用 / Circle element with method chaining

**链式方法 / Chaining Methods:**
```go
type CircleElement struct {
    // 继承基础样式方法 / Inherits base style methods
}

// Fill 设置填充颜色
func (c *CircleElement) Fill(color interface{}) *CircleElement

// Stroke 设置描边颜色
func (c *CircleElement) Stroke(color interface{}) *CircleElement

// StrokeWidth 设置描边宽度
func (c *CircleElement) StrokeWidth(width float64) *CircleElement

// Transform 设置变换
func (c *CircleElement) Transform(transform string) *CircleElement
```

### 椭圆 / Ellipse

```go
// Ellipse 创建椭圆元素
// Ellipse creates an ellipse element
func (s *SVG) Ellipse(cx, cy, rx, ry float64) *EllipseElement
```

**参数 / Parameters:**
- `cx float64`: 椭圆中心X坐标 / Ellipse center X coordinate
- `cy float64`: 椭圆中心Y坐标 / Ellipse center Y coordinate
- `rx float64`: X轴半径 / X-axis radius
- `ry float64`: Y轴半径 / Y-axis radius

### 直线 / Line

```go
// Line 创建直线元素
// Line creates a line element
func (s *SVG) Line(x1, y1, x2, y2 float64) *LineElement
```

**参数 / Parameters:**
- `x1 float64`: 起点X坐标 / Start point X coordinate
- `y1 float64`: 起点Y坐标 / Start point Y coordinate
- `x2 float64`: 终点X坐标 / End point X coordinate
- `y2 float64`: 终点Y坐标 / End point Y coordinate

### 折线 / Polyline

```go
// Polyline 创建折线元素
// Polyline creates a polyline element
func (s *SVG) Polyline(points string) *PolylineElement
```

**参数 / Parameters:**
- `points string`: 点坐标字符串，格式："x1,y1 x2,y2 x3,y3" / Point coordinates string

### 多边形 / Polygon

```go
// Polygon 创建多边形元素
// Polygon creates a polygon element
func (s *SVG) Polygon(points string) *PolygonElement
```

**参数 / Parameters:**
- `points string`: 点坐标字符串，格式："x1,y1 x2,y2 x3,y3" / Point coordinates string

### 路径 / Path

```go
// Path 创建路径元素
// Path creates a path element
func (s *SVG) Path(d string) *PathElement
```

**参数 / Parameters:**
- `d string`: 路径数据字符串 / Path data string

**路径命令 / Path Commands:**
- `M x,y`: 移动到 / Move to
- `L x,y`: 直线到 / Line to
- `H x`: 水平线到 / Horizontal line to
- `V y`: 垂直线到 / Vertical line to
- `C x1,y1 x2,y2 x,y`: 三次贝塞尔曲线 / Cubic Bézier curve
- `Q x1,y1 x,y`: 二次贝塞尔曲线 / Quadratic Bézier curve
- `A rx,ry rotation large-arc,sweep x,y`: 椭圆弧 / Elliptical arc
- `Z`: 闭合路径 / Close path

## 📝 文本API / Text API

### 文本元素 / Text Element

```go
// Text 创建文本元素
// Text creates a text element
func (s *SVG) Text(x, y float64, content string) *TextElement
```

**参数 / Parameters:**
- `x float64`: 文本X坐标 / Text X coordinate
- `y float64`: 文本Y坐标 / Text Y coordinate
- `content string`: 文本内容 / Text content

**返回值 / Returns:**
- `*TextElement`: 文本元素，支持链式调用 / Text element with method chaining

**链式方法 / Chaining Methods:**
```go
type TextElement struct {
    // 继承基础样式方法 / Inherits base style methods
}

// FontSize 设置字体大小
// FontSize sets font size
func (t *TextElement) FontSize(size float64) *TextElement

// FontFamily 设置字体族
// FontFamily sets font family
func (t *TextElement) FontFamily(family string) *TextElement

// FontWeight 设置字体粗细
// FontWeight sets font weight
func (t *TextElement) FontWeight(weight string) *TextElement

// FontStyle 设置字体样式
// FontStyle sets font style
func (t *TextElement) FontStyle(style string) *TextElement

// TextAnchor 设置文本锚点
// TextAnchor sets text anchor
func (t *TextElement) TextAnchor(anchor string) *TextElement

// TextDecoration 设置文本装饰
// TextDecoration sets text decoration
func (t *TextElement) TextDecoration(decoration string) *TextElement

// Fill 设置文本颜色
// Fill sets text color
func (t *TextElement) Fill(color interface{}) *TextElement
```

**字体粗细值 / Font Weight Values:**
- `"100"` - `"900"`: 数字粗细 / Numeric weights
- `"normal"`: 正常粗细 / Normal weight
- `"bold"`: 粗体 / Bold weight
- `"lighter"`: 更细 / Lighter weight
- `"bolder"`: 更粗 / Bolder weight

**文本锚点值 / Text Anchor Values:**
- `"start"`: 起始对齐 / Start alignment
- `"middle"`: 居中对齐 / Middle alignment
- `"end"`: 结束对齐 / End alignment

**文本装饰值 / Text Decoration Values:**
- `"none"`: 无装饰 / No decoration
- `"underline"`: 下划线 / Underline
- `"overline"`: 上划线 / Overline
- `"line-through"`: 删除线 / Line through

## 🎨 样式API / Style API

### 颜色系统 / Color System

支持多种颜色格式：

Supports multiple color formats:

```go
// 1. color.RGBA 结构体 / color.RGBA struct
color.RGBA{255, 0, 0, 255}  // 红色 / Red

// 2. 十六进制字符串 / Hex string
"#FF0000"    // 红色 / Red
"#F00"       // 红色简写 / Red shorthand

// 3. RGB字符串 / RGB string
"rgb(255, 0, 0)"           // 红色 / Red
"rgba(255, 0, 0, 1.0)"     // 带透明度的红色 / Red with alpha

// 4. 颜色名称 / Color names
"red"        // 红色 / Red
"blue"       // 蓝色 / Blue
"green"      // 绿色 / Green
"transparent" // 透明 / Transparent
"none"       // 无颜色 / No color
```

### 描边样式 / Stroke Styles

```go
// StrokeDashArray 设置虚线样式
// StrokeDashArray sets dash pattern
func (e *Element) StrokeDashArray(pattern string) *Element

// StrokeLineCap 设置线帽样式
// StrokeLineCap sets line cap style
func (e *Element) StrokeLineCap(cap string) *Element

// StrokeLineJoin 设置线连接样式
// StrokeLineJoin sets line join style
func (e *Element) StrokeLineJoin(join string) *Element
```

**虚线模式 / Dash Patterns:**
- `"5,5"`: 5像素线段，5像素间隔 / 5px dash, 5px gap
- `"10,5,5,5"`: 复杂虚线模式 / Complex dash pattern
- `"none"`: 实线 / Solid line

**线帽样式 / Line Cap Styles:**
- `"butt"`: 平头 / Butt cap
- `"round"`: 圆头 / Round cap
- `"square"`: 方头 / Square cap

**线连接样式 / Line Join Styles:**
- `"miter"`: 尖角连接 / Miter join
- `"round"`: 圆角连接 / Round join
- `"bevel"`: 斜角连接 / Bevel join

## 🔄 变换API / Transform API

### 变换函数 / Transform Functions

```go
// Transform 设置变换
// Transform sets transformation
func (e *Element) Transform(transform string) *Element
```

**变换类型 / Transform Types:**

```go
// 平移 / Translation
"translate(50, 100)"        // 平移50,100 / Translate by 50,100
"translateX(50)"            // X轴平移 / X-axis translation
"translateY(100)"           // Y轴平移 / Y-axis translation

// 缩放 / Scaling
"scale(2)"                  // 等比缩放2倍 / Uniform scale by 2
"scale(2, 0.5)"             // X轴2倍，Y轴0.5倍 / Scale X by 2, Y by 0.5
"scaleX(2)"                 // X轴缩放 / X-axis scaling
"scaleY(0.5)"               // Y轴缩放 / Y-axis scaling

// 旋转 / Rotation
"rotate(45)"                // 绕原点旋转45度 / Rotate 45° around origin
"rotate(45 100 100)"        // 绕点(100,100)旋转45度 / Rotate 45° around (100,100)

// 倾斜 / Skewing
"skewX(30)"                 // X轴倾斜30度 / Skew X by 30°
"skewY(15)"                 // Y轴倾斜15度 / Skew Y by 15°

// 组合变换 / Combined transforms
"translate(50, 50) rotate(45) scale(1.5)"  // 多重变换 / Multiple transforms
```

## 📁 分组API / Group API

### 分组元素 / Group Element

```go
// Group 创建分组
// Group creates a group
func (s *SVG) Group() *GroupElement
```

**返回值 / Returns:**
- `*GroupElement`: 分组元素，支持链式调用 / Group element with method chaining

**链式方法 / Chaining Methods:**
```go
type GroupElement struct {
    // 继承基础样式方法 / Inherits base style methods
}

// Add 添加子元素
// Add adds child element
func (g *GroupElement) Add(element Element) *GroupElement

// Transform 设置分组变换
// Transform sets group transformation
func (g *GroupElement) Transform(transform string) *GroupElement

// Class 设置CSS类
// Class sets CSS class
func (g *GroupElement) Class(class string) *GroupElement

// ID 设置元素ID
// ID sets element ID
func (g *GroupElement) ID(id string) *GroupElement
```

**示例 / Example:**
```go
group := canvas.Group().
    Transform("translate(100, 100) rotate(45)").
    Class("my-group")

group.Add(canvas.Rect(0, 0, 50, 50).Fill("red"))
group.Add(canvas.Circle(25, 25, 10).Fill("blue"))
```

## 💾 输入输出API / I/O API

### 保存方法 / Save Methods

```go
// SaveSVG 保存为SVG文件
// SaveSVG saves as SVG file
func (s *SVG) SaveSVG(filename string) error

// SavePNG 保存为PNG文件
// SavePNG saves as PNG file
func (s *SVG) SavePNG(filename string) error

// SaveJPEG 保存为JPEG文件
// SaveJPEG saves as JPEG file
func (s *SVG) SaveJPEG(filename string, quality int) error

// ToSVGString 转换为SVG字符串
// ToSVGString converts to SVG string
func (s *SVG) ToSVGString() string

// ToImage 转换为image.Image
// ToImage converts to image.Image
func (s *SVG) ToImage() (image.Image, error)
```

**参数说明 / Parameter Description:**
- `filename string`: 文件路径 / File path
- `quality int`: JPEG质量(1-100) / JPEG quality (1-100)

### 加载方法 / Load Methods

```go
// LoadSVG 从文件加载SVG
// LoadSVG loads SVG from file
func LoadSVG(filename string) (*SVG, error)

// ParseSVG 从字符串解析SVG
// ParseSVG parses SVG from string
func ParseSVG(svgContent string) (*SVG, error)
```

## 🎬 动画API / Animation API

### 动画构建器 / Animation Builder

```go
// NewAnimationBuilder 创建动画构建器
// NewAnimationBuilder creates animation builder
func NewAnimationBuilder(width, height int) *AnimationBuilder
```

**动画构建器方法 / Animation Builder Methods:**
```go
type AnimationBuilder struct {
    // 私有字段... / Private fields...
}

// SetFrameCount 设置帧数
// SetFrameCount sets frame count
func (ab *AnimationBuilder) SetFrameCount(count int) *AnimationBuilder

// SetFrameRate 设置帧率
// SetFrameRate sets frame rate
func (ab *AnimationBuilder) SetFrameRate(fps int) *AnimationBuilder

// SetDuration 设置持续时间
// SetDuration sets duration
func (ab *AnimationBuilder) SetDuration(seconds float64) *AnimationBuilder
```

### 动画配置 / Animation Configuration

```go
type AnimationConfig struct {
    Duration   float64     // 动画持续时间(秒) / Animation duration in seconds
    Easing     EasingFunc  // 缓动函数 / Easing function
    Background color.Color // 背景颜色 / Background color
    Loop       bool        // 是否循环 / Whether to loop
}
```

### 缓动函数 / Easing Functions

```go
type EasingFunc func(t float64) float64

// 预定义缓动函数 / Predefined easing functions
var (
    Linear        EasingFunc // 线性 / Linear
    EaseIn        EasingFunc // 缓入 / Ease in
    EaseOut       EasingFunc // 缓出 / Ease out
    EaseInOut     EasingFunc // 缓入缓出 / Ease in-out
    EaseInOutQuad EasingFunc // 二次缓入缓出 / Quadratic ease in-out
)
```

### 预设动画 / Preset Animations

```go
// CreateRotatingShapes 创建旋转图形动画
// CreateRotatingShapes creates rotating shapes animation
func (ab *AnimationBuilder) CreateRotatingShapes(config AnimationConfig) *AnimationBuilder

// CreateColorfulParticles 创建彩色粒子动画
// CreateColorfulParticles creates colorful particles animation
func (ab *AnimationBuilder) CreateColorfulParticles(config AnimationConfig) *AnimationBuilder

// CreatePulsingCircles 创建脉冲圆形动画
// CreatePulsingCircles creates pulsing circles animation
func (ab *AnimationBuilder) CreatePulsingCircles(config AnimationConfig) *AnimationBuilder

// CreateWaveAnimation 创建波浪动画
// CreateWaveAnimation creates wave animation
func (ab *AnimationBuilder) CreateWaveAnimation(config AnimationConfig) *AnimationBuilder

// SaveToGIF 保存为GIF文件
// SaveToGIF saves as GIF file
func (ab *AnimationBuilder) SaveToGIF(filename string) error
```

## 🔧 高级API / Advanced API

### SVGBuilder 高级构建器 / SVGBuilder Advanced Builder

```go
// NewSVGBuilder 创建高级SVG构建器
// NewSVGBuilder creates advanced SVG builder
func NewSVGBuilder(width, height int) *SVGBuilder
```

**高级构建器方法 / Advanced Builder Methods:**
```go
type SVGBuilder struct {
    // 私有字段... / Private fields...
}

// SetBackground 设置背景
// SetBackground sets background
func (sb *SVGBuilder) SetBackground(color color.Color) *SVGBuilder

// AddRect 添加矩形
// AddRect adds rectangle
func (sb *SVGBuilder) AddRect(x, y, width, height float64) *SVGBuilder

// AddCircle 添加圆形
// AddCircle adds circle
func (sb *SVGBuilder) AddCircle(cx, cy, r float64) *SVGBuilder

// AddText 添加文本
// AddText adds text
func (sb *SVGBuilder) AddText(x, y float64, content string) *SVGBuilder

// BeginGroup 开始分组
// BeginGroup begins group
func (sb *SVGBuilder) BeginGroup() *SVGBuilder

// EndGroup 结束分组
// EndGroup ends group
func (sb *SVGBuilder) EndGroup() *SVGBuilder

// Build 构建SVG
// Build builds SVG
func (sb *SVGBuilder) Build() *SVG
```

### 字体系统 / Font System

```go
// FontMetrics 字体度量
// FontMetrics font metrics
type FontMetrics struct {
    Ascent     float64 // 上升高度 / Ascent height
    Descent    float64 // 下降高度 / Descent height
    LineHeight float64 // 行高 / Line height
    CapHeight  float64 // 大写字母高度 / Cap height
    XHeight    float64 // 小写字母高度 / X height
}

// GetFontMetrics 获取字体度量
// GetFontMetrics gets font metrics
func GetFontMetrics(fontFamily string, fontSize float64) FontMetrics

// MeasureText 测量文本尺寸
// MeasureText measures text dimensions
func MeasureText(text, fontFamily string, fontSize float64) (width, height float64)
```

## 🎯 实用工具 / Utilities

### 颜色工具 / Color Utilities

```go
// ParseColor 解析颜色字符串
// ParseColor parses color string
func ParseColor(colorStr string) (color.Color, error)

// ColorToHex 颜色转十六进制
// ColorToHex converts color to hex
func ColorToHex(c color.Color) string

// ColorToRGBA 颜色转RGBA
// ColorToRGBA converts color to RGBA
func ColorToRGBA(c color.Color) color.RGBA
```

### 几何工具 / Geometry Utilities

```go
// Point 点结构
// Point structure
type Point struct {
    X, Y float64
}

// Bounds 边界结构
// Bounds structure
type Bounds struct {
    X, Y, Width, Height float64
}

// CalculateBounds 计算元素边界
// CalculateBounds calculates element bounds
func CalculateBounds(elements []Element) Bounds

// PointInBounds 检查点是否在边界内
// PointInBounds checks if point is within bounds
func PointInBounds(point Point, bounds Bounds) bool
```

### 路径工具 / Path Utilities

```go
// PathBuilder 路径构建器
// PathBuilder path builder
type PathBuilder struct {
    // 私有字段... / Private fields...
}

// NewPathBuilder 创建路径构建器
// NewPathBuilder creates path builder
func NewPathBuilder() *PathBuilder

// MoveTo 移动到
// MoveTo moves to point
func (pb *PathBuilder) MoveTo(x, y float64) *PathBuilder

// LineTo 直线到
// LineTo draws line to point
func (pb *PathBuilder) LineTo(x, y float64) *PathBuilder

// CurveTo 曲线到
// CurveTo draws curve to point
func (pb *PathBuilder) CurveTo(x1, y1, x2, y2, x, y float64) *PathBuilder

// Close 闭合路径
// Close closes path
func (pb *PathBuilder) Close() *PathBuilder

// Build 构建路径字符串
// Build builds path string
func (pb *PathBuilder) Build() string
```

## ⚠️ 错误处理 / Error Handling

### 错误类型 / Error Types

```go
// SVGError SVG错误类型
// SVGError SVG error type
type SVGError struct {
    Code    int    // 错误代码 / Error code
    Message string // 错误消息 / Error message
    Cause   error  // 原因错误 / Cause error
}

// Error 实现error接口
// Error implements error interface
func (e *SVGError) Error() string
```

### 常见错误代码 / Common Error Codes

```go
const (
    ErrInvalidDimensions = 1001 // 无效尺寸 / Invalid dimensions
    ErrInvalidColor      = 1002 // 无效颜色 / Invalid color
    ErrFileNotFound      = 1003 // 文件未找到 / File not found
    ErrParseError        = 1004 // 解析错误 / Parse error
    ErrRenderError       = 1005 // 渲染错误 / Render error
)
```

## 📊 性能指标 / Performance Metrics

### 内存使用 / Memory Usage

- **基础SVG对象**: ~1KB / Basic SVG object: ~1KB
- **每个基础元素**: ~100-200B / Per basic element: ~100-200B
- **文本元素**: ~200-500B / Text element: ~200-500B
- **复杂路径**: ~500B-2KB / Complex path: ~500B-2KB

### 渲染性能 / Rendering Performance

- **简单图形**: <1ms / Simple shapes: <1ms
- **复杂路径**: 1-10ms / Complex paths: 1-10ms
- **大量元素**: 线性增长 / Many elements: linear growth
- **动画帧**: 10-50ms / Animation frame: 10-50ms

### 文件大小 / File Sizes

- **SVG文件**: 通常比PNG小50-80% / SVG files: typically 50-80% smaller than PNG
- **PNG输出**: 取决于复杂度和尺寸 / PNG output: depends on complexity and size
- **GIF动画**: 通常1-10MB / GIF animations: typically 1-10MB

## 🔗 相关资源 / Related Resources

- [快速入门指南](QUICK_START.md) - 库的基础使用方法
- [基础教程](BASIC_TUTORIAL.md) - 详细的功能教程
- [示例集合](EXAMPLES.md) - 丰富的代码示例
- [最佳实践指南](BEST_PRACTICES.md) - 开发建议和技巧
- [动画构建器文档](ANIMATION_BUILDER_README.md) - 高级动画功能

---

**版本信息 / Version Info**: v1.0.0  
**最后更新 / Last Updated**: 2024年12月  
**兼容性 / Compatibility**: Go 1.16+