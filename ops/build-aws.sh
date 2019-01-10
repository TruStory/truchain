#!/usr/bin/env bash

PATH=/home/ubuntu/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/games:/usr/local/games:/snap/bin

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

# copy files
cp bin/truchaind $GOPATH/bin
cp bin/trucli $GOPATH/bin
cp .chain/bootstrap.csv ~/.truchaind/bootstrap.csv

# start truchain daemon
sudo systemctl start truchaind.service
