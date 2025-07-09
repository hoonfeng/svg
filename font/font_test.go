package font

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
)

// TestFontStyles 测试各种字体样式的渲染效果
func TestFontStyles(t *testing.T) {
	// 创建文本渲染器
	renderer := NewSVGTextRenderer()
	
	// 创建测试图像
	img := image.NewRGBA(image.Rect(0, 0, 800, 600))
	
	// 填充白色背景
	for y := 0; y < 600; y++ {
		for x := 0; x < 800; x++ {
			img.Set(x, y, color.RGBA{255, 255, 255, 255})
		}
	}
	
	// 测试文本
	testText := "Hello World 你好世界"
	
	// 测试不同的字体样式
	testCases := []struct {
		name       string
		fontWeight FontWeight
		fontStyle  FontStyle
		y          float64
	}{
		{"Normal", FontWeightNormal, FontStyleNormal, 50},
		{"Bold", FontWeightBold, FontStyleNormal, 100},
		{"Italic", FontWeightNormal, FontStyleItalic, 150},
		{"Oblique", FontWeightNormal, FontStyleOblique, 200},
		{"Bold Italic", FontWeightBold, FontStyleItalic, 250},
		{"Bold Oblique", FontWeightBold, FontStyleOblique, 300},
		{"Light", FontWeightLight, FontStyleNormal, 350},
		{"Medium", FontWeightMedium, FontStyleNormal, 400},
		{"Semibold", FontWeightSemibold, FontStyleNormal, 450},
		{"Black", FontWeightBlack, FontStyleNormal, 500},
		{"100", FontWeight100, FontStyleNormal, 550},
	}
	
	for _, tc := range testCases {
		// 创建文本样式
		style := &TextStyle{
			FontFamily: "sans-serif",
			FontSize:   24,
			FontWeight: tc.fontWeight,
			FontStyle:  tc.fontStyle,
			Fill:       &image.Uniform{color.RGBA{0, 0, 0, 255}},
		}
		
		// 渲染标签文本
		labelStyle := &TextStyle{
			FontFamily: "sans-serif",
			FontSize:   16,
			FontWeight: FontWeightNormal,
			FontStyle:  FontStyleNormal,
			Fill:       &image.Uniform{color.RGBA{128, 128, 128, 255}},
		}
		
		err := renderer.RenderText(img, tc.name+":", 20, tc.y-5, labelStyle)
		if err != nil {
			t.Logf("Warning: Failed to render label for %s: %v", tc.name, err)
		}
		
		// 渲染测试文本
		err = renderer.RenderText(img, testText, 150, tc.y, style)
		if err != nil {
			t.Errorf("Failed to render %s text: %v", tc.name, err)
		}
	}
	
	// 保存测试结果图像
	file, err := os.Create("font_styles_test.png")
	if err != nil {
		t.Fatalf("Failed to create output file: %v", err)
	}
	defer file.Close()
	
	err = png.Encode(file, img)
	if err != nil {
		t.Fatalf("Failed to encode PNG: %v", err)
	}
	
	t.Log("Font styles test completed. Check font_styles_test.png for results.")
}

// TestFontWeightMapping 测试字体权重映射
func TestFontWeightMapping(t *testing.T) {
	testCases := []struct {
		input    FontWeight
		expected FontWeight
	}{
		{FontWeight100, FontWeightLight},
		{FontWeight200, FontWeightLight},
		{FontWeight300, FontWeightLight},
		{FontWeight400, FontWeightNormal},
		{FontWeight500, FontWeightMedium},
		{FontWeight600, FontWeightSemibold},
		{FontWeight700, FontWeightBold},
		{FontWeight800, FontWeightBlack},
		{FontWeight900, FontWeightBlack},
		{FontWeightBold, FontWeightBold},
		{FontWeightNormal, FontWeightNormal},
	}
	
	for _, tc := range testCases {
		result := normalizeFontWeight(tc.input)
		if result != tc.expected {
			t.Errorf("normalizeFontWeight(%s) = %s, expected %s", tc.input, result, tc.expected)
		}
	}
}

// TestFontWeightIntensity 测试字体权重强度
func TestFontWeightIntensity(t *testing.T) {
	testCases := []struct {
		weight   FontWeight
		expected float64
	}{
		{FontWeight100, 0.1},
		{FontWeight400, 0.4},
		{FontWeight700, 0.7},
		{FontWeight900, 0.9},
		{FontWeightNormal, 0.4},
		{FontWeightBold, 0.7},
		{FontWeightLight, 0.3},
		{FontWeightBlack, 0.9},
	}
	
	for _, tc := range testCases {
		result := getFontWeightIntensity(tc.weight)
		if result != tc.expected {
			t.Errorf("getFontWeightIntensity(%s) = %f, expected %f", tc.weight, result, tc.expected)
		}
	}
}

// BenchmarkFontRendering 性能测试
func BenchmarkFontRendering(b *testing.B) {
	renderer := NewSVGTextRenderer()
	img := image.NewRGBA(image.Rect(0, 0, 400, 300))
	style := &TextStyle{
		FontFamily: "sans-serif",
		FontSize:   16,
		FontWeight: FontWeightNormal,
		FontStyle:  FontStyleNormal,
		Fill:       &image.Uniform{color.RGBA{0, 0, 0, 255}},
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		renderer.RenderText(img, "Benchmark Test", 10, 50, style)
	}
}

// BenchmarkBoldRendering 粗体渲染性能测试
func BenchmarkBoldRendering(b *testing.B) {
	renderer := NewSVGTextRenderer()
	img := image.NewRGBA(image.Rect(0, 0, 400, 300))
	style := &TextStyle{
		FontFamily: "sans-serif",
		FontSize:   16,
		FontWeight: FontWeightBold,
		FontStyle:  FontStyleNormal,
		Fill:       &image.Uniform{color.RGBA{0, 0, 0, 255}},
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		renderer.RenderText(img, "Bold Benchmark", 10, 50, style)
	}
}

// BenchmarkItalicRendering 斜体渲染性能测试
func BenchmarkItalicRendering(b *testing.B) {
	renderer := NewSVGTextRenderer()
	img := image.NewRGBA(image.Rect(0, 0, 400, 300))
	style := &TextStyle{
		FontFamily: "sans-serif",
		FontSize:   16,
		FontWeight: FontWeightNormal,
		FontStyle:  FontStyleItalic,
		Fill:       &image.Uniform{color.RGBA{0, 0, 0, 255}},
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		renderer.RenderText(img, "Italic Benchmark", 10, 50, style)
	}
}