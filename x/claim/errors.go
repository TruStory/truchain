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
	ErrorCodeAddressNotAuthorised        CodeType = 109
	ErrorCodeJSONParsing                 CodeType = 110
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

// ErrUnknownClaim throws an error on an unknown claim id
func ErrUnknownClaim(id uint64) sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		ErrorCodeInvalidID,
		fmt.Sprintf("Unknown claim id: %d", id))
}

// ErrInvalidCommunityID throws an error on invalid community id
func ErrInvalidCommunityID(id string) sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		ErrorCodeClaimsWithCommunityNotFound,
		fmt.Sprintf("Invalid community id: %s", id))
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

// ErrAddressNotAuthorised throws an error when the address is not admin
func ErrAddressNotAuthorised() sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		ErrorCodeAddressNotAuthorised,
		"This address is not authorised to perform this action.")
}

// ErrJSONParse throws an error on failed JSON parsing
func ErrJSONParse(err error) sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		ErrorCodeJSONParsing,
		"JSON parsing error: "+err.Error())
}
