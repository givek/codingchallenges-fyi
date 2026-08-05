package server_test

import (
	"sync"
	"testing"

	"github.com/givek/codingchallenges-fyi/load-balancer-go/internal/server"
)

func TestServer_SetActiveAndIsActive(t *testing.T) {
	s := server.NewServer(9090, false)

	if s.IsActive() {
		t.Errorf("Expected server to be inactive initially")
	}

	s.SetActive(true)
	if !s.IsActive() {
		t.Errorf("Expected server to be active after setting to true")
	}

	s.SetActive(false)
	if s.IsActive() {
		t.Errorf("Expected server to be inactive after setting to false")
	}
}

func TestServer_ConcurrentAccess(t *testing.T) {
	s := server.NewServer(8080, false)

	var wg sync.WaitGroup
	workers := 1_000
	iterations := 1_000

	for i := range workers {
		wg.Add(2)

		go func(val bool) {
			defer wg.Done()
			for range iterations {
				s.SetActive(val)
			}
		}(i%2 == 0)

		go func() {
			defer wg.Done()
			for range iterations {
				_ = s.IsActive()
			}
		}()
	}

	wg.Wait()
}
