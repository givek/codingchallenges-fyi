package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/givek/codingchallenges-fyi/load-balancer-go/internal/utils"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	userArgs := os.Args[1:]

	if len(userArgs) < 1 {
		logger.Error("Server port is requied argument")
		os.Exit(1)
	}

	portStr := strings.TrimSpace(userArgs[0])
	port, err := strconv.Atoi(portStr)
	if err != nil {
		logger.Error(
			"Expected an integer port number",
			slog.String("portStr", portStr),
			utils.ErrAttr(err),
		)
		os.Exit(1)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		utils.LogReqDetails(r, logger)

		w.Write([]byte(fmt.Sprintf("Hello From Backend Server running on port: %v\n", port)))

		logger.Info("Replied with a hello message")
	})

	if err := http.ListenAndServe(fmt.Sprintf(":%v", port), mux); err != nil {
		logger.Error(
			"Failed to start test server",
			slog.Int("port", port),
			utils.ErrAttr(err),
		)
		os.Exit(1)
	}
}
