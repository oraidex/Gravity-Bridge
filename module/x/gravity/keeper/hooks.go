package keeper

import (
	"context"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/slashing/types"
)

// Wrapper struct
type Hooks struct {
	k Keeper
}

// Implements StakingHooks interface
var _ types.StakingHooks = Hooks{}

// Create new gravity hooks
func (k Keeper) Hooks() Hooks {
	// if startup is mis-ordered in app.go this hook will halt
	// the chain when called. Keep this check to make such a mistake
	// obvious
	if k.storeKey == nil {
		panic("Hooks initialized before GravityKeeper!")
	}
	return Hooks{k}
}

func (h Hooks) AfterValidatorBeginUnbonding(ctx context.Context, _ sdk.ConsAddress, _ sdk.ValAddress) error {

	// When Validator starts Unbonding, Persist the block height in the store
	// Later in endblocker, check if there is at least one validator who started unbonding and create a valset request.
	// The reason for creating valset requests in endblock is to create only one valset request per block,
	// if multiple validators starts unbonding at same block.

	// this hook IS called for jailing or unbonding triggered by users but it IS NOT called for jailing triggered
	// in the endblocker therefore we call the keeper function ourselves there.
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	h.k.SetLastUnBondingBlockHeight(sdkCtx, uint64(sdkCtx.BlockHeight()))
	return nil
}

func (h Hooks) BeforeDelegationCreated(_ context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	return nil
}
func (h Hooks) AfterValidatorCreated(ctx context.Context, valAddr sdk.ValAddress) error { return nil }
func (h Hooks) BeforeValidatorModified(_ context.Context, _ sdk.ValAddress) error       { return nil }
func (h Hooks) AfterValidatorBonded(_ context.Context, _ sdk.ConsAddress, _ sdk.ValAddress) error {
	return nil
}

func (h Hooks) BeforeDelegationRemoved(_ context.Context, _ sdk.AccAddress, _ sdk.ValAddress) error {
	return nil
}
func (h Hooks) AfterValidatorRemoved(ctx context.Context, _ sdk.ConsAddress, valAddr sdk.ValAddress) error {
	return nil
}
func (h Hooks) BeforeValidatorSlashed(ctx context.Context, valAddr sdk.ValAddress, fraction sdkmath.LegacyDec) error {
	return nil
}
func (h Hooks) BeforeDelegationSharesModified(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	return nil
}
func (h Hooks) AfterDelegationModified(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	return nil
}
func (h Hooks) AfterUnbondingInitiated(_ context.Context, _ uint64) error {
	return nil
}
