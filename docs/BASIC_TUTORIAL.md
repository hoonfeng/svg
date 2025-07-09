# SVG库基础教程 / Basic Tutorial

## 📖 教程概述 / Tutorial Overview

本教程将深入介绍SVG库的核心功能，包括所有基础图形、样式系统、文本处理和路径绘制。完成本教程后，您将能够创建复杂的SVG图形。

This tutorial provides an in-depth introduction to the core features of the SVG library, including all basic shapes, styling system, text processing, and path drawing. After completing this tutorial, you'll be able to create complex SVG graphics.

## 📐 基础图形绘制 / Basic Shape Drawing

### 矩形 (Rectangle)

矩形是最基础的SVG图形之一，支持圆角、填充、描边等属性。

Rectangles are one of the most basic SVG shapes, supporting rounded corners, fill, stroke, and other attributes.

```go
package main

import (
    "image/color"
    "svg"
)

func main() {
    canvas := svg.New(500, 400)
    canvas.SetBackground(color.RGBA{245, 245, 245, 255})
    
    // 基本矩形 / Basic rectangle
    canvas.Rect(50, 50, 100, 80).Fill(color.RGBA{255, 100, 100, 255})
    
    // 带描边的矩形 / Rectangle with stroke
    canvas.Rect(200, 50, 100, 80).
        Fill(color.RGBA{100, 255, 100, 255}).
        Stroke(color.RGBA{0, 0, 0, 255}).
        StrokeWidth(2)
    
    // 圆角矩形 / Rounded rectangle
    canvas.Rect(350, 50, 100, 80).
        Fill(color.RGBA{100, 100, 255, 255}).
        Rx(15).Ry(15)
    
    // 半透明矩形 / Semi-transparent rectangle
    canvas.Rect(50, 180, 100, 80).
        Fill(color.RGBA{255, 255, 0, 128})
    
    // 只有描边的矩形 / Stroke-only rectangle
    canvas.Rect(200, 180, 100, 80).
        Fill(color.RGBA{0, 0, 0, 0}). // 透明填充 / Transparent fill
        Stroke(color.RGBA{255, 0, 255, 255}).
        StrokeWidth(3)
    
    // 虚线描边矩形 / Dashed stroke rectangle
    canvas.Rect(350, 180, 100, 80).
        Fill(color.RGBA{255, 255, 255, 255}).
        Stroke(color.RGBA{0, 0, 0, 255}).
        StrokeWidth(2).
        StrokeDashArray("5,5")
    
    canvas.SaveSVG("rectangles_demo.svg")
    canvas.SavePNG("rectangles_demo.png")
}
```

### 圆形和椭圆 (Circle and Ellipse)

圆形和椭圆用于创建圆润的图形元素。

Circles and ellipses are used to create rounded graphic elements.

```go
func drawCirclesAndEllipses() {
    canvas := svg.New(600, 400)
    canvas.SetBackground(color.RGBA{250, 250, 250, 255})
    
    // 基本圆形 / Basic circle
    canvas.Circle(100, 100, 40).Fill(color.RGBA{255, 0, 0, 255})
    
    // 带描边的圆形 / Circle with stroke
    canvas.Circle(250, 100, 40).
        Fill(color.RGBA{255, 255, 255, 255}).
        Stroke(color.RGBA{0, 0, 255, 255}).
        StrokeWidth(4)
    
    // 渐变填充圆形 (需要定义渐变) / Gradient filled circle
    canvas.Circle(400, 100, 40).Fill(color.RGBA{255, 200, 0, 255})
    
    // 基本椭圆 / Basic ellipse
    canvas.Ellipse(100, 250, 60, 30).Fill(color.RGBA{0, 255, 0, 255})
    
    // 旋转的椭圆 / Rotated ellipse
    canvas.Ellipse(250, 250, 60, 30).
        Fill(color.RGBA{255, 0, 255, 255}).
        Transform("rotate(45 250 250)")
    
    // 多重描边椭圆 / Multiple stroke ellipse
    canvas.Ellipse(400, 250, 60, 30).
        Fill(color.RGBA{0, 255, 255, 128}).
        Stroke(color.RGBA{0, 0, 0, 255}).
        StrokeWidth(2)
    
    canvas.SaveSVG("circles_ellipses_demo.svg")
    canvas.SavePNG("circles_ellipses_demo.png")
}
```

### 直线和折线 (Line and Polyline)

直线和折线用于创建线性图形和连接元素。

Lines and polylines are used to create linear graphics and connect elements.

```go
func drawLinesAndPolylines() {
    canvas := svg.New(600, 400)
    canvas.SetBackground(color.RGBA{255, 255, 255, 255})
    
    // 基本直线 / Basic line
    canvas.Line(50, 50, 200, 50).
        Stroke(color.RGBA{0, 0, 0, 255}).
        StrokeWidth(2)
    
    // 彩色粗线 / Colored thick line
    canvas.Line(50, 100, 200, 100).
        Stroke(color.RGBA{255, 0, 0, 255}).
        StrokeWidth(5)
    
    // 虚线 / Dashed line
    canvas.Line(50, 150, 200, 150).
        Stroke(color.RGBA{0, 0, 255, 255}).
        StrokeWidth(3).
        StrokeDashArray("10,5")
    
    // 点线 / Dotted line
    canvas.Line(50, 200, 200, 200).
        Stroke(color.RGBA{0, 255, 0, 255}).
        StrokeWidth(3).
        StrokeDashArray("2,3")
    
    // 折线 / Polyline
    points := "300,50 350,100 400,80 450,120 500,90"
    canvas.Polyline(points).
        Fill("none").
        Stroke(color.RGBA{255, 100, 0, 255}).
        StrokeWidth(3)
    
    // 闭合折线 (多边形) / Closed polyline (polygon)
    polygonPoints := "300,200 350,180 400,200 380,250 320,250"
    canvas.Polygon(polygonPoints).
        Fill(color.RGBA{100, 200, 255, 128}).
        Stroke(color.RGBA{0, 0, 100, 255}).
        StrokeWidth(2)
    
    canvas.SaveSVG("lines_polylines_demo.svg")
    canvas.SavePNG("lines_polylines_demo.png")
}
```

### 多边形 (Polygon)

多边形用于创建复杂的几何图形。

Polygons are used to create complex geometric shapes.

```go
func drawPolygons() {
    canvas := svg.New(600, 400)
    canvas.SetBackground(color.RGBA{240, 240, 240, 255})
    
    // 三角形 / Triangle
    triangle := "100,50 50,150 150,150"
    canvas.Polygon(triangle).
        Fill(color.RGBA{255, 0, 0, 255}).
        Stroke(color.RGBA{0, 0, 0, 255}).
        StrokeWidth(2)
    
    // 五角星 / Pentagon star
    star := "300,50 310,80 340,80 320,100 330,130 300,115 270,130 280,100 260,80 290,80"
    canvas.Polygon(star).
        Fill(color.RGBA{255, 215, 0, 255}).
        Stroke(color.RGBA{255, 140, 0, 255}).
        StrokeWidth(2)
    
    // 六边形 / Hexagon
    hexagon := "500,80 530,100 530,140 500,160 470,140 470,100"
    canvas.Polygon(hexagon).
        Fill(color.RGBA{0, 255, 0, 255}).
        Stroke(color.RGBA{0, 150, 0, 255}).
        StrokeWidth(3)
    
    // 复杂多边形 / Complex polygon
    complex := "100,250 150,200 200,220 180,270 220,300 170,320 120,300 80,280"
    canvas.Polygon(complex).
        Fill(color.RGBA{128, 0, 128, 200}).
        Stroke(color.RGBA{64, 0, 64, 255}).
        StrokeWidth(2)
    
    canvas.SaveSVG("polygons_demo.svg")
    canvas.SavePNG("polygons_demo.png")
}
```

## 🎨 样式和颜色 / Styles and Colors

### 颜色系统 / Color System

SVG库支持多种颜色格式和表示方法。

The SVG library supports various color formats and representation methods.

```go
func demonstrateColors() {
    canvas := svg.New(700, 500)
    
    // RGB颜色 / RGB colors
    red := color.RGBA{255, 0, 0, 255}
    green := color.RGBA{0, 255, 0, 255}
    blue := color.RGBA{0, 0, 255, 255}
    
    // RGBA颜色 (带透明度) / RGBA colors (with transparency)
    transparentRed := color.RGBA{255, 0, 0, 128}
    transparentGreen := color.RGBA{0, 255, 0, 128}
    transparentBlue := color.RGBA{0, 0, 255, 128}
    
    // 基本颜色示例 / Basic color examples
    canvas.Rect(50, 50, 80, 60).Fill(red)
    canvas.Rect(150, 50, 80, 60).Fill(green)
    canvas.Rect(250, 50, 80, 60).Fill(blue)
    
    // 透明度示例 / Transparency examples
    canvas.Rect(50, 150, 80, 60).Fill(transparentRed)
    canvas.Rect(150, 150, 80, 60).Fill(transparentGreen)
    canvas.Rect(250, 150, 80, 60).Fill(transparentBlue)
    
    // 重叠透明图形 / Overlapping transparent shapes
    canvas.Circle(400, 100, 40).Fill(color.RGBA{255, 0, 0, 100})
    canvas.Circle(430, 100, 40).Fill(color.RGBA{0, 255, 0, 100})
    canvas.Circle(415, 130, 40).Fill(color.RGBA{0, 0, 255, 100})
    
    // 灰度颜色 / Grayscale colors
    for i := 0; i < 10; i++ {
        gray := uint8(i * 25)
        canvas.Rect(float64(50+i*50), 250, 40, 40).
            Fill(color.RGBA{gray, gray, gray, 255})
    }
    
    // 颜色渐变效果 (通过多个图形模拟) / Color gradient effect (simulated with multiple shapes)
    for i := 0; i < 20; i++ {
        r := uint8(255 - i*12)
        g := uint8(i * 12)
        canvas.Rect(float64(50+i*25), 350, 25, 40).
            Fill(color.RGBA{r, g, 0, 255})
    }
    
    canvas.SaveSVG("colors_demo.svg")
    canvas.SavePNG("colors_demo.png")
}
```

### 描边样式 / Stroke Styles

描边是图形的轮廓线，可以设置颜色、宽度、样式等属性。

Strokes are the outline of shapes, with customizable color, width, style, and other attributes.

```go
func demonstrateStrokes() {
    canvas := svg.New(600, 500)
    canvas.SetBackground(color.RGBA{250, 250, 250, 255})
    
    // 不同宽度的描边 / Different stroke widths
    for i := 1; i <= 5; i++ {
        canvas.Circle(float64(80+i*100), 80, 30).
            Fill(color.RGBA{255, 255, 255, 255}).
            Stroke(color.RGBA{0, 0, 0, 255}).
            StrokeWidth(float64(i))
    }
    
    // 不同颜色的描边 / Different stroke colors
    colors := []color.RGBA{
        {255, 0, 0, 255},   // 红色 / Red
        {0, 255, 0, 255},   // 绿色 / Green
        {0, 0, 255, 255},   // 蓝色 / Blue
        {255, 255, 0, 255}, // 黄色 / Yellow
        {255, 0, 255, 255}, // 洋红 / Magenta
    }
    
    for i, c := range colors {
        canvas.Rect(float64(80+i*100), 150, 60, 40).
            Fill(color.RGBA{255, 255, 255, 255}).
            Stroke(c).
            StrokeWidth(3)
    }
    
    // 虚线样式 / Dash patterns
    dashPatterns := []string{
        "5,5",     // 短虚线 / Short dashes
        "10,5",    // 长虚线 / Long dashes
        "15,5,5,5", // 长短交替 / Long-short alternating
        "2,3",     // 点线 / Dotted
        "20,5,5,5,5,5", // 复杂模式 / Complex pattern
    }
    
    for i, pattern := range dashPatterns {
        canvas.Line(80, float64(250+i*30), 520, float64(250+i*30)).
            Stroke(color.RGBA{0, 0, 0, 255}).
            StrokeWidth(3).
            StrokeDashArray(pattern)
    }
    
    // 线端样式 / Line cap styles
    canvas.Line(80, 420, 180, 420).
        Stroke(color.RGBA{255, 0, 0, 255}).
        StrokeWidth(10).
        StrokeLineCap("butt") // 平端 / Flat end
    
    canvas.Line(220, 420, 320, 420).
        Stroke(color.RGBA{0, 255, 0, 255}).
        StrokeWidth(10).
        StrokeLineCap("round") // 圆端 / Round end
    
    canvas.Line(360, 420, 460, 420).
        Stroke(color.RGBA{0, 0, 255, 255}).
        StrokeWidth(10).
        StrokeLineCap("square") // 方端 / Square end
    
    canvas.SaveSVG("strokes_demo.svg")
    canvas.SavePNG("strokes_demo.png")
}
```

### 填充样式 / Fill Styles

填充是图形内部的颜色或图案。

Fill is the color or pattern inside shapes.

```go
func demonstrateFills() {
    canvas := svg.New(600, 400)
    canvas.SetBackground(color.RGBA{240, 240, 240, 255})
    
    // 纯色填充 / Solid color fills
    canvas.Rect(50, 50, 80, 60).Fill(color.RGBA{255, 100, 100, 255})
    canvas.Rect(150, 50, 80, 60).Fill(color.RGBA{100, 255, 100, 255})
    canvas.Rect(250, 50, 80, 60).Fill(color.RGBA{100, 100, 255, 255})
    
    // 透明填充 / Transparent fills
    canvas.Circle(100, 180, 40).Fill(color.RGBA{255, 0, 0, 100})
    canvas.Circle(130, 180, 40).Fill(color.RGBA{0, 255, 0, 100})
    canvas.Circle(115, 210, 40).Fill(color.RGBA{0, 0, 255, 100})
    
    // 无填充 (只有描边) / No fill (stroke only)
    canvas.Rect(350, 50, 80, 60).
        Fill(color.RGBA{0, 0, 0, 0}). // 透明填充 / Transparent fill
        Stroke(color.RGBA{0, 0, 0, 255}).
        StrokeWidth(2)
    
    // 填充和描边组合 / Fill and stroke combination
    canvas.Circle(400, 180, 40).
        Fill(color.RGBA{255, 255, 0, 255}).
        Stroke(color.RGBA{255, 0, 0, 255}).
        StrokeWidth(4)
    
    // 渐变色模拟 (使用多个形状) / Gradient simulation (using multiple shapes)
    for i := 0; i < 16; i++ {
        intensity := uint8(255 - i*15)
        canvas.Rect(float64(50+i*30), 300, 30, 40).
            Fill(color.RGBA{intensity, 100, 255-intensity, 255})
    }
    
    canvas.SaveSVG("fills_demo.svg")
    canvas.SavePNG("fills_demo.png")
}
```

## 📝 文本处理 / Text Processing

### 基本文本绘制 / Basic Text Drawing

文本是SVG中重要的元素，支持多种字体和样式设置。

Text is an important element in SVG, supporting various fonts and style settings.

```go
func demonstrateBasicText() {
    canvas := svg.New(800, 600)
    canvas.SetBackground(color.RGBA{255, 255, 255, 255})
    
    // 基本文本 / Basic text
    canvas.Text(50, 50, "Hello, SVG World!").
        FontSize(24).
        Fill(color.RGBA{0, 0, 0, 255})
    
    // 不同大小的文本 / Different text sizes
    sizes := []float64{12, 16, 20, 24, 32, 40}
    for i, size := range sizes {
        canvas.Text(50, float64(100+i*50), fmt.Sprintf("字体大小 %.0f", size)).
            FontSize(size).
            Fill(color.RGBA{0, 0, 0, 255})
    }
    
    // 不同颜色的文本 / Different colored text
    colors := []color.RGBA{
        {255, 0, 0, 255},   // 红色 / Red
        {0, 255, 0, 255},   // 绿色 / Green
        {0, 0, 255, 255},   // 蓝色 / Blue
        {255, 165, 0, 255}, // 橙色 / Orange
        {128, 0, 128, 255}, // 紫色 / Purple
    }
    
    for i, c := range colors {
        canvas.Text(400, float64(100+i*40), "彩色文本").
            FontSize(20).
            Fill(c)
    }
    
    canvas.SaveSVG("basic_text_demo.svg")
    canvas.SavePNG("basic_text_demo.png")
}
```

### 字体样式 / Font Styles

```go
func demonstrateFontStyles() {
    canvas := svg.New(800, 700)
    canvas.SetBackground(color.RGBA{250, 250, 250, 255})
    
    // 字体族 / Font families
    fontFamilies := []string{
        "Arial",
        "Times New Roman",
        "Courier New",
        "Helvetica",
        "Georgia",
    }
    
    for i, family := range fontFamilies {
        canvas.Text(50, float64(50+i*40), fmt.Sprintf("字体: %s", family)).
            FontSize(18).
            FontFamily(family).
            Fill(color.RGBA{0, 0, 0, 255})
    }
    
    // 字体粗细 / Font weights
    weights := []string{"normal", "bold", "lighter", "bolder"}
    for i, weight := range weights {
        canvas.Text(400, float64(50+i*40), fmt.Sprintf("粗细: %s", weight)).
            FontSize(18).
            FontWeight(weight).
            Fill(color.RGBA{0, 0, 0, 255})
    }
    
    // 字体样式 / Font styles
    styles := []string{"normal", "italic", "oblique"}
    for i, style := range styles {
        canvas.Text(50, float64(300+i*40), fmt.Sprintf("样式: %s", style)).
            FontSize(18).
            FontStyle(style).
            Fill(color.RGBA{0, 0, 0, 255})
    }
    
    // 文本装饰 / Text decorations
    decorations := []string{"none", "underline", "overline", "line-through"}
    for i, decoration := range decorations {
        canvas.Text(400, float64(300+i*40), fmt.Sprintf("装饰: %s", decoration)).
            FontSize(18).
            TextDecoration(decoration).
            Fill(color.RGBA{0, 0, 0, 255})
    }
    
    canvas.SaveSVG("font_styles_demo.svg")
    canvas.SavePNG("font_styles_demo.png")
}
```

### 文本对齐 / Text Alignment

```go
func demonstrateTextAlignment() {
    canvas := svg.New(600, 400)
    canvas.SetBackground(color.RGBA{255, 255, 255, 255})
    
    // 绘制参考线 / Draw reference lines
    centerX := 300.0
    canvas.Line(centerX, 50, centerX, 350).
        Stroke(color.RGBA{200, 200, 200, 255}).
        StrokeWidth(1).
        StrokeDashArray("2,2")
    
    // 文本锚点 / Text anchors
    anchors := []string{"start", "middle", "end"}
    for i, anchor := range anchors {
        y := float64(100 + i*60)
        
        // 绘制锚点 / Draw anchor point
        canvas.Circle(centerX, y, 3).Fill(color.RGBA{255, 0, 0, 255})
        
        // 绘制文本 / Draw text
        canvas.Text(centerX, y, fmt.Sprintf("文本锚点: %s", anchor)).
            FontSize(16).
            TextAnchor(anchor).
            Fill(color.RGBA{0, 0, 0, 255})
    }
    
    // 基线对齐 / Baseline alignment
    baselines := []string{"alphabetic", "middle", "hanging"}
    for i, baseline := range baselines {
        y := float64(280 + i*30)
        
        // 绘制基线 / Draw baseline
        canvas.Line(50, y, 550, y).
            Stroke(color.RGBA{200, 200, 200, 255}).
            StrokeWidth(1).
            StrokeDashArray("2,2")
        
        // 绘制文本 / Draw text
        canvas.Text(100, y, fmt.Sprintf("基线: %s", baseline)).
            FontSize(16).
            AlignmentBaseline(baseline).
            Fill(color.RGBA{0, 0, 0, 255})
    }
    
    canvas.SaveSVG("text_alignment_demo.svg")
    canvas.SavePNG("text_alignment_demo.png")
}
```

## 🛤️ 路径基础 / Path Basics

### 路径语法入门 / Path Syntax Introduction

路径是SVG中最强大的绘图工具，可以创建任意复杂的图形。

Paths are the most powerful drawing tool in SVG, capable of creating arbitrarily complex graphics.

```go
func demonstrateBasicPaths() {
    canvas := svg.New(800, 600)
    canvas.SetBackground(color.RGBA{250, 250, 250, 255})
    
    // 简单直线路径 / Simple line path
    canvas.Path("M 50 50 L 200 50").
        Fill("none").
        Stroke(color.RGBA{0, 0, 0, 255}).
        StrokeWidth(2)
    
    // 折线路径 / Polyline path
    canvas.Path("M 50 100 L 100 80 L 150 120 L 200 100").
        Fill("none").
        Stroke(color.RGBA{255, 0, 0, 255}).
        StrokeWidth(2)
    
    // 闭合路径 / Closed path
    canvas.Path("M 50 150 L 100 130 L 150 170 L 200 150 Z").
        Fill(color.RGBA{255, 255, 0, 128}).
        Stroke(color.RGBA{255, 0, 0, 255}).
        StrokeWidth(2)
    
    // 曲线路径 / Curved path
    canvas.Path("M 50 200 Q 125 150 200 200").
        Fill("none").
        Stroke(color.RGBA{0, 0, 255, 255}).
        StrokeWidth(3)
    
    // 贝塞尔曲线 / Bezier curve
    canvas.Path("M 50 250 C 75 200, 125 300, 200 250").
        Fill("none").
        Stroke(color.RGBA{0, 255, 0, 255}).
        StrokeWidth(3)
    
    // 圆弧 / Arc
    canvas.Path("M 50 300 A 50 30 0 0 1 200 300").
        Fill("none").
        Stroke(color.RGBA{255, 0, 255, 255}).
        StrokeWidth(3)
    
    canvas.SaveSVG("basic_paths_demo.svg")
    canvas.SavePNG("basic_paths_demo.png")
}
```

### 复杂路径示例 / Complex Path Examples

```go
func demonstrateComplexPaths() {
    canvas := svg.New(800, 600)
    canvas.SetBackground(color.RGBA{255, 255, 255, 255})
    
    // 心形 / Heart shape
    heartPath := "M 400,200 C 400,180 380,160 350,160 C 320,160 300,180 300,200 C 300,180 280,160 250,160 C 220,160 200,180 200,200 C 200,240 400,320 400,320 C 400,320 600,240 600,200 C 600,180 580,160 550,160 C 520,160 500,180 500,200 C 500,180 480,160 450,160 C 420,160 400,180 400,200 Z"
    canvas.Path(heartPath).
        Fill(color.RGBA{255, 100, 100, 255}).
        Stroke(color.RGBA{200, 0, 0, 255}).
        StrokeWidth(2)
    
    // 星形 / Star shape
    starPath := "M 150,50 L 160,80 L 190,80 L 170,100 L 180,130 L 150,115 L 120,130 L 130,100 L 110,80 L 140,80 Z"
    canvas.Path(starPath).
        Fill(color.RGBA{255, 215, 0, 255}).
        Stroke(color.RGBA{255, 140, 0, 255}).
        StrokeWidth(2)
    
    // 花朵形状 / Flower shape
    flowerPath := "M 650,150 C 650,130 630,110 600,110 C 570,110 550,130 550,150 C 530,150 510,170 510,200 C 510,230 530,250 550,250 C 550,270 570,290 600,290 C 630,290 650,270 650,250 C 670,250 690,230 690,200 C 690,170 670,150 650,150 Z"
    canvas.Path(flowerPath).
        Fill(color.RGBA{255, 192, 203, 255}).
        Stroke(color.RGBA{255, 20, 147, 255}).
        StrokeWidth(2)
    
    // 波浪线 / Wave line
    wavePath := "M 50,400 Q 100,350 150,400 T 250,400 T 350,400 T 450,400 T 550,400 T 650,400 T 750,400"
    canvas.Path(wavePath).
        Fill("none").
        Stroke(color.RGBA{0, 100, 255, 255}).
        StrokeWidth(4)
    
    // 螺旋线 / Spiral
    spiralPath := "M 400,450 Q 410,440 420,450 Q 430,470 410,480 Q 380,490 370,470 Q 360,440 390,430 Q 430,420 450,450 Q 470,490 430,510 Q 370,530 340,490 Q 310,440 360,410 Q 420,380 470,420"
    canvas.Path(spiralPath).
        Fill("none").
        Stroke(color.RGBA{128, 0, 128, 255}).
        StrokeWidth(3)
    
    canvas.SaveSVG("complex_paths_demo.svg")
    canvas.SavePNG("complex_paths_demo.png")
}
```

## 🔧 实用示例 / Practical Examples

### 创建图标 / Creating Icons

```go
func createIcons() {
    canvas := svg.New(600, 400)
    canvas.SetBackground(color.RGBA{245, 245, 245, 255})
    
    // 主页图标 / Home icon
    canvas.Path("M 80,120 L 120,80 L 160,120 L 160,160 L 80,160 Z").
        Fill(color.RGBA{100, 150, 255, 255}).
        Stroke(color.RGBA{0, 100, 200, 255}).
        StrokeWidth(2)
    canvas.Rect(100, 140, 40, 20).Fill(color.RGBA{139, 69, 19, 255})
    
    // 设置图标 / Settings icon
    canvas.Circle(250, 120, 30).
        Fill("none").
        Stroke(color.RGBA{100, 100, 100, 255}).
        StrokeWidth(4)
    canvas.Circle(250, 120, 8).Fill(color.RGBA{100, 100, 100, 255})
    for i := 0; i < 8; i++ {
        angle := float64(i) * 45
        canvas.Line(250, 120, 250, 120).
            Transform(fmt.Sprintf("rotate(%f 250 120)", angle)).
            Stroke(color.RGBA{100, 100, 100, 255}).
            StrokeWidth(2)
    }
    
    // 邮件图标 / Mail icon
    canvas.Rect(320, 100, 80, 50).
        Fill(color.RGBA{255, 255, 255, 255}).
        Stroke(color.RGBA{200, 200, 200, 255}).
        StrokeWidth(2)
    canvas.Path("M 320,100 L 360,125 L 400,100").
        Fill("none").
        Stroke(color.RGBA{200, 200, 200, 255}).
        StrokeWidth(2)
    
    // 搜索图标 / Search icon
    canvas.Circle(480, 120, 20).
        Fill("none").
        Stroke(color.RGBA{0, 0, 0, 255}).
        StrokeWidth(3)
    canvas.Line(495, 135, 510, 150).
        Stroke(color.RGBA{0, 0, 0, 255}).
        StrokeWidth(3)
    
    canvas.SaveSVG("icons_demo.svg")
    canvas.SavePNG("icons_demo.png")
}
```

### 简单图表 / Simple Charts

```go
func createSimpleChart() {
    canvas := svg.New(600, 400)
    canvas.SetBackground(color.RGBA{255, 255, 255, 255})
    
    // 绘制坐标轴 / Draw axes
    canvas.Line(50, 350, 550, 350).
        Stroke(color.RGBA{0, 0, 0, 255}).
        StrokeWidth(2) // X轴 / X-axis
    canvas.Line(50, 50, 50, 350).
        Stroke(color.RGBA{0, 0, 0, 255}).
        StrokeWidth(2) // Y轴 / Y-axis
    
    // 数据点 / Data points
    data := []float64{80, 120, 200, 150, 250, 180, 300}
    colors := []color.RGBA{
        {255, 0, 0, 255},
        {0, 255, 0, 255},
        {0, 0, 255, 255},
        {255, 255, 0, 255},
        {255, 0, 255, 255},
        {0, 255, 255, 255},
        {255, 165, 0, 255},
    }
    
    // 绘制柱状图 / Draw bar chart
    for i, value := range data {
        x := float64(80 + i*60)
        height := value
        y := 350 - height
        
        canvas.Rect(x, y, 40, height).
            Fill(colors[i]).
            Stroke(color.RGBA{0, 0, 0, 255}).
            StrokeWidth(1)
        
        // 添加数值标签 / Add value labels
        canvas.Text(x+20, y-10, fmt.Sprintf("%.0f", value)).
            FontSize(12).
            TextAnchor("middle").
            Fill(color.RGBA{0, 0, 0, 255})
    }
    
    // 添加标题 / Add title
    canvas.Text(300, 30, "简单柱状图").
        FontSize(20).
        TextAnchor("middle").
        FontWeight("bold").
        Fill(color.RGBA{0, 0, 0, 255})
    
    canvas.SaveSVG("simple_chart_demo.svg")
    canvas.SavePNG("simple_chart_demo.png")
}
```

## 📚 总结 / Summary

通过本教程，您已经学会了：

Through this tutorial, you have learned:

### ✅ 掌握的技能 / Skills Mastered

1. **基础图形绘制** / **Basic Shape Drawing**
   - 矩形、圆形、椭圆、直线、多边形
   - Rectangles, circles, ellipses, lines, polygons

2. **样式系统** / **Styling System**
   - 颜色设置、填充、描边、透明度
   - Color settings, fills, strokes, transparency

3. **文本处理** / **Text Processing**
   - 字体设置、样式、对齐方式
   - Font settings, styles, alignment

4. **路径绘制** / **Path Drawing**
   - 基本路径语法、复杂图形创建
   - Basic path syntax, complex shape creation

5. **实用应用** / **Practical Applications**
   - 图标制作、简单图表
   - Icon creation, simple charts

### 🎯 下一步学习 / Next Steps

现在您可以继续学习：

Now you can continue learning:

1. **进阶教程** - 动画、高级API、自定义功能
   **Advanced Tutorial** - Animations, advanced APIs, custom features

2. **示例集合** - 更多实用示例和应用场景
   **Examples Collection** - More practical examples and use cases

3. **API参考** - 完整的接口文档
   **API Reference** - Complete interface documentation

### 💡 学习建议 / Learning Tips

1. **多练习** - 尝试创建不同类型的图形
   **Practice More** - Try creating different types of graphics

2. **实验样式** - 组合不同的颜色和样式
   **Experiment with Styles** - Combine different colors and styles

3. **查看源码** - 理解库的内部实现
   **Read Source Code** - Understand the library's internal implementation

4. **参考示例** - 学习最佳实践
   **Reference Examples** - Learn best practices

---

🎉 **恭喜完成基础教程！** / **Congratulations on completing the basic tutorial!**

您现在已经具备了创建复杂SVG图形的基础知识。继续探索更高级的功能，创造出更精彩的作品！

You now have the foundational knowledge to create complex SVG graphics. Continue exploring more advanced features to create even more amazing works!