package path

// 解析各种命令的函数
func parseMoveToAbs(tokenizer *pathTokenizer) (Command, error) {
	x, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	y, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	return &MoveToCommand{X: x, Y: y, Relative: false}, nil
}

func parseMoveToRel(tokenizer *pathTokenizer) (Command, error) {
	dx, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	dy, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	return &MoveToCommand{X: dx, Y: dy, Relative: true}, nil
}

func parseLineToAbs(tokenizer *pathTokenizer) (Command, error) {
	x, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	y, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	return &LineToCommand{X: x, Y: y, Relative: false}, nil
}

func parseLineToRel(tokenizer *pathTokenizer) (Command, error) {
	dx, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	dy, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	return &LineToCommand{X: dx, Y: dy, Relative: true}, nil
}

func parseHorizontalLineToAbs(tokenizer *pathTokenizer) (Command, error) {
	x, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	return &HorizontalLineToCommand{X: x, Relative: false}, nil
}

func parseHorizontalLineToRel(tokenizer *pathTokenizer) (Command, error) {
	dx, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	return &HorizontalLineToCommand{X: dx, Relative: true}, nil
}

func parseVerticalLineToAbs(tokenizer *pathTokenizer) (Command, error) {
	y, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	return &VerticalLineToCommand{Y: y, Relative: false}, nil
}

func parseVerticalLineToRel(tokenizer *pathTokenizer) (Command, error) {
	dy, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	return &VerticalLineToCommand{Y: dy, Relative: true}, nil
}

func parseCubicCurveToAbs(tokenizer *pathTokenizer) (Command, error) {
	x1, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	y1, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	x2, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	y2, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	x, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	y, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	return &CubicCurveToCommand{X1: x1, Y1: y1, X2: x2, Y2: y2, X: x, Y: y, Relative: false}, nil
}

func parseCubicCurveToRel(tokenizer *pathTokenizer) (Command, error) {
	dx1, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	dy1, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	dx2, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	dy2, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	dx, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	dy, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	return &CubicCurveToCommand{X1: dx1, Y1: dy1, X2: dx2, Y2: dy2, X: dx, Y: dy, Relative: true}, nil
}

func parseSmoothCubicCurveToAbs(tokenizer *pathTokenizer) (Command, error) {
	x2, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	y2, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	x, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	y, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	return &SmoothCubicCurveToCommand{X2: x2, Y2: y2, X: x, Y: y, Relative: false}, nil
}

func parseSmoothCubicCurveToRel(tokenizer *pathTokenizer) (Command, error) {
	dx2, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	dy2, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	dx, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	dy, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	return &SmoothCubicCurveToCommand{X2: dx2, Y2: dy2, X: dx, Y: dy, Relative: true}, nil
}

func parseQuadraticCurveToAbs(tokenizer *pathTokenizer) (Command, error) {
	x1, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	y1, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	x, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	y, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	return &QuadraticCurveToCommand{X1: x1, Y1: y1, X: x, Y: y, Relative: false}, nil
}

func parseQuadraticCurveToRel(tokenizer *pathTokenizer) (Command, error) {
	dx1, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	dy1, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	dx, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	dy, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	return &QuadraticCurveToCommand{X1: dx1, Y1: dy1, X: dx, Y: dy, Relative: true}, nil
}

func parseSmoothQuadraticCurveToAbs(tokenizer *pathTokenizer) (Command, error) {
	x, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	y, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	return &SmoothQuadraticCurveToCommand{X: x, Y: y, Relative: false}, nil
}

func parseSmoothQuadraticCurveToRel(tokenizer *pathTokenizer) (Command, error) {
	dx, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	dy, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	return &SmoothQuadraticCurveToCommand{X: dx, Y: dy, Relative: true}, nil
}

func parseArcToAbs(tokenizer *pathTokenizer) (Command, error) {
	rx, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	ry, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	xAxisRotation, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	largeArc, err := tokenizer.nextBool()
	if err != nil {
		return nil, err
	}

	sweep, err := tokenizer.nextBool()
	if err != nil {
		return nil, err
	}

	x, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	y, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	return &ArcToAbs{RX: rx, RY: ry, XAxisRotation: xAxisRotation, LargeArc: largeArc, Sweep: sweep, X: x, Y: y}, nil
}

func parseArcToRel(tokenizer *pathTokenizer) (Command, error) {
	rx, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	ry, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	xAxisRotation, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	largeArc, err := tokenizer.nextBool()
	if err != nil {
		return nil, err
	}

	sweep, err := tokenizer.nextBool()
	if err != nil {
		return nil, err
	}

	dx, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	dy, err := tokenizer.nextFloat()
	if err != nil {
		return nil, err
	}

	return &ArcToRel{RX: rx, RY: ry, XAxisRotation: xAxisRotation, LargeArc: largeArc, Sweep: sweep, DX: dx, DY: dy}, nil
}

func parseClosePath() (Command, error) {
	return &ClosePathCommand{}, nil
}
