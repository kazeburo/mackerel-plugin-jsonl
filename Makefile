VERSION=0.0.1
LDFLAGS=-ldflags "-w -s -X main.version=${VERSION}"
all: mackerel-plugin-jsonl

.PHONY: mackerel-plugin-jsonl

mackerel-plugin-jsonl: cmd/mackerel-plugin-jsonl/*.go
	go build $(LDFLAGS) -o mackerel-plugin-jsonl cmd/mackerel-plugin-jsonl/*.go

linux: cmd/mackerel-plugin-jsonl/*.go
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o mackerel-plugin-jsonl cmd/*.go

check:
	go test -v ./...


