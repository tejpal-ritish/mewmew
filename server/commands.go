package main

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

var valid_commands = map[string]struct{}{
	"get":  {},
	"set":  {},
	"del":  {},
	"ping": {},
	"bgsave": {},
}

type Command struct {
	Type string
	Args []string
}

func (c *Command) ExecCommand() (any, error) {
	switch c.Type {
	case "set":
		if len(c.Args) < 2 {
			return nil, errors.New("key or value cannot be empty")
		}

		key := c.Args[0]
		val := c.Args[1]

		cache[key] = val

		log := fmt.Sprintf("SET %s %s", key, val)

		if len(c.Args) == 3 {
			exp := c.Args[2]

			log += fmt.Sprintf(" EXPIRY %s", exp)

			go func(exp string) {
				secs, _ := strconv.Atoi(exp)
				time.Sleep(time.Duration(secs) * time.Second)

				log = fmt.Sprintf("DEL %s", key)
				AppendToAOF(log)

				delete(cache, key)
			}(exp)
		}

		AppendToAOF(log)
		return nil, nil

	case "get":
		key := c.Args[0]

		val, exists := cache[key]
		if !exists {
			return nil, errors.New("key does not exist")
		}

		return val, nil

	case "del":
		key := c.Args[0]
		delete(cache, key)

		log := fmt.Sprintf("DEL %s", key)
		AppendToAOF(log)

	case "ping":
		return "ping", nil
	}

	return nil, nil
}
