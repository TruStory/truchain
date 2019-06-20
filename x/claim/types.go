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
	StoreKey     = ModuleName
)

// Association list keys
var (
	communityKey = sdk.NewKVStoreKey("community")
)

// Claim stores data about a claim
type Claim struct {
	ID              uint64         `json:"id"`
	CommunityID     uint64         `json:"community_id"`
	Body            string         `json:"body"`
	Creator         sdk.AccAddress `json:"creator"`
	Source          url.URL        `json:"source,omitempty"`
	TotalStakers    uint64         `json:"total_stakers,omitempty"`
	TotalBacked     sdk.Coin       `json:"total_backed,omitempty"`
	TotalChallenged sdk.Coin       `json:"total_challenged,omitempty"`
	CreatedTime     time.Time      `json:"created_time"`
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
