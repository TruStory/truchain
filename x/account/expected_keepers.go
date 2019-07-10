package account

import (
	bankexported "github.com/TruStory/truchain/x/bank/exported"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BankKeeper is the expected bank keeper interface for this module
type BankKeeper interface {
	AddCoin(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coin,
		referenceID uint64, txType bankexported.TransactionType, setters ...bankexported.TransactionSetter) (sdk.Coins, sdk.Error)
}
