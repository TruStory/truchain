package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewCategoryCoin generates a new category coin
func NewCategoryCoin(toDenom string, from sdk.Coin) sdk.Coin {
	rate := ExchangeRate(from, toDenom)

	return sdk.NewCoin(
		toDenom,
		sdk.NewDecFromInt(from.Amount).Mul(rate).RoundInt())
}

// ExchangeRate exchanges coins from one denom to the other
func ExchangeRate(from sdk.Coin, toDenom string) sdk.Dec {

	if from.Denom == toDenom {
		return sdk.NewDec(1)
	}
	return sdk.NewDec(1)
}
