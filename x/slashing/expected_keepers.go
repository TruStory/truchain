package slashing

import (
	"time"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Stake is the palceholder here until Staking module is done
type Stake struct {
    ID           uint64
    ArgumentID   uint64
    Type         StakeType
    Amount       sdk.Coin
    Creator      sdk.AccAddress
    CreatedTime  time.Time
    EndTime      time.Time
}

// StakeType enum
type StakeType int

// StakeKeeper is the expected Staking keeper interface for this module
type StakeKeeper interface {
	Stake(ctx sdk.Context, id uint64) (Stake, sdk.Error)
}