package backing

import (
	"encoding/json"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler creates a new handler for all TruStory messages
func NewHandler(k Keeper) sdk.Handler {
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

func handleBackStoryMsg(ctx sdk.Context, k Keeper, msg BackStoryMsg) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	id, err := k.Create(
		ctx,
		msg.StoryID,
		msg.Amount,
		0,
		msg.Argument,
		msg.Creator)
	if err != nil {
		return err.Result()
	}

	story, err := k.storyKeeper.Story(ctx, msg.StoryID)
	if err != nil {
		return err.Result()
	}

	return result(ctx, k, msg.StoryID, id, msg.Creator, story.Creator, msg.Amount)
}

func handleLikeArgumentMsg(ctx sdk.Context, k Keeper, msg LikeBackingArgumentMsg) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	backingID, err := k.LikeArgument(ctx, msg.ArgumentID, msg.Creator, msg.Amount)
	if err != nil {
		return err.Result()
	}

	backing, err := k.Backing(ctx, backingID)
	if err != nil {
		err.Result()
	}
	argument, err := k.argumentKeeper.Argument(ctx, msg.ArgumentID)
	if err != nil {
		return err.Result()
	}
	return result(ctx, k, backing.StoryID(), backingID, msg.Creator, argument.Creator, msg.Amount)
}

func result(ctx sdk.Context, k Keeper, storyID, backingID int64, from, to sdk.AccAddress, amount sdk.Coin) sdk.Result {
	resultData := app.StakeNotificationResult{
		MsgResult: app.MsgResult{ID: backingID},
		Amount:    amount,
		StoryID:   storyID,
		From:      from,
		To:        to,
	}

	resultBytes, jsonErr := json.Marshal(resultData)
	if jsonErr != nil {
		panic(jsonErr)
	}

	return sdk.Result{
		Data: resultBytes,
		Tags: app.PushTag,
	}
}
