#!/bin/bash
set -eE -o functrace

failure() {
  local lineno=$1
  local msg=$2
  echo "Failed at $lineno: $msg"
}
trap 'failure ${LINENO} "$BASH_COMMAND"' ERR

source scripts/serve_env.sh

# Stop if it is already running
if pgrep -x "$BINARY" >/dev/null; then
    echo "Terminating $BINARY..."
    pkill $BINARY
    sleep 5 # To avoid removing the folder to be any issue
fi

# Also stop Stargaze if already running
if pgrep -x "$STARGAZE_BINARY" >/dev/null; then
    echo "Terminating $STARGAZE_BINARY..."
    pkill $STARGAZE_BINARY
    sleep 5 # To avoid removing the folder to be any issue
fi

# Also stop evm if already running
lsof -ti :8545 | xargs --no-run-if-empty kill

if [ -d $CHAIN_DIR ]; then
  echo "Removing previous data from $CHAIN_DIR..."
  rm -rf $CHAIN_DIR &> /dev/null
fi

if [ -d $STARGAZE_CHAIN_DIR ]; then
  echo "Removing previous data from $STARGAZE_CHAIN_DIR..."
  rm -rf $STARGAZE_CHAIN_DIR &> /dev/null
fi

if [ -d $GBT_DIR ]; then
  echo "Removing previous data from $GBT_DIR..."
  rm -rf $GBT_DIR &> /dev/null
fi

# Add directories for gravity chain, exit if an error occurs
if ! mkdir -p $CHAIN_DIR 2>/dev/null; then
    echo "Failed to create gravity chain folder. Aborting..."
    exit 1
fi

# Add directories for stargaze chain, exit if an error occurs
if ! mkdir -p $STARGAZE_CHAIN_DIR 2>/dev/null; then
    echo "Failed to create stargaze chain folder. Aborting..."
    exit 1
fi

echo "Initializing $CHAIN_ID..."
$BINARY init test --home $CHAIN_DIR --chain-id=$CHAIN_ID

echo "Initializing $STARGAZE_CHAIN_ID..."
$STARGAZE_BINARY init test --home $STARGAZE_CHAIN_DIR --chain-id=$STARGAZE_CHAIN_ID

echo "Adding genesis accounts for gravity..."
echo "$ALICE_MNEMONIC" | $BINARY keys add alice --home $CHAIN_DIR --recover --keyring-backend=test
echo "$BOB_MNEMONIC" | $BINARY keys add bob --home $CHAIN_DIR --recover --keyring-backend=test
echo "$VALIDATOR_MNEMONIC" | $BINARY keys add validator --home $CHAIN_DIR --recover --keyring-backend=test
echo "$GRAVITY_RELAY_ACCOUNT" | $BINARY keys add gravity_relay_account --home $CHAIN_DIR --recover --keyring-backend=test
echo "$ORCHESTRATOR_MNEMONIC" | $BINARY keys add orchestrator --home $CHAIN_DIR --recover --keyring-backend=test
$BINARY add-genesis-account $($BINARY --home $CHAIN_DIR keys show alice --keyring-backend test -a) 100000000000stake  --home $CHAIN_DIR
$BINARY add-genesis-account $($BINARY --home $CHAIN_DIR keys show bob --keyring-backend test -a) 100000000000stake  --home $CHAIN_DIR
$BINARY add-genesis-account $($BINARY --home $CHAIN_DIR keys show validator --keyring-backend test -a) 100000000000stake  --home $CHAIN_DIR
$BINARY add-genesis-account $($BINARY --home $CHAIN_DIR keys show gravity_relay_account --keyring-backend test -a) 100000000000stake  --home $CHAIN_DIR
$BINARY add-genesis-account $($BINARY --home $CHAIN_DIR keys show orchestrator --keyring-backend test -a) 100000000000stake  --home $CHAIN_DIR
$BINARY gentx validator 7000000000stake $VALIDATOR_ETH_ADDRESS $ORCHESTRATOR_ADDRESS --home $CHAIN_DIR --chain-id $CHAIN_ID --keyring-backend test
$BINARY collect-gentxs --home $CHAIN_DIR

echo "Adding genesis accounts for stargaze..."
echo "$KAARE_MNEMONIC" | $STARGAZE_BINARY keys add kaare --home $STARGAZE_CHAIN_DIR --recover --keyring-backend=test
echo "$KENT_ROGER_MNEMONIC" | $STARGAZE_BINARY keys add kent --home $STARGAZE_CHAIN_DIR --recover --keyring-backend=test
echo "$STARGAZE_VALIDATOR_MNEMONIC" | $STARGAZE_BINARY keys add stargaze_validator --home $STARGAZE_CHAIN_DIR --recover --keyring-backend=test
echo "$STARGAZE_RELAY_ACCOUNT" | $STARGAZE_BINARY keys add stargaze_relay_account --home $STARGAZE_CHAIN_DIR --recover --keyring-backend=test
$STARGAZE_BINARY add-genesis-account $($STARGAZE_BINARY --home $STARGAZE_CHAIN_DIR keys show kaare --keyring-backend test -a) 100000000000stake  --home $STARGAZE_CHAIN_DIR
$STARGAZE_BINARY add-genesis-account $($STARGAZE_BINARY --home $STARGAZE_CHAIN_DIR keys show kent --keyring-backend test -a) 100000000000stake  --home $STARGAZE_CHAIN_DIR
$STARGAZE_BINARY add-genesis-account $($STARGAZE_BINARY --home $STARGAZE_CHAIN_DIR keys show stargaze_validator --keyring-backend test -a) 100000000000stake  --home $STARGAZE_CHAIN_DIR
$STARGAZE_BINARY add-genesis-account $($STARGAZE_BINARY --home $STARGAZE_CHAIN_DIR keys show stargaze_relay_account --keyring-backend test -a) 100000000000stake  --home $STARGAZE_CHAIN_DIR
$STARGAZE_BINARY gentx stargaze_validator 7000000000stake --home $STARGAZE_CHAIN_DIR --chain-id $STARGAZE_CHAIN_ID --keyring-backend test
$STARGAZE_BINARY collect-gentxs --home $STARGAZE_CHAIN_DIR

echo "Changing config (defaults and ports in app.toml and config.toml files) for gravity..."
sed -i -e 's/"nativeHRP": "osmo"/"nativeHRP": "gravity"/g' $CHAIN_DIR/config/genesis.json
sed -i -e 's/"bridge_ethereum_address": "0x0000000000000000000000000000000000000000",/"bridge_ethereum_address": "'"$GRAVITY_CONTRACT_ADDRESS"'",/g' $CHAIN_DIR/config/genesis.json
sed -i -e 's/"bridge_erc721_ethereum_address": "",/"bridge_erc721_ethereum_address": "'"$GRAVITY_ERC721_CONTRACT_ADDRESS"'",/g' $CHAIN_DIR/config/genesis.json
sed -i -e 's/"voting_period": "172800s"/"voting_period": "60s"/g' $CHAIN_DIR/config/genesis.json
sed -i -e 's#"tcp://0.0.0.0:26656"#"tcp://0.0.0.0:'"$P2P_PORT"'"#g' $CHAIN_DIR/config/config.toml
sed -i -e 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:'"$RPC_PORT"'"#g' $CHAIN_DIR/config/config.toml
sed -i -e 's/timeout_commit = "5s"/timeout_commit = "1s"/g' $CHAIN_DIR/config/config.toml
sed -i -e 's/timeout_propose = "3s"/timeout_propose = "1s"/g' $CHAIN_DIR/config/config.toml
sed -i -e 's/index_all_keys = false/index_all_keys = true/g' $CHAIN_DIR/config/config.toml
sed -i -e 's/enable = false/enable = true/g' $CHAIN_DIR/config/app.toml
sed -i -e 's/swagger = false/swagger = true/g' $CHAIN_DIR/config/app.toml
sed -i -e 's#"tcp://0.0.0.0:1317"#"tcp://0.0.0.0:'"$REST_PORT"'"#g' $CHAIN_DIR/config/app.toml
sed -i -e 's#":8080"#":'"$ROSETTA_PORT"'"#g' $CHAIN_DIR/config/app.toml
sed -i -e 's/enable-unsafe-cors = false/enable-unsafe-cors = true/g' $CHAIN_DIR/config/app.toml
sed -i -e 's/enabled-unsafe-cors = false/enable-unsafe-cors = true/g' $CHAIN_DIR/config/app.toml
sed -i.bak -e "s/^minimum-gas-prices *=.*/minimum-gas-prices = \"0.025stake\"/" $CHAIN_DIR/config/app.toml

echo "Changing config (defaults and ports in app.toml and config.toml files) for stargaze..."
sed -i -e 's/"voting_period": "172800s"/"voting_period": "60s"/g' $STARGAZE_CHAIN_DIR/config/genesis.json
sed -i -e 's#"tcp://0.0.0.0:26656"#"tcp://0.0.0.0:'"$STARGAZE_P2P_PORT"'"#g' $STARGAZE_CHAIN_DIR/config/config.toml
sed -i -e 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:'"$STARGAZE_RPC_PORT"'"#g' $STARGAZE_CHAIN_DIR/config/config.toml
sed -i -e 's/timeout_commit = "5s"/timeout_commit = "1s"/g' $STARGAZE_CHAIN_DIR/config/config.toml
sed -i -e 's/timeout_propose = "3s"/timeout_propose = "1s"/g' $STARGAZE_CHAIN_DIR/config/config.toml
sed -i -e 's/index_all_keys = false/index_all_keys = true/g' $STARGAZE_CHAIN_DIR/config/config.toml
sed -i -e 's/enable = false/enable = true/g' $STARGAZE_CHAIN_DIR/config/app.toml
sed -i -e 's/swagger = false/swagger = true/g' $STARGAZE_CHAIN_DIR/config/app.toml
sed -i -e 's#"tcp://0.0.0.0:1317"#"tcp://0.0.0.0:'"$STARGAZE_REST_PORT"'"#g' $STARGAZE_CHAIN_DIR/config/app.toml
sed -i -e 's#":8080"#":'"$STARGAZE_ROSETTA_PORT"'"#g' $STARGAZE_CHAIN_DIR/config/app.toml
sed -i -e 's/enable-unsafe-cors = false/enable-unsafe-cors = true/g' $STARGAZE_CHAIN_DIR/config/app.toml
sed -i -e 's/enabled-unsafe-cors = false/enable-unsafe-cors = true/g' $STARGAZE_CHAIN_DIR/config/app.toml
sed -i.bak -e "s/^minimum-gas-prices *=.*/minimum-gas-prices = \"0.025stake\"/" $STARGAZE_CHAIN_DIR/config/app.toml

echo "Starting $CHAIN_ID in $CHAIN_DIR..."
echo "Creating log file at $LOG_FILE_PATH"
$BINARY start --home $CHAIN_DIR --pruning=nothing --rpc.unsafe --grpc.address="0.0.0.0:$GRPC_PORT" --grpc-web.address="0.0.0.0:$GRPC_WEB" > $LOG_FILE_PATH 2>&1 &

sleep 3

if ! $BINARY --home $CHAIN_DIR --node tcp://:$RPC_PORT status; then
  echo "Gravity failed to start"
  exit 1
fi

echo ""
echo "----------- Gravity Config -------------"
echo "RPC: tcp://0.0.0.0:$RPC_PORT"
echo "REST: tcp://0.0.0.0:$REST_PORT"
echo "chain-id: $CHAIN_ID"
echo ""

echo "Starting $STARGAZE_CHAIN_ID in $STARGAZE_CHAIN_DIR..."
echo "Creating log file at $STARGAZE_LOG_FILE_PATH"
$STARGAZE_BINARY start --home $STARGAZE_CHAIN_DIR --pruning=nothing --rpc.unsafe --grpc.address="0.0.0.0:$STARGAZE_GRPC_PORT" --grpc-web.address="0.0.0.0:$STARGAZE_GRPC_WEB" > $STARGAZE_LOG_FILE_PATH 2>&1 &

sleep 3

if ! $STARGAZE_BINARY --home $STARGAZE_CHAIN_DIR --node tcp://:$STARGAZE_RPC_PORT status; then
  echo "Stargaze failed to start"
  exit 1
fi

echo ""
echo "----------- Stargaze Config -------------"
echo "RPC: tcp://0.0.0.0:$STARGAZE_RPC_PORT"
echo "REST: tcp://0.0.0.0:$STARGAZE_REST_PORT"
echo "chain-id: $STARGAZE_STARGAZE_CHAIN_ID"
echo ""

echo "-------- Chains started! --------"

echo "Starting evm chain..."
echo "Creating log file at $EVM_LOG_FILE_PATH"
cd solidity
npm run evm > $EVM_LOG_FILE_PATH 2>&1 &

sleep 3

./scripts/contract-deployer.sh
cd ..

echo "-------- EVM started and contracts deployed! --------"

echo "Setting up orchestrator..."
./orchestrator/target/debug/gbt --home $GBT_DIR init
./orchestrator/target/debug/gbt --home $GBT_DIR keys set-ethereum-key --key $VALIDATOR_ETH_PRIVATE_KEY
./orchestrator/target/debug/gbt --home $GBT_DIR keys set-orchestrator-key --phrase "$ORCHESTRATOR_MNEMONIC"
#./orchestrator/target/debug/gbt --home $GBT_DIR keys register-orchestrator-address --validator-phrase "$VALIDATOR_MNEMONIC" --fees 0stake
sed -i '/\[orchestrator\]/a check_eth_rpc = false' $GBT_DIR/config.toml

echo "Starting orchestrator..."
echo "Creating log file at $GBT_LOG_FILE_PATH"
./orchestrator/target/debug/gbt --home $GBT_DIR orchestrator --fees 0stake > $GBT_LOG_FILE_PATH 2>&1 &

sleep 5
echo "-------- Orchestrator started! --------"