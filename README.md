# TruStory Sidechain

### API

[TODO]

Waiting on https://github.com/cosmos/cosmos-sdk/issues/1081 to implement API.

### Running

#### Local

BUIDL the binaries for the client apps:

`make buidl`

(NOTE: `make build` also works)

This creates:

`bin/trucli`: TruStory command-line client and lite client
`bin/truchaind`: TruStory server node

#### Deployment

`trucli`, the light client, will ideally run on it's own machine. It will handle all
API requets, and communicate via RPC with `truchaind`.

`truchaind`, will initially run as a single Cosmos node, but eventually as a zone of many nodes.


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
        │   ├── keeper.go
        │   ├── keeper_test.go
        │   ├── story.go
        │   ├── story_queue.go
        │   ├── story_queue_test.go
        │   ├── story_test.go
        │   ├── tick.go
        │   ├── tick_test.go
        │   ├── vote.go
        │   └── vote_test.go
        ├── handler.go
        ├── handler_test.go
        └── types
            ├── errors.go
            ├── msg.go
            ├── msg_test.go
            └── types.go
```

It compiles into two binaries, `trucli` (lite client) and `truchaind` (dapp chain). The lite client is responsible for responding to API requests from clients wanting to access or modify data on the dapp chain. The dapp chain is responsible for responding to requests from the lite client, including querying and storing data.

### Messages

These are the messages needed to modify state on the TruStory blockchain.

- `PlaceBondMsg`: to place a bond on a story
- `AddCommentMsg`: to add a comment to a story
- `SubmitEvidenceMsg`: to submit evidence for a story
- `SubmitStoryMsg`: to submit a story
- `VoteMsg`: to vote on a story

Each of these messages mutates state in a key-value store (`KVStore`) on each node. The values in `KVStore` are binary encoded using [Amino](https://github.com/tendermint/go-amino), a library based on [Protocol Buffers](https://developers.google.com/protocol-buffers/). These values are Go types that can be serialized and deserialized.

### Types

Currently the supported types are: `Bond`, `Evidence`, `Comment`, `Story`, and `Vote`. Using protobufs allows types to be forwards and backwards compatible, allowing multiple versions of them to co-exist.

### Keepers

Keepers are wrappers around `KVStore` that handle reading and writing. They do the main heavy lifting of the app.
