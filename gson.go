package gson

import "fmt"

func parseBytes(bytes []byte) (interface{}, error) {
	var result interface{}
	pos := 0
	for pos < len(bytes) {
		if bytes[pos] == '"' {
			startString := pos + 1
			endOfString := len(bytes) - 1
			result = string(bytes[startString:endOfString])
			pos = len(bytes)
		} else if len(bytes)-pos >= 4 && bytes[pos] == 'n' && bytes[pos+1] == 'u' && bytes[pos+2] == 'l' && bytes[pos+3] == 'l' {
			result = nil
			pos += 4
		} else if len(bytes)-pos >= 4 && bytes[pos] == 't' && bytes[pos+1] == 'r' && bytes[pos+2] == 'u' && bytes[pos+3] == 'e' {
			result = true
			pos += 4
		} else if len(bytes)-pos >= 5 && bytes[pos] == 'f' && bytes[pos+1] == 'a' && bytes[pos+2] == 'l' && bytes[pos+3] == 's' && bytes[pos+4] == 'e' {
			result = false
			pos += 5
		} else {
			return nil, fmt.Errorf("unknown token")
		}
	}
	return result, nil
}

func Parse(text string) (interface{}, error) {
	return parseBytes([]byte(text))
}
