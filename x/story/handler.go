package story

import (
	"net/url"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler creates a new handler for all TruStory messages
func NewHandler(k WriteKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case SubmitStoryMsg:
			return handleSubmitStoryMsg(ctx, k, msg)
		case FlagStoryMsg:
			return handleFlagStoryMsg(ctx, k, msg)
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

	// parse url from string
	sourceURL, urlError := url.ParseRequestURI(msg.Source)
	if urlError != nil {
		return ErrInvalidSourceURL(msg.Source).Result()
	}

	id, err := k.Create(
		ctx,
		msg.Argument,
		msg.Body,
		msg.CategoryID,
		msg.Creator,
		[]Evidence{},
		*sourceURL,
		msg.StoryType)
	if err != nil {
		return err.Result()
	}

	return app.Result(id)
}

func handleFlagStoryMsg(ctx sdk.Context, k WriteKeeper, msg FlagStoryMsg) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	// get story
	story, err := k.Story(ctx, msg.StoryID)
	if err != nil {
		err.Result()
	}

	if story.Flagged != true {
		story.Flagged = true
		k.UpdateStory(ctx, story)
	}

	return app.Result(story.ID)
}
