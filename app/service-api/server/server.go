package server

import (
	"fmt"
	"sync"
)

type server struct {
	wg   sync.WaitGroup
	port int
}

func NewServer(port int) *server {
	return &server{
		port: port,
	}
}

func (s *server) Background(fn func()) {
	s.wg.Add(1)

	go func() {
		defer s.wg.Done()

		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("background error: %s", err)
			}
		}()

		fn()
	}()
}
