package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

func logReqDetails(r *http.Request) {
	log.Println("Received request from ", r.RemoteAddr)
	log.Println(r.Method, r.URL.Path, r.Proto)
	log.Println("Host: ", r.Host)
	log.Println("User-Agent: ", r.UserAgent())

	acceptHeader := r.Header.Get("Accept")
	if acceptHeader != "" {
		log.Println("Accept: ", acceptHeader)
	}
}

func dummyServer(port int) {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logReqDetails(r)

		w.Write([]byte(fmt.Sprintf("Hello From Backend Server running on port: %v\n", port)))

		log.Println("Replied with a hello message")
	})

	if err := http.ListenAndServe(fmt.Sprintf(":%v", port), mux); err != nil {
		log.Fatal(err)
	}
}

func main() {
	userArgs := os.Args[1:]

	if len(userArgs) > 1 {
		flag := strings.TrimSpace(userArgs[0])

		const testServerFlag = "-ts"

		if flag == testServerFlag {
			port, err := strconv.Atoi(strings.TrimSpace(userArgs[1]))
			if err != nil {
				// TODO: There needs to be a helpful messsage
				// 	along with the err.
				log.Fatal(err)
			}
			dummyServer(port)
			return
		} else {
			log.Fatal("Unsupported Flag: ", flag)
		}
	}

	port := ":80"

	mux := http.NewServeMux()

	servers := []int{8080, 8081}
	// TODO: Golang maps are not conc safe??
	serversMap := map[int]bool{
		8080: true,
		8081: true,
	}
	currServerIdx := 0

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logReqDetails(r)

		httpClient := &http.Client{
			Timeout: time.Second * 5,
		}

		serverPort := servers[currServerIdx]
		currServerIdx = (currServerIdx + 1) % len(servers)

		newBaseURL, err := url.Parse(fmt.Sprintf("http://localhost:%v", serverPort))
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

	go func() {
		for t := range time.Tick(time.Second * 2) {
			for k, _ := range serversMap {
				res, err := http.Get(fmt.Sprintf("http://localhost:%v/health-check", k))
				if err != nil {
					// TODO: log.Fatal is bad, we need something better with more context
					log.Println("Got ERR HEALTH CHECK", err)
					serversMap[k] = false
					servers = slices.DeleteFunc(servers, func(p int) bool {
						return p == k
					})
					continue
				}

				if res.StatusCode != http.StatusOK {
					serversMap[k] = false
					servers = slices.DeleteFunc(servers, func(p int) bool {
						return p == k
					})
				} else if serversMap[k] == false {
					serversMap[k] = true
					servers = append(servers, k)
				}
			}

			log.Println("Req done: ", servers, t)
		}
	}()

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal(err)
	}
}
