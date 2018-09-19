PACKAGES=$(shell go list ./...)

all: build test

buidl: build

build:
	go build -o bin/trucli cmd/trucli/main.go && go build -o bin/truchaind cmd/truchaind/main.go

test:
	@go test $(PACKAGES)

benchmark:
	@go test -bench=. $(PACKAGES)

.PHONY: all build test benchmark