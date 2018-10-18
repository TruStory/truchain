# TruChain

### Installation

1. Install `go` by following the [official docs](https://golang.org/doc/install). Remember to set your `$GOPATH`, `$GOBIN`, and `$PATH` environment variables.


2. Now let's install truchain.

``` bash
mkdir -p $GOPATH/src/github.com/trustory`

cd $GOPATH/src/github.com/trustory`

git clone https://github.com/TruStory/truchain.git`

cd truchain && git checkout master`
```

### Running

1. Install dependencies

`make update_vendor_deps`

This should install all dependencies in `/vendor`.

2. Buidl the binaries for the client apps:

`make buidl`

(NOTE: `make build` also works)

This creates:

`bin/trucli`: TruStory command-line client and lite client

`bin/truchaind`: TruStory server node

`trucli`, the light client, will ideally run on it's own machine. It will handle all
API requests, and communicate via RPC with `truchaind`.

`truchaind`, will initially run as a single Cosmos node, but eventually as a zone of many nodes.

3. Create genesis file (one-time only)

    a. `bin/truchaind init`

    b. Edit `~/.truchaind/config/genesis.json` and add the following account to the `"accounts"` array:

    ```json
    {
      "address": "cosmos1w3e82cmgv95kuctrvdex2emfwd68yctjpzp3mr",
      "coins": [{"denom": "trustake", "amount": "123456"}]
    }
    ```

4. Start blockchain

`bin/truchaind start`

### Architecture

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
    ├── backing
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
    ├── category
    │   ├── codec.go
    │   ├── errors.go
    │   ├── handler.go
    │   ├── handler_test.go
    │   ├── keeper.go
    │   ├── keeper_test.go
    │   ├── msg.go
    │   ├── msg_test.go
    │   ├── test_common.go
    │   └── types.go
    └── story
        ├── codec.go
        ├── errors.go
        ├── handler.go
        ├── handler_test.go
        ├── keeper.go
        ├── keeper_test.go
        ├── msg.go
        ├── msg_test.go
        ├── test_common.go
        └── types.go
```

It compiles into two binaries, `trucli` (lite client) and `truchaind` (dapp chain). The lite client is responsible for responding to API requests from clients wanting to access or modify data on the dapp chain. The dapp chain is responsible for responding to requests from the lite client, such as querying and storing data.

Each main feature of TruChain is implemented as a separate module that lives under `x/`. Each module has it's own types for data storage, "keepers" for reading and writing this data, `Msg` types that communicate with the blockchain, and handlers that route messages.

#### Key-Value Stores

Because the current Cosmos SDK data store is built on key-value storage, database operations are more explicit than a relational or even NoSQL database. Lists and queues must be made for data that needs to be retrieved.

Keepers handle all reads and writes from key-value storage. There's a separate keeper for each data type:

Stories key-value store:

* stories:1 -> `Story`
* stories:2 -> `Story`
* stories:len -> 2

Backings key-value store:

* backings:1 -> `Backing`
* backings:2 -> `Backing`
* backings:3 -> `Backing`
* backings:len -> 3
* backings:queue:unexpired -> [1, 2, 3]

Each module provides a `ReadKeeper`, `WriteKeeper`, and `ReadWriteKeeper`. Other modules should get passed the appropriate keeper for it's needs. For example, if a module doesn't need to create categories, but only read them, it should get passed a category `ReadKeeper`.

All data in stores are binary encoded using [Amino](https://github.com/tendermint/go-amino) for efficient storage in a Merkle tree. Keepers handle marshalling and umarshalling data between its binary encoding and Go data type.

### Messages

These are the messages needed to modify state on TruChain.

- `CreateCategoryMsg` creates a category for a story
- `BackStoryMsg`: to back a story
- `SubmitEvidenceMsg`: to submit evidence for a story
- `SubmitStoryMsg`: to submit a story

Each message is routed to a handler which performs operations on the key-value stores.

### Data Types

Currently the supported types are: `Back`, `Evidence`, `Category`, and `Story`. Using Amino/protobufs allows types to be forwards and backwards compatible, allowing multiple versions of them to co-exist.

