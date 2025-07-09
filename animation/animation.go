package animation

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
	"time"

	"github.com/hoonfeng/svg/attributes"
	"github.com/hoonfeng/svg/types"
)

// Animation 动画接口
// Animation interface defines the basic animation operations
type Animation interface {
	// Start 开始动画
	Start()
	// Pause 暂停动画
	Pause()
	// Resume 恢复动画
	Resume()
	// Stop 停止动画
	Stop()
	// Reset 重置动画
	Reset()
	// Update 更新动画状态
	Update(deltaTime float64)
	// Duration 返回动画持续时间
	Duration() float64
	// IsRunning 返回动画是否正在运行
	IsRunning() bool
}

// Easing 表示动画缓动函数
type Easing func(t float64) float64

// 预定义的缓动函数
var (
	// Linear 线性缓动
	Linear Easing = func(t float64) float64 {
		return t
	}

	// EaseInQuad 二次方缓入
	EaseInQuad Easing = func(t float64) float64 {
		return t * t
	}

	// EaseOutQuad 二次方缓出
	EaseOutQuad Easing = func(t float64) float64 {
		return t * (2 - t)
	}

	// EaseInOutQuad 二次方缓入缓出
	EaseInOutQuad Easing = func(t float64) float64 {
		if t < 0.5 {
			return 2 * t * t
		}
		return -1 + (4-2*t)*t
	}

	// EaseInCubic 三次方缓入
	EaseInCubic Easing = func(t float64) float64 {
		return t * t * t
	}

	// EaseOutCubic 三次方缓出
	EaseOutCubic Easing = func(t float64) float64 {
		return 1 + (t-1)*(t-1)*(t-1)
	}

	// EaseInOutCubic 三次方缓入缓出
	EaseInOutCubic Easing = func(t float64) float64 {
		if t < 0.5 {
			return 4 * t * t * t
		}
		return 1 + (t-1)*(2*t-2)*(2*t-2)
	}
)

// BaseAnimation 是所有动画的基础结构
type BaseAnimation struct {
	duration      float64 // 持续时间（秒）
	delay         float64 // 延迟时间（秒）
	currentTime   float64 // 当前时间（秒）
	isRunning     bool    // 是否正在运行
	isCompleted   bool    // 是否已完成
	easing        Easing  // 缓动函数
	onComplete    func()  // 完成回调
	repeatCount   int     // 重复次数（-1表示无限重复）
	currentRepeat int     // 当前重复次数
	autoReverse   bool    // 是否自动反向
	isReversed    bool    // 是否反向播放
}

// NewBaseAnimation 创建一个新的基础动画
func NewBaseAnimation(duration float64) *BaseAnimation {
	return &BaseAnimation{
		duration:      duration,
		delay:         0,
		currentTime:   0,
		isRunning:     false,
		isCompleted:   false,
		easing:        Linear,
		repeatCount:   0,
		currentRepeat: 0,
		autoReverse:   false,
		isReversed:    false,
	}
}

// Start 开始动画
func (a *BaseAnimation) Start() {
	a.isRunning = true
	a.isCompleted = false
	a.currentTime = 0
	a.currentRepeat = 0
	a.isReversed = false
}

// Pause 暂停动画
func (a *BaseAnimation) Pause() {
	a.isRunning = false
}

// Resume 恢复动画
func (a *BaseAnimation) Resume() {
	if !a.isCompleted {
		a.isRunning = true
	}
}

// Stop 停止动画
func (a *BaseAnimation) Stop() {
	a.isRunning = false
	a.isCompleted = true
}

// Reset 重置动画
func (a *BaseAnimation) Reset() {
	a.currentTime = 0
	a.isRunning = false
	a.isCompleted = false
	a.currentRepeat = 0
	a.isReversed = false
}

// Duration 返回动画持续时间
func (a *BaseAnimation) Duration() float64 {
	return a.duration
}

// SetDuration 设置动画持续时间
func (a *BaseAnimation) SetDuration(duration float64) {
	a.duration = duration
}

// IsRunning 返回动画是否正在运行
func (a *BaseAnimation) IsRunning() bool {
	return a.isRunning
}

// IsCompleted 返回动画是否已完成
func (a *BaseAnimation) IsCompleted() bool {
	return a.isCompleted
}

// SetEasing 设置缓动函数
func (a *BaseAnimation) SetEasing(easing Easing) {
	a.easing = easing
}

// SetDelay 设置延迟时间
func (a *BaseAnimation) SetDelay(delay float64) {
	a.delay = delay
}

// SetRepeatCount 设置重复次数
func (a *BaseAnimation) SetRepeatCount(count int) {
	a.repeatCount = count
}

// SetAutoReverse 设置是否自动反向
func (a *BaseAnimation) SetAutoReverse(autoReverse bool) {
	a.autoReverse = autoReverse
}

// OnComplete 设置动画完成回调
func (a *BaseAnimation) OnComplete(callback func()) {
	a.onComplete = callback
}

// Update 更新动画状态
func (a *BaseAnimation) Update(deltaTime float64) {
	if !a.isRunning || a.isCompleted {
		return
	}

	// 处理延迟
	if a.currentTime < a.delay {
		a.currentTime += deltaTime
		return
	}

	// 更新当前时间
	a.currentTime += deltaTime

	// 计算进度
	progress := (a.currentTime - a.delay) / a.duration

	// 检查是否完成一次循环
	if progress >= 1.0 {
		// 处理重复
		if a.repeatCount == -1 || a.currentRepeat < a.repeatCount {
			a.currentRepeat++
			a.currentTime = a.delay + float64(int64(a.currentTime-a.delay)%int64(a.duration))

			// 处理自动反向
			if a.autoReverse {
				a.isReversed = !a.isReversed
			}
		} else {
			// 动画完成
			a.isRunning = false
			a.isCompleted = true
			progress = 1.0

			// 调用完成回调
			if a.onComplete != nil {
				a.onComplete()
			}
		}
	}

	// 应用缓动函数
	easedProgress := a.easing(progress)

	// 如果是反向播放，反转进度
	if a.isReversed {
		easedProgress = 1.0 - easedProgress
	}

	// 应用动画效果（由子类实现）
	a.apply(easedProgress)
}

// apply 应用动画效果（由子类实现）
func (a *BaseAnimation) apply(progress float64) {
	// 空实现，由子类重写
}

// PropertyAnimation 属性动画
type PropertyAnimation struct {
	*BaseAnimation
	element   types.Element // 目标元素
	property  string        // 属性名
	fromValue string        // 起始值
	toValue   string        // 结束值
	valueType string        // 值类型（如"color", "length", "number"等）
}

// NewPropertyAnimation 创建一个新的属性动画
func NewPropertyAnimation(element types.Element, property, fromValue, toValue string, duration float64) *PropertyAnimation {
	return &PropertyAnimation{
		BaseAnimation: NewBaseAnimation(duration),
		element:       element,
		property:      property,
		fromValue:     fromValue,
		toValue:       toValue,
		valueType:     detectValueType(fromValue, toValue),
	}
}

// Update 更新属性动画
func (a *PropertyAnimation) Update(deltaTime float64) {
	if !a.isRunning || a.isCompleted {
		return
	}

	// 处理延迟
	if a.currentTime < a.delay {
		a.currentTime += deltaTime
		return
	}

	// 更新当前时间
	a.currentTime += deltaTime

	// 计算进度
	progress := (a.currentTime - a.delay) / a.duration

	// 检查是否完成一次循环
	if progress >= 1.0 {
		// 处理重复
		if a.repeatCount == -1 || a.currentRepeat < a.repeatCount {
			a.currentRepeat++
			a.currentTime = a.delay + float64(int64(a.currentTime-a.delay)%int64(a.duration))

			// 处理自动反向
			if a.autoReverse {
				a.isReversed = !a.isReversed
			}
		} else {
			// 动画完成
			a.isRunning = false
			a.isCompleted = true
			progress = 1.0

			// 调用完成回调
			if a.onComplete != nil {
				a.onComplete()
			}
		}
	}

	// 应用缓动函数
	easedProgress := a.easing(progress)

	// 如果是反向播放，反转进度
	if a.isReversed {
		easedProgress = 1.0 - easedProgress
	}

	// 应用动画效果
	a.apply(easedProgress)
}

// apply 应用属性动画
func (a *PropertyAnimation) apply(progress float64) {
	// 根据值类型插值
	var value string

	switch a.valueType {
	case "number":
		value = interpolateNumber(a.fromValue, a.toValue, progress)
	case "length":
		value = interpolateLength(a.fromValue, a.toValue, progress)
	case "color":
		value = interpolateColor(a.fromValue, a.toValue, progress)
	default:
		// 对于不支持插值的类型，在过程中间切换值
		if progress < 0.5 {
			value = a.fromValue
		} else {
			value = a.toValue
		}
	}

	// 设置元素属性
	a.element.SetAttribute(a.property, value)
}

// TransformAnimation 变换动画
type TransformAnimation struct {
	*BaseAnimation
	element       types.Element         // 目标元素
	fromTransform *attributes.Transform // 起始变换
	toTransform   *attributes.Transform // 结束变换
}

// NewTransformAnimation 创建一个新的变换动画
func NewTransformAnimation(element types.Element, fromTransform, toTransform *attributes.Transform, duration float64) *TransformAnimation {
	return &TransformAnimation{
		BaseAnimation: NewBaseAnimation(duration),
		element:       element,
		fromTransform: fromTransform,
		toTransform:   toTransform,
	}
}

// apply 应用变换动画
func (a *TransformAnimation) apply(progress float64) {
	// 插值变换矩阵的各个属性
	fromMatrix := a.fromTransform.GetMatrix()
	toMatrix := a.toTransform.GetMatrix()

	// 对矩阵的每个元素进行插值
	interpolatedMatrix := &attributes.Matrix{
		A: fromMatrix.A + (toMatrix.A-fromMatrix.A)*progress,
		B: fromMatrix.B + (toMatrix.B-fromMatrix.B)*progress,
		C: fromMatrix.C + (toMatrix.C-fromMatrix.C)*progress,
		D: fromMatrix.D + (toMatrix.D-fromMatrix.D)*progress,
		E: fromMatrix.E + (toMatrix.E-fromMatrix.E)*progress,
		F: fromMatrix.F + (toMatrix.F-fromMatrix.F)*progress,
	}

	// 创建新的Transform对象并设置矩阵变换
	interpolatedTransform := attributes.NewTransform()
	interpolatedTransform.Matrix(
		interpolatedMatrix.A, interpolatedMatrix.B,
		interpolatedMatrix.C, interpolatedMatrix.D,
		interpolatedMatrix.E, interpolatedMatrix.F,
	)

	// 将插值后的变换转换为SVG transform属性字符串
	transformStr := interpolatedTransform.ToString()

	// 设置元素的transform属性
	a.element.SetAttribute("transform", transformStr)
}

// KeyframeAnimation 关键帧动画
type KeyframeAnimation struct {
	*BaseAnimation
	element   types.Element      // 目标元素
	property  string             // 属性名
	keyframes map[float64]string // 关键帧（时间点 -> 值）
	valueType string             // 值类型
}

// NewKeyframeAnimation 创建一个新的关键帧动画
func NewKeyframeAnimation(element types.Element, property string, duration float64) *KeyframeAnimation {
	return &KeyframeAnimation{
		BaseAnimation: NewBaseAnimation(duration),
		element:       element,
		property:      property,
		keyframes:     make(map[float64]string),
		valueType:     "unknown",
	}
}

// AddKeyframe 添加关键帧
func (a *KeyframeAnimation) AddKeyframe(time float64, value string) {
	a.keyframes[time] = value

	// 更新值类型
	if a.valueType == "unknown" && len(a.keyframes) > 0 {
		a.valueType = detectValueType(value, value)
	}
}

// apply 应用关键帧动画
func (a *KeyframeAnimation) apply(progress float64) {
	// 找到当前进度对应的关键帧
	var prevTime float64 = 0
	var nextTime float64 = 1
	var prevValue, nextValue string

	for time, value := range a.keyframes {
		if time <= progress && time > prevTime {
			prevTime = time
			prevValue = value
		}
		if time > progress && time < nextTime {
			nextTime = time
			nextValue = value
		}
	}

	// 如果没有找到合适的关键帧，使用第一个和最后一个
	if prevValue == "" {
		for time, value := range a.keyframes {
			if time < prevTime || prevValue == "" {
				prevTime = time
				prevValue = value
			}
		}
	}

	if nextValue == "" {
		for time, value := range a.keyframes {
			if time > nextTime || nextValue == "" {
				nextTime = time
				nextValue = value
			}
		}
	}

	// 计算关键帧之间的插值
	var value string

	if nextTime == prevTime {
		value = prevValue
	} else {
		segmentProgress := (progress - prevTime) / (nextTime - prevTime)

		switch a.valueType {
		case "number":
			value = interpolateNumber(prevValue, nextValue, segmentProgress)
		case "length":
			value = interpolateLength(prevValue, nextValue, segmentProgress)
		case "color":
			value = interpolateColor(prevValue, nextValue, segmentProgress)
		default:
			value = prevValue
		}
	}

	// 设置元素属性
	a.element.SetAttribute(a.property, value)
}

// AnimationGroup 动画组
type AnimationGroup struct {
	*BaseAnimation
	animations []Animation // 子动画
}

// NewAnimationGroup 创建一个新的动画组
func NewAnimationGroup() *AnimationGroup {
	return &AnimationGroup{
		BaseAnimation: NewBaseAnimation(0),
		animations:    make([]Animation, 0),
	}
}

// AddAnimation 添加子动画
func (g *AnimationGroup) AddAnimation(animation Animation) {
	g.animations = append(g.animations, animation)

	// 更新持续时间
	if animation.Duration() > g.duration {
		g.duration = animation.Duration()
	}
}

// Start 开始所有子动画
func (g *AnimationGroup) Start() {
	g.BaseAnimation.Start()

	for _, animation := range g.animations {
		animation.Start()
	}
}

// Pause 暂停所有子动画
func (g *AnimationGroup) Pause() {
	g.BaseAnimation.Pause()

	for _, animation := range g.animations {
		animation.Pause()
	}
}

// Resume 恢复所有子动画
func (g *AnimationGroup) Resume() {
	g.BaseAnimation.Resume()

	for _, animation := range g.animations {
		animation.Resume()
	}
}

// Stop 停止所有子动画
func (g *AnimationGroup) Stop() {
	g.BaseAnimation.Stop()

	for _, animation := range g.animations {
		animation.Stop()
	}
}

// Reset 重置所有子动画
func (g *AnimationGroup) Reset() {
	g.BaseAnimation.Reset()

	for _, animation := range g.animations {
		animation.Reset()
	}
}

// Update 更新所有子动画
func (g *AnimationGroup) Update(deltaTime float64) {
	g.BaseAnimation.Update(deltaTime)

	if !g.isRunning {
		return
	}

	allCompleted := true

	for _, animation := range g.animations {
		animation.Update(deltaTime)
		if animation.IsRunning() {
			allCompleted = false
		}
	}

	// 如果所有子动画都完成了，标记组动画为完成
	if allCompleted && !g.isCompleted {
		g.isCompleted = true
		g.isRunning = false

		if g.onComplete != nil {
			g.onComplete()
		}
	}
}

// SequentialAnimationGroup 顺序动画组
type SequentialAnimationGroup struct {
	*AnimationGroup
	currentIndex int // 当前正在播放的动画索引
}

// NewSequentialAnimationGroup 创建一个新的顺序动画组
func NewSequentialAnimationGroup() *SequentialAnimationGroup {
	return &SequentialAnimationGroup{
		AnimationGroup: NewAnimationGroup(),
		currentIndex:   0,
	}
}

// Start 开始第一个子动画
func (g *SequentialAnimationGroup) Start() {
	g.BaseAnimation.Start()
	g.currentIndex = 0

	if len(g.animations) > 0 {
		g.animations[0].Start()
	} else {
		g.isCompleted = true
		g.isRunning = false
	}
}

// Update 更新当前子动画
func (g *SequentialAnimationGroup) Update(deltaTime float64) {
	g.BaseAnimation.Update(deltaTime)

	if !g.isRunning || g.isCompleted {
		return
	}

	// 更新当前动画
	if g.currentIndex < len(g.animations) {
		currentAnim := g.animations[g.currentIndex]
		currentAnim.Update(deltaTime)

		// 如果当前动画完成，移动到下一个
		if !currentAnim.IsRunning() {
			g.currentIndex++

			// 如果还有下一个动画，开始它
			if g.currentIndex < len(g.animations) {
				g.animations[g.currentIndex].Start()
			} else {
				// 所有动画都完成了
				g.isCompleted = true
				g.isRunning = false

				if g.onComplete != nil {
					g.onComplete()
				}
			}
		}
	}
}

// AnimationManager 动画管理器
type AnimationManager struct {
	animations []Animation
	lastTime   time.Time
	isRunning  bool
}

// NewAnimationManager 创建一个新的动画管理器
func NewAnimationManager() *AnimationManager {
	return &AnimationManager{
		animations: make([]Animation, 0),
		lastTime:   time.Now(),
		isRunning:  false,
	}
}

// AddAnimation 添加动画
func (m *AnimationManager) AddAnimation(animation Animation) {
	m.animations = append(m.animations, animation)
}

// RemoveAnimation 移除动画
func (m *AnimationManager) RemoveAnimation(animation Animation) {
	for i, anim := range m.animations {
		if anim == animation {
			m.animations = append(m.animations[:i], m.animations[i+1:]...)
			return
		}
	}
}

// Start 开始所有动画
func (m *AnimationManager) Start() {
	m.isRunning = true
	m.lastTime = time.Now()

	for _, animation := range m.animations {
		animation.Start()
	}
}

// Stop 停止所有动画
func (m *AnimationManager) Stop() {
	m.isRunning = false

	for _, animation := range m.animations {
		animation.Stop()
	}
}

// Update 更新所有动画
func (m *AnimationManager) Update() {
	if !m.isRunning {
		return
	}

	now := time.Now()
	deltaTime := now.Sub(m.lastTime).Seconds()
	m.lastTime = now

	// 更新所有动画
	for i := 0; i < len(m.animations); i++ {
		animation := m.animations[i]
		animation.Update(deltaTime)

		// 如果动画已完成且不是无限重复，移除它
		if !animation.IsRunning() {
			m.animations = append(m.animations[:i], m.animations[i+1:]...)
			i--
		}
	}
}

// 辅助函数

// detectValueType 检测值类型
func detectValueType(fromValue, toValue string) string {
	// 检查是否是颜色
	if (strings.HasPrefix(fromValue, "#") || strings.HasPrefix(fromValue, "rgb")) &&
		(strings.HasPrefix(toValue, "#") || strings.HasPrefix(toValue, "rgb")) {
		return "color"
	}

	// 检查是否是带单位的长度
	if containsLengthUnit(fromValue) && containsLengthUnit(toValue) {
		return "length"
	}

	// 检查是否是数字
	_, err1 := strconv.ParseFloat(fromValue, 64)
	_, err2 := strconv.ParseFloat(toValue, 64)
	if err1 == nil && err2 == nil {
		return "number"
	}

	return "string"
}

// containsLengthUnit 检查字符串是否包含长度单位
func containsLengthUnit(s string) bool {
	units := []string{"px", "pt", "pc", "mm", "cm", "in", "%", "em", "ex", "rem"}
	for _, unit := range units {
		if strings.HasSuffix(s, unit) {
			// 检查前面的部分是否是数字
			numPart := s[:len(s)-len(unit)]
			_, err := strconv.ParseFloat(numPart, 64)
			return err == nil
		}
	}
	return false
}

// interpolateNumber 插值数字
func interpolateNumber(from, to string, progress float64) string {
	fromVal, _ := strconv.ParseFloat(from, 64)
	toVal, _ := strconv.ParseFloat(to, 64)

	result := fromVal + (toVal-fromVal)*progress
	return fmt.Sprintf("%g", result)
}

// interpolateLength 插值长度
func interpolateLength(from, to string, progress float64) string {
	// 提取数值和单位
	var fromVal, toVal float64
	var fromUnit, toUnit string

	for _, unit := range []string{"px", "pt", "pc", "mm", "cm", "in", "%", "em", "ex", "rem"} {
		if strings.HasSuffix(from, unit) {
			fromUnit = unit
			fromVal, _ = strconv.ParseFloat(from[:len(from)-len(unit)], 64)
			break
		}
	}

	for _, unit := range []string{"px", "pt", "pc", "mm", "cm", "in", "%", "em", "ex", "rem"} {
		if strings.HasSuffix(to, unit) {
			toUnit = unit
			toVal, _ = strconv.ParseFloat(to[:len(to)-len(unit)], 64)
			break
		}
	}

	// 如果单位不同，使用第一个单位
	unit := fromUnit
	if unit == "" {
		unit = toUnit
	}

	// 插值数值
	result := fromVal + (toVal-fromVal)*progress
	return fmt.Sprintf("%g%s", result, unit)
}

// interpolateColor 插值颜色
func interpolateColor(from, to string, progress float64) string {
	// 解析颜色
	fromColor, _ := attributes.ParseColor(from)
	toColor, _ := attributes.ParseColor(to)

	// 提取RGBA分量
	fr, fg, fb, fa := fromColor.RGBA()
	tr, tg, tb, ta := toColor.RGBA()

	// 将16位值转换为8位
	fr, fg, fb, fa = fr>>8, fg>>8, fb>>8, fa>>8
	tr, tg, tb, ta = tr>>8, tg>>8, tb>>8, ta>>8

	// 插值
	r := uint8(float64(fr) + float64(tr-fr)*progress)
	g := uint8(float64(fg) + float64(tg-fg)*progress)
	b := uint8(float64(fb) + float64(tb-fb)*progress)
	a := uint8(float64(fa) + float64(ta-fa)*progress)

	// 创建结果颜色
	result := attributes.ColorToHex(color.RGBA{r, g, b, a})
	return result
}
