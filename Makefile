PACKAGES=$(shell go list ./...)

MODULES = backing category challenge game story vote

CHAIN_DIR = ./.chain

define \n


endef

benchmark:
	@go test -bench=. $(PACKAGES)

buidl: build

build: build_cli build_daemon

br: build_daemon run_daemon

bwr: build_daemon wipe_chain run_daemon

build_cli:
	go build -o bin/trucli cmd/trucli/main.go

build_daemon:
	go build -o bin/truchaind cmd/truchaind/main.go

init:
	cp -R ./dev/empty_truchain $(CHAIN_DIR)

run_daemon:
	bin/truchaind --home $(CHAIN_DIR) start

wipe_chain:
	bin/truchaind --home $(CHAIN_DIR) unsafe-reset-all

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
