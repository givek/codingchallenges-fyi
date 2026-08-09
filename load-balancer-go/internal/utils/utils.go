package utils

import (
	"log/slog"
	"net/http"
)

func ErrAttr(err error) slog.Attr {
	return slog.Any("error", err)
}

func LogReqDetails(r *http.Request, logger *slog.Logger) {
	attrs := []any{
		slog.String("remote_addr", r.RemoteAddr),
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
		slog.String("proto", r.Proto),
		slog.String("host", r.Host),
		slog.String("user_agent", r.UserAgent()),
	}

	acceptHeader := r.Header.Get("Accept")
	if acceptHeader != "" {
		attrs = append(attrs, slog.String("accept", acceptHeader))
	}

	logger.Info("Received request", attrs...)
}
