package tcpsender

import (
    "fmt"
    "net"
)

func SendMessage(message []byte) error {
	connection, err := net.Dial("tcp", ":42069")
	if err != nil{
		return fmt.Errorf("error dialing tcp connection: %v", err)
	}
	defer connection.Close()

	_, err = connection.Write(message)
	if err != nil  {
		return fmt.Errorf("error writing message: %v", err)
	}

	fmt.Print("sending message...")
	return nil
}