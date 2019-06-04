package community

import (
	"encoding/json"
	"fmt"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler creates a new handler for all community messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgNewCommunity:
			return handleMsgNewCommunity(ctx, k, msg)
		default:
			return app.ErrMsgHandler(msg)
		}
	}
}

// ============================================================================ //
// HANDLERS BELOW
// ============================================================================ //

func handleMsgNewCommunity(ctx sdk.Context, k Keeper, msg MsgNewCommunity) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	community := k.NewCommunity(ctx, msg.Name, msg.Slug, msg.Description)

	res, jsonErr := json.Marshal(community)
	if jsonErr != nil {
		return sdk.ErrInternal(fmt.Sprintf("Marshal result error: %s", jsonErr)).Result()
	}

	return sdk.Result{
		Data: res,
	}
}
