package main

import (
	"testing"
)

func TestParseRESP(t *testing.T) {
	input := []byte("*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$1\r\n2\r\n")

	output := string(ParseRESP(input))
	if output[:4] == "-ERR" {
		t.Error(output)
		return
	}
	t.Log(string(output))
}

func TestResponseEncoder(t *testing.T) {
	val := 2
	output := ResponseEncoder(val)

	str := string(output)

	t.Log(str)
}