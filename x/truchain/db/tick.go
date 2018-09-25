package db

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewResponseEndBlock checks stories and generates a ResponseEndBlock.
// It is called at the end of every block, and processes any timing-related
// acitivities within the app.
func (k TruKeeper) NewResponseEndBlock(ctx sdk.Context) abci.ResponseEndBlock {
	err := checkStory(ctx, k)
	if err != nil {
		panic(err)
	}

	return abci.ResponseEndBlock{}
}

// ============================================================================

// checkStory checks if the story reached the end of the voting period
// and handles distributing rewards. It calls itself recursively until
// all stories in the in-progress state are processed, or until there
// are no more stories to process.
func checkStory(ctx sdk.Context, k TruKeeper) sdk.Error {
	panic("test")
}
