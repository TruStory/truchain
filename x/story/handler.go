package story

import (
	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler creates a new handler for all TruStory messages
func NewHandler(k WriteKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case SubmitStoryMsg:
			return handleSubmitStoryMsg(ctx, k, msg)
		default:
			return app.ErrMsgHandler(msg)
		}
	}
}

// ============================================================================

func handleSubmitStoryMsg(ctx sdk.Context, k WriteKeeper, msg SubmitStoryMsg) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	id, err := k.NewStory(ctx, msg.Body, msg.CategoryID, msg.Creator, msg.StoryType)
	if err != nil {
		return err.Result()
	}

	return app.Result(id)
}
