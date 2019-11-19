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
	@go build $(BUILD_FLAGS) -o bin/truchaincli cmd/truchaincli/*.go

build_daemon:
	@go build $(BUILD_FLAGS) -o bin/truchaind cmd/truchaind/*.go

build-linux:
	GOOS=linux GOARCH=amd64 go build -mod=readonly $(BUILD_FLAGS) -o build/truchaind cmd/truchaind/*.go
	GOOS=linux GOARCH=amd64 go build -mod=readonly $(BUILD_FLAGS) -o build/truchaincli cmd/truchaincli/*.go

doc:
	@echo "--> Wait a few seconds and visit http://localhost:6060/pkg/github.com/TruStory/truchain/"
	godoc -http=:6060

export:
	@bin/truchaind export

create-wallet:
	bin/truchaincli keys add validator --home ~/.octopus

init:
	rm -rf ~/.truchaind
	bin/truchaind init trunode
	bin/truchaind add-genesis-account $(shell bin/truchaincli keys show validator -a --home ~/.octopus) 10000000000utru
	sed -i -e 's/registrar.*/registrar": "$(shell bin/truchaincli keys show validator -a --home ~/.octopus)",/' ~/.truchaind/config/genesis.json
	sed -i -e 's/community_admins.*/community_admins": ["$(shell bin/truchaincli keys show validator -a --home ~/.octopus)"]/' ~/.truchaind/config/genesis.json
	sed -i -e 's/claim_admins.*/community_admins": ["$(shell bin/truchaincli keys show validator -a --home ~/.octopus)"]/' ~/.truchaind/config/genesis.json
	sed -i -e 's/staking_admins.*/staking_admins": ["$(shell bin/truchaincli keys show validator -a --home ~/.octopus)"],/' ~/.truchaind/config/genesis.json
	sed -i -e 's/slash_admins.*/slash_admins": ["$(shell bin/truchaincli keys show validator -a --home ~/.octopus)"],/' ~/.truchaind/config/genesis.json
	sed -i -e 's/reward_broker_address.*/reward_broker_address": "$(shell bin/truchaincli keys show validator -a --home ~/.octopus)"/' ~/.truchaind/config/genesis.json
	bin/truchaind gentx --name=validator --amount 10000000000utru --home-client ~/.octopus
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
	bin/truchaind start --inv-check-period 10 --log_level "main:info,state:info,*:error,app:info,account:info,trubank2:info,claim:info,community:info,truslashing:info,trustaking:info"

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
