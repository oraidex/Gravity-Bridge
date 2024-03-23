# txhash=$(gravity tx gravity send-to-eth 0x0deB52499C2e9F3921c631cb6Ad3522C576d5484 12500000000000000000oraib0x0000000000000000000000000000000000C0FFEE 1000000000000000000oraib0x0000000000000000000000000000000000C0FFEE 2500000000000000oraib0x0000000000000000000000000000000000C0FFEE oraib --fees 10uoraib  --from orchestrator1 --chain-id testing --home $HOME/.gravity/validator1 --keyring-backend test -y | grep -oP 'txhash: \K\S+')

txhash=$(gravity tx gravity request-batch oraib0x0000000000000000000000000000000000C0FFEE oraib --fees 10uoraib  --from orchestrator1 --chain-id testing --home $HOME/.gravity/validator1 --keyring-backend test -y | grep -oP 'txhash: \K\S+')

# In ra txhash để kiểm tra
echo "txhash: $txhash"

sleep 10
if echo "$(curl -s "http://127.0.0.1:26657/tx?hash=0x$txhash&prove=true" | jq -r '.result.tx_result.log')" | grep -q "batched_tx_ids"
then
    echo "Test success, batched_tx_ids already exist on event"
    # Thêm các tác vụ khác nếu chuỗi con được tìm thấy
else
    echo "Test failed, batched_tx_ids does not exist on event"
    exit 1
fi

# gravity tx gravity request-batch oraib0x0000000000000000000000000000000000C0FFEE oraib --fees 10uoraib  --from orchestrator1 --chain-id testing --home $HOME/.gravity/validator1 --keyring-backend test

# add-evm-chain [evm-chain-name] [evm-chain-prefix] [evm-chain-net-version] [evm-chain-gravity-id] [evm-chain-bridge-eth-address] [title] [initial-deposit] [description]
# gravity tx gravity add-evm-chain bsc oraib 1 1 0x8858eeb3dfffa017d4bce9801d340d36cf895ccf "BSC MAINNET" 10000000uoraib "BSC MAINNET" --fees 2uoraib  --from validator1 --chain-id testing --home $HOME/.gravity/validator1 --keyring-backend test

# evmChain := types.EvmChainData{
#     EvmChain:           types.EvmChain{EvmChainPrefix: p.EvmChainPrefix, EvmChainName: p.EvmChainName, EvmChainNetVersion: p.EvmChainNetVersion},
#     GravityNonces:      types.GravityNonces{},
#     Valsets:            []types.Valset{},
#     ValsetConfirms:     []types.MsgValsetConfirm{},
#     Batches:            []types.OutgoingTxBatch{},
#     BatchConfirms:      []types.MsgConfirmBatch{},
#     LogicCalls:         []types.OutgoingLogicCall{},
#     LogicCallConfirms:  []types.MsgConfirmLogicCall{},
#     Attestations:       []types.Attestation{},
#     DelegateKeys:       []types.MsgSetOrchestratorAddress{},
#     Erc20ToDenoms:      []types.ERC20ToDenom{},
#     UnbatchedTransfers: []types.OutgoingTransferTx{},
# }
