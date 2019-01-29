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

doc:
	@echo "--> Wait a few seconds and visit http://localhost:6060/pkg/github.com/TruStory/truchain/"
	godoc -http=:6060

debug:
	bin/truchaind --home $(CHAIN_DIR) --log_level "app:info,story:info,backing:info,challenge:info,game:info,users:info,vote:info,*:error" start

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

go_test:
	@go test $(PACKAGES)

set_registrar:
	mkdir -p $(HOME)/.truchaind
	cp $(CHAIN_DIR)/registrar.key $(HOME)/.truchaind/registrar.key

set_seed_data:
	mkdir -p $(HOME)/.truchaind
	cp $(CHAIN_DIR)/bootstrap.csv $(HOME)/.truchaind/bootstrap.csv

set_env_vars:
	mkdir -p $(HOME)/.truchaind
	cp .env.example $(CHAIN_DIR)/.env
	cp .env.example $(HOME)/.truchaind/.env

test: set_registrar set_seed_data set_env_vars go_test

test_cover: set_registrar set_seed_data set_env_vars
	@go test $(PACKAGES) -v -timeout 30m -race -coverprofile=coverage.txt -covermode=atomic

update_deps:
	@echo "--> Running dep ensure"
	@rm -rf .vendor-new
	@dep ensure -v

.PHONY: benchmark buidl build check dep_graph test test_cover update_deps
