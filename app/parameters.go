package app

import (
  sdk "github.com/cosmos/cosmos-sdk/types"
  "github.com/cosmos/cosmos-sdk/x/auth"
)

// TODO: Make values in this file configurable from environment [notduncansmith]

const (
  AppName    = "TruChain"
  StakeDenom = "trusteak"
  Hostname   = "0.0.0.0"
  Portname   = "8080"
)

// TODO: Update with actual user initial coins [notduncansmith]
var initialCoins = sdk.Coins{
  sdk.Coin{Amount: sdk.NewInt(123456), Denom: StakeDenom},
}

// TODO: Use more accurate gas estimate [notduncansmith]
var registrationFee = auth.StdFee{
  Amount: sdk.Coins{sdk.Coin{Amount: sdk.NewInt(1), Denom: StakeDenom}},
  Gas:    10000,
}
