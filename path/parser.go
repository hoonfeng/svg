package path

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

// SVGPath 表示SVG路径
type SVGPath struct {
	Commands []Command
}

// ParsePath 解析SVG路径数据
func ParsePath(data string) (*SVGPath, error) {
	// 创建路径对象
	path := &SVGPath{
		Commands: []Command{},
	}

	// 解析路径数据
	tokens, err := tokenizePath(data)
	if err != nil {
		return nil, err
	}

	// 解析命令
	commands, err := parseCommands(tokens)
	if err != nil {
		return nil, err
	}

	path.Commands = commands
	return path, nil
}

// tokenizePath 将路径数据分解为标记
func tokenizePath(data string) ([]string, error) {
	// 预处理数据
	// 在命令字符前后添加空格
	for _, cmd := range "MmLlHhVvCcSsQqTtAaZz" {
		data = strings.ReplaceAll(data, string(cmd), " "+string(cmd)+" ")
	}

	// 将逗号替换为空格
	data = strings.ReplaceAll(data, ",", " ")
	
	// 处理连续的负数：在负号前添加空格（除了字符串开头和已有空格后的负号）
	result := ""
	for i, char := range data {
		if char == '-' && i > 0 && data[i-1] != ' ' && data[i-1] != ',' {
			result += " -"
		} else {
			result += string(char)
		}
	}
	data = result

	// 将连续的空格替换为单个空格
	for strings.Contains(data, "  ") {
		data = strings.ReplaceAll(data, "  ", " ")
	}

	// 去除首尾空格
	data = strings.TrimSpace(data)

	// 分割为标记
	tokens := strings.Fields(data)

	return tokens, nil
}

// parseCommands 解析命令标记
func parseCommands(tokens []string) ([]Command, error) {
	commands := []Command{}

	// 当前命令类型
	var currentCmd rune

	// 参数列表
	params := []float64{}

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		// 检查是否是命令字符
		if len(token) == 1 && strings.ContainsAny(token, "MmLlHhVvCcSsQqTtAaZz") {
			// 如果有参数，处理前一个命令
			if len(params) > 0 && currentCmd != 0 {
				cmd, err := createCommand(currentCmd, params)
				if err != nil {
					return nil, err
				}
				commands = append(commands, cmd...)
				params = []float64{}
			}

			// 更新当前命令
			currentCmd = rune(token[0])
			
			// 对于Z命令（闭合路径），立即处理，因为它没有参数
			if currentCmd == 'Z' || currentCmd == 'z' {
				cmd, err := createCommand(currentCmd, []float64{})
				if err != nil {
					return nil, err
				}
				commands = append(commands, cmd...)
				currentCmd = 0 // 重置当前命令
			}
		} else {
			// 解析参数
			param, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return nil, err
			}
			params = append(params, param)

			// 检查是否有足够的参数来创建命令
			if currentCmd != 0 {
				paramCount := getParamCount(currentCmd)
				if paramCount > 0 && len(params) >= paramCount {
					cmd, err := createCommand(currentCmd, params[:paramCount])
					if err != nil {
						return nil, err
					}
					commands = append(commands, cmd...)

					// 对于某些命令（如M、L等），后续的坐标对应于相同类型的命令
					if currentCmd == 'M' {
						currentCmd = 'L'
					} else if currentCmd == 'm' {
						currentCmd = 'l'
					}

					// 移除已处理的参数
					params = params[paramCount:]
				}
			}
		}
	}

	// 处理最后一个命令
	if len(params) > 0 && currentCmd != 0 {
		cmd, err := createCommand(currentCmd, params)
		if err != nil {
			return nil, err
		}
		commands = append(commands, cmd...)
	}

	return commands, nil
}

// getParamCount 获取命令的参数数量
func getParamCount(cmd rune) int {
	switch cmd {
	case 'M', 'm', 'L', 'l', 'T', 't':
		return 2
	case 'H', 'h', 'V', 'v':
		return 1
	case 'C', 'c':
		return 6
	case 'S', 's', 'Q', 'q':
		return 4
	case 'A', 'a':
		return 7
	case 'Z', 'z':
		return 0
	default:
		return -1
	}
}

// createCommand 创建命令对象
func createCommand(cmd rune, params []float64) ([]Command, error) {
	isRelative := unicode.IsLower(cmd)

	switch unicode.ToUpper(cmd) {
	case 'M': // 移动
		if len(params) < 2 {
			return nil, errors.New("M命令需要至少2个参数")
		}
		return []Command{&MoveToCommand{X: params[0], Y: params[1], Relative: isRelative}}, nil

	case 'L': // 直线
		if len(params) < 2 {
			return nil, errors.New("L命令需要至少2个参数")
		}
		return []Command{&LineToCommand{X: params[0], Y: params[1], Relative: isRelative}}, nil

	case 'H': // 水平线
		if len(params) < 1 {
			return nil, errors.New("H命令需要至少1个参数")
		}
		return []Command{&HorizontalLineToCommand{X: params[0], Relative: isRelative}}, nil

	case 'V': // 垂直线
		if len(params) < 1 {
			return nil, errors.New("V命令需要至少1个参数")
		}
		return []Command{&VerticalLineToCommand{Y: params[0], Relative: isRelative}}, nil

	case 'C': // 三次贝塞尔曲线
		if len(params) < 6 {
			return nil, errors.New("C命令需要至少6个参数")
		}
		return []Command{&CubicCurveToCommand{
			X1: params[0], Y1: params[1],
			X2: params[2], Y2: params[3],
			X: params[4], Y: params[5],
			Relative: isRelative,
		}}, nil

	case 'S': // 平滑三次贝塞尔曲线
		if len(params) < 4 {
			return nil, errors.New("S命令需要至少4个参数")
		}
		return []Command{&SmoothCubicCurveToCommand{
			X2: params[0], Y2: params[1],
			X: params[2], Y: params[3],
			Relative: isRelative,
		}}, nil

	case 'Q': // 二次贝塞尔曲线
		if len(params) < 4 {
			return nil, errors.New("Q命令需要至少4个参数")
		}
		return []Command{&QuadraticCurveToCommand{
			X1: params[0], Y1: params[1],
			X: params[2], Y: params[3],
			Relative: isRelative,
		}}, nil

	case 'T': // 平滑二次贝塞尔曲线
		if len(params) < 2 {
			return nil, errors.New("T命令需要至少2个参数")
		}
		return []Command{&SmoothQuadraticCurveToCommand{
			X: params[0], Y: params[1],
			Relative: isRelative,
		}}, nil

	case 'A': // 椭圆弧
		if len(params) < 7 {
			return nil, errors.New("A命令需要至少7个参数")
		}
		return []Command{&ArcToCommand{
			RX: params[0], RY: params[1],
			XAxisRotation: params[2],
			LargeArc:      params[3] != 0,
			Sweep:         params[4] != 0,
			X:             params[5], Y: params[6],
			Relative: isRelative,
		}}, nil

	case 'Z': // 闭合路径
		return []Command{&ClosePathCommand{}}, nil

	default:
		return nil, errors.New("未知的命令: " + string(cmd))
	}
}
