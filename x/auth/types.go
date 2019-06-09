package auth

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkAuth "github.com/cosmos/cosmos-sdk/x/auth"
)

// Defines auth module constants
const (
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
)

// EarnedCoin is a representation of a coin associated with a community, or "earned trustake".
type EarnedCoin struct {
	sdk.Coin

	CommunityID int64
}

// EarnedCoins is a collection of EarnedCoins
type EarnedCoins []EarnedCoin

// AppAccount is the main account for a TruStory user.
type AppAccount struct {
	sdkAuth.BaseAccount

	EarnedStake EarnedCoins
	SlashCount  int
	IsJailed    bool
	JailEndTime time.Time
	CreatedTime time.Time
}
