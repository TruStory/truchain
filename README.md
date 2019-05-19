# TruChain

[![CircleCI](https://circleci.com/gh/TruStory/truchain.svg?style=svg&circle-token=0cea219dcac9bd6962a057d85c4a319613c6680e)](https://circleci.com/gh/TruStory/truchain)
[![codecov](https://codecov.io/gh/TruStory/truchain/branch/master/graph/badge.svg?token=jh3muAAEBs)](https://codecov.io/gh/TruStory/truchain)

## Installation

1. Install Go by following the [official docs](https://golang.org/doc/install). Remember to set your `$GOPATH`, `$GOBIN`, and `$PATH` environment variables.

**Go version must be 1.12+**.

2. Now let's install truchain.

``` bash
mkdir -p $GOPATH/src/github.com/TruStory
cd $GOPATH/src/github.com/TruStory
git clone https://github.com/TruStory/truchain.git
cd truchain && git checkout master
```

## Running

### Install dependencies

```
make update_deps
```

### Build binaries

```
make buidl
```

This creates:

`./bin/truchaind`: TruStory blockchain daemon

`./bin/truchaincli`: TruStory blockchain client. Used for creating keys and lightweight interaction with the chain and underlying Tendermint node.


### Create genesis file and setup private validator

```
make init
```

This creates `genesis.json` in the chain config folder that contains all the initial parameters for bootstrapping a node. It resets the chain so it starts with a clean slate with block 0.

### Start chain

```
make start
```

## Architecture

TruChain is an application-specific blockchain built with [Cosmos SDK](https://cosmos.network/sdk).

Each main feature of TruChain is implemented as a separate module that lives under `x/`. Each module has it's own types for data storage, _keepers_ for reading and writing this data, `Msg` types that communicate with the blockchain, and _handlers_ that route messages.

Each module has it's own [README](x/README.md).

### Key-Value Storage

Because the current Cosmos SDK data store is built on key-value storage, database operations are more explicit than a relational or even NoSQL database. Lists and queues must be made for data that needs to be retrieved.

Keepers handle all reads and writes from key-value storage. There's a separate keeper for each module.

All data in stores are binary encoded using [Amino](https://github.com/tendermint/go-amino) for efficient storage in a Merkle tree. Keepers handle marshalling and umarshalling data between its binary encoding and Go data type.

### Block Handlers

Most chain operations are executed based on changes to data at certain block times. After each block (`EndBlock`), a queue of story ids is checked:

#### Story Queue

A queue of all new stories that haven't expired. These are stories in the pending state. When they expire they are handled in [`x/expiration`](x/expiration/README.md), where rewards and interested are distributed.

## Testing

```
make test
```

## Generate Documentation

```
make doc
```
