package voting

import (
	"reflect"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Vote errors reserve 1100 ~ 1199.
const (
	DefaultCodespace sdk.CodespaceType = "voting"

	CodeUnknownVote        sdk.CodeType = 1301
	CodeRewardPoolNotEmpty sdk.CodeType = 1302
)

// ErrInvalidVote returns an unknown Vote type error
func ErrInvalidVote(vote interface{}, msg ...string) sdk.Error {
	if mType := reflect.TypeOf(vote); mType != nil {
		errMsg := "Unrecognized Vote type: " + mType.Name() + strings.Join(msg, ",")
		return sdk.NewError(DefaultCodespace, CodeUnknownVote, errMsg)
	}

	return sdk.NewError(DefaultCodespace, CodeUnknownVote, "Unknown Vote type")
}

// ErrNonEmptyRewardPool throws when a reward pool is not empty
func ErrNonEmptyRewardPool(amount sdk.Coin) sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		CodeRewardPoolNotEmpty,
		"Reward pool not empty: "+amount.String())
}
