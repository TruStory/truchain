![logo](./logo.jpg)

[![CircleCI](https://circleci.com/gh/TruStory/truchain.svg?style=svg&circle-token=0cea219dcac9bd6962a057d85c4a319613c6680e)](https://circleci.com/gh/TruStory/truchain)
[![codecov](https://codecov.io/gh/TruStory/truchain/branch/master/graph/badge.svg?token=jh3muAAEBs)](https://codecov.io/gh/TruStory/truchain)
[![Go Report Card](https://goreportcard.com/badge/github.com/trustory/truchain)](https://goreportcard.com/report/github.com/trustory/truchain)
[![](https://tokei.rs/b1/github/XAMPPRocky/tokei)](https://github.com/TruStory/truchain)
[![API Reference](https://godoc.org/github.com/TruStory/truchain/types?status.svg)](https://godoc.org/github.com/TruStory/truchain/types)

TruChain is the application-specific blockchain that powers [beta.trustory.io](https://beta.trustory.io/).

## Installation

1. Install Go by following the [official docs](https://golang.org/doc/install). 

**Go version must be 1.13+**.

2. Install truchain binaries:

```sh
git clone https://github.com/TruStory/truchain.git
cd truchain && git checkout master
make install
```

This creates:

`truchaind`: TruStory blockchain daemon

`truchaincli`: TruStory blockchain client. Used for creating keys and lightweight interaction with the chain and underlying Tendermint node.

## Getting Started

## Run a single node

```sh
# Build the binaries
make build

# Create a wallet and save the mnemonic and passphrase
make create-wallet

# Initialize configuration files and genesis file
# Enter passphrase from above
make init

# Start the chain
make start
```

## Run a local 4-node testnet

A 4-node local testnet can be created with Docker Compose.

NOTE: You will not be able to register accounts because each node won't have a registrar key setup. This restriction will go away after client-side signing.

```sh
# Build daemon for linux so it can run inside a Docker container
make build-linux

# Create 4-nodes with their own genesis files and configuration
make localnet-start
```

Go to each config file in `build/nodeN/truchaind/config/config.toml` and replace these:

```toml
laddr = "tcp://0.0.0.0:26657"
addr_book_strict = false
```

```sh
# stop and restart
make localnet-stop && make localnet-start

# Tail logs
docker-compose logs -f
```

## Run a full node

TruChain can be run as a full node, syncing it's state with another node or validator. First follow the instructions above to install and setup a single node.

```sh
# Initialize another chain with a new moniker but same chain-id
truchaind init <moniker-2> --chain-id betanet-1 --home ~/.devnet

# Copy the genesis file from the first node
scp ubuntu@devnet:/home/ubuntu/.truchaind/config/genesis.json ~/.devnet/config/

# Get the node id of the first node
truchaincli status

# Add first node to `persistent_peers` in config.toml
sed -i -e 's/persistent_peers.*/persistent_peers = "[ip_address]:26656"/' ~/.devnet/config/config.toml

# Optional: Add verbose logging
sed -i -e 's/log_level.*/log_level = "main:info,state:info,*:error,app:info,account:info,trubank:info,claim:info,community:info,truslashing:info,trustaking:info"/' ~/.devnet/config/config.toml

# Start the node
truchaind start --home ~/.devnet
```

If the first node has many blocks, it could take several minutes for the first sync to complete. Now you will have two nodes running in lockstep!

## Testing

```sh
# Run linter
make check

# Run tests
make test
```

## API Documentation

```sh
# Generate a website with documentation
make doc
```
