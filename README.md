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

`make update_deps`

2. Buidl the binaries for the client apps:

`make buidl`

NOTE: On macOS Mojave, you might have to run `export CGO_ENABLED=1; export CC=gcc`.

This creates:

`bin/trucli`: TruStory command-line client and lite client

`bin/truchaind`: TruStory server node

`trucli`, the light client, will ideally run on it's own machine. It will communicates with `truchaind` via RPC.

`truchaind` will initially run as a single Cosmos node, but eventually as a zone of many nodes. It includes an HTTP server which handles all API requests.

3. Create registrar key

Create a file in the root of this project named `registrar.key` containing the upper-case-hex encoding of a secp256k1 private key. Example:

```
6D5A20923CB334E4950C32C344842FE0DCBAC559FB04AE97D64B49ACC81BC1FB
```

This is the private key to the account added in step 3b.

4. Copy genesis file and setup private validator

Copy `genesis.json` from secrets repo into `.chain/config`.

Run `make wipe_chain`. 

Open the `genesis.json` file and in the "validators" section overwrite the "address" and "pub_key"->"value" fields with the corresponding values from `.chain/config/priv_validator.json`.

5. Start blockchain

`make run_daemon`

You can wipe the chain, build, and start using the alias `make bwr`.

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

#### Key-Value Stores

Because the current Cosmos SDK data store is built on key-value storage, database operations are more explicit than a relational or even NoSQL database. Lists and queues must be made for data that needs to be retrieved.

Keepers handle all reads and writes from key-value storage. There's a separate keeper for each module.

Each module provides a `ReadKeeper`, `WriteKeeper`, and `ReadWriteKeeper`. Other modules should get passed the appropriate keeper for it's needs. For example, if a module doesn't need to create categories, but only read them, it should get passed a category `ReadKeeper`.

All data in stores are binary encoded using [Amino](https://github.com/tendermint/go-amino) for efficient storage in a Merkle tree. Keepers handle marshalling and umarshalling data between its binary encoding and Go data type.

### Recipes

For step-by-step instructions for common tasks like adding a new message type or query, see [./dev/RECIPES.md](./dev/RECIPES.md).

### Testing

`make test`
