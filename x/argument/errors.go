package argument

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Story errors reserve 1500 ~ 1599.
const (
	DefaultCodespace sdk.CodespaceType = "argument"

	CodeInvalidArgumentID sdk.CodeType = 1501
)

// ErrInvalidArgumentID throws an error on invalid title
func ErrInvalidArgumentID() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInvalidArgumentID, "Invalid argument ID")
}
