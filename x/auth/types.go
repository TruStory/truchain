package auth

import (
	"fmt"
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

	CommunityID uint64
}

// EarnedCoins is a collection of EarnedCoins
type EarnedCoins []EarnedCoin

// AppAccount is the main account for a TruStory user.
type AppAccount struct {
	sdkAuth.BaseAccount

	ID          uint64
	EarnedStake EarnedCoins
	SlashCount  int
	IsJailed    bool
	JailEndTime time.Time
	CreatedTime time.Time
}

func (appAccount AppAccount) String() string {
	return fmt.Sprintf("AppAccount <%d %s %s>", appAccount.ID, appAccount.BaseAccount.Address, appAccount.BaseAccount.PubKey)
}
