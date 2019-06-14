package slashing

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Defines slashing module constants
const (
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
)

// Slashing stores data about a slashing
type Slashing struct {
	ID              int64
	Creator         sdk.AccAddress
	CreatedTime     time.Time
}
