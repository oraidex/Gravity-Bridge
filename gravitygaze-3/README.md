
1. Create two Cosmos accounts - one for validator operator and another for orchestrator operator.
2. Create an Ethereum account for validator operator. Going to be used to sign messages to be executed on Gravity contract.
3. Proceed with creating validator. The only difference is that you have to specify orchestarator operator cosmos address and validator operator eth address additionally.
4. Install Rust and built orchestrator using command bellow while in `/orchestrator`. It will produce `gbt` binary in target directory.
```
cargo build --release
```
5. Initialize orchestrator config
```
gbt init
```
6. Set orchestrator operator key, so orchestrator can sign and send cosmos transactions
```
gbt keys set-orchestrator-key --phrase <MNEMONIC>
```
7. Set ethereum key, so it can be used by orchestrator
```
gbt keys set-ethereum-key --key <ETH_PRIV_KEY>
```
8. Run your Cosmos node first. Peers below
```
persistent_peers = "478286f225a7ddfe70b663ad016007f2a66dfddd@51.159.152.113:26656"
```
9.  Run orchestrator (preferred way is to create a systemd service)
```
gbt orchestrator --fees 200ugraviton --ethereum-key <ETH_VALIDATOR_PRIV_KEY> --cosmos-phrase <COSMOS_VALIDATOR_MNEMONIC> --gravityerc721-contract-address 0xDfb55f6da42484229DeC88698f032479d8C5c590 --gravity-contract-address 0xaE995E6377729CD1E987c86170159dC2395eC0c7 --ethereum-rpc "https://rpc.sepolia.org/"
```