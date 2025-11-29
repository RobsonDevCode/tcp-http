package reader

import (
    "errors"
    "fmt"
    "io"
    "tcp-http/Internal/contracts"
)

func RequestFromReader(reader io.Reader) (*contracts.Request, error) {
    request := &contracts.Request{
        State: contracts.StateInit,
    }
    buffer := make([]byte, 1024)
    bufferLength := 0
    for !request.IsDone() && !request.IsError() {
        n, err := reader.Read(buffer[bufferLength:])
        if err != nil {
            return nil, errors.Join(fmt.Errorf("error while reading from buffer at buffer index %v", bufferLength), err)
        }

        bufferLength += n
        readN, err := request.ParseRequest(buffer[:bufferLength])
        if err != nil {
            return nil, err
        }

        copy(buffer, buffer[readN:bufferLength])
        bufferLength -= readN
    }

    return request, nil
}
