package contracts

import "fmt"

type Request struct {
    RequestLine RequestLine
    State       ParserState
}

type ParserState string

const (
    StateInit ParserState = "init"

    StateError ParserState = "error"
    StateDone  ParserState = "done"
)

var ErrorRequestInErrorState = fmt.Errorf("request is in error state")

func (r *Request) ParseRequest(data []byte) (int, error) {
    read := 0
outer:
    for {
        switch r.State {
        case StateDone:
            return read, nil

        case StateInit:
            requestLine, n, err := parseRequestLine(data[read:])
            if err != nil {
                r.State = StateError
                return 0, err
            }

            if n == 0 {
                break outer
            }

            r.RequestLine = *requestLine
            read += n

            r.State = StateDone
            break outer
        case StateError:
            return 0, ErrorRequestInErrorState
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
