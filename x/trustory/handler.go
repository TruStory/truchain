package trustory

import (
	"reflect"

	db "github.com/TruStory/trucoin/x/trustory/db"
	ts "github.com/TruStory/trucoin/x/trustory/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler creates a new handler for all TruStory messages
func NewHandler(sk db.StoryKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case ts.SubmitStoryMsg:
			return handleSubmitStoryMsg(ctx, sk, msg)
		case ts.VoteMsg:
			return handleVoteMsg(ctx, sk, msg)
		default:
			errMsg := "Unrecognized Msg type: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// handleSubmitStoryMsg handles the logic of a SubmitStoryMsg
func handleSubmitStoryMsg(ctx sdk.Context, sk db.StoryKeeper, msg ts.SubmitStoryMsg) sdk.Result {
	err := msg.ValidateBasic()
	if err != nil {
		return err.Result()
	}

	storyID, err := sk.AddStory(ctx, msg.Body, msg.Creator)
	if err != nil {
		panic(err)
	}

	data, error := sk.Cdc.MarshalBinary(storyID)
	if error != nil {
		panic(error)
	}

	return sdk.Result{Data: data}
}

func handleVoteMsg(ctx sdk.Context, sk db.StoryKeeper, msg ts.VoteMsg) sdk.Result {
	// err := msg.ValidateBasic()
	// if err != nil {
	// 	return err.Result()
	// }

	return sdk.Result{}
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
func checkStory(ctx sdk.Context, sk db.StoryKeeper) sdk.Error {
	return nil
}
