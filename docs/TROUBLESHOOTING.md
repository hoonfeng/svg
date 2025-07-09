# SVG库故障排除指南 / Troubleshooting Guide

## 📖 概述 / Overview

本指南提供了使用SVG库时可能遇到的常见问题、错误信息和解决方案。按照问题类型分类，便于快速定位和解决问题。

This guide provides common issues, error messages, and solutions you may encounter when using the SVG library. Organized by problem type for quick identification and resolution.

## 🚨 常见错误 / Common Errors

### 1. 编译错误 / Compilation Errors

#### 错误：找不到包 / Error: Package Not Found

```
error: cannot find package "github.com/yourproject/svg"
```

**原因 / Cause:**
- Go模块路径配置错误 / Incorrect Go module path configuration
- 依赖未正确安装 / Dependencies not properly installed

**解决方案 / Solution:**
```bash
# 1. 初始化Go模块
# 1. Initialize Go module
go mod init your-project-name

# 2. 添加依赖
# 2. Add dependency
go mod tidy

# 3. 如果是本地开发，使用replace指令
# 3. For local development, use replace directive
echo 'replace github.com/yourproject/svg => ./svg' >> go.mod
```

#### 错误：类型不匹配 / Error: Type Mismatch

```
error: cannot use color.RGBA literal as color.Color value
```

**原因 / Cause:**
- 颜色类型转换问题 / Color type conversion issue

**解决方案 / Solution:**
```go
// ❌ 错误写法 / Wrong way
canvas.Rect(0, 0, 100, 100).Fill(color.RGBA{255, 0, 0, 255})

// ✅ 正确写法 / Correct way
canvas.Rect(0, 0, 100, 100).Fill("#FF0000")
// 或者 / Or
canvas.Rect(0, 0, 100, 100).Fill("red")
// 或者 / Or
canvas.Rect(0, 0, 100, 100).Fill("rgb(255, 0, 0)")
```

### 2. 运行时错误 / Runtime Errors

#### 错误：空指针异常 / Error: Nil Pointer Exception

```
panic: runtime error: invalid memory address or nil pointer dereference
```

**常见原因和解决方案 / Common Causes and Solutions:**

```go
// 原因1：未初始化SVG画布
// Cause 1: SVG canvas not initialized
// ❌ 错误
var canvas *SVG
canvas.Rect(0, 0, 100, 100) // panic!

// ✅ 正确
canvas := svg.New(800, 600)
canvas.Rect(0, 0, 100, 100)

// 原因2：元素为nil
// Cause 2: Element is nil
// ❌ 错误
var rect *RectElement
rect.Fill("red") // panic!

// ✅ 正确
rect := canvas.Rect(0, 0, 100, 100)
rect.Fill("red")

// 原因3：链式调用中断
// Cause 3: Method chaining interrupted
// ❌ 错误
canvas.Rect(0, 0, 100, 100).
    Fill("invalid-color"). // 返回nil
    Stroke("black")       // panic!

// ✅ 正确 - 分步验证
rect := canvas.Rect(0, 0, 100, 100)
if rect != nil {
    rect.Fill("red")
    rect.Stroke("black")
}
```

#### 错误：文件操作失败 / Error: File Operation Failed

```
error: open output.svg: permission denied
```

**解决方案 / Solution:**
```go
// 检查文件权限和路径
// Check file permissions and path
func SafeSaveFile(canvas *SVG, filename string) error {
    // 1. 检查目录是否存在
    // 1. Check if directory exists
    dir := filepath.Dir(filename)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return fmt.Errorf("failed to create directory: %w", err)
    }
    
    // 2. 检查文件是否可写
    // 2. Check if file is writable
    if _, err := os.Stat(filename); err == nil {
        // 文件存在，检查权限
        // File exists, check permissions
        file, err := os.OpenFile(filename, os.O_WRONLY, 0)
        if err != nil {
            return fmt.Errorf("file not writable: %w", err)
        }
        file.Close()
    }
    
    // 3. 保存文件
    // 3. Save file
    return canvas.SaveSVG(filename)
}
```

### 3. 渲染问题 / Rendering Issues

#### 问题：SVG显示空白 / Issue: SVG Shows Blank

**可能原因和解决方案 / Possible Causes and Solutions:**

```go
// 原因1：元素超出画布范围
// Cause 1: Elements outside canvas bounds
// ❌ 问题代码
canvas := svg.New(100, 100)
canvas.Rect(200, 200, 50, 50) // 超出画布范围

// ✅ 解决方案
canvas := svg.New(300, 300)
canvas.Rect(200, 200, 50, 50) // 在画布范围内

// 原因2：元素颜色与背景相同
// Cause 2: Element color same as background
// ❌ 问题代码
canvas := svg.New(100, 100)
canvas.SetBackground(color.RGBA{255, 255, 255, 255}) // 白色背景
canvas.Rect(10, 10, 50, 50).Fill("white") // 白色矩形，看不见

// ✅ 解决方案
canvas.Rect(10, 10, 50, 50).Fill("red").Stroke("black")

// 原因3：元素尺寸为0
// Cause 3: Element size is 0
// ❌ 问题代码
canvas.Rect(10, 10, 0, 0) // 尺寸为0

// ✅ 解决方案
canvas.Rect(10, 10, 50, 50) // 正确尺寸
```

#### 问题：文本不显示 / Issue: Text Not Showing

```go
// 调试文本显示问题
// Debug text display issues
func DebugTextRendering(canvas *SVG) {
    // 1. 检查文本位置
    // 1. Check text position
    text := canvas.Text(50, 50, "Hello World")
    
    // 2. 设置明显的样式
    // 2. Set obvious styles
    text.Fill("red").
         FontSize(20).
         FontFamily("Arial").
         FontWeight("bold")
    
    // 3. 添加背景矩形用于定位
    // 3. Add background rectangle for positioning
    canvas.Rect(45, 30, 100, 30).
           Fill("yellow").
           Stroke("black")
    
    // 4. 检查文本边界
    // 4. Check text bounds
    bounds := text.GetBounds()
    fmt.Printf("Text bounds: x=%.2f, y=%.2f, w=%.2f, h=%.2f\n",
        bounds.X, bounds.Y, bounds.Width, bounds.Height)
}
```

## 🎨 样式问题 / Style Issues

### 1. 颜色问题 / Color Issues

#### 问题：颜色不生效 / Issue: Colors Not Working

```go
// 颜色格式验证
// Color format validation
func ValidateAndFixColors() {
    // ❌ 常见错误格式
    // ❌ Common wrong formats
    badColors := []string{
        "FF0000",      // 缺少#号
        "#GGHHII",     // 无效十六进制
        "rgb(300,0,0)", // 值超出范围
        "rgba(255,0,0)", // 缺少alpha值
    }
    
    // ✅ 正确格式
    // ✅ Correct formats
    goodColors := []string{
        "#FF0000",
        "#F00",
        "rgb(255, 0, 0)",
        "rgba(255, 0, 0, 1.0)",
        "red",
        "transparent",
    }
    
    // 验证颜色
    // Validate colors
    for _, colorStr := range goodColors {
        if _, err := ParseColor(colorStr); err != nil {
            fmt.Printf("Invalid color: %s, error: %v\n", colorStr, err)
        } else {
            fmt.Printf("Valid color: %s\n", colorStr)
        }
    }
}
```

#### 问题：透明度不正确 / Issue: Incorrect Transparency

```go
// 透明度处理
// Transparency handling
func HandleTransparency() {
    canvas := svg.New(200, 200)
    
    // 方法1：使用RGBA
    // Method 1: Use RGBA
    canvas.Rect(10, 10, 50, 50).Fill("rgba(255, 0, 0, 0.5)") // 50%透明
    
    // 方法2：使用opacity属性
    // Method 2: Use opacity attribute
    canvas.Rect(70, 10, 50, 50).
           Fill("red").
           SetAttribute("opacity", "0.5")
    
    // 方法3：使用fill-opacity
    // Method 3: Use fill-opacity
    canvas.Rect(130, 10, 50, 50).
           Fill("red").
           SetAttribute("fill-opacity", "0.5")
}
```

### 2. 字体问题 / Font Issues

#### 问题：字体不显示或显示错误 / Issue: Font Not Displaying or Wrong Font

```go
// 字体回退机制
// Font fallback mechanism
func SetupFontFallback(text *TextElement) {
    // 设置字体族回退
    // Set font family fallback
    text.FontFamily("'Custom Font', Arial, sans-serif")
    
    // 检查系统字体
    // Check system fonts
    systemFonts := GetSystemFonts()
    if len(systemFonts) > 0 {
        text.FontFamily(systemFonts[0])
    }
}

// 获取系统字体列表
// Get system font list
func GetSystemFonts() []string {
    switch runtime.GOOS {
    case "windows":
        return []string{"Segoe UI", "Arial", "Tahoma"}
    case "darwin":
        return []string{"San Francisco", "Helvetica", "Arial"}
    case "linux":
        return []string{"Ubuntu", "DejaVu Sans", "Liberation Sans"}
    default:
        return []string{"Arial", "sans-serif"}
    }
}
```

## 🔄 动画问题 / Animation Issues

### 1. GIF生成问题 / GIF Generation Issues

#### 问题：GIF文件损坏或无法播放 / Issue: GIF File Corrupted or Won't Play

```go
// GIF生成调试
// GIF generation debugging
func DebugGIFGeneration() {
    builder := NewAnimationBuilder(400, 300)
    
    // 1. 检查帧数设置
    // 1. Check frame count settings
    frameCount := 30
    if frameCount < 2 {
        log.Fatal("Frame count must be at least 2")
    }
    
    // 2. 检查帧率设置
    // 2. Check frame rate settings
    frameRate := 10
    if frameRate < 1 || frameRate > 50 {
        log.Fatal("Frame rate must be between 1 and 50")
    }
    
    // 3. 设置合理的参数
    // 3. Set reasonable parameters
    builder.SetFrameCount(frameCount).
            SetFrameRate(frameRate)
    
    // 4. 检查动画配置
    // 4. Check animation configuration
    config := AnimationConfig{
        Duration:   2.0, // 2秒
        Easing:     EaseInOut,
        Background: color.RGBA{255, 255, 255, 255},
        Loop:       true,
    }
    
    // 5. 生成动画
    // 5. Generate animation
    builder.CreateRotatingShapes(config)
    
    // 6. 保存并检查文件
    // 6. Save and check file
    filename := "debug_animation.gif"
    if err := builder.SaveToGIF(filename); err != nil {
        log.Fatalf("Failed to save GIF: %v", err)
    }
    
    // 7. 验证文件大小
    // 7. Verify file size
    if stat, err := os.Stat(filename); err == nil {
        if stat.Size() == 0 {
            log.Fatal("Generated GIF file is empty")
        }
        fmt.Printf("GIF file size: %d bytes\n", stat.Size())
    }
}
```

#### 问题：动画播放速度不正确 / Issue: Animation Speed Incorrect

```go
// 动画速度调试
// Animation speed debugging
func FixAnimationSpeed() {
    // 计算正确的帧率和持续时间关系
    // Calculate correct frame rate and duration relationship
    duration := 3.0    // 3秒
    frameCount := 60   // 60帧
    frameRate := int(float64(frameCount) / duration) // 20 FPS
    
    fmt.Printf("Duration: %.1fs, Frames: %d, FPS: %d\n", 
        duration, frameCount, frameRate)
    
    builder := NewAnimationBuilder(400, 300)
    builder.SetFrameCount(frameCount).
            SetFrameRate(frameRate)
    
    config := AnimationConfig{
        Duration: duration,
        Easing:   Linear, // 使用线性缓动便于测试
    }
    
    builder.CreateRotatingShapes(config)
    builder.SaveToGIF("speed_test.gif")
}
```

### 2. 缓动函数问题 / Easing Function Issues

```go
// 自定义缓动函数调试
// Custom easing function debugging
func TestEasingFunctions() {
    // 测试缓动函数的输入输出
    // Test easing function input/output
    easingFuncs := map[string]EasingFunc{
        "Linear":        Linear,
        "EaseIn":        EaseIn,
        "EaseOut":       EaseOut,
        "EaseInOut":     EaseInOut,
        "EaseInOutQuad": EaseInOutQuad,
    }
    
    for name, easing := range easingFuncs {
        fmt.Printf("\n%s easing function test:\n", name)
        
        // 测试关键点
        // Test key points
        testPoints := []float64{0.0, 0.25, 0.5, 0.75, 1.0}
        
        for _, t := range testPoints {
            result := easing(t)
            fmt.Printf("  t=%.2f -> %.3f\n", t, result)
            
            // 验证输出范围
            // Validate output range
            if result < 0 || result > 1 {
                fmt.Printf("  WARNING: Output out of range [0,1]\n")
            }
        }
    }
}
```

## 💾 文件操作问题 / File Operation Issues

### 1. 保存问题 / Save Issues

#### 问题：文件保存失败 / Issue: File Save Failed

```go
// 安全的文件保存
// Safe file saving
func SafeFileSave(canvas *SVG, filename string) error {
    // 1. 验证文件名
    // 1. Validate filename
    if filename == "" {
        return errors.New("filename cannot be empty")
    }
    
    // 2. 检查文件扩展名
    // 2. Check file extension
    ext := strings.ToLower(filepath.Ext(filename))
    validExts := []string{".svg", ".png", ".jpg", ".jpeg"}
    
    isValid := false
    for _, validExt := range validExts {
        if ext == validExt {
            isValid = true
            break
        }
    }
    
    if !isValid {
        return fmt.Errorf("unsupported file extension: %s", ext)
    }
    
    // 3. 创建目录
    // 3. Create directory
    dir := filepath.Dir(filename)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return fmt.Errorf("failed to create directory: %w", err)
    }
    
    // 4. 检查磁盘空间
    // 4. Check disk space
    if err := checkDiskSpace(dir); err != nil {
        return fmt.Errorf("insufficient disk space: %w", err)
    }
    
    // 5. 保存文件
    // 5. Save file
    switch ext {
    case ".svg":
        return canvas.SaveSVG(filename)
    case ".png":
        return canvas.SavePNG(filename)
    case ".jpg", ".jpeg":
        return canvas.SaveJPEG(filename, 90) // 90%质量
    default:
        return fmt.Errorf("unsupported format: %s", ext)
    }
}

// 检查磁盘空间
// Check disk space
func checkDiskSpace(dir string) error {
    // 简化的磁盘空间检查
    // Simplified disk space check
    const minFreeSpace = 10 * 1024 * 1024 // 10MB
    
    // 这里应该实现实际的磁盘空间检查
    // Actual disk space check should be implemented here
    // 为了示例，我们假设有足够空间
    // For example purposes, we assume sufficient space
    
    return nil
}
```

### 2. 加载问题 / Load Issues

#### 问题：SVG文件解析失败 / Issue: SVG File Parse Failed

```go
// 安全的SVG加载
// Safe SVG loading
func SafeLoadSVG(filename string) (*SVG, error) {
    // 1. 检查文件是否存在
    // 1. Check if file exists
    if _, err := os.Stat(filename); os.IsNotExist(err) {
        return nil, fmt.Errorf("file does not exist: %s", filename)
    }
    
    // 2. 检查文件大小
    // 2. Check file size
    stat, err := os.Stat(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to get file info: %w", err)
    }
    
    const maxFileSize = 50 * 1024 * 1024 // 50MB
    if stat.Size() > maxFileSize {
        return nil, fmt.Errorf("file too large: %d bytes (max: %d)", 
            stat.Size(), maxFileSize)
    }
    
    // 3. 读取文件内容
    // 3. Read file content
    content, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to read file: %w", err)
    }
    
    // 4. 验证文件格式
    // 4. Validate file format
    if !strings.Contains(string(content), "<svg") {
        return nil, errors.New("file does not appear to be a valid SVG")
    }
    
    // 5. 解析SVG
    // 5. Parse SVG
    svg, err := ParseSVG(string(content))
    if err != nil {
        return nil, fmt.Errorf("failed to parse SVG: %w", err)
    }
    
    return svg, nil
}
```

## 🔧 性能问题 / Performance Issues

### 1. 内存使用过高 / High Memory Usage

```go
// 内存使用监控
// Memory usage monitoring
func MonitorMemoryUsage() {
    var m runtime.MemStats
    
    // 获取初始内存状态
    // Get initial memory state
    runtime.ReadMemStats(&m)
    initialAlloc := m.Alloc
    
    fmt.Printf("Initial memory: %d KB\n", initialAlloc/1024)
    
    // 创建大量元素
    // Create many elements
    canvas := svg.New(1000, 1000)
    
    for i := 0; i < 10000; i++ {
        canvas.Rect(float64(i%1000), float64(i/1000), 1, 1).Fill("red")
        
        // 每1000个元素检查一次内存
        // Check memory every 1000 elements
        if i%1000 == 0 {
            runtime.ReadMemStats(&m)
            currentAlloc := m.Alloc
            fmt.Printf("Elements: %d, Memory: %d KB (+%d KB)\n", 
                i, currentAlloc/1024, (currentAlloc-initialAlloc)/1024)
        }
    }
    
    // 强制垃圾回收
    // Force garbage collection
    runtime.GC()
    runtime.ReadMemStats(&m)
    finalAlloc := m.Alloc
    
    fmt.Printf("Final memory after GC: %d KB\n", finalAlloc/1024)
}

// 内存优化建议
// Memory optimization suggestions
func OptimizeMemoryUsage() {
    // 1. 使用对象池
    // 1. Use object pools
    pool := &ElementPool{}
    
    // 2. 批量处理
    // 2. Batch processing
    const batchSize = 1000
    
    // 3. 及时清理不需要的引用
    // 3. Clean up unnecessary references promptly
    canvas := svg.New(1000, 1000)
    
    for batch := 0; batch < 10; batch++ {
        elements := make([]*RectElement, batchSize)
        
        // 创建元素
        // Create elements
        for i := 0; i < batchSize; i++ {
            elements[i] = canvas.Rect(float64(i), 0, 1, 1)
        }
        
        // 处理元素
        // Process elements
        for _, element := range elements {
            element.Fill("red")
        }
        
        // 清理引用
        // Clear references
        for i := range elements {
            elements[i] = nil
        }
        elements = nil
        
        // 定期垃圾回收
        // Periodic garbage collection
        if batch%5 == 0 {
            runtime.GC()
        }
    }
}
```

### 2. 渲染速度慢 / Slow Rendering Speed

```go
// 渲染性能优化
// Rendering performance optimization
func OptimizeRenderingSpeed() {
    canvas := svg.New(1000, 1000)
    
    start := time.Now()
    
    // 方法1：使用分组减少DOM操作
    // Method 1: Use groups to reduce DOM operations
    group := canvas.Group()
    
    for i := 0; i < 1000; i++ {
        rect := canvas.Rect(float64(i%100)*10, float64(i/100)*10, 8, 8)
        rect.Fill(fmt.Sprintf("hsl(%d, 70%%, 50%%)", i%360))
        group.Add(rect)
    }
    
    groupTime := time.Since(start)
    fmt.Printf("Group method: %v\n", groupTime)
    
    // 方法2：预分配切片
    // Method 2: Pre-allocate slices
    start = time.Now()
    
    canvas2 := svg.New(1000, 1000)
    elements := make([]Element, 0, 1000) // 预分配容量
    
    for i := 0; i < 1000; i++ {
        rect := canvas2.Rect(float64(i%100)*10, float64(i/100)*10, 8, 8)
        rect.Fill(fmt.Sprintf("hsl(%d, 70%%, 50%%)", i%360))
        elements = append(elements, rect)
    }
    
    preAllocTime := time.Since(start)
    fmt.Printf("Pre-allocation method: %v\n", preAllocTime)
    
    // 方法3：批量样式设置
    // Method 3: Batch style setting
    start = time.Now()
    
    canvas3 := svg.New(1000, 1000)
    
    // 创建样式类
    // Create style classes
    styleGroup := canvas3.Group().Class("batch-style")
    
    for i := 0; i < 1000; i++ {
        rect := canvas3.Rect(float64(i%100)*10, float64(i/100)*10, 8, 8)
        rect.Class(fmt.Sprintf("color-%d", i%10))
        styleGroup.Add(rect)
    }
    
    batchTime := time.Since(start)
    fmt.Printf("Batch styling method: %v\n", batchTime)
}
```

## 🔍 调试技巧 / Debugging Tips

### 1. 启用调试模式 / Enable Debug Mode

```go
// 调试配置
// Debug configuration
type DebugConfig struct {
    ShowBounds     bool
    ShowGrid       bool
    LogOperations  bool
    ValidateInputs bool
    ProfileMemory  bool
}

// 调试渲染器
// Debug renderer
func CreateDebugRenderer(config DebugConfig) *DebugRenderer {
    return &DebugRenderer{
        config: config,
        logger: log.New(os.Stdout, "[SVG-DEBUG] ", log.LstdFlags),
    }
}

type DebugRenderer struct {
    config DebugConfig
    logger *log.Logger
}

func (dr *DebugRenderer) RenderElement(canvas *SVG, element Element) {
    if dr.config.LogOperations {
        dr.logger.Printf("Rendering element: %T", element)
    }
    
    if dr.config.ValidateInputs {
        if err := dr.validateElement(element); err != nil {
            dr.logger.Printf("Validation error: %v", err)
            return
        }
    }
    
    if dr.config.ShowBounds {
        dr.addBoundsVisualization(canvas, element)
    }
    
    canvas.Add(element)
    
    if dr.config.ProfileMemory {
        var m runtime.MemStats
        runtime.ReadMemStats(&m)
        dr.logger.Printf("Memory after rendering: %d KB", m.Alloc/1024)
    }
}

func (dr *DebugRenderer) validateElement(element Element) error {
    bounds := element.GetBounds()
    
    if bounds.Width <= 0 || bounds.Height <= 0 {
        return fmt.Errorf("invalid element dimensions: %.2fx%.2f", 
            bounds.Width, bounds.Height)
    }
    
    if bounds.X < -10000 || bounds.Y < -10000 || 
       bounds.X > 10000 || bounds.Y > 10000 {
        return fmt.Errorf("element position out of reasonable range: (%.2f, %.2f)", 
            bounds.X, bounds.Y)
    }
    
    return nil
}
```

### 2. 性能分析 / Performance Profiling

```go
// 性能分析工具
// Performance profiling tools
func ProfileSVGGeneration() {
    // 启用CPU分析
    // Enable CPU profiling
    cpuFile, err := os.Create("cpu.prof")
    if err != nil {
        log.Fatal(err)
    }
    defer cpuFile.Close()
    
    pprof.StartCPUProfile(cpuFile)
    defer pprof.StopCPUProfile()
    
    // 启用内存分析
    // Enable memory profiling
    defer func() {
        memFile, err := os.Create("mem.prof")
        if err != nil {
            log.Fatal(err)
        }
        defer memFile.Close()
        
        runtime.GC()
        pprof.WriteHeapProfile(memFile)
    }()
    
    // 执行性能测试
    // Execute performance test
    canvas := svg.New(1000, 1000)
    
    for i := 0; i < 10000; i++ {
        canvas.Rect(float64(i%1000), float64(i/1000), 1, 1).
               Fill(fmt.Sprintf("hsl(%d, 50%%, 50%%)", i%360))
    }
    
    canvas.SaveSVG("performance_test.svg")
}

// 使用方法：
// Usage:
// go run main.go
// go tool pprof cpu.prof
// go tool pprof mem.prof
```

## 📞 获取帮助 / Getting Help

### 1. 错误报告 / Error Reporting

当遇到问题时，请提供以下信息：

When reporting issues, please provide the following information:

```go
// 系统信息收集
// System information collection
func CollectSystemInfo() {
    fmt.Println("=== System Information ===")
    fmt.Printf("Go Version: %s\n", runtime.Version())
    fmt.Printf("OS: %s\n", runtime.GOOS)
    fmt.Printf("Architecture: %s\n", runtime.GOARCH)
    fmt.Printf("CPU Cores: %d\n", runtime.NumCPU())
    
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    fmt.Printf("Memory: %d MB\n", m.Sys/1024/1024)
    
    fmt.Println("\n=== SVG Library Information ===")
    fmt.Printf("Library Version: %s\n", GetVersion())
    
    fmt.Println("\n=== Error Details ===")
    // 在这里添加具体的错误信息
    // Add specific error information here
}
```

### 2. 常用调试命令 / Common Debug Commands

```bash
# 检查Go环境
# Check Go environment
go version
go env

# 运行测试
# Run tests
go test -v ./...

# 运行基准测试
# Run benchmarks
go test -bench=. -benchmem

# 检查代码质量
# Check code quality
go vet ./...
golint ./...

# 查看依赖
# View dependencies
go mod graph
go mod why -m github.com/yourproject/svg
```

### 3. 社区资源 / Community Resources

- **文档**: 查看完整的API文档和教程
- **示例**: 参考示例代码库
- **论坛**: 在开发者社区提问
- **GitHub**: 提交Issue和Pull Request

---

## 📋 快速检查清单 / Quick Checklist

遇到问题时，请按以下顺序检查：

When encountering issues, check in the following order:

- [ ] **环境检查** / **Environment Check**
  - [ ] Go版本是否兼容 / Go version compatibility
  - [ ] 依赖是否正确安装 / Dependencies properly installed
  - [ ] 模块路径是否正确 / Module path correct

- [ ] **代码检查** / **Code Check**
  - [ ] 变量是否正确初始化 / Variables properly initialized
  - [ ] 参数是否在有效范围内 / Parameters within valid range
  - [ ] 颜色格式是否正确 / Color format correct
  - [ ] 文件路径是否存在 / File path exists

- [ ] **性能检查** / **Performance Check**
  - [ ] 元素数量是否过多 / Too many elements
  - [ ] 内存使用是否正常 / Memory usage normal
  - [ ] 文件大小是否合理 / File size reasonable

- [ ] **输出检查** / **Output Check**
  - [ ] SVG语法是否正确 / SVG syntax correct
  - [ ] 文件是否成功保存 / File saved successfully
  - [ ] 渲染结果是否符合预期 / Rendering result as expected

---

**版本信息 / Version Info**: v1.0.0  
**最后更新 / Last Updated**: 2024年12月  
**兼容性 / Compatibility**: Go 1.16+