package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/givek/codingchallenges-fyi/load-balancer-go/internal/lb"
	"github.com/givek/codingchallenges-fyi/load-balancer-go/internal/server"
	"github.com/givek/codingchallenges-fyi/load-balancer-go/internal/utils"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	servers := []*server.Server{
		server.NewServer(8080, true),
		server.NewServer(8081, true),
		server.NewServer(8082, true),
	}

	lb := lb.NewLoadBalancer(servers, 80, time.Second*15, logger)

	if err := lb.Start(); err != nil {
		logger.Error(
			"Failed to start load balancer",
			utils.ErrAttr(err),
		)
		os.Exit(1)
	}
}
