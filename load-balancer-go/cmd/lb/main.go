package main

import (
	"log"

	"github.com/givek/codingchallenges-fyi/load-balancer-go/internal/lb"
	"github.com/givek/codingchallenges-fyi/load-balancer-go/internal/server"
)

func main() {
	servers := []*server.Server{
		server.NewServer(8080, true),
		server.NewServer(8081, true),
		server.NewServer(8082, true),
	}

	lb := lb.NewLoadBalancer(servers, 80)

	if err := lb.Start(); err != nil {
		log.Fatal(err)
	}
}
