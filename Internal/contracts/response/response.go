package response

import (
	"fmt"
	"io"
	"strconv"
	"tcp-http/Internal/contracts/headers"
	"tcp-http/Internal/contracts/request"
)

type Response struct {
}

type StatusCode int

const (
	StatusOK                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

type HandlerError struct {
	StatusCode StatusCode
	Message    string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func WithDefaultHeaders(contentLength int) *headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", strconv.Itoa(contentLength))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")

	return h
}

func WriteHeaders(w io.Writer, h *headers.Headers) error {
	b := []byte{}
	h.ForEach(func(n, v string) {
		b = fmt.Appendf(b, "%s: %s\r\n", n, v)
	})

	b = fmt.Append(b, "\r\n")
	_, err := w.Write(b)
	return err
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	var statusLine []byte
	switch statusCode {
	case StatusOK:
		statusLine = []byte("HTTP/1.1 200 OK\r\n")
	case StatusBadRequest:
		statusLine = []byte("HTTP/1.1 400 Bad Reqeuest\r\n")
	case StatusInternalServerError:
		statusLine = []byte("Http/1.1 500 Internal Server Error\r\n")
	default:
		return fmt.Errorf("not supported status code")
	}

	_, err := w.Write(statusLine)
	return err
}
