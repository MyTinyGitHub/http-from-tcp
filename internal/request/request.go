package request

import (
	"fmt"
	"http-from-tcp/internal/headers"
	"io"
	"strings"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	Headers     map[string]string
	//Body  []byte
}

func (r RequestLine) Print() {
	fmt.Println("Request line:")
	fmt.Printf("- Method: %v\n", r.Method)
	fmt.Printf("- Target: %v\n", r.RequestTarget)
	fmt.Printf("- Version: %v\n", r.HttpVersion)
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	header, err := io.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	bytesRead, requestLine, err := parseRequestLine(header)
	if err != nil {
		return nil, err
	}

	h, err := parseHeader(bytesRead, header)
	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: requestLine,
		Headers:     h,
	}, nil
}

func parseHeader(startOffset int, data []byte) (headers.Headers, error) {
	header := headers.Headers{}
	offset := startOffset

	for {
		readBytes, done, err := header.Parse(data[offset:])
		if err != nil {
			return nil, err
		}

		if done {
			break
		}

		offset += readBytes
	}

	return header, nil
}

func parseRequestLine(data []byte) (int, RequestLine, error) {
	lines := strings.Split(string(data), "\r\n")

	if len(lines) < 1 {
		return 0, RequestLine{}, fmt.Errorf("invalid header")
	}

	fmt.Println(lines[0])
	rLine := strings.Split(lines[0], " ")
	if len(rLine) != 3 {
		return 0, RequestLine{}, fmt.Errorf("missing request line fields")
	}

	requestLine := RequestLine{
		HttpVersion:   strings.Split(rLine[2], "/")[1],
		RequestTarget: rLine[1],
		Method:        rLine[0],
	}

	return len(lines[0]) + 2, requestLine, nil
}
