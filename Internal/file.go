package stringextensions

import (
	"bytes"
	"io"
	"log"
	"strings"
)

func GetLines(file io.ReadCloser) <- chan string{
	out := make(chan string, 1)
	go func() {
		defer handleClose(file)
		defer close(out)
		var builder strings.Builder
		for {
			data := make([]byte, 8)
			chunkedText, err := file.Read(data)
			if err != nil{
				break
			}

			data = data[:chunkedText]
			if i := bytes.IndexByte(data, '\n'); i != -1 {
				builder.Write(data[:i])

				out <- builder.String()
				builder.Reset()
				data = data[i + 1:]
			}

			builder.WriteString(string(data))
		}

		result := builder.String()
		if len(result) != 0 {
			out <- result
		}
	}()

	return out
}

func handleClose(f io.ReadCloser) {
	err := f.Close()
	if err != nil {
		log.Fatalf("error while closing reader: %v", err)
	}
}
