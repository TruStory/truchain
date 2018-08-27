package trustory

import (
	"reflect"

	db "github.com/TruStory/trucoin/x/trustory/db"
	ts "github.com/TruStory/trucoin/x/trustory/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

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

// handleSubmitStoryMsg handles the logic of a SubmitStoryMsg
func handleSubmitStoryMsg(ctx sdk.Context, k db.TruKeeper, msg ts.SubmitStoryMsg) sdk.Result {
	err := msg.ValidateBasic()
	if err != nil {
		return err.Result()
	}

	storyID, err := k.AddStory(ctx, msg.Body, msg.Category, msg.Creator, msg.StoryType)
	if err != nil {
		panic(err)
	}

	data, error := k.Cdc.MarshalBinary(storyID)
	if error != nil {
		panic(error)
	}

	return sdk.Result{Data: data}
}

func handleVoteMsg(ctx sdk.Context, k db.TruKeeper, msg ts.VoteMsg) sdk.Result {
	err := msg.ValidateBasic()
	if err != nil {
		return err.Result()
	}

	voteID, err := k.VoteStory(ctx, msg.StoryID, msg.Creator, msg.Vote, msg.Stake)
	if err != nil {
		panic(err)
	}

	data, error := k.Cdc.MarshalBinary(voteID)
	if error != nil {
		panic(error)
	}

	return sdk.Result{Data: data}
}

// NewEndBlocker checks stories and generates an EndBlocker
// func NewEndBlocker(k Keeper) sdk.EndBlocker {
// 	return func(ctx sdk.Context, req abci.RequestEndBlock) (res abci.ResponseEndBlock) {
// 		err := checkStory(ctx, k)
// 		if err != nil {
// 			panic(err)
// 		}
// 		return
// 	}
// }

// checkStory checks if the story reached the end of the voting period
// and handles the logic of ending voting
func checkStory(ctx sdk.Context, k db.TruKeeper) sdk.Error {
	return nil
}
