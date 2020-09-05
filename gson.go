package gson

import (
	"fmt"
	"strconv"
)

func expectToken(bytes []byte, pos int, expect byte) bool {
	return pos < len(bytes) && bytes[pos] == expect
}

func consumeToken(bytes []byte, pos int, expect string) (string, error) {
	l := len(expect)
	if len(bytes) < pos+l {
		return "", fmt.Errorf("unknown token")
	}
	token := string(bytes[pos:(pos + l)])
	if token == expect {
		return token, nil
	}

	return "", fmt.Errorf("expect %s, but got %s", expect, token)
}

func isWhitespace(b byte) bool {
	return b == ' ' || b == '\n' || b == '\r' || b == '\t'
}

func skipWhitespace(bytes []byte, pos int) int {
	for pos < len(bytes) && isWhitespace(bytes[pos]) {
		pos++
	}
	return pos
}

func isDigit(b byte) bool {
	return '0' <= b && b <= '9'
}

func isNumber(b byte) bool {
	return b == '-' || isDigit(b)
}

func parseNumber(bytes []byte, pos int) (interface{}, int, error) {
	begin := pos
	end := pos + 1
	if bytes[pos] == '-' {
		end += 1
	}
	for end < len(bytes) && isDigit(bytes[end]) {
		end += 1
	}

	if end == len(bytes) || bytes[end] != '.' {
		ret, err := strconv.Atoi(string(bytes[begin:end]))
		if err != nil {
			return nil, pos, err
		}
		return ret, end, nil
	}

	//fraction
	end += 1
	for end < len(bytes) && isDigit(bytes[end]) {
		end += 1
	}

	ret, err := strconv.ParseFloat(string(bytes[begin:end]), 64)
	if err != nil {
		return nil, pos, err
	}
	return ret, end, nil

	//exponent
}

func parseValue(bytes []byte, pos int) (interface{}, int, error) {
	// case whitespace -> value -> whitespace
	pos = skipWhitespace(bytes, pos)
	v, newpos, err := nextToken(bytes, pos)
	if err != nil {
		return v, newpos, err
	}
	newpos = skipWhitespace(bytes, newpos)
	return v, newpos, nil
}

func nextToken(bytes []byte, pos int) (interface{}, int, error) {
	if pos >= len(bytes) {
		return nil, pos, fmt.Errorf("try to scan next token but got end of input")
	}

	switch bytes[pos] {
	case '"':
		return parseString(bytes, pos)
	case 'n':
		return parseNull(bytes, pos)
	case 't':
		return parseTrue(bytes, pos)
	case 'f':
		return parseFalse(bytes, pos)
	case '[':
		return parseArray(bytes, pos)
	default:
		if isNumber(bytes[pos]) {
			return parseNumber(bytes, pos)
		}

		return nil, -1, fmt.Errorf("unknown token")
	}
}

func parseNull(bytes []byte, pos int) (interface{}, int, error) {
	if _, err := consumeToken(bytes, pos, "null"); err != nil {
		return nil, pos, err
	}
	return nil, pos + 4, nil
}

func parseTrue(bytes []byte, pos int) (interface{}, int, error) {
	if _, err := consumeToken(bytes, pos, "true"); err != nil {
		return nil, pos, err
	}
	return true, pos + 4, nil
}

func parseFalse(bytes []byte, pos int) (interface{}, int, error) {
	if _, err := consumeToken(bytes, pos, "false"); err != nil {
		return nil, pos, err
	}
	return false, pos + 5, nil
}

func parseString(bytes []byte, pos int) (interface{}, int, error) {
	start := pos + 1
	var end int
	for end = start; !expectToken(bytes, end, '"'); end++ {
	}
	result := string(bytes[start:end])
	return result, end + 1, nil
}

func parseArrayElem(bytes []byte, pos int) (interface{}, int, error) {
	if !expectToken(bytes, pos, ',') {
		return nil, 0, fmt.Errorf("expect ',' but got %v", bytes[pos])
	}
	return parseValue(bytes, pos+1)
}

func parseArray(bytes []byte, pos int) (interface{}, int, error) {
	pos += 1
	ret := make([]interface{}, 0)
	var val interface{}
	var err error
	pos = skipWhitespace(bytes, pos)
	// case: '[' -> whitespace -> ']'
	if expectToken(bytes, pos, ']') {
		return ret, pos + 1, nil
	}

	// case: '[' -> value
	val, pos, err = parseValue(bytes, pos)
	if err != nil {
		return nil, pos, err
	}
	if pos == len(bytes) {
		return nil, pos, fmt.Errorf("not found ']'")
	}
	ret = append(ret, val)

	for !expectToken(bytes, pos, ']') {
		// case: ',' -> value
		val, pos, err = parseArrayElem(bytes, pos)
		if err != nil {
			return nil, pos, err
		}

		if pos == len(bytes) {
			return nil, pos, fmt.Errorf("not found ']'")
		}
		ret = append(ret, val)
	}
	return ret, pos + 1, nil
}

func Parse(text string) (interface{}, error) {
	result, _, err := parseValue([]byte(text), 0)
	return result, err
}
