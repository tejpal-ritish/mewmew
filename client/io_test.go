package main

import "testing"

func TestEncodeInput(t *testing.T) {
	input := "SET key 2"
	output := EncodeInput(input)

	t.Log(string(output))
}
