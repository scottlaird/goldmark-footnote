.PHONY: test lint

test:
	go test -race ./...

lint:
	golangci-lint run -c .golangci.yml ./...
