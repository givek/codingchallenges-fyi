package server

import "sync"

type Server struct {
	// TODO: Maybe one of the servers is not on the local network?
	// 	- Is this even valid concern?
	Port   int
	active bool
	mu     sync.RWMutex
}

func NewServer(port int, active bool) *Server {
	return &Server{
		Port:   port,
		active: active,
	}
}

func (s *Server) IsActive() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.active
}

func (s *Server) SetActive(active bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.active = active
}
