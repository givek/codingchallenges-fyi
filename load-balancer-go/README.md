# Load Balancer in Go

Small learning project for building an HTTP load balancer in Go.

## Current Check

- `go test ./...` passes.
- `go test -race ./...` passes.
- `go build ./cmd/lb` and `go build ./cmd/ts` complete in the current environment.

## Priority Checklist

### High Priority

- [x] Return HTTP error responses from load balancer error paths.
  - In `internal/lb/lb.go`, `handleReq` no longer exits the whole process on request errors, which is good.
  - Error paths now write explicit responses instead of returning an implicit empty `200 OK`.
  - Current status codes: `503 Service Unavailable` when no backend is available, `500 Internal Server Error` for unexpected internal setup errors, and `502 Bad Gateway` when the selected backend fails.

- [ ] Make the injected logger nil-safe.
  - `NewLoadBalancer` accepts `logger *slog.Logger`, so a caller can accidentally pass `nil`.
  - Either default to `slog.Default()` when `logger == nil`, or treat nil as invalid and return an error from the constructor.

- [ ] Add tests for failure behavior.
  - Done: test what happens when all servers are inactive.
  - Test what happens when a selected backend is unreachable.
  - Test all request error paths return useful HTTP status codes instead of an empty `200 OK`.

### Medium Priority

- [ ] Decide how you want to handle the upstream `Host` header.
  - `clonedReq.URL.Host` changes where the proxy connects, but `clonedReq.Host` controls the `Host` header sent to the backend.
  - Some backend applications route or validate requests based on `Host`.
  - Try logging `r.Host` in the backend and compare behavior with and without setting `clonedReq.Host`.

- [ ] Reuse one `http.Client` instead of creating a new one per request.
  - `http.Client` is safe for concurrent use.
  - Keeping one client on `LoadBalancer` lets Go reuse connections efficiently.

- [ ] Add proxy behavior tests using `httptest.Server`.
  - Check that requests rotate across active backends.
  - Check that inactive backends are skipped.
  - Check that path and query strings are preserved.
  - Check that response status codes and headers are copied back.

- [ ] Consider adding retries or failover.
  - A backend can pass the last health check and still fail during a real request.
  - One learning exercise: if the chosen backend fails, mark it inactive and try the next active backend once.

### Lower Priority

- [ ] Use lower-case error strings.
  - Go convention is error strings like `all servers are inactive`, not `All servers are inactive`.

- [ ] Consider replacing `time.Tick` with `time.NewTicker`.
  - `time.Tick` cannot be stopped.
  - A `time.NewTicker` is easier to clean up later if you add graceful shutdown with `context.Context`.

- [x] Add a `test` target to the `Makefile`.
  - `make test` now runs the full test suite with the race detector.

- [ ] Avoid committing generated binaries.
  - `bin/lb` and `bin/ts` look like build outputs.
  - Usually these go in `.gitignore`, and each developer builds them locally.

- [ ] Consider making backend servers configurable.
  - Right now the load balancer has hard-coded backend ports.
  - Later, you could read them from command-line flags, an environment variable, or a config file.

- [ ] Look at `net/http/httputil.ReverseProxy` after finishing the manual version.
  - Your current manual proxy code is good for learning.
  - After that, compare it with the standard library reverse proxy to see which edge cases it handles for you.
