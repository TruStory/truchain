package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// NewCategoryCoin creates a new category coin type
func NewCategoryCoin(toDenom string, from sdk.Coin) sdk.Coin {
	rate := exchangeRate(from, toDenom)

	return sdk.NewCoin(
		toDenom,
		sdk.NewDecFromInt(from.Amount).Mul(rate).RoundInt())
}

// SwapForCategoryCoin swaps any coin for category coin
func SwapForCategoryCoin(
	ctx sdk.Context,
	bankKeeper bank.Keeper,
	from sdk.Coin,
	toDenom string,
	user sdk.AccAddress) (coin sdk.Coin, err sdk.Error) {

	coin = NewCategoryCoin(toDenom, from)

	err = SwapCoin(ctx, bankKeeper, from, coin, user)

	return
}

// SwapCoin replaces one coin with another
func SwapCoin(
	ctx sdk.Context,
	bankKeeper bank.Keeper,
	from sdk.Coin,
	to sdk.Coin,
	user sdk.AccAddress) (err sdk.Error) {

	_, _, err = bankKeeper.SubtractCoins(ctx, user, sdk.Coins{from})
	if err != nil {
		return
	}

	_, _, err = bankKeeper.AddCoins(ctx, user, sdk.Coins{to})

	return
}

// naïve implementation: 1 trustake = 1 category coin
// TODO [Shane]: https://github.com/TruStory/truchain/issues/21
func exchangeRate(from sdk.Coin, toDenom string) sdk.Dec {

	if from.Denom == toDenom {
		return sdk.NewDec(1)
	}
	return sdk.NewDec(1)
}
