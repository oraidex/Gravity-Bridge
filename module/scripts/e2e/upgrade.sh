#!/bin/bash

# setup the network using the old binary

OLD_VERSION=${OLD_VERSION:-"v1.0.3"}
ARGS="--chain-id testing -y --keyring-backend test --fees 200uoraib --gas auto --gas-adjustment 1.5 -b block"
NEW_VERSION=${NEW_VERSION:-"txidevent"}
VALIDATOR_HOME=${VALIDATOR_HOME:-"$HOME/.gravity/validator1"}

# kill all running binaries
pkill gravity && sleep 2

# download current production binary
current_dir=$PWD
rm -rf ../../gravity-bridge-old/ && git clone https://github.com/Oraichain/Gravity-Bridge.git ../../gravity-bridge-old/ && cd ../../gravity-bridge-old/module && git checkout $OLD_VERSION && go mod tidy && make install

cd $current_dir

# setup local network
sh $PWD/scripts/multi_local_node.sh

# create new upgrade proposal
UPGRADE_HEIGHT=${UPGRADE_HEIGHT:-19}
gravity tx gov submit-proposal software-upgrade $NEW_VERSION --title "foobar" --description "foobar"  --from validator1 --upgrade-height $UPGRADE_HEIGHT --upgrade-info "x" --deposit 10000000uoraib $ARGS --home $VALIDATOR_HOME
gravity tx gov vote 1 yes --from validator1 --home "$HOME/.gravity/validator1" $ARGS && gravity tx gov vote 1 yes --from validator2 --home "$HOME/.gravity/validator2" $ARGS

# sleep to wait til the proposal passes
echo "Sleep til the proposal passes..."
sleep 12

# Check if latest height is less than the upgrade height
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
sleep 7

# send test uoraib to a test account
gravity tx bank send $(gravity keys show validator1 -a --keyring-backend=test --home=$HOME/.gravity/validator1) oraib14n3tx8s5ftzhlxvq0w5962v60vd82h305kec0j 50000uoraib --keyring-backend=test --home=$HOME/.gravity/validator1 --chain-id=testing --broadcast-mode block --gas 200000 --fees 2uoraib --node http://localhost:26657 --yes

height_before=$(curl --no-progress-meter http://localhost:1317/blocks/latest | jq '.block.header.height | tonumber')

re='^[0-9]+([.][0-9]+)?$'
if ! [[ $height_before =~ $re ]] ; then
   echo "error: Not a number" >&2; exit 1
fi

sleep 30

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