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

## Run a single node

TruChain currently needs a _registrar_ account to sign new user registration messages.

```
# Add a new key named registrar
make registrar

# Initialize configuration files and genesis file
make init

# Edit genesis file to rename bond_denom value to "trusteak"
# Also add the registrar address to "account->params->registrar"
vi ~/.truchaind/config/genesis.json

# Collect genesis transactions
make gentx

# Start the chain
make start
```

## Run a local 4-node testnet

A 4-node local testnet can be created with Docker Compose.

NOTE: You will not be able to register accounts because each node won't have a registrar key setup. This restriction will go away after client-side signing.

```
# Build daemon for linux so it can run inside a Docker container
make build-linux

# Create 4-nodes with their own genesis files and configuration
make localnet-start

# Tail Docker logs
docker logs -f truchaindnodeN
```

## Run a full node on another machine

TruChain can be run as a full node, syncing it's state with another node. First follow the instructions above to install and setup a single node.

```
# Initialize another chain with a new moniker but same chain-id
truchaincli init <moniker-2> --chain-id betanet-1

# Copy the genesis file from the first node
scp ubuntu@devnet:/home/ubuntu/.truchaind/config/genesis.json ~/.truchaind/config/

# Add first node to `persistent_peers` in config.toml
vi ~/.truchaind/config/config.toml
# persistent_peers = "first_node_id@first_node_ip:26656"
persistent_peers = "3ebaf6ae8000af5e233ce2d3158776f7245e5ae0@ec2-54-183-49-244.us-west-1.compute.amazonaws.com:26656"

# Start the second node
truchaind start
```

If the first node has many blocks, it could take several minutes for the first sync to complete. Now you will have two nodes running in lockstep!

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

## Upgrades

To migrate between different genesis schemas, use the python script in `contrib/migration`.

```sh
python 2-stories-to-claims.py \
    --exported-genesis exported-genesis.json \
    --chain-id devnet-1 > genesis.json 
```
