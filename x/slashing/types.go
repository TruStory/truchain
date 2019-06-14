package slashing

import (
	"fmt"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Defines slashing module constants
const (
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
)

// Slash stores data about a slashing
type Slash struct {
	ID        uint64
	StakeID   uint64
	Creator   sdk.AccAddress
	Timestamp app.Timestamp
}

func (s Slash) String() string {
	return fmt.Sprintf("Slash <%d %d %s>", s.ID, s.StakeID, s.Creator)
}
