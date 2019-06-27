package staking

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Staking errors reserve 500 ~ 599.
const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	ErrorCodeInvalidStakeType         sdk.CodeType = 501
	ErrorCodeAccountJailed            sdk.CodeType = 502
	ErrorCodeInvalidBodyLength        sdk.CodeType = 503
	ErrorCodeInvalidSummaryLength     sdk.CodeType = 504
	ErrorCodeUnknownArgument          sdk.CodeType = 505
	ErrorCodeUnknownStake             sdk.CodeType = 506
	ErrorCodeDuplicateStake           sdk.CodeType = 507
	ErrorCodeMaxNumOfArgumentsReached sdk.CodeType = 508
	ErrorCodeMaxAmountStakingReached  sdk.CodeType = 509
)

// ErrCodeAccountJailed throws an error is in jailed status when performing actions.
func ErrCodeAccountJailed(acc sdk.AccAddress) sdk.Error {
	return sdk.NewError(DefaultCodespace,
		ErrorCodeAccountJailed,
		fmt.Sprintf("Account is jailed %s", acc.String()),
	)
}

// ErrCodeInvalidStakeType throws an error when an invalid stake type is
func ErrCodeInvalidStakeType(stakeType StakeType) sdk.Error {
	return sdk.NewError(DefaultCodespace,
		ErrorCodeInvalidStakeType,
		fmt.Sprintf("Invalid stake type %s", stakeType.String()),
	)
}

// ErrCodeInvalidBodyLength throws an error when an invalid body length.
func ErrCodeInvalidBodyLength() sdk.Error {
	return sdk.NewError(DefaultCodespace,
		ErrorCodeInvalidBodyLength,
		"Invalid body length ",
	)
}

// ErrCodeInvalidSummaryLength throws an error when an invalid body length.
func ErrCodeInvalidSummaryLength() sdk.Error {
	return sdk.NewError(DefaultCodespace,
		ErrorCodeInvalidSummaryLength,
		"Invalid summary length ",
	)
}

// ErrCodeUnknownArgument throws an error when an invalid argument id
func ErrCodeUnknownArgument(argumentID uint64) sdk.Error {
	return sdk.NewError(DefaultCodespace,
		ErrorCodeUnknownArgument,
		fmt.Sprintf("Unknown argument id %d", argumentID),
	)
}

// ErrCodeUnknownStake throws an error when an invalid stake id
func ErrCodeUnknownStake(argumentID uint64) sdk.Error {
	return sdk.NewError(DefaultCodespace,
		ErrorCodeUnknownStake,
		fmt.Sprintf("Unknown stake id %d", argumentID),
	)
}

// ErrCodeDuplicateStake throws an error when you already staked.
func ErrCodeDuplicateStake(argumentID uint64) sdk.Error {
	return sdk.NewError(DefaultCodespace,
		ErrorCodeDuplicateStake,
		fmt.Sprintf("Already staked for argument id %d", argumentID),
	)
}

// ErrCodeMaxNumOfArgumentsReached throws an error when you already staked.
func ErrCodeMaxNumOfArgumentsReached(max int) sdk.Error {
	return sdk.NewError(DefaultCodespace,
		ErrorCodeMaxNumOfArgumentsReached,
		fmt.Sprintf("You have reached max number of %d arguments per claim", max),
	)
}

// ErrCodeMaxAmountStakingReached throws an error when you already staked.
func ErrCodeMaxAmountStakingReached(days int) sdk.Error {
	return sdk.NewError(DefaultCodespace,
		ErrorCodeMaxAmountStakingReached,
		fmt.Sprintf("You have reached the max amout for staking for a period of %d hours", days),
	)
}
