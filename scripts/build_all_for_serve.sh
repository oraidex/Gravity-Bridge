#!/bin/bash
set -e

cd scripts

if [ ! -d ics721 ]; then
git clone https://github.com/public-awesome/ics721.git
fi
cd ics721
git pull
just optimize
cd ..

if [ ! -d cw-nfts ]; then
git clone https://github.com/CosmWasm/cw-nfts.git
fi
cd cw-nfts
git pull
# Make sure you have cargo make installed: $ cargo install --force cargo-make
cargo install cargo-make
cargo make optimize
cd ..

cd ..

cd module
make
cd ..

cd solidity
npm i
npm run typechain
npm run compile-deployer
cd ..

cd orchestrator
cargo build --all