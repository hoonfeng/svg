// 高级动画构建器 / Advanced Animation Builder
// 提供简化的API来创建常见的动画效果 / Provides simplified API for creating common animation effects
package svg

import (
	"fmt"
	"image/color"
	"math"
	"os"
	"path/filepath"

	"github.com/hoonfeng/svg/elements"
	"github.com/hoonfeng/svg/renderer"
	"github.com/hoonfeng/svg/types"
)

// AnimationBuilder 动画构建器 / Animation Builder
type AnimationBuilder struct {
	width      int                   // 画布宽度 / Canvas width
	height     int                   // 画布高度 / Canvas height
	frameCount int                   // 帧数 / Frame count
	frameRate  int                   // 帧率(fps) / Frame rate (fps)
	frames     []*types.Document     // 动画帧 / Animation frames
	renderer   *renderer.GIFRenderer // GIF渲染器 / GIF renderer
}

// NewAnimationBuilder 创建新的动画构建器 / Create new animation builder
func NewAnimationBuilder(width, height int) *AnimationBuilder {
	return &AnimationBuilder{
		width:      width,
		height:     height,
		frameCount: 60,
		frameRate:  30,
		frames:     make([]*types.Document, 0),
		renderer:   renderer.NewGIFRenderer(),
	}
}

// SetFrameCount 设置帧数 / Set frame count
func (ab *AnimationBuilder) SetFrameCount(count int) *AnimationBuilder {
	ab.frameCount = count
	return ab
}

// SetFrameRate 设置帧率 / Set frame rate
func (ab *AnimationBuilder) SetFrameRate(fps int) *AnimationBuilder {
	ab.frameRate = fps
	return ab
}

// AnimationConfig 动画配置 / Animation configuration
type AnimationConfig struct {
	Duration   float64    // 动画持续时间(秒) / Animation duration (seconds)
	Easing     EasingFunc // 缓动函数 / Easing function
	Background color.RGBA // 背景颜色 / Background color
	Loop       bool       // 是否循环 / Whether to loop
}

// EasingFunc 缓动函数类型 / Easing function type
type EasingFunc func(t float64) float64

// 预定义缓动函数 / Predefined easing functions
var (
	// Linear 线性缓动 / Linear easing
	Linear EasingFunc = func(t float64) float64 {
		return t
	}

	// EaseInOut 缓入缓出 / Ease in-out
	EaseInOut EasingFunc = func(t float64) float64 {
		if t < 0.5 {
			return 2 * t * t
		}
		return -1 + (4-2*t)*t
	}

	// EaseInQuad 二次方缓入 / Quadratic ease in
	EaseInQuad EasingFunc = func(t float64) float64 {
		return t * t
	}

	// EaseOutQuad 二次方缓出 / Quadratic ease out
	EaseOutQuad EasingFunc = func(t float64) float64 {
		return t * (2 - t)
	}

	// EaseInOutQuad 二次方缓入缓出 / Quadratic ease in-out
	EaseInOutQuad EasingFunc = func(t float64) float64 {
		if t < 0.5 {
			return 2 * t * t
		}
		return -1 + (4-2*t)*t
	}

	// EaseInCubic 三次方缓入 / Cubic ease in
	EaseInCubic EasingFunc = func(t float64) float64 {
		return t * t * t
	}

	// EaseOutCubic 三次方缓出 / Cubic ease out
	EaseOutCubic EasingFunc = func(t float64) float64 {
		t--
		return t*t*t + 1
	}

	// Bounce 弹跳效果 / Bounce effect
	Bounce EasingFunc = func(t float64) float64 {
		if t < 1/2.75 {
			return 7.5625 * t * t
		} else if t < 2/2.75 {
			t -= 1.5 / 2.75
			return 7.5625*t*t + 0.75
		} else if t < 2.5/2.75 {
			t -= 2.25 / 2.75
			return 7.5625*t*t + 0.9375
		} else {
			t -= 2.625 / 2.75
			return 7.5625*t*t + 0.984375
		}
	}
)

// CreateRotatingShapes 创建旋转图形动画 / Create rotating shapes animation
func (ab *AnimationBuilder) CreateRotatingShapes(config AnimationConfig) *AnimationBuilder {
	ab.frames = make([]*types.Document, ab.frameCount)

	for i := 0; i < ab.frameCount; i++ {
		// 计算动画进度 / Calculate animation progress
		progress := float64(i) / float64(ab.frameCount-1)
		easedProgress := config.Easing(progress)
		angle := easedProgress * 2 * math.Pi

		// 创建SVG文档 / Create SVG document
		doc := types.NewDocument(ab.width, ab.height)
		doc.SetViewBox(0, 0, float64(ab.width), float64(ab.height))

		// 添加背景 / Add background
		background := elements.NewRect(0, 0, float64(ab.width), float64(ab.height))
		background.SetAttribute("fill", fmt.Sprintf("#%02x%02x%02x", config.Background.R, config.Background.G, config.Background.B))
		doc.AppendElement(background)

		// 添加旋转的图形 / Add rotating shapes
		centerX := float64(ab.width) / 2
		centerY := float64(ab.height) / 2

		// 旋转的圆形 / Rotating circles
		for j := 0; j < 6; j++ {
			shapeAngle := angle + float64(j)*math.Pi/3
			x := centerX + 100*math.Cos(shapeAngle)
			y := centerY + 100*math.Sin(shapeAngle)
			radius := 20 + 10*math.Sin(angle*2+float64(j))

			// 彩色圆形 / Colorful circles
			hue := int((progress*360 + float64(j)*60)) % 360
			shapeColor := hslToRGBA(hue, 80, 60)

			circle := elements.NewCircle(x, y, radius)
			circle.SetAttribute("fill", fmt.Sprintf("#%02x%02x%02x", shapeColor.R, shapeColor.G, shapeColor.B))
			doc.AppendElement(circle)
		}

		ab.frames[i] = doc
	}

	return ab
}

// CreateColorfulParticles 创建彩色粒子动画 / Create colorful particles animation
func (ab *AnimationBuilder) CreateColorfulParticles(config AnimationConfig) *AnimationBuilder {
	ab.frames = make([]*types.Document, ab.frameCount)

	for i := 0; i < ab.frameCount; i++ {
		// 计算动画进度 / Calculate animation progress
		progress := float64(i) / float64(ab.frameCount-1)
		easedProgress := config.Easing(progress)
		angle := easedProgress * 4 * math.Pi // 更快的旋转 / Faster rotation

		// 创建SVG文档 / Create SVG document
		doc := types.NewDocument(ab.width, ab.height)
		doc.SetViewBox(0, 0, float64(ab.width), float64(ab.height))

		// 动态背景 / Dynamic background
		hue := int(progress * 360)
		bgColor := hslToRGBA(hue, 30, 20)
		background := elements.NewRect(0, 0, float64(ab.width), float64(ab.height))
		background.SetAttribute("fill", fmt.Sprintf("#%02x%02x%02x", bgColor.R, bgColor.G, bgColor.B))
		doc.AppendElement(background)

		// 添加粒子效果 / Add particle effects
		centerX := float64(ab.width) / 2
		centerY := float64(ab.height) / 2

		for p := 0; p < 30; p++ {
			particleAngle := angle + float64(p)*math.Pi/15
			radius := 50 + float64(p)*8 + 30*math.Sin(angle+float64(p))
			x := centerX + radius*math.Cos(particleAngle)
			y := centerY + radius*math.Sin(particleAngle)
			size := 3 + 5*math.Sin(angle*2+float64(p))

			// 彩色粒子 / Colorful particles
			particleHue := (hue + p*12) % 360
			particleColor := hslToRGBA(particleHue, 90, 70)

			circle := elements.NewCircle(x, y, size)
			circle.SetAttribute("fill", fmt.Sprintf("#%02x%02x%02x", particleColor.R, particleColor.G, particleColor.B))
			doc.AppendElement(circle)
		}

		ab.frames[i] = doc
	}

	return ab
}

// CreatePulsingCircles 创建脉冲圆形动画 / Create pulsing circles animation
func (ab *AnimationBuilder) CreatePulsingCircles(config AnimationConfig) *AnimationBuilder {
	ab.frames = make([]*types.Document, ab.frameCount)

	for i := 0; i < ab.frameCount; i++ {
		// 计算动画进度 / Calculate animation progress
		progress := float64(i) / float64(ab.frameCount-1)
		easedProgress := config.Easing(progress)

		// 创建SVG文档 / Create SVG document
		doc := types.NewDocument(ab.width, ab.height)
		doc.SetViewBox(0, 0, float64(ab.width), float64(ab.height))

		// 添加背景 / Add background
		background := elements.NewRect(0, 0, float64(ab.width), float64(ab.height))
		background.SetAttribute("fill", fmt.Sprintf("#%02x%02x%02x", config.Background.R, config.Background.G, config.Background.B))
		doc.AppendElement(background)

		// 添加脉冲圆形 / Add pulsing circles
		centerX := float64(ab.width) / 2
		centerY := float64(ab.height) / 2

		for ring := 0; ring < 5; ring++ {
			// 每个圆环有不同的脉冲相位 / Each ring has different pulse phase
			phase := easedProgress*2*math.Pi + float64(ring)*math.Pi/2
			baseRadius := 30 + float64(ring)*40
			pulse := 0.5 + 0.5*math.Sin(phase)
			radius := baseRadius * (0.5 + 0.5*pulse)

			// 颜色随脉冲变化 / Color changes with pulse
			hue := int((progress*360 + float64(ring)*72)) % 360
			alpha := uint8(100 + 155*pulse) // 透明度变化 / Alpha changes
			ringColor := hslToRGBA(hue, 80, 60)
			ringColor.A = alpha

			circle := elements.NewCircle(centerX, centerY, radius)
			circle.SetAttribute("fill", "none")
			circle.SetAttribute("stroke", fmt.Sprintf("rgba(%d,%d,%d,%.2f)", ringColor.R, ringColor.G, ringColor.B, float64(alpha)/255.0))
			circle.SetAttribute("stroke-width", "3")
			doc.AppendElement(circle)
		}

		ab.frames[i] = doc
	}

	return ab
}

// CreateWaveAnimation 创建波浪动画 / Create wave animation
func (ab *AnimationBuilder) CreateWaveAnimation(config AnimationConfig) *AnimationBuilder {
	ab.frames = make([]*types.Document, ab.frameCount)

	for i := 0; i < ab.frameCount; i++ {
		// 计算动画进度 / Calculate animation progress
		progress := float64(i) / float64(ab.frameCount-1)
		easedProgress := config.Easing(progress)
		phase := easedProgress * 2 * math.Pi

		// 创建SVG文档 / Create SVG document
		doc := types.NewDocument(ab.width, ab.height)
		doc.SetViewBox(0, 0, float64(ab.width), float64(ab.height))

		// 添加背景 / Add background
		background := elements.NewRect(0, 0, float64(ab.width), float64(ab.height))
		background.SetAttribute("fill", fmt.Sprintf("#%02x%02x%02x", config.Background.R, config.Background.G, config.Background.B))
		doc.AppendElement(background)

		// 创建波浪路径 / Create wave path
		pathData := fmt.Sprintf("M0,%d", ab.height/2)
		for x := 0; x <= ab.width; x += 5 {
			y := float64(ab.height)/2 + 50*math.Sin(float64(x)*0.02+phase)
			pathData += fmt.Sprintf(" L%d,%.2f", x, y)
		}
		pathData += fmt.Sprintf(" L%d,%d L0,%d Z", ab.width, ab.height, ab.height)

		// 波浪颜色 / Wave color
		hue := int((progress * 360)) % 360
		waveColor := hslToRGBA(hue, 70, 50)

		path := elements.NewPath(pathData)
		path.SetAttribute("fill", fmt.Sprintf("#%02x%02x%02x", waveColor.R, waveColor.G, waveColor.B))
		path.SetAttribute("opacity", "0.8")
		doc.AppendElement(path)

		ab.frames[i] = doc
	}

	return ab
}

// SaveToGIF 保存为GIF文件 / Save to GIF file
func (ab *AnimationBuilder) SaveToGIF(filename string) error {
	if len(ab.frames) == 0 {
		return fmt.Errorf("没有动画帧可保存")
	}

	// 确保输出目录存在 / Ensure output directory exists
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %v", err)
	}

	// 计算帧延迟 / Calculate frame delay
	frameDelay := 100 / ab.frameRate // 1/100秒为单位 / In 1/100 second units

	// 渲染GIF / Render GIF
	return ab.renderer.RenderAnimationToGIF(ab.frames, ab.width, ab.height, frameDelay, filename)
}

// GetFrameCount 获取帧数 / Get frame count
func (ab *AnimationBuilder) GetFrameCount() int {
	return len(ab.frames)
}

// GetDuration 获取动画时长(秒) / Get animation duration (seconds)
func (ab *AnimationBuilder) GetDuration() float64 {
	return float64(len(ab.frames)) / float64(ab.frameRate)
}

// hslToRGBA HSL转RGBA / HSL to RGBA
func hslToRGBA(h, s, l int) color.RGBA {
	// 确保值在有效范围内 / Ensure values are in valid range
	h = h % 360
	if s > 100 {
		s = 100
	}
	if l > 100 {
		l = 100
	}

	// 转换为RGB / Convert to RGB
	r, g, b := hslToRGB(float64(h), float64(s)/100.0, float64(l)/100.0)

	return color.RGBA{
		R: uint8(r * 255),
		G: uint8(g * 255),
		B: uint8(b * 255),
		A: 255,
	}
}

// hslToRGB HSL转RGB / HSL to RGB
func hslToRGB(h, s, l float64) (float64, float64, float64) {
	h = h / 360.0

	var r, g, b float64

	if s == 0 {
		r, g, b = l, l, l // 灰度 / Grayscale
	} else {
		hue2rgb := func(p, q, t float64) float64 {
			if t < 0 {
				t += 1
			}
			if t > 1 {
				t -= 1
			}
			if t < 1.0/6.0 {
				return p + (q-p)*6*t
			}
			if t < 1.0/2.0 {
				return q
			}
			if t < 2.0/3.0 {
				return p + (q-p)*(2.0/3.0-t)*6
			}
			return p
		}

		var q float64
		if l < 0.5 {
			q = l * (1 + s)
		} else {
			q = l + s - l*s
		}
		p := 2*l - q

		r = hue2rgb(p, q, h+1.0/3.0)
		g = hue2rgb(p, q, h)
		b = hue2rgb(p, q, h-1.0/3.0)
	}

	return r, g, b
}
