package auth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BankKeeper is the expected bank keeper interface for this module
type BankKeeper interface {
	NewTransaction(ctx sdk.Context, to sdk.AccAddress, coins sdk.Coins) bool
}