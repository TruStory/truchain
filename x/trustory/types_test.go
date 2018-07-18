package trustory

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestNewSubmitStoryMsg(t *testing.T) {
	goodBody := "Jae Kwon invented Tendermint"
	addr1 := sdk.AccAddress([]byte{1, 2})
	emptyStr := ""
	emptyAddr := sdk.AccAddress{}

	cases := []struct {
		valid bool
		ssMsg SubmitStoryMsg
	}{
		{true, NewSubmitStoryMsg(goodBody, addr1)},
		{false, NewSubmitStoryMsg(emptyStr, addr1)},
		{false, NewSubmitStoryMsg(goodBody, emptyAddr)},
	}

	for i, msg := range cases {
		err := msg.ssMsg.ValidateBasic()
		if msg.valid {
			assert.Nil(t, err, "%d: %+v", i, err)
		} else {
			assert.NotNil(t, err, "%d", i)
		}
	}
}

func TestNewVoteMsg(t *testing.T) {
	addr1 := sdk.AccAddress([]byte{1, 2})
	emptyStr := ""
	emptyAddr := sdk.AccAddress{}
	yay := "Yes"
	nay := "No"

	var posStoryID int64 = 3
	var negStoryID int64 = -8

	cases := []struct {
		valid   bool
		voteMsg VoteMsg
	}{
		{true, NewVoteMsg(posStoryID, yay, addr1)},
		{true, NewVoteMsg(posStoryID, nay, addr1)},

		{false, NewVoteMsg(negStoryID, yay, addr1)},
		{false, NewVoteMsg(posStoryID, emptyStr, addr1)},
		{false, NewVoteMsg(posStoryID, yay, emptyAddr)},
	}

	for i, msg := range cases {
		err := msg.voteMsg.ValidateBasic()
		if msg.valid {
			assert.Nil(t, err, "%d: %+v", i, err)
			// GetSigners
			assert.Len(t, msg.voteMsg.GetSigners(), 1)
			assert.Equal(t, msg.voteMsg.GetSigners()[0], msg.voteMsg.Voter)
			// GetSignBytes
			assert.NotPanics(t, assert.PanicTestFunc(func() {
				msg.voteMsg.GetSignBytes()
			}))
			assert.NotNil(t, msg.voteMsg.GetSignBytes())
		} else {
			fmt.Print(err)
			assert.NotNil(t, err, "%d", i)
		}
	}
}
