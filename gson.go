package gson

func parseBytes(bytes []byte) (interface{}, error) {
	endOfString := len(bytes) - 1
	return string(bytes[1:endOfString]), nil
}

func Parse(text string) (interface{}, error) {
	return parseBytes([]byte(text))
}
