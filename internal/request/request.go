package request

import (
	"bytes"
	"fmt"
	"http-from-tcp/internal/headers"
	"io"
	"strconv"
	"strings"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte

	state requestState
}

type requestState int

const (
	requestStateInitialized requestState = iota
	requestStateParsingHeaders
	requestStateParsingBody
	requestStateDone
)

func (r Request) Print() {
	fmt.Println("Request line:")
	fmt.Printf("- Method: %v\n", r.RequestLine.Method)
	fmt.Printf("- Target: %v\n", r.RequestLine.RequestTarget)
	fmt.Printf("- Version: %v\n", r.RequestLine.HttpVersion)
	fmt.Println("Headers:")
	for key, value := range r.Headers {
		fmt.Printf("- %v: %v\n", key, value)
	}
	fmt.Println("Body:")
	fmt.Println(string(r.Body))
}

const bufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	buffer := make([]byte, bufferSize, bufferSize)
	currentOffset := 0
	req := &Request{
		state:   requestStateInitialized,
		Headers: headers.NewHeaders(),
	}

	for req.state != requestStateDone {
		if currentOffset >= len(buffer) {
			newBuffer := make([]byte, len(buffer)*2)
			copy(newBuffer, buffer)
			buffer = newBuffer
		}

		readBytes, err := reader.Read(buffer[currentOffset:])
		if err != nil {
			if err == io.EOF {
				bytes, err := req.parse(buffer[:currentOffset])

				if err != nil {
					return nil, err
				}

				if bytes <= 0 {
					return nil, fmt.Errorf("unexpected end of reader %d", bytes)
				}

				return req, nil
			}
			return nil, err
		}

		currentOffset += readBytes

		numBytesParsed, err := req.parse(buffer[:currentOffset])
		if err != nil {
			return nil, err
		}

		copy(buffer, buffer[numBytesParsed:])
		currentOffset -= numBytesParsed
	}

	return req, nil
}

func (req *Request) parse(data []byte) (int, error) {
	if req.state == requestStateInitialized {
		return parseRequestLine(req, data)
	} else if req.state == requestStateParsingHeaders {
		return parseHeader(req, data)
	} else if req.state == requestStateParsingBody {
		return parseBody(req, data)
	}

	return 0, nil
}

func parseBody(req *Request, data []byte) (int, error) {
	contentLength, ok := req.Headers["content-length"]
	if !ok {
		req.state = requestStateDone
		return 2, nil
	}

	bodySize, _ := strconv.Atoi(contentLength)

	req.Body = append(req.Body, data...)
	if len(req.Body) == bodySize {
		req.state = requestStateDone
	}

	return len(data), nil
}

func parseHeader(req *Request, data []byte) (int, error) {
	readBytes, done, err := req.Headers.Parse(data)
	if err != nil {
		return 0, err
	}

	if done {
		req.state = requestStateParsingBody
	}

	return readBytes, nil
}

func parseRequestLine(req *Request, data []byte) (int, error) {
	idx := bytes.Index(data, []byte("\r\n"))

	if idx == -1 {
		return 0, nil
	}

	requestLineText := string(data[:idx])
	parts := strings.Split(requestLineText, " ")
	if len(parts) != 3 {
		return 0, fmt.Errorf("missing request line fields")
	}

	req.RequestLine = RequestLine{
		HttpVersion:   strings.Split(parts[2], "/")[1],
		RequestTarget: parts[1],
		Method:        parts[0],
	}

	req.state = requestStateParsingHeaders

	return idx + 2, nil
}
