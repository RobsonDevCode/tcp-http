package request

import (
    "bytes"
    "fmt"
    "strings"
    "tcp-http/Internal/constants/seperatorConstants"
)

type RequestLine struct {
    HttpVersion   string
    RequestTarget string
    Method        string
}

func (r *RequestLine) ValidHttp(version string) (bool, *string) {
    parsedVersion := strings.TrimPrefix(version, "HTTP/")
    if parsedVersion == "" {
        return false, nil
    }

    return true, &parsedVersion
}

func parseRequestLine(b []byte) (*RequestLine, int, error) {
    idx := bytes.Index(b, seperatorConstants.RNSEPERATOR)
    if idx == -1 {
        return nil, 0, nil
    }

    startLine := b[:idx]
    read := idx + len(seperatorConstants.RNSEPERATOR)

    parts := bytes.Split(startLine, []byte(" "))

    if len(parts) != 3 {
        return nil, 0, fmt.Errorf("malformed request line")
    }

    requestLine := &RequestLine{
        Method:        string(parts[0]),
        RequestTarget: string(parts[1]),
    }

    valid, parsedVersion := requestLine.ValidHttp(string(parts[2]))

    if !valid || parsedVersion == nil {
        return nil, 0, fmt.Errorf("unsupported http version, only supports version 1.1")
    }

    requestLine.HttpVersion = *parsedVersion
    return requestLine, read, nil
}
