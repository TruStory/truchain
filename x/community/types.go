package community

import (
	"fmt"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Defines module constants
const (
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
)

// Community represents the state of a community on TruStory
type Community struct {
	ID               uint64
	Name             string
	Slug             string
	Description      string
	TotalEarnedStake sdk.Coin
	Timestamp        app.Timestamp
}

func (c Community) String() string {
	return fmt.Sprintf("Community <%d %s %s %s>", c.ID, c.Name, c.Slug, c.Description)
}
