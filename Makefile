# Create a Makefile for the project using the following command: lint, test, bench, vet.

.PHONY: lint test bench vet

# Prerequisites:
# 	- golangci-lint (https://golangci-lint.run/usage/install/#local-installation)
lint:
	golangci-lint run ./...

test:
	go test -v ./...

bench:
	go test -bench=. -benchmem -count=1 ./...

vet:
	go vet ./...
