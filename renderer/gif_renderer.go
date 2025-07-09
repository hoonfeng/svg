package renderer

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"os"

	"github.com/hoonfeng/svg/types"
)

// GIFRenderer GIF动画渲染器
type GIFRenderer struct {
	imageRenderer *ImageRenderer
}

// NewGIFRenderer 创建新的GIF渲染器
func NewGIFRenderer() *GIFRenderer {
	return &GIFRenderer{
		imageRenderer: NewImageRenderer(),
	}
}

// RenderAnimationFrames 渲染动画帧序列为GIF
func (g *GIFRenderer) RenderAnimationFrames(frames []*types.Document, width, height int, delay int) (*gif.GIF, error) {
	if len(frames) == 0 {
		return nil, fmt.Errorf("没有动画帧可渲染")
	}

	gifAnim := &gif.GIF{
		Image: make([]*image.Paletted, 0, len(frames)),
		Delay: make([]int, 0, len(frames)),
	}

	for i, frame := range frames {
		// 渲染单帧为RGBA图像
		rgbaImg, err := g.imageRenderer.Render(frame, width, height)
		if err != nil {
			return nil, fmt.Errorf("渲染第%d帧失败: %v", i, err)
		}

		// 转换为调色板图像
		palettedImg := convertToPaletted(rgbaImg)

		gifAnim.Image = append(gifAnim.Image, palettedImg)
		gifAnim.Delay = append(gifAnim.Delay, delay) // 延迟时间（1/100秒为单位）
	}

	return gifAnim, nil
}

// SaveGIF 保存GIF动画到文件
func (g *GIFRenderer) SaveGIF(gifAnim *gif.GIF, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("创建文件失败: %v", err)
	}
	defer file.Close()

	return gif.EncodeAll(file, gifAnim)
}

// convertToPaletted 将RGBA图像转换为调色板图像
func convertToPaletted(src *image.RGBA) *image.Paletted {
	bounds := src.Bounds()
	palette := make(color.Palette, 0, 256)
	colorMap := make(map[color.RGBA]uint8)

	// 收集所有颜色
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := src.RGBAAt(x, y)
			if _, exists := colorMap[c]; !exists && len(palette) < 256 {
				palette = append(palette, c)
				colorMap[c] = uint8(len(palette) - 1)
			}
		}
	}

	// 如果颜色太多，使用简化调色板
	if len(palette) == 0 {
		palette = color.Palette{
			color.RGBA{0, 0, 0, 0},         // 透明
			color.RGBA{255, 255, 255, 255}, // 白色
			color.RGBA{0, 0, 0, 255},       // 黑色
		}
	}

	// 创建调色板图像
	palettedImg := image.NewPaletted(bounds, palette)

	// 复制像素数据
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := src.RGBAAt(x, y)
			if idx, exists := colorMap[c]; exists {
				palettedImg.SetColorIndex(x, y, idx)
			} else {
				// 找到最接近的颜色
				closestIdx := palette.Index(c)
				palettedImg.SetColorIndex(x, y, uint8(closestIdx))
			}
		}
	}

	return palettedImg
}

// RenderAnimationToGIF 直接渲染动画序列为GIF文件
func (g *GIFRenderer) RenderAnimationToGIF(frames []*types.Document, width, height int, delay int, filename string) error {
	gifAnim, err := g.RenderAnimationFrames(frames, width, height, delay)
	if err != nil {
		return err
	}

	return g.SaveGIF(gifAnim, filename)
}
