package voting

import (
	"fmt"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TokenVote defines a simple token vote on a story
type TokenVote struct {
	app.Vote `json:"vote"`
}

// ID implements `Voter`
func (v TokenVote) ID() int64 {
	return v.Vote.ID
}

// StoryID implements `Voter`
func (v TokenVote) StoryID() int64 {
	return v.Vote.ID
}

// Amount implements `Voter`
func (v TokenVote) Amount() sdk.Coin {
	return v.Vote.Amount
}

// Creator implements `Voter`
func (v TokenVote) Creator() sdk.AccAddress {
	return v.Vote.Creator
}

// VoteChoice implements `Voter`
func (v TokenVote) VoteChoice() bool {
	return v.Vote.Vote
}

type poll struct {
	trueVotes  []app.Voter
	falseVotes []app.Voter
}

func (p poll) String() string {
	return fmt.Sprintf(
		"Poll results:\n True votes: %v\n False votes: %v",
		p.trueVotes, p.falseVotes)
}
