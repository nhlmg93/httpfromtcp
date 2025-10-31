package headers

import (
	"bytes"
	"fmt"
)

type Headers map[string]string

const crlf = "\r\n"

func NewHeaders() Headers {
	return make(Headers)
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
		n += idx + len(crlf)
		h[name] = value
	}
	return n, done, err
}
