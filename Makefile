PACKAGES=$(shell go list ./...)

MODULES = argument backing category challenge expiration stake story

CHAIN_DIR = ./.chain

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')

ldflags = -X github.com/TruStory/truchain/vendor/github.com/cosmos/cosmos-sdk/version.Name=truchaind \
	-X github.com/TruStory/truchain/vendor/github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/TruStory/truchain/vendor/github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT)

BUILD_FLAGS := -ldflags '$(ldflags)'

define \n


endef

benchmark:
	@go test -bench=. $(PACKAGES)

buidl: build

build: build_daemon

build_cli:
	@go build -o bin/trucli cmd/trucli/main.go

build_daemon:
	@go build $(BUILD_FLAGS) -o bin/truchaind cmd/truchaind/main.go

build-linux:
	GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o build/truchaind cmd/truchaind/main.go

doc:
	@echo "--> Wait a few seconds and visit http://localhost:6060/pkg/github.com/TruStory/truchain/"
	godoc -http=:6060

export:
	bin/truchaind --home $(CHAIN_DIR) export

init:
	bin/truchaind --home $(CHAIN_DIR) init trunode

reset:
	bin/truchaind --home $(CHAIN_DIR) unsafe-reset-all

restart: build_daemon reset start

start:
	bin/truchaind --home $(CHAIN_DIR) --log_level "main:info,state:info,*:error,app:info,argument:info,backing:info,category:info,challenge:info,expiration:info,stake:info,stories:info" start

check:
	@echo "--> Running golangci"
	@golangci-lint run --tests=false --skip-files=\\btest_common.go

dep_graph: ; $(foreach dir, $(MODULES), godepgraph -s -novendor github.com/TruStory/truchain/x/$(dir) | dot -Tpng -o x/$(dir)/dep.png${\n})

install_tools_macos:
	brew install dep && brew upgrade dep
	brew install golangci/tap/golangci-lint
	brew upgrade golangci/tap/golangci-lint

go_test:
	@go test $(PACKAGES)

test: go_test

test_cover:
	@go test $(PACKAGES) -v -timeout 30m -race -coverprofile=coverage.txt -covermode=atomic
	@go tool cover -html=coverage.txt

update_deps:
	@echo "--> Running dep ensure"
	@rm -rf .vendor-new
	@dep ensure -v -vendor-only

########################################
### Local validator nodes using docker and docker-compose

build-docker-truchaindnode:
	$(MAKE) -C networks/local

# Run a 4-node testnet locally
localnet-start: localnet-stop
	@if ! [ -f build/node0/truchaind/config/genesis.json ]; then docker run --rm -v $(CURDIR)/build:/truchaind:Z trustory/truchaindnode testnet --v 4 -o . --starting-ip-address 192.168.10.2 ; fi
	docker-compose up -d

# Stop testnet
localnet-stop:
	docker-compose down

########################################

.PHONY: benchmark buidl build check dep_graph test test_cover update_deps \
build-docker-truchaindnode localnet-start localnet-stop
