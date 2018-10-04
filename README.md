# TruChain

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

`./truchaind init`

4. Start blockchain

`./truchaind start`

### Architecture

TruChain is dapp chain built with the [Cosmos SDK](https://cosmos.network/sdk) that runs on the [Cosmos Network](https://cosmos.network).

Project layout:

```
├── app
│   ├── app.go
│   └── app_test.go
├── cmd
│   ├── trucli
│   │   └── main.go
│   └── trucoind
│       └── main.go
├── types
│   └── account.go
└── x
    └── truchain
        ├── client
        │   ├── cli
        │   │   └── trustory.go
        │   └── rest
        │       └── trustory.go
        ├── db
        │   ├── back.go
        │   ├── back_queue.go
        │   ├── back_queue_test.go
        │   ├── back_test.go
        │   ├── keeper.go
        │   ├── keeper_test.go
        │   ├── story.go
        │   ├── story_test.go
        │   ├── test_common.go
        │   ├── tick.go
        │   └── tick_test.go
        ├── handler.go
        ├── handler_test.go
        └── types
            ├── back.go
            ├── back_test.go
            ├── comment.go
            ├── comment_test.go
            ├── errors.go
            ├── evidence.go
            ├── evidence_test.go
            ├── msg.go
            ├── story.go
            ├── story_test.go
            └── types.go
```

It compiles into two binaries, `trucli` (lite client) and `truchaind` (dapp chain). The lite client is responsible for responding to API requests from clients wanting to access or modify data on the dapp chain. The dapp chain is responsible for responding to requests from the lite client, such as querying and storing data.

#### Key-Value Stores

Because the current Cosmos SDK data store is built on key-value storage, database operations are more explicit than a relational or even NoSQL database. Lists and queues must be made for data that needs to be retrieved.

`TruKeeper` handles all reads and writes from key-value storage. It contains a separate store for each data type:

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

All data in stores are binary encoded using [Amino](https://github.com/tendermint/go-amino) for efficient storage in a Merkle tree. `TruKeeper` handles marshalling and umarshalling data between its binary encoding and Go data type.

### Messages

These are the messages needed to modify state on TruChain.

- `BackStoryMsg`: to back a story
- `AddCommentMsg`: to add a comment to a story
- `SubmitEvidenceMsg`: to submit evidence for a story
- `SubmitStoryMsg`: to submit a story

Each message is routed to a handler which performs operations on the key-value stores.

### Data Types

Currently the supported types are: `Back`, `Evidence`, `Comment`, and `Story`. Using Amino/protobufs allows types to be forwards and backwards compatible, allowing multiple versions of them to co-exist.

