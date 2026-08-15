.PHONY: build test run-ts stop-ts

build:
	go build -o bin/lb cmd/lb/main.go
	go build -o bin/ts cmd/ts/main.go

test:
	go test ./... --race

run-ts:
	$(MAKE) build
	@mkdir -p logs
	@echo "starting test server 8080"
	./bin/ts 8080 > ./logs/ts_8080.log 2>&1 &
	@echo "starting test server 8081"
	./bin/ts 8081 > ./logs/ts_8081.log 2>&1 &
	@echo "starting test server 8082"
	./bin/ts 8082 > ./logs/ts_8082.log 2>&1 &

stop-ts:
	@pkill -f "./bin/ts" || true
