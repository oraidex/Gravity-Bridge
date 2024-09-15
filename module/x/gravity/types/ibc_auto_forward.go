package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// ValidateBasic checks the ForeignReceiver is valid and foreign, the Amount is non-zero, the IbcChannel is
// non-empty, and the EventNonce is non-zero
func (p PendingIbcAutoForward) ValidateBasic() error {

	if p.ForeignReceiver == "" {
		return errorsmod.Wrapf(ErrInvalid, "ForeignReceiver is empty")
	}

	if p.Token.IsZero() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "Token must be non-zero")
	}

	if p.IbcChannel == "" {
		return errorsmod.Wrap(ErrInvalid, "IbcChannel must not be an empty string")
	}

	if p.EventNonce == 0 {
		return errorsmod.Wrap(ErrInvalid, "EventNonce must be non-zero")
	}

	return nil
}
