package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
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

func dummyServer() {
	port := ":8080"

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logReqDetails(r)

		w.Write([]byte("Hello From Backend Server\n"))

		log.Println("Replied with a hello message")
	})

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal(err)
	}
}

func main() {
	userArgs := os.Args[1:]

	if len(userArgs) > 0 {
		flag := strings.TrimSpace(userArgs[0])

		const testServerFlag = "-test-server"

		if flag == testServerFlag {
			dummyServer()
			return
		} else {
			log.Fatal("Unsupported Flag: ", flag)
		}
	}

	port := ":80"

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logReqDetails(r)

		httpClient := &http.Client{
			Timeout: time.Second * 5,
		}

		newBaseURL, err := url.Parse("http://localhost:8080")
		if err != nil {
			log.Fatal(err)
		}

		// TODO: Need to study r.Host vs r.URL.Host
		// - Is this a hack?
		// - Would it be better to set a custom header: X-Client-Host?

		clonedReq := r.Clone(context.Background())
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

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal(err)
	}
}
