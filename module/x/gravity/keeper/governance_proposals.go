package keeper

import (
	"fmt"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	v3 "github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravity/migrations/v3"
	"github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravity/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

// this file contains code related to custom governance proposals

func RegisterProposalTypes() {
	// use of prefix stripping to prevent a typo between the proposal we check
	// and the one we register, any issues with the registration string will prevent
	// the proposal from working. We must check for double registration so that cli commands
	// submitting these proposals work.
	// For some reason the cli code is run during app.go startup, but of course app.go is not
	// run during operation of one off tx commands, so we need to run this 'twice'
	prefix := "gravity/"
	metadata := "gravity/IBCMetadata"
	if !govtypes.IsValidProposalType(strings.TrimPrefix(metadata, prefix)) {
		govtypes.RegisterProposalType(types.ProposalTypeIBCMetadata)
	}
	unhalt := "gravity/UnhaltBridge"
	if !govtypes.IsValidProposalType(strings.TrimPrefix(unhalt, prefix)) {
		govtypes.RegisterProposalType(types.ProposalTypeUnhaltBridge)

	}
	airdrop := "gravity/Airdrop"
	if !govtypes.IsValidProposalType(strings.TrimPrefix(airdrop, prefix)) {
		govtypes.RegisterProposalType(types.ProposalTypeAirdrop)

	}
	addEvmChain := "gravity/AddEvmChain"
	if !govtypes.IsValidProposalType(strings.TrimPrefix(addEvmChain, prefix)) {
		govtypes.RegisterProposalType(types.ProposalTypeAddEvmChain)

	}

	removeEvmChain := "gravity/RemoveEvmChain"
	if !govtypes.IsValidProposalType(strings.TrimPrefix(removeEvmChain, prefix)) {
		govtypes.RegisterProposalType(types.ProposalTypeRemoveEvmChain)

	}

	monitoredERC20Tokens := "gravity/MonitoredERC20Tokens"
	if !govtypes.IsValidProposalType(strings.TrimPrefix(monitoredERC20Tokens, prefix)) {
		govtypes.RegisterProposalType(types.ProposalTypeMonitoredERC20Tokens)

	}
}

func NewGravityProposalHandler(k Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.UnhaltBridgeProposal:
			return k.HandleUnhaltBridgeProposal(ctx, c)
		case *types.AirdropProposal:
			return k.HandleAirdropProposal(ctx, c)
		case *types.IBCMetadataProposal:
			return k.HandleIBCMetadataProposal(ctx, c)
		case *types.AddEvmChainProposal:
			return k.HandleAddEvmChainProposal(ctx, c)
		case *types.RemoveEvmChainProposal:
			return k.HandleRemoveEvmChainProposal(ctx, c)
		case *types.MonitoredERC20TokensProposal:
			return k.HandleMonitoredERC20TokensProposal(ctx, c)

		default:
			return errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized Gravity proposal content type: %T", c)
		}
	}
}

// Unhalt Bridge specific functions

// In the event the bridge is halted and governance has decided to reset oracle
// history, we roll back oracle history and reset the parameters
func (k Keeper) HandleUnhaltBridgeProposal(ctx sdk.Context, p *types.UnhaltBridgeProposal) error {
	ctx.Logger().Info("Gov vote passed: Resetting oracle history", "nonce", p.TargetNonce)
	pruneAttestationsAfterNonce(ctx, p.EvmChainPrefix, k, p.TargetNonce)
	return nil
}

// In the event we need to add new evm chains, we can create a new proposal
func (k Keeper) HandleAddEvmChainProposal(ctx sdk.Context, p *types.AddEvmChainProposal) error {

	// isEvmChainExist := k.GetEvmChainData(ctx, p.EvmChainPrefix)
	// if isEvmChainExist != nil {
	// 	return errorsmod.Wrap(types.ErrInvalid, "The proposed EVM Chain already exists on-chain. Cannot re-add it!")
	// }

	// evmChains := k.GetEvmChains(ctx)
	// for _, chain := range evmChains {
	// 	if chain.EvmChainNetVersion == p.EvmChainNetVersion {
	// 		return errorsmod.Wrap(types.ErrInvalid, "The proposed EVM Chain net version already exists on-chain. Cannot add a new chain with the same net version")
	// 	}
	// }

	ctx.Logger().Info("Gov vote passed: Adding new EVM chain", "evm chain prefix", p.EvmChainPrefix)
	evmChain := types.EvmChainData{
		EvmChain:           types.EvmChain{EvmChainPrefix: p.EvmChainPrefix, EvmChainName: p.EvmChainName, EvmChainNetVersion: p.EvmChainNetVersion},
		GravityNonces:      types.GravityNonces{},
		Valsets:            []types.Valset{},
		ValsetConfirms:     []types.MsgValsetConfirm{},
		Batches:            []types.OutgoingTxBatch{},
		BatchConfirms:      []types.MsgConfirmBatch{},
		LogicCalls:         []types.OutgoingLogicCall{},
		LogicCallConfirms:  []types.MsgConfirmLogicCall{},
		Attestations:       []types.Attestation{},
		DelegateKeys:       []types.MsgSetOrchestratorAddress{},
		Erc20ToDenoms:      []types.ERC20ToDenom{},
		UnbatchedTransfers: []types.OutgoingTransferTx{},
	}
	k.SetEvmChainData(ctx, evmChain.EvmChain)

	chainPrefix := p.EvmChainPrefix
	k.SetLatestValsetNonce(ctx, chainPrefix, evmChain.GravityNonces.LatestValsetNonce)
	k.setLastObservedEventNonce(ctx, chainPrefix, evmChain.GravityNonces.LastObservedNonce)
	k.SetLastSlashedValsetNonce(ctx, chainPrefix, evmChain.GravityNonces.LastSlashedValsetNonce)
	k.SetLastSlashedBatchBlock(ctx, chainPrefix, evmChain.GravityNonces.LastSlashedBatchBlock)
	k.SetLastSlashedLogicCallBlock(ctx, chainPrefix, evmChain.GravityNonces.LastSlashedLogicCallBlock)
	k.SetLastObservedEvmChainBlockHeight(ctx, chainPrefix, evmChain.GravityNonces.LastObservedEvmBlockHeight)
	k.setID(ctx, evmChain.GravityNonces.LastTxPoolId, types.AppendChainPrefix(types.KeyLastTXPoolID, chainPrefix))
	k.setID(ctx, evmChain.GravityNonces.LastBatchId, types.AppendChainPrefix(types.KeyLastOutgoingBatchID, chainPrefix))

	initBridgeDataFromGenesis(ctx, k, evmChain)

	// check bridge address, if invalid then we set default to 0x0
	finalEthAddress := p.BridgeEthereumAddress
	err := types.ValidateEthAddress(p.BridgeEthereumAddress)
	if err != nil {
		finalEthAddress = "0x0000000000000000000000000000000000000000"
	}

	// update param to match with the new evm chain
	params := k.GetParams(ctx)
	evmChainParam := &types.EvmChainParam{
		EvmChainPrefix:           p.EvmChainPrefix,
		GravityId:                p.GravityId,
		ContractSourceHash:       "",
		BridgeEthereumAddress:    finalEthAddress,
		BridgeChainId:            p.EvmChainNetVersion,
		AverageEthereumBlockTime: 15000,
		BridgeActive:             true,
		EthereumBlacklist:        []string{},
	}

	var evmChainParams []*types.EvmChainParam

	exists := false
	for _, param := range params.EvmChainParams {
		if param.EvmChainPrefix == evmChainParam.EvmChainPrefix {
			evmChainParams = append(evmChainParams, evmChainParam)
			exists = true
		} else {
			evmChainParams = append(evmChainParams, param)
		}
	}
	if !exists {
		evmChainParams = append(evmChainParams, evmChainParam)
	}
	params.EvmChainParams = evmChainParams
	k.SetParams(ctx, params)
	return nil
}

// In the event we need to remove an evm chains, we can create a new proposal, but call remove evm chain method from store migration
func (k Keeper) HandleRemoveEvmChainProposal(ctx sdk.Context, p *types.RemoveEvmChainProposal) error {
	// remove params for current evm chain first
	var evmChainParams []*types.EvmChainParam

	params := k.GetParams(ctx)
	for _, param := range params.EvmChainParams {
		if param.EvmChainPrefix == p.EvmChainPrefix {
			continue
		}
		evmChainParams = append(evmChainParams, param)

	}

	params.EvmChainParams = evmChainParams
	k.SetParams(ctx, params)
	return v3.RemoveEvmChainFromStore(ctx, k.storeKey, k.cdc, p.EvmChainPrefix)
}

// Iterate over all attestations currently being voted on in order of nonce
// and prune those that are older than nonceCutoff
func pruneAttestationsAfterNonce(ctx sdk.Context, evmChainPrefix string, k Keeper, nonceCutoff uint64) {
	// Decide on the most recent nonce we can actually roll back to
	lastObserved := k.GetLastObservedEventNonce(ctx, evmChainPrefix)
	if nonceCutoff < lastObserved || nonceCutoff == 0 {
		ctx.Logger().Error("Attempted to reset to a nonce before the last \"observed\" event, which is not allowed", "lastObserved", lastObserved, "nonce", nonceCutoff)
		return
	}

	// Get relevant event nonces
	attmap, keys := k.GetAttestationMapping(ctx, evmChainPrefix)

	// Discover all affected validators whose LastEventNonce must be reset to nonceCutoff
	power, err := k.StakingKeeper.GetBondedValidatorsByPower(ctx)
	if err != nil {
		panic(err)
	}

	numValidators := len(power)
	// void and setMember are necessary for sets to work
	type void struct{}
	var setMember void
	// Initialize a Set of validators
	affectedValidatorsSet := make(map[string]void, numValidators)

	// Delete all reverted attestations, keeping track of the validators who attested to any of them
	for _, nonce := range keys {
		for _, att := range attmap[nonce] {
			// we delete all attestations earlier than the cutoff event nonce
			if nonce > nonceCutoff {
				ctx.Logger().Info(fmt.Sprintf("Deleting attestation at height %v", att.Height))
				for _, vote := range att.Votes {
					if _, ok := affectedValidatorsSet[vote]; !ok { // if set does not contain vote
						affectedValidatorsSet[vote] = setMember // add key to set
					}
				}

				k.DeleteAttestation(ctx, att)
			}
		}
	}

	// Reset the last event nonce for all validators affected by history deletion
	for vote := range affectedValidatorsSet {
		val, err := sdk.ValAddressFromBech32(vote)
		if err != nil {
			panic(errorsmod.Wrap(err, "invalid validator address affected by bridge reset"))
		}
		valLastNonce := k.GetLastEventNonceByValidator(ctx, evmChainPrefix, val)
		if valLastNonce > nonceCutoff {
			ctx.Logger().Info("Resetting validator's last event nonce due to bridge unhalt", "validator", vote, "lastEventNonce", valLastNonce, "resetNonce", nonceCutoff)
			k.SetLastEventNonceByValidator(ctx, evmChainPrefix, val, nonceCutoff)
		}
	}
}

// Allows governance to deploy an airdrop to a provided list of addresses
func (k Keeper) HandleAirdropProposal(ctx sdk.Context, p *types.AirdropProposal) error {
	ctx.Logger().Info("Gov vote passed: Performing airdrop")
	startingSupply := k.bankKeeper.GetSupply(ctx, p.Denom)

	validateDenom := sdk.ValidateDenom(p.Denom)
	if validateDenom != nil {
		ctx.Logger().Info("Airdrop failed to execute invalid denom!")
		return errorsmod.Wrap(types.ErrInvalid, "Invalid airdrop denom")
	}

	feePool, err := k.DistKeeper.FeePool.Get(ctx)
	if err != nil {
		panic(err)
	}
	feePoolAmount := feePool.CommunityPool.AmountOf(p.Denom)

	airdropTotal := sdkmath.NewInt(0)
	for _, v := range p.Amounts {
		airdropTotal = airdropTotal.Add(sdkmath.NewIntFromUint64(v))
	}

	totalRequiredDecCoin := sdk.NewDecCoinFromCoin(sdk.NewCoin(p.Denom, airdropTotal))

	// check that we have enough tokens in the community pool to actually execute
	// this airdrop with the provided recipients list
	totalRequiredDec := totalRequiredDecCoin.Amount
	if totalRequiredDec.GT(feePoolAmount) {
		ctx.Logger().Info("Airdrop failed to execute insufficient tokens in the community pool!")
		return errorsmod.Wrap(types.ErrInvalid, "Insufficient tokens in community pool")
	}

	// we're packing addresses as 20 bytes rather than valid bech32 in order to maximize participants
	// so if the recipients list is not a multiple of 20 it must be invalid
	numRecipients := len(p.Recipients) / 20
	if len(p.Recipients)%20 != 0 || numRecipients != len(p.Amounts) {
		ctx.Logger().Info("Airdrop failed to execute invalid recipients")
		return errorsmod.Wrap(types.ErrInvalid, "Invalid recipients")
	}

	parsedRecipients := make([]sdk.AccAddress, len(p.Recipients)/20)
	for i := 0; i < numRecipients; i++ {
		indexStart := i * 20
		indexEnd := indexStart + 20
		addr := p.Recipients[indexStart:indexEnd]
		parsedRecipients[i] = addr
	}

	// check again, just in case the above modulo math is somehow wrong or spoofed
	if len(parsedRecipients) != len(p.Amounts) {
		ctx.Logger().Info("Airdrop failed to execute invalid recipients")
		return errorsmod.Wrap(types.ErrInvalid, "Invalid recipients")
	}

	// the total amount actually sent in dec coins
	totalSent := sdkmath.LegacyNewDec(0)
	for i, addr := range parsedRecipients {
		usersAmount := p.Amounts[i]
		usersIntAmount := sdkmath.NewIntFromUint64(usersAmount)
		usersDecAmount := sdkmath.LegacyNewDecFromInt(usersIntAmount)
		err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, disttypes.ModuleName, addr, sdk.NewCoins(sdk.NewCoin(p.Denom, usersIntAmount)))
		// if there is no error we add to the total actually sent
		if err == nil {
			totalSent = totalSent.Add(usersDecAmount)
		} else {
			// return an err to prevent execution from finishing, this will prevent the changes we
			// have made so far from taking effect the governance proposal will instead time out
			ctx.Logger().Info("invalid address in airdrop! not executing", "address", addr)
			return err
		}
	}

	if !totalRequiredDecCoin.Amount.Equal(totalSent) {
		ctx.Logger().Info("Airdrop failed to execute Invalid amount sent", "sent", totalRequiredDecCoin.Amount, "expected", totalSent)
		return errorsmod.Wrap(types.ErrInvalid, "Invalid amount sent")
	}

	newCoins, InvalidModuleBalance := feePool.CommunityPool.SafeSub(sdk.NewDecCoins(totalRequiredDecCoin))
	// this shouldn't ever happen because we check that we have enough before starting
	// but lets be conservative.
	if InvalidModuleBalance {
		return errorsmod.Wrap(types.ErrInvalid, "internal error!")
	}
	feePool.CommunityPool = newCoins
	err = k.DistKeeper.FeePool.Set(ctx, feePool)
	if err != nil {
		return err
	}

	endingSupply := k.bankKeeper.GetSupply(ctx, p.Denom)
	if !startingSupply.Equal(endingSupply) {
		return errorsmod.Wrap(types.ErrInvalid, "total chain supply has changed!")
	}

	return nil
}

// handles a governance proposal for setting the metadata of an IBC token, this takes the normal
// metadata struct with one key difference, the base unit must be set as the ibc path string in order
// for setting the denom metadata to work.
func (k Keeper) HandleIBCMetadataProposal(ctx sdk.Context, p *types.IBCMetadataProposal) error {
	ctx.Logger().Info("Gov vote passed: Setting IBC Metadata", "denom", p.IbcDenom)

	// checks if the provided token denom is a proper IBC token, not a native token.
	if !strings.HasPrefix(p.IbcDenom, "ibc/") && !strings.HasPrefix(p.IbcDenom, "IBC/") {
		ctx.Logger().Info("invalid denom for metadata proposal", "denom", p.IbcDenom)
		return errorsmod.Wrap(types.ErrInvalid, "Target denom is not an IBC token")
	}

	// check that our base unit is the IBC token name on this chain. This makes setting/loading denom
	// metadata work out, as SetDenomMetadata uses the base denom as an index
	if p.Metadata.Base != p.IbcDenom {
		ctx.Logger().Info("invalid metadata for metadata proposal must be the same as IBCDenom", "base", p.Metadata.Base)
		return errorsmod.Wrap(types.ErrInvalid, "Metadata base must be the same as the IBC denom!")
	}

	// outsource validating this to the bank validation function
	metadataErr := p.Metadata.Validate()
	if metadataErr != nil {
		ctx.Logger().Info("invalid metadata for metadata proposal", "validation error", metadataErr)
		return errorsmod.Wrap(metadataErr, "Invalid metadata")

	}

	// if metadata already exists then changing it is only a good idea if we have not already deployed an ERC20
	// for this denom if we have we can't change it
	_, metadataExists := k.bankKeeper.GetDenomMetaData(ctx, p.IbcDenom)
	_, erc20RepresentationExists := k.GetCosmosOriginatedERC20(ctx, p.EvmChainPrefix, p.IbcDenom)
	if metadataExists && erc20RepresentationExists {
		ctx.Logger().Info("invalid trying to set metadata when ERC20 has already been deployed")
		return errorsmod.Wrap(types.ErrInvalid, "Metadata can only be changed before ERC20 is created")
	}

	// write out metadata, this will update existing metadata if no erc20 has been deployed
	k.bankKeeper.SetDenomMetaData(ctx, p.Metadata)

	return nil
}

// handles a governance proposal for setting the metadata of an IBC token, this takes the normal
// metadata struct with one key difference, the base unit must be set as the ibc path string in order
// for setting the denom metadata to work.
func (k Keeper) HandleMonitoredERC20TokensProposal(ctx sdk.Context, p *types.MonitoredERC20TokensProposal) error {
	ctx.Logger().Info("Gov vote passed: Setting Monitored ERC20 Tokens", "tokens", p.Tokens)

	// checks each token to see if it is a valid address
	if err := p.ValidateBasic(); err != nil {
		return fmt.Errorf("Invalid MonitoredERC20TokensProposal: %v", err)
	}
	var tokens []types.EthAddress
	for _, t := range p.Tokens {
		// Address validation already occurred, so we can ignore address errors
		addr, _ := types.NewEthAddress(t)
		// Check that any cosmos originated denoms have an ERC20 representation
		denom, exists := k.GetCosmosOriginatedDenom(ctx, p.EvmChainPrefix, *addr)
		if exists && len(denom) > 0 {
			// The ERC20 is a cosmos originated denom, check that the token has been bridged
			contract, exists := k.GetCosmosOriginatedERC20(ctx, p.EvmChainPrefix, denom)
			if !exists {
				return fmt.Errorf(
					"Invalid MonitoredERC20TokensProposal: ERC20 token %v is cosmos originated (%v), but no ERC20 representation has been registered?",
					addr, denom,
				)
			}
			if contract.GetAddress().String() != addr.GetAddress().String() {
				return fmt.Errorf(
					"Invalid MonitoredERC20TokensProposal: ERC20 token %v is cosmos originated (%v), but the registered representation is different than expected %v",
					addr, denom, contract,
				)
			}
		}

		// If the above checks pass, add it to the list of contracts to use in the store
		tokens = append(tokens, *addr)
	}
	k.setMonitoredERC20Tokens(ctx, tokens)

	return nil
}
