package request

import (
    "fmt"
    "tcp-http/Internal/contracts/headers"
)

type Request struct {
    RequestLine RequestLine
    Headers     *headers.Headers
    State       ParserState
    Body        string
}

type ParserState string

const (
    StateInit    ParserState = "init"
    StateHeaders ParserState = "headers"
    StateBody    ParserState = "body"
    StateError   ParserState = "error"
    StateDone    ParserState = "done"
)

var ErrorRequestInErrorState = fmt.Errorf("request is in error state")

func NewRequest() *Request {
    return &Request{
        State:   StateInit,
        Headers: headers.NewHeaders(),
        Body:    "",
    }
}
func (r *Request) ParseRequest(data []byte) (int, error) {
    read := 0
outer:
    for {
        currentData := data[read:]
        if len(currentData) == 0 {
            break outer
        }

        switch r.State {
        case StateDone:
            break outer

        case StateInit:
            requestLine, n, err := parseRequestLine(currentData)
            if err != nil {
                r.State = StateError
                return 0, err
            }

            if n == 0 {
                break outer
            }

            r.RequestLine = *requestLine
            read += n
            r.State = StateHeaders

        case StateHeaders:

            n, done, err := r.Headers.Parse(currentData)
            if err != nil {
                r.State = StateError
                return 0, err
            }

            if n == 0 {
                break outer
            }
            read += n
            if done {
                r.State = StateBody
            }

        case StateBody:
            contentLength := headers.GetInt(r.Headers, "content-length", 0)
            if contentLength == 0 {
                r.State = StateDone
                break outer
            }

            remaining := min(contentLength-len(r.Body), len(currentData))
            r.Body += string(currentData[:remaining])
            read += remaining

            if len(r.Body) == contentLength {
                r.State = StateDone
            }
        case StateError:
            return 0, ErrorRequestInErrorState

        default:
            panic("we dont expect this")
        }
    }

    return read, nil
}

func (r *Request) IsDone() bool {
    return r.State == StateDone
}

func (r *Request) IsError() bool {
    return r.State == StateError
}
