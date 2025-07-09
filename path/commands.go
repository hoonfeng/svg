package path

import (
	"fmt"
	"math"
	"strconv"

	"github.com/hoonfeng/svg/types"
)

// PathContext 表示路径上下文，用于跟踪当前点和控制点
type PathContext struct {
	CurrentPoint types.Point
	StartPoint   types.Point
	PrevControl  types.Point // 上一个控制点（用于平滑曲线）
	Points       []types.Point
}

// NewPathContext 创建新的路径上下文
func NewPathContext() *PathContext {
	return &PathContext{
		Points: []types.Point{},
	}
}

// MoveToCommand 表示移动命令
type MoveToCommand struct {
	X, Y     float64
	Relative bool
}

func (c *MoveToCommand) Execute(ctx *PathContext, precision float64) {
	if c.Relative {
		ctx.CurrentPoint.X += c.X
		ctx.CurrentPoint.Y += c.Y
	} else {
		ctx.CurrentPoint.X = c.X
		ctx.CurrentPoint.Y = c.Y
	}

	ctx.StartPoint = ctx.CurrentPoint
	ctx.Points = append(ctx.Points, ctx.CurrentPoint)
}

func (c *MoveToCommand) String() string {
	if c.Relative {
		return "m " + strconv.FormatFloat(c.X, 'f', -1, 64) + " " + strconv.FormatFloat(c.Y, 'f', -1, 64)
	}
	return "M " + strconv.FormatFloat(c.X, 'f', -1, 64) + " " + strconv.FormatFloat(c.Y, 'f', -1, 64)
}

// LineToCommand 表示直线命令
type LineToCommand struct {
	X, Y     float64
	Relative bool
}

func (c *LineToCommand) Execute(ctx *PathContext, precision float64) {
	var endPoint types.Point
	if c.Relative {
		endPoint = types.Point{
			X: ctx.CurrentPoint.X + c.X,
			Y: ctx.CurrentPoint.Y + c.Y,
		}
	} else {
		endPoint = types.Point{
			X: c.X,
			Y: c.Y,
		}
	}

	ctx.Points = append(ctx.Points, endPoint)
	ctx.CurrentPoint = endPoint
}

func (c *LineToCommand) String() string {
	if c.Relative {
		return "l " + strconv.FormatFloat(c.X, 'f', -1, 64) + " " + strconv.FormatFloat(c.Y, 'f', -1, 64)
	}
	return "L " + strconv.FormatFloat(c.X, 'f', -1, 64) + " " + strconv.FormatFloat(c.Y, 'f', -1, 64)
}

// HorizontalLineToCommand 表示水平线命令
type HorizontalLineToCommand struct {
	X        float64
	Relative bool
}

func (c *HorizontalLineToCommand) Execute(ctx *PathContext, precision float64) {
	var endPoint types.Point
	if c.Relative {
		endPoint = types.Point{
			X: ctx.CurrentPoint.X + c.X,
			Y: ctx.CurrentPoint.Y,
		}
	} else {
		endPoint = types.Point{
			X: c.X,
			Y: ctx.CurrentPoint.Y,
		}
	}

	ctx.Points = append(ctx.Points, endPoint)
	ctx.CurrentPoint = endPoint
}

func (c *HorizontalLineToCommand) String() string {
	if c.Relative {
		return "h " + strconv.FormatFloat(c.X, 'f', -1, 64)
	}
	return "H " + strconv.FormatFloat(c.X, 'f', -1, 64)
}

// VerticalLineToCommand 表示垂直线命令
type VerticalLineToCommand struct {
	Y        float64
	Relative bool
}

func (c *VerticalLineToCommand) Execute(ctx *PathContext, precision float64) {
	var endPoint types.Point
	if c.Relative {
		endPoint = types.Point{
			X: ctx.CurrentPoint.X,
			Y: ctx.CurrentPoint.Y + c.Y,
		}
	} else {
		endPoint = types.Point{
			X: ctx.CurrentPoint.X,
			Y: c.Y,
		}
	}

	ctx.Points = append(ctx.Points, endPoint)
	ctx.CurrentPoint = endPoint
}

func (c *VerticalLineToCommand) String() string {
	if c.Relative {
		return "v " + strconv.FormatFloat(c.Y, 'f', -1, 64)
	}
	return "V " + strconv.FormatFloat(c.Y, 'f', -1, 64)
}

// CubicCurveToCommand 表示三次贝塞尔曲线命令
type CubicCurveToCommand struct {
	X1, Y1   float64 // 控制点1
	X2, Y2   float64 // 控制点2
	X, Y     float64 // 终点
	Relative bool
}

func (c *CubicCurveToCommand) Execute(ctx *PathContext, precision float64) {
	var startPoint = ctx.CurrentPoint
	var control1, control2, endPoint types.Point

	if c.Relative {
		control1 = types.Point{
			X: startPoint.X + c.X1,
			Y: startPoint.Y + c.Y1,
		}
		control2 = types.Point{
			X: startPoint.X + c.X2,
			Y: startPoint.Y + c.Y2,
		}
		endPoint = types.Point{
			X: startPoint.X + c.X,
			Y: startPoint.Y + c.Y,
		}
	} else {
		control1 = types.Point{X: c.X1, Y: c.Y1}
		control2 = types.Point{X: c.X2, Y: c.Y2}
		endPoint = types.Point{X: c.X, Y: c.Y}
	}

	// 智能自适应flatness计算 / Intelligent adaptive flatness calculation
	// 考虑控制点偏离程度和曲线长度 / Consider control point deviation and curve length
	curveLength := math.Sqrt(math.Pow(endPoint.X-startPoint.X, 2) + math.Pow(endPoint.Y-startPoint.Y, 2))
	controlDist1 := math.Sqrt(math.Pow(control1.X-startPoint.X, 2) + math.Pow(control1.Y-startPoint.Y, 2))
	controlDist2 := math.Sqrt(math.Pow(control2.X-endPoint.X, 2) + math.Pow(control2.Y-endPoint.Y, 2))
	maxControlDist := math.Max(controlDist1, controlDist2)
	// 基于曲线复杂度的智能flatness：控制点偏离越大，需要更精细的平滑
	// Intelligent flatness based on curve complexity: greater control point deviation requires finer smoothing
	complexityFactor := math.Min(10.0, maxControlDist/math.Max(1.0, curveLength))
	flatness := math.Min(2.0, math.Max(0.05, 0.5/complexityFactor)) // 更精细的自适应范围 / More refined adaptive range
	bezierPoints := adaptiveCubicBezierFlattening(startPoint, control1, control2, endPoint, flatness)
	// 跳过起点，因为它已经在路径中 / Skip start point as it's already in the path
	if len(bezierPoints) > 1 {
		ctx.Points = append(ctx.Points, bezierPoints[1:]...)
	}
	ctx.CurrentPoint = endPoint
	ctx.PrevControl = control2 // 保存最后一个控制点 / Save last control point
}

func (c *CubicCurveToCommand) String() string {
	if c.Relative {
		return fmt.Sprintf("c %.2f %.2f %.2f %.2f %.2f %.2f",
			c.X1, c.Y1, c.X2, c.Y2, c.X, c.Y)
	}
	return fmt.Sprintf("C %.2f %.2f %.2f %.2f %.2f %.2f",
		c.X1, c.Y1, c.X2, c.Y2, c.X, c.Y)
}

// SmoothCubicCurveToCommand 表示平滑三次贝塞尔曲线命令
type SmoothCubicCurveToCommand struct {
	X2, Y2   float64 // 控制点2
	X, Y     float64 // 终点
	Relative bool
}

func (c *SmoothCubicCurveToCommand) Execute(ctx *PathContext, precision float64) {
	var startPoint = ctx.CurrentPoint
	var control1, control2, endPoint types.Point

	// 计算反射的控制点1 / Calculate reflected control point 1
	if ctx.PrevControl != (types.Point{}) {
		control1 = types.Point{
			X: 2*startPoint.X - ctx.PrevControl.X,
			Y: 2*startPoint.Y - ctx.PrevControl.Y,
		}
	} else {
		control1 = startPoint
	}

	if c.Relative {
		control2 = types.Point{
			X: startPoint.X + c.X2,
			Y: startPoint.Y + c.Y2,
		}
		endPoint = types.Point{
			X: startPoint.X + c.X,
			Y: startPoint.Y + c.Y,
		}
	} else {
		control2 = types.Point{X: c.X2, Y: c.Y2}
		endPoint = types.Point{X: c.X, Y: c.Y}
	}

	// 智能自适应flatness计算 / Intelligent adaptive flatness calculation
	// 考虑控制点偏离程度和曲线长度 / Consider control point deviation and curve length
	curveLength := math.Sqrt(math.Pow(endPoint.X-startPoint.X, 2) + math.Pow(endPoint.Y-startPoint.Y, 2))
	controlDist1 := math.Sqrt(math.Pow(control1.X-startPoint.X, 2) + math.Pow(control1.Y-startPoint.Y, 2))
	controlDist2 := math.Sqrt(math.Pow(control2.X-endPoint.X, 2) + math.Pow(control2.Y-endPoint.Y, 2))
	maxControlDist := math.Max(controlDist1, controlDist2)
	// 基于曲线复杂度的智能flatness：控制点偏离越大，需要更精细的平滑
	// Intelligent flatness based on curve complexity: greater control point deviation requires finer smoothing
	complexityFactor := math.Min(10.0, maxControlDist/math.Max(1.0, curveLength))
	flatness := math.Min(2.0, math.Max(0.05, 0.5/complexityFactor)) // 更精细的自适应范围 / More refined adaptive range
	bezierPoints := adaptiveCubicBezierFlattening(startPoint, control1, control2, endPoint, flatness)
	// 跳过起点，因为它已经在路径中 / Skip start point as it's already in the path
	if len(bezierPoints) > 1 {
		ctx.Points = append(ctx.Points, bezierPoints[1:]...)
	}
	ctx.CurrentPoint = endPoint
	ctx.PrevControl = control2
}

func (c *SmoothCubicCurveToCommand) String() string {
	if c.Relative {
		return fmt.Sprintf("s %.2f %.2f %.2f %.2f",
			c.X2, c.Y2, c.X, c.Y)
	}
	return fmt.Sprintf("S %.2f %.2f %.2f %.2f",
		c.X2, c.Y2, c.X, c.Y)
}

// QuadraticCurveToCommand 表示二次贝塞尔曲线命令
type QuadraticCurveToCommand struct {
	X1, Y1   float64 // 控制点
	X, Y     float64 // 终点
	Relative bool
}

func (c *QuadraticCurveToCommand) Execute(ctx *PathContext, precision float64) {
	var startPoint = ctx.CurrentPoint
	var control, endPoint types.Point

	if c.Relative {
		control = types.Point{
			X: startPoint.X + c.X1,
			Y: startPoint.Y + c.Y1,
		}
		endPoint = types.Point{
			X: startPoint.X + c.X,
			Y: startPoint.Y + c.Y,
		}
	} else {
		control = types.Point{X: c.X1, Y: c.Y1}
		endPoint = types.Point{X: c.X, Y: c.Y}
	}

	// 智能自适应flatness计算 / Intelligent adaptive flatness calculation
	// 考虑控制点偏离程度和曲线长度 / Consider control point deviation and curve length
	curveLength := math.Sqrt(math.Pow(endPoint.X-startPoint.X, 2) + math.Pow(endPoint.Y-startPoint.Y, 2))
	// 计算控制点到起点-终点直线的距离 / Calculate distance from control point to start-end line
	midPoint := types.Point{X: (startPoint.X + endPoint.X) / 2, Y: (startPoint.Y + endPoint.Y) / 2}
	controlDist := math.Sqrt(math.Pow(control.X-midPoint.X, 2) + math.Pow(control.Y-midPoint.Y, 2))
	// 基于曲线复杂度的智能flatness：控制点偏离越大，需要更精细的平滑
	// Intelligent flatness based on curve complexity: greater control point deviation requires finer smoothing
	complexityFactor := math.Min(10.0, controlDist/math.Max(1.0, curveLength/2))
	flatness := math.Min(2.0, math.Max(0.05, 0.5/complexityFactor)) // 更精细的自适应范围 / More refined adaptive range
	bezierPoints := adaptiveQuadraticBezierFlattening(startPoint, control, endPoint, flatness)
	// 跳过起点，因为它已经在路径中 / Skip start point as it's already in the path
	if len(bezierPoints) > 1 {
		ctx.Points = append(ctx.Points, bezierPoints[1:]...)
	}
	ctx.CurrentPoint = endPoint
	ctx.PrevControl = control
}

func (c *QuadraticCurveToCommand) String() string {
	if c.Relative {
		return fmt.Sprintf("q %.2f %.2f %.2f %.2f",
			c.X1, c.Y1, c.X, c.Y)
	}
	return fmt.Sprintf("Q %.2f %.2f %.2f %.2f",
		c.X1, c.Y1, c.X, c.Y)
}

// SmoothQuadraticCurveToCommand 表示平滑二次贝塞尔曲线命令
type SmoothQuadraticCurveToCommand struct {
	X, Y     float64 // 终点
	Relative bool
}

func (c *SmoothQuadraticCurveToCommand) Execute(ctx *PathContext, precision float64) {
	var startPoint = ctx.CurrentPoint
	var control, endPoint types.Point

	// 计算反射的控制点 / Calculate reflected control point
	if ctx.PrevControl != (types.Point{}) {
		control = types.Point{
			X: 2*startPoint.X - ctx.PrevControl.X,
			Y: 2*startPoint.Y - ctx.PrevControl.Y,
		}
	} else {
		control = startPoint
	}

	if c.Relative {
		endPoint = types.Point{
			X: startPoint.X + c.X,
			Y: startPoint.Y + c.Y,
		}
	} else {
		endPoint = types.Point{X: c.X, Y: c.Y}
	}

	// 智能自适应flatness计算 / Intelligent adaptive flatness calculation
	// 考虑控制点偏离程度和曲线长度 / Consider control point deviation and curve length
	curveLength := math.Sqrt(math.Pow(endPoint.X-startPoint.X, 2) + math.Pow(endPoint.Y-startPoint.Y, 2))
	// 计算控制点到起点-终点直线的距离 / Calculate distance from control point to start-end line
	midPoint := types.Point{X: (startPoint.X + endPoint.X) / 2, Y: (startPoint.Y + endPoint.Y) / 2}
	controlDist := math.Sqrt(math.Pow(control.X-midPoint.X, 2) + math.Pow(control.Y-midPoint.Y, 2))
	// 基于曲线复杂度的智能flatness：控制点偏离越大，需要更精细的平滑
	// Intelligent flatness based on curve complexity: greater control point deviation requires finer smoothing
	complexityFactor := math.Min(10.0, controlDist/math.Max(1.0, curveLength/2))
	flatness := math.Min(2.0, math.Max(0.05, 0.5/complexityFactor)) // 更精细的自适应范围 / More refined adaptive range
	bezierPoints := adaptiveQuadraticBezierFlattening(startPoint, control, endPoint, flatness)
	// 跳过起点，因为它已经在路径中 / Skip start point as it's already in the path
	if len(bezierPoints) > 1 {
		ctx.Points = append(ctx.Points, bezierPoints[1:]...)
	}
	ctx.CurrentPoint = endPoint
	ctx.PrevControl = control
}

func (c *SmoothQuadraticCurveToCommand) String() string {
	if c.Relative {
		return fmt.Sprintf("t %.2f %.2f", c.X, c.Y)
	}
	return fmt.Sprintf("T %.2f %.2f", c.X, c.Y)
}

// boolToInt 将布尔值转换为0或1
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// ArcToCommand 表示椭圆弧命令
type ArcToCommand struct {
	RX, RY        float64 // 半径
	XAxisRotation float64 // x轴旋转角度（度）
	LargeArc      bool    // 大弧标志
	Sweep         bool    // 扫掠标志
	X, Y          float64 // 终点
	Relative      bool
}

func (c *ArcToCommand) Execute(ctx *PathContext, precision float64) {
	// 使用完整的椭圆弧实现，转换为贝塞尔曲线 / Use complete elliptical arc implementation, convert to Bezier curves
	startPoint := ctx.CurrentPoint
	var endPoint types.Point

	if c.Relative {
		endPoint = types.Point{
			X: startPoint.X + c.X,
			Y: startPoint.Y + c.Y,
		}
	} else {
		endPoint = types.Point{
			X: c.X,
			Y: c.Y,
		}
	}

	// 如果起点和终点相同，不需要绘制弧 / If start and end points are the same, no need to draw arc
	if startPoint.X == endPoint.X && startPoint.Y == endPoint.Y {
		return
	}

	// 转换为中心参数化 / Convert to center parameterization
	rx, ry := math.Abs(c.RX), math.Abs(c.RY)
	if rx == 0 || ry == 0 {
		// 如果半径为0，绘制直线 / If radius is 0, draw a line
		ctx.Points = append(ctx.Points, endPoint)
		ctx.CurrentPoint = endPoint
		ctx.PrevControl = types.Point{}
		return
	}

	xAxisRot := c.XAxisRotation * math.Pi / 180

	// 计算中间参数 / Calculate intermediate parameters
	dx := (startPoint.X - endPoint.X) / 2
	dy := (startPoint.Y - endPoint.Y) / 2
	x1 := math.Cos(xAxisRot)*dx + math.Sin(xAxisRot)*dy
	y1 := -math.Sin(xAxisRot)*dx + math.Cos(xAxisRot)*dy

	// 校正半径 / Correct radii
	lambda := (x1*x1)/(rx*rx) + (y1*y1)/(ry*ry)
	if lambda > 1 {
		rx *= math.Sqrt(lambda)
		ry *= math.Sqrt(lambda)
	}

	// 计算中心点 / Calculate center point
	sign := 1.0
	if c.LargeArc == c.Sweep {
		sign = -1.0
	}
	sqrt_val := (rx*rx*ry*ry - rx*rx*y1*y1 - ry*ry*x1*x1) / (rx*rx*y1*y1 + ry*ry*x1*x1)
	if sqrt_val < 0 {
		sqrt_val = 0
	}
	coeff := sign * math.Sqrt(sqrt_val)
	cx1 := coeff * rx * y1 / ry
	cy1 := -coeff * ry * x1 / rx

	cx := math.Cos(xAxisRot)*cx1 - math.Sin(xAxisRot)*cy1 + (startPoint.X+endPoint.X)/2
	cy := math.Sin(xAxisRot)*cx1 + math.Cos(xAxisRot)*cy1 + (startPoint.Y+endPoint.Y)/2

	// 计算起始角和终止角 / Calculate start and end angles
	angle := func(ux, uy, vx, vy float64) float64 {
		dot := ux*vx + uy*vy
		det := ux*vy - uy*vx
		return math.Atan2(det, dot)
	}
	theta1 := angle(1, 0, (x1-cx1)/rx, (y1-cy1)/ry)
	dtheta := angle((x1-cx1)/rx, (y1-cy1)/ry, (-x1-cx1)/rx, (-y1-cy1)/ry)

	if !c.Sweep && dtheta > 0 {
		dtheta -= 2 * math.Pi
	} else if c.Sweep && dtheta < 0 {
		dtheta += 2 * math.Pi
	}

	// 将椭圆弧转换为贝塞尔曲线段 / Convert elliptical arc to Bezier curve segments
	segments := int(math.Ceil(math.Abs(dtheta) / (math.Pi / 2)))
	if segments == 0 {
		segments = 1
	}
	delta := dtheta / float64(segments)
	t := (8.0 / 3.0) * math.Sin(delta/4) * math.Sin(delta/4) / math.Sin(delta/2)

	for i := 0; i < segments; i++ {
		cosTheta1 := math.Cos(theta1)
		sinTheta1 := math.Sin(theta1)
		cosTheta2 := math.Cos(theta1 + delta)
		sinTheta2 := math.Sin(theta1 + delta)

		// 计算贝塞尔控制点 / Calculate Bezier control points
		p0 := types.Point{
			X: cx + rx*cosTheta1*math.Cos(xAxisRot) - ry*sinTheta1*math.Sin(xAxisRot),
			Y: cy + rx*cosTheta1*math.Sin(xAxisRot) + ry*sinTheta1*math.Cos(xAxisRot),
		}
		p1 := types.Point{
			X: cx + rx*(cosTheta1-t*sinTheta1)*math.Cos(xAxisRot) - ry*(sinTheta1+t*cosTheta1)*math.Sin(xAxisRot),
			Y: cy + rx*(cosTheta1-t*sinTheta1)*math.Sin(xAxisRot) + ry*(sinTheta1+t*cosTheta1)*math.Cos(xAxisRot),
		}
		p2 := types.Point{
			X: cx + rx*(cosTheta2+t*sinTheta2)*math.Cos(xAxisRot) - ry*(sinTheta2-t*cosTheta2)*math.Sin(xAxisRot),
			Y: cy + rx*(cosTheta2+t*sinTheta2)*math.Sin(xAxisRot) + ry*(sinTheta2-t*cosTheta2)*math.Cos(xAxisRot),
		}
		p3 := types.Point{
			X: cx + rx*cosTheta2*math.Cos(xAxisRot) - ry*sinTheta2*math.Sin(xAxisRot),
			Y: cy + rx*cosTheta2*math.Sin(xAxisRot) + ry*sinTheta2*math.Cos(xAxisRot),
		}

		// 使用更精细的flatness值进行自适应平滑化 / Use more refined flatness value for adaptive flattening
		// 根据曲线复杂度动态调整flatness / Dynamically adjust flatness based on curve complexity
		curveLength := math.Sqrt(math.Pow(p3.X-p0.X, 2) + math.Pow(p3.Y-p0.Y, 2))
		flatness := math.Min(1.0, math.Max(0.1, curveLength/100.0)) // 基于曲线长度的自适应flatness / Adaptive flatness based on curve length
		bezierPoints := adaptiveCubicBezierFlattening(p0, p1, p2, p3, flatness)
		// 跳过起点（除了第一段）/ Skip start point (except for first segment)
		if i == 0 {
			ctx.Points = append(ctx.Points, bezierPoints...)
		} else {
			ctx.Points = append(ctx.Points, bezierPoints[1:]...)
		}

		theta1 += delta
	}

	ctx.CurrentPoint = endPoint
	ctx.PrevControl = types.Point{} // 重置控制点 / Reset control point
}

func (c *ArcToCommand) String() string {
	if c.Relative {
		return fmt.Sprintf("a %.2f %.2f %.2f %d %d %.2f %.2f",
			c.RX, c.RY, c.XAxisRotation, boolToInt(c.LargeArc), boolToInt(c.Sweep), c.X, c.Y)
	}
	return fmt.Sprintf("A %.2f %.2f %.2f %d %d %.2f %.2f",
		c.RX, c.RY, c.XAxisRotation, boolToInt(c.LargeArc), boolToInt(c.Sweep), c.X, c.Y)
}

// ClosePathCommand 表示闭合路径命令
type ClosePathCommand struct{}

func (c *ClosePathCommand) Execute(ctx *PathContext, precision float64) {
	if ctx.CurrentPoint != ctx.StartPoint {
		// 添加直线回到起点
		ctx.Points = append(ctx.Points, ctx.StartPoint)
		ctx.CurrentPoint = ctx.StartPoint
	}
}

func (c *ClosePathCommand) String() string {
	return "Z" // 闭合路径命令没有相对/绝对之分，总是使用大写Z
}

// ArcToAbs 表示绝对坐标的椭圆弧命令
type ArcToAbs struct {
	RX, RY        float64 // 半径
	XAxisRotation float64 // x轴旋转角度（度）
	LargeArc      bool    // 大弧标志
	Sweep         bool    // 扫掠标志
	X, Y          float64 // 终点
	DX, DY        float64 // 相对坐标
}

func (c *ArcToAbs) Execute(ctx *PathContext, precision float64) {
	startPoint := ctx.CurrentPoint
	endPoint := types.Point{X: c.X, Y: c.Y}

	// 转换为中心参数化
	rx, ry := math.Abs(c.RX), math.Abs(c.RY)
	xAxisRot := c.XAxisRotation * math.Pi / 180

	// 计算中间参数
	dx := (startPoint.X - endPoint.X) / 2
	dy := (startPoint.Y - endPoint.Y) / 2
	x1 := math.Cos(xAxisRot)*dx + math.Sin(xAxisRot)*dy
	y1 := -math.Sin(xAxisRot)*dx + math.Cos(xAxisRot)*dy

	// 校正半径
	lambda := (x1*x1)/(rx*rx) + (y1*y1)/(ry*ry)
	if lambda > 1 {
		rx *= math.Sqrt(lambda)
		ry *= math.Sqrt(lambda)
	}

	// 计算中心点
	sign := 1.0
	if c.LargeArc == c.Sweep {
		sign = -1.0
	}
	cx := sign * math.Sqrt((rx*rx*ry*ry-rx*rx*y1*y1-ry*ry*x1*x1)/(rx*rx*y1*y1+ry*ry*x1*x1))
	if math.IsNaN(cx) {
		cx = 0
	}
	cx = math.Cos(xAxisRot)*cx*rx/ry*y1 - math.Sin(xAxisRot)*cx*ry/rx*x1 + (startPoint.X+endPoint.X)/2
	cy := math.Sin(xAxisRot)*cx*rx/ry*y1 + math.Cos(xAxisRot)*cx*ry/rx*x1 + (startPoint.Y+endPoint.Y)/2

	// 计算起始角和终止角
	angle := func(u, v float64) float64 {
		return math.Atan2(v, u)
	}
	theta1 := angle((x1-cx)/rx, (y1-cy)/ry)
	theta2 := angle((-x1-cx)/rx, (-y1-cy)/ry)

	// 将椭圆弧转换为贝塞尔曲线段
	segments := int(math.Ceil(math.Abs(theta2-theta1) / (math.Pi / 2)))
	delta := (theta2 - theta1) / float64(segments)
	t := 8.0 / 3.0 * math.Sin(delta/4) * math.Sin(delta/4) / math.Sin(delta/2)

	for i := 0; i < segments; i++ {
		cos1 := math.Cos(theta1)
		sin1 := math.Sin(theta1)
		cos2 := math.Cos(theta1 + delta)
		sin2 := math.Sin(theta1 + delta)

		// 计算贝塞尔控制点
		point1 := types.Point{
			X: cx + rx*cos1 - (rx*sin1)*t,
			Y: cy + ry*sin1 + (ry*cos1)*t,
		}
		point2 := types.Point{
			X: cx + rx*cos2 + (rx*sin2)*t,
			Y: cy + ry*sin2 - (ry*cos2)*t,
		}
		point3 := types.Point{
			X: cx + rx*cos2,
			Y: cy + ry*sin2,
		}

		// 添加点到路径
		if i == 0 {
			ctx.Points = append(ctx.Points, point1)
		}
		ctx.Points = append(ctx.Points, point2, point3)
		theta1 += delta
	}

	ctx.CurrentPoint = endPoint
	ctx.PrevControl = types.Point{}
}

func (c *ArcToAbs) String() string {
	return fmt.Sprintf("A %.2f %.2f %.2f %d %d %.2f %.2f",
		c.RX, c.RY, c.XAxisRotation, boolToInt(c.LargeArc), boolToInt(c.Sweep), c.X, c.Y)
}

// ArcToRel 表示相对坐标的椭圆弧命令
type ArcToRel struct {
	RX, RY        float64 // 半径
	XAxisRotation float64 // x轴旋转角度（度）
	LargeArc      bool    // 大弧标志
	Sweep         bool    // 扫掠标志
	X, Y          float64 // 终点
	DX, DY        float64 // 相对坐标
}

func (c *ArcToRel) Execute(ctx *PathContext, precision float64) {
	startPoint := ctx.CurrentPoint
	endPoint := types.Point{
		X: startPoint.X + c.X,
		Y: startPoint.Y + c.Y,
	}

	// 转换为中心参数化
	rx, ry := math.Abs(c.RX), math.Abs(c.RY)
	xAxisRot := c.XAxisRotation * math.Pi / 180

	// 计算中间参数
	dx := (startPoint.X - endPoint.X) / 2
	dy := (startPoint.Y - endPoint.Y) / 2
	x1 := math.Cos(xAxisRot)*dx + math.Sin(xAxisRot)*dy
	y1 := -math.Sin(xAxisRot)*dx + math.Cos(xAxisRot)*dy

	// 校正半径
	lambda := (x1*x1)/(rx*rx) + (y1*y1)/(ry*ry)
	if lambda > 1 {
		rx *= math.Sqrt(lambda)
		ry *= math.Sqrt(lambda)
	}

	// 计算中心点
	sign := 1.0
	if c.LargeArc == c.Sweep {
		sign = -1.0
	}
	cx := sign * math.Sqrt((rx*rx*ry*ry-rx*rx*y1*y1-ry*ry*x1*x1)/(rx*rx*y1*y1+ry*ry*x1*x1))
	if math.IsNaN(cx) {
		cx = 0
	}
	cx = math.Cos(xAxisRot)*cx*rx/ry*y1 - math.Sin(xAxisRot)*cx*ry/rx*x1 + (startPoint.X+endPoint.X)/2
	cy := math.Sin(xAxisRot)*cx*rx/ry*y1 + math.Cos(xAxisRot)*cx*ry/rx*x1 + (startPoint.Y+endPoint.Y)/2

	// 计算起始角和终止角
	angle := func(u, v float64) float64 {
		return math.Atan2(v, u)
	}
	theta1 := angle((x1-cx)/rx, (y1-cy)/ry)
	theta2 := angle((-x1-cx)/rx, (-y1-cy)/ry)

	// 将椭圆弧转换为贝塞尔曲线段
	segments := int(math.Ceil(math.Abs(theta2-theta1) / (math.Pi / 2)))
	delta := (theta2 - theta1) / float64(segments)
	t := 8.0 / 3.0 * math.Sin(delta/4) * math.Sin(delta/4) / math.Sin(delta/2)

	for i := 0; i < segments; i++ {
		cos1 := math.Cos(theta1)
		sin1 := math.Sin(theta1)
		cos2 := math.Cos(theta1 + delta)
		sin2 := math.Sin(theta1 + delta)

		// 计算贝塞尔控制点
		point1 := types.Point{
			X: cx + rx*cos1 - (rx*sin1)*t,
			Y: cy + ry*sin1 + (ry*cos1)*t,
		}
		point2 := types.Point{
			X: cx + rx*cos2 + (rx*sin2)*t,
			Y: cy + ry*sin2 - (ry*cos2)*t,
		}
		point3 := types.Point{
			X: cx + rx*cos2,
			Y: cy + ry*sin2,
		}

		// 添加点到路径
		if i == 0 {
			ctx.Points = append(ctx.Points, point1)
		}
		ctx.Points = append(ctx.Points, point2, point3)
		theta1 += delta
	}

	ctx.CurrentPoint = endPoint
	ctx.PrevControl = types.Point{}
}

func (c *ArcToRel) String() string {
	return fmt.Sprintf("a %.2f %.2f %.2f %d %d %.2f %.2f",
		c.RX, c.RY, c.XAxisRotation, boolToInt(c.LargeArc), boolToInt(c.Sweep), c.X, c.Y)
}
