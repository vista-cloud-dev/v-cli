# v — the VistA developer-tools umbrella CLI.
BIN     ?= v
GOFLAGS := -trimpath

build:
	go build $(GOFLAGS) -o dist/$(BIN) .

test:
	go test -race -cover ./...

vet:
	go vet ./...

lint:
	golangci-lint run ./...

check: vet lint test build

clean:
	rm -rf dist/
