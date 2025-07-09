package main

import (
	"fmt"
	"image/png"
	"os"
	"time"

	"github.com/hoonfeng/svg/animation"
	"github.com/hoonfeng/svg/elements"
	"github.com/hoonfeng/svg/renderer"
	"github.com/hoonfeng/svg/types"
)

func main() {
	// 创建一个示例SVG文档
	doc := createSampleSVG()

	// 保存SVG文件
	if err := saveSVG(doc, "example.svg"); err != nil {
		fmt.Printf("保存SVG失败: %v\n", err)
		return
	}
	fmt.Println("SVG文件已保存为 example.svg")

	// 渲染为PNG图像
	if err := renderToPNG(doc, "example.png", 800, 600); err != nil {
		fmt.Printf("渲染PNG失败: %v\n", err)
		return
	}
	fmt.Println("PNG图像已保存为 example.png")

	// 演示动画
	animateExample()
}

// createSampleSVG 创建一个示例SVG文档
func createSampleSVG() *types.Document {
	// 创建SVG文档
	doc := types.NewDocument(800, 600)
	doc.SetViewBox(0, 0, 800, 600)

	// 添加背景矩形
	background := elements.NewRect(0, 0, 800, 600)
	background.SetAttribute("fill", "#f0f0f0")
	doc.AppendElement(background)

	// 添加一个圆形
	circle := elements.NewCircle(400, 300, 100)
	circle.SetAttribute("fill", "#ff6600")
	circle.SetAttribute("stroke", "#333333")
	circle.SetAttribute("stroke-width", "5")
	doc.AppendElement(circle)

	// 添加一个矩形
	rect := elements.NewRect(200, 150, 150, 100)
	rect.SetAttribute("fill", "#3366cc")
	rect.SetAttribute("stroke", "#333333")
	rect.SetAttribute("stroke-width", "3")
	rect.SetAttribute("rx", "10")
	rect.SetAttribute("ry", "10")
	doc.AppendElement(rect)

	// 添加一个椭圆
	ellipse := elements.NewEllipse(600, 200, 80, 40)
	ellipse.SetAttribute("fill", "#33cc33")
	ellipse.SetAttribute("stroke", "#333333")
	ellipse.SetAttribute("stroke-width", "2")
	doc.AppendElement(ellipse)

	// 添加一条线
	line := elements.NewLine(100, 400, 700, 450)
	line.SetAttribute("stroke", "#cc3333")
	line.SetAttribute("stroke-width", "4")
	line.SetAttribute("stroke-linecap", "round")
	doc.AppendElement(line)

	// 添加一个多边形
	polygon := elements.NewPolygon([]types.Point{
		{X: 400, Y: 450},
		{X: 450, Y: 500},
		{X: 400, Y: 550},
		{X: 350, Y: 500},
	})
	polygon.SetAttribute("fill", "#9933cc")
	polygon.SetAttribute("stroke", "#333333")
	polygon.SetAttribute("stroke-width", "2")
	doc.AppendElement(polygon)

	// 添加文本
	text := elements.NewText(400, 100, "SVG示例")
	text.SetAttribute("font-family", "Arial, sans-serif")
	text.SetAttribute("font-size", "36")
	text.SetAttribute("font-weight", "bold")
	text.SetAttribute("text-anchor", "middle")
	text.SetAttribute("fill", "#333333")
	doc.AppendElement(text)

	// 添加一个组
	group := elements.NewGroup()
	group.SetAttribute("transform", "translate(150, 450)")

	// 在组中添加一些小圆
	for i := 0; i < 5; i++ {
		smallCircle := elements.NewCircle(float64(i*40), 0, 15)
		smallCircle.SetAttribute("fill", fmt.Sprintf("#%02x%02x%02x", 100+i*30, 150, 200))
		group.AppendChild(smallCircle)
	}

	doc.AppendElement(group)

	// 添加一个路径
	path := elements.NewPath("M100,300 C150,200 250,200 300,300 S450,400 500,300")
	path.SetAttribute("fill", "none")
	path.SetAttribute("stroke", "#996633")
	path.SetAttribute("stroke-width", "4")
	doc.AppendElement(path)

	return doc
}

// saveSVG 保存SVG文档到文件
func saveSVG(doc *types.Document, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return doc.WriteTo(file)
}

// renderToPNG 将SVG渲染为PNG图像
func renderToPNG(doc *types.Document, filename string, width, height int) error {
	// 创建渲染器
	r := renderer.NewImageRenderer()

	// 渲染SVG
	img, err := r.Render(doc, width, height)
	if err != nil {
		return err
	}

	// 保存为PNG
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}

// animateExample 演示动画功能
func animateExample() {
	// 创建SVG文档
	doc := types.NewDocument(800, 600)
	doc.SetViewBox(0, 0, 800, 600)

	// 添加背景
	background := elements.NewRect(0, 0, 800, 600)
	background.SetAttribute("fill", "#f0f0f0")
	doc.AppendElement(background)

	// 添加一个圆形用于动画
	circle := elements.NewCircle(100, 300, 50)
	circle.SetAttribute("fill", "#ff6600")
	circle.SetID("animatedCircle")
	doc.AppendElement(circle)

	// 创建动画管理器
	animManager := animation.NewAnimationManager()

	// 创建一个属性动画，改变圆的x坐标
	positionAnim := animation.NewPropertyAnimation(circle, "cx", "100", "700", 2.0)
	positionAnim.SetEasing(animation.EaseInOutQuad)

	// 创建一个属性动画，改变圆的颜色
	colorAnim := animation.NewPropertyAnimation(circle, "fill", "#ff6600", "#3366cc", 2.0)
	colorAnim.SetEasing(animation.EaseInOutCubic)

	// 创建一个顺序动画组
	seqGroup := animation.NewSequentialAnimationGroup()
	seqGroup.AddAnimation(positionAnim)
	seqGroup.AddAnimation(colorAnim)

	// 设置动画组为循环播放
	seqGroup.SetRepeatCount(-1)
	seqGroup.SetAutoReverse(true)

	// 添加到动画管理器
	animManager.AddAnimation(seqGroup)

	// 开始动画
	animManager.Start()

	// 模拟动画循环
	fmt.Println("动画演示开始，按Ctrl+C停止...")

	// 创建一个渲染器
	r := renderer.NewImageRenderer()

	// 每帧更新和渲染
	frameCount := 0
	for {
		// 更新动画
		animManager.Update()

		// 每10帧保存一次图像
		if frameCount%10 == 0 {
			// 渲染当前帧
			img, err := r.Render(doc, 800, 600)
			if err != nil {
				fmt.Printf("渲染失败: %v\n", err)
				continue
			}

			// 保存为PNG
			filename := fmt.Sprintf("animation_frame_%03d.png", frameCount/10)
			file, err := os.Create(filename)
			if err != nil {
				fmt.Printf("创建文件失败: %v\n", err)
				continue
			}

			if err := png.Encode(file, img); err != nil {
				fmt.Printf("编码PNG失败: %v\n", err)
			}

			file.Close()
			fmt.Printf("保存动画帧: %s\n", filename)
		}

		frameCount++

		// 限制帧率
		time.Sleep(33 * time.Millisecond) // ~30fps

		// 演示60帧后退出
		if frameCount >= 60 {
			break
		}
	}

	fmt.Println("动画演示结束")
}
