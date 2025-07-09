# SVG Library for Go

一个功能强大的Go语言SVG图形库，支持图形绘制、动画制作、文本渲染和高级样式设置。

A powerful SVG graphics library for Go, supporting shape drawing, animation creation, text rendering, and advanced styling.

## ✨ 特性 Features

- 🎨 **丰富的图形绘制** - 支持矩形、圆形、椭圆、直线、折线、多边形、路径等基础图形
- 🎬 **动画支持** - 内置动画构建器，支持旋转、缩放、移动、颜色变化等动画效果
- 📝 **文本渲染** - 完整的文本渲染系统，支持多种字体、样式和对齐方式
- 🎨 **样式系统** - 支持颜色、渐变、描边、填充、透明度等样式设置
- 🔄 **变换支持** - 支持平移、旋转、缩放、倾斜等2D变换
- 📊 **高级API** - 提供链式调用的高级API，简化复杂图形的创建
- 🖼️ **多格式输出** - 支持SVG、PNG、GIF等多种格式输出
- 🚀 **高性能** - 优化的渲染引擎，支持大规模图形处理

## 📦 安装 Installation

```bash
go get github.com/hoonfeng/svg
```

## 🚀 快速开始 Quick Start

### 基础图形绘制

```go
package main

import (
    "image/color"
    "svg"
)

func main() {
    // 创建SVG画布
    canvas := svg.New(400, 300)
    
    // 设置背景
    canvas.Background(color.RGBA{255, 255, 255, 255})
    
    // 绘制红色圆形
    canvas.Circle(200, 150, 50).
        Fill(color.RGBA{255, 0, 0, 255}).
        Stroke(color.RGBA{0, 0, 0, 255}).
        StrokeWidth(2).
        End()
    
    // 绘制蓝色矩形
    canvas.Rect(100, 100, 200, 100).
        Fill(color.RGBA{0, 0, 255, 255}).
        Opacity(0.7).
        End()
    
    // 添加文本
    canvas.Text(200, 50, "Hello SVG!").
        Fill(color.RGBA{0, 0, 0, 255}).
        FontSize(24).
        FontFamily("Arial").
        TextAnchor("middle").
        End()
    
    // 保存为SVG文件
    canvas.Save("example.svg")
    
    // 保存为PNG文件
    canvas.SavePNG("example.png", 400, 300)
}
```

### 动画制作

```go
package main

import (
    "image/color"
    "svg/animation"
)

func main() {
    // 创建动画构建器
    builder := animation.NewAnimationBuilder(400, 300, 3.0) // 3秒动画
    
    // 设置背景
    builder.SetBackground(color.RGBA{240, 240, 240, 255})
    
    // 创建旋转的彩色圆形
    for i := 0; i < 60; i++ {
        frame := float64(i) / 60.0
        
        // 添加帧
        builder.AddFrame(frame)
        
        // 旋转角度
        angle := frame * 360
        
        // 绘制旋转的圆形
        builder.Circle(200, 150, 30).
            Fill(animation.HSL(int(angle), 80, 60)).
            Transform(fmt.Sprintf("rotate(%.1f 200 150)", angle)).
            End()
    }
    
    // 生成GIF动画
    builder.SaveGIF("rotating_circle.gif")
}
```

## 支持的SVG元素

- 矩形 (Rect)
- 圆形 (Circle)
- 椭圆 (Ellipse)
- 线段 (Line)
- 折线 (Polyline)
- 多边形 (Polygon)
- 路径 (Path)
- 文本 (Text)
- 组 (Group)

## 路径命令支持

- 移动 (M, m)
- 直线 (L, l)
- 水平线 (H, h)
- 垂直线 (V, v)
- 三次贝塞尔曲线 (C, c)
- 平滑三次贝塞尔曲线 (S, s)
- 二次贝塞尔曲线 (Q, q)
- 平滑二次贝塞尔曲线 (T, t)
- 椭圆弧 (A, a)
- 闭合路径 (Z, z)

## 📚 文档 Documentation

- [📖 快速入门指南](docs/QUICKSTART.md) - 快速上手SVG库
- [📘 基础教程](docs/BASIC_TUTORIAL.md) - 详细的基础功能教程
- [📙 API参考](docs/API.md) - 完整的API文档
- [📗 示例集合](docs/EXAMPLES.md) - 丰富的代码示例
- [📕 最佳实践](docs/BEST_PRACTICES.md) - 性能优化和代码质量指南
- [🔧 故障排除](docs/TROUBLESHOOTING.md) - 常见问题解决方案
- [📋 文档中心](docs/README.md) - 完整文档导航

## 🏗️ 项目结构 Project Structure

```
svg/
├── animation/          # 动画模块
├── api/               # 高级API
├── attributes/        # 属性和样式
├── cmd/               # 命令行工具和示例
├── elements/          # SVG元素定义
├── examples/          # 示例代码
├── font/              # 字体和文本渲染
├── io/                # 输入输出处理
├── output/            # 输出文件目录
├── parser/            # XML解析器
├── path/              # 路径处理
├── renderer/          # 渲染引擎
├── test_animation/    # 动画测试
├── test_renderer/     # 渲染测试
├── transform/         # 变换处理
├── types/             # 类型定义
├── docs/              # 文档文件
└── svg.go             # 主要API
```

## 🎯 主要功能 Main Features

### 图形绘制
- ✅ 基础图形：矩形、圆形、椭圆、直线、折线、多边形
- ✅ 复杂路径：贝塞尔曲线、弧线、复合路径
- ✅ 分组和嵌套：支持元素分组和层次结构

### 样式系统
- ✅ 颜色支持：RGB、RGBA、HSL、命名颜色
- ✅ 填充和描边：实色、渐变、图案填充
- ✅ 透明度和混合模式
- ✅ 滤镜效果

### 文本处理
- ✅ 多字体支持：系统字体、自定义字体
- ✅ 文本样式：粗体、斜体、下划线
- ✅ 文本对齐：左对齐、居中、右对齐
- ✅ 文本路径：沿路径排列文本

### 动画功能
- ✅ 关键帧动画：位置、旋转、缩放、颜色
- ✅ 缓动函数：线性、贝塞尔、弹性等
- ✅ 循环和延迟控制
- ✅ GIF导出

### 输出格式
- ✅ SVG：矢量格式，可缩放
- ✅ PNG：高质量位图
- ✅ GIF：动画支持
- ✅ 批量处理

## 🔧 系统要求 Requirements

- Go 1.18 或更高版本
- 支持的操作系统：Windows、macOS、Linux
- 可选依赖：
  - 字体文件（用于自定义字体）
  - 图像处理库（用于高级滤镜）

## 🤝 贡献 Contributing

我们欢迎所有形式的贡献！请查看 [CONTRIBUTING.md](CONTRIBUTING.md) 了解详细信息。

### 开发环境设置

```bash
# 克隆仓库
git clone https://github.com/hoonfeng/svg.git
cd svg

# 安装依赖
go mod tidy

# 运行测试
go test ./...

# 运行示例
go run examples/animation_builder_demo.go
```

## 📄 许可证 License

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🙏 致谢 Acknowledgments

- 感谢所有贡献者的努力
- 特别感谢开源社区的支持
- 参考了多个优秀的图形库设计

## 📞 联系我们 Contact

- 🐛 Issues: [GitHub Issues](https://github.com/hoonfeng/svg/issues)
- 💬 Discussions: [GitHub Discussions](https://github.com/hoonfeng/svg/discussions)

## 🔗 相关链接 Related Links

- [SVG规范](https://www.w3.org/TR/SVG2/)
- [Go语言官网](https://golang.org/)
- [图形编程资源](https://github.com/topics/graphics)

---

⭐ 如果这个项目对你有帮助，请给我们一个星标！

⭐ If this project helps you, please give us a star!