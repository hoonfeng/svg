package svg

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/hoonfeng/svg/animation"
	"github.com/hoonfeng/svg/types"
)

// AnimationBuilder 动画构建器 / Animation builder
type AnimationBuilder struct {
	svg    *SVG
	frames int
	fps    int
}

// NewAnimationBuilder 创建动画构建器 / Create animation builder
func (s *SVG) NewAnimationBuilder() *AnimationBuilder {
	return &AnimationBuilder{
		svg:    s,
		frames: 30,
		fps:    30,
	}
}

// SetFrames 设置帧数 / Set frame count
func (ab *AnimationBuilder) SetFrames(frames int) *AnimationBuilder {
	ab.frames = frames
	return ab
}

// SetFPS 设置帧率 / Set frame rate
func (ab *AnimationBuilder) SetFPS(fps int) *AnimationBuilder {
	ab.fps = fps
	return ab
}

// AnimationConfig 动画配置 / Animation configuration
type AnimationConfig struct {
	Duration   float64                    // 动画持续时间(秒) / Animation duration (seconds)
	Easing     func(float64) float64      // 缓动函数 / Easing function
	OnProgress func(progress float64)     // 进度回调 / Progress callback
	OnComplete func()                     // 完成回调 / Complete callback
}

// 预定义缓动函数 / Predefined easing functions
var (
	// Linear 线性缓动 / Linear easing
	Linear = func(t float64) float64 {
		return t
	}

	// EaseInOut 缓入缓出 / Ease in-out
	EaseInOut = func(t float64) float64 {
		if t < 0.5 {
			return 2 * t * t
		}
		return -1 + (4-2*t)*t
	}

	// EaseIn 缓入 / Ease in
	EaseIn = func(t float64) float64 {
		return t * t
	}

	// EaseOut 缓出 / Ease out
	EaseOut = func(t float64) float64 {
		return t * (2 - t)
	}

	// Bounce 弹跳 / Bounce
	Bounce = func(t float64) float64 {
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

	// Elastic 弹性 / Elastic
	Elastic = func(t float64) float64 {
		if t == 0 || t == 1 {
			return t
		}
		p := 0.3
		s := p / 4
		return math.Pow(2, -10*t) * math.Sin((t-s)*(2*math.Pi)/p) + 1
	}
)

// CreateRotatingShapes 创建旋转图形动画 / Create rotating shapes animation
func (ab *AnimationBuilder) CreateRotatingShapes() error {
	centerX, centerY := ab.svg.GetWidth()/2, ab.svg.GetHeight()/2
	radius := math.Min(centerX, centerY) * 0.3

	// 创建多个旋转的图形 / Create multiple rotating shapes
	shapeCount := 8
	colors := []color.Color{
		color.RGBA{255, 0, 0, 255},   // 红色 / Red
		color.RGBA{0, 255, 0, 255},   // 绿色 / Green
		color.RGBA{0, 0, 255, 255},   // 蓝色 / Blue
		color.RGBA{255, 255, 0, 255}, // 黄色 / Yellow
		color.RGBA{255, 0, 255, 255}, // 洋红 / Magenta
		color.RGBA{0, 255, 255, 255}, // 青色 / Cyan
		color.RGBA{255, 128, 0, 255}, // 橙色 / Orange
		color.RGBA{128, 0, 255, 255}, // 紫色 / Purple
	}

	for frame := 0; frame < ab.frames; frame++ {
		progress := float64(frame) / float64(ab.frames-1)
		angle := progress * 2 * math.Pi

		// 清除画布 / Clear canvas
		ab.svg.Background(color.RGBA{255, 255, 255, 255})

		for i := 0; i < shapeCount; i++ {
			shapeAngle := angle + float64(i)*2*math.Pi/float64(shapeCount)
			x := centerX + radius*math.Cos(shapeAngle)
			y := centerY + radius*math.Sin(shapeAngle)

			// 根据形状索引选择不同的图形 / Choose different shapes based on index
			switch i % 3 {
			case 0:
				ab.svg.Circle(x, y, 20).Fill(colors[i%len(colors)]).End()
			case 1:
				ab.svg.Rect(x-15, y-15, 30, 30).Fill(colors[i%len(colors)]).End()
			case 2:
				ab.svg.Star(x, y, 25, 10, 5).Fill(colors[i%len(colors)]).End()
			}
		}

		// 保存帧 / Save frame
		filename := fmt.Sprintf("output/frame_%03d.svg", frame)
		if err := ab.svg.Save(filename); err != nil {
			return err
		}
	}

	return nil
}

// CreateColorfulParticles 创建彩色粒子动画 / Create colorful particle animation
func (ab *AnimationBuilder) CreateColorfulParticles() error {
	rand.Seed(time.Now().UnixNano())

	// 粒子结构 / Particle structure
	type Particle struct {
		X, Y   float64
		VX, VY float64
		Color  color.Color
		Size   float64
		Life   float64
	}

	// 初始化粒子 / Initialize particles
	particleCount := 50
	particles := make([]Particle, particleCount)
	centerX, centerY := ab.svg.GetWidth()/2, ab.svg.GetHeight()/2

	for i := range particles {
		angle := rand.Float64() * 2 * math.Pi
		speed := rand.Float64()*3 + 1
		particles[i] = Particle{
			X:     centerX,
			Y:     centerY,
			VX:    math.Cos(angle) * speed,
			VY:    math.Sin(angle) * speed,
			Color: color.RGBA{uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(rand.Intn(256)), 255},
			Size:  rand.Float64()*10 + 5,
			Life:  1.0,
		}
	}

	for frame := 0; frame < ab.frames; frame++ {
		// 清除画布 / Clear canvas
		ab.svg.Background(color.RGBA{0, 0, 0, 255})

		// 更新和绘制粒子 / Update and draw particles
		for i := range particles {
			p := &particles[i]

			// 更新位置 / Update position
			p.X += p.VX
			p.Y += p.VY

			// 更新生命值 / Update life
			p.Life -= 1.0 / float64(ab.frames)

			// 重置粒子如果生命值耗尽 / Reset particle if life is exhausted
			if p.Life <= 0 || p.X < 0 || p.X > ab.svg.GetWidth() || p.Y < 0 || p.Y > ab.svg.GetHeight() {
				angle := rand.Float64() * 2 * math.Pi
				speed := rand.Float64()*3 + 1
				p.X = centerX
				p.Y = centerY
				p.VX = math.Cos(angle) * speed
				p.VY = math.Sin(angle) * speed
				p.Life = 1.0
				p.Color = color.RGBA{uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(rand.Intn(256)), 255}
			}

			// 绘制粒子 / Draw particle
			if p.Life > 0 {
				// 根据生命值调整透明度 / Adjust opacity based on life
				r, g, b, _ := p.Color.RGBA()
				alpha := uint8(float64(255) * p.Life)
				particleColor := color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), alpha}

				ab.svg.Circle(p.X, p.Y, p.Size*p.Life).Fill(particleColor).End()
			}
		}

		// 保存帧 / Save frame
		filename := fmt.Sprintf("output/particle_%03d.svg", frame)
		if err := ab.svg.Save(filename); err != nil {
			return err
		}
	}

	return nil
}