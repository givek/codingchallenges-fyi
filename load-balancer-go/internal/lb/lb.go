package lb

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/givek/codingchallenges-fyi/load-balancer-go/internal/utils"
)

type Server struct {
	// TODO: Maybe one of the servers is not on the local network?
	// 	- Is this even valid concern?
	port   int
	active bool
}

func NewServer(port int, active bool) *Server {
	return &Server{port, active}
}

type LoadBalancer struct {
	servers []*Server
	idx     int
}

func NewLoadBalancer(servers []*Server) *LoadBalancer {
	return &LoadBalancer{servers: servers, idx: 0}
}

func (lb *LoadBalancer) healthCheckServers(interval time.Duration) {
	for t := range time.Tick(interval) {
		for _, s := range lb.servers {
			res, err := http.Get(fmt.Sprintf("http://localhost:%v/health-check", s.port))
			if err != nil {
				// TODO: log.Fatal is bad, we need something better with more context
				log.Println("Got ERR HEALTH CHECK", err)
				s.active = false
				continue
			}

			if res.StatusCode != http.StatusOK {
				s.active = false
			} else {
				s.active = true
			}
		}
		log.Println("Req done: ", t)
	}
}

func (lb *LoadBalancer) getNextServer() (*Server, error) {
	servers := lb.servers
	serversLen := len(servers)

	for i := 0; i < serversLen; i++ {
		idx := (lb.idx + i) % serversLen

		s := servers[idx]

		if s.active {
			lb.idx = (idx + 1) % serversLen
			return s, nil
		}
	}

	return nil, fmt.Errorf("All servers are inactive")
}

func (lb *LoadBalancer) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		utils.LogReqDetails(r)

		httpClient := &http.Client{
			Timeout: time.Second * 5,
		}

		server, err := lb.getNextServer()
		if err != nil {
			// TODO: should we return a 500 response or crash?
			log.Fatal(err)
		}

		newBaseURL, err := url.Parse(fmt.Sprintf("http://localhost:%v", server.port))
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
	})

	go lb.healthCheckServers(2 * time.Second)

	port := ":80"
	if err := http.ListenAndServe(port, mux); err != nil {
		return err
	}

	return nil
}
