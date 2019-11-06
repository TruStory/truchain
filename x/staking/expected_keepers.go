package staking

import (
	"time"

	"github.com/TruStory/truchain/x/account"
	bankexported "github.com/TruStory/truchain/x/bank/exported"
	"github.com/TruStory/truchain/x/claim"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type AccountKeeper interface {
	IsJailed(ctx sdk.Context, address sdk.AccAddress) (bool, sdk.Error)
	UnJail(ctx sdk.Context, address sdk.AccAddress) sdk.Error
	IterateAppAccounts(ctx sdk.Context, cb func(acc account.AppAccount) (stop bool))
}

type ClaimKeeper interface {
	Claim(ctx sdk.Context, id uint64) (claim claim.Claim, ok bool)
	AddBackingStake(ctx sdk.Context, id uint64, stake sdk.Coin) sdk.Error
	AddChallengeStake(ctx sdk.Context, id uint64, stake sdk.Coin) sdk.Error
	SubtractBackingStake(ctx sdk.Context, id uint64, stake sdk.Coin) sdk.Error
	SubtractChallengeStake(ctx sdk.Context, id uint64, stake sdk.Coin) sdk.Error
	SetFirstArgumentTime(ctx sdk.Context, id uint64, firstArgumentTime time.Time) sdk.Error
}

// BankKeeper is the expected bank keeper interface for this module
type BankKeeper interface {
	AddCoin(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coin,
		referenceID uint64, txType bankexported.TransactionType, setters ...bankexported.TransactionSetter) (sdk.Coins, sdk.Error)
	GetCoins(ctx sdk.Context, address sdk.AccAddress) sdk.Coins
	SubtractCoin(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coin,
		referenceID uint64, txType TransactionType, setters ...bankexported.TransactionSetter) (sdk.Coins, sdk.Error)
	TransactionsByAddress(ctx sdk.Context, address sdk.AccAddress, filterSetters ...bankexported.Filter) []bankexported.Transaction
	IterateUserTransactions(sdk.Context, sdk.AccAddress, bool, func(tx bankexported.Transaction) bool)
}
