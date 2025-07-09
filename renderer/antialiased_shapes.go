package renderer

import (
	"image"
	"image/color"
	"math"
)

// AntiAliasedRenderer 抗锯齿渲染器 / Anti-aliased renderer
type AntiAliasedRenderer struct {
	*ImageRenderer
}

// NewAntiAliasedRenderer 创建抗锯齿渲染器 / Create anti-aliased renderer
func NewAntiAliasedRenderer() *AntiAliasedRenderer {
	return &AntiAliasedRenderer{
		ImageRenderer: NewImageRenderer(),
	}
}

// DrawAntiAliasedCircle 绘制抗锯齿圆形描边 / Draw anti-aliased circle stroke
func (r *AntiAliasedRenderer) DrawAntiAliasedCircle(img *image.RGBA, centerX, centerY float64, radius float64, c color.RGBA) {
	// 默认描边宽度 / Default stroke width
	strokeWidth := 1.0

	// 计算边界框 / Calculate bounding box
	minX := int(math.Floor(centerX - radius - strokeWidth - 1))
	maxX := int(math.Ceil(centerX + radius + strokeWidth + 1))
	minY := int(math.Floor(centerY - radius - strokeWidth - 1))
	maxY := int(math.Ceil(centerY + radius + strokeWidth + 1))

	// 遍历边界框内的每个像素 / Iterate through each pixel in bounding box
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			// 使用距离场计算描边覆盖率 / Use distance field to calculate stroke coverage
			coverage := calculateCircleStrokeCoverageWithSupersampling(float64(x), float64(y), centerX, centerY, radius, strokeWidth, 4)

			if coverage > 0 {
				// 混合颜色 / Blend color
				blendedColor := blendColors(getPixelColor(img, x, y), c, coverage)
				DrawPixel(img, x, y, blendedColor)
			}
		}
	}
}

// FillAntiAliasedCircle 填充抗锯齿圆形 / Fill anti-aliased circle
func (r *AntiAliasedRenderer) FillAntiAliasedCircle(img *image.RGBA, centerX, centerY float64, radius float64, c color.RGBA) {
	// 计算边界框 / Calculate bounding box
	minX := int(math.Floor(centerX - radius - 1))
	maxX := int(math.Ceil(centerX + radius + 1))
	minY := int(math.Floor(centerY - radius - 1))
	maxY := int(math.Ceil(centerY + radius + 1))

	// 遍历边界框内的每个像素 / Iterate through each pixel in bounding box
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			// 使用超采样计算覆盖率 / Use supersampling to calculate coverage
			coverage := calculateCircleCoverageWithSupersampling(float64(x), float64(y), centerX, centerY, radius, 4)

			if coverage > 0 {
				// 混合颜色 / Blend color
				blendedColor := blendColors(getPixelColor(img, x, y), c, coverage)
				DrawPixel(img, x, y, blendedColor)
			}
		}
	}
}

// DrawAntiAliasedEllipse 绘制抗锯齿椭圆描边 / Draw anti-aliased ellipse stroke
func (r *AntiAliasedRenderer) DrawAntiAliasedEllipse(img *image.RGBA, centerX, centerY, radiusX, radiusY float64, c color.RGBA) {
	// 默认描边宽度 / Default stroke width
	strokeWidth := 1.0

	// 计算边界框 / Calculate bounding box
	maxRadius := math.Max(radiusX, radiusY)
	minX := int(math.Floor(centerX - maxRadius - strokeWidth - 1))
	maxX := int(math.Ceil(centerX + maxRadius + strokeWidth + 1))
	minY := int(math.Floor(centerY - maxRadius - strokeWidth - 1))
	maxY := int(math.Ceil(centerY + maxRadius + strokeWidth + 1))

	// 遍历边界框内的每个像素 / Iterate through each pixel in bounding box
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			// 使用距离场计算椭圆描边覆盖率 / Use distance field to calculate ellipse stroke coverage
			coverage := calculateEllipseStrokeCoverageWithSupersampling(float64(x), float64(y), centerX, centerY, radiusX, radiusY, strokeWidth, 4)

			if coverage > 0 {
				// 混合颜色 / Blend color
				blendedColor := blendColors(getPixelColor(img, x, y), c, coverage)
				DrawPixel(img, x, y, blendedColor)
			}
		}
	}
}

// FillAntiAliasedEllipse 填充抗锯齿椭圆 / Fill anti-aliased ellipse
func (r *AntiAliasedRenderer) FillAntiAliasedEllipse(img *image.RGBA, centerX, centerY, radiusX, radiusY float64, c color.RGBA) {
	// 计算边界框 / Calculate bounding box
	minX := int(math.Floor(centerX - radiusX - 1))
	maxX := int(math.Ceil(centerX + radiusX + 1))
	minY := int(math.Floor(centerY - radiusY - 1))
	maxY := int(math.Ceil(centerY + radiusY + 1))

	// 遍历边界框内的每个像素 / Iterate through each pixel in bounding box
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			// 计算椭圆内部覆盖率 / Calculate ellipse interior coverage
			coverage := calculateEllipseInteriorCoverage(float64(x), float64(y), centerX, centerY, radiusX, radiusY, 4)

			if coverage > 0 {
				// 混合颜色 / Blend color
				blendedColor := blendColors(getPixelColor(img, x, y), c, coverage)
				DrawPixel(img, x, y, blendedColor)
			}
		}
	}
}

// calculateCircleCoverageWithSupersampling 使用超采样计算圆形覆盖率 / Calculate circle coverage with supersampling
func calculateCircleCoverageWithSupersampling(pixelX, pixelY, centerX, centerY, radius float64, samples int) float64 {
	insideCount := 0
	totalSamples := samples * samples
	step := 1.0 / float64(samples)

	for i := 0; i < samples; i++ {
		for j := 0; j < samples; j++ {
			// 计算子像素位置 / Calculate sub-pixel position
			sampleX := pixelX + (float64(i)+0.5)*step
			sampleY := pixelY + (float64(j)+0.5)*step

			// 检查是否在圆内 / Check if inside circle
			dx := sampleX - centerX
			dy := sampleY - centerY
			distance := math.Sqrt(dx*dx + dy*dy)

			if distance <= radius {
				insideCount++
			}
		}
	}

	return float64(insideCount) / float64(totalSamples)
}

// DrawAntiAliasedCircle 绘制抗锯齿圆形（描边模式）
func DrawAntiAliasedCircle(img *image.RGBA, centerX, centerY, radius float64, color color.RGBA, strokeWidth float64) {
	// 计算边界框，考虑描边宽度
	halfStroke := strokeWidth / 2
	minX := int(math.Floor(centerX - radius - halfStroke))
	maxX := int(math.Ceil(centerX + radius + halfStroke))
	minY := int(math.Floor(centerY - radius - halfStroke))
	maxY := int(math.Ceil(centerY + radius + halfStroke))

	// 遍历边界框内的每个像素
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			// 计算覆盖率
			coverage := calculateCircleStrokeCoverageWithSupersampling(float64(x), float64(y), centerX, centerY, radius, strokeWidth, 4)

			if coverage > 0 {
				// 混合颜色
				blendedColor := blendColors(getPixelColor(img, x, y), color, coverage)
				DrawPixel(img, x, y, blendedColor)
			}
		}
	}
}

// DrawAntiAliasedCircleWithWidth 绘制带描边宽度的抗锯齿圆形
func (r *AntiAliasedRenderer) DrawAntiAliasedCircleWithWidth(img *image.RGBA, centerX, centerY, radius float64, color color.RGBA, strokeWidth float64) {
	DrawAntiAliasedCircle(img, centerX, centerY, radius, color, strokeWidth)
}

// calculateCircleStrokeCoverageWithSupersampling 使用超采样计算圆形描边覆盖率 / Calculate circle stroke coverage with supersampling
func calculateCircleStrokeCoverageWithSupersampling(pixelX, pixelY, centerX, centerY, radius, strokeWidth float64, samples int) float64 {
	onStrokeCount := 0
	totalSamples := samples * samples
	step := 1.0 / float64(samples)

	for i := 0; i < samples; i++ {
		for j := 0; j < samples; j++ {
			// 计算子像素位置 / Calculate sub-pixel position
			sampleX := pixelX + (float64(i)+0.5)*step
			sampleY := pixelY + (float64(j)+0.5)*step

			// 计算到圆心的距离 / Calculate distance to circle center
			dx := sampleX - centerX
			dy := sampleY - centerY
			distance := math.Sqrt(dx*dx + dy*dy)

			// 检查是否在描边区域内 / Check if in stroke area
			if distance >= radius-strokeWidth/2 && distance <= radius+strokeWidth/2 {
				onStrokeCount++
			}
		}
	}

	return float64(onStrokeCount) / float64(totalSamples)
}

// calculateEllipseStrokeCoverageWithSupersampling 使用超采样计算椭圆描边覆盖率 / Calculate ellipse stroke coverage with supersampling
func calculateEllipseStrokeCoverageWithSupersampling(pixelX, pixelY, centerX, centerY, radiusX, radiusY, strokeWidth float64, samples int) float64 {
	onStrokeCount := 0
	totalSamples := samples * samples
	step := 1.0 / float64(samples)

	for i := 0; i < samples; i++ {
		for j := 0; j < samples; j++ {
			// 计算子像素位置 / Calculate sub-pixel position
			sampleX := pixelX + (float64(i)+0.5)*step
			sampleY := pixelY + (float64(j)+0.5)*step

			// 计算椭圆距离场 / Calculate ellipse distance field
			dx := sampleX - centerX
			dy := sampleY - centerY

			// 使用椭圆距离场近似 / Use ellipse distance field approximation
			p := math.Sqrt((dx*dx)/(radiusX*radiusX) + (dy*dy)/(radiusY*radiusY))
			distanceToEdge := math.Abs(p-1.0) * math.Min(radiusX, radiusY)

			// 检查是否在描边区域内 / Check if in stroke area
			if distanceToEdge <= strokeWidth/2 {
				onStrokeCount++
			}
		}
	}

	return float64(onStrokeCount) / float64(totalSamples)
}

// calculateEllipseInteriorCoverage 计算椭圆内部覆盖率 / Calculate ellipse interior coverage
func calculateEllipseInteriorCoverage(pixelX, pixelY, centerX, centerY, radiusX, radiusY float64, samples int) float64 {
	insideCount := 0
	totalSamples := samples * samples
	step := 1.0 / float64(samples)

	for i := 0; i < samples; i++ {
		for j := 0; j < samples; j++ {
			// 计算子像素位置 / Calculate sub-pixel position
			sampleX := pixelX + (float64(i)+0.5)*step
			sampleY := pixelY + (float64(j)+0.5)*step

			// 计算椭圆方程值 / Calculate ellipse equation value
			dx := sampleX - centerX
			dy := sampleY - centerY
			ellipseValue := (dx*dx)/(radiusX*radiusX) + (dy*dy)/(radiusY*radiusY)

			// 检查是否在椭圆内 / Check if inside ellipse
			if ellipseValue <= 1.0 {
				insideCount++
			}
		}
	}

	return float64(insideCount) / float64(totalSamples)
}

// DrawAntiAliasedEllipse 绘制抗锯齿椭圆（描边模式）
func DrawAntiAliasedEllipse(img *image.RGBA, centerX, centerY, radiusX, radiusY float64, color color.RGBA, strokeWidth float64) {
	// 计算边界框，考虑描边宽度
	halfStroke := strokeWidth / 2
	maxRadius := math.Max(radiusX, radiusY)
	minX := int(math.Floor(centerX - maxRadius - halfStroke))
	maxX := int(math.Ceil(centerX + maxRadius + halfStroke))
	minY := int(math.Floor(centerY - maxRadius - halfStroke))
	maxY := int(math.Ceil(centerY + maxRadius + halfStroke))

	// 遍历边界框内的每个像素
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			// 计算覆盖率
			coverage := calculateEllipseStrokeCoverageWithSupersampling(float64(x), float64(y), centerX, centerY, radiusX, radiusY, strokeWidth, 4)

			if coverage > 0 {
				// 混合颜色
				blendedColor := blendColors(getPixelColor(img, x, y), color, coverage)
				DrawPixel(img, x, y, blendedColor)
			}
		}
	}
}

// DrawAntiAliasedEllipseWithWidth 绘制带描边宽度的抗锯齿椭圆
func (r *AntiAliasedRenderer) DrawAntiAliasedEllipseWithWidth(img *image.RGBA, centerX, centerY, radiusX, radiusY float64, color color.RGBA, strokeWidth float64) {
	DrawAntiAliasedEllipse(img, centerX, centerY, radiusX, radiusY, color, strokeWidth)
}

// getPixelColor 获取像素颜色
func getPixelColor(img *image.RGBA, x, y int) color.RGBA {
	if x < 0 || y < 0 || x >= img.Bounds().Dx() || y >= img.Bounds().Dy() {
		return color.RGBA{0, 0, 0, 0}
	}
	return img.RGBAAt(x, y)
}

// blendColors 混合两种颜色
func blendColors(bg, fg color.RGBA, alpha float64) color.RGBA {
	if alpha <= 0 {
		return bg
	}
	if alpha >= 1 {
		return fg
	}

	alpha16 := uint16(alpha * 65535)
	invAlpha := 65535 - alpha16

	r := (uint16(bg.R)*invAlpha + uint16(fg.R)*alpha16) / 65535
	g := (uint16(bg.G)*invAlpha + uint16(fg.G)*alpha16) / 65535
	b := (uint16(bg.B)*invAlpha + uint16(fg.B)*alpha16) / 65535
	a := (uint16(bg.A)*invAlpha + uint16(fg.A)*alpha16) / 65535

	return color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: uint8(a),
	}
}

// smoothStep 平滑步函数 / Smooth step function
func smoothStep(t float64) float64 {
	if t <= 0 {
		return 0
	}
	if t >= 1 {
		return 1
	}
	return t * t * (3 - 2*t)
}
