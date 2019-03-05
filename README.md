# TruChain

[![CircleCI](https://circleci.com/gh/TruStory/truchain.svg?style=svg&circle-token=0cea219dcac9bd6962a057d85c4a319613c6680e)](https://circleci.com/gh/TruStory/truchain)
[![codecov](https://codecov.io/gh/TruStory/truchain/branch/master/graph/badge.svg?token=jh3muAAEBs)](https://codecov.io/gh/TruStory/truchain)

## Installation

1. Install `go` by following the [official docs](https://golang.org/doc/install). Remember to set your `$GOPATH`, `$GOBIN`, and `$PATH` environment variables.


2. Now let's install truchain.

``` bash
# Go dependencies are case-sensitive.
# The directory name must have the characters T and S in capital to match the name of this repo.
mkdir -p $GOPATH/src/github.com/TruStory

cd $GOPATH/src/github.com/TruStory

git clone https://github.com/TruStory/truchain.git

cd truchain && git checkout master
```

3. Setup and configure Postgres for off-chain storage

See the db module [README](x/db/README.md).


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

`bin/trucli`: TruStory command-line client and lite client

`bin/truchaind`: TruStory blockchain daemon

`trucli`, the light client, will ideally run on it's own machine. It will communicates with `truchaind` via RPC.

`truchaind` will initially run as a single Cosmos node, but eventually as a zone of many nodes. It includes an HTTP server which handles all API requests.

### Create genesis file and setup private validator
```
make init
make reset
```

This creates `genesis.json` in the chain config folder that contains all the initial parameters for bootstrapping a node. It resets the chain so it starts with a clean slate with block 1.

Open the `genesis.json` file and in the "validators" section overwrite the "address" and "pub_key"->"value" fields with the corresponding values from `.chain/config/priv_validator.json`.

### Start chain

```
make start
```

## GraphQL Queries
You can reach your client at `http://localhost:3030/graphiql/`

Sample query:
```
query StoriesQuery {
  stories {
    body
    creator {
      address
    }
    argument
    source
    game {
      votingEndTime
    }
  }
}
```

## Architecture

TruChain is dapp chain built with the [Cosmos SDK](https://cosmos.network/sdk) that runs on the [Cosmos Network](https://cosmos.network).

Project layout:

```
├── app
│   ├── app.go
│   └── app_test.go
├── bin
│   ├── truchaind
│   └── trucli
├── cmd
│   ├── truchaind
│   │   └── main.go
│   └── trucli
│       └── main.go
├── types
│   ├── account.go
│   ├── handler.go
│   ├── keeper.go
│   └── msg.go
├── vendor
| ...
└── x
    ├── [MODULE]
    │   ├── codec.go
    │   ├── errors.go
    │   ├── handler.go
    │   ├── handler_test.go
    │   ├── keeper.go
    │   ├── keeper_queue.go
    │   ├── keeper_queue_test.go
    │   ├── keeper_test.go
    │   ├── msg.go
    │   ├── msg_test.go
    │   ├── test_common.go
    │   ├── tick.go
    │   ├── tick_test.go
    │   └── types.go
```

It compiles into two binaries, `trucli` (lite client) and `truchaind` (dapp chain). The lite client is responsible for responding to API requests from clients wanting to access or modify data on the dapp chain. The dapp chain is responsible for responding to requests from the lite client, such as querying and storing data.

Each main feature of TruChain is implemented as a separate module that lives under `x/`. Each module has it's own types for data storage, "keepers" for reading and writing this data, `Msg` types that communicate with the blockchain, and handlers that route messages.

Each module has it's own [README](x/README.md).

### On-chain Data Store

Because the current Cosmos SDK data store is built on key-value storage, database operations are more explicit than a relational or even NoSQL database. Lists and queues must be made for data that needs to be retrieved.

Keepers handle all reads and writes from key-value storage. There's a separate keeper for each module.

Each module provides a `ReadKeeper`, `WriteKeeper`, and `ReadWriteKeeper`. Other modules should get passed the appropriate keeper for it's needs. For example, if a module doesn't need to create categories, but only read them, it should get passed a category `ReadKeeper`.

All data in stores are binary encoded using [Amino](https://github.com/tendermint/go-amino) for efficient storage in a Merkle tree. Keepers handle marshalling and umarshalling data between its binary encoding and Go data type.

### Reactive Architecture

Most chain operations are executed based on changes to data. Worker queues are checked every time a new block is produced:

#### Pending Story Queue

A queue of all new stories that haven't expired or gone through voting yet. These are stories in the pending state.

#### Expiring Story Queue

Stories that expire before going into voting are pushed into this list. The queue is processed after each block to distribute rewards to backers and return funds to challengers.

#### Challenged Story Queue

Handles the lifecycle of voting on a story (validation game). Upon completion of a voting game, funds are distributed to winners, and removed from losers.

## Testing

```
make test
```

If you run the tests your `./chain/.env` will be replaced. Make sure to replace it before running the client again.

## Generate Documentation

```
make doc
```
