package account

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankexported "github.com/TruStory/truchain/x/bank/exported"
)

// BankKeeper is the expected bank keeper interface for this module
type BankKeeper interface {
	AddCoin(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coin,
		referenceID uint64, txType bankexported.TransactionType) (sdk.Coins, sdk.Error)
}
