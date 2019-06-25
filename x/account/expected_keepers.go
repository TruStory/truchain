package account


import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BankKeeper is the expected bank keeper interface for this module
type BankKeeper interface {
	AddCoin(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coin,
		referenceID uint64, txType int) (sdk.Coins, sdk.Error)
}
