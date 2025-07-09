package renderer

import (
	"image"
	"image/color"
	"math"
)

// CreateImage 创建指定大小和背景色的图像
func CreateImage(width, height int, background color.Color) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 填充背景色
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, background)
		}
	}

	return img
}

// DrawPixel 在图像上绘制像素
func DrawPixel(img *image.RGBA, x, y int, c color.Color) {
	// 检查边界
	if x < 0 || y < 0 || x >= img.Bounds().Dx() || y >= img.Bounds().Dy() {
		return
	}

	img.Set(x, y, c)
}

// DrawLine 在图像上绘制直线（使用Bresenham算法）
func DrawLine(img *image.RGBA, x0, y0, x1, y1 int, c color.Color) {
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	sx, sy := 1, 1

	if x0 >= x1 {
		sx = -1
	}

	if y0 >= y1 {
		sy = -1
	}

	err := dx - dy

	for {
		DrawPixel(img, x0, y0, c)

		if x0 == x1 && y0 == y1 {
			break
		}

		e2 := 2 * err

		if e2 > -dy {
			err -= dy
			x0 += sx
		}

		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}

// DrawAntiAliasedLine 绘制抗锯齿直线
func DrawAntiAliasedLine(img *image.RGBA, x0, y0, x1, y1 float64, c color.Color, strokeWidth float64) {
	// 使用改进的超采样抗锯齿算法
	dx := x1 - x0
	dy := y1 - y0
	length := math.Sqrt(dx*dx + dy*dy)
	
	// 如果线条太短，直接绘制点
	if length < 0.5 {
		DrawAntiAliasedPixel(img, x0, y0, c, 1.0)
		return
	}
	
	// 计算线条的单位向量（暂时不需要，注释掉避免编译错误）
	// ux := dx / length
	// uy := dy / length
	// 法向量暂时不需要，注释掉避免编译错误
	// nx := -uy // 法向量
	// ny := ux
	
	// 计算线条的边界框
	minX := math.Min(x0, x1) - strokeWidth/2 - 1
	maxX := math.Max(x0, x1) + strokeWidth/2 + 1
	minY := math.Min(y0, y1) - strokeWidth/2 - 1
	maxY := math.Max(y0, y1) + strokeWidth/2 + 1
	
	// 遍历边界框内的每个像素
	for py := int(math.Floor(minY)); py <= int(math.Ceil(maxY)); py++ {
		for px := int(math.Floor(minX)); px <= int(math.Ceil(maxX)); px++ {
			// 使用4x4超采样计算覆盖率
			coverage := calculateLineCoverage(float64(px), float64(py), x0, y0, x1, y1, strokeWidth)
			if coverage > 0 {
				blendPixelWithCoverage(img, px, py, c, coverage)
			}
		}
	}
}

// DrawAntiAliasedPixel 绘制抗锯齿像素
func DrawAntiAliasedPixel(img *image.RGBA, x, y float64, c color.Color, alpha float64) {
	px := int(math.Floor(x))
	py := int(math.Floor(y))
	
	// 检查边界
	if px < 0 || py < 0 || px >= img.Bounds().Dx() || py >= img.Bounds().Dy() {
		return
	}
	
	// 限制alpha值
	if alpha > 1.0 {
		alpha = 1.0
	}
	if alpha <= 0.0 {
		return
	}
	
	// 直接使用blendPixelWithCoverage函数，保持一致性
	blendPixelWithCoverage(img, px, py, c, alpha)
}

// calculateLineCoverage 计算像素在线条中的覆盖率 / Calculate pixel coverage in line
func calculateLineCoverage(px, py, x0, y0, x1, y1, strokeWidth float64) float64 {
	// 使用4x4超采样 / Use 4x4 supersampling
	count := 0
	for sy := 0; sy < 4; sy++ {
		for sx := 0; sx < 4; sx++ {
			// 计算子像素坐标 / Calculate sub-pixel coordinates
			sampleX := px + (float64(sx)+0.5)/4.0
			sampleY := py + (float64(sy)+0.5)/4.0
			
			// 计算点到线段的距离 / Calculate distance from point to line segment
			dist := pointToLineDistance(sampleX, sampleY, x0, y0, x1, y1)
			
			// 如果距离小于线宽的一半，则该子像素被覆盖 / If distance is less than half stroke width, sub-pixel is covered
			if dist <= strokeWidth/2.0 {
				count++
			}
		}
	}
	
	// 返回覆盖率 / Return coverage ratio
	return float64(count) / 16.0
}

// pointToLineDistance 计算点到线段的最短距离 / Calculate shortest distance from point to line segment
func pointToLineDistance(px, py, x0, y0, x1, y1 float64) float64 {
	// 线段向量 / Line segment vector
	dx := x1 - x0
	dy := y1 - y0
	
	// 如果线段长度为0，返回点到点的距离 / If line segment length is 0, return point-to-point distance
	if dx == 0 && dy == 0 {
		return math.Sqrt((px-x0)*(px-x0) + (py-y0)*(py-y0))
	}
	
	// 计算投影参数 / Calculate projection parameter
	t := ((px-x0)*dx + (py-y0)*dy) / (dx*dx + dy*dy)
	
	// 限制t在[0,1]范围内 / Clamp t to [0,1] range
	if t < 0 {
		t = 0
	} else if t > 1 {
		t = 1
	}
	
	// 计算最近点 / Calculate closest point
	closestX := x0 + t*dx
	closestY := y0 + t*dy
	
	// 返回距离 / Return distance
	return math.Sqrt((px-closestX)*(px-closestX) + (py-closestY)*(py-closestY))
}

// fpart 返回浮点数的小数部分
func fpart(x float64) float64 {
	return x - math.Floor(x)
}

// DrawRect 在图像上绘制矩形
func DrawRect(img *image.RGBA, x, y, width, height int, c color.Color, fill bool) {
	if fill {
		// 填充矩形
		for j := y; j < y+height; j++ {
			for i := x; i < x+width; i++ {
				DrawPixel(img, i, j, c)
			}
		}
	} else {
		// 绘制矩形边框
		DrawLine(img, x, y, x+width, y, c)
		DrawLine(img, x+width, y, x+width, y+height, c)
		DrawLine(img, x+width, y+height, x, y+height, c)
		DrawLine(img, x, y+height, x, y, c)
	}
}

// DrawCircle 在图像上绘制圆形（支持填充和描边）/ Draw circle with fill or stroke support
func DrawCircle(img *image.RGBA, centerX, centerY, radius int, c color.Color, fill bool) {
	if fill {
		DrawAntiAliasedFilledCircle(img, centerX, centerY, radius, c)
	} else {
		DrawAntiAliasedCircleOutline(img, centerX, centerY, radius, c)
	}
}

// DrawAntiAliasedFilledCircle 绘制抗锯齿填充圆形 / Draw anti-aliased filled circle
func DrawAntiAliasedFilledCircle(img *image.RGBA, centerX, centerY, radius int, c color.Color) {
	// 使用超采样抗锯齿算法 / Use supersampling anti-aliasing
	for y := centerY - radius - 1; y <= centerY + radius + 1; y++ {
		for x := centerX - radius - 1; x <= centerX + radius + 1; x++ {
			// 计算覆盖率 / Calculate coverage
			coverage := calculateCircleCoverage(float64(x), float64(y), float64(centerX), float64(centerY), float64(radius))
			if coverage > 0 {
				blendPixelWithCoverage(img, x, y, c, coverage)
			}
		}
	}
}

// DrawAntiAliasedCircleOutline 绘制抗锯齿圆形轮廓 / Draw anti-aliased circle outline
func DrawAntiAliasedCircleOutline(img *image.RGBA, centerX, centerY, radius int, c color.Color) {
	// 使用距离场抗锯齿算法 / Use distance field anti-aliasing
	strokeWidth := 1.0 // 默认描边宽度
	innerRadius := float64(radius) - strokeWidth/2
	outerRadius := float64(radius) + strokeWidth/2
	
	// 扩展绘制范围以包含抗锯齿边缘
	extendedRadius := int(outerRadius + 2)
	
	for y := centerY - extendedRadius; y <= centerY + extendedRadius; y++ {
		for x := centerX - extendedRadius; x <= centerX + extendedRadius; x++ {
			// 计算轮廓覆盖率
			coverage := calculateCircleOutlineCoverage(float64(x), float64(y), float64(centerX), float64(centerY), innerRadius, outerRadius)
			if coverage > 0 {
				blendPixelWithCoverage(img, x, y, c, coverage)
			}
		}
	}
}

// FillCircle 在图像上绘制填充圆形
func FillCircle(img *image.RGBA, centerX, centerY, radius int, c color.Color) {
	for y := -radius; y <= radius; y++ {
		for x := -radius; x <= radius; x++ {
			if x*x+y*y <= radius*radius {
				DrawPixel(img, centerX+x, centerY+y, c)
			}
		}
	}
}

// DrawEllipse 在图像上绘制椭圆（支持填充和描边）/ Draw ellipse with fill or stroke support
func DrawEllipse(img *image.RGBA, centerX, centerY, radiusX, radiusY int, c color.Color, fill bool) {
	if fill {
		DrawAntiAliasedFilledEllipse(img, centerX, centerY, radiusX, radiusY, c)
	} else {
		DrawAntiAliasedEllipseOutline(img, centerX, centerY, radiusX, radiusY, c)
	}
}

// DrawAntiAliasedFilledEllipse 绘制抗锯齿填充椭圆 / Draw anti-aliased filled ellipse
func DrawAntiAliasedFilledEllipse(img *image.RGBA, centerX, centerY, radiusX, radiusY int, c color.Color) {
	// 使用超采样抗锯齿算法 / Use supersampling anti-aliasing
	maxRadius := int(math.Max(float64(radiusX), float64(radiusY)))
	for y := centerY - maxRadius - 1; y <= centerY + maxRadius + 1; y++ {
		for x := centerX - maxRadius - 1; x <= centerX + maxRadius + 1; x++ {
			// 计算覆盖率 / Calculate coverage
			coverage := calculateEllipseCoverage(float64(x), float64(y), float64(centerX), float64(centerY), float64(radiusX), float64(radiusY))
			if coverage > 0 {
				blendPixelWithCoverage(img, x, y, c, coverage)
			}
		}
	}
}

// DrawAntiAliasedEllipseOutline 绘制抗锯齿椭圆轮廓 / Draw anti-aliased ellipse outline
func DrawAntiAliasedEllipseOutline(img *image.RGBA, centerX, centerY, radiusX, radiusY int, c color.Color) {
	// 使用距离场抗锯齿算法 / Use distance field anti-aliasing
	strokeWidth := 1.0 // 默认描边宽度
	
	// 计算椭圆的内外边界
	innerRadiusX := float64(radiusX) - strokeWidth/2
	innerRadiusY := float64(radiusY) - strokeWidth/2
	outerRadiusX := float64(radiusX) + strokeWidth/2
	outerRadiusY := float64(radiusY) + strokeWidth/2
	
	// 扩展绘制范围
	maxRadius := int(math.Max(outerRadiusX, outerRadiusY) + 2)
	
	for y := centerY - maxRadius; y <= centerY + maxRadius; y++ {
		for x := centerX - maxRadius; x <= centerX + maxRadius; x++ {
			// 计算椭圆轮廓覆盖率
			coverage := calculateEllipseOutlineCoverage(float64(x), float64(y), float64(centerX), float64(centerY), innerRadiusX, innerRadiusY, outerRadiusX, outerRadiusY)
			if coverage > 0 {
				blendPixelWithCoverage(img, x, y, c, coverage)
			}
		}
	}
}

// calculateEllipseCoverage 计算像素在椭圆中的覆盖率 / Calculate pixel coverage in ellipse
func calculateEllipseCoverage(x, y, centerX, centerY, radiusX, radiusY float64) float64 {
	// 使用4x4超采样 / Use 4x4 supersampling
	samples := 4
	coveredSamples := 0
	
	for sy := 0; sy < samples; sy++ {
		for sx := 0; sx < samples; sx++ {
			// 计算采样点位置 / Calculate sample point position
			sampleX := x + (float64(sx)+0.5)/float64(samples) - 0.5
			sampleY := y + (float64(sy)+0.5)/float64(samples) - 0.5
			
			// 计算椭圆方程 / Calculate ellipse equation
			dx := sampleX - centerX
			dy := sampleY - centerY
			ellipseValue := (dx*dx)/(radiusX*radiusX) + (dy*dy)/(radiusY*radiusY)
			
			// 检查是否在椭圆内 / Check if inside ellipse
			if ellipseValue <= 1.0 {
				coveredSamples++
			}
		}
	}
	
	return float64(coveredSamples) / float64(samples*samples)
}

// calculateCircleCoverage 计算像素在圆形中的覆盖率 / Calculate pixel coverage in circle
func calculateCircleCoverage(x, y, centerX, centerY, radius float64) float64 {
	// 使用4x4超采样 / Use 4x4 supersampling
	samples := 4
	coveredSamples := 0
	
	for sy := 0; sy < samples; sy++ {
		for sx := 0; sx < samples; sx++ {
			// 计算采样点位置 / Calculate sample point position
			sampleX := x + (float64(sx)+0.5)/float64(samples) - 0.5
			sampleY := y + (float64(sy)+0.5)/float64(samples) - 0.5
			
			// 计算到圆心的距离 / Calculate distance to center
			dx := sampleX - centerX
			dy := sampleY - centerY
			dist := math.Sqrt(dx*dx + dy*dy)
			
			// 检查是否在圆内 / Check if inside circle
			if dist <= radius {
				coveredSamples++
			}
		}
	}
	
	return float64(coveredSamples) / float64(samples*samples)
}

// calculateCircleOutlineCoverage 计算像素在圆形轮廓中的覆盖率 / Calculate pixel coverage in circle outline
func calculateCircleOutlineCoverage(x, y, centerX, centerY, innerRadius, outerRadius float64) float64 {
	// 使用4x4超采样 / Use 4x4 supersampling
	samples := 4
	coveredSamples := 0
	
	for sy := 0; sy < samples; sy++ {
		for sx := 0; sx < samples; sx++ {
			// 计算采样点位置 / Calculate sample point position
			sampleX := x + (float64(sx)+0.5)/float64(samples) - 0.5
			sampleY := y + (float64(sy)+0.5)/float64(samples) - 0.5
			
			// 计算到圆心的距离 / Calculate distance to center
			dx := sampleX - centerX
			dy := sampleY - centerY
			dist := math.Sqrt(dx*dx + dy*dy)
			
			// 检查是否在轮廓范围内 / Check if inside outline range
			if dist >= innerRadius && dist <= outerRadius {
				coveredSamples++
			}
		}
	}
	
	return float64(coveredSamples) / float64(samples*samples)
}

// calculateEllipseOutlineCoverage 计算像素在椭圆轮廓中的覆盖率 / Calculate pixel coverage in ellipse outline
func calculateEllipseOutlineCoverage(x, y, centerX, centerY, innerRadiusX, innerRadiusY, outerRadiusX, outerRadiusY float64) float64 {
	// 使用4x4超采样 / Use 4x4 supersampling
	samples := 4
	coveredSamples := 0
	
	for sy := 0; sy < samples; sy++ {
		for sx := 0; sx < samples; sx++ {
			// 计算采样点位置 / Calculate sample point position
			sampleX := x + (float64(sx)+0.5)/float64(samples) - 0.5
			sampleY := y + (float64(sy)+0.5)/float64(samples) - 0.5
			
			// 计算椭圆方程值 / Calculate ellipse equation values
			dx := sampleX - centerX
			dy := sampleY - centerY
			
			// 计算内椭圆和外椭圆的方程值
			innerValue := (dx*dx)/(innerRadiusX*innerRadiusX) + (dy*dy)/(innerRadiusY*innerRadiusY)
			outerValue := (dx*dx)/(outerRadiusX*outerRadiusX) + (dy*dy)/(outerRadiusY*outerRadiusY)
			
			// 检查是否在轮廓范围内（在外椭圆内但在内椭圆外）
			if outerValue <= 1.0 && innerValue >= 1.0 {
				coveredSamples++
			}
		}
	}
	
	return float64(coveredSamples) / float64(samples*samples)
}

// blendPixelWithCoverage 根据覆盖率混合像素 / Blend pixel with coverage
func blendPixelWithCoverage(img *image.RGBA, x, y int, c color.Color, coverage float64) {
	if x < 0 || y < 0 || x >= img.Bounds().Dx() || y >= img.Bounds().Dy() {
		return
	}
	
	// 限制覆盖率范围 / Clamp coverage range
	if coverage <= 0.0 {
		return
	}
	if coverage > 1.0 {
		coverage = 1.0
	}
	
	// 获取当前像素颜色 / Get current pixel color
	existingColor := img.RGBAAt(x, y)
	r1, g1, b1, a1 := float64(existingColor.R), float64(existingColor.G), float64(existingColor.B), float64(existingColor.A)
	
	// 获取新颜色 / Get new color
	r2, g2, b2, a2 := c.RGBA()
	newR := float64(r2 >> 8)
	newG := float64(g2 >> 8)
	newB := float64(b2 >> 8)
	newA := float64(a2 >> 8)
	
	// 使用覆盖率进行Alpha混合 / Use coverage for alpha blending
	alpha := coverage
	invAlpha := 1.0 - alpha
	
	// 计算最终颜色 / Calculate final color
	finalR := uint8(r1*invAlpha + newR*alpha)
	finalG := uint8(g1*invAlpha + newG*alpha)
	finalB := uint8(b1*invAlpha + newB*alpha)
	finalA := uint8(math.Max(a1, newA*coverage))
	
	img.SetRGBA(x, y, color.RGBA{R: finalR, G: finalG, B: finalB, A: finalA})
}

// abs 返回整数的绝对值
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
