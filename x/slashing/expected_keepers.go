package slashing

import (
	"net/url"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Stake is the palceholder here until Staking module is done
type Stake struct {
	ID          uint64
	ArgumentID  uint64
	ClaimID     uint64
	Amount      sdk.Coin
	Creator     sdk.AccAddress
	CreatedTime time.Time
	EndTime     time.Time
	SlashCount  int
}

// StakeKeeper is the expected Staking keeper interface for this module
type StakeKeeper interface {
	Stake(ctx sdk.Context, id uint64) (Stake, sdk.Error)
	SlashCountByID(ctx sdk.Context, id uint64) (int, sdk.Error)
	IncrementSlashCount(ctx sdk.Context, id uint64) (sdk.Error)
}

// // AppAccount is the palceholder here until Auth module is done
// type AppAccount struct {
// 	Address          sdk.AccAddress
// 	IsJailed bool
// }

// // AppAccountKeeper is the expected AppAccount keeper interface for this module
// type AppAccountKeeper interface {
// 	AppAccount(ctx sdk.Context, address sdk.AccAddress) (auth.AppAccount, sdk.Error)
// 	IsJailed(ctx sdk.Context, address sdk.AccAddress) (bool, sdk.Error)
// 	JailUntil(ctx sdk.Context, address sdk.AccAddress, until time.Time) sdk.Error
// }

// Claim is the placeholder here until claim module is done
type Claim struct {
    ID                  uint64
    CommunityID         uint64
    Body                string
    Creator             sdk.AccAddress
    Source              url.URL
    TotalParticipants   int64
    TotalBackingStake   sdk.Coin
    TotalChallengeStake sdk.Coin
    CreatedTime         time.Time
}

// ClaimKeeper is the expected Claim keeper interface for this module
type ClaimKeeper interface {
	Claim(ctx sdk.Context, id uint64) (Claim, sdk.Error)
}