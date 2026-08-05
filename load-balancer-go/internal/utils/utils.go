package utils

import (
	"log"
	"net/http"
)

func LogReqDetails(r *http.Request) {
	log.Println("Received request from ", r.RemoteAddr)
	log.Println(r.Method, r.URL.Path, r.Proto)
	log.Println("Host: ", r.Host)
	log.Println("User-Agent: ", r.UserAgent())

	acceptHeader := r.Header.Get("Accept")
	if acceptHeader != "" {
		log.Println("Accept: ", acceptHeader)
	}
}
