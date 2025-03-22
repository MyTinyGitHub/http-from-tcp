package headers

import (
	"fmt"
	"slices"
	"strings"
	"unicode"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(map[string]string)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	if ok := strings.Contains(string(data), "\r\n"); !ok {
		return 0, false, nil
	}

	header := strings.SplitAfter(string(data), "\r\n")[0]

	crlfIndex := strings.Index(header, "\r\n")
	if crlfIndex == 0 {
		return len(header), true, nil
	}

	trimmed := strings.Trim(header, " ")
	split := strings.SplitN(trimmed, ":", 2)

	if len(split) != 2 {
		return 0, false, fmt.Errorf("incorrect format missing : ")
	}

	if strings.HasSuffix(split[0], " ") {
		return 0, false, fmt.Errorf("incorrect format, empty space before : ")
	}

	value := strings.Trim(split[1], "\r\n")
	value = strings.Trim(value, " ")

	key := strings.ToLower(split[0])
	if strings.ContainsFunc(key, isValid) {
		return 0, false, fmt.Errorf("header contains symbols that are not allowed")
	}

	val, ok := h[key]
	if ok {
		value = val + ", " + value
	}

	h[key] = value

	return len(header), false, nil
}

var allowedSpecialChars = []rune{'!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '`'}

func isValid(r rune) bool {
	return !unicode.IsDigit(r) && !unicode.IsLetter(r) && !slices.Contains(allowedSpecialChars, r)
}
