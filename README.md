# TruChain

[![CircleCI](https://circleci.com/gh/TruStory/truchain.svg?style=svg&circle-token=0cea219dcac9bd6962a057d85c4a319613c6680e)](https://circleci.com/gh/TruStory/truchain)
[![codecov](https://codecov.io/gh/TruStory/truchain/branch/master/graph/badge.svg?token=jh3muAAEBs)](https://codecov.io/gh/TruStory/truchain)

TruChain is an application-specific blockchain built with [Cosmos SDK](https://cosmos.network/sdk). It powers [trustory.io](https://www.trustory.io/).

## Installation

1. Install Go by following the [official docs](https://golang.org/doc/install). 

**Go version must be 1.12+**.

2. Now let's install truchain.

```bash
mkdir -p github.com/TruStory
cd TruStory
git clone https://github.com/TruStory/truchain.git
cd truchain && git checkout master
```

## Getting Started

### Download dependencies

```
make download
```

### Install

```
make install
```

This creates:

`truchaind`: TruStory blockchain daemon

`truchaincli`: TruStory blockchain client. Used for creating keys and lightweight interaction with the chain and underlying Tendermint node.

### Setup genesis file

TruChain currently needs a _registrar_ account to sign new user registration messages.

```
# Initialize configuration files and genesis file
truchaind init trustory --chain-id localnet-1

# Configure the CLI to eliminate need for chain-id flag
truchaincli config chain-id truchain --home ~/.octopus

# Add a new key named registrar
truchaincli keys add registrar --home ~/.octopus

# Add the genesis account on the chain, giving it some trusteak
truchaind add-genesis-account $(truchaincli keys show registrar -a --home ~/.octopus) 1000000000trusteak

# Add genesis transactions that creates a test validator when the chain starts
truchaind gentx --name=registrar --amount 100000000trusteak --home-client ~/.octopus
truchaind collect-gentxs

# Start the chain
truchaind start
```

## Architecture

Each main feature of TruChain is implemented as a separate module that lives under `x/`. Each module has it's own types for data storage, _keepers_ for reading and writing this data, `Msg` types that communicate with the blockchain, and _handlers_ that route messages.

Each module has it's own [README](x/README.md).

### Key-Value Storage

Keepers handle all reads and writes from key-value storage. There's a separate keeper for each module.

All data in stores are binary encoded using [Amino](https://github.com/tendermint/go-amino) for efficient storage in a Merkle tree. Keepers handle marshalling and umarshalling data between its binary encoding and Go data type.

### Block Handlers

Most chain operations are executed based on changes to data at certain block times. After each block (`EndBlock`), a queue of story ids is checked:

#### Story Queue

A queue of all new stories that haven't expired. These are stories in the pending state. When they expire they are handled in [`x/expiration`](x/expiration/README.md), where rewards and interested are distributed.

## Local single-node development

TruChain can be setup for local development.

```
# Create a genesis file in .chain
make init

# Start the chain with logging for all modules turned on
make start
```

## Local 4-node testnet

A 4-node local testnet can be created with Docker Compose.

```
# Build daemon for linux so it can run inside a Docker container
make build-linux

# Create 4-nodes with their own genesis files and configuration
make localnet-start

# Tail chain logs within Docker Compose
docker-compose logs --tail=0 --follow
```

## Testing

```
# Run linter
make check

# Run tests
make test
```

## API Documentation

```
# Generate a website with documentation
make doc
```
