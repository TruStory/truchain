#!/usr/bin/env bash
PATH=/home/ubuntu/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/games:/usr/local/games:/snap/bin
GOPATH=/home/ubuntu/go
cd "${GOPATH}/src/github.com/TruStory/truchain"

# get repo
git pull || echo "git pull truchain failed."

# build node
make update_deps
make buidl

# stop truchaind daemon
sudo systemctl stop truchaind.service

# copy binaries
cp bin/truchaind $GOPATH/bin
cp bin/trucli $GOPATH/bin

# allow truchain to run on priviledged ports :80/:443
sudo setcap CAP_NET_BIND_SERVICE=+eip /home/ubuntu/go/bin/truchaind

# upgrade octopus
cd "/home/ubuntu/octopus"
git pull || echo "git pull octopus failed."
make db_init
make db_migrate

# start truchain daemon
sudo systemctl start truchaind.service

