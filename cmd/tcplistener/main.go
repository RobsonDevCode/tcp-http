package main

import (
	"fmt"
	"log"
	"net"
	"tcp-http/Internal/reader"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("error creating lister: %v", err)
	}

	fmt.Printf("Listening...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("error opening connetion: %v \n", err)
		}

		request, err := reader.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("could not request from reader: %v", err)
		}

		fmt.Printf("Request line: \n")
		fmt.Printf("- Method: %s\n", request.RequestLine.Method)
		fmt.Printf("- Target: %s\n", request.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", request.RequestLine.HttpVersion)
		request.Headers.ForEach(func(n, v string) {
			fmt.Printf("- %s: %s\n", n, v)
		})

		fmt.Printf("Body: \n")
		fmt.Printf("%s \n", request.Body)
	}
}
