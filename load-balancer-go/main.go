package main

import (
	"log"
	"net/http"
)

func main() {
	port := ":80"

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request from ", r.RemoteAddr)
		log.Println(r.Method, r.URL.Path, r.Proto)
		log.Println("Host: ", r.Host)
		log.Println("User-Agent: ", r.UserAgent())

		acceptHeader := r.Header.Get("Accept")
		if acceptHeader != "" {
			log.Println("Accept: ", acceptHeader)
		}
	})

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal(err)
	}
}
