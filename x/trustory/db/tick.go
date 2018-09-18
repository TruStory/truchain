package db

import (
	ts "github.com/TruStory/trucoin/x/trustory/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

const maxNumVotes = 10

// NewResponseEndBlock checks stories and generates a ResponseEndBlock.
// It is called at the end of every block, and processes any timing-related
// acitivities within the app -- currently handling the end of the voting
// period.
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
	story, err := k.ActiveStoryQueueHead(ctx)
	if err != nil {
		if err.Code() == ts.CodeActiveStoryQueueEmpty {
			return nil
		}
		return err
	}

	// story reached the end of the voting period
	if ctx.BlockHeader().Time.After(story.VoteEnd) {
		k.ActiveStoryQueuePop(ctx)

		story.Round++

		// check if we have enough votes to proceed
		votes := k.GetActiveVotes(ctx, story.ID)

		// didn't achieve max number of votes
		// mark story as unverifiable and return coins
		if len(votes) < maxNumVotes {
			story.State = ts.Unverifiable
			k.UpdateStory(ctx, story)

			err = returnCoins(ctx, k, story.Escrow, votes)
			if err != nil {
				return err
			}

			// process next story
			return checkStory(ctx, k)
		}

		// reset active votes list
		defer k.SetActiveVotes(ctx, story.ID, []int64{})

		// tally and get votes
		yesVotes := []ts.Vote{}
		noVotes := []ts.Vote{}
		for _, voteID := range votes {
			vote, err := k.GetVote(ctx, voteID)
			if err != nil {
				return err
			}
			if vote.Vote == true {
				yesVotes = append(yesVotes, vote)
			} else {
				noVotes = append(noVotes, vote)
			}
		}

		// determine if we have a supermajority win
		superMajority := 0.66 * maxNumVotes
		if float64(len(yesVotes)) > superMajority || float64(len(noVotes)) > superMajority {
			story.State = ts.Validated
		} else {
			story.State = ts.Unverifiable
		}

		// reward winning voters
		if story.State == ts.Validated {
			if len(yesVotes) > len(noVotes) {
				err := rewardWinners(ctx, k, story.Escrow, story.Category, yesVotes)
				if err != nil {
					return err
				}
			} else {
				err := rewardWinners(ctx, k, story.Escrow, story.Category, noVotes)
				if err != nil {
					return err
				}
			}
		}

		// update story with changes, persist in keeper
		k.UpdateStory(ctx, story)

		// process next in queue
		return checkStory(ctx, k)
	}

	return nil
}

// returnCoins returns coins back to voters for unverified stories
func returnCoins(ctx sdk.Context, k TruKeeper, escrow sdk.AccAddress, voteIDs []int64) sdk.Error {
	for _, voteID := range voteIDs {
		vote, err := k.GetVote(ctx, voteID)
		if err != nil {
			return err
		}

		// return coins back to voter
		_, err = k.ck.SendCoins(ctx, escrow, vote.Creator, vote.Amount)
		if err != nil {
			return err
		}
	}

	return nil
}

// rewardWinners rewards winners of the voting process. It calculates the winning
// amount and distributes coins from the escrow account evenly to all the winners.
func rewardWinners(
	ctx sdk.Context,
	k TruKeeper,
	escrowAddr sdk.AccAddress,
	category ts.StoryCategory,
	winners []ts.Vote) sdk.Error {

	// retrieve coins from escrow account
	// coin denom is category slug, i.e: "stablecoins"
	denom := category.Slug()
	escrowAmount := k.ck.GetCoins(ctx, escrowAddr).AmountOf(denom)

	// calculate winning amount
	numWinners := int64(len(winners))
	winnerAmount := escrowAmount.Div(sdk.NewInt(numWinners))
	amt := sdk.NewCoin(denom, winnerAmount)

	// reward winners
	for _, vote := range winners {
		_, err := k.ck.SendCoins(ctx, escrowAddr, vote.Creator, sdk.Coins{amt})
		if err != nil {
			return err
		}
	}
	return nil
}
