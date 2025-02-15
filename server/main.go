package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
)

var cache = map[string]any{}

type Server struct {
	Clients map[string]net.Conn
	mu      sync.Mutex
}

func main() {
	l, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal("error starting server: ", err)
	}

	log.Println("server ready to accept connections")

	s := Server{
		Clients: make(map[string]net.Conn),
	}

	s.LoadSnapshot()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	stop := make(chan struct{})
	var once sync.Once

	go func() {
		<-sig
		log.Println("shutting down server")

		l.Close()

		once.Do(func() {
			close(stop)
		})

		s.mu.Lock()
		for _, conn := range s.Clients {
			conn.Write([]byte("biebie\r\n"))
			conn.Close()
		}
		s.Clients = nil
		s.mu.Unlock()
	}()

	go s.HandlePersistence(stop)

	for {
		conn, err := l.Accept()
		if err != nil {
			select {
			case <-stop:
				log.Println("server shutting down gracefully")
				return
			default:
				log.Println("error accepting connection: ", err)
			}
		}

		s.mu.Lock()
		s.Clients[conn.RemoteAddr().String()] = conn
		s.mu.Unlock()

		go HandleConnection(conn)
	}
}
