package staking

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Error string

func (e Error) Error() string { return string(e) }

// Staking errors reserve 500 ~ 599.
const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	ErrorCodeInvalidStakeType                sdk.CodeType = 501
	ErrorCodeAccountJailed                   sdk.CodeType = 502
	ErrorCodeInvalidBodyLength               sdk.CodeType = 503
	ErrorCodeInvalidSummaryLength            sdk.CodeType = 504
	ErrorCodeUnknownArgument                 sdk.CodeType = 505
	ErrorCodeUnknownStake                    sdk.CodeType = 506
	ErrorCodeDuplicateStake                  sdk.CodeType = 507
	ErrorCodeMaxNumOfArgumentsReached        sdk.CodeType = 508
	ErrorCodeMaxAmountStakingReached         sdk.CodeType = 509
	ErrorCodeInvalidQueryParams              sdk.CodeType = 510
	ErrorCodeJSONParsing                     sdk.CodeType = 511
	ErrorCodeUnknownClaim                    sdk.CodeType = 512
	ErrorCodeUnknownStakeType                sdk.CodeType = 513
	ErrorCodeCannotEditArgumentAlreadyStaked sdk.CodeType = 514
	ErrorCodeCannotEditArgumentWrongCreator  sdk.CodeType = 515
	ErrorCodeMinBalance                      sdk.CodeType = 516
	ErrorCodeAddressNotAuthorised            sdk.CodeType = 517
)

// GenesisErrors
const (
	ErrInvalidArgumentStakeDenom = Error("invalid denomination for argument stake")
	ErrInvalidUpvoteStakeDenom   = Error("invalid denomination for upvote stake")
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

// ErrCodeUnknownClaim throws an error when an invalid claim id
func ErrCodeUnknownClaim(claimID uint64) sdk.Error {
	return sdk.NewError(DefaultCodespace,
		ErrorCodeUnknownClaim,
		fmt.Sprintf("Unknown claim id %d", claimID),
	)
}

// ErrCodeUnknownStake throws an error when an invalid stake id
func ErrCodeUnknownStake(stakeID uint64) sdk.Error {
	return sdk.NewError(DefaultCodespace,
		ErrorCodeUnknownStake,
		fmt.Sprintf("Unknown stake id %d", stakeID),
	)
}

// ErrCodeUnknownStakeType throws an error when an invalid stake type
func ErrCodeUnknownStakeType() sdk.Error {
	return sdk.NewError(DefaultCodespace,
		ErrorCodeUnknownStakeType,
		fmt.Sprintf("Unknown stake type"),
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

// ErrCodeCannotEditArgumentAlreadyStaked throws an error when an argument cannot be edited because it has already been staked
func ErrCodeCannotEditArgumentAlreadyStaked(argumentID uint64) sdk.Error {
	return sdk.NewError(DefaultCodespace,
		ErrorCodeCannotEditArgumentAlreadyStaked,
		fmt.Sprintf("This argument cannot be edited because someone else has already agreed to it"),
	)
}

// ErrCodeCannotEditArgumentWrongCreator throws an error when an argument cannot be edited because the edit is not coming from the creator
func ErrCodeCannotEditArgumentWrongCreator(argumentID uint64) sdk.Error {
	return sdk.NewError(DefaultCodespace,
		ErrorCodeCannotEditArgumentWrongCreator,
		fmt.Sprintf("This argument cannot be edited because you are not the writer of the Argument"),
	)
}

// ErrCodeMaxAmountStakingReached throws an error when you already staked.
func ErrCodeMaxAmountStakingReached() sdk.Error {
	return sdk.NewError(DefaultCodespace,
		ErrorCodeMaxAmountStakingReached,
		fmt.Sprintf("You have reached the maximum amount for staking"),
	)
}

// ErrCodeMinBalance throws an error when you have the minimum balance.
func ErrCodeMinBalance() sdk.Error {
	return sdk.NewError(DefaultCodespace,
		ErrorCodeMinBalance,
		fmt.Sprintf("You don't have the minimum balance required for staking"),
	)
}

// ErrInvalidQueryParams throws an error when the transaction type is invalid.
func ErrInvalidQueryParams(err error) sdk.Error {
	return sdk.NewError(DefaultCodespace,
		ErrorCodeInvalidQueryParams,
		fmt.Sprintf("Invalid query params  %s", err.Error()),
	)
}

// ErrJSONParse throws an error on failed JSON parsing
func ErrJSONParse(err error) sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		ErrorCodeJSONParsing,
		"JSON parsing error: "+err.Error())
}

// ErrAddressNotAuthorised throws an error when the address is not admin
func ErrAddressNotAuthorised() sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		ErrorCodeAddressNotAuthorised,
		"This creator is not authorised to perform this action.")
}
