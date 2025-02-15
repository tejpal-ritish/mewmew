package main

import (
	"bufio"
	"fmt"
	"strings"
)

func EncodeInput(input string) []byte {
	input = strings.TrimSpace(input)
	arr := strings.Split(input, " ")

	buf := ""
	size := len(arr)

	buf += fmt.Sprintf("*%d\r\n", size)
	for _, word := range arr {
		_buf := ""
		_buf += fmt.Sprintf("$%d\r\n", len(word))
		_buf += fmt.Sprintf("%s\r\n", word)

		buf += _buf
	}

	return []byte(buf)
}

func DecodeMessage(buf []byte) {
	reader := bufio.NewReader(strings.NewReader(string(buf)))

	for {
		b, err := reader.ReadByte()
		if err != nil {
			break
		}

		if b == '\r' {
			reader.ReadByte()
			fmt.Println()
		}

		fmt.Print(string(b))
	}

	fmt.Println()
}