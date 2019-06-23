package claim

import (
	"fmt"
	"net/url"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler creates a new handler
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgCreateClaim:
			return handleMsgCreateClaim(ctx, keeper, msg)
		case MsgDeleteClaim:
			return handleMsgDeleteClaim(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized claim message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgCreateClaim(ctx sdk.Context, keeper Keeper, msg MsgCreateClaim) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	// parse url from string
	sourceURL, urlError := url.Parse(msg.Source)
	if urlError != nil {
		return ErrInvalidSourceURL(msg.Source).Result()
	}

	claim, err := keeper.SubmitClaim(ctx, msg.Body, msg.CommunityID, msg.Creator, *sourceURL)
	if err != nil {
		return err.Result()
	}

	res, codecErr := ModuleCodec.MarshalBinaryBare(claim)
	if codecErr != nil {
		return sdk.ErrInternal(fmt.Sprintf("Marshal result error: %s", codecErr)).Result()
	}

	return sdk.Result{
		Data: res,
	}
}

func handleMsgDeleteClaim(ctx sdk.Context, keeper Keeper, msg MsgDeleteClaim) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	// _ = keeper.Delete(ctx, msg.ID)

	return sdk.Result{}
}
