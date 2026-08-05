package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/givek/codingchallenges-fyi/load-balancer-go/internal/utils"
)

func main() {
	userArgs := os.Args[1:]

	if len(userArgs) < 1 {
		log.Fatal("Server port is requied argument")
		return
	}

	port, err := strconv.Atoi(strings.TrimSpace(userArgs[0]))
	if err != nil {
		log.Fatal("Expected an integer port number got: ", port, err)
		return
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		utils.LogReqDetails(r)

		w.Write([]byte(fmt.Sprintf("Hello From Backend Server running on port: %v\n", port)))

		log.Println("Replied with a hello message")
	})

	if err := http.ListenAndServe(fmt.Sprintf(":%v", port), mux); err != nil {
		log.Fatal(err)
	}
}
