MODULE := github.com/maz/k8s-mini-explorer
BINARY := k8s-explorer
LDFLAGS := -s -w

.PHONY: build test run serve docker tidy clean

build:
	go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY) ./cmd/k8s-explorer

test:
	go test ./... -count=1 -cover

run: build
	./bin/$(BINARY) --help

serve: build
	./bin/$(BINARY) serve --port 8080

docker:
	docker compose build

tidy:
	go mod tidy

clean:
	rm -rf bin/ dist/
