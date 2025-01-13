#!/bin/bash
set -ux

# setup the network using the old binary

OLD_VERSION=${OLD_VERSION:-"v1.0.7"}
ARGS="--chain-id testing -y --keyring-backend test --fees 200uoraib --gas auto --gas-adjustment 1.5 -b block"
NEW_VERSION=${NEW_VERSION:-"monitorevent"}
VALIDATOR_HOME=${VALIDATOR_HOME:-"$HOME/.gravity/validator1"}

# kill all running binaries
pkill gravity && sleep 2

# download current production binary
current_dir=$PWD
rm -rf ../../gravity-bridge-old/ && git clone https://github.com/Oraichain/Gravity-Bridge.git ../../gravity-bridge-old/ && cd ../../gravity-bridge-old/module && git checkout $OLD_VERSION && go mod tidy && make install

cd $current_dir

# setup local network
sh $PWD/scripts/multi_local_node.sh

# test send from cosmos to eth
#gravity tx gravity send-to-eth 0x0000000000000000000000000000000000071001 12500000000000000000oraib0x0000000000000000000000000000000000C0FFEE 1000000000000000000oraib0x0000000000000000000000000000000000C0FFEE 2500000000000000oraib0x0000000000000000000000000000000000C0FFEE oraib --fees 10uoraib  --from orchestrator1 --chain-id testing --home $HOME/.gravity/validator1 --keyring-backend test -y -b block

# test send request batch, can not set this on because client cli of this request-batch
# at v1.0.3 version is wrong format
# gravity tx gravity request-batch oraib0x0000000000000000000000000000000000C0FFEE oraib --fees 10uoraib  --from orchestrator1 --chain-id testing --home $HOME/.gravity/validator1 --keyring-backend test -y

# create new upgrade proposal, turn this on when we need real proposal
UPGRADE_HEIGHT=${UPGRADE_HEIGHT:-19}
gravity tx gov submit-proposal software-upgrade $NEW_VERSION --title "foobar" --description "foobar"  --from validator1 --upgrade-height $UPGRADE_HEIGHT --upgrade-info "x" --deposit 10000000uoraib $ARGS --home $VALIDATOR_HOME
gravity tx gov vote 1 yes --from validator1 --home "$HOME/.gravity/validator1" $ARGS && gravity tx gov vote 1 yes --from validator2 --home "$HOME/.gravity/validator2" $ARGS

# # sleep to wait til the proposal passes
echo "Sleep til the proposal passes..."
sleep 12

# # Check if latest height is less than the upgrade height
latest_height=$(curl --no-progress-meter http://localhost:1317/blocks/latest | jq '.block.header.height | tonumber')
while [ $latest_height -lt $UPGRADE_HEIGHT ];
do
   sleep 5
   ((latest_height=$(curl --no-progress-meter http://localhost:1317/blocks/latest | jq '.block.header.height | tonumber')))
   echo $latest_height
done

# kill all processes
pkill gravity

# install new binary for the upgrade
echo "install new binary"
make install

# re-run all validators. All should run
screen -S validator1 -d -m gravity start --home=$HOME/.gravity/validator1 --minimum-gas-prices=0.00001uoraib
screen -S validator2 -d -m gravity start --home=$HOME/.gravity/validator2 --minimum-gas-prices=0.00001uoraib
screen -S validator3 -d -m gravity start --home=$HOME/.gravity/validator3 --minimum-gas-prices=0.00001uoraib

# sleep a bit for the network to start 
echo "Sleep to wait for the network to start..."
sleep 10

# test send from cosmos to eth
gravity tx gravity send-to-eth 0x0000000000000000000000000000000000071001 12500000000000000000oraib0x0000000000000000000000000000000000C0FFEE 1000000000000000000oraib0x0000000000000000000000000000000000C0FFEE 2500000000000000oraib0x0000000000000000000000000000000000C0FFEE oraib --fees 10000uoraib  --from orchestrator1 --chain-id testing --home $HOME/.gravity/validator1 --keyring-backend test -y -b block --gas 400000000

# total of 3 tx ids
gravity tx gravity send-to-eth 0x0000000000000000000000000000000000071001 12500000000000000000oraib0x0000000000000000000000000000000000C0FFEE 1000000000000000000oraib0x0000000000000000000000000000000000C0FFEE 2500000000000000oraib0x0000000000000000000000000000000000C0FFEE oraib --fees 10000uoraib  --from orchestrator1 --chain-id testing --home $HOME/.gravity/validator1 --keyring-backend test -y -b block --gas 400000000

# test send request batch
txhash=$(gravity tx gravity request-batch oraib0x0000000000000000000000000000000000C0FFEE oraib --fees 10000uoraib --gas 400000000 --from orchestrator1 --chain-id testing --home $HOME/.gravity/validator1 --keyring-backend test -y | grep -o 'txhash: [^ ]*' | awk '{print $2}')

echo "Wait 10s for tx $txhash is executed on blockchain..."
sleep 10
echo "$(curl -s "http://127.0.0.1:26657/tx?hash=0x$txhash&prove=true" | jq -r '.result.tx_result.log')" | grep -q "batched_tx_ids"
if echo "$(curl -s "http://127.0.0.1:26657/tx?hash=0x$txhash&prove=true" | jq -r '.result.tx_result.log')" | grep -q "batched_tx_ids"
then
    echo "Testcase success, batched_tx_ids already exist on event"
else
    echo "Testcase failed, batched_tx_ids does not exist on event"
    exit 1
fi

height_before=$(curl --no-progress-meter http://localhost:1317/blocks/latest | jq '.block.header.height | tonumber')

re='^[0-9]+([.][0-9]+)?$'
if ! [[ $height_before =~ $re ]] ; then
   echo "error: Not a number" >&2; exit 1
fi

sleep 30

height_after=$(curl --no-progress-meter http://localhost:1317/blocks/latest | jq '.block.header.height | tonumber')

height_after=$(curl --no-progress-meter http://localhost:1317/blocks/latest | jq '.block.header.height | tonumber')

if ! [[ $height_after =~ $re ]] ; then
   echo "error: Not a number" >&2; exit 1
fi

if [ $height_after -gt $height_before ]
then
echo "Test Passed"
else
echo "Test Failed"
fi

inflation=$(curl --no-progress-meter http://localhost:1317/cosmos/mint/v1beta1/inflation | jq '.inflation | tonumber')
if ! [[ $inflation =~ $re ]] ; then
   echo "Error: Cannot query inflation => Potentially missing Go GRPC backport" >&2;
   echo "Tests Failed"; exit 1
fi
