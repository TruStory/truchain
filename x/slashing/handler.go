package slashing

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler creates a new handler for slashing module
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		errMsg := fmt.Sprintf("Unrecognized slashing message type: %T", msg)
		return sdk.ErrUnknownRequest(errMsg).Result()
	}
}