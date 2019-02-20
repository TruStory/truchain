#!/usr/bin/env bash

PATH=/home/ubuntu/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/games:/usr/local/games:/snap/bin

# go (no not like the language) to repo
GOPATH=/home/ubuntu/go
cd "${GOPATH}/src/github.com/TruStory/truchain"

COMMIT_BEFORE_PULL=$(git rev-parse HEAD)

# get repo
git pull > /dev/null

COMMIT_AFTER_PULL=$(git rev-parse HEAD)

# if new commits on branch, update and restart truchain service
if [ "$COMMIT_BEFORE_PULL" != "$COMMIT_AFTER_PULL" ]
then

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

fi
