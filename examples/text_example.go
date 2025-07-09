package main

import (
	"fmt"
	"image/png"
	"os"

	"github.com/hoonfeng/svg/elements"
	"github.com/hoonfeng/svg/renderer"
	"github.com/hoonfeng/svg/types"
)

func main() {
	// 创建SVG文档
	doc := types.NewDocument(400, 300)
	doc.SetViewBox(0, 0, 400, 300)

	// 创建不同样式的文本元素

	// 1. 基本文本
	text1 := elements.NewText(50, 50, "Hello SVG!")
	text1.SetFontFamily("sans-serif")
	text1.SetFontSize(24)
	text1.SetFill("black")
	doc.AppendElement(text1)

	// 2. 居中对齐的文本
	text2 := elements.NewText(200, 100, "Centered Text")
	text2.SetFontFamily("serif")
	text2.SetFontSize(20)
	text2.SetTextAnchor("middle")
	text2.SetFill("blue")
	doc.AppendElement(text2)

	// 3. 右对齐的文本
	text3 := elements.NewText(350, 150, "Right Aligned")
	text3.SetFontFamily("monospace")
	text3.SetFontSize(18)
	text3.SetTextAnchor("end")
	text3.SetFill("red")
	doc.AppendElement(text3)

	// 4. 粗体文本
	text4 := elements.NewText(50, 200, "Bold Text")
	text4.SetFontFamily("sans-serif")
	text4.SetFontSize(22)
	text4.SetFontWeight("bold")
	text4.SetFill("green")
	doc.AppendElement(text4)

	// 5. 带描边的文本
	text5 := elements.NewText(200, 250, "Stroked Text")
	text5.SetFontFamily("sans-serif")
	text5.SetFontSize(20)
	text5.SetFill("yellow")
	text5.SetStroke("black")
	text5.SetStrokeWidth(1)
	text5.SetTextAnchor("middle")
	doc.AppendElement(text5)

	// 渲染为图像
	renderer := renderer.NewImageRenderer()
	img, err := renderer.Render(doc, 400, 300)
	if err != nil {
		fmt.Printf("渲染错误: %v\n", err)
		return
	}

	// 保存图像
	file, err := os.Create("text_example.png")
	if err != nil {
		fmt.Printf("创建文件错误: %v\n", err)
		return
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		fmt.Printf("编码图像错误: %v\n", err)
		return
	}

	fmt.Println("SVG文本示例已保存为 text_example.png")
	fmt.Println("该示例展示了以下SVG标准文本功能:")
	fmt.Println("- 不同字体族 (sans-serif, serif, monospace)")
	fmt.Println("- 不同字体大小")
	fmt.Println("- 文本锚点 (start, middle, end)")
	fmt.Println("- 字体粗细 (normal, bold)")
	fmt.Println("- 填充和描边颜色")
	fmt.Println("- 描边宽度")
}
