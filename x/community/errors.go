package community

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Community errors reserve 800 ~ 899.
const (
	DefaultCodespace sdk.CodespaceType = StoreKey

	CodeCommunityNotFound sdk.CodeType = 801
)

// ErrCommunityNotFound throws an error when the searched category is not found
func ErrCommunityNotFound(id int64) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeCommunityNotFound, "Community not found with ID: "+fmt.Sprintf("%d", id))
}
