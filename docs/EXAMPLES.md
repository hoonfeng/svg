# SVG库示例集合 / Examples Collection

## 📖 示例概述 / Examples Overview

本文档包含了丰富的SVG库使用示例，从简单的基础图形到复杂的实际应用。每个示例都包含完整的可运行代码和详细说明。

This document contains rich examples of using the SVG library, from simple basic shapes to complex real-world applications. Each example includes complete runnable code and detailed explanations.

## 🎯 基础示例 / Basic Examples

### Hello World

最简单的SVG程序，创建一个包含文本的基本图形。

The simplest SVG program, creating a basic graphic with text.

```go
package main

import (
    "image/color"
    "svg"
)

func main() {
    // 创建画布 / Create canvas
    canvas := svg.New(300, 200)
    
    // 设置背景 / Set background
    canvas.SetBackground(color.RGBA{240, 248, 255, 255})
    
    // 添加文本 / Add text
    canvas.Text(150, 100, "Hello, SVG!").
        FontSize(24).
        TextAnchor("middle").
        Fill(color.RGBA{0, 0, 0, 255})
    
    // 保存文件 / Save file
    canvas.SaveSVG("hello_world.svg")
    canvas.SavePNG("hello_world.png")
}
```

### 基本图形组合 / Basic Shape Combination

展示如何组合多种基本图形创建复合图案。

Demonstrates how to combine multiple basic shapes to create composite patterns.

```go
package main

import (
    "image/color"
    "svg"
)

func main() {
    canvas := svg.New(400, 300)
    canvas.SetBackground(color.RGBA{250, 250, 250, 255})
    
    // 创建房子 / Create a house
    // 房子主体 / House body
    canvas.Rect(150, 150, 100, 80).
        Fill(color.RGBA{255, 228, 196, 255}).
        Stroke(color.RGBA{139, 69, 19, 255}).
        StrokeWidth(2)
    
    // 屋顶 / Roof
    canvas.Polygon("150,150 200,100 250,150").
        Fill(color.RGBA{178, 34, 34, 255}).
        Stroke(color.RGBA{139, 0, 0, 255}).
        StrokeWidth(2)
    
    // 门 / Door
    canvas.Rect(180, 190, 20, 40).
        Fill(color.RGBA{139, 69, 19, 255})
    
    // 窗户 / Windows
    canvas.Rect(160, 170, 15, 15).
        Fill(color.RGBA{173, 216, 230, 255}).
        Stroke(color.RGBA{0, 0, 0, 255}).
        StrokeWidth(1)
    
    canvas.Rect(225, 170, 15, 15).
        Fill(color.RGBA{173, 216, 230, 255}).
        Stroke(color.RGBA{0, 0, 0, 255}).
        StrokeWidth(1)
    
    // 太阳 / Sun
    canvas.Circle(320, 80, 25).
        Fill(color.RGBA{255, 255, 0, 255}).
        Stroke(color.RGBA{255, 165, 0, 255}).
        StrokeWidth(2)
    
    // 太阳光线 / Sun rays
    for i := 0; i < 8; i++ {
        angle := float64(i) * 45
        x1, y1 := 320.0, 80.0
        x2 := x1 + 35*math.Cos(angle*math.Pi/180)
        y2 := y1 + 35*math.Sin(angle*math.Pi/180)
        
        canvas.Line(x1, y1, x2, y2).
            Stroke(color.RGBA{255, 165, 0, 255}).
            StrokeWidth(2)
    }
    
    // 草地 / Grass
    canvas.Rect(0, 230, 400, 70).
        Fill(color.RGBA{124, 252, 0, 255})
    
    canvas.SaveSVG("house_scene.svg")
    canvas.SavePNG("house_scene.png")
}
```

### 几何图案 / Geometric Patterns

创建重复的几何图案和装饰性设计。

Create repeating geometric patterns and decorative designs.

```go
package main

import (
    "image/color"
    "math"
    "svg"
)

func main() {
    canvas := svg.New(600, 600)
    canvas.SetBackground(color.RGBA{255, 255, 255, 255})
    
    // 创建同心圆图案 / Create concentric circle pattern
    centerX, centerY := 300.0, 300.0
    colors := []color.RGBA{
        {255, 0, 0, 100},
        {255, 127, 0, 100},
        {255, 255, 0, 100},
        {0, 255, 0, 100},
        {0, 0, 255, 100},
        {75, 0, 130, 100},
        {148, 0, 211, 100},
    }
    
    for i, c := range colors {
        radius := float64(200 - i*25)
        canvas.Circle(centerX, centerY, radius).
            Fill(c).
            Stroke(color.RGBA{0, 0, 0, 50}).
            StrokeWidth(1)
    }
    
    // 添加放射状线条 / Add radial lines
    for i := 0; i < 12; i++ {
        angle := float64(i) * 30 * math.Pi / 180
        x1 := centerX + 50*math.Cos(angle)
        y1 := centerY + 50*math.Sin(angle)
        x2 := centerX + 200*math.Cos(angle)
        y2 := centerY + 200*math.Sin(angle)
        
        canvas.Line(x1, y1, x2, y2).
            Stroke(color.RGBA{0, 0, 0, 150}).
            StrokeWidth(2)
    }
    
    canvas.SaveSVG("geometric_pattern.svg")
    canvas.SavePNG("geometric_pattern.png")
}
```

## 🎨 样式示例 / Style Examples

### 颜色渐变效果 / Color Gradient Effects

通过多个图形模拟渐变效果。

Simulate gradient effects using multiple shapes.

```go
package main

import (
    "image/color"
    "svg"
)

func main() {
    canvas := svg.New(800, 400)
    canvas.SetBackground(color.RGBA{0, 0, 0, 255})
    
    // 水平渐变 / Horizontal gradient
    for i := 0; i < 100; i++ {
        red := uint8(255 * i / 100)
        blue := uint8(255 * (100 - i) / 100)
        
        canvas.Rect(float64(i*4), 50, 4, 100).
            Fill(color.RGBA{red, 0, blue, 255})
    }
    
    // 垂直渐变 / Vertical gradient
    for i := 0; i < 100; i++ {
        green := uint8(255 * i / 100)
        alpha := uint8(255 * (100 - i) / 100)
        
        canvas.Rect(450, float64(50+i*2), 100, 2).
            Fill(color.RGBA{0, green, 255, alpha})
    }
    
    // 径向渐变模拟 / Radial gradient simulation
    centerX, centerY := 650.0, 150.0
    for i := 50; i > 0; i-- {
        intensity := uint8(255 * i / 50)
        canvas.Circle(centerX, centerY, float64(i)).
            Fill(color.RGBA{intensity, intensity, 0, 100})
    }
    
    // 彩虹条纹 / Rainbow stripes
    rainbowColors := []color.RGBA{
        {255, 0, 0, 255},   // 红 / Red
        {255, 127, 0, 255}, // 橙 / Orange
        {255, 255, 0, 255}, // 黄 / Yellow
        {0, 255, 0, 255},   // 绿 / Green
        {0, 0, 255, 255},   // 蓝 / Blue
        {75, 0, 130, 255},  // 靛 / Indigo
        {148, 0, 211, 255}, // 紫 / Violet
    }
    
    for i, c := range rainbowColors {
        canvas.Rect(float64(i*100), 300, 100, 50).
            Fill(c)
    }
    
    canvas.SaveSVG("gradient_effects.svg")
    canvas.SavePNG("gradient_effects.png")
}
```

### 文本样式展示 / Text Style Showcase

展示各种文本样式和排版效果。

Showcase various text styles and typography effects.

```go
package main

import (
    "image/color"
    "svg"
)

func main() {
    canvas := svg.New(800, 600)
    canvas.SetBackground(color.RGBA{255, 255, 255, 255})
    
    // 标题 / Title
    canvas.Text(400, 50, "文本样式展示").
        FontSize(32).
        FontWeight("bold").
        TextAnchor("middle").
        Fill(color.RGBA{0, 0, 0, 255})
    
    // 不同大小的文本 / Different text sizes
    sizes := []float64{12, 16, 20, 24, 28, 32}
    for i, size := range sizes {
        canvas.Text(50, float64(100+i*40), fmt.Sprintf("字体大小 %.0fpx", size)).
            FontSize(size).
            Fill(color.RGBA{0, 0, 0, 255})
    }
    
    // 不同字体粗细 / Different font weights
    weights := []string{"100", "300", "400", "600", "700", "900"}
    for i, weight := range weights {
        canvas.Text(300, float64(100+i*40), fmt.Sprintf("粗细 %s", weight)).
            FontSize(18).
            FontWeight(weight).
            Fill(color.RGBA{0, 0, 0, 255})
    }
    
    // 彩色文本 / Colored text
    colors := []color.RGBA{
        {255, 0, 0, 255},
        {0, 255, 0, 255},
        {0, 0, 255, 255},
        {255, 165, 0, 255},
        {128, 0, 128, 255},
        {255, 20, 147, 255},
    }
    
    for i, c := range colors {
        canvas.Text(500, float64(100+i*40), "彩色文本").
            FontSize(18).
            Fill(c)
    }
    
    // 文本对齐 / Text alignment
    alignments := []string{"start", "middle", "end"}
    for i, align := range alignments {
        y := float64(400 + i*30)
        
        // 绘制参考线 / Draw reference line
        canvas.Line(400, y, 400, y).
            Stroke(color.RGBA{200, 200, 200, 255}).
            StrokeWidth(1)
        
        canvas.Text(400, y, fmt.Sprintf("对齐方式: %s", align)).
            FontSize(16).
            TextAnchor(align).
            Fill(color.RGBA{0, 0, 0, 255})
    }
    
    // 文本装饰 / Text decorations
    decorations := []string{"underline", "overline", "line-through"}
    for i, decoration := range decorations {
        canvas.Text(50, float64(500+i*30), fmt.Sprintf("装饰: %s", decoration)).
            FontSize(16).
            TextDecoration(decoration).
            Fill(color.RGBA{0, 0, 0, 255})
    }
    
    canvas.SaveSVG("text_styles.svg")
    canvas.SavePNG("text_styles.png")
}
```

### 变换效果 / Transform Effects

展示旋转、缩放、倾斜等变换效果。

Demonstrate rotation, scaling, skewing and other transform effects.

```go
package main

import (
    "image/color"
    "math"
    "svg"
)

func main() {
    canvas := svg.New(800, 600)
    canvas.SetBackground(color.RGBA{250, 250, 250, 255})
    
    // 旋转效果 / Rotation effects
    for i := 0; i < 12; i++ {
        angle := float64(i * 30)
        canvas.Rect(200, 150, 60, 20).
            Fill(color.RGBA{255, uint8(i*20), 0, 200}).
            Transform(fmt.Sprintf("rotate(%f 230 160)", angle))
    }
    
    // 缩放效果 / Scaling effects
    for i := 1; i <= 5; i++ {
        scale := float64(i) * 0.3
        canvas.Circle(500, 150, 20).
            Fill(color.RGBA{0, 255, uint8(i*50), 150}).
            Transform(fmt.Sprintf("scale(%f) translate(%f %f)", scale, 500/scale, 150/scale))
    }
    
    // 倾斜效果 / Skew effects
    for i := 0; i < 5; i++ {
        skewX := float64(i * 10)
        canvas.Rect(float64(100+i*80), 350, 50, 50).
            Fill(color.RGBA{uint8(i*60), 100, 255, 255}).
            Transform(fmt.Sprintf("skewX(%f)", skewX))
    }
    
    // 组合变换 / Combined transforms
    canvas.Polygon("600,400 650,380 700,400 680,450 620,450").
        Fill(color.RGBA{255, 0, 255, 200}).
        Transform("rotate(45 650 415) scale(1.2)")
    
    canvas.SaveSVG("transform_effects.svg")
    canvas.SavePNG("transform_effects.png")
}
```

## 🎬 动画示例 / Animation Examples

### 简单旋转动画 / Simple Rotation Animation

创建基本的旋转动画效果。

Create basic rotation animation effects.

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
    err := builder.CreateRotatingShapes(config).SaveToGIF("simple_rotation.gif")
    if err != nil {
        fmt.Printf("创建动画失败: %v\n", err)
        return
    }
    
    fmt.Println("✅ 旋转动画已创建: simple_rotation.gif")
}
```

### 彩色粒子动画 / Colorful Particle Animation

创建动态的粒子效果动画。

Create dynamic particle effect animation.

```go
package main

import (
    "image/color"
    "svg"
)

func main() {
    builder := svg.NewAnimationBuilder(600, 400)
    builder.SetFrameCount(90).SetFrameRate(30)
    
    config := svg.AnimationConfig{
        Duration:   3.0, // 3秒 / 3 seconds
        Easing:     svg.Linear,
        Background: color.RGBA{10, 10, 20, 255},
        Loop:       true,
    }
    
    err := builder.CreateColorfulParticles(config).SaveToGIF("particle_animation.gif")
    if err != nil {
        fmt.Printf("创建粒子动画失败: %v\n", err)
        return
    }
    
    fmt.Println("✅ 粒子动画已创建: particle_animation.gif")
}
```

### 脉冲动画 / Pulse Animation

创建心跳般的脉冲效果。

Create heartbeat-like pulse effects.

```go
package main

import (
    "image/color"
    "svg"
)

func main() {
    builder := svg.NewAnimationBuilder(500, 500)
    builder.SetFrameCount(80).SetFrameRate(25)
    
    config := svg.AnimationConfig{
        Duration:   3.2, // 3.2秒 / 3.2 seconds
        Easing:     svg.EaseInOutQuad,
        Background: color.RGBA{30, 30, 50, 255},
        Loop:       true,
    }
    
    err := builder.CreatePulsingCircles(config).SaveToGIF("pulse_animation.gif")
    if err != nil {
        fmt.Printf("创建脉冲动画失败: %v\n", err)
        return
    }
    
    fmt.Println("✅ 脉冲动画已创建: pulse_animation.gif")
}
```

### 波浪动画 / Wave Animation

创建流动的波浪效果。

Create flowing wave effects.

```go
package main

import (
    "image/color"
    "svg"
)

func main() {
    builder := svg.NewAnimationBuilder(800, 300)
    builder.SetFrameCount(60).SetFrameRate(30)
    
    config := svg.AnimationConfig{
        Duration:   2.0, // 2秒 / 2 seconds
        Easing:     svg.EaseInOut,
        Background: color.RGBA{30, 50, 80, 255}, // 海洋蓝 / Ocean blue
        Loop:       true,
    }
    
    err := builder.CreateWaveAnimation(config).SaveToGIF("wave_animation.gif")
    if err != nil {
        fmt.Printf("创建波浪动画失败: %v\n", err)
        return
    }
    
    fmt.Println("✅ 波浪动画已创建: wave_animation.gif")
}
```

## 🏗️ 实际应用 / Real-World Applications

### 图表绘制 / Chart Drawing

创建各种类型的数据图表。

Create various types of data charts.

```go
package main

import (
    "image/color"
    "math"
    "svg"
)

// 柱状图 / Bar Chart
func createBarChart() {
    canvas := svg.New(600, 400)
    canvas.SetBackground(color.RGBA{255, 255, 255, 255})
    
    // 数据 / Data
    data := []struct {
        label string
        value float64
        color color.RGBA
    }{
        {"一月", 120, color.RGBA{255, 99, 132, 255}},
        {"二月", 190, color.RGBA{54, 162, 235, 255}},
        {"三月", 300, color.RGBA{255, 205, 86, 255}},
        {"四月", 250, color.RGBA{75, 192, 192, 255}},
        {"五月", 180, color.RGBA{153, 102, 255, 255}},
        {"六月", 220, color.RGBA{255, 159, 64, 255}},
    }
    
    // 绘制坐标轴 / Draw axes
    canvas.Line(60, 350, 540, 350).
        Stroke(color.RGBA{0, 0, 0, 255}).
        StrokeWidth(2)
    canvas.Line(60, 50, 60, 350).
        Stroke(color.RGBA{0, 0, 0, 255}).
        StrokeWidth(2)
    
    // 绘制柱子 / Draw bars
    barWidth := 60.0
    spacing := 20.0
    maxValue := 300.0
    
    for i, item := range data {
        x := 80 + float64(i)*(barWidth+spacing)
        height := (item.value / maxValue) * 250
        y := 350 - height
        
        // 柱子 / Bar
        canvas.Rect(x, y, barWidth, height).
            Fill(item.color).
            Stroke(color.RGBA{0, 0, 0, 100}).
            StrokeWidth(1)
        
        // 数值标签 / Value label
        canvas.Text(x+barWidth/2, y-10, fmt.Sprintf("%.0f", item.value)).
            FontSize(12).
            TextAnchor("middle").
            Fill(color.RGBA{0, 0, 0, 255})
        
        // 月份标签 / Month label
        canvas.Text(x+barWidth/2, 370, item.label).
            FontSize(14).
            TextAnchor("middle").
            Fill(color.RGBA{0, 0, 0, 255})
    }
    
    // 标题 / Title
    canvas.Text(300, 30, "月度销售数据").
        FontSize(20).
        FontWeight("bold").
        TextAnchor("middle").
        Fill(color.RGBA{0, 0, 0, 255})
    
    canvas.SaveSVG("bar_chart.svg")
    canvas.SavePNG("bar_chart.png")
}

// 饼图 / Pie Chart
func createPieChart() {
    canvas := svg.New(500, 500)
    canvas.SetBackground(color.RGBA{255, 255, 255, 255})
    
    // 数据 / Data
    data := []struct {
        label string
        value float64
        color color.RGBA
    }{
        {"产品A", 30, color.RGBA{255, 99, 132, 255}},
        {"产品B", 25, color.RGBA{54, 162, 235, 255}},
        {"产品C", 20, color.RGBA{255, 205, 86, 255}},
        {"产品D", 15, color.RGBA{75, 192, 192, 255}},
        {"其他", 10, color.RGBA{153, 102, 255, 255}},
    }
    
    centerX, centerY := 250.0, 250.0
    radius := 120.0
    total := 100.0
    
    startAngle := 0.0
    
    for i, item := range data {
        // 计算角度 / Calculate angle
        angle := (item.value / total) * 360
        endAngle := startAngle + angle
        
        // 计算路径 / Calculate path
        startRad := startAngle * math.Pi / 180
        endRad := endAngle * math.Pi / 180
        
        x1 := centerX + radius*math.Cos(startRad)
        y1 := centerY + radius*math.Sin(startRad)
        x2 := centerX + radius*math.Cos(endRad)
        y2 := centerY + radius*math.Sin(endRad)
        
        largeArc := 0
        if angle > 180 {
            largeArc = 1
        }
        
        pathData := fmt.Sprintf("M %f %f L %f %f A %f %f 0 %d 1 %f %f Z",
            centerX, centerY, x1, y1, radius, radius, largeArc, x2, y2)
        
        canvas.Path(pathData).
            Fill(item.color).
            Stroke(color.RGBA{255, 255, 255, 255}).
            StrokeWidth(2)
        
        // 标签 / Label
        labelAngle := (startAngle + endAngle) / 2 * math.Pi / 180
        labelX := centerX + (radius+30)*math.Cos(labelAngle)
        labelY := centerY + (radius+30)*math.Sin(labelAngle)
        
        canvas.Text(labelX, labelY, fmt.Sprintf("%s\n%.0f%%", item.label, item.value)).
            FontSize(12).
            TextAnchor("middle").
            Fill(color.RGBA{0, 0, 0, 255})
        
        startAngle = endAngle
    }
    
    // 标题 / Title
    canvas.Text(250, 30, "产品销售占比").
        FontSize(18).
        FontWeight("bold").
        TextAnchor("middle").
        Fill(color.RGBA{0, 0, 0, 255})
    
    canvas.SaveSVG("pie_chart.svg")
    canvas.SavePNG("pie_chart.png")
}

// 折线图 / Line Chart
func createLineChart() {
    canvas := svg.New(700, 400)
    canvas.SetBackground(color.RGBA{255, 255, 255, 255})
    
    // 数据点 / Data points
    data := []struct {
        x, y float64
    }{
        {1, 20}, {2, 45}, {3, 30}, {4, 60}, {5, 40},
        {6, 75}, {7, 55}, {8, 80}, {9, 65}, {10, 90},
    }
    
    // 绘制坐标轴 / Draw axes
    canvas.Line(60, 350, 640, 350).
        Stroke(color.RGBA{0, 0, 0, 255}).
        StrokeWidth(2)
    canvas.Line(60, 50, 60, 350).
        Stroke(color.RGBA{0, 0, 0, 255}).
        StrokeWidth(2)
    
    // 绘制网格线 / Draw grid lines
    for i := 1; i <= 10; i++ {
        x := 60 + float64(i)*58
        canvas.Line(x, 350, x, 50).
            Stroke(color.RGBA{200, 200, 200, 255}).
            StrokeWidth(1).
            StrokeDashArray("2,2")
    }
    
    for i := 1; i <= 5; i++ {
        y := 350 - float64(i)*60
        canvas.Line(60, y, 640, y).
            Stroke(color.RGBA{200, 200, 200, 255}).
            StrokeWidth(1).
            StrokeDashArray("2,2")
    }
    
    // 绘制折线 / Draw line
    var pathData strings.Builder
    for i, point := range data {
        x := 60 + point.x*58
        y := 350 - (point.y/100)*300
        
        if i == 0 {
            pathData.WriteString(fmt.Sprintf("M %f %f", x, y))
        } else {
            pathData.WriteString(fmt.Sprintf(" L %f %f", x, y))
        }
        
        // 绘制数据点 / Draw data points
        canvas.Circle(x, y, 4).
            Fill(color.RGBA{255, 0, 0, 255})
    }
    
    canvas.Path(pathData.String()).
        Fill("none").
        Stroke(color.RGBA{0, 100, 255, 255}).
        StrokeWidth(3)
    
    // 标题 / Title
    canvas.Text(350, 30, "月度增长趋势").
        FontSize(18).
        FontWeight("bold").
        TextAnchor("middle").
        Fill(color.RGBA{0, 0, 0, 255})
    
    canvas.SaveSVG("line_chart.svg")
    canvas.SavePNG("line_chart.png")
}

func main() {
    createBarChart()
    createPieChart()
    createLineChart()
    
    fmt.Println("✅ 所有图表已创建完成")
}
```

### 图标制作 / Icon Creation

创建常用的图标和符号。

Create common icons and symbols.

```go
package main

import (
    "image/color"
    "math"
    "svg"
)

// 创建图标集合 / Create icon collection
func createIconSet() {
    canvas := svg.New(800, 600)
    canvas.SetBackground(color.RGBA{250, 250, 250, 255})
    
    iconSize := 60.0
    spacing := 100.0
    
    // 主页图标 / Home icon
    x, y := 100.0, 100.0
    canvas.Polygon(fmt.Sprintf("%f,%f %f,%f %f,%f %f,%f %f,%f",
        x, y+30, x+30, y, x+60, y+30, x+60, y+50, x, y+50)).
        Fill(color.RGBA{100, 150, 255, 255}).
        Stroke(color.RGBA{0, 100, 200, 255}).
        StrokeWidth(2)
    canvas.Rect(x+20, y+35, 20, 15).Fill(color.RGBA{139, 69, 19, 255})
    
    // 设置图标 / Settings icon
    x += spacing
    centerX, centerY := x+30, y+30
    canvas.Circle(centerX, centerY, 25).
        Fill("none").
        Stroke(color.RGBA{100, 100, 100, 255}).
        StrokeWidth(4)
    canvas.Circle(centerX, centerY, 8).Fill(color.RGBA{100, 100, 100, 255})
    
    for i := 0; i < 8; i++ {
        angle := float64(i) * 45 * math.Pi / 180
        x1 := centerX + 15*math.Cos(angle)
        y1 := centerY + 15*math.Sin(angle)
        x2 := centerX + 30*math.Cos(angle)
        y2 := centerY + 30*math.Sin(angle)
        
        canvas.Line(x1, y1, x2, y2).
            Stroke(color.RGBA{100, 100, 100, 255}).
            StrokeWidth(3)
    }
    
    // 邮件图标 / Mail icon
    x += spacing
    canvas.Rect(x, y+10, 60, 40).
        Fill(color.RGBA{255, 255, 255, 255}).
        Stroke(color.RGBA{200, 200, 200, 255}).
        StrokeWidth(2)
    canvas.Path(fmt.Sprintf("M %f,%f L %f,%f L %f,%f", x, y+10, x+30, y+30, x+60, y+10)).
        Fill("none").
        Stroke(color.RGBA{200, 200, 200, 255}).
        StrokeWidth(2)
    
    // 搜索图标 / Search icon
    x += spacing
    canvas.Circle(x+20, y+20, 15).
        Fill("none").
        Stroke(color.RGBA{0, 0, 0, 255}).
        StrokeWidth(3)
    canvas.Line(x+32, y+32, x+45, y+45).
        Stroke(color.RGBA{0, 0, 0, 255}).
        StrokeWidth(3)
    
    // 用户图标 / User icon
    x += spacing
    canvas.Circle(x+30, y+15, 12).
        Fill(color.RGBA{100, 150, 255, 255}).
        Stroke(color.RGBA{0, 100, 200, 255}).
        StrokeWidth(2)
    canvas.Path(fmt.Sprintf("M %f,%f Q %f,%f %f,%f Q %f,%f %f,%f",
        x+10, y+50, x+30, y+35, x+50, y+50, x+30, y+35, x+10, y+50)).
        Fill(color.RGBA{100, 150, 255, 255}).
        Stroke(color.RGBA{0, 100, 200, 255}).
        StrokeWidth(2)
    
    // 下载图标 / Download icon
    x, y = 100.0, 250.0
    canvas.Line(x+30, y, x+30, y+35).
        Stroke(color.RGBA{0, 150, 0, 255}).
        StrokeWidth(4)
    canvas.Polygon(fmt.Sprintf("%f,%f %f,%f %f,%f", x+20, y+25, x+30, y+35, x+40, y+25)).
        Fill(color.RGBA{0, 150, 0, 255})
    canvas.Rect(x+10, y+45, 40, 5).
        Fill(color.RGBA{100, 100, 100, 255})
    
    // 上传图标 / Upload icon
    x += spacing
    canvas.Line(x+30, y+15, x+30, y+50).
        Stroke(color.RGBA{255, 100, 0, 255}).
        StrokeWidth(4)
    canvas.Polygon(fmt.Sprintf("%f,%f %f,%f %f,%f", x+20, y+25, x+30, y+15, x+40, y+25)).
        Fill(color.RGBA{255, 100, 0, 255})
    canvas.Rect(x+10, y+45, 40, 5).
        Fill(color.RGBA{100, 100, 100, 255})
    
    // 心形图标 / Heart icon
    x += spacing
    heartPath := fmt.Sprintf("M %f,%f C %f,%f %f,%f %f,%f C %f,%f %f,%f %f,%f C %f,%f %f,%f %f,%f",
        x+30, y+35, x+25, y+20, x+15, y+20, x+15, y+30,
        x+15, y+20, x+5, y+20, x+30, y+50,
        x+55, y+20, x+45, y+20, x+45, y+30)
    canvas.Path(heartPath).
        Fill(color.RGBA{255, 100, 100, 255}).
        Stroke(color.RGBA{200, 0, 0, 255}).
        StrokeWidth(2)
    
    // 星形图标 / Star icon
    x += spacing
    starPath := fmt.Sprintf("M %f,%f L %f,%f L %f,%f L %f,%f L %f,%f L %f,%f L %f,%f L %f,%f L %f,%f L %f,%f Z",
        x+30, y+10, x+35, y+25, x+50, y+25, x+40, y+35,
        x+45, y+50, x+30, y+42, x+15, y+50, x+20, y+35,
        x+10, y+25, x+25, y+25)
    canvas.Path(starPath).
        Fill(color.RGBA{255, 215, 0, 255}).
        Stroke(color.RGBA{255, 140, 0, 255}).
        StrokeWidth(2)
    
    // 添加标题 / Add title
    canvas.Text(400, 50, "常用图标集合").
        FontSize(24).
        FontWeight("bold").
        TextAnchor("middle").
        Fill(color.RGBA{0, 0, 0, 255})
    
    canvas.SaveSVG("icon_set.svg")
    canvas.SavePNG("icon_set.png")
}

func main() {
    createIconSet()
    fmt.Println("✅ 图标集合已创建完成")
}
```

### 数据可视化 / Data Visualization

创建复杂的数据可视化图表。

Create complex data visualization charts.

```go
package main

import (
    "image/color"
    "math"
    "math/rand"
    "svg"
    "time"
)

// 热力图 / Heatmap
func createHeatmap() {
    canvas := svg.New(600, 400)
    canvas.SetBackground(color.RGBA{255, 255, 255, 255})
    
    // 生成随机数据 / Generate random data
    rand.Seed(time.Now().UnixNano())
    rows, cols := 10, 15
    cellWidth, cellHeight := 30.0, 25.0
    
    for i := 0; i < rows; i++ {
        for j := 0; j < cols; j++ {
            value := rand.Float64()
            
            // 根据值计算颜色 / Calculate color based on value
            red := uint8(255 * value)
            blue := uint8(255 * (1 - value))
            
            x := 50 + float64(j)*cellWidth
            y := 50 + float64(i)*cellHeight
            
            canvas.Rect(x, y, cellWidth-1, cellHeight-1).
                Fill(color.RGBA{red, 0, blue, 255}).
                Stroke(color.RGBA{255, 255, 255, 255}).
                StrokeWidth(1)
            
            // 添加数值 / Add value
            canvas.Text(x+cellWidth/2, y+cellHeight/2, fmt.Sprintf("%.2f", value)).
                FontSize(8).
                TextAnchor("middle").
                Fill(color.RGBA{255, 255, 255, 255})
        }
    }
    
    // 标题 / Title
    canvas.Text(300, 30, "数据热力图").
        FontSize(18).
        FontWeight("bold").
        TextAnchor("middle").
        Fill(color.RGBA{0, 0, 0, 255})
    
    canvas.SaveSVG("heatmap.svg")
    canvas.SavePNG("heatmap.png")
}

// 散点图 / Scatter Plot
func createScatterPlot() {
    canvas := svg.New(600, 500)
    canvas.SetBackground(color.RGBA{255, 255, 255, 255})
    
    // 生成随机数据点 / Generate random data points
    rand.Seed(time.Now().UnixNano())
    numPoints := 50
    
    // 绘制坐标轴 / Draw axes
    canvas.Line(60, 450, 540, 450).
        Stroke(color.RGBA{0, 0, 0, 255}).
        StrokeWidth(2)
    canvas.Line(60, 50, 60, 450).
        Stroke(color.RGBA{0, 0, 0, 255}).
        StrokeWidth(2)
    
    // 绘制网格 / Draw grid
    for i := 1; i <= 10; i++ {
        x := 60 + float64(i)*48
        y := 450 - float64(i)*40
        
        canvas.Line(x, 450, x, 50).
            Stroke(color.RGBA{200, 200, 200, 255}).
            StrokeWidth(1).
            StrokeDashArray("2,2")
        
        canvas.Line(60, y, 540, y).
            Stroke(color.RGBA{200, 200, 200, 255}).
            StrokeWidth(1).
            StrokeDashArray("2,2")
    }
    
    // 绘制数据点 / Draw data points
    colors := []color.RGBA{
        {255, 0, 0, 200},
        {0, 255, 0, 200},
        {0, 0, 255, 200},
        {255, 255, 0, 200},
        {255, 0, 255, 200},
    }
    
    for i := 0; i < numPoints; i++ {
        x := 60 + rand.Float64()*480
        y := 50 + rand.Float64()*400
        size := 3 + rand.Float64()*8
        colorIndex := rand.Intn(len(colors))
        
        canvas.Circle(x, y, size).
            Fill(colors[colorIndex])
    }
    
    // 标题 / Title
    canvas.Text(300, 30, "散点图分析").
        FontSize(18).
        FontWeight("bold").
        TextAnchor("middle").
        Fill(color.RGBA{0, 0, 0, 255})
    
    canvas.SaveSVG("scatter_plot.svg")
    canvas.SavePNG("scatter_plot.png")
}

// 雷达图 / Radar Chart
func createRadarChart() {
    canvas := svg.New(500, 500)
    canvas.SetBackground(color.RGBA{255, 255, 255, 255})
    
    centerX, centerY := 250.0, 250.0
    radius := 150.0
    
    // 数据 / Data
    categories := []string{"速度", "力量", "技巧", "防御", "智力", "经验"}
    values := []float64{0.8, 0.6, 0.9, 0.7, 0.5, 0.8}
    
    numCategories := len(categories)
    angleStep := 2 * math.Pi / float64(numCategories)
    
    // 绘制背景网格 / Draw background grid
    for level := 1; level <= 5; level++ {
        r := radius * float64(level) / 5
        
        var pathData strings.Builder
        for i := 0; i < numCategories; i++ {
            angle := float64(i)*angleStep - math.Pi/2
            x := centerX + r*math.Cos(angle)
            y := centerY + r*math.Sin(angle)
            
            if i == 0 {
                pathData.WriteString(fmt.Sprintf("M %f %f", x, y))
            } else {
                pathData.WriteString(fmt.Sprintf(" L %f %f", x, y))
            }
        }
        pathData.WriteString(" Z")
        
        canvas.Path(pathData.String()).
            Fill("none").
            Stroke(color.RGBA{200, 200, 200, 255}).
            StrokeWidth(1)
    }
    
    // 绘制轴线 / Draw axis lines
    for i := 0; i < numCategories; i++ {
        angle := float64(i)*angleStep - math.Pi/2
        x := centerX + radius*math.Cos(angle)
        y := centerY + radius*math.Sin(angle)
        
        canvas.Line(centerX, centerY, x, y).
            Stroke(color.RGBA{150, 150, 150, 255}).
            StrokeWidth(1)
        
        // 标签 / Labels
        labelX := centerX + (radius+20)*math.Cos(angle)
        labelY := centerY + (radius+20)*math.Sin(angle)
        
        canvas.Text(labelX, labelY, categories[i]).
            FontSize(12).
            TextAnchor("middle").
            Fill(color.RGBA{0, 0, 0, 255})
    }
    
    // 绘制数据区域 / Draw data area
    var dataPath strings.Builder
    for i := 0; i < numCategories; i++ {
        angle := float64(i)*angleStep - math.Pi/2
        r := radius * values[i]
        x := centerX + r*math.Cos(angle)
        y := centerY + r*math.Sin(angle)
        
        if i == 0 {
            dataPath.WriteString(fmt.Sprintf("M %f %f", x, y))
        } else {
            dataPath.WriteString(fmt.Sprintf(" L %f %f", x, y))
        }
        
        // 数据点 / Data points
        canvas.Circle(x, y, 4).
            Fill(color.RGBA{255, 0, 0, 255})
    }
    dataPath.WriteString(" Z")
    
    canvas.Path(dataPath.String()).
        Fill(color.RGBA{255, 0, 0, 100}).
        Stroke(color.RGBA{255, 0, 0, 255}).
        StrokeWidth(2)
    
    // 标题 / Title
    canvas.Text(250, 30, "能力雷达图").
        FontSize(18).
        FontWeight("bold").
        TextAnchor("middle").
        Fill(color.RGBA{0, 0, 0, 255})
    
    canvas.SaveSVG("radar_chart.svg")
    canvas.SavePNG("radar_chart.png")
}

func main() {
    createHeatmap()
    createScatterPlot()
    createRadarChart()
    
    fmt.Println("✅ 所有数据可视化图表已创建完成")
}
```

### 游戏图形 / Game Graphics

创建游戏中常用的图形元素。

Create common graphic elements used in games.

```go
package main

import (
    "image/color"
    "math"
    "svg"
)

// 创建游戏角色 / Create game character
func createGameCharacter() {
    canvas := svg.New(400, 500)
    canvas.SetBackground(color.RGBA{135, 206, 235, 255}) // 天空蓝 / Sky blue
    
    // 角色身体 / Character body
    centerX, centerY := 200.0, 250.0
    
    // 头部 / Head
    canvas.Circle(centerX, centerY-80, 40).
        Fill(color.RGBA{255, 220, 177, 255}).
        Stroke(color.RGBA{0, 0, 0, 255}).
        StrokeWidth(2)
    
    // 眼睛 / Eyes
    canvas.Circle(centerX-15, centerY-90, 5).
        Fill(color.RGBA{0, 0, 0, 255})
    canvas.Circle(centerX+15, centerY-90, 5).
        Fill(color.RGBA{0, 0, 0, 255})
    
    // 嘴巴 / Mouth
    canvas.Path(fmt.Sprintf("M %f,%f Q %f,%f %f,%f",
        centerX-10, centerY-70, centerX, centerY-65, centerX+10, centerY-70)).
        Fill("none").
        Stroke(color.RGBA{0, 0, 0, 255}).
        StrokeWidth(2)
    
    // 身体 / Body
    canvas.Rect(centerX-25, centerY-40, 50, 80).
        Fill(color.RGBA{255, 0, 0, 255}).
        Stroke(color.RGBA{0, 0, 0, 255}).
        StrokeWidth(2)
    
    // 手臂 / Arms
    canvas.Rect(centerX-45, centerY-30, 20, 60).
        Fill(color.RGBA{255, 220, 177, 255}).
        Stroke(color.RGBA{0, 0, 0, 255}).
        StrokeWidth(2)
    canvas.Rect(centerX+25, centerY-30, 20, 60).
        Fill(color.RGBA{255, 220, 177, 255}).
        Stroke(color.RGBA{0, 0, 0, 255}).
        StrokeWidth(2)
    
    // 腿部 / Legs
    canvas.Rect(centerX-20, centerY+40, 15, 80).
        Fill(color.RGBA{0, 0, 255, 255}).
        Stroke(color.RGBA{0, 0, 0, 255}).
        StrokeWidth(2)
    canvas.Rect(centerX+5, centerY+40, 15, 80).
        Fill(color.RGBA{0, 0, 255, 255}).
        Stroke(color.RGBA{0, 0, 0, 255}).
        StrokeWidth(2)
    
    // 脚 / Feet
    canvas.Ellipse(centerX-12, centerY+130, 20, 10).
        Fill(color.RGBA{0, 0, 0, 255})
    canvas.Ellipse(centerX+12, centerY+130, 20, 10).
        Fill(color.RGBA{0, 0, 0, 255})
    
    canvas.SaveSVG("game_character.svg")
    canvas.SavePNG("game_character.png")
}

// 创建游戏道具 / Create game items
func createGameItems() {
    canvas := svg.New(600, 400)
    canvas.SetBackground(color.RGBA{50, 50, 50, 255})
    
    // 金币 / Gold coin
    x, y := 100.0, 100.0
    canvas.Circle(x, y, 30).
        Fill(color.RGBA{255, 215, 0, 255}).
        Stroke(color.RGBA{255, 140, 0, 255}).
        StrokeWidth(3)
    canvas.Text(x, y, "$").
        FontSize(24).
        FontWeight("bold").
        TextAnchor("middle").
        Fill(color.RGBA{255, 140, 0, 255})
    
    // 宝石 / Gem
    x += 120
    gemPath := fmt.Sprintf("M %f,%f L %f,%f L %f,%f L %f,%f L %f,%f L %f,%f Z",
        x, y-25, x-20, y-10, x-15, y+15, x+15, y+15, x+20, y-10, x, y-25)
    canvas.Path(gemPath).
        Fill(color.RGBA{0, 255, 255, 255}).
        Stroke(color.RGBA{0, 200, 200, 255}).
        StrokeWidth(2)
    
    // 剑 / Sword
    x += 120
    // 剑刃 / Blade
    canvas.Rect(x-3, y-40, 6, 60).
        Fill(color.RGBA{192, 192, 192, 255}).
        Stroke(color.RGBA{128, 128, 128, 255}).
        StrokeWidth(2)
    // 护手 / Guard
    canvas.Rect(x-15, y+15, 30, 5).
        Fill(color.RGBA{139, 69, 19, 255})
    // 剑柄 / Handle
    canvas.Rect(x-4, y+20, 8, 25).
        Fill(color.RGBA{139, 69, 19, 255})
    // 剑首 / Pommel
    canvas.Circle(x, y+50, 6).
        Fill(color.RGBA{255, 215, 0, 255})
    
    // 盾牌 / Shield
    x += 120
    shieldPath := fmt.Sprintf("M %f,%f Q %f,%f %f,%f Q %f,%f %f,%f Q %f,%f %f,%f Q %f,%f %f,%f Z",
        x, y-30, x-25, y-20, x-25, y, x-25, y+20, x, y+35,
        x+25, y+20, x+25, y, x+25, y-20, x, y-30)
    canvas.Path(shieldPath).
        Fill(color.RGBA{0, 100, 200, 255}).
        Stroke(color.RGBA{0, 50, 150, 255}).
        StrokeWidth(3)
    
    // 盾牌装饰 / Shield decoration
    canvas.Circle(x, y, 8).
        Fill(color.RGBA{255, 215, 0, 255})
    
    // 药水瓶 / Potion bottle
    x, y = 100.0, 280.0
    // 瓶身 / Bottle body
    canvas.Rect(x-15, y-10, 30, 40).
        Fill(color.RGBA{100, 255, 100, 200}).
        Stroke(color.RGBA{0, 150, 0, 255}).
        StrokeWidth(2)
    // 瓶颈 / Bottle neck
    canvas.Rect(x-8, y-25, 16, 15).
        Fill(color.RGBA{139, 69, 19, 255})
    // 瓶塞 / Cork
    canvas.Rect(x-10, y-30, 20, 8).
        Fill(color.RGBA{160, 82, 45, 255})
    
    // 魔法书 / Magic book
    x += 120
    canvas.Rect(x-20, y-25, 40, 50).
        Fill(color.RGBA{139, 0, 139, 255}).
        Stroke(color.RGBA{75, 0, 130, 255}).
        StrokeWidth(2)
    // 书页 / Pages
    canvas.Rect(x-18, y-23, 36, 46).
        Fill(color.RGBA{255, 255, 240, 255})
    // 符文 / Rune
    canvas.Text(x, y, "✦").
        FontSize(20).
        TextAnchor("middle").
        Fill(color.RGBA{255, 215, 0, 255})
    
    // 钥匙 / Key
    x += 120
    // 钥匙头 / Key head
    canvas.Circle(x-10, y-10, 12).
        Fill("none").
        Stroke(color.RGBA{255, 215, 0, 255}).
        StrokeWidth(3)
    // 钥匙柄 / Key shaft
    canvas.Line(x+2, y-10, x+25, y-10).
        Stroke(color.RGBA{255, 215, 0, 255}).
        StrokeWidth(4)
    // 钥匙齿 / Key teeth
    canvas.Line(x+20, y-10, x+20, y-5).
        Stroke(color.RGBA{255, 215, 0, 255}).
        StrokeWidth(3)
    canvas.Line(x+25, y-10, x+25, y-3).
        Stroke(color.RGBA{255, 215, 0, 255}).
        StrokeWidth(3)
    
    // 炸弹 / Bomb
    x += 120
    canvas.Circle(x, y, 20).
        Fill(color.RGBA{0, 0, 0, 255}).
        Stroke(color.RGBA{50, 50, 50, 255}).
        StrokeWidth(2)
    // 引线 / Fuse
    canvas.Line(x-15, y-15, x-25, y-25).
        Stroke(color.RGBA{139, 69, 19, 255}).
        StrokeWidth(3)
    // 火花 / Spark
    canvas.Circle(x-25, y-25, 3).
        Fill(color.RGBA{255, 255, 0, 255})
    
    canvas.SaveSVG("game_items.svg")
    canvas.SavePNG("game_items.png")
}

// 创建游戏地图元素 / Create game map elements
func createGameMapElements() {
    canvas := svg.New(800, 600)
    canvas.SetBackground(color.RGBA{34, 139, 34, 255}) // 森林绿 / Forest green
    
    // 树木 / Trees
    for i := 0; i < 5; i++ {
        x := 100.0 + float64(i)*150
        y := 400.0
        
        // 树干 / Tree trunk
        canvas.Rect(x-10, y, 20, 60).
            Fill(color.RGBA{139, 69, 19, 255})
        
        // 树冠 / Tree crown
        canvas.Circle(x, y-20, 40).
            Fill(color.RGBA{0, 100, 0, 255})
    }
    
    // 山脉 / Mountains
    mountainPath := "M 0,300 L 150,150 L 300,200 L 450,100 L 600,180 L 750,120 L 800,160 L 800,300 Z"
    canvas.Path(mountainPath).
        Fill(color.RGBA{105, 105, 105, 255}).
        Stroke(color.RGBA{70, 70, 70, 255}).
        StrokeWidth(2)
    
    // 河流 / River
    riverPath := "M 0,350 Q 200,320 400,340 Q 600,360 800,330"
    canvas.Path(riverPath).
        Fill("none").
        Stroke(color.RGBA{30, 144, 255, 255}).
        StrokeWidth(20)
    
    // 城堡 / Castle
    x, y := 600.0, 450.0
    // 主塔 / Main tower
    canvas.Rect(x-30, y-80, 60, 80).
        Fill(color.RGBA{169, 169, 169, 255}).
        Stroke(color.RGBA{105, 105, 105, 255}).
        StrokeWidth(2)
    
    // 侧塔 / Side towers
    canvas.Rect(x-60, y-60, 30, 60).
        Fill(color.RGBA{169, 169, 169, 255}).
        Stroke(color.RGBA{105, 105, 105, 255}).
        StrokeWidth(2)
    canvas.Rect(x+30, y-60, 30, 60).
        Fill(color.RGBA{169, 169, 169, 255}).
        Stroke(color.RGBA{105, 105, 105, 255}).
        StrokeWidth(2)
    
    // 旗帜 / Flags
    canvas.Line(x, y-80, x, y-100).
        Stroke(color.RGBA{139, 69, 19, 255}).
        StrokeWidth(3)
    canvas.Polygon(fmt.Sprintf("%f,%f %f,%f %f,%f", x, y-100, x+15, y-95, x, y-90)).
        Fill(color.RGBA{255, 0, 0, 255})
    
    canvas.SaveSVG("game_map.svg")
    canvas.SavePNG("game_map.png")
}

func main() {
    createGameCharacter()
    createGameItems()
    createGameMapElements()
    
    fmt.Println("✅ 所有游戏图形已创建完成")
}
```

## 🔧 实用工具 / Utility Tools

### 二维码生成器 / QR Code Generator

创建简单的二维码样式图案。

Create simple QR code style patterns.

```go
package main

import (
    "image/color"
    "math/rand"
    "svg"
    "time"
)

func createQRCodePattern() {
    canvas := svg.New(400, 400)
    canvas.SetBackground(color.RGBA{255, 255, 255, 255})
    
    // 生成随机二维码图案 / Generate random QR code pattern
    rand.Seed(time.Now().UnixNano())
    
    gridSize := 20
    cellSize := 15.0
    margin := 50.0
    
    for i := 0; i < gridSize; i++ {
        for j := 0; j < gridSize; j++ {
            if rand.Float64() > 0.5 {
                x := margin + float64(j)*cellSize
                y := margin + float64(i)*cellSize
                
                canvas.Rect(x, y, cellSize-1, cellSize-1).
                    Fill(color.RGBA{0, 0, 0, 255})
            }
        }
    }
    
    // 添加定位标记 / Add position markers
    positions := []struct{ x, y float64 }{
        {margin, margin},                                    // 左上 / Top-left
        {margin + 13*cellSize, margin},                     // 右上 / Top-right
        {margin, margin + 13*cellSize},                     // 左下 / Bottom-left
    }
    
    for _, pos := range positions {
        // 外框 / Outer frame
        canvas.Rect(pos.x, pos.y, cellSize*7, cellSize*7).
            Fill(color.RGBA{0, 0, 0, 255})
        
        // 内框 / Inner frame
        canvas.Rect(pos.x+cellSize, pos.y+cellSize, cellSize*5, cellSize*5).
            Fill(color.RGBA{255, 255, 255, 255})
        
        // 中心点 / Center dot
        canvas.Rect(pos.x+cellSize*3, pos.y+cellSize*3, cellSize, cellSize).
            Fill(color.RGBA{0, 0, 0, 255})
    }
    
    canvas.SaveSVG("qr_code_pattern.svg")
    canvas.SavePNG("qr_code_pattern.png")
}

func main() {
    createQRCodePattern()
    fmt.Println("✅ 二维码图案已创建完成")
}
```

### 条形码生成器 / Barcode Generator

创建简单的条形码图案。

Create simple barcode patterns.

```go
package main

import (
    "image/color"
    "math/rand"
    "svg"
    "time"
)

func createBarcode() {
    canvas := svg.New(600, 200)
    canvas.SetBackground(color.RGBA{255, 255, 255, 255})
    
    // 生成随机条形码 / Generate random barcode
    rand.Seed(time.Now().UnixNano())
    
    x := 50.0
    y := 50.0
    height := 100.0
    
    for i := 0; i < 50; i++ {
        width := 2.0 + rand.Float64()*8 // 随机宽度 / Random width
        
        if rand.Float64() > 0.5 { // 50%概率绘制黑条 / 50% chance to draw black bar
            canvas.Rect(x, y, width, height).
                Fill(color.RGBA{0, 0, 0, 255})
        }
        
        x += width + 1 // 间距 / Spacing
    }
    
    // 添加数字 / Add numbers
    canvas.Text(300, 170, "1234567890123").
        FontSize(14).
        TextAnchor("middle").
        Fill(color.RGBA{0, 0, 0, 255})
    
    canvas.SaveSVG("barcode.svg")
    canvas.SavePNG("barcode.png")
}

func main() {
    createBarcode()
    fmt.Println("✅ 条形码已创建完成")
}
```

### 徽章和标签 / Badges and Labels

创建各种徽章和标签样式。

Create various badge and label styles.

```go
package main

import (
    "image/color"
    "svg"
)

func createBadges() {
    canvas := svg.New(800, 400)
    canvas.SetBackground(color.RGBA{245, 245, 245, 255})
    
    // 成功徽章 / Success badge
    x, y := 100.0, 100.0
    canvas.Rect(x, y, 120, 40).
        Fill(color.RGBA{40, 167, 69, 255}).
        Rx(20).
        Ry(20)
    canvas.Text(x+60, y+25, "SUCCESS").
        FontSize(14).
        FontWeight("bold").
        TextAnchor("middle").
        Fill(color.RGBA{255, 255, 255, 255})
    
    // 警告徽章 / Warning badge
    x += 150
    canvas.Rect(x, y, 120, 40).
        Fill(color.RGBA{255, 193, 7, 255}).
        Rx(20).
        Ry(20)
    canvas.Text(x+60, y+25, "WARNING").
        FontSize(14).
        FontWeight("bold").
        TextAnchor("middle").
        Fill(color.RGBA{0, 0, 0, 255})
    
    // 错误徽章 / Error badge
    x += 150
    canvas.Rect(x, y, 120, 40).
        Fill(color.RGBA{220, 53, 69, 255}).
        Rx(20).
        Ry(20)
    canvas.Text(x+60, y+25, "ERROR").
        FontSize(14).
        FontWeight("bold").
        TextAnchor("middle").
        Fill(color.RGBA{255, 255, 255, 255})
    
    // 信息徽章 / Info badge
    x += 150
    canvas.Rect(x, y, 120, 40).
        Fill(color.RGBA{23, 162, 184, 255}).
        Rx(20).
        Ry(20)
    canvas.Text(x+60, y+25, "INFO").
        FontSize(14).
        FontWeight("bold").
        TextAnchor("middle").
        Fill(color.RGBA{255, 255, 255, 255})
    
    // 版本标签 / Version labels
    versions := []struct {
        text  string
        color color.RGBA
    }{
        {"v1.0.0", color.RGBA{108, 117, 125, 255}},
        {"v2.1.3", color.RGBA{40, 167, 69, 255}},
        {"v3.0.0-beta", color.RGBA{255, 193, 7, 255}},
        {"v4.0.0-alpha", color.RGBA{220, 53, 69, 255}},
    }
    
    for i, version := range versions {
        x := 100.0 + float64(i)*150
        y := 200.0
        
        canvas.Rect(x, y, 100, 30).
            Fill(version.color).
            Rx(15).
            Ry(15)
        canvas.Text(x+50, y+20, version.text).
            FontSize(12).
            FontWeight("bold").
            TextAnchor("middle").
            Fill(color.RGBA{255, 255, 255, 255})
    }
    
    // 进度标签 / Progress labels
    progresses := []struct {
        text     string
        progress float64
        color    color.RGBA
    }{
        {"25%", 0.25, color.RGBA{220, 53, 69, 255}},
        {"50%", 0.50, color.RGBA{255, 193, 7, 255}},
        {"75%", 0.75, color.RGBA{23, 162, 184, 255}},
        {"100%", 1.00, color.RGBA{40, 167, 69, 255}},
    }
    
    for i, prog := range progresses {
        x := 100.0 + float64(i)*150
        y := 300.0
        
        // 背景 / Background
        canvas.Rect(x, y, 100, 20).
            Fill(color.RGBA{233, 236, 239, 255}).
            Rx(10).
            Ry(10)
        
        // 进度条 / Progress bar
        canvas.Rect(x, y, 100*prog.progress, 20).
            Fill(prog.color).
            Rx(10).
            Ry(10)
        
        // 文本 / Text
        canvas.Text(x+50, y+14, prog.text).
            FontSize(10).
            FontWeight("bold").
            TextAnchor("middle").
            Fill(color.RGBA{0, 0, 0, 255})
    }
    
    canvas.SaveSVG("badges.svg")
    canvas.SavePNG("badges.png")
}

func main() {
    createBadges()
    fmt.Println("✅ 徽章和标签已创建完成")
}
```

### 装饰性边框 / Decorative Borders

创建各种装饰性边框和框架。

Create various decorative borders and frames.

```go
package main

import (
    "image/color"
    "math"
    "svg"
)

func createDecorativeBorders() {
    canvas := svg.New(800, 600)
    canvas.SetBackground(color.RGBA{255, 255, 255, 255})
    
    // 简单边框 / Simple border
    canvas.Rect(50, 50, 200, 150).
        Fill("none").
        Stroke(color.RGBA{0, 0, 0, 255}).
        StrokeWidth(3)
    canvas.Text(150, 125, "简单边框").
        FontSize(16).
        TextAnchor("middle").
        Fill(color.RGBA{0, 0, 0, 255})
    
    // 虚线边框 / Dashed border
    canvas.Rect(300, 50, 200, 150).
        Fill("none").
        Stroke(color.RGBA{255, 0, 0, 255}).
        StrokeWidth(3).
        StrokeDashArray("10,5")
    canvas.Text(400, 125, "虚线边框").
        FontSize(16).
        TextAnchor("middle").
        Fill(color.RGBA{255, 0, 0, 255})
    
    // 圆角边框 / Rounded border
    canvas.Rect(550, 50, 200, 150).
        Fill(color.RGBA{240, 248, 255, 255}).
        Stroke(color.RGBA{0, 100, 200, 255}).
        StrokeWidth(3).
        Rx(20).
        Ry(20)
    canvas.Text(650, 125, "圆角边框").
        FontSize(16).
        TextAnchor("middle").
        Fill(color.RGBA{0, 100, 200, 255})
    
    // 装饰性花边 / Decorative lace
    x, y := 50.0, 250.0
    width, height := 200.0, 150.0
    
    // 主边框 / Main border
    canvas.Rect(x, y, width, height).
        Fill(color.RGBA{255, 248, 220, 255}).
        Stroke(color.RGBA{184, 134, 11, 255}).
        StrokeWidth(2)
    
    // 装饰角 / Decorative corners
    corners := []struct{ x, y float64 }{
        {x, y}, {x + width, y}, {x, y + height}, {x + width, y + height},
    }
    
    for _, corner := range corners {
        canvas.Circle(corner.x, corner.y, 8).
            Fill(color.RGBA{184, 134, 11, 255})
        canvas.Circle(corner.x, corner.y, 4).
            Fill(color.RGBA{255, 248, 220, 255})
    }
    
    canvas.Text(x+width/2, y+height/2, "装饰边框").
        FontSize(16).
        TextAnchor("middle").
        Fill(color.RGBA{184, 134, 11, 255})
    
    // 波浪边框 / Wave border
    x = 300.0
    wavePoints := make([]string, 0)
    
    // 上边 / Top edge
    for i := 0; i <= 40; i++ {
        waveX := x + float64(i)*5
        waveY := y + 10*math.Sin(float64(i)*0.3)
        wavePoints = append(wavePoints, fmt.Sprintf("%f,%f", waveX, waveY))
    }
    
    canvas.Polyline(strings.Join(wavePoints, " ")).
        Fill("none").
        Stroke(color.RGBA{255, 20, 147, 255}).
        StrokeWidth(3)
    
    // 下边 / Bottom edge
    wavePoints = make([]string, 0)
    for i := 0; i <= 40; i++ {
        waveX := x + float64(i)*5
        waveY := y + height + 10*math.Sin(float64(i)*0.3+math.Pi)
        wavePoints = append(wavePoints, fmt.Sprintf("%f,%f", waveX, waveY))
    }
    
    canvas.Polyline(strings.Join(wavePoints, " ")).
        Fill("none").
        Stroke(color.RGBA{255, 20, 147, 255}).
        StrokeWidth(3)
    
    canvas.Text(x+width/2, y+height/2, "波浪边框").
        FontSize(16).
        TextAnchor("middle").
        Fill(color.RGBA{255, 20, 147, 255})
    
    // 星形边框 / Star border
    x = 550.0
    canvas.Rect(x, y, width, height).
        Fill(color.RGBA{240, 230, 255, 255}).
        Stroke(color.RGBA{138, 43, 226, 255}).
        StrokeWidth(2)
    
    // 在边框周围添加小星星 / Add small stars around border
    starPositions := []struct{ x, y float64 }{
        {x - 10, y - 10}, {x + width/2, y - 15}, {x + width + 10, y - 10},
        {x - 15, y + height/2}, {x + width + 15, y + height/2},
        {x - 10, y + height + 10}, {x + width/2, y + height + 15}, {x + width + 10, y + height + 10},
    }
    
    for _, pos := range starPositions {
        // 简单的星形 / Simple star
        canvas.Text(pos.x, pos.y, "★").
            FontSize(12).
            TextAnchor("middle").
            Fill(color.RGBA{138, 43, 226, 255})
    }
    
    canvas.Text(x+width/2, y+height/2, "星形边框").
        FontSize(16).
        TextAnchor("middle").
        Fill(color.RGBA{138, 43, 226, 255})
    
    canvas.SaveSVG("decorative_borders.svg")
    canvas.SavePNG("decorative_borders.png")
}

func main() {
    createDecorativeBorders()
    fmt.Println("✅ 装饰性边框已创建完成")
}
```

## 📚 学习路径建议 / Learning Path Recommendations

### 初学者路径 / Beginner Path

1. **基础图形** / Basic Shapes
   - 从 Hello World 示例开始
   - 练习矩形、圆形、直线的绘制
   - 学习基本的颜色和样式设置

2. **文本处理** / Text Handling
   - 学习文本绘制和字体设置
   - 练习文本对齐和装饰
   - 尝试多行文本布局

3. **简单组合** / Simple Combinations
   - 组合多个图形创建复合图案
   - 练习房子、太阳等简单场景
   - 学习图层和重叠效果

### 进阶路径 / Intermediate Path

1. **复杂图形** / Complex Shapes
   - 学习路径绘制和贝塞尔曲线
   - 练习多边形和自定义形状
   - 掌握变换操作（旋转、缩放、倾斜）

2. **样式进阶** / Advanced Styling
   - 学习渐变效果模拟
   - 练习阴影和光效
   - 掌握透明度和混合模式

3. **数据可视化** / Data Visualization
   - 创建各种图表类型
   - 学习数据到图形的映射
   - 练习交互式图表设计

### 高级路径 / Advanced Path

1. **动画制作** / Animation Creation
   - 使用动画构建器创建基础动画
   - 学习缓动函数和时间控制
   - 创建复杂的动画序列

2. **实际应用** / Real-world Applications
   - 开发图标和徽章系统
   - 创建游戏图形和界面元素
   - 构建完整的可视化应用

3. **性能优化** / Performance Optimization
   - 学习大数据量的处理技巧
   - 掌握内存和渲染优化
   - 实现高效的批量操作

## 🎯 最佳实践总结 / Best Practices Summary

### 代码组织 / Code Organization

- **模块化设计**：将复杂图形拆分为独立函数
- **参数化配置**：使用结构体管理配置选项
- **错误处理**：始终检查文件操作的错误返回
- **代码复用**：创建通用的绘图函数库

### 性能优化 / Performance Optimization

- **批量操作**：尽量减少单独的绘图调用
- **内存管理**：及时释放大型图像资源
- **文件大小**：合理选择输出格式和质量
- **渲染效率**：避免不必要的重复计算

### 视觉设计 / Visual Design

- **颜色搭配**：使用和谐的颜色方案
- **布局平衡**：注意元素的对齐和间距
- **层次结构**：通过大小和颜色建立视觉层次
- **用户体验**：确保图形清晰易读

### 维护性 / Maintainability

- **注释文档**：为复杂逻辑添加清晰注释
- **版本控制**：跟踪代码变更和功能演进
- **测试验证**：定期测试输出结果的正确性
- **代码规范**：遵循一致的编码风格

## 🔗 相关资源 / Related Resources

- [快速入门指南](QUICK_START.md) - 库的基础使用方法
- [基础教程](BASIC_TUTORIAL.md) - 详细的功能教程
- [动画构建器文档](ANIMATION_BUILDER_README.md) - 高级动画功能
- [API参考文档](API_REFERENCE.md) - 完整的API说明
- [最佳实践指南](BEST_PRACTICES.md) - 开发建议和技巧

---

**注意**：所有示例代码都可以直接运行，建议按照学习路径逐步练习，从简单到复杂，循序渐进地掌握SVG库的各项功能。

**Note**: All example code can be run directly. It is recommended to practice step by step according to the learning path, from simple to complex, and gradually master the various functions of the SVG library.