.PHONY: build run clean test lint format install release docker-build docker-run

VERSION  := $(shell cat .version 2>/dev/null || echo "dev")
COMMIT   := $(shell git rev-parse --short=8 HEAD 2>/dev/null || echo "unknown")
LDFLAGS  := -s -w
LDFLAGS  += -X main.Version=$(VERSION)
LDFLAGS  += -X main.Commit=$(COMMIT)
LDFLAGS  += -X main.PublisherName=rkriad585
LDFLAGS  += -X main.PublisherEmail=rkriad585@gmail.com

BINARY   := neovector
GOFLAGS  := -trimpath
OUTPUT   := $(BINARY)

build:
	go build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(OUTPUT) .

run: build
	./$(OUTPUT) $(ARGS)

clean:
	rm -f $(OUTPUT)
	rm -f $(OUTPUT).exe
	rm -rf bin/

test:
	go test ./... -v -count=1

lint:
	go vet ./...

format:
	go fmt ./...

install:
	go install $(GOFLAGS) -ldflags "$(LDFLAGS)" .

release: clean
	GOOS=linux   GOARCH=amd64 go build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o bin/$(BINARY)-linux-amd64 .
	GOOS=linux   GOARCH=arm64 go build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o bin/$(BINARY)-linux-arm64 .
	GOOS=darwin  GOARCH=amd64 go build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o bin/$(BINARY)-darwin-amd64 .
	GOOS=darwin  GOARCH=arm64 go build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o bin/$(BINARY)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o bin/$(BINARY)-windows-amd64.exe .
	GOOS=windows GOARCH=arm64 go build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o bin/$(BINARY)-windows-arm64.exe .

docker-build:
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		-t rkriad585/$(BINARY):latest \
		-t rkriad585/$(BINARY):$(VERSION) \
		.

docker-run: docker-build
	docker run --rm rkriad585/$(BINARY):latest $(ARGS)
