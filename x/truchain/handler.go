package truchain

import (
	"encoding/binary"
	"reflect"
	"time"

	db "github.com/TruStory/truchain/x/truchain/db"
	ts "github.com/TruStory/truchain/x/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const votingPeriod = 24 // hours

// NewHandler creates a new handler for all TruStory messages
func NewHandler(k db.TruKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case ts.SubmitStoryMsg:
			return handleSubmitStoryMsg(ctx, k, msg)
		case ts.VoteMsg:
			return handleVoteMsg(ctx, k, msg)
		default:
			errMsg := "Unrecognized Msg type: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// ============================================================================

// handleSubmitStoryMsg handles the logic of a SubmitStoryMsg
func handleSubmitStoryMsg(ctx sdk.Context, k db.TruKeeper, msg ts.SubmitStoryMsg) sdk.Result {
	err := msg.ValidateBasic()
	if err != nil {
		return err.Result()
	}

	// calculate voting period
	voteStart := ctx.BlockHeader().Time
	voteEnd := voteStart.Add(time.Hour * time.Duration(votingPeriod))

	storyID, err := k.AddStory(ctx, msg.Body, msg.Category, msg.Creator, msg.Escrow, msg.StoryType, voteStart, voteEnd)
	if err != nil {
		panic(err)
	}

	return sdk.Result{Data: i2b(storyID)}
}

func handleVoteMsg(ctx sdk.Context, k db.TruKeeper, msg ts.VoteMsg) sdk.Result {
	err := msg.ValidateBasic()
	if err != nil {
		return err.Result()
	}

	voteID, err := k.VoteStory(ctx, msg.StoryID, msg.Creator, msg.Vote, sdk.Coins{msg.Amount})
	if err != nil {
		panic(err)
	}

	return sdk.Result{Data: i2b(voteID)}
}

// i2b converts an int64 into a byte array
func i2b(x int64) []byte {
	var b [binary.MaxVarintLen64]byte
	return b[:binary.PutVarint(b[:], x)]
}
