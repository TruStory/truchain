package claim

import (
	"encoding/json"
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

func TestHandleMsgAddAdmin(t *testing.T) {
	ctx, keeper := mockDB()
	handler := NewHandler(keeper)
	assert.NotNil(t, handler) // assert handler is present

	admin := getFakeAdmin()
	creator := keeper.GetParams(ctx).ClaimAdmins[0]
	msg := NewMsgAddAdmin(admin, creator)
	assert.NotNil(t, msg) // assert msgs can be created

	result := handler(ctx, msg)
	var success bool
	err := json.Unmarshal(result.Data, &success)
	assert.NoError(t, err)
	assert.Equal(t, success, true)
}

func TestHandleMsgRemoveAdmin(t *testing.T) {
	ctx, keeper := mockDB()
	handler := NewHandler(keeper)
	assert.NotNil(t, handler) // assert handler is present

	admin := getFakeAdmin()
	creator := keeper.GetParams(ctx).ClaimAdmins[0]
	msg := NewMsgRemoveAdmin(admin, creator)
	assert.NotNil(t, msg) // assert msgs can be created

	result := handler(ctx, msg)
	var success bool
	err := json.Unmarshal(result.Data, &success)
	assert.NoError(t, err)
	assert.Equal(t, success, true)
}

func TestMsgEditClaim(t *testing.T) {
	ctx, keeper := mockDB()

	handler := NewHandler(keeper)
	assert.NotNil(t, handler)

	claim := createFakeClaim(ctx, keeper)
	updatedBody := "If change is the only constant, why immutability is the future of technology?"
	editor := keeper.GetParams(ctx).ClaimAdmins[0]

	msg := NewMsgEditClaim(claim.ID, updatedBody, editor)
	assert.NotNil(t, msg)

	res := handler(ctx, msg)
	assert.NotNil(t, res)
	assert.True(t, res.IsOK())

	var updated Claim
	ModuleCodec.UnmarshalJSON(res.Data, &updated)
	assert.Equal(t, updated.ID, claim.ID)
	assert.Equal(t, updated.Body, updatedBody)
}
