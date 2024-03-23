#!/bin/bash
set -ux

# always returns true so set -e doesn't exit if it is not running.
killall gravity || true
rm -rf $HOME/.gravity/
killall screen

# make four orai directories
mkdir $HOME/.gravity
mkdir $HOME/.gravity/validator1

# init all three validators
gravity init --chain-id=testing validator1 --home=$HOME/.gravity/validator1

# create keys for all three validators
gravity keys add validator1 --keyring-backend=test --home=$HOME/.gravity/validator1
gravity keys add orchestrator1 --keyring-backend=test --home=$HOME/.gravity/validator1
gravity eth_keys add --keyring-backend=test --home=$HOME/.gravity/validator1

update_genesis () {    
    cat $HOME/.gravity/validator1/config/genesis.json | jq "$1" > $HOME/.gravity/validator1/config/tmp_genesis.json && mv $HOME/.gravity/validator1/config/tmp_genesis.json $HOME/.gravity/validator1/config/genesis.json
}

# change staking denom to uoraib
update_genesis '.app_state["staking"]["params"]["bond_denom"]="uoraib"'

# create validator node 1s
gravity add-genesis-account $(gravity keys show validator1 -a --keyring-backend=test --home=$HOME/.gravity/validator1) 1000000000000uoraib,1000000000000stake,1000000000000000000000x0000000000000000000000000000000000C0FFEE  --home=$HOME/.gravity/validator1
gravity add-genesis-account $(gravity keys show orchestrator1 -a --keyring-backend=test --home=$HOME/.gravity/validator1) 1000000000000uoraib,1000000000000stake,100000000000000000000oraib0x0000000000000000000000000000000000C0FFEE  --home=$HOME/.gravity/validator1
gravity gentx validator1 500000000uoraib 0x$(cat $HOME/.gravity/validator1/UTC--* | jq -r '.address') $(gravity keys show orchestrator1 -a --bech acc  --keyring-backend=test --home=$HOME/.gravity/validator1) --keyring-backend=test --home=$HOME/.gravity/validator1 --chain-id=testing
gravity collect-gentxs --home=$HOME/.gravity/validator1
# gravity validate-genesis --home=$HOME/.gravity/validator1

# update native hrp
update_genesis '.app_state["bech32ibc"]["nativeHRP"]="oraib"'
# update staking genesis
update_genesis '.app_state["staking"]["params"]["unbonding_time"]="240s"'
# update crisis variable to uoraib
update_genesis '.app_state["crisis"]["constant_fee"]["denom"]="uoraib"'
# udpate gov genesis
update_genesis '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="uoraib"'
# update mint genesis
update_genesis '.app_state["mint"]["params"]["mint_denom"]="uoraib"'
update_genesis '.app_state["gov"]["voting_params"]["voting_period"]="30s"'
update_genesis '.app_state["gravity"]["evm_chains"]=[{"evm_chain": {"evm_chain_prefix": "oraib","evm_chain_name": "Binance Smart Chain","evm_chain_net_version": "56"},gravity_nonces: {},valsets: [],valset_confirms: [],batches: [],batch_confirms: [],logic_calls: [],logic_call_confirms: [],attestations: [],delegate_keys: [],erc20_to_denoms: [],unbatched_transfers: []}]'
# port key (validator1 uses default ports)
# validator1 1317, 9090, 9091, 26658, 26657, 26656, 6060

# change app.toml values
VALIDATOR1_APP_TOML=$HOME/.gravity/validator1/config/app.toml

# change config.toml values
VALIDATOR1_CONFIG=$HOME/.gravity/validator1/config/config.toml

# # Pruning - comment this configuration if you want to run upgrade script
pruning="custom"
pruning_keep_recent="5"
pruning_keep_every="10"
pruning_interval="10000"

sed -i -e "s%^enable *=.*%enable = true%; " $VALIDATOR1_APP_TOML
sed -i -e "s%^pruning *=.*%pruning = \"$pruning\"%; " $VALIDATOR1_APP_TOML
sed -i -e "s%^pruning-keep-recent *=.*%pruning-keep-recent = \"$pruning_keep_recent\"%; " $VALIDATOR1_APP_TOML
sed -i -e "s%^pruning-keep-every *=.*%pruning-keep-every = \"$pruning_keep_every\"%; " $VALIDATOR1_APP_TOML
sed -i -e "s%^pruning-interval *=.*%pruning-interval = \"$pruning_interval\"%; " $VALIDATOR1_APP_TOML

# # state sync  - comment this configuration if you want to run upgrade script
snapshot_interval="10"
snapshot_keep_recent="2"

sed -i -e "s%^snapshot-interval *=.*%snapshot-interval = \"$snapshot_interval\"%; " $VALIDATOR1_APP_TOML
sed -i -e "s%^snapshot-keep-recent *=.*%snapshot-keep-recent = \"$snapshot_keep_recent\"%; " $VALIDATOR1_APP_TOML

# validator1
sed -i -E 's|allow_duplicate_ip = false|allow_duplicate_ip = true|g' $VALIDATOR1_CONFIG

# start all three validators
screen -S validator1 -d -m gravity start --home=$HOME/.gravity/validator1 --minimum-gas-prices=0.00001uoraib

echo "1 Validator is up and running!"