package parameters

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

// TODO: Make values in this file configurable from environment [notduncansmith]

const (
	// AppName is the name of the Cosmos app
	AppName = "TruChain"

	// StakeDenom is the name of the main staking currency (will be "trustake" on mainnet launch)
	StakeDenom = "trusteak"

	// Hostname is the address the app's HTTP server will bind to
	Hostname = "0.0.0.0"

	// Portname is the port the app's HTTP server will bind to
	Portname = "1337"
)

// InitialCredAmount is the initial amount of cred for categories
var InitialCredAmount = sdk.NewInt(1000000000)

// InitialTruStake is an `sdk.Coins` representing the balance a new user is granted upon registration
var InitialTruStake = sdk.Coin{Amount: sdk.NewInt(10000000000000), Denom: StakeDenom}

// RegistrationFee is an `auth.StdFee` representing the coin and gas cost of registering a new account
// TODO: Use more accurate gas estimate [notduncansmith]
var RegistrationFee = auth.StdFee{
	Amount: sdk.Coins{sdk.Coin{Amount: sdk.NewInt(1), Denom: StakeDenom}},
	Gas:    10000,
}

// Fee is for spam prevention and validator rewards
var Fee = sdk.Coins{
	sdk.Coin{Amount: sdk.NewInt(10), Denom: StakeDenom},
}

// Feature flags
const (
	FeeFlag = iota
	BootstrapFlag
)

// Features sets flags on features to turn on/off during testnet
var Features = map[int]bool{
	FeeFlag:       false,
	BootstrapFlag: false,
}
