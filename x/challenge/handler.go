package challenge

import (
	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler creates a new handler for all challenge messages
func NewHandler(k WriteKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case CreateChallengeMsg:
			return handleCreateChallengeMsg(ctx, k, msg)
		case LikeChallengeArgumentMsg:
			return handleLikeArgumentMsg(ctx, k, msg)
		default:
			return app.ErrMsgHandler(msg)
		}
	}
}

// handles a message to create a challenge
func handleCreateChallengeMsg(
	ctx sdk.Context, k WriteKeeper, msg CreateChallengeMsg) sdk.Result {

	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	id, err := k.Create(
		ctx, msg.StoryID, msg.Amount, 0, msg.Argument, msg.Creator, false)
	if err != nil {
		return err.Result()
	}

	return app.Result(id)
}

func handleLikeArgumentMsg(ctx sdk.Context, k WriteKeeper, msg LikeChallengeArgumentMsg) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	challengeID, err := k.LikeArgument(ctx, msg.ArgumentID, msg.Creator, msg.Amount)
	if err != nil {
		return err.Result()
	}

	return app.Result(challengeID)
}
