package query

import (
	"github.com/TruStory/truchain/x/argument"
	"github.com/TruStory/truchain/x/backing"
	"github.com/TruStory/truchain/x/challenge"
	"github.com/TruStory/truchain/x/story"
	"github.com/TruStory/truchain/x/trubank"
)

// Keeper re
type Keeper struct {
	argumentKeeper  argument.Keeper
	backingKeeper   backing.ReadKeeper
	challengeKeeper challenge.ReadKeeper
	storyKeeper     story.ReadKeeper
	trubankKeeper   trubank.Keeper
}

// NewKeeper creates a new keeper with write and read access
func NewKeeper(
	argumentKeeper argument.Keeper,
	backingKeeper backing.ReadKeeper,
	challengeKeeper challenge.ReadKeeper,
	storyKeeper story.ReadKeeper,
	trubankKeeper trubank.Keeper) Keeper {
	return Keeper{
		argumentKeeper,
		backingKeeper,
		challengeKeeper,
		storyKeeper,
		trubankKeeper,
	}
}
