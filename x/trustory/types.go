package trustory

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Bond placed on a story
type Bond struct {
	ID           int64          `json:"id"`
	StoryID      int64          `json:"story_id"`
	Stake        sdk.Coins      `json:"amount"`
	CreatedBlock int64          `json:"created_block"`
	Creator      sdk.AccAddress `json:"creator"`
	Period       time.Time      `json:"period"`
}

// NewBond creates a new bond
func NewBond(
	id int64,
	storyID int64,
	stake sdk.Coins,
	createdBlock int64,
	creator sdk.AccAddress,
	period time.Time) Bond {
	return Bond{
		ID:           id,
		StoryID:      storyID,
		Stake:        stake,
		CreatedBlock: createdBlock,
		Creator:      creator,
		Period:       period,
	}
}

// Comment for a story
type Comment struct {
	ID      int64          `json:"id"`
	StoryID int64          `json:"story_id"`
	Body    string         `json:"body"`
	Creator sdk.AccAddress `json:"creator"`
}

// NewComment creates a new comment for a given story
func NewComment(id int64, storyID int64, body string, creator sdk.AccAddress) Comment {
	return Comment{
		ID:      id,
		StoryID: storyID,
		Body:    body,
		Creator: creator,
	}
}

// Evidence for a story
type Evidence struct {
	ID      int64          `json:"id"`
	StoryID int64          `json:"story_id"`
	Creator sdk.AccAddress `json:"creator"`
	URI     string         `json:"uri"`
}

// NewEvidence creates new evidence for a story
func NewEvidence(id int64, storyID int64, creator sdk.AccAddress, uri string) Evidence {
	return Evidence{
		ID:      id,
		StoryID: storyID,
		Creator: creator,
		URI:     uri,
	}
}

// Story type
type Story struct {
	ID           int64            `json:"id"`
	BondIDs      []int64          `json:"bond_i_ds,omitempty"`
	CommentIDs   []int64          `json:"comment_i_ds,omitempty"`
	EvidenceIDs  []int64          `json:"evidence_i_ds,omitempty"`
	Thread       []int64          `json:"thread,omitempty"`
	VoteIDs      []int64          `json:"vote_i_ds"`
	Body         string           `json:"body"`
	Category     string           `json:"category"`
	CreatedBlock int64            `json:"created_block"`
	Creator      sdk.AccAddress   `json:"creator"`
	Expiration   time.Time        `json:"expiration,omitempty"`
	Rank         float64          `json:"rank"`
	State        string           `json:"state"`
	SubmitBlock  int64            `json:"submit_block"`
	StoryType    string           `json:"type"`
	UpdatedBlock int64            `json:"updated_block"`
	Users        []sdk.AccAddress `json:"users"`
}

// NewStory creates a new story
func NewStory(
	id int64,
	bondIDs []int64,
	commentIDs []int64,
	evidenceIDs []int64,
	thread []int64,
	voteIDs []int64,
	body string,
	category string,
	createdBlock int64,
	creator sdk.AccAddress,
	expiration time.Time,
	rank float64,
	state string,
	submitBlock int64,
	storyType string,
	updatedBlock int64,
	users []sdk.AccAddress) Story {
	return Story{
		ID:           id,
		BondIDs:      bondIDs,
		CommentIDs:   commentIDs,
		EvidenceIDs:  evidenceIDs,
		Thread:       thread,
		VoteIDs:      voteIDs,
		Body:         body,
		Category:     category,
		CreatedBlock: createdBlock,
		Creator:      creator,
		Expiration:   expiration,
		Rank:         rank,
		State:        "Created",
		SubmitBlock:  submitBlock,
		StoryType:    storyType,
		UpdatedBlock: updatedBlock,
		Users:        users,
	}
}

// Vote for a story
type Vote struct {
	ID           int64          `json:"id"`
	CreatedBlock int64          `json:"created_block"`
	Creator      sdk.AccAddress `json:"creator"`
	StoryID      int64          `json:"story_id"`
	Vote         bool           `json:"vote"`
}

// NewVote creates a new vote for a story
func NewVote(
	id int64,
	storyID int64,
	createdBlock int64,
	creator sdk.AccAddress,
	vote bool) Vote {
	return Vote{
		ID:           id,
		StoryID:      storyID,
		CreatedBlock: createdBlock,
		Creator:      creator,
		Vote:         vote,
	}
}
