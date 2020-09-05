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

	//exponent
	if end < len(bytes) && (bytes[end] == 'e' || bytes[end] == 'E') {
		end += 1
		if end == len(bytes) || (bytes[end] != '+' && bytes[end] != '-') {
			return nil, end, fmt.Errorf("invalid number format: %s", string(bytes[begin:end]))
		}

		end += 1
		for end < len(bytes) && isDigit(bytes[end]) {
			end += 1
		}
	}

	ret, err := strconv.ParseFloat(string(bytes[begin:end]), 64)
	if err != nil {
		return nil, pos, err
	}
	return ret, end, nil

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

func parseArrayElem(bytes []byte, pos int, ret []interface{}) ([]interface{}, int, error) {
	val, pos, err := parseValue(bytes, pos)
	if err != nil {
		return nil, pos, err
	}

	if pos == len(bytes) {
		return nil, pos, fmt.Errorf("not found ']'")
	}
	ret = append(ret, val)
	return ret, pos, nil
}

func parseArray(bytes []byte, pos int) (interface{}, int, error) {
	pos += 1
	ret := make([]interface{}, 0)
	var err error
	pos = skipWhitespace(bytes, pos)
	// case: '[' -> whitespace -> ']'
	if expectToken(bytes, pos, ']') {
		return ret, pos + 1, nil
	}

	// case: '[' -> value
	ret, pos, err = parseArrayElem(bytes, pos, ret)
	if err != nil {
		return nil, pos, err
	}

	for !expectToken(bytes, pos, ']') {
		// case: ',' -> value
		if !expectToken(bytes, pos, ',') {
			return nil, 0, fmt.Errorf("expect ',' but got %s", string(bytes[pos]))
		}

		ret, pos, err = parseArrayElem(bytes, pos+1, ret)
		if err != nil {
			return nil, pos, err
		}
	}

	return ret, pos + 1, nil
}

func parseObjectItem(bytes []byte, pos int, ret map[string]interface{}) (map[string]interface{}, int, error) {
	pos = skipWhitespace(bytes, pos)
	var key, val interface{}
	var err error
	key, pos, err = parseString(bytes, pos)
	if err != nil {
		return nil, pos, err
	}
	pos = skipWhitespace(bytes, pos)
	if bytes[pos] != ':' {
		return nil, pos, fmt.Errorf("expect ':', but got %s", string(bytes[pos]))
	}
	pos += 1
	val, pos, err = parseValue(bytes, pos)
	if err != nil {
		return nil, pos, err
	}

	if pos == len(bytes) {
		return nil, pos, fmt.Errorf("not found '}'")
	}

	ret[key.(string)] = val
	return ret, pos, nil
}

func parseObject(bytes []byte, pos int) (interface{}, int, error) {
	pos += 1
	ret := make(map[string]interface{}, 0)
	pos = skipWhitespace(bytes, pos)
	// case: '{' -> whitespace -> '}'
	if expectToken(bytes, pos, '}') {
		return ret, pos + 1, nil
	}

	var err error
	ret, pos, err = parseObjectItem(bytes, pos, ret)
	if err != nil {
		return nil, pos, err
	}

	for !expectToken(bytes, pos, '}') {
		if !expectToken(bytes, pos, ',') {
			return nil, pos, fmt.Errorf("expect ',', but got %s", string(bytes[pos]))
		}

		ret, pos, err = parseObjectItem(bytes, pos+1, ret)
		if err != nil {
			return nil, pos, err
		}
	}
	return ret, pos + 1, nil
}

func Parse(text string) (interface{}, error) {
	result, _, err := parseValue([]byte(text), 0)
	return result, err
}
