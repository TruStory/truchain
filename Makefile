PACKAGES=$(shell GO111MODULE=on go list -mod=readonly ./...)

MODULES = argument backing category challenge expiration stake story

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=truchaind \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT)

BUILD_FLAGS := -ldflags '$(ldflags)'

define \n


endef

benchmark:
	@go test -mod=readonly -bench=. $(PACKAGES)

buidl: build

build: build_cli build_daemon

download:
	go mod download

build_cli:
	@go build -mod=readonly $(BUILD_FLAGS) -o bin/truchaincli cmd/truchaincli/main.go

build_daemon:
	@go build -mod=readonly $(BUILD_FLAGS) -o bin/truchaind cmd/truchaind/*.go

build-linux:
	GOOS=linux GOARCH=amd64 go build -mod=readonly $(BUILD_FLAGS) -o build/truchaind cmd/truchaind/*.go
	GOOS=linux GOARCH=amd64 go build -mod=readonly $(BUILD_FLAGS) -o build/truchaincli cmd/truchaincli/main.go

doc:
	@echo "--> Wait a few seconds and visit http://localhost:6060/pkg/github.com/TruStory/truchain/"
	godoc -http=:6060

export:
	bin/truchaind export

registrar:
	bin/truchaincli keys add registrar --home ~/.octopus

init:
	rm -rf ~/.truchaind
	bin/truchaind init trunode
	bin/truchaind add-genesis-account $(shell bin/truchaincli keys show registrar -a --home ~/.octopus) 1000000000trusteak

gentx:
	bin/truchaind gentx --name=registrar --amount 100000000trusteak --home-client ~/.octopus
	bin/truchaind collect-gentxs

install:
	@go install $(BUILD_FLAGS) ./cmd/truchaind
	@go install $(BUILD_FLAGS) ./cmd/truchaincli
	@echo "Installed truchaind and truchaincli ..."
	@truchaind version --long

reset:
	bin/truchaind unsafe-reset-all

restart: build_daemon reset start

start:
	bin/truchaind start --log_level "main:info,state:info,*:error,app:info,argument:info,backing:info,category:info,challenge:info,expiration:info,stake:info,stories:info"

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

version:
	@bin/truchaind version --long

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

.PHONY: benchmark buidl build build_cli build_daemon check dep_graph test test_cover update_deps \
build-docker-truchaindnode localnet-start localnet-stop
