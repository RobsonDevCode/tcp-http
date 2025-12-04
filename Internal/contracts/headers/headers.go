package headers

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"tcp-http/Internal/constants/seperatorConstants"
)

type Headers struct {
	headers map[string]string
}

func NewHeaders() *Headers {
	return &Headers{
		headers: map[string]string{},
	}
}

func (h *Headers) TryGet(name string) (string, bool) {
	str, exists := h.headers[strings.ToLower(name)]
	return str, exists
}

func (h *Headers) Set(name, value string) {
	h.headers[strings.ToLower(name)] = value
}

func (h *Headers) TrySet(name, value string) {
	name = strings.ToLower(name)
	prev, exists := h.headers[name]
	if exists {
		h.headers[name] = fmt.Sprintf("%s,%s", prev, value)
		return
	}

	h.headers[name] = value
}

func (h *Headers) TryReplace(name, value string) error {
	n := strings.ToLower(name)
	_, exists := h.headers[name]
	if !exists {
		return fmt.Errorf("cannot replace value with value: %s, key %s does not exist", value, name)
	}
	h.headers[n] = value

	return nil
}

func (h *Headers) Parse(data []byte) (int, bool, error) {
	read := 0
	done := false

	for {
		idx := bytes.Index(data[read:], seperatorConstants.RNSEPERATOR)
		if idx == -1 {
			break
		}

		//empty header
		if idx == 0 {
			done = true
			read += len(seperatorConstants.RNSEPERATOR)
			break
		}

		name, value, err := parseHeader(data[read : read+idx])
		if err != nil {
			return 0, false, err
		}

		if !isToken([]byte(name)) {
			return 0, false, fmt.Errorf("malformed header name")
		}

		read += idx + len(seperatorConstants.RNSEPERATOR)

		h.TrySet(name, value)
	}

	return read, done, nil
}

func (h *Headers) ForEach(cb func(n, v string)) {
	for n, v := range h.headers {
		cb(n, v)
	}
}

func GetInt(headers *Headers, name string, defaultValue int) int {
	valueStr, exists := headers.TryGet(name)
	if !exists {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}
func parseHeader(fieldLine []byte) (string, string, error) {
	parts := bytes.SplitN(fieldLine, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("malformed field line")
	}

	name := parts[0]
	value := bytes.TrimSpace(parts[1])

	if bytes.HasSuffix(name, []byte(" ")) {
		return "", "", fmt.Errorf("malformed field name")
	}

	return string(name), string(value), nil
}

func isToken(str []byte) bool {
	if len(str) == 0 {
		return false
	}

	for _, ch := range str {
		found := false
		switch {
		case ch >= 'A' && ch <= 'Z' ||
			ch >= 'a' && ch <= 'z' ||
			ch >= '0' && ch <= '9':
			found = true
		case ch == '!', ch == '#', ch == '$', ch == '%', ch == '&',
			ch == '\'', ch == '*', ch == '+', ch == '-', ch == '.',
			ch == '^', ch == '_', ch == '`', ch == '|', ch == '~':
			found = true
		}
		if !found {
			return false
		}
	}
	return true
}
