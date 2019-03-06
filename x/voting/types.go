package voting

import (
	"fmt"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type poll struct {
	trueVotes  []app.Voter
	falseVotes []app.Voter
}

// VoteResults are the vote results for a story
type VoteResults struct {
	ID                  int64         `json:"id,omitempty"`
	BackedCredTotal     sdk.Int       `json:"backed_cred_total,omitempty"`
	ChallengedCredTotal sdk.Int       `json:"challenged_cred_total,omitempty"`
	Timestamp           app.Timestamp `json:"timestamp,omitempty"`
}

func (p poll) String() string {
	return fmt.Sprintf(
		"Poll results:\n True votes: %v\n False votes: %v",
		p.trueVotes, p.falseVotes)
}
