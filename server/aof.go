package main

import (
	"os"
)

func AppendToAOF(msg string) error {
	file, err := os.OpenFile("data/appendonly.aof", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.WriteString(msg + "\r\n")
	if err != nil {
		return err
	}

	return nil
}
