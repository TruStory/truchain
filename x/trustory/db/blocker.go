package db

import (
	"fmt"

	ts "github.com/TruStory/trucoin/x/trustory/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

const maxNumVotes = 10

// NewResponseEndBlock checks stories and generates a ResponseEndBlock
func (k TruKeeper) NewResponseEndBlock(ctx sdk.Context) abci.ResponseEndBlock {
	err := checkStory(ctx, k)
	if err != nil {
		panic(err)
	}
	return abci.ResponseEndBlock{}
}

// ============================================================================

// checkStory checks if the story reached the end of the voting period
// and handles distributing rewards
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

		fmt.Printf("hello2")

		story.Round++

		// check if we have enough votes to proceed
		votes, err := k.GetActiveVotes(ctx, story.ID)

		fmt.Printf("num votes %d", len(votes))

		if err != nil {
			// process next story
			return checkStory(ctx, k)
		}

		fmt.Printf("hello3")

		// didn't achieve max number of votes
		// mark story as unverifiable and return coins
		if len(votes) < maxNumVotes {
			story.State = ts.Unverifiable
			err := k.UpdateStory(ctx, story)
			if err != nil {
				return err
			}

			fmt.Printf("hello4")

			err = returnCoins(ctx, k, story.Escrow, votes)
			if err != nil {
				return err
			}

			fmt.Printf("hello5")

			// process next story
			return checkStory(ctx, k)
		}

		fmt.Printf("hello6")

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

		superMajority := 0.66 * maxNumVotes
		if float64(len(yesVotes)) > superMajority || float64(len(noVotes)) > superMajority {
			story.State = ts.Validated
		} else {
			story.State = ts.Unverifiable
		}

		// reward winning voters
		if story.State == ts.Validated {
			if len(yesVotes) > len(noVotes) {
				err := rewardWinners(ctx, k, story.Escrow, story.Category, yesVotes, noVotes)
				if err != nil {
					return err
				}
			} else {
				err := rewardWinners(ctx, k, story.Escrow, story.Category, noVotes, yesVotes)
				if err != nil {
					return err
				}
			}
		}

		// update story with changes, persist in keeper
		err = k.UpdateStory(ctx, story)
		if err != nil {
			return err
		}

		// process next in queue
		return checkStory(ctx, k)
	}

	// TODO: shouldn't return nil, since it will panic
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
		fmt.Print(vote.Creator)
		fmt.Print(vote.Amount)
		_, err = k.ck.SendCoins(ctx, escrow, vote.Creator, vote.Amount)
		if err != nil {
			return err
		}
	}

	return nil
}

// rewardWinners rewards winners of the voting process
func rewardWinners(
	ctx sdk.Context,
	k TruKeeper,
	escrowAddr sdk.AccAddress,
	category ts.StoryCategory,
	win []ts.Vote,
	lose []ts.Vote) sdk.Error {

	for _, vote := range lose {
		_, err := k.ck.SendCoins(ctx, vote.Creator, escrowAddr, vote.Amount)
		if err != nil {
			return err
		}
	}

	// calculate winning amount
	numWinners := int64(len(win))
	denom := category.Slug() // coin denom is category slug, i.e: "stablecoins"
	escrowAmount := k.ck.GetCoins(ctx, escrowAddr).AmountOf(denom)
	winnerAmount := escrowAmount.Div(sdk.NewInt(numWinners))
	amt := sdk.NewCoin(denom, winnerAmount)

	// reward winners
	for _, vote := range win {
		_, err := k.ck.SendCoins(ctx, escrowAddr, vote.Creator, sdk.Coins{amt})
		if err != nil {
			return err
		}
	}
	return nil
}
