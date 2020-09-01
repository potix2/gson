package gson

import "fmt"

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
	pos = skipWhitespace(bytes, pos)
	return nextToken(bytes, pos)
}

func nextToken(bytes []byte, pos int) (interface{}, int, error) {
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
	end := start
	for end = start; bytes[end] != '"'; end++ {
	}
	result := string(bytes[start:end])
	return result, end + 1, nil
}

func parseArray(bytes []byte, pos int) (interface{}, int, error) {
	ret := make([]interface{}, 0)
	for pos = skipWhitespace(bytes, pos+1); pos < len(bytes) && bytes[pos] != ']'; pos = skipWhitespace(bytes, pos) {
		var val interface{}
		var err error
		val, pos, err = parseValue(bytes, pos)
		if err != nil {
			return nil, pos, err
		}
		ret = append(ret, val)
	}

	if bytes[pos-1] != ']' && len(bytes) == pos {
		return nil, pos, fmt.Errorf("not found ']'")
	}
	return ret, pos + 1, nil
}

func Parse(text string) (interface{}, error) {
	result, _, err := parseValue([]byte(text), 0)
	return result, err
}
