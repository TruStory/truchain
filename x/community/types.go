package community

import (
	"fmt"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Community represents the state of a community on TruStory
type Community struct {
	ID               int64
	Name             string
	Slug             string
	Description      string
	TotalEarnedStake sdk.Coin
	Timestamp        app.Timestamp
}

func (c Community) String() string {
	return fmt.Sprintf("Community <%d %s %s %s>", c.ID, c.Name, c.Slug, c.Description)
}

// MsgParams holds data for community parameters
type MsgParams struct {
	MinNameLen        int
	MaxNameLen        int
	MinSlugLen        int
	MaxSlugLen        int
	MaxDescriptionLen int
}

// DefaultMsgParams creates a new MsgParams type with the defaults
func DefaultMsgParams() MsgParams {
	return MsgParams{
		MinNameLen:        5,
		MaxNameLen:        25,
		MinSlugLen:        3,
		MaxSlugLen:        15,
		MaxDescriptionLen: 140,
	}
}
