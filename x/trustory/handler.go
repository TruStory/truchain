package trustory

import (
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler creates a new handler for all TruStory type messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case SubmitStoryMsg:
			return handleSubmitStoryMsg(ctx, k, msg)
		case VoteMsg:
			return handleVoteMsg(ctx, k, msg)
		default:
			errMsg := "Unrecognized Msg type: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// handleSubmitStoryMsg handles the logic of a SubmitStoryMsg
func handleSubmitStoryMsg(ctx sdk.Context, k Keeper, msg SubmitStoryMsg) sdk.Result {
	err := msg.ValidateBasic()
	if err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

func handleVoteMsg(ctx sdk.Context, k Keeper, msg VoteMsg) sdk.Result {
	err := msg.ValidateBasic()
	if err != nil {
		return err.Result()
	}

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
func checkStory(ctx sdk.Context, k Keeper) sdk.Error {
	return nil
}
