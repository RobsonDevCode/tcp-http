package main

import (
	"fmt"
	"log"
	"net"
	"os"
	stringextensions "tcp-http/Internal/extensions"
	tcpsender "tcp-http/cmd/tcp"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil{
		log.Fatalf("error creating lister: %v", err)
	}

	fmt.Printf("Listening...")
	message, err := os.ReadFile("message.txt")
	if err != nil{
		log.Fatalf("error reading file: %v", err)
	}

	if err := tcpsender.SendMessage(message); err != nil{
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil{
			log.Fatalf("error opening connetion: %v \n", err)
		}
		firstLine := true
		for line := range stringextensions.GetLines(conn){
			if firstLine {
				fmt.Printf("\nread: %s\n", line)
				firstLine = false
				continue
			}
			fmt.Printf("read: %s\n", line)
		}
	}
}
