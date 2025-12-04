package reader

import (
    "github.com/stretchr/testify/require"
    "io"
    "net/http"
    "strings"
    "testing"
)

type chunkReader struct {
    data            string
    numBytesPerRead int
    pos             int
}

func TestRequestEmptyLineParse(t *testing.T) {
    const emptyRequest = "GET / HTTP/1.1\r\nHost: localhost:42069 \r\nUser-Agent:curl/7.81.0\r\nAccept: */*\r\n\r\n"

    reader := &chunkReader{
        data:            emptyRequest,
        numBytesPerRead: 3,
    }

    r, err := RequestFromReader(reader)
    require.NoError(t, err)
    require.NotNil(t, r)

    require.Equal(t, http.MethodGet, r.RequestLine.Method)
    require.Equal(t, "/", r.RequestLine.RequestTarget)
    require.Equal(t, "1.1", r.RequestLine.HttpVersion)

    const badCoffeeRequest = "/coffee HTTP/1.1\r\nHost: localhost:42069 \r\nUser-Agent:curl/7.81.0\r\nAccept: */*\r\n\r\n"

    _, err = RequestFromReader(strings.NewReader(badCoffeeRequest))
    require.Error(t, err)
}

func TestRequestLineParseCoffeeEndpoint(t *testing.T) {
    const coffeeRequest = "GET /coffee HTTP/1.1\r\nHost: localhost:42069 \r\nUser-Agent:curl/7.81.0\r\nAccept: */*\r\n\r\n"

    reader := &chunkReader{
        data:            coffeeRequest,
        numBytesPerRead: 1,
    }

    r, err := RequestFromReader(reader)
    require.NoError(t, err)
    require.NotNil(t, r)

    require.Equal(t, http.MethodGet, r.RequestLine.Method)
    require.Equal(t, "/coffee", r.RequestLine.RequestTarget)
    require.Equal(t, "1.1", r.RequestLine.HttpVersion)
}

func TestBadRequest(t *testing.T) {
    const badCoffeeRequest = "/coffee HTTP/1.1\r\nHost: localhost:42069 \r\nUser-Agent:curl/7.81.0\r\nAccept: */*\r\n\r\n"

    _, err := RequestFromReader(strings.NewReader(badCoffeeRequest))
    require.Error(t, err)
}

func TestRequestParsingWithHeadersShouldReturnHeaders(t *testing.T) {
    const emptyRequest = "GET / HTTP/1.1\r\nHost: localhost:42069 \r\nUser-Agent:curl/7.81.0\r\nAccept: */*\r\n\r\n"
    reader := &chunkReader{
        data:            emptyRequest,
        numBytesPerRead: 3,
    }

    r, err := RequestFromReader(reader)

    require.NoError(t, err)
    require.NotNil(t, r)

    host, _ := r.Headers.TryGet("host")
    userAgent, _ := r.Headers.TryGet("user-agent")
    accept, _ := r.Headers.TryGet("accept")

    require.Equal(t, "localhost:42069", host)
    require.Equal(t, "curl/7.81.0", userAgent)
    require.Equal(t, "*/*", accept)
}

func TestRequestParsingWithMalformedHeadersShouldReturnErr(t *testing.T) {
    const invalidRequest = "GET / HTTP/1.1\r\nHost : localhost:42069\r\n\r\n"
    reader := &chunkReader{
        data:            invalidRequest,
        numBytesPerRead: 3,
    }

    r, err := RequestFromReader(reader)
    require.Error(t, err)
    require.Nil(t, r)
}

func TestRequestParsingWithBodyShouldReturnBody(t *testing.T) {
    reader := &chunkReader{
        data: "POST /submit HTTP/1.1\r\n" +
            "Host: localhost:42069\r\n" +
            "Content-Length: 13\r\n" +
            "\r\n" +
            "hello world!\n",
        numBytesPerRead: 3,
    }

    r, err := RequestFromReader(reader)

    require.NoError(t, err)
    require.NotNil(t, r)
    require.Equal(t, "hello world!\n", r.Body)
}

func (cr *chunkReader) Read(p []byte) (n int, err error) {
    if cr.pos >= len(cr.data) {
        return 0, io.EOF
    }

    endIndex := cr.pos + cr.numBytesPerRead
    if endIndex > len(cr.data) {
        endIndex = len(cr.data)
    }

    n = copy(p, cr.data[cr.pos:endIndex])
    cr.pos += n
    if n > cr.numBytesPerRead {
        n = cr.numBytesPerRead
        cr.pos -= n - cr.numBytesPerRead
    }

    return n, nil
}
