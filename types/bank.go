package types

import (
	params "github.com/TruStory/truchain/parameters"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// CategoryCoinFromTruStake creates category coins by burning trustake
func CategoryCoinFromTruStake(
	ctx sdk.Context, bankKeeper bank.Keeper, denom string, amount sdk.Coin, userAddr sdk.AccAddress) (
	catCoin sdk.Coin, err sdk.Error) {

	// naïve implementation: 1 trustake = 1 category coin
	// TODO [Shane]: https://github.com/TruStory/truchain/issues/21
	conversionRate := sdk.NewDec(1)

	// mint new category coins
	catCoin = sdk.NewCoin(
		denom,
		sdk.NewDecFromInt(amount.Amount).Mul(conversionRate).RoundInt())

	// burn equivalent trustake
	trustake := sdk.Coins{sdk.NewCoin(params.StakeDenom, catCoin.Amount)}
	_, _, err = bankKeeper.SubtractCoins(ctx, userAddr, trustake)
	if err != nil {
		return
	}

	return catCoin, nil
}
