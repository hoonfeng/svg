package path

// Command 接口定义了SVG路径命令的基本行为
type Command interface {
	Execute(ctx *PathContext, precision float64)
	String() string
}
