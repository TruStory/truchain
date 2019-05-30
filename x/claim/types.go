package claim

import (
	"fmt"
	"net/url"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Defines module constants
const (
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
)

// Claim stores data about a claim
type Claim struct {
	ID                uint64
	CommunityID       uint64
	Body              string
	Creator           sdk.AccAddress
	Source            url.URL
	TotalParticipants uint64
	TotalBacked       sdk.Coin
	TotalChallenged   sdk.Coin
	CreatedTime       time.Time
}

// NewClaim creates a new claim object
func NewClaim(id, communityID uint64, body string, creator sdk.AccAddress, source url.URL, createdTime time.Time) Claim {
	return Claim{
		ID:          id,
		CommunityID: communityID,
		Body:        body,
		Creator:     creator,
		Source:      source,
		CreatedTime: createdTime,
	}
}

func (c Claim) String() string {
	return fmt.Sprintf("Claim <%d %s %s>", c.ID, c.Body, c.CreatedTime)
}
