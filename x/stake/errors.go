package stake

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Stake errors reserve 1600 ~ 1699.
const (
	DefaultCodespace sdk.CodespaceType = StoreKey

	CodeInvalidStoryState sdk.CodeType = 1601
)

// ErrInvalidStoryState throws when story not pending
func ErrInvalidStoryState(state string) sdk.Error {
	msg := "Story can only be staked when it's pending. Current state is: %s"

	return sdk.NewError(
		DefaultCodespace, CodeInvalidStoryState, fmt.Sprintf(msg, state))
}
