#!/usr/bin/env bash

# go (no not like the language) to repo
GOPATH=/home/ubuntu/go
cd "${GOPATH}/src/github.com/TruStory/truchain"

# get repo
git pull || echo "git pull failed."

# build node
make update_deps
make buidl

# stop truchaind daemon
sudo systemctl stop truchaind.service

# copy binaries
cp bin/truchaind $GOPATH/bin
cp bin/trucli $GOPATH/bin

# start truchain daemon
sudo systemctl start truchaind.service
