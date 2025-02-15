package main

import (
	"io"
	"log"
	"net"
)

func HandleConnection(conn net.Conn) {
	log.Printf("Client connected: %s\n", conn.RemoteAddr().String())
	defer conn.Close()

	for {
		buf := make([]byte, 1024)

		_, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Printf("%s: connection closed", conn.RemoteAddr().String())
				break
			}

			log.Printf("%s: error reading from connection: %s", conn.RemoteAddr().String(), err)
		}

		res := ParseRESP(buf)
		log.Println(string(res))
		if err != nil {
			conn.Write([]byte("+Incorrect\rCommand\r\n"))
		} else {
			conn.Write(res)
		}
	}
}
