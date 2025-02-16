lint:
	./run_lint.sh

test:
	go test ./internal/... -coverprofile=coverage.out -coverpkg=./... 
	go tool cover -func=coverage.out
	
intergation:
	cd ./tests/integration && go test -v ./...

