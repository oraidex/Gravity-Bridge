#!/bin/bash

# Common
ROOT_DIR=/tmp

# EVM
EVM_LOG_FILE_PATH=$ROOT_DIR/evm-local-1.log

# RELAYER
RELAYER_BINARY=rly
RELAYER_HOME_DIR=$ROOT_DIR/relayer-stuff
RELAYER_LOG_FILE_PATH=$RELAYER_HOME_DIR/relayer.log

# Gravity
BINARY=gravity
CHAIN_ID=gravity-local-1

CHAIN_DIR=$ROOT_DIR/$CHAIN_ID
LOG_FILE_PATH=$CHAIN_DIR/$CHAIN_ID.log

# TODO: FIND AND DOCUMENT THE ADDRESS OF THE ACCOUNTS
ALICE_MNEMONIC="clock post desk civil pottery foster expand merit dash seminar song memory figure uniform spice circle try happy obvious trash crime hybrid hood cushion" # gravity18hl5c9xn5dze2g50uaw0l2mr02ew57zkwe6mzg
BOB_MNEMONIC="angry twist harsh drastic left brass behave host shove marriage fall update business leg direct reward object ugly security warm tuna model broccoli choice" # gravity1qnk2n4nlkpw9xfqntladh74w6ujtulwn6z46cd
VALIDATOR_MNEMONIC="banner spread envelope side kite person disagree path silver will brother under couch edit food venture squirrel civil budget number acquire point work mass" # gravity1m9l358xunhhwds0568za49mzhvuxx9ux8fxne9
GRAVITY_RELAY_ACCOUNT="minor fetch reward clean pepper agree online oppose enroll claw mimic stable around thrive lyrics deer unknown dutch fee enhance pact horse misery electric" # gravity12xhn82e4ykp43grlgq6l52yy9e9lccypr874q6

P2P_PORT=26656
RPC_PORT=26657
REST_PORT=1317
ROSETTA_PORT=8080
GRPC_PORT=9090
GRPC_WEB=9091

GRAVITY_CONTRACT_ADDRESS=0x7580bFE88Dd3d07947908FAE12d95872a260F2D8
GRAVITY_ERC721_CONTRACT_ADDRESS=0xD50c0953a99325d01cca655E57070F1be4983b6b

# Orchestrator
GBT_DIR=$ROOT_DIR/gbt-local-1
GBT_LOG_FILE_PATH=$GBT_DIR/gbt-local-1.log

VALIDATOR_ETH_ADDRESS="0xbf660843528035a5a4921534e156a27e64b231fe"
VALIDATOR_ETH_PRIVATE_KEY="0xb1bab011e03a9862664706fc3bbaa1b16651528e5f0e7fbfcbfdd8be302a13e7"
ORCHESTRATOR_MNEMONIC="file lamp sunny powder judge fatigue wool target kit mimic neck debris liar miracle crime weapon lucky tongue gorilla goose rib unique flee satisfy" # gravity1ay00mn99whhjmzteuwdlv2negp3kgwqaqk5869
ORCHESTRATOR_ADDRESS=gravity1ay00mn99whhjmzteuwdlv2negp3kgwqaqk5869

# STARGAZE
STARGAZE_BINARY=starsd
STARGAZE_CHAIN_ID=stargaze-local-1

STARGAZE_CHAIN_DIR=$ROOT_DIR/$STARGAZE_CHAIN_ID
STARGAZE_LOG_FILE_PATH=$CHAIN_DIR/$CHAIN_ID.log

# TODO: FIND AND DOCUMENT THE ADDRESS OF THE ACCOUNTS
KAARE_MNEMONIC="skate height year measure reunion toast onion canal cupboard innocent dash develop spend pottery wine nest orchard vibrant zebra climb cricket carbon unhappy color" # stars1j0hkmu8rklcewz4g0wclxlzf4tzhlx00a9apjl
KENT_ROGER_MNEMONIC="name rose pill armor surprise position vague suggest rescue april evidence all silly original dignity wet seven lazy slogan smoke genre cost faith royal" # stars1yza3nu6qz8kwn67tgd395s8yedpq45vf5pfczp
STARGAZE_VALIDATOR_MNEMONIC="traffic occur lens age swim tilt canvas train stairs leg base inmate vessel trigger abstract thunder whale resist summer popular nature move original tired" # stars16vrj6an5f7pl0g8gl8qxex7rts2vucq4gw0m7c
<<<<<<< HEAD
STARGAZE_RELAY_ACCOUNT="result impulse book sand wedding mass top ritual swing assault claw mind outside hand kind detect gasp large radar relief wool tank taxi item" # stars16vrj6an5f7pl0g8gl8qxex7rts2vucq4gw0m7c
=======
STARGAZE_RELAY_ACCOUNT="result impulse book sand wedding mass top ritual swing assault claw mind outside hand kind detect gasp large radar relief wool tank taxi item" # stars1sl3mpymvv70z5xyxfe7fvlqqjvm4v4hlrmuq2e
>>>>>>> 1afca4e1bb393f95473ad06d99899f3582308119

STARGAZE_P2P_PORT=36656
STARGAZE_RPC_PORT=36657
STARGAZE_REST_PORT=2317
STARGAZE_ROSETTA_PORT=9080
STARGAZE_GRPC_PORT=9991
STARGAZE_GRPC_WEB=9992