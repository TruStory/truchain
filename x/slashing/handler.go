package slashing

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/slashing/tags"
)

// NewHandler creates a new handler for slashing module
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSlashArgument:
			return handleMsgSlashArgument(ctx, keeper, msg)
		case MsgUpdateParams:
			return handleMsgUpdateParams(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized slashing message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgSlashArgument(ctx sdk.Context, k Keeper, msg MsgSlashArgument) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	slash, punishmentResults, err := k.CreateSlash(ctx, msg.ArgumentID, msg.SlashType, msg.SlashReason, msg.SlashDetailedReason, msg.Creator)
	if err != nil {
		return err.Result()
	}

	res, jsonErr := ModuleCodec.MarshalJSON(slash)
	if jsonErr != nil {
		return sdk.ErrInternal(fmt.Sprintf("Marshal result error: %s", jsonErr)).Result()
	}

	argument, ok := k.stakingKeeper.Argument(ctx, slash.ArgumentID)
	if !ok {
		return sdk.ErrInternal("couldn't find argument").Result()
	}
	isJailed, err := k.accountKeeper.IsJailed(ctx, argument.Creator)
	if err != nil {
		return err.Result()
	}
	resultTags := append(app.PushTag,
		sdk.NewTags(
			tags.Category, tags.TxCategory,
			tags.Action, tags.ActionCreateSlash,
			tags.ArgumentCreator, argument.Creator.String(),
		)...,
	)
	if isJailed {
		resultTags = append(resultTags, sdk.NewTags(tags.ArgumentCreatorJailed, "jailed")...)
	}

	if len(punishmentResults) > 0 {
		json, jsonErr := json.Marshal(punishmentResults)
		if jsonErr != nil {
			return sdk.ErrInternal(fmt.Sprintf("Marshal result error: %s", jsonErr)).Result()
		}
		resultTags = append(resultTags, sdk.NewTags(tags.SlashResults, json)...)
	}
	return sdk.Result{
		Data: res,
		Tags: resultTags,
	}
}

func handleMsgUpdateParams(ctx sdk.Context, k Keeper, msg MsgUpdateParams) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	err := k.UpdateParams(ctx, msg.Updates)
	if err != nil {
		return err.Result()
	}

	res, jsonErr := json.Marshal(true)
	if jsonErr != nil {
		return sdk.ErrInternal(fmt.Sprintf("Marshal result error: %s", jsonErr)).Result()
	}

	return sdk.Result{
		Data: res,
	}
}
