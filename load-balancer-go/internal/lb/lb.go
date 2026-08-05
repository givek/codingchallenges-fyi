package lb

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/givek/codingchallenges-fyi/load-balancer-go/internal/server"
	"github.com/givek/codingchallenges-fyi/load-balancer-go/internal/utils"
)

type LoadBalancer struct {
	port    int
	servers []*server.Server
	idx     int
	mu      sync.Mutex
}

func NewLoadBalancer(servers []*server.Server, port int) *LoadBalancer {
	return &LoadBalancer{
		servers: servers,
		idx:     0,
		port:    port,
	}
}

func (lb *LoadBalancer) healthCheckServers(interval time.Duration) {
	client := http.Client{Timeout: 2 * time.Second}

	for _ = range time.Tick(interval) {
		for _, s := range lb.servers {
			res, err := client.Get(fmt.Sprintf("http://localhost:%v/health-check", s.Port))
			if err != nil {
				log.Printf("Health check failed for port %d: %v", s.Port, err)
				s.SetActive(false)
				continue
			}

			res.Body.Close()

			if res.StatusCode != http.StatusOK {
				s.SetActive(false)
			} else {
				s.SetActive(true)
			}
		}
	}
}

func (lb *LoadBalancer) getNextServer() (*server.Server, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	servers := lb.servers
	serversLen := len(servers)

	for i := 0; i < serversLen; i++ {
		idx := (lb.idx + i) % serversLen

		s := servers[idx]

		if s.IsActive() {
			lb.idx = (idx + 1) % serversLen
			return s, nil
		}
	}

	return nil, fmt.Errorf("All servers are inactive")
}
func (lb *LoadBalancer) handleReq(w http.ResponseWriter, r *http.Request) {
	utils.LogReqDetails(r)

	httpClient := &http.Client{
		Timeout: time.Second * 5,
	}

	server, err := lb.getNextServer()
	if err != nil {
		// TODO: should we return a 500 response or crash?
		log.Fatal(err)
	}

	newBaseURL, err := url.Parse(fmt.Sprintf("http://localhost:%v", server.Port))
	if err != nil {
		log.Fatal(err)
	}

	// TODO: Need to study r.Host vs r.URL.Host
	// - Is this a hack?
	// - Would it be better to set a custom header: X-Client-Host?

	// Use r.Context() intead of context.Background(), because if
	// the client disconnects or aborts their HTTP request halfway
	// through, our proxy server will still keep the connection open
	// and waste resources processing.
	clonedReq := r.Clone(r.Context())
	// clonedReq := r.Clone(context.Background())
	// clonedReq.URL = newBaseURL

	clonedReq.URL.Scheme = newBaseURL.Scheme
	clonedReq.URL.Host = newBaseURL.Host

	// clonedReq.Host = newBaseURL.Host

	clonedReq.RequestURI = ""

	serverRes, err := httpClient.Do(clonedReq)
	if err != nil {
		log.Fatal(err)
	}

	for headerName, headerVals := range serverRes.Header {
		for _, headerVal := range headerVals {
			w.Header().Add(headerName, headerVal)
		}
	}

	w.WriteHeader(serverRes.StatusCode)

	defer serverRes.Body.Close()
	if _, err := io.Copy(w, serverRes.Body); err != nil {
		log.Fatal(err)
	}
}

func (lb *LoadBalancer) Start() error {
	go lb.healthCheckServers(2 * time.Second)

	mux := http.NewServeMux()

	mux.HandleFunc("/", lb.handleReq)

	port := fmt.Sprintf(":%v", lb.port)

	log.Printf("Starting load balancer on port %s", port)

	if err := http.ListenAndServe(port, mux); err != nil {
		return err
	}

	return nil
}
