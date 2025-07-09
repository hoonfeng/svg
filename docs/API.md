# API 参考文档 / API Reference

本文档提供了SVG库的完整API参考。

This document provides a complete API reference for the SVG library.

## 目录 / Table of Contents

- [核心API / Core API](#核心api--core-api)
- [图形元素 / Graphic Elements](#图形元素--graphic-elements)
- [文本处理 / Text Processing](#文本处理--text-processing)
- [样式系统 / Style System](#样式系统--style-system)
- [动画系统 / Animation System](#动画系统--animation-system)
- [渲染器 / Renderer](#渲染器--renderer)
- [路径处理 / Path Processing](#路径处理--path-processing)
- [变换 / Transforms](#变换--transforms)

## 核心API / Core API

### SVG 结构体 / SVG Struct

```go
type SVG struct {
    Width  float64
    Height float64
    // 其他字段...
}
```

#### 创建新的SVG / Create New SVG

```go
// New 创建一个新的SVG实例 / Creates a new SVG instance
func New(width, height float64) *SVG
```

**参数 / Parameters:**
- `width`: SVG画布宽度 / SVG canvas width
- `height`: SVG画布高度 / SVG canvas height

**返回值 / Returns:**
- `*SVG`: 新的SVG实例 / New SVG instance

**示例 / Example:**
```go
svg := svg.New(800, 600)
```

#### 保存SVG / Save SVG

```go
// Save 将SVG保存到文件 / Saves SVG to file
func (s *SVG) Save(filename string) error
```

**参数 / Parameters:**
- `filename`: 文件名 / Filename

**返回值 / Returns:**
- `error`: 错误信息（如果有）/ Error (if any)

**示例 / Example:**
```go
err := svg.Save("output.svg")
if err != nil {
    log.Fatal(err)
}
```

## 图形元素 / Graphic Elements

### 圆形 / Circle

```go
// Circle 创建圆形元素 / Creates a circle element
func (s *SVG) Circle(x, y, r float64) *CircleElement
```

**参数 / Parameters:**
- `x`: 圆心X坐标 / Center X coordinate
- `y`: 圆心Y坐标 / Center Y coordinate
- `r`: 半径 / Radius

**返回值 / Returns:**
- `*CircleElement`: 圆形元素 / Circle element

**示例 / Example:**
```go
circle := svg.Circle(100, 100, 50)
circle.Fill("red").Stroke("black", 2)
```

### 矩形 / Rectangle

```go
// Rect 创建矩形元素 / Creates a rectangle element
func (s *SVG) Rect(x, y, width, height float64) *RectElement
```

**参数 / Parameters:**
- `x`: 左上角X坐标 / Top-left X coordinate
- `y`: 左上角Y坐标 / Top-left Y coordinate
- `width`: 宽度 / Width
- `height`: 高度 / Height

**返回值 / Returns:**
- `*RectElement`: 矩形元素 / Rectangle element

**示例 / Example:**
```go
rect := svg.Rect(50, 50, 200, 100)
rect.Fill("blue").Stroke("black", 1)
```

### 线条 / Line

```go
// Line 创建线条元素 / Creates a line element
func (s *SVG) Line(x1, y1, x2, y2 float64) *LineElement
```

**参数 / Parameters:**
- `x1`: 起点X坐标 / Start X coordinate
- `y1`: 起点Y坐标 / Start Y coordinate
- `x2`: 终点X坐标 / End X coordinate
- `y2`: 终点Y坐标 / End Y coordinate

**返回值 / Returns:**
- `*LineElement`: 线条元素 / Line element

**示例 / Example:**
```go
line := svg.Line(0, 0, 100, 100)
line.Stroke("black", 2)
```

### 椭圆 / Ellipse

```go
// Ellipse 创建椭圆元素 / Creates an ellipse element
func (s *SVG) Ellipse(cx, cy, rx, ry float64) *EllipseElement
```

**参数 / Parameters:**
- `cx`: 中心X坐标 / Center X coordinate
- `cy`: 中心Y坐标 / Center Y coordinate
- `rx`: X轴半径 / X-axis radius
- `ry`: Y轴半径 / Y-axis radius

**返回值 / Returns:**
- `*EllipseElement`: 椭圆元素 / Ellipse element

**示例 / Example:**
```go
ellipse := svg.Ellipse(150, 100, 80, 50)
ellipse.Fill("green").Stroke("black", 1)
```

## 文本处理 / Text Processing

### 文本元素 / Text Element

```go
// Text 创建文本元素 / Creates a text element
func (s *SVG) Text(x, y float64, content string) *TextElement
```

**参数 / Parameters:**
- `x`: 文本X坐标 / Text X coordinate
- `y`: 文本Y坐标 / Text Y coordinate
- `content`: 文本内容 / Text content

**返回值 / Returns:**
- `*TextElement`: 文本元素 / Text element

**示例 / Example:**
```go
text := svg.Text(100, 50, "Hello, SVG!")
text.FontFamily("Arial").FontSize(24).Fill("black")
```

### 文本样式方法 / Text Style Methods

```go
// FontFamily 设置字体族 / Sets font family
func (t *TextElement) FontFamily(family string) *TextElement

// FontSize 设置字体大小 / Sets font size
func (t *TextElement) FontSize(size float64) *TextElement

// FontWeight 设置字体粗细 / Sets font weight
func (t *TextElement) FontWeight(weight string) *TextElement

// FontStyle 设置字体样式 / Sets font style
func (t *TextElement) FontStyle(style string) *TextElement
```

**示例 / Example:**
```go
text := svg.Text(100, 100, "样式文本")
text.FontFamily("Microsoft YaHei")
    .FontSize(20)
    .FontWeight("bold")
    .FontStyle("italic")
    .Fill("red")
```

## 样式系统 / Style System

### 填充和描边 / Fill and Stroke

```go
// Fill 设置填充颜色 / Sets fill color
func (e *Element) Fill(color string) *Element

// Stroke 设置描边 / Sets stroke
func (e *Element) Stroke(color string, width float64) *Element

// StrokeWidth 设置描边宽度 / Sets stroke width
func (e *Element) StrokeWidth(width float64) *Element

// StrokeDashArray 设置虚线样式 / Sets dash array
func (e *Element) StrokeDashArray(dashArray string) *Element
```

**颜色格式 / Color Formats:**
- 颜色名称 / Color names: `"red"`, `"blue"`, `"green"`
- 十六进制 / Hex: `"#FF0000"`, `"#00FF00"`
- RGB: `"rgb(255, 0, 0)"`
- RGBA: `"rgba(255, 0, 0, 0.5)"`

**示例 / Example:**
```go
circle := svg.Circle(100, 100, 50)
circle.Fill("rgba(255, 0, 0, 0.7)")
      .Stroke("#000000", 2)
      .StrokeDashArray("5,5")
```

### 透明度 / Opacity

```go
// Opacity 设置透明度 / Sets opacity
func (e *Element) Opacity(opacity float64) *Element

// FillOpacity 设置填充透明度 / Sets fill opacity
func (e *Element) FillOpacity(opacity float64) *Element

// StrokeOpacity 设置描边透明度 / Sets stroke opacity
func (e *Element) StrokeOpacity(opacity float64) *Element
```

**示例 / Example:**
```go
rect := svg.Rect(50, 50, 100, 100)
rect.Fill("red").Opacity(0.5)
```

## 动画系统 / Animation System

### 动画构建器 / Animation Builder

```go
// NewAnimationBuilder 创建动画构建器 / Creates animation builder
func NewAnimationBuilder(width, height int) *AnimationBuilder

// AddFrame 添加帧 / Adds frame
func (ab *AnimationBuilder) AddFrame(svg *SVG, duration time.Duration) *AnimationBuilder

// SaveGIF 保存为GIF / Saves as GIF
func (ab *AnimationBuilder) SaveGIF(filename string) error
```

**示例 / Example:**
```go
builder := animation.NewAnimationBuilder(400, 400)

for i := 0; i < 60; i++ {
    svg := svg.New(400, 400)
    angle := float64(i) * 6 // 每帧旋转6度
    
    circle := svg.Circle(200, 200, 50)
    circle.Fill("red")
    circle.Transform(fmt.Sprintf("rotate(%f 200 200)", angle))
    
    builder.AddFrame(svg, 50*time.Millisecond)
}

builder.SaveGIF("rotation.gif")
```

### 动画属性 / Animation Properties

```go
// Animate 添加动画 / Adds animation
func (e *Element) Animate(attributeName string, values []string, duration time.Duration) *Element

// AnimateTransform 添加变换动画 / Adds transform animation
func (e *Element) AnimateTransform(transformType string, values []string, duration time.Duration) *Element
```

**示例 / Example:**
```go
circle := svg.Circle(100, 100, 50)
circle.Fill("red")
circle.Animate("r", []string{"10", "50", "10"}, 2*time.Second)
```

## 渲染器 / Renderer

### PNG渲染 / PNG Rendering

```go
// RenderToPNG 渲染为PNG / Renders to PNG
func (s *SVG) RenderToPNG(filename string, width, height int) error

// RenderToPNGWithDPI 以指定DPI渲染为PNG / Renders to PNG with specified DPI
func (s *SVG) RenderToPNGWithDPI(filename string, dpi float64) error
```

**参数 / Parameters:**
- `filename`: 输出文件名 / Output filename
- `width`: 输出宽度（像素）/ Output width (pixels)
- `height`: 输出高度（像素）/ Output height (pixels)
- `dpi`: 分辨率 / Resolution

**示例 / Example:**
```go
// 渲染为800x600的PNG
svg.RenderToPNG("output.png", 800, 600)

// 以300 DPI渲染
svg.RenderToPNGWithDPI("high_res.png", 300)
```

### 渲染选项 / Render Options

```go
type RenderOptions struct {
    Width       int     // 输出宽度 / Output width
    Height      int     // 输出高度 / Output height
    DPI         float64 // 分辨率 / Resolution
    Background  string  // 背景颜色 / Background color
    Quality     int     // 质量 (1-100) / Quality (1-100)
}

// RenderWithOptions 使用选项渲染 / Renders with options
func (s *SVG) RenderWithOptions(filename string, options RenderOptions) error
```

**示例 / Example:**
```go
options := RenderOptions{
    Width:      1920,
    Height:     1080,
    DPI:        300,
    Background: "white",
    Quality:    95,
}

svg.RenderWithOptions("high_quality.png", options)
```

## 路径处理 / Path Processing

### 路径元素 / Path Element

```go
// Path 创建路径元素 / Creates path element
func (s *SVG) Path(d string) *PathElement
```

**参数 / Parameters:**
- `d`: 路径数据 / Path data

**返回值 / Returns:**
- `*PathElement`: 路径元素 / Path element

**路径命令 / Path Commands:**
- `M x,y`: 移动到 / Move to
- `L x,y`: 直线到 / Line to
- `H x`: 水平线到 / Horizontal line to
- `V y`: 垂直线到 / Vertical line to
- `C x1,y1 x2,y2 x,y`: 三次贝塞尔曲线 / Cubic Bézier curve
- `Q x1,y1 x,y`: 二次贝塞尔曲线 / Quadratic Bézier curve
- `A rx,ry rotation large-arc-flag,sweep-flag x,y`: 椭圆弧 / Elliptical arc
- `Z`: 闭合路径 / Close path

**示例 / Example:**
```go
// 创建一个心形
heartPath := "M12,21.35l-1.45-1.32C5.4,15.36,2,12.28,2,8.5 C2,5.42,4.42,3,7.5,3c1.74,0,3.41,0.81,4.5,2.09C13.09,3.81,14.76,3,16.5,3 C19.58,3,22,5.42,22,8.5c0,3.78-3.4,6.86-8.55,11.54L12,21.35z"
path := svg.Path(heartPath)
path.Fill("red").Stroke("black", 1)
```

### 路径构建器 / Path Builder

```go
type PathBuilder struct {
    // 内部字段...
}

// NewPathBuilder 创建路径构建器 / Creates path builder
func NewPathBuilder() *PathBuilder

// MoveTo 移动到指定点 / Moves to specified point
func (pb *PathBuilder) MoveTo(x, y float64) *PathBuilder

// LineTo 画直线到指定点 / Draws line to specified point
func (pb *PathBuilder) LineTo(x, y float64) *PathBuilder

// CurveTo 画三次贝塞尔曲线 / Draws cubic Bézier curve
func (pb *PathBuilder) CurveTo(x1, y1, x2, y2, x, y float64) *PathBuilder

// Close 闭合路径 / Closes path
func (pb *PathBuilder) Close() *PathBuilder

// String 获取路径字符串 / Gets path string
func (pb *PathBuilder) String() string
```

**示例 / Example:**
```go
builder := NewPathBuilder()
builder.MoveTo(100, 100)
       .LineTo(200, 100)
       .LineTo(150, 200)
       .Close()

path := svg.Path(builder.String())
path.Fill("blue")
```

## 变换 / Transforms

### 变换方法 / Transform Methods

```go
// Transform 设置变换 / Sets transform
func (e *Element) Transform(transform string) *Element

// Translate 平移 / Translates
func (e *Element) Translate(x, y float64) *Element

// Rotate 旋转 / Rotates
func (e *Element) Rotate(angle float64, cx, cy float64) *Element

// Scale 缩放 / Scales
func (e *Element) Scale(sx, sy float64) *Element

// Skew 倾斜 / Skews
func (e *Element) SkewX(angle float64) *Element
func (e *Element) SkewY(angle float64) *Element
```

**示例 / Example:**
```go
// 使用变换字符串
rect := svg.Rect(50, 50, 100, 100)
rect.Transform("rotate(45 100 100) scale(1.5)")

// 使用便捷方法
circle := svg.Circle(100, 100, 50)
circle.Rotate(30, 100, 100).Scale(1.2, 0.8)
```

### 变换矩阵 / Transform Matrix

```go
// Matrix 设置变换矩阵 / Sets transform matrix
func (e *Element) Matrix(a, b, c, d, e, f float64) *Element
```

**参数 / Parameters:**
- `a, b, c, d, e, f`: 变换矩阵参数 / Transform matrix parameters

**示例 / Example:**
```go
// 自定义变换矩阵
rect := svg.Rect(0, 0, 100, 100)
rect.Matrix(1, 0, 0, 1, 50, 50) // 平移(50, 50)
```

## 错误处理 / Error Handling

库中的大多数方法都返回错误值，应该适当处理：

Most methods in the library return error values that should be handled appropriately:

```go
// 正确的错误处理 / Proper error handling
svg := svg.New(800, 600)
circle := svg.Circle(400, 300, 100)
circle.Fill("red")

if err := svg.Save("output.svg"); err != nil {
    log.Printf("保存SVG失败: %v", err)
    return
}

if err := svg.RenderToPNG("output.png", 800, 600); err != nil {
    log.Printf("渲染PNG失败: %v", err)
    return
}
```

## 性能优化 / Performance Optimization

### 批量操作 / Batch Operations

```go
// 批量添加元素 / Batch add elements
func (s *SVG) AddElements(elements ...Element) *SVG

// 批量设置样式 / Batch set styles
func (s *SVG) SetGlobalStyle(css string) *SVG
```

### 内存管理 / Memory Management

```go
// 清理资源 / Clean up resources
func (s *SVG) Cleanup()

// 重置SVG / Reset SVG
func (s *SVG) Reset() *SVG
```

## 最佳实践 / Best Practices

1. **资源管理 / Resource Management**
   ```go
   defer svg.Cleanup() // 确保资源被清理
   ```

2. **错误处理 / Error Handling**
   ```go
   if err := svg.Save(filename); err != nil {
       return fmt.Errorf("保存失败: %w", err)
   }
   ```

3. **性能优化 / Performance Optimization**
   ```go
   // 对于大量元素，使用批量操作
   elements := make([]Element, 0, 1000)
   // ... 创建元素
   svg.AddElements(elements...)
   ```

4. **内存使用 / Memory Usage**
   ```go
   // 对于大型SVG，考虑分块处理
   for i := 0; i < totalElements; i += batchSize {
       // 处理一批元素
       svg.ProcessBatch(i, min(i+batchSize, totalElements))
   }
   ```

---

## 更多信息 / More Information

- [快速开始指南](../README.md#快速开始) / [Quick Start Guide](../README.md#quick-start)
- [示例代码](../examples/) / [Example Code](../examples/)
- [贡献指南](../CONTRIBUTING.md) / [Contributing Guide](../CONTRIBUTING.md)
- [更新日志](../CHANGELOG.md) / [Changelog](../CHANGELOG.md)

如有问题或建议，请在GitHub上创建issue。

For questions or suggestions, please create an issue on GitHub.