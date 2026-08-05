package main

import (
	"log"

	"github.com/givek/codingchallenges-fyi/load-balancer-go/internal/lb"
)

func main() {
	servers := []*lb.Server{
		lb.NewServer(8080, true),
		lb.NewServer(8081, true),
		lb.NewServer(8082, true),
	}

	lb := lb.NewLoadBalancer(servers)

	if err := lb.Start(); err != nil {
		log.Fatal(err)
	}
}
