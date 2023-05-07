#!/bin/bash

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