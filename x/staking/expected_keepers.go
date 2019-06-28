package staking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/TruStory/truchain/x/claim"
)

type AccountKeeper interface {
	IsJailed(ctx sdk.Context, address sdk.AccAddress) (bool, sdk.Error)
	UnJail(ctx sdk.Context, address sdk.AccAddress) sdk.Error
}

type ClaimKeeper interface {
	Claim(ctx sdk.Context, id uint64) (claim claim.Claim, ok bool)
}
