package backing

import (
	"fmt"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler creates a new handler for all TruStory messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case BackStoryMsg:
			return handleBackStoryMsg(ctx, k, msg)
		default:
			return app.ErrMsgHandler(msg)
		}
	}
}

// ============================================================================

func handleBackStoryMsg(ctx sdk.Context, k Keeper, msg BackStoryMsg) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	id, err := k.Create(
		ctx,
		msg.StoryID,
		msg.Amount,
		msg.Argument,
		msg.Creator,
		false)
	if err != nil {
		return err.Result()
	}

	pushBody := fmt.Sprintf("%s backed your story for %s", msg.Creator, msg.Amount)

	story, err := k.storyKeeper.Story(ctx, msg.StoryID)
	if err != nil {
		return err.Result()
	}

	tags := sdk.NewTags(
		"push.type", []byte("normal"),
		"push.body", []byte(pushBody),
		"push.from", msg.Creator.Bytes(),
		"push.to", story.Creator.Bytes(),
	)

	return sdk.Result{
		Data: k.GetCodec().MustMarshalBinaryLengthPrefixed(id),
		Tags: tags,
	}
}
