package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"

	"github.com/hoonfeng/svg/font" // 导入字体包
)

func main() {
	fmt.Println("字体样式演示程序 / Font Styles Demo")

	// 创建输出目录
	outputDir := "output"
	os.MkdirAll(outputDir, 0755)

	// 创建文本渲染器
	renderer := font.NewSVGTextRenderer()

	// 演示1: 基本字体样式
	createBasicStylesDemo(renderer, outputDir)

	// 演示2: 字体权重对比
	createWeightComparisonDemo(renderer, outputDir)

	// 演示3: 斜体样式对比
	createItalicComparisonDemo(renderer, outputDir)

	// 演示4: 混合样式展示
	createMixedStylesDemo(renderer, outputDir)

	fmt.Println("演示完成！请查看 output 目录中的图像文件。")
}

// createBasicStylesDemo 创建基本字体样式演示
func createBasicStylesDemo(renderer *font.SVGTextRenderer, outputDir string) {
	img := image.NewRGBA(image.Rect(0, 0, 1000, 400))
	fillBackground(img, color.RGBA{255, 255, 255, 255})

	// 标题
	titleStyle := &font.TextStyle{
		FontFamily: "sans-serif",
		FontSize:   32,
		FontWeight: font.FontWeightBold,
		FontStyle:  font.FontStyleNormal,
		Fill:       &image.Uniform{color.RGBA{0, 0, 0, 255}},
	}
	renderer.RenderText(img, "基本字体样式演示 / Basic Font Styles Demo", 50, 50, titleStyle)

	// 基本样式演示
	styles := []struct {
		name   string
		weight font.FontWeight
		style  font.FontStyle
		y      float64
	}{
		{"普通 / Normal", font.FontWeightNormal, font.FontStyleNormal, 120},
		{"粗体 / Bold", font.FontWeightBold, font.FontStyleNormal, 160},
		{"斜体 / Italic", font.FontWeightNormal, font.FontStyleItalic, 200},
		{"倾斜 / Oblique", font.FontWeightNormal, font.FontStyleOblique, 240},
		{"粗斜体 / Bold Italic", font.FontWeightBold, font.FontStyleItalic, 280},
		{"粗倾斜 / Bold Oblique", font.FontWeightBold, font.FontStyleOblique, 320},
	}

	for _, s := range styles {
		// 标签
		labelStyle := &font.TextStyle{
			FontFamily: "sans-serif",
			FontSize:   16,
			FontWeight: font.FontWeightNormal,
			FontStyle:  font.FontStyleNormal,
			Fill:       &image.Uniform{color.RGBA{128, 128, 128, 255}},
		}
		renderer.RenderText(img, s.name+":", 50, s.y-5, labelStyle)

		// 示例文本
		textStyle := &font.TextStyle{
			FontFamily: "sans-serif",
			FontSize:   24,
			FontWeight: s.weight,
			FontStyle:  s.style,
			Fill:       &image.Uniform{color.RGBA{0, 0, 0, 255}},
		}
		renderer.RenderText(img, "The quick brown fox jumps over the lazy dog. 快速的棕色狐狸跳过懒狗。", 250, s.y, textStyle)
	}

	saveImage(img, filepath.Join(outputDir, "basic_styles.png"))
}

// createWeightComparisonDemo 创建字体权重对比演示
func createWeightComparisonDemo(renderer *font.SVGTextRenderer, outputDir string) {
	img := image.NewRGBA(image.Rect(0, 0, 1000, 600))
	fillBackground(img, color.RGBA{255, 255, 255, 255})

	// 标题
	titleStyle := &font.TextStyle{
		FontFamily: "sans-serif",
		FontSize:   32,
		FontWeight: font.FontWeightBold,
		FontStyle:  font.FontStyleNormal,
		Fill:       &image.Uniform{color.RGBA{0, 0, 0, 255}},
	}
	renderer.RenderText(img, "字体权重对比 / Font Weight Comparison", 50, 50, titleStyle)

	// 权重演示
	weights := []struct {
		name   string
		weight font.FontWeight
		y      float64
	}{
		{"100 - Thin", font.FontWeight100, 120},
		{"200 - Extra Light", font.FontWeight200, 160},
		{"300 - Light", font.FontWeight300, 200},
		{"400 - Normal", font.FontWeight400, 240},
		{"500 - Medium", font.FontWeight500, 280},
		{"600 - Semibold", font.FontWeight600, 320},
		{"700 - Bold", font.FontWeight700, 360},
		{"800 - Extra Bold", font.FontWeight800, 400},
		{"900 - Black", font.FontWeight900, 440},
	}

	for _, w := range weights {
		// 标签
		labelStyle := &font.TextStyle{
			FontFamily: "sans-serif",
			FontSize:   16,
			FontWeight: font.FontWeightNormal,
			FontStyle:  font.FontStyleNormal,
			Fill:       &image.Uniform{color.RGBA{128, 128, 128, 255}},
		}
		renderer.RenderText(img, w.name+":", 50, w.y-5, labelStyle)

		// 示例文本
		textStyle := &font.TextStyle{
			FontFamily: "sans-serif",
			FontSize:   24,
			FontWeight: w.weight,
			FontStyle:  font.FontStyleNormal,
			Fill:       &image.Uniform{color.RGBA{0, 0, 0, 255}},
		}
		renderer.RenderText(img, "Font Weight Example 字体权重示例", 250, w.y, textStyle)
	}

	saveImage(img, filepath.Join(outputDir, "weight_comparison.png"))
}

// createItalicComparisonDemo 创建斜体样式对比演示
func createItalicComparisonDemo(renderer *font.SVGTextRenderer, outputDir string) {
	img := image.NewRGBA(image.Rect(0, 0, 1000, 300))
	fillBackground(img, color.RGBA{255, 255, 255, 255})

	// 标题
	titleStyle := &font.TextStyle{
		FontFamily: "sans-serif",
		FontSize:   32,
		FontWeight: font.FontWeightBold,
		FontStyle:  font.FontStyleNormal,
		Fill:       &image.Uniform{color.RGBA{0, 0, 0, 255}},
	}
	renderer.RenderText(img, "斜体样式对比 / Italic Styles Comparison", 50, 50, titleStyle)

	// 斜体样式演示
	styles := []struct {
		name  string
		style font.FontStyle
		y     float64
	}{
		{"Normal - 正常", font.FontStyleNormal, 120},
		{"Italic - 斜体 (15°)", font.FontStyleItalic, 160},
		{"Oblique - 倾斜 (12°)", font.FontStyleOblique, 200},
	}

	for _, s := range styles {
		// 标签
		labelStyle := &font.TextStyle{
			FontFamily: "sans-serif",
			FontSize:   16,
			FontWeight: font.FontWeightNormal,
			FontStyle:  font.FontStyleNormal,
			Fill:       &image.Uniform{color.RGBA{128, 128, 128, 255}},
		}
		renderer.RenderText(img, s.name+":", 50, s.y-5, labelStyle)

		// 示例文本
		textStyle := &font.TextStyle{
			FontFamily: "sans-serif",
			FontSize:   28,
			FontWeight: font.FontWeightNormal,
			FontStyle:  s.style,
			Fill:       &image.Uniform{color.RGBA{0, 0, 0, 255}},
		}
		renderer.RenderText(img, "Italic and Oblique styles comparison 斜体和倾斜样式对比", 300, s.y, textStyle)
	}

	saveImage(img, filepath.Join(outputDir, "italic_comparison.png"))
}

// createMixedStylesDemo 创建混合样式演示
func createMixedStylesDemo(renderer *font.SVGTextRenderer, outputDir string) {
	img := image.NewRGBA(image.Rect(0, 0, 1200, 500))
	fillBackground(img, color.RGBA{250, 250, 250, 255})

	// 标题
	titleStyle := &font.TextStyle{
		FontFamily: "sans-serif",
		FontSize:   36,
		FontWeight: font.FontWeightBlack,
		FontStyle:  font.FontStyleNormal,
		Fill:       &image.Uniform{color.RGBA{0, 0, 0, 255}},
	}
	renderer.RenderText(img, "混合字体样式展示 / Mixed Font Styles Showcase", 50, 60, titleStyle)

	// 创建一个文本布局示例
	texts := []struct {
		content string
		x, y    float64
		size    float64
		weight  font.FontWeight
		style   font.FontStyle
		color   color.RGBA
	}{
		{"这是一个标题", 50, 130, 32, font.FontWeightBold, font.FontStyleNormal, color.RGBA{0, 0, 0, 255}},
		{"This is a subtitle in italic", 50, 170, 24, font.FontWeightMedium, font.FontStyleItalic, color.RGBA{64, 64, 64, 255}},
		{"正文内容使用普通字体，便于阅读。", 50, 210, 18, font.FontWeightNormal, font.FontStyleNormal, color.RGBA{0, 0, 0, 255}},
		{"Body text uses normal font for easy reading.", 50, 240, 18, font.FontWeightNormal, font.FontStyleNormal, color.RGBA{0, 0, 0, 255}},
		{"重要信息", 50, 280, 20, font.FontWeightBold, font.FontStyleNormal, color.RGBA{200, 0, 0, 255}},
		{"Important information in bold red", 150, 280, 20, font.FontWeightBold, font.FontStyleNormal, color.RGBA{200, 0, 0, 255}},
		{"引用文本通常使用斜体", 70, 320, 16, font.FontWeightNormal, font.FontStyleItalic, color.RGBA{100, 100, 100, 255}},
		{"Quoted text usually uses italic style", 70, 345, 16, font.FontWeightNormal, font.FontStyleItalic, color.RGBA{100, 100, 100, 255}},
		{"超粗体标题", 50, 390, 28, font.FontWeightBlack, font.FontStyleNormal, color.RGBA{0, 0, 0, 255}},
		{"Ultra Bold Heading", 200, 390, 28, font.FontWeightBlack, font.FontStyleNormal, color.RGBA{0, 0, 0, 255}},
		{"细体文本", 50, 430, 16, font.FontWeightLight, font.FontStyleNormal, color.RGBA{128, 128, 128, 255}},
		{"Light weight text", 150, 430, 16, font.FontWeightLight, font.FontStyleNormal, color.RGBA{128, 128, 128, 255}},
	}

	for _, t := range texts {
		style := &font.TextStyle{
			FontFamily: "sans-serif",
			FontSize:   t.size,
			FontWeight: t.weight,
			FontStyle:  t.style,
			Fill:       &image.Uniform{t.color},
		}
		renderer.RenderText(img, t.content, t.x, t.y, style)
	}

	saveImage(img, filepath.Join(outputDir, "mixed_styles.png"))
}

// fillBackground 填充背景色
func fillBackground(img *image.RGBA, c color.RGBA) {
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			img.Set(x, y, c)
		}
	}
}

// saveImage 保存图像到文件
func saveImage(img image.Image, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("创建文件失败: %v\n", err)
		return
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		fmt.Printf("编码PNG失败: %v\n", err)
		return
	}

	fmt.Printf("已保存: %s\n", filename)
}
