package auth

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Defines auth module constants
const (
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
)

// Auth stores data about a auth
type Auth struct {
	ID              int64
	Creator         sdk.AccAddress
	CreatedTime     time.Time
}
