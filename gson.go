package gson

import "fmt"

func expectToken(bytes []byte, pos int, expect byte) bool {
	if pos >= len(bytes) {
		return false
	}

	return bytes[pos] == expect
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

	if bytes[pos] == '"' {
		return parseString(bytes, pos)
	} else if bytes[pos] == 'n' {
		return parseNull(bytes, pos)
	} else if bytes[pos] == 't' {
		return parseTrue(bytes, pos)
	} else if bytes[pos] == 'f' {
		return parseFalse(bytes, pos)
	} else if bytes[pos] == '[' {
		return parseArray(bytes, pos)
	} else {
		return nil, -1, fmt.Errorf("unknown token")
	}
}

func parseNull(bytes []byte, pos int) (interface{}, int, error) {
	if len(bytes)-pos >= 4 && bytes[pos] == 'n' && bytes[pos+1] == 'u' && bytes[pos+2] == 'l' && bytes[pos+3] == 'l' {
		return nil, pos + 4, nil
	} else {
		return nil, pos, fmt.Errorf("unknown token")
	}
}

func parseTrue(bytes []byte, pos int) (interface{}, int, error) {
	if len(bytes)-pos >= 4 && bytes[pos] == 't' && bytes[pos+1] == 'r' && bytes[pos+2] == 'u' && bytes[pos+3] == 'e' {
		return true, pos + 4, nil
	} else {
		return nil, pos, fmt.Errorf("unknown token")
	}
}

func parseFalse(bytes []byte, pos int) (interface{}, int, error) {
	if len(bytes)-pos >= 5 && bytes[pos] == 'f' && bytes[pos+1] == 'a' && bytes[pos+2] == 'l' && bytes[pos+3] == 's' && bytes[pos+4] == 'e' {
		return false, pos + 5, nil
	} else {
		return nil, pos, fmt.Errorf("unknown token")
	}
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
