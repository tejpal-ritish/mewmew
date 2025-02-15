package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func (s *Server) HandlePersistence(stop chan struct{}) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := s.Snapshot(); err != nil {
				log.Println("Error in Snapshot:", err)
			}
		case <-stop:
			return
		}
	}
}

func (s *Server) Snapshot() error {
	err := os.MkdirAll("./data", 0755)
	if err != nil {
		return err
	}

	file, err := os.OpenFile("data/dump.rdb", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	s.mu.Lock()
	defer s.mu.Unlock()

	data := make([]byte, 0, 2048)
	data = append(data, []byte("#!RDB FILE\r\n")...)

	for key, val := range cache {
		data = append(data, []byte(fmt.Sprintf("%s\r%s\r\n", key, val))...)
	}

	err = gob.NewEncoder(file).Encode(data)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) LoadSnapshot() error {
	if _, err := os.Stat("data/dump.rdb"); os.IsNotExist(err) {
		log.Println("No existing dump file found, starting with an empty cache.")
		return nil
	}

	log.Println("Loading snapshot from dump file.")

	file, err := os.Open("data/dump.rdb")
	if err != nil {
		return err
	}

	defer file.Close()

	var data []byte
	err = gob.NewDecoder(file).Decode(&data)
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if !bytes.HasPrefix(data, []byte("#!RDB FILE\r\n")) {
		return fmt.Errorf("invalid dump file")
	}

	reader := bufio.NewReader(bytes.NewReader(data[12:]))
	for {
		key, err := reader.ReadBytes('\r')
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		key = key[:len(key)-1]

		val, err := reader.ReadBytes('\r')
		if err != nil {
			return err
		}
		val = val[:len(val)-1]

		cache[string(key)] = string(val)

		reader.ReadByte()
	}

	return nil
}
