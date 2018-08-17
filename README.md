# TruStory Cosmos app

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
`bin/trucoind`: TruStory server node

#### Deployment

`trucli`, the light client, will ideally run on it's own machine. It will handle all
API requets, and communicate via RPC with `trucoind`.

`trucoind`, will initially run as a single Cosmos node, but eventually as a zone of many nodes.

### Messages

These are the messages needed to modify state on the TruStory blockchain.

- `PlaceBondMsg`: to place a bond on a story
- `AddCommentMsg`: to add a comment to a story
- `SubmitEvidenceMsg`: to submit evidence for a story
- `SubmitStoryMsg`: to submit a story
- `VoteMsg`: to vote on a story
