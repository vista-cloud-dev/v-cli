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

# Regenerate the umbrella registry (v-cli-platform.md §5) → dist/v-registry.json.
registry:
	UPDATE_GOLDEN=1 go test -run Registry .

# Drift gate: fail if dist/v-registry.json is stale vs the pinned domains' contracts.
check-registry:
	go test -run Registry .
