# Run a new redis container
run_redis: kill_redis
	docker run --name redis-local-instance -p 6379:6379 -d redis

# Kill any running redis container
kill_redis:
	 docker container rm -f redis-local-instance

run_go_tests:
	go mod tidy
	go test

# Spin up a new redis container, execute all Go tests and then kill the container
test: run_redis run_go_tests kill_redis

# Compile and run the Go code
go_run:
	go mod tidy
	go build -o server
	./server

run: run_redis go_run
