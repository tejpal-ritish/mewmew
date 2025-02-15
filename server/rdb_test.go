package main

import (
	"net"
	"testing"
)

func TestLoadSnapshot(t *testing.T) {
	s := Server{
		Clients: make(map[string]net.Conn),
	}

	err := s.LoadSnapshot()
	if err != nil {
		t.Error(err)
	}

	for key, val := range cache {
		t.Logf("key: %s, val: %v\n", key, val)
	}
}
