package account

import (
	"github.com/TruStory/truchain/x/bank"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BankKeeper is the expected bank keeper interface for this module
type BankKeeper interface {
	AddCoin(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coin,
		referenceID uint64, txType bank.TransactionType) (sdk.Coins, sdk.Error)
}
