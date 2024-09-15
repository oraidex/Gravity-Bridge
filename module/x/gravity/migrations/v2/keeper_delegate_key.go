package v2

import (
	"github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravity/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetEthAddressByValidator returns the eth address for a given gravity validator
func GetEthAddressByValidator(ctx sdk.Context, validator sdk.ValAddress, store sdk.KVStore) (ethAddress *types.EthAddress, found bool) {
	if err := sdk.VerifyAddressFormat(validator); err != nil {
		panic(errorsmod.Wrap(err, "invalid validator address"))
	}

	ethAddr := store.Get(GetEthAddressByValidatorKey(validator))
	if ethAddr == nil {
		return nil, false
	}

	addr, err := types.NewEthAddressFromBytes(ethAddr)
	if err != nil {
		return nil, false
	}
	return addr, true
}
