package keeper

import (
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravity/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RemoveEvmChainFromStore performs in-place store migrations to remove an evm chain
func RemoveEvmChainFromStore(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec, evmChainPrefix string) error {

	ctx.Logger().Info("Pleiades Upgrade: Beginning the migrations for the gravity module")
	store := ctx.KVStore(storeKey)

	// delete evmChain data
	store.Delete(types.GetEvmChainKey(evmChainPrefix))

	// single key with chain
	removeKeyPrefixFromEvm(store, types.KeyLastOutgoingBatchID, evmChainPrefix)
	removeKeyPrefixFromEvm(store, types.LastObservedEventNonceKey, evmChainPrefix)
	removeKeyPrefixFromEvm(store, types.LastObservedEvmBlockHeightKey, evmChainPrefix)
	removeKeyPrefixFromEvm(store, types.KeyLastTXPoolID, evmChainPrefix)
	removeKeyPrefixFromEvm(store, types.LastSlashedValsetNonce, evmChainPrefix)
	removeKeyPrefixFromEvm(store, types.LatestValsetNonce, evmChainPrefix)
	removeKeyPrefixFromEvm(store, types.LastSlashedBatchBlock, evmChainPrefix)
	removeKeyPrefixFromEvm(store, types.LastSlashedLogicCallBlock, evmChainPrefix)
	removeKeyPrefixFromEvm(store, types.LastObservedValsetKey, evmChainPrefix)

	// multi key with chain
	removeKeysPrefixFromEvm(store, types.ValsetRequestKey, evmChainPrefix)
	removeKeysPrefixFromEvm(store, types.ValsetConfirmKey, evmChainPrefix)
	removeKeysPrefixFromEvm(store, types.OracleAttestationKey, evmChainPrefix)
	removeKeysPrefixFromEvm(store, types.OutgoingTXPoolKey, evmChainPrefix)
	removeKeysPrefixFromEvm(store, types.OutgoingTxBatchKey, evmChainPrefix)
	removeKeysPrefixFromEvm(store, types.BatchConfirmKey, evmChainPrefix)
	removeKeysPrefixFromEvm(store, types.LastEventNonceByValidatorKey, evmChainPrefix)
	removeKeysPrefixFromEvm(store, types.KeyOutgoingLogicCall, evmChainPrefix)
	removeKeysPrefixFromEvm(store, types.KeyOutgoingLogicConfirm, evmChainPrefix)
	removeKeysPrefixFromEvm(store, types.DenomToERC20Key, evmChainPrefix)
	removeKeysPrefixFromEvm(store, types.ERC20ToDenomKey, evmChainPrefix)
	removeKeysPrefixFromEvm(store, types.PastEvmSignatureCheckpointKey, evmChainPrefix)
	// PendingIbcAutoForwards is only existed in v3
	removeKeysPrefixFromEvm(store, types.PendingIbcAutoForwards, evmChainPrefix)

	return nil
}

func removeKeyPrefixFromEvm(store storetypes.KVStore, key []byte, evmChainPrefix string) {
	store.Delete(types.AppendChainPrefix(key, evmChainPrefix))
}

func removeKeysPrefixFromEvm(store storetypes.KVStore, key []byte, evmChainPrefix string) {
	keyPrefix := types.AppendChainPrefix(key, evmChainPrefix)
	prefixStore := prefix.NewStore(store, keyPrefix)
	storeIter := prefixStore.Iterator(nil, nil)
	defer storeIter.Close()

	for ; storeIter.Valid(); storeIter.Next() {
		prefixStore.Delete(storeIter.Key())
	}
}
