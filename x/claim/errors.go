package claim

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CodeType is the error code type for the module
type CodeType = sdk.CodeType

// Claim error types
const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	ErrorCodeInvalidBodyTooShort         CodeType = 101
	ErrorCodeInvalidBodyTooLong          CodeType = 102
	ErrorCodeInvalidID                   CodeType = 103
	ErrorCodeNotFound                    CodeType = 104
	ErrorCodeInvalidSType                CodeType = 105
	ErrorCodeClaimsWithCommunityNotFound CodeType = 106
	ErrorCodeInvalidSourceURL            CodeType = 107
	ErrorCodeCreatorJailed               CodeType = 108
	ErrorCodeJSONParsing                 CodeType = 109
)

// ErrInvalidBodyTooShort throws an error on invalid claim body
func ErrInvalidBodyTooShort(body string) sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		ErrorCodeInvalidBodyTooShort,
		"Invalid claim body, too short: "+body)
}

// ErrInvalidBodyTooLong throws an error on invalid claim body
func ErrInvalidBodyTooLong() sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		ErrorCodeInvalidBodyTooLong,
		"Invalid claim body, too long")
}

// ErrInvalidID throws an error on invalid claim body
func ErrInvalidID() sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		ErrorCodeInvalidID,
		"Invalid claim ID")
}

// ErrInvalidCommunityID throws an error on invalid community id
func ErrInvalidCommunityID(id uint64) sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		ErrorCodeClaimsWithCommunityNotFound,
		fmt.Sprintf("Invalid community id: %d", id))
}

// ErrInvalidSourceURL throws an error when a URL in invalid
func ErrInvalidSourceURL(url string) sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		ErrorCodeInvalidSourceURL,
		"Invalid source URL: "+url)
}

// ErrCreatorJailed throws an error on jailed creator
func ErrCreatorJailed(addr sdk.AccAddress) sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		ErrorCodeCreatorJailed,
		"Creator cannot be jailed: "+addr.String())
}

// ErrJSONParse throws an error on failed JSON parsing
func ErrJSONParse(err error) sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		ErrorCodeJSONParsing,
		"JSON parsing error: "+err.Error())
}
