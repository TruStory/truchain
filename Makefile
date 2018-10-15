PACKAGES=$(shell go list ./...)

MODULES = backing category challenge story

define \n


endef

benchmark:
	@go test -bench=. $(PACKAGES)

buidl: build

build:
	go build -o bin/trucli cmd/trucli/main.go && go build -o bin/truchaind cmd/truchaind/main.go

check:
	gometalinter ./...

dep_graph: ; $(foreach dir, $(MODULES), godepgraph -s -novendor github.com/TruStory/truchain/x/$(dir) | dot -Tpng -o x/$(dir)/dep.png${\n})

install_tools_macos:
	brew install dep && brew upgrade dep
	brew tap alecthomas/homebrew-tap && brew install gometalinter

test:
	@go test $(PACKAGES)

test_cover:
	go test $(PACKAGES) -v -timeout 30m -race -covermode=atomic

update_deps:
	@echo "--> Running dep ensure"
	@rm -rf .vendor-new
	@dep ensure -v

.PHONY: benchmark buidl build check dep_graph test test_cover update_deps
