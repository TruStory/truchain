package backing

import (
	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler creates a new handler for all TruStory messages
func NewHandler(k WriteKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case BackStoryMsg:
			return handleBackStoryMsg(ctx, k, msg)
		case LikeBackingArgumentMsg:
			return handleLikeArgumentMsg(ctx, k, msg)
		default:
			return app.ErrMsgHandler(msg)
		}
	}
}

func handleBackStoryMsg(ctx sdk.Context, k WriteKeeper, msg BackStoryMsg) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	id, err := k.Create(
		ctx,
		msg.StoryID,
		msg.Amount,
		0,
		msg.Argument,
		msg.Creator,
		false)
	if err != nil {
		return err.Result()
	}

	return app.Result(id)
}

func handleLikeArgumentMsg(ctx sdk.Context, k WriteKeeper, msg LikeBackingArgumentMsg) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	backingID, err := k.LikeArgument(ctx, msg.ArgumentID, msg.Creator, msg.Amount)
	if err != nil {
		return err.Result()
	}

	return app.Result(backingID)
}
