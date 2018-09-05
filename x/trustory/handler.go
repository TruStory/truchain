package trustory

import (
	"encoding/binary"
	"reflect"
	"time"

	db "github.com/TruStory/trucoin/x/trustory/db"
	ts "github.com/TruStory/trucoin/x/trustory/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

const votingPeriod = 24 // hours
const maxNumVotes = 10

// NewHandler creates a new handler for all TruStory messages
func NewHandler(k db.TruKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case ts.SubmitStoryMsg:
			return handleSubmitStoryMsg(ctx, k, msg)
		case ts.VoteMsg:
			return handleVoteMsg(ctx, k, msg)
		default:
			errMsg := "Unrecognized Msg type: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// type EndBlocker func(ctx Context, req abci.RequestEndBlock) abci.ResponseEndBlock

// NewEndBlocker checks stories and generates an EndBlocker
func NewEndBlocker(k db.TruKeeper) sdk.EndBlocker {
	return func(ctx sdk.Context, req abci.RequestEndBlock) (res abci.ResponseEndBlock) {
		err := checkStory(ctx, k)
		if err != nil {
			panic(err)
		}
		return abci.ResponseEndBlock{}
	}
}

// ============================================================================

// checkStory checks if the story reached the end of the voting period
// and handles distributing rewards
func checkStory(ctx sdk.Context, k db.TruKeeper) sdk.Error {
	story, err := k.ActiveStoryQueueHead(ctx)
	if err != nil {
		return err
	}

	// story reached the end of the voting period
	if ctx.BlockHeader().Time.After(story.VoteEnd) {
		k.ActiveStoryQueuePop(ctx)

		story.Round++

		// check if we have enough votes to proceed
		votes, err := k.GetActiveVotes(ctx, story.ID)
		if err != nil {
			// process next story
			return checkStory(ctx, k)
		}

		// didn't achieve max number of votes
		// mark story as unverifiable and return coins
		if len(votes) < maxNumVotes {
			story.State = ts.Unverifiable
			err := k.UpdateStory(ctx, story)
			if err != nil {
				return err
			}
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
	return nil
}

// handleSubmitStoryMsg handles the logic of a SubmitStoryMsg
func handleSubmitStoryMsg(ctx sdk.Context, k db.TruKeeper, msg ts.SubmitStoryMsg) sdk.Result {
	err := msg.ValidateBasic()
	if err != nil {
		return err.Result()
	}

	// calculate voting period
	voteStart := ctx.BlockHeader().Time
	voteEnd := voteStart.Add(time.Hour * time.Duration(votingPeriod))

	storyID, err := k.AddStory(ctx, msg.Body, msg.Category, msg.Creator, msg.StoryType, voteStart, voteEnd)
	if err != nil {
		panic(err)
	}

	return sdk.Result{Data: i2b(storyID)}
}

func handleVoteMsg(ctx sdk.Context, k db.TruKeeper, msg ts.VoteMsg) sdk.Result {
	err := msg.ValidateBasic()
	if err != nil {
		return err.Result()
	}

	voteID, err := k.VoteStory(ctx, msg.StoryID, msg.Creator, msg.Vote, sdk.Coins{msg.Amount})
	if err != nil {
		panic(err)
	}

	return sdk.Result{Data: i2b(voteID)}
}

// i2b converts an int64 into a byte array
func i2b(x int64) []byte {
	var b [binary.MaxVarintLen64]byte
	return b[:binary.PutVarint(b[:], x)]
}

// returnCoins returns coins back to voters for unverified stories
func returnCoins(ctx sdk.Context, k db.TruKeeper, escrow sdk.AccAddress, voteIDs []int64) sdk.Error {
	for _, voteID := range voteIDs {
		vote, err := k.GetVote(ctx, voteID)
		if err != nil {
			return err
		}
		// return coins back to voter
		_, err = k.Ck.SendCoins(ctx, escrow, vote.Creator, vote.Amount)
		if err != nil {
			return err
		}
	}

	return nil
}

// rewardWinners rewards winners of the voting process
func rewardWinners(
	ctx sdk.Context,
	k db.TruKeeper,
	escrowAddr sdk.AccAddress,
	category ts.StoryCategory,
	win []ts.Vote,
	lose []ts.Vote) sdk.Error {

	for _, vote := range lose {
		_, err := k.Ck.SendCoins(ctx, vote.Creator, escrowAddr, vote.Amount)
		if err != nil {
			return err
		}
	}

	// calculate winning amount
	numWinners := int64(len(win))
	escrowAmount := k.Am.GetAccount(ctx, escrowAddr).GetCoins().AmountOf(category.Slug())
	winnerAmount := escrowAmount.Div(sdk.NewInt(numWinners))
	amt := sdk.NewCoin(category.Slug(), winnerAmount)

	// reward winners
	for _, vote := range win {
		_, err := k.Ck.SendCoins(ctx, escrowAddr, vote.Creator, sdk.Coins{amt})
		if err != nil {
			return err
		}
	}
	return nil
}
