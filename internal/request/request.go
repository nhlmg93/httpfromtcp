package request

import (
	"errors"
	"io"
	"strings"
)

const (
	requestLine = iota
)

type Method = string

const (
	GET    Method = "GET"
	POST   Method = "POST"
	DELETE Method = "DELETE"
	PUT    Method = "PUT"
	PATCH  Method = "PATCH"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	rawReq, err := io.ReadAll(reader)
	if err != nil {
		return &Request{}, err
	}

	reqStr := string(rawReq)
	reqSplit := strings.Split(reqStr, "\r\n")
	reqLine, err := parseRequestLine(reqSplit[requestLine])
	if err != nil {
		return nil, err
	}

	req := &Request{}
	req.RequestLine = reqLine
	return req, nil
}

func parseRequestLine(req string) (RequestLine, error) {
	reqLine := RequestLine{}
	requestLinePieces := strings.Split(req, " ")

	method, reqTarget, httpVersion, err := validateRequestLine(requestLinePieces)
	if err != nil {
		return reqLine, err
	}

	reqLine.Method = method
	reqLine.RequestTarget = reqTarget
	httpVersion = strings.Split(httpVersion, "/")[1]
	reqLine.HttpVersion = httpVersion

	return reqLine, nil
}

func validateRequestLine(reqLine []string) (method, reqTarget, httpVersion string, err error) {
	if len(reqLine) != 3 {
		return "", "", "", errors.New("ERROR: Request Line malformed. Request Line length of 3 expected")
	}

	method = reqLine[0]
	reqTarget = reqLine[1]
	httpVersion = reqLine[2]

	var isMethod = false
	for _, m := range []Method{GET, PATCH, POST, PUT, DELETE} {
		if method == m {
			isMethod = true
		}
	}
	if isMethod == false {
		return "", "", "", errors.New("ERROR: Request Line invalid 'Method'")
	}

	if !strings.Contains(reqTarget, "/") {
		return "", "", "", errors.New("ERROR: Request Line invalid 'PATH'")
	}

	expectedVersion := "HTTP/1.1"
	if httpVersion != expectedVersion {
		return "", "", "", errors.New("ERROR: Request Line invalid 'HTTP Version'")
	}

	return method, reqTarget, httpVersion, nil
}
