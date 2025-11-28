package main

import (
	"fmt"
	"log"
	"net"
	"os"
	stringextensions "tcp-http/Internal"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil{
		log.Fatalf("error creating lister: %v", err)
	}

	fmt.Printf("Listening...")
	if err := sendMessage(); err != nil{
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil{
			log.Fatalf("error opening connetion: %v \n", err)
		}
		firstLine := true
		for line := range stringextensions.GetLines(conn){
			if(firstLine){
				fmt.Printf("\nread: %s\n", line)
				firstLine = false
			}
			fmt.Printf("read: %s\n", line)
		}
	}
}

func sendMessage() error {
	connection, err := net.Dial("tcp", ":42069")
	if err != nil{
		return fmt.Errorf("error dialing tcp connection: %v", err)
	}
	defer connection.Close()

	message, err := os.ReadFile("message.txt")
	if err != nil{
		return fmt.Errorf("error couldnt open file: %v", err)
	}

	_, err = connection.Write(message)
	if err != nil  {
		return fmt.Errorf("error writing message: %v", err)
	}

	fmt.Print("sending message...")
	return nil
}

