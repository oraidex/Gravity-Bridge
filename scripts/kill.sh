#!/bin/bash
set -eE -o functrace

failure() {
  local lineno=$1
  local msg=$2
  echo "Failed at $lineno: $msg"
}
trap 'failure ${LINENO} "$BASH_COMMAND"' ERR

source scripts/serve_env.sh

echo "Killing all servers and removing data..."

# Stop if it is already running
if pgrep -x "$BINARY" >/dev/null; then
    echo "Terminating $BINARY..."
    pkill $BINARY
    sleep 2 # To avoid removing the folder to be any issue
fi

# Also stop Stargaze if already running
if pgrep -x "$STARGAZE_BINARY" >/dev/null; then
    echo "Terminating $STARGAZE_BINARY..."
    pkill $STARGAZE_BINARY
    sleep 2 # To avoid removing the folder to be any issue
fi

# TODO: Also stop gbt if already running

# Also stop rly if already running
if pgrep -x "$RELAYER_BINARY" >/dev/null; then
    echo "Terminating $RELAYER_BINARY..."
    pkill $RELAYER_BINARY
    sleep 2 # To avoid removing the folder to be any issue
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

if [ -d $RELAYER_HOME_DIR ]; then
  echo "Removing previous data from $RELAYER_HOME_DIR..."
  rm -rf $RELAYER_HOME_DIR &> /dev/null
fi