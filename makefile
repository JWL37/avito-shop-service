lint:
	./run_lint.sh

test:
	go test ./... -coverprofile=coverage.out -coverpkg=./internal/... 
	go tool cover -func=coverage.out
