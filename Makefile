PACKAGES=$(shell go list ./...)

benchmark:
	@go test -bench=. $(PACKAGES)

buidl: build

build:
	go build -o bin/trucli cmd/trucli/main.go && go build -o bin/truchaind cmd/truchaind/main.go

test:
	@go test $(PACKAGES)

update_vendor_deps:
	@echo "--> Running dep ensure"
	@rm -rf .vendor-new
	@dep ensure -v

.PHONY: build test benchmark