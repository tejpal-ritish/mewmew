package main

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

// ParseRESP is used to process the validity of a user command request
// Currently only supports bulk string arrays validation
// Supported commands: GET, SET, DEL, PING
func ParseRESP(input []byte) []byte {
	_input := string(input)
	reader := bufio.NewReader(strings.NewReader(_input))

	// Check if bulk string array
	b, _ := reader.ReadByte()
	if b != '*' {
		return []byte("-ERR Protocol Error: not a valid bulk string\r\n")
	}

	// size of command array
	size, _ := reader.ReadByte()
	strsize, _ := strconv.ParseInt(string(size), 10, 64)

	// Read "\r\n"
	reader.ReadByte()
	reader.ReadByte()

	// Store command type and args
	c := &Command{}

	// Reading individual bulk strings
	for i := range strsize {
		if b, _ = reader.ReadByte(); b != '$' {
			return []byte("-ERR Protocol Error: not a valid bulk string\r\n")
		}

		_size := 0

		// Determining size to read for string with size > 9 bytes
		for {
			b, _ = reader.ReadByte()
			if b == '\r' {
				break
			}

			n, _ := strconv.Atoi(string(b))

			_size = _size*10 + n
		}

		// Read "\r\n"
		reader.ReadByte()

		// Read actual string content
		buf := make([]byte, 0, _size)
		for range _size {
			b, _ = reader.ReadByte()
			buf = append(buf, b)
		}

		// Validate command type
		if i == 0 {
			valid := ValidateCommandType(buf)
			if !valid {
				return []byte("-ERR Protocol Error: not a valid command type\r\n")
			}

			c.Type = strings.ToLower(string(buf)) 
			c.Args = make([]string, 0, strsize-1)
		} else {
			c.Args = append(c.Args, string(buf))
		}

		reader.ReadByte()
		reader.ReadByte()
	}

	val, err := c.ExecCommand()
	if err != nil {
		return []byte(fmt.Sprintf("-ERR: %s\r\n", err))
	}

	if val == "ping" {
		return []byte("+PONG\r\n")
	}

	if val != nil {
		return ResponseEncoder(val)
	}

	return []byte("+OK\r\n")
}

func ValidateCommandType(command []byte) bool {
	_, exists := valid_commands[strings.ToLower(string(command))]
	return exists
}

// For now enforces only int & string values
// TODO: fix single digit int value encoding
func ResponseEncoder(val any) []byte {
	res := ""

	switch v := val.(type) {
	case string:
		size := len(v)
		res += fmt.Sprintf("$%d\r\n%s\r\n", size, v)
	case int:
		strval := strconv.Itoa(v)
		size := len(strval)
		res += fmt.Sprintf("$%d\r\n%s\r\n", size, strval)
	}

	return []byte(res)
}
