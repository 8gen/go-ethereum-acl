#!/bin/bash

set -e
rm -rf datadir/node*/geth -v
make geth
./build/bin/geth init --datadir datadir/node1 datadir/node1/genesis.json
./build/bin/geth init --datadir datadir/node2 datadir/node1/genesis.json
./build/bin/geth init --datadir datadir/node3 datadir/node1/genesis.json
./build/bin/geth init --datadir datadir/node4 datadir/node1/genesis.json
foreman start
killall -9 geth
