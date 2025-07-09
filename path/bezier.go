package path

import (
	"math"

	"github.com/hoonfeng/svg/types"
)

// adaptiveCubicBezierFlattening 使用自适应算法将三次贝塞尔曲线平滑化为点列表
func adaptiveCubicBezierFlattening(p0, p1, p2, p3 types.Point, flatness float64) []types.Point {
	// 初始化点列表
	points := []types.Point{p0}

	// 递归平滑化
	recursiveCubicBezier(p0, p1, p2, p3, flatness, &points)

	// 添加终点
	points = append(points, p3)

	return points
}

// recursiveCubicBezier 递归平滑化三次贝塞尔曲线
func recursiveCubicBezier(p0, p1, p2, p3 types.Point, flatness float64, points *[]types.Point) {
	// 改进的误差计算方法 - 使用更精确的曲率估计
	// Improved error calculation method - using more accurate curvature estimation
	dx := p3.X - p0.X
	dy := p3.Y - p0.Y
	length := dx*dx + dy*dy

	// 如果起点和终点太近，直接返回
	// If start and end points are too close, return directly
	if length < 1e-10 {
		return
	}

	// 计算控制点到起点-终点直线的距离
	// Calculate distance from control points to start-end line
	d1 := ((p1.X-p0.X)*dy - (p1.Y-p0.Y)*dx)
	d2 := ((p2.X-p0.X)*dy - (p2.Y-p0.Y)*dx)

	// 使用绝对值和更合理的阈值
	// Use absolute values and more reasonable threshold
	maxDist := math.Max(math.Abs(d1), math.Abs(d2))
	maxDist = maxDist * maxDist / length

	// 如果误差小于阈值，则不需要继续细分
	// If error is less than threshold, no need to continue subdivision
	if maxDist <= flatness*flatness {
		return
	}

	// 细分曲线（de Casteljau算法）
	p01 := types.Point{X: (p0.X + p1.X) / 2, Y: (p0.Y + p1.Y) / 2}
	p12 := types.Point{X: (p1.X + p2.X) / 2, Y: (p1.Y + p2.Y) / 2}
	p23 := types.Point{X: (p2.X + p3.X) / 2, Y: (p2.Y + p3.Y) / 2}
	p012 := types.Point{X: (p01.X + p12.X) / 2, Y: (p01.Y + p12.Y) / 2}
	p123 := types.Point{X: (p12.X + p23.X) / 2, Y: (p12.Y + p23.Y) / 2}
	p0123 := types.Point{X: (p012.X + p123.X) / 2, Y: (p012.Y + p123.Y) / 2}

	// 递归处理两个子曲线
	recursiveCubicBezier(p0, p01, p012, p0123, flatness, points)
	*points = append(*points, p0123)
	recursiveCubicBezier(p0123, p123, p23, p3, flatness, points)
}

// adaptiveQuadraticBezierFlattening 使用自适应算法将二次贝塞尔曲线平滑化为点列表
func adaptiveQuadraticBezierFlattening(p0, p1, p2 types.Point, flatness float64) []types.Point {
	// 初始化点列表
	points := []types.Point{p0}

	// 递归平滑化
	recursiveQuadraticBezier(p0, p1, p2, flatness, &points)

	// 添加终点
	points = append(points, p2)

	return points
}

// recursiveQuadraticBezier 递归平滑化二次贝塞尔曲线
func recursiveQuadraticBezier(p0, p1, p2 types.Point, flatness float64, points *[]types.Point) {
	// 改进的误差计算方法 - 使用更精确的曲率估计
	// Improved error calculation method - using more accurate curvature estimation
	dx := p2.X - p0.X
	dy := p2.Y - p0.Y
	length := dx*dx + dy*dy

	// 如果起点和终点太近，直接返回
	// If start and end points are too close, return directly
	if length < 1e-10 {
		return
	}

	// 计算控制点到起点-终点直线的距离
	// Calculate distance from control point to start-end line
	d := ((p1.X-p0.X)*dy - (p1.Y-p0.Y)*dx)
	maxDist := math.Abs(d) * math.Abs(d) / length

	// 如果误差小于阈值，则不需要继续细分
	// If error is less than threshold, no need to continue subdivision
	if maxDist <= flatness*flatness {
		return
	}

	// 细分曲线（de Casteljau算法）
	p01 := types.Point{X: (p0.X + p1.X) / 2, Y: (p0.Y + p1.Y) / 2}
	p12 := types.Point{X: (p1.X + p2.X) / 2, Y: (p1.Y + p2.Y) / 2}
	p012 := types.Point{X: (p01.X + p12.X) / 2, Y: (p01.Y + p12.Y) / 2}

	// 递归处理两个子曲线
	recursiveQuadraticBezier(p0, p01, p012, flatness, points)
	*points = append(*points, p012)
	recursiveQuadraticBezier(p012, p12, p2, flatness, points)
}

// FlattenPath 将路径平滑化为点列表
func (p *SVGPath) FlattenPath(precision float64) []types.Point {
	// 创建路径上下文
	ctx := NewPathContext()

	// 执行所有命令
	for _, cmd := range p.Commands {
		cmd.Execute(ctx, precision)
	}

	return ctx.Points
}

// FlattenSubPaths 将路径分解为多个子路径
func (p *SVGPath) FlattenSubPaths(precision float64) [][]types.Point {
	subPaths := [][]types.Point{}
	ctx := NewPathContext()
	subPathStartIndex := 0

	for _, cmd := range p.Commands {
		// 检查是否是闭合路径命令
		if _, isClose := cmd.(*ClosePathCommand); isClose {
			// 执行闭合命令
			cmd.Execute(ctx, precision)
			// 提取当前子路径的点
			if len(ctx.Points) > subPathStartIndex {
				currentSubPath := ctx.Points[subPathStartIndex:]
				if len(currentSubPath) >= 2 { // 降低要求，允许线条
					// 复制子路径点
					subPath := make([]types.Point, len(currentSubPath))
					copy(subPath, currentSubPath)
					subPaths = append(subPaths, subPath)
				}
				// 更新下一个子路径的起始索引
				subPathStartIndex = len(ctx.Points)
			}
		} else {
			// 执行其他命令
			cmd.Execute(ctx, precision)
		}
	}

	// 处理最后一个未闭合的子路径
	if len(ctx.Points) > subPathStartIndex {
		currentSubPath := ctx.Points[subPathStartIndex:]
		if len(currentSubPath) >= 2 { // 降低要求，允许线条
			// 复制子路径点
			subPath := make([]types.Point, len(currentSubPath))
			copy(subPath, currentSubPath)
			subPaths = append(subPaths, subPath)
		}
	}

	return subPaths
}

// GetSubPathCloseInfo 获取每个子路径的闭合信息
func (p *SVGPath) GetSubPathCloseInfo() []bool {
	closeInfo := []bool{}
	ctx := NewPathContext()
	subPathStartIndex := 0
	currentSubPathClosed := false

	for _, cmd := range p.Commands {
		// 检查是否是移动命令（开始新的子路径）
		if _, isMove := cmd.(*MoveToCommand); isMove {
			// 如果之前有子路径，记录其闭合状态
			if len(ctx.Points) > subPathStartIndex {
				currentSubPath := ctx.Points[subPathStartIndex:]
				if len(currentSubPath) >= 2 {
					closeInfo = append(closeInfo, currentSubPathClosed)
				}
				subPathStartIndex = len(ctx.Points)
			}
			currentSubPathClosed = false // 重置闭合状态
		}

		// 检查是否是闭合路径命令
		if _, isClose := cmd.(*ClosePathCommand); isClose {
			currentSubPathClosed = true
			// 执行闭合命令
			cmd.Execute(ctx, 0.001)
			// 记录当前子路径的闭合状态
			if len(ctx.Points) > subPathStartIndex {
				currentSubPath := ctx.Points[subPathStartIndex:]
				if len(currentSubPath) >= 2 {
					closeInfo = append(closeInfo, true)
				}
				subPathStartIndex = len(ctx.Points)
			}
			currentSubPathClosed = false // 重置为下一个子路径
		} else {
			// 执行其他命令
			cmd.Execute(ctx, 0.001) // 使用固定精度
		}
	}

	// 处理最后一个未闭合的子路径
	if len(ctx.Points) > subPathStartIndex {
		currentSubPath := ctx.Points[subPathStartIndex:]
		if len(currentSubPath) >= 2 {
			closeInfo = append(closeInfo, currentSubPathClosed)
		}
	}

	return closeInfo
}
