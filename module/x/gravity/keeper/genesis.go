package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravity/types"
)

func initBridgeDataFromGenesis(ctx sdk.Context, k Keeper, data types.EvmChainData) {
	// reset valsets in state
	evmChainPrefix := data.EvmChain.EvmChainPrefix
	highest := uint64(0)
	for _, vs := range data.Valsets {
		if vs.Nonce > highest {
			highest = vs.Nonce
		}
		k.StoreValset(ctx, evmChainPrefix, vs)
	}
	k.SetLatestValsetNonce(ctx, evmChainPrefix, highest)

	// reset valset confirmations in state
	for _, conf := range data.ValsetConfirms {
		k.SetValsetConfirm(ctx, conf)
	}

	// reset batches in state
	for _, batch := range data.Batches {
		// TODO: block height?
		intBatch, err := batch.ToInternal()
		if err != nil {
			panic(sdkerrors.Wrapf(err, "unable to make batch internal: %v", batch))
		}
		k.StoreBatch(ctx, evmChainPrefix, *intBatch)
	}

	// reset batch confirmations in state
	for _, conf := range data.BatchConfirms {
		conf := conf
		k.SetBatchConfirm(ctx, &conf)
	}

	// reset logic calls in state
	for _, call := range data.LogicCalls {
		k.SetOutgoingLogicCall(ctx, evmChainPrefix, call)
	}

	// reset logic call confirmations in state
	for _, conf := range data.LogicCallConfirms {
		conf := conf
		k.SetLogicCallConfirm(ctx, &conf)
	}
<<<<<<< HEAD
=======
}

// InitGenesis starts a chain from a genesis state
func InitGenesis(ctx sdk.Context, k Keeper, data types.GenesisState) {
	k.SetParams(ctx, *data.Params)

	// restore various nonces, this MUST match GravityNonces in genesis
	k.SetLatestValsetNonce(ctx, data.GravityNonces.LatestValsetNonce)
	k.setLastObservedEventNonce(ctx, data.GravityNonces.LastObservedNonce, types.GravityContractNonce)
	k.setLastObservedEventNonce(ctx, data.GravityNonces.LastErc721ObservedNonce, types.ERC721ContractNonce)
	k.SetLastSlashedValsetNonce(ctx, data.GravityNonces.LastSlashedValsetNonce)
	k.SetLastSlashedBatchBlock(ctx, data.GravityNonces.LastSlashedBatchBlock)
	k.SetLastSlashedLogicCallBlock(ctx, data.GravityNonces.LastSlashedLogicCallBlock)
	k.setID(ctx, data.GravityNonces.LastTxPoolId, []byte(types.KeyLastTXPoolID))
	k.setID(ctx, data.GravityNonces.LastBatchId, []byte(types.KeyLastOutgoingBatchID))

	initBridgeDataFromGenesis(ctx, k, data)
>>>>>>> 81057dc97ff3a6f3702fca99300ddbb3a7011770

	// reset pool transactions in state
	for _, tx := range data.UnbatchedTransfers {
		intTx, err := tx.ToInternal()
		if err != nil {
			panic(sdkerrors.Wrapf(err, "invalid unbatched tx: %v", tx))
		}
		if err := k.addUnbatchedTX(ctx, evmChainPrefix, intTx); err != nil {
			panic(err)
		}
	}

	// reset attestations in state
	for _, att := range data.Attestations {
		att := att
		claim, err := k.UnpackAttestationClaim(&att)
		if err != nil {
			panic("couldn't cast to claim")
		}

		// block height?
		if claim.GetEthBlockHeight() == 0 {
			panic("eth block height can not be zero")
		}
		hash, err := claim.ClaimHash()
		if err != nil {
			panic(fmt.Errorf("error when computing ClaimHash for %v", hash))
		}
<<<<<<< HEAD
		// these claims always have the same evmChainPrefix even not set
		claim.SetEvmChainPrefix(evmChainPrefix)

		k.SetAttestation(ctx, claim.GetEvmChainPrefix(), claim.GetEventNonce(), hash, &att)
=======
		k.SetAttestation(ctx, claim.GetEventNonce(), hash, &att, types.GravityContractNonce)
	}
	// reset erc721 attestations in state
	for _, att := range data.Erc721Attestations {
		att := att
		claim, err := k.UnpackAttestationClaim(&att)
		if err != nil {
			panic("couldn't cast to claim")
		}

		// TODO: block height?
		hash, err := claim.ClaimHash()
		if err != nil {
			panic(fmt.Errorf("error when computing ClaimHash for %v", hash))
		}
		k.SetAttestation(ctx, claim.GetEventNonce(), hash, &att, types.ERC721ContractNonce)
>>>>>>> 81057dc97ff3a6f3702fca99300ddbb3a7011770
	}

	// reset attestation state of specific validators
	// this must be done after the above to be correct
	for _, att := range data.Attestations {
		att := att
		claim, err := k.UnpackAttestationClaim(&att)
		if err != nil {
			panic("couldn't cast to claim")
		}
		// reconstruct the latest event nonce for every validator
		// if somehow this genesis state is saved when all attestations
		// have been cleaned up GetLastEventNonceByValidator handles that case
		//
		// if we were to save and load the last event nonce for every validator
		// then we would need to carry that state forever across all chain restarts
		// but since we've already had to handle the edge case of new validators joining
		// while all attestations have already been cleaned up we can do this instead and
		// not carry around every validators event nonce counter forever.
		for _, vote := range att.Votes {
			val, err := sdk.ValAddressFromBech32(vote)
			if err != nil {
				panic(err)
			}
<<<<<<< HEAD
			last := k.GetLastEventNonceByValidator(ctx, evmChainPrefix, val)
			if claim.GetEventNonce() > last {
				k.SetLastEventNonceByValidator(ctx, evmChainPrefix, val, claim.GetEventNonce())
=======
			last := k.GetLastEventNonceByValidator(ctx, val, types.GravityContractNonce)
			if claim.GetEventNonce() > last {
				k.SetLastEventNonceByValidator(ctx, val, claim.GetEventNonce(), types.GravityContractNonce)
			}
		}
	}

	// reset erc721 attestation state of specific validators
	// this must be done after the above to be correct
	for _, att := range data.Erc721Attestations {
		att := att
		claim, err := k.UnpackAttestationClaim(&att)
		if err != nil {
			panic("couldn't cast to claim")
		}
		/*
			reconstruct the latest event nonce for every validator
			if somehow this genesis state is saved when all attestations
			have been cleaned up GetLastEventNonceByValidator handles that case

			if we were to save and load the last event nonce for every validator
			then we would need to carry that state forever across all chain restarts
			but since we've already had to handle the edge case of new validators joining
			while all attestations have already been cleaned up we can do this instead and
			not carry around every validator's event nonce counter forever.
		*/
		for _, vote := range att.Votes {
			val, err := sdk.ValAddressFromBech32(vote)
			if err != nil {
				panic(err)
			}
			last := k.GetLastEventNonceByValidator(ctx, val, types.ERC721ContractNonce)
			if claim.GetEventNonce() > last {
				k.SetLastEventNonceByValidator(ctx, val, claim.GetEventNonce(), types.ERC721ContractNonce)
>>>>>>> 81057dc97ff3a6f3702fca99300ddbb3a7011770
			}
		}
	}

	// reset delegate keys in state
	if hasDuplicates(data.DelegateKeys) {
		panic("Duplicate delegate key found in Genesis!")
	}
	for _, keys := range data.DelegateKeys {
		err := keys.ValidateBasic()
		if err != nil {
			panic("Invalid delegate key in Genesis!")
		}
		val, err := sdk.ValAddressFromBech32(keys.Validator)
		if err != nil {
			panic(err)
		}
		evmAddr, err := types.NewEthAddress(keys.EthAddress)
		if err != nil {
			panic(err)
		}

		orch, err := sdk.AccAddressFromBech32(keys.Orchestrator)
		if err != nil {
			panic(err)
		}

		// set the orchestrator address
		k.SetOrchestratorValidator(ctx, val, orch)
		// set the ethereum address
		k.SetEvmAddressForValidator(ctx, val, *evmAddr)
	}

	// populate state with cosmos originated denom-erc20 mapping
	for i, item := range data.Erc20ToDenoms {
		ethAddr, err := types.NewEthAddress(item.Erc20)
		if err != nil {
			panic(fmt.Errorf("invalid erc20 address in Erc20ToDenoms for item %d: %s", i, item.Erc20))
		}
		k.setCosmosOriginatedDenomToERC20(ctx, evmChainPrefix, item.Denom, *ethAddr)
	}

	// now that we have the denom-erc20 mapping we need to validate
	// that the valset reward is possible and cosmos originated remove
	// this if you want a non-cosmos originated reward
	valsetReward := k.GetParams(ctx).ValsetReward
	if valsetReward.IsValid() && !valsetReward.IsZero() {
		_, exists := k.GetCosmosOriginatedERC20(ctx, evmChainPrefix, valsetReward.Denom)
		if !exists {
			panic("Invalid Cosmos originated denom for valset reward")
		}
	}
}

// InitGenesis starts a chain from a genesis state
func InitGenesis(ctx sdk.Context, k Keeper, data types.GenesisState) {
	k.SetParams(ctx, *data.Params)

	for _, evmChain := range data.EvmChains {
		// restore various nonces, this MUST match GravityNonces in genesis
		chainPrefix := evmChain.EvmChain.EvmChainPrefix
		k.SetLatestValsetNonce(ctx, chainPrefix, evmChain.GravityNonces.LatestValsetNonce)
		k.setLastObservedEventNonce(ctx, chainPrefix, evmChain.GravityNonces.LastObservedNonce)
		k.SetLastSlashedValsetNonce(ctx, chainPrefix, evmChain.GravityNonces.LastSlashedValsetNonce)
		k.SetLastSlashedBatchBlock(ctx, chainPrefix, evmChain.GravityNonces.LastSlashedBatchBlock)
		k.SetLastSlashedLogicCallBlock(ctx, chainPrefix, evmChain.GravityNonces.LastSlashedLogicCallBlock)
		k.SetLastObservedEvmChainBlockHeight(ctx, chainPrefix, evmChain.GravityNonces.LastObservedEvmBlockHeight)
		k.setID(ctx, evmChain.GravityNonces.LastTxPoolId, types.AppendChainPrefix(types.KeyLastTXPoolID, chainPrefix))
		k.setID(ctx, evmChain.GravityNonces.LastBatchId, types.AppendChainPrefix(types.KeyLastOutgoingBatchID, chainPrefix))
		k.SetEvmChainData(ctx, evmChain.EvmChain)

		initBridgeDataFromGenesis(ctx, k, evmChain)
	}

	for _, forward := range data.PendingErc721IbcAutoForwards {
		err := k.addPendingERC721PendingIbcAutoForward(ctx, forward)
		if err != nil {
			panic(fmt.Errorf("unable to restore pending ibc auto forward (%v) to store: %v", forward, err))
		}
	}
}

func hasDuplicates(d []types.MsgSetOrchestratorAddress) bool {
	ethMap := make(map[string]struct{}, len(d))
	orchMap := make(map[string]struct{}, len(d))
	// creates a hashmap then ensures that the hashmap and the array
	// have the same length, this acts as an O(n) duplicates check
	for i := range d {
		ethMap[d[i].EthAddress] = struct{}{}
		orchMap[d[i].Orchestrator] = struct{}{}
	}
	return len(ethMap) != len(d) || len(orchMap) != len(d)
}

// ExportGenesis exports all the state needed to restart the chain
// from the current state of the chain
func ExportGenesis(ctx sdk.Context, k Keeper) types.GenesisState {
<<<<<<< HEAD
	p := k.GetParams(ctx)

	chains := k.GetEvmChains(ctx)
	evmChains := make([]types.EvmChainData, len(chains))
=======
	var (
		p                           = k.GetParams(ctx)
		calls                       = k.GetOutgoingLogicCalls(ctx)
		batches                     = k.GetOutgoingTxBatches(ctx)
		valsets                     = k.GetValsets(ctx)
		attmap, attKeys             = k.GetAttestationMapping(ctx, types.GravityContractNonce)
		erc721AttMap, erc721AttKeys = k.GetAttestationMapping(ctx, types.ERC721ContractNonce)
		vsconfs                     = []types.MsgValsetConfirm{}
		batchconfs                  = []types.MsgConfirmBatch{}
		callconfs                   = []types.MsgConfirmLogicCall{}
		attestations                = []types.Attestation{}
		erc721Attestations                = []types.Attestation{}
		delegates                   = k.GetDelegateKeys(ctx)
		erc20ToDenoms               = []types.ERC20ToDenom{}
		unbatchedTransfers          = k.GetUnbatchedTransactions(ctx)
		pendingForwards             = k.PendingIbcAutoForwards(ctx, 0)
		erc721PendingForwards       = k.PendingERC721IbcAutoForwards(ctx, 0)
	)
	var forwards []types.PendingIbcAutoForward
	for _, forward := range pendingForwards {
		forwards = append(forwards, *forward)
	}

	var erc721Forwards []types.PendingERC721IbcAutoForward
	for _, forward := range erc721PendingForwards {
		erc721Forwards = append(erc721Forwards, *forward)
	}

	// export valset confirmations from state
	for _, vs := range valsets {
		// TODO: set height = 0?
		vsconfs = append(vsconfs, k.GetValsetConfirms(ctx, vs.Nonce)...)
	}
>>>>>>> 81057dc97ff3a6f3702fca99300ddbb3a7011770

	for ci, evmChain := range chains {
		calls := k.GetOutgoingLogicCalls(ctx, evmChain.EvmChainPrefix)
		batches := k.GetOutgoingTxBatches(ctx, evmChain.EvmChainPrefix)
		valsets := k.GetValsets(ctx, evmChain.EvmChainPrefix)
		attmap, attKeys := k.GetAttestationMapping(ctx, evmChain.EvmChainPrefix)
		vsconfs := []types.MsgValsetConfirm{}
		batchconfs := []types.MsgConfirmBatch{}
		callconfs := []types.MsgConfirmLogicCall{}
		attestations := []types.Attestation{}
		delegates := k.GetDelegateKeys(ctx)
		erc20ToDenoms := []types.ERC20ToDenom{}
		unbatchedTransfers := k.GetUnbatchedTransactions(ctx, evmChain.EvmChainPrefix)

		// export valset confirmations from state
		for _, vs := range valsets {
			// TODO: set height = 0?
			vsconfs = append(vsconfs, k.GetValsetConfirms(ctx, evmChain.EvmChainPrefix, vs.Nonce)...)
		}

		// export batch confirmations from state
		extBatches := make([]types.OutgoingTxBatch, len(batches))
		for i, batch := range batches {
			// TODO: set height = 0?
			batchconfs = append(batchconfs,
				k.GetBatchConfirmByNonceAndTokenContract(ctx, evmChain.EvmChainPrefix, batch.BatchNonce, batch.TokenContract)...)
			extBatches[i] = batch.ToExternal()
		}

<<<<<<< HEAD
		// export logic call confirmations from state
		for _, call := range calls {
			// TODO: set height = 0?
			callconfs = append(callconfs,
				k.GetLogicConfirmsByInvalidationIDAndNonce(ctx, evmChain.EvmChainPrefix, call.InvalidationId, call.InvalidationNonce)...)
		}
=======
	// export erc721 attestations from state
	for _, key := range erc721AttKeys {
		// TODO: set height = 0?
		erc721Attestations = append(erc721Attestations, erc721AttMap[key]...)
	}

	// export erc20 to denom relations
	k.IterateERC20ToDenom(ctx, func(key []byte, erc20ToDenom *types.ERC20ToDenom) bool {
		erc20ToDenoms = append(erc20ToDenoms, *erc20ToDenom)
		return false
	})
>>>>>>> 81057dc97ff3a6f3702fca99300ddbb3a7011770

		// export attestations from state
		for _, key := range attKeys {
			// TODO: set height = 0?
			attestations = append(attestations, attmap[key]...)
		}

		// export erc20 to denom relations
		k.IterateERC20ToDenom(ctx, evmChain.EvmChainPrefix, func(key []byte, erc20ToDenom *types.ERC20ToDenom) bool {
			erc20ToDenoms = append(erc20ToDenoms, *erc20ToDenom)
			return false
		})

		unbatchedTxs := make([]types.OutgoingTransferTx, len(unbatchedTransfers))
		for i, v := range unbatchedTransfers {
			unbatchedTxs[i] = v.ToExternal()
		}

		evmChains[ci] = types.EvmChainData{
			EvmChain: types.EvmChain{
				EvmChainNetVersion: evmChain.EvmChainNetVersion,
				EvmChainPrefix:     evmChain.EvmChainPrefix,
				EvmChainName:       evmChain.EvmChainName,
			},
			GravityNonces: types.GravityNonces{
				LatestValsetNonce:          k.GetLatestValsetNonce(ctx, evmChain.EvmChainPrefix),
				LastObservedNonce:          k.GetLastObservedEventNonce(ctx, evmChain.EvmChainPrefix),
				LastSlashedValsetNonce:     k.GetLastSlashedValsetNonce(ctx, evmChain.EvmChainPrefix),
				LastSlashedBatchBlock:      k.GetLastSlashedBatchBlock(ctx, evmChain.EvmChainPrefix),
				LastSlashedLogicCallBlock:  k.GetLastSlashedLogicCallBlock(ctx, evmChain.EvmChainPrefix),
				LastObservedEvmBlockHeight: k.GetLastObservedEvmChainBlockHeight(ctx, evmChain.EvmChainPrefix).EthereumBlockHeight,
				LastTxPoolId:               k.getID(ctx, types.AppendChainPrefix(types.KeyLastTXPoolID, evmChain.EvmChainPrefix)),
				LastBatchId:                k.getID(ctx, types.AppendChainPrefix(types.KeyLastOutgoingBatchID, evmChain.EvmChainPrefix)),
			},
			Valsets:            valsets,
			ValsetConfirms:     vsconfs,
			Batches:            extBatches,
			BatchConfirms:      batchconfs,
			LogicCalls:         calls,
			LogicCallConfirms:  callconfs,
			Attestations:       attestations,
			DelegateKeys:       delegates,
			Erc20ToDenoms:      erc20ToDenoms,
			UnbatchedTransfers: unbatchedTxs,
		}
	}

	return types.GenesisState{
<<<<<<< HEAD
		Params:    &p,
		EvmChains: evmChains,
=======
		Params: &p,
		GravityNonces: types.GravityNonces{
			LatestValsetNonce:         k.GetLatestValsetNonce(ctx),
			LastObservedNonce:         k.GetLastObservedEventNonce(ctx, types.GravityContractNonce),
			LastErc721ObservedNonce:   k.GetLastObservedEventNonce(ctx, types.ERC721ContractNonce),
			LastSlashedValsetNonce:    k.GetLastSlashedValsetNonce(ctx),
			LastSlashedBatchBlock:     k.GetLastSlashedBatchBlock(ctx),
			LastSlashedLogicCallBlock: k.GetLastSlashedLogicCallBlock(ctx),
			LastTxPoolId:              k.getID(ctx, types.KeyLastTXPoolID),
			LastBatchId:               k.getID(ctx, types.KeyLastOutgoingBatchID),
		},
		Valsets:                      valsets,
		ValsetConfirms:               vsconfs,
		Batches:                      extBatches,
		BatchConfirms:                batchconfs,
		LogicCalls:                   calls,
		LogicCallConfirms:            callconfs,
		Attestations:                 attestations,
		DelegateKeys:                 delegates,
		Erc20ToDenoms:                erc20ToDenoms,
		UnbatchedTransfers:           unbatchedTxs,
		PendingIbcAutoForwards:       forwards,
		PendingErc721IbcAutoForwards: erc721Forwards,
		Erc721Attestations:           erc721Attestations,
>>>>>>> 81057dc97ff3a6f3702fca99300ddbb3a7011770
	}
}
