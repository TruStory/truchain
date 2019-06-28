package claim

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestMsgCreateClaim(t *testing.T) {
	ctx, keeper := mockDB()

	handler := NewHandler(keeper)
	assert.NotNil(t, handler)

	communityID := "crypto"
	body := "fake story body with minimum length"
	creator := sdk.AccAddress([]byte{1, 2})
	source := "http://trustory.io"
	msg := NewMsgCreateClaim(communityID, body, creator, source)
	assert.NotNil(t, msg)

	res := handler(ctx, msg)
	assert.NotNil(t, res)
	assert.True(t, res.IsOK())

	var claim Claim
	ModuleCodec.UnmarshalJSON(res.Data, &claim)
	assert.Equal(t, uint64(1), claim.ID)
	assert.Equal(t, body, claim.Body)
}
