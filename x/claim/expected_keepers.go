package claim

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccountKeeper is the expected account keeper interface for this module
type AccountKeeper interface {
	IsJailed(ctx sdk.Context, addr sdk.AccAddress) (bool, sdk.Error)
}
