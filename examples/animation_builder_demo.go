// 高级动画构建器演示程序 / Advanced Animation Builder Demo
package main

import (
	"fmt"
	"image/color"
	"log"
	"os"

	"github.com/hoonfeng/svg"
)

func main() {
	fmt.Println("🎬 开始演示高级动画构建器...")
	fmt.Println("调试: 程序开始执行")

	// 创建输出目录 / Create output directory
	outputDir := "./output/animation_builder"
	fmt.Printf("调试: 创建输出目录: %s\n", outputDir)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Printf("创建输出目录失败: %v", err)
		return
	}
	fmt.Println("调试: 输出目录创建成功")

	// 演示1: 旋转图形动画 / Demo 1: Rotating shapes animation
	demoRotatingShapes(outputDir)

	// 演示2: 彩色粒子动画 / Demo 2: Colorful particles animation
	demoColorfulParticles(outputDir)

	// 演示3: 脉冲圆形动画 / Demo 3: Pulsing circles animation
	demoPulsingCircles(outputDir)

	// 演示4: 波浪动画 / Demo 4: Wave animation
	demoWaveAnimation(outputDir)

	fmt.Println("✅ 所有动画演示完成！")
}

// 演示旋转图形动画 / Demo rotating shapes animation
func demoRotatingShapes(outputDir string) {
	fmt.Println("🔄 生成旋转图形动画...")
	fmt.Println("调试: 开始创建动画构建器")

	// 创建动画构建器 / Create animation builder
	builder := svg.NewAnimationBuilder(400, 400)
	fmt.Println("调试: 动画构建器创建成功")
	builder.SetFrameCount(60).SetFrameRate(30)
	fmt.Println("调试: 设置帧数和帧率完成")

	// 配置动画 / Configure animation
	config := svg.AnimationConfig{
		Duration:   2.0, // 2秒 / 2 seconds
		Easing:     svg.EaseInOut,
		Background: color.RGBA{20, 20, 40, 255}, // 深蓝背景 / Dark blue background
		Loop:       true,
	}

	// 创建动画并保存 / Create animation and save
	filename := fmt.Sprintf("%s/rotating_shapes.gif", outputDir)
	err := builder.CreateRotatingShapes(config).SaveToGIF(filename)
	if err != nil {
		log.Printf("❌ 保存旋转图形动画失败: %v", err)
		return
	}

	fmt.Printf("✅ 旋转图形动画已保存: %s\n", filename)
	fmt.Printf("   帧数: %d, 时长: %.1f秒\n", builder.GetFrameCount(), builder.GetDuration())
}

// 演示彩色粒子动画 / Demo colorful particles animation
func demoColorfulParticles(outputDir string) {
	fmt.Println("✨ 生成彩色粒子动画...")

	// 创建动画构建器 / Create animation builder
	builder := svg.NewAnimationBuilder(600, 400)
	builder.SetFrameCount(90).SetFrameRate(30)

	// 配置动画 / Configure animation
	config := svg.AnimationConfig{
		Duration:   3.0, // 3秒 / 3 seconds
		Easing:     svg.Linear,
		Background: color.RGBA{10, 10, 20, 255}, // 深色背景 / Dark background
		Loop:       true,
	}

	// 创建动画并保存 / Create animation and save
	filename := fmt.Sprintf("%s/colorful_particles.gif", outputDir)
	err := builder.CreateColorfulParticles(config).SaveToGIF(filename)
	if err != nil {
		log.Printf("❌ 保存彩色粒子动画失败: %v", err)
		return
	}

	fmt.Printf("✅ 彩色粒子动画已保存: %s\n", filename)
	fmt.Printf("   帧数: %d, 时长: %.1f秒\n", builder.GetFrameCount(), builder.GetDuration())
}

// 演示脉冲圆形动画 / Demo pulsing circles animation
func demoPulsingCircles(outputDir string) {
	fmt.Println("💓 生成脉冲圆形动画...")

	// 创建动画构建器 / Create animation builder
	builder := svg.NewAnimationBuilder(500, 500)
	builder.SetFrameCount(80).SetFrameRate(25)

	// 配置动画 / Configure animation
	config := svg.AnimationConfig{
		Duration:   3.2, // 3.2秒 / 3.2 seconds
		Easing:     svg.EaseInOutQuad,
		Background: color.RGBA{5, 5, 15, 255}, // 很深的背景 / Very dark background
		Loop:       true,
	}

	// 创建动画并保存 / Create animation and save
	filename := fmt.Sprintf("%s/pulsing_circles.gif", outputDir)
	err := builder.CreatePulsingCircles(config).SaveToGIF(filename)
	if err != nil {
		log.Printf("❌ 保存脉冲圆形动画失败: %v", err)
		return
	}

	fmt.Printf("✅ 脉冲圆形动画已保存: %s\n", filename)
	fmt.Printf("   帧数: %d, 时长: %.1f秒\n", builder.GetFrameCount(), builder.GetDuration())
}

// 演示波浪动画 / Demo wave animation
func demoWaveAnimation(outputDir string) {
	fmt.Println("🌊 生成波浪动画...")

	// 创建动画构建器 / Create animation builder
	builder := svg.NewAnimationBuilder(800, 300)
	builder.SetFrameCount(60).SetFrameRate(30)

	// 配置动画 / Configure animation
	config := svg.AnimationConfig{
		Duration:   2.0, // 2秒 / 2 seconds
		Easing:     svg.EaseInOut,
		Background: color.RGBA{30, 50, 80, 255}, // 海洋蓝背景 / Ocean blue background
		Loop:       true,
	}

	// 创建动画并保存 / Create animation and save
	filename := fmt.Sprintf("%s/wave_animation.gif", outputDir)
	err := builder.CreateWaveAnimation(config).SaveToGIF(filename)
	if err != nil {
		log.Printf("❌ 保存波浪动画失败: %v", err)
		return
	}

	fmt.Printf("✅ 波浪动画已保存: %s\n", filename)
	fmt.Printf("   帧数: %d, 时长: %.1f秒\n", builder.GetFrameCount(), builder.GetDuration())
}
