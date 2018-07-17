package trustory

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestNewSubmitStoryMsg(t *testing.T) {
	goodBody := "Jae Kwon invented Tendermint"
	addr1 := sdk.Address([]byte{1, 2})

	cases := []struct {
		valid bool
		ssMsg SubmitStoryMsg
	}{
		{true, NewSubmitStoryMsg(goodBody, addr1)},
	}

	for i, msg := range cases {
		err := msg.ssMsg.ValidateBasic()
		if msg.valid {
			assert.Nil(t, err, "%d: %+v", i, err)
		} else {
			assert.NotNil(r, err, "%d", i)
		}
	}
}
