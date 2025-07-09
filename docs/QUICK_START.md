# SVG库快速入门指南 / Quick Start Guide

## 🚀 5分钟快速开始 / 5-Minute Quick Start

欢迎使用Go SVG库！这个指南将帮助您在5分钟内创建第一个SVG图形。

Welcome to the Go SVG library! This guide will help you create your first SVG graphic in 5 minutes.

## 📦 安装和导入 / Installation and Import

### 前置要求 / Prerequisites

- Go 1.18 或更高版本 / Go 1.18 or higher
- 基本的Go语言知识 / Basic Go programming knowledge

### 导入库 / Import Library

```go
package main

import (
    "image/color"
    "svg"
)
```

## 🎯 第一个SVG程序 / Your First SVG Program

让我们创建一个简单的"Hello SVG"程序：

Let's create a simple "Hello SVG" program:

```go
package main

import (
    "fmt"
    "image/color"
    "svg"
)

func main() {
    // 创建SVG画布 / Create SVG canvas
    canvas := svg.New(400, 300)
    
    // 设置背景颜色 / Set background color
    canvas.SetBackground(color.RGBA{240, 248, 255, 255}) // 淡蓝色 / Light blue
    
    // 添加一个红色圆形 / Add a red circle
    canvas.Circle(200, 150, 50).Fill(color.RGBA{255, 0, 0, 255})
    
    // 添加文本 / Add text
    canvas.Text(200, 200, "Hello SVG!").FontSize(24).TextAnchor("middle")
    
    // 保存为SVG文件 / Save as SVG file
    err := canvas.SaveSVG("hello.svg")
    if err != nil {
        fmt.Printf("保存失败: %v\n", err)
        return
    }
    
    // 保存为PNG图片 / Save as PNG image
    err = canvas.SavePNG("hello.png")
    if err != nil {
        fmt.Printf("保存PNG失败: %v\n", err)
        return
    }
    
    fmt.Println("✅ SVG文件已创建: hello.svg")
    fmt.Println("✅ PNG文件已创建: hello.png")
}
```

运行这个程序：

Run this program:

```bash
go run main.go
```

您将得到两个文件：
- `hello.svg` - 矢量图形文件
- `hello.png` - 位图文件

You will get two files:
- `hello.svg` - Vector graphics file
- `hello.png` - Bitmap file

## 🖼️ 基本图形绘制 / Basic Shape Drawing

### 矩形 / Rectangle

```go
// 创建矩形 / Create rectangle
canvas.Rect(50, 50, 100, 80).Fill(color.RGBA{0, 255, 0, 255})

// 带圆角的矩形 / Rectangle with rounded corners
canvas.Rect(200, 50, 100, 80).Fill(color.RGBA{0, 0, 255, 255}).Rx(10)
```

### 圆形和椭圆 / Circle and Ellipse

```go
// 圆形 / Circle
canvas.Circle(100, 200, 40).Fill(color.RGBA{255, 255, 0, 255})

// 椭圆 / Ellipse
canvas.Ellipse(250, 200, 60, 30).Fill(color.RGBA{255, 0, 255, 255})
```

### 直线 / Line

```go
// 直线 / Line
canvas.Line(50, 300, 350, 300).Stroke(color.RGBA{0, 0, 0, 255}).StrokeWidth(2)
```

### 完整示例 / Complete Example

```go
package main

import (
    "image/color"
    "svg"
)

func main() {
    // 创建画布 / Create canvas
    canvas := svg.New(400, 400)
    canvas.SetBackground(color.RGBA{250, 250, 250, 255})
    
    // 绘制彩色矩形 / Draw colorful rectangles
    canvas.Rect(50, 50, 80, 60).Fill(color.RGBA{255, 100, 100, 255})
    canvas.Rect(150, 50, 80, 60).Fill(color.RGBA{100, 255, 100, 255})
    canvas.Rect(250, 50, 80, 60).Fill(color.RGBA{100, 100, 255, 255})
    
    // 绘制圆形 / Draw circles
    canvas.Circle(90, 180, 30).Fill(color.RGBA{255, 200, 0, 255})
    canvas.Circle(190, 180, 30).Fill(color.RGBA{255, 0, 200, 255})
    canvas.Circle(290, 180, 30).Fill(color.RGBA{0, 255, 200, 255})
    
    // 绘制线条 / Draw lines
    canvas.Line(50, 250, 350, 250).Stroke(color.RGBA{100, 100, 100, 255}).StrokeWidth(3)
    canvas.Line(200, 280, 200, 350).Stroke(color.RGBA{100, 100, 100, 255}).StrokeWidth(3)
    
    // 添加标题 / Add title
    canvas.Text(200, 30, "基本图形示例").FontSize(20).TextAnchor("middle").Fill(color.RGBA{50, 50, 50, 255})
    
    // 保存文件 / Save files
    canvas.SaveSVG("basic_shapes.svg")
    canvas.SavePNG("basic_shapes.png")
}
```

## 🎨 样式设置 / Style Settings

### 颜色设置 / Color Settings

```go
// RGB颜色 / RGB colors
red := color.RGBA{255, 0, 0, 255}
green := color.RGBA{0, 255, 0, 255}
blue := color.RGBA{0, 0, 255, 255}

// 半透明颜色 / Semi-transparent colors
transparentRed := color.RGBA{255, 0, 0, 128}

// 使用颜色 / Use colors
canvas.Circle(100, 100, 50).Fill(red)
canvas.Rect(200, 50, 100, 100).Fill(transparentRed)
```

### 描边设置 / Stroke Settings

```go
// 设置描边颜色和宽度 / Set stroke color and width
canvas.Circle(150, 150, 40).
    Fill(color.RGBA{255, 255, 255, 255}).  // 白色填充 / White fill
    Stroke(color.RGBA{0, 0, 0, 255}).      // 黑色描边 / Black stroke
    StrokeWidth(3)                         // 描边宽度3 / Stroke width 3
```

### 文本样式 / Text Styles

```go
// 基本文本 / Basic text
canvas.Text(200, 100, "普通文本").FontSize(16)

// 粗体文本 / Bold text
canvas.Text(200, 130, "粗体文本").FontSize(16).FontWeight("bold")

// 居中文本 / Centered text
canvas.Text(200, 160, "居中文本").FontSize(16).TextAnchor("middle")

// 彩色文本 / Colored text
canvas.Text(200, 190, "彩色文本").FontSize(16).Fill(color.RGBA{255, 0, 0, 255})
```

## 💾 保存和导出 / Save and Export

### 保存为SVG / Save as SVG

```go
// 保存SVG文件 / Save SVG file
err := canvas.SaveSVG("my_drawing.svg")
if err != nil {
    fmt.Printf("保存SVG失败: %v\n", err)
}
```

### 保存为PNG / Save as PNG

```go
// 保存PNG文件 / Save PNG file
err := canvas.SavePNG("my_drawing.png")
if err != nil {
    fmt.Printf("保存PNG失败: %v\n", err)
}
```

### 保存为JPEG / Save as JPEG

```go
// 保存JPEG文件 / Save JPEG file
err := canvas.SaveJPEG("my_drawing.jpg")
if err != nil {
    fmt.Printf("保存JPEG失败: %v\n", err)
}
```

### 自定义尺寸导出 / Custom Size Export

```go
// 以指定尺寸保存PNG / Save PNG with specified size
err := canvas.SavePNGWithSize("large_image.png", 800, 600)
if err != nil {
    fmt.Printf("保存大尺寸PNG失败: %v\n", err)
}
```

## 🎬 简单动画 / Simple Animation

创建一个简单的旋转动画：

Create a simple rotation animation:

```go
package main

import (
    "image/color"
    "svg"
)

func main() {
    // 使用动画构建器 / Use animation builder
    builder := svg.NewAnimationBuilder(400, 400)
    builder.SetFrameCount(60).SetFrameRate(30)
    
    // 配置动画 / Configure animation
    config := svg.AnimationConfig{
        Duration:   2.0, // 2秒 / 2 seconds
        Easing:     svg.EaseInOut,
        Background: color.RGBA{20, 20, 40, 255},
        Loop:       true,
    }
    
    // 创建旋转图形动画 / Create rotating shapes animation
    err := builder.CreateRotatingShapes(config).SaveToGIF("rotation.gif")
    if err != nil {
        fmt.Printf("创建动画失败: %v\n", err)
        return
    }
    
    fmt.Println("✅ 动画已创建: rotation.gif")
}
```

## 🔧 实用技巧 / Useful Tips

### 1. 方法链 / Method Chaining

```go
// 可以链式调用方法 / You can chain methods
canvas.Circle(200, 200, 50).
    Fill(color.RGBA{255, 0, 0, 255}).
    Stroke(color.RGBA{0, 0, 0, 255}).
    StrokeWidth(2)
```

### 2. 颜色常量 / Color Constants

```go
// 使用预定义颜色 / Use predefined colors
var (
    Red    = color.RGBA{255, 0, 0, 255}
    Green  = color.RGBA{0, 255, 0, 255}
    Blue   = color.RGBA{0, 0, 255, 255}
    White  = color.RGBA{255, 255, 255, 255}
    Black  = color.RGBA{0, 0, 0, 255}
)

canvas.Circle(100, 100, 50).Fill(Red)
```

### 3. 错误处理 / Error Handling

```go
// 总是检查错误 / Always check errors
if err := canvas.SaveSVG("output.svg"); err != nil {
    fmt.Printf("保存失败: %v\n", err)
    return
}
```

### 4. 文件路径 / File Paths

```go
// 创建输出目录 / Create output directory
os.MkdirAll("output", 0755)

// 保存到指定目录 / Save to specific directory
canvas.SaveSVG("output/my_drawing.svg")
canvas.SavePNG("output/my_drawing.png")
```

## 📚 下一步学习 / Next Steps

恭喜！您已经掌握了SVG库的基础用法。接下来可以学习：

Congratulations! You've mastered the basics of the SVG library. Next, you can learn:

1. **基础教程** - 深入学习所有图形元素和样式
   **Basic Tutorial** - Learn all graphic elements and styles in depth

2. **进阶教程** - 学习动画、高级API和自定义功能
   **Advanced Tutorial** - Learn animations, advanced APIs and custom features

3. **示例集合** - 查看更多实用示例
   **Examples Collection** - View more practical examples

4. **API参考** - 查阅完整的API文档
   **API Reference** - Consult complete API documentation

## ❓ 常见问题 / FAQ

### Q: 如何设置画布大小？ / How to set canvas size?

```go
// 创建指定大小的画布 / Create canvas with specified size
canvas := svg.New(800, 600) // 宽800，高600 / Width 800, Height 600
```

### Q: 如何创建透明背景？ / How to create transparent background?

```go
// 不设置背景颜色即为透明 / No background color means transparent
canvas := svg.New(400, 300)
// 或者显式设置透明背景 / Or explicitly set transparent background
canvas.SetBackground(color.RGBA{0, 0, 0, 0})
```

### Q: 支持哪些图片格式？ / What image formats are supported?

支持的格式：
- SVG (矢量格式)
- PNG (支持透明)
- JPEG (不支持透明)
- GIF (动画格式)

Supported formats:
- SVG (vector format)
- PNG (supports transparency)
- JPEG (no transparency)
- GIF (animation format)

### Q: 如何调试渲染问题？ / How to debug rendering issues?

```go
// 先保存为SVG查看结构 / Save as SVG first to check structure
canvas.SaveSVG("debug.svg")

// 然后保存为PNG查看渲染结果 / Then save as PNG to see rendering result
canvas.SavePNG("debug.png")
```

---

🎉 **恭喜您完成了快速入门！** / **Congratulations on completing the quick start!**

现在您已经可以创建基本的SVG图形了。继续探索更多功能，创造出精彩的图形作品！

Now you can create basic SVG graphics. Continue exploring more features to create amazing graphic works!