package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:3000")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	for {
		fmt.Printf("> ")

		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Println("error reading input: ", err)
			continue
		}

		if input == "exit\n" {
			break
		}

		command := EncodeInput(input)
		conn.Write(command)

		// Write server response to the terminal
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("error reading from server: ", err)
			continue
		}

		DecodeMessage(buf[:n])
	}
}
