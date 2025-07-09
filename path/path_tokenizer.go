package path

import (
	"strconv"
	"strings"
)

type pathTokenizer struct {
	data   string
	pos    int
	length int
}

func newPathTokenizer(data string) *pathTokenizer {
	return &pathTokenizer{
		data:   strings.TrimSpace(data),
		pos:    0,
		length: len(data),
	}
}

func (t *pathTokenizer) nextFloat() (float64, error) {
	start := t.pos
	for t.pos < t.length && !isWhitespace(t.data[t.pos]) && t.data[t.pos] != ',' {
		t.pos++
	}

	str := t.data[start:t.pos]
	t.skipSeparators()

	return strconv.ParseFloat(str, 64)
}

func (t *pathTokenizer) nextBool() (bool, error) {
	start := t.pos
	for t.pos < t.length && !isWhitespace(t.data[t.pos]) && t.data[t.pos] != ',' {
		t.pos++
	}

	str := t.data[start:t.pos]
	t.skipSeparators()

	return strconv.ParseBool(str)
}

func (t *pathTokenizer) skipSeparators() {
	for t.pos < t.length && (isWhitespace(t.data[t.pos]) || t.data[t.pos] == ',') {
		t.pos++
	}
}

func isWhitespace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}
