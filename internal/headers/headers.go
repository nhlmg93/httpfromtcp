package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

const crlf = "\r\n"

func NewHeaders() Headers {
	return make(Headers)
}
func (h Headers) Get(name string) (string, bool) {
	v, ok := h[strings.ToLower(name)]
	return v, ok
}
func (h Headers) Set(name, value string) {
	name = strings.ToLower(name)
	if v, ok := h[name]; ok {
		h[name] = v + ", " + value

	} else {
		h[name] = value
	}
}
func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	n = 0
	done = false
	for {
		idx := bytes.Index(data[n:], []byte(crlf))
		if idx == -1 {
			break
		}
		if idx == 0 {
			done = true
			n += len(crlf)
			break
		}
		name, value, err := parseHeader(data[n : n+idx])
		if err != nil {
			return 0, false, err
		}

		if !isValidToken([]byte(name)) {
			return 0, false, fmt.Errorf("malformed header name")
		}
		n += idx + len(crlf)
		h.Set(name, value)
	}
	return n, done, err
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

func isValidToken(str []byte) bool {
	for _, ch := range str {
		found := false
		switch ch {
		case '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~':
			found = true
		}
		if ch >= 'A' && ch <= 'Z' ||
			ch >= 'a' && ch <= 'z' ||
			ch >= '0' && ch <= '9' {
			found = true
		}
		if !found {
			return false
		}
	}
	return true
}
