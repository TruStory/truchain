package community

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Defines module constants
const (
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
	StoreKey     = ModuleName
)

// Community represents the state of a community on TruStory
type Community struct {
	ID               uint64    `json:"id"`
	Name             string    `json:"name"`
	Slug             string    `json:"slug"`
	Description      string    `json:"description,omitempty"`
	TotalEarnedStake sdk.Coin  `json:"total_earned_stake"`
	CreatedTime      time.Time `json:"created_time,omitempty"`
}

func (c Community) String() string {
	return fmt.Sprintf("Community <%d %s %s %s>", c.ID, c.Name, c.Slug, c.Description)
}
