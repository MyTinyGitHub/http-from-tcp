package request

import (
	"fmt"
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
	//Headers   map[string]string
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

	lines := strings.Split(string(header), "\r\n")
	if len(lines) < 1 {
		return nil, fmt.Errorf("invalid header")
	}

	fmt.Println(lines[0])
	rLine := strings.Split(lines[0], " ")
	if len(rLine) != 3 {
		return nil, fmt.Errorf("missing request line fields")
	}

	requestLine := RequestLine{
		HttpVersion:   strings.Split(rLine[2], "/")[1],
		RequestTarget: rLine[1],
		Method:        rLine[0],
	}

	return &Request{
		RequestLine: requestLine,
	}, nil
}
