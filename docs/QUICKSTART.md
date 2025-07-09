# 快速开始指南 / Quick Start Guide

欢迎使用SVG库！本指南将帮助您快速上手，在几分钟内创建您的第一个SVG图形。

Welcome to the SVG library! This guide will help you get started quickly and create your first SVG graphics in just a few minutes.

## 📋 前提条件 / Prerequisites

- Go 1.18 或更高版本 / Go 1.18 or higher
- 基本的Go编程知识 / Basic Go programming knowledge

## 🚀 安装 / Installation

### 1. 初始化Go模块 / Initialize Go Module

```bash
mkdir my-svg-project
cd my-svg-project
go mod init my-svg-project
```

### 2. 安装SVG库 / Install SVG Library

```bash
go get github.com/hoonfeng/svg
```

### 3. 验证安装 / Verify Installation

创建一个简单的测试文件 `main.go`：

Create a simple test file `main.go`:

```go
package main

import (
    "fmt"
    "github.com/hoonfeng/svg"
)

func main() {
    fmt.Println("SVG库安装成功！")
    fmt.Println("SVG library installed successfully!")
}
```

运行测试 / Run test:
```bash
go run main.go
```

## 🎨 第一个SVG图形 / Your First SVG

让我们创建一个简单的SVG图形，包含圆形、矩形和文本。

Let's create a simple SVG graphic with a circle, rectangle, and text.

### 示例代码 / Example Code

```go
package main

import (
    "log"
    "github.com/hoonfeng/svg"
)

func main() {
    // 创建800x600的SVG画布
    // Create an 800x600 SVG canvas
    s := svg.New(800, 600)
    
    // 添加背景矩形
    // Add background rectangle
    background := s.Rect(0, 0, 800, 600)
    background.Fill("#f0f8ff") // 淡蓝色背景 / Light blue background
    
    // 添加标题文本
    // Add title text
    title := s.Text(400, 50, "我的第一个SVG / My First SVG")
    title.FontFamily("Arial")
         .FontSize(32)
         .FontWeight("bold")
         .Fill("#333333")
         .TextAnchor("middle") // 居中对齐 / Center align
    
    // 添加红色圆形
    // Add red circle
    circle := s.Circle(200, 200, 80)
    circle.Fill("#ff6b6b")
          .Stroke("#d63031", 3)
          .Opacity(0.8)
    
    // 添加蓝色矩形
    // Add blue rectangle
    rect := s.Rect(350, 120, 160, 160)
    rect.Fill("#74b9ff")
        .Stroke("#0984e3", 3)
        .Rx(20) // 圆角 / Rounded corners
        .Ry(20)
    
    // 添加绿色椭圆
    // Add green ellipse
    ellipse := s.Ellipse(600, 200, 100, 60)
    ellipse.Fill("#55a3ff")
           .Stroke("#00b894", 3)
    
    // 添加描述文本
    // Add description text
    desc := s.Text(400, 350, "圆形、矩形和椭圆 / Circle, Rectangle and Ellipse")
    desc.FontFamily("Arial")
        .FontSize(18)
        .Fill("#666666")
        .TextAnchor("middle")
    
    // 添加一条装饰线
    // Add a decorative line
    line := s.Line(100, 400, 700, 400)
    line.Stroke("#ddd", 2)
        .StrokeDashArray("10,5")
    
    // 保存SVG文件
    // Save SVG file
    if err := s.Save("my_first_svg.svg"); err != nil {
        log.Fatal("保存SVG失败:", err)
    }
    
    // 渲染为PNG图像
    // Render to PNG image
    if err := s.RenderToPNG("my_first_svg.png", 800, 600); err != nil {
        log.Fatal("渲染PNG失败:", err)
    }
    
    fmt.Println("✅ SVG和PNG文件创建成功！")
    fmt.Println("✅ SVG and PNG files created successfully!")
}
```

### 运行代码 / Run the Code

```bash
go run main.go
```

运行后，您将在当前目录下看到两个文件：
- `my_first_svg.svg` - SVG矢量图形文件
- `my_first_svg.png` - PNG位图文件

After running, you will see two files in the current directory:
- `my_first_svg.svg` - SVG vector graphics file
- `my_first_svg.png` - PNG bitmap file

## 🎬 创建简单动画 / Creating Simple Animation

让我们创建一个旋转的彩色圆形动画。

Let's create a rotating colorful circle animation.

```go
package main

import (
    "fmt"
    "log"
    "time"
    "github.com/hoonfeng/svg"
    "github.com/hoonfeng/svg/animation"
)

func main() {
    // 创建动画构建器
    // Create animation builder
    builder := animation.NewAnimationBuilder(400, 400)
    
    // 创建60帧动画（1秒，60FPS）
    // Create 60 frames animation (1 second, 60FPS)
    for i := 0; i < 60; i++ {
        // 创建新的SVG帧
        // Create new SVG frame
        s := svg.New(400, 400)
        
        // 添加背景
        // Add background
        bg := s.Rect(0, 0, 400, 400)
        bg.Fill("#1a1a2e")
        
        // 计算旋转角度
        // Calculate rotation angle
        angle := float64(i) * 6 // 每帧旋转6度 / 6 degrees per frame
        
        // 创建旋转的圆形
        // Create rotating circle
        circle := s.Circle(200, 200, 50)
        
        // 根据角度变化颜色
        // Change color based on angle
        hue := int(angle) % 360
        color := fmt.Sprintf("hsl(%d, 70%%, 60%%)", hue)
        
        circle.Fill(color)
              .Stroke("white", 2)
              .Transform(fmt.Sprintf("rotate(%f 200 200)", angle))
        
        // 添加中心点
        // Add center point
        center := s.Circle(200, 200, 5)
        center.Fill("white")
        
        // 添加轨迹圆
        // Add orbit circle
        orbit := s.Circle(200, 200, 50)
        orbit.Fill("none")
             .Stroke("rgba(255,255,255,0.3)", 1)
             .StrokeDashArray("5,5")
        
        // 添加帧到动画
        // Add frame to animation
        builder.AddFrame(s, 16*time.Millisecond) // ~60 FPS
    }
    
    // 保存为GIF动画
    // Save as GIF animation
    if err := builder.SaveGIF("rotating_circle.gif"); err != nil {
        log.Fatal("保存GIF失败:", err)
    }
    
    fmt.Println("🎬 动画GIF创建成功！")
    fmt.Println("🎬 Animation GIF created successfully!")
}
```

## 📝 文本和字体示例 / Text and Font Example

展示如何使用不同的字体和样式。

Demonstrate how to use different fonts and styles.

```go
package main

import (
    "log"
    "github.com/hoonfeng/svg"
)

func main() {
    s := svg.New(800, 600)
    
    // 背景
    // Background
    bg := s.Rect(0, 0, 800, 600)
    bg.Fill("#f8f9fa")
    
    // 标题
    // Title
    title := s.Text(400, 60, "字体样式演示 / Font Style Demo")
    title.FontFamily("Arial")
         .FontSize(36)
         .FontWeight("bold")
         .Fill("#2d3436")
         .TextAnchor("middle")
    
    // 不同字体族示例
    // Different font family examples
    fonts := []struct {
        family string
        text   string
        y      float64
    }{
        {"Arial", "Arial字体 - 现代无衬线字体 / Arial Font - Modern Sans-serif", 120},
        {"Times New Roman", "Times字体 - 经典衬线字体 / Times Font - Classic Serif", 160},
        {"Courier New", "Courier字体 - 等宽字体 / Courier Font - Monospace", 200},
        {"Microsoft YaHei", "微软雅黑 - 中文字体 / Microsoft YaHei - Chinese Font", 240},
    }
    
    for _, font := range fonts {
        text := s.Text(50, font.y, font.text)
        text.FontFamily(font.family)
            .FontSize(20)
            .Fill("#636e72")
    }
    
    // 字体样式示例
    // Font style examples
    styles := []struct {
        weight string
        style  string
        text   string
        y      float64
    }{
        {"normal", "normal", "正常样式 / Normal Style", 320},
        {"bold", "normal", "粗体样式 / Bold Style", 360},
        {"normal", "italic", "斜体样式 / Italic Style", 400},
        {"bold", "italic", "粗斜体样式 / Bold Italic Style", 440},
    }
    
    for _, style := range styles {
        text := s.Text(50, style.y, style.text)
        text.FontFamily("Arial")
            .FontSize(24)
            .FontWeight(style.weight)
            .FontStyle(style.style)
            .Fill("#2d3436")
    }
    
    // 彩色文本示例
    // Colorful text examples
    colors := []string{"#e17055", "#74b9ff", "#55a3ff", "#fd79a8", "#fdcb6e"}
    for i, color := range colors {
        text := s.Text(50+float64(i*140), 520, fmt.Sprintf("彩色%d", i+1))
        text.FontFamily("Arial")
            .FontSize(28)
            .FontWeight("bold")
            .Fill(color)
    }
    
    // 保存文件
    // Save files
    if err := s.Save("font_demo.svg"); err != nil {
        log.Fatal(err)
    }
    
    if err := s.RenderToPNG("font_demo.png", 800, 600); err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("📝 字体演示文件创建成功！")
    fmt.Println("📝 Font demo files created successfully!")
}
```

## 🎨 路径和曲线 / Paths and Curves

学习如何创建复杂的路径和曲线。

Learn how to create complex paths and curves.

```go
package main

import (
    "log"
    "math"
    "github.com/hoonfeng/svg"
)

func main() {
    s := svg.New(800, 600)
    
    // 背景
    bg := s.Rect(0, 0, 800, 600)
    bg.Fill("#2d3436")
    
    // 标题
    title := s.Text(400, 40, "路径和曲线演示 / Paths and Curves Demo")
    title.FontFamily("Arial")
         .FontSize(28)
         .FontWeight("bold")
         .Fill("white")
         .TextAnchor("middle")
    
    // 1. 简单路径 - 三角形
    // Simple path - triangle
    trianglePath := "M 100 100 L 200 100 L 150 50 Z"
    triangle := s.Path(trianglePath)
    triangle.Fill("#74b9ff")
            .Stroke("#0984e3", 2)
    
    // 2. 贝塞尔曲线
    // Bézier curves
    curvePath := "M 300 100 Q 400 50 500 100 T 700 100"
    curve := s.Path(curvePath)
    curve.Fill("none")
         .Stroke("#55a3ff", 3)
    
    // 3. 复杂路径 - 心形
    // Complex path - heart shape
    heartPath := "M 400 200 C 400 180, 380 160, 360 160 C 340 160, 320 180, 320 200 C 320 220, 400 280, 400 280 C 400 280, 480 220, 480 200 C 480 180, 460 160, 440 160 C 420 160, 400 180, 400 200 Z"
    heart := s.Path(heartPath)
    heart.Fill("#fd79a8")
         .Stroke("#e84393", 2)
    
    // 4. 使用路径构建器创建波浪线
    // Create wave using path builder
    builder := svg.NewPathBuilder()
    builder.MoveTo(50, 350)
    
    // 创建正弦波
    // Create sine wave
    for x := 0; x <= 700; x += 10 {
        y := 350 + 50*math.Sin(float64(x)*math.Pi/100)
        builder.LineTo(float64(50+x), y)
    }
    
    wave := s.Path(builder.String())
    wave.Fill("none")
        .Stroke("#00b894", 3)
        .StrokeDashArray("5,5")
    
    // 5. 螺旋线
    // Spiral
    spiralBuilder := svg.NewPathBuilder()
    centerX, centerY := 400.0, 450.0
    spiralBuilder.MoveTo(centerX, centerY)
    
    for i := 0; i < 360*3; i += 5 {
        angle := float64(i) * math.Pi / 180
        radius := float64(i) / 10
        x := centerX + radius*math.Cos(angle)
        y := centerY + radius*math.Sin(angle)
        spiralBuilder.LineTo(x, y)
    }
    
    spiral := s.Path(spiralBuilder.String())
    spiral.Fill("none")
          .Stroke("#fdcb6e", 2)
    
    // 添加说明文本
    // Add description text
    descriptions := []struct {
        text string
        x, y float64
    }{
        {"三角形 / Triangle", 150, 130},
        {"贝塞尔曲线 / Bézier Curve", 500, 130},
        {"心形 / Heart", 400, 320},
        {"正弦波 / Sine Wave", 400, 380},
        {"螺旋线 / Spiral", 400, 550},
    }
    
    for _, desc := range descriptions {
        text := s.Text(desc.x, desc.y, desc.text)
        text.FontFamily("Arial")
            .FontSize(14)
            .Fill("#ddd")
            .TextAnchor("middle")
    }
    
    // 保存文件
    if err := s.Save("paths_demo.svg"); err != nil {
        log.Fatal(err)
    }
    
    if err := s.RenderToPNG("paths_demo.png", 800, 600); err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("🎨 路径演示文件创建成功！")
    fmt.Println("🎨 Paths demo files created successfully!")
}
```

## 🔧 常见问题 / Common Issues

### 1. 字体问题 / Font Issues

**问题**: 中文字体显示不正确
**解决**: 确保系统安装了相应的中文字体

**Issue**: Chinese fonts not displaying correctly
**Solution**: Ensure the system has the appropriate Chinese fonts installed

```go
// 使用系统字体
// Use system fonts
text.FontFamily("Microsoft YaHei, SimHei, sans-serif")
```

### 2. 渲染问题 / Rendering Issues

**问题**: PNG渲染失败
**解决**: 检查输出目录权限和磁盘空间

**Issue**: PNG rendering fails
**Solution**: Check output directory permissions and disk space

```go
// 添加错误处理
// Add error handling
if err := svg.RenderToPNG("output.png", 800, 600); err != nil {
    log.Printf("渲染失败: %v", err)
    // 处理错误...
}
```

### 3. 性能问题 / Performance Issues

**问题**: 大型SVG处理缓慢
**解决**: 使用批量操作和适当的缓存

**Issue**: Large SVG processing is slow
**Solution**: Use batch operations and appropriate caching

```go
// 批量添加元素
// Batch add elements
elements := make([]svg.Element, 0, 1000)
for i := 0; i < 1000; i++ {
    circle := svg.Circle(float64(i%100*8), float64(i/100*8), 2)
    elements = append(elements, circle)
}
svg.AddElements(elements...)
```

## 📚 下一步 / Next Steps

现在您已经掌握了基础知识，可以探索更多高级功能：

Now that you've mastered the basics, you can explore more advanced features:

1. **高级动画** / **Advanced Animation**
   - 关键帧动画 / Keyframe animation
   - 缓动函数 / Easing functions
   - 复杂动画序列 / Complex animation sequences

2. **样式系统** / **Style System**
   - CSS样式 / CSS styles
   - 渐变和图案 / Gradients and patterns
   - 滤镜效果 / Filter effects

3. **交互功能** / **Interactive Features**
   - 事件处理 / Event handling
   - 动态更新 / Dynamic updates
   - 用户交互 / User interaction

4. **性能优化** / **Performance Optimization**
   - 大数据可视化 / Big data visualization
   - 内存管理 / Memory management
   - 渲染优化 / Rendering optimization

## 📖 更多资源 / More Resources

- [完整API文档](API.md) / [Complete API Documentation](API.md)
- [示例代码](../examples/) / [Example Code](../examples/)
- [贡献指南](../CONTRIBUTING.md) / [Contributing Guide](../CONTRIBUTING.md)
- [GitHub仓库](https://github.com/hoonfeng/svg) / [GitHub Repository](https://github.com/hoonfeng/svg)

## 🤝 获取帮助 / Getting Help

如果您遇到问题或有疑问：

If you encounter issues or have questions:

- 📧 [创建GitHub Issue](https://github.com/hoonfeng/svg/issues)
- 💬 [参与讨论](https://github.com/hoonfeng/svg/discussions)
- 📖 查看[API文档](API.md)

---

祝您使用愉快！🎉

Happy coding! 🎉