package lb

import (
	"log/slog"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/givek/codingchallenges-fyi/load-balancer-go/internal/server"
)

func TestLoadBalancerConcurrentRoundRobin(t *testing.T) {
	servers := []*server.Server{
		server.NewServer(8080, true),
		server.NewServer(8081, true),
		server.NewServer(8082, true),
		server.NewServer(8083, true),
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	lb := NewLoadBalancer(servers, 80, 15*time.Second, logger)

	var wg sync.WaitGroup
	workers := 10_000 * len(servers)

	var (
		count_8080 atomic.Int64
		count_8081 atomic.Int64
		count_8082 atomic.Int64
		count_8083 atomic.Int64
	)

	wg.Add(workers)

	for range workers {
		go func() {
			defer wg.Done()
			s, err := lb.getNextServer()
			if err != nil {
				t.Errorf("Failed to get next server - %v\n", err)
			}

			if s.Port == 8080 {
				count_8080.Add(1)
			}
			if s.Port == 8081 {
				count_8081.Add(1)
			}
			if s.Port == 8082 {
				count_8082.Add(1)
			}
			if s.Port == 8083 {
				count_8083.Add(1)
			}
		}()
	}

	wg.Wait()

	expectedReqPerServer := int64(workers) / int64(len(servers))
	if count_8080.Load() != expectedReqPerServer ||
		count_8081.Load() != expectedReqPerServer ||
		count_8082.Load() != expectedReqPerServer ||
		count_8083.Load() != expectedReqPerServer {
		t.Errorf("No round robin\n")
	}
}
