package request

import (
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

// each line of an http request should end with these
var REQUEST_LINE_END = "\r\n"

func (r *RequestLine) IsValidHTTP() bool {
	for _, val := range r.Method {
		if !unicode.IsLetter(val) {
			return false
		}
	}

	// grab the http version parts
	httpParts := strings.Split(r.HttpVersion, "/")
	if len(httpParts) != 2 || httpParts[0] != "HTTP" || httpParts[1] != "1.1" {
		return false
	}

	return true
}

// parse out a line of an http request
func parseRequestLine(b string) (*RequestLine, string, error) {
	// grab the index of the end of the line
	idx := strings.Index(b, REQUEST_LINE_END)
	if idx == -1 {
		return nil, b, nil
	}

	// grab the first line
	startLine := b[:idx]
	// this is the rest of the request
	remainingMsg := b[idx+len(REQUEST_LINE_END):]

	// grab the three parts of first line, Method, path, spec
	lines := strings.Split(startLine, " ")
	if len(lines) != 3 {
		return nil, b, fmt.Errorf("malformed request")
	}

	// build up the request line
	reqLine := &RequestLine{
		Method:        lines[0],
		RequestTarget: lines[1],
		HttpVersion:   lines[2],
	}

	if !reqLine.IsValidHTTP() {
		return reqLine, b, fmt.Errorf("invalid HTTP request")
	}
	// we return the line of the request, the rest of the request, nil
	return reqLine, remainingMsg, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	// grab the data from reader
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("unable to io.ReadAll: %w", err)
	}

	// grab a string from the data
	str := string(data)
	reqLine, _, err := parseRequestLine(str)

	// make the Request
	req := &Request{
		RequestLine: *reqLine,
	}

	return req, nil
}
