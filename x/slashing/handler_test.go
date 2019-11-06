package slashing

import (
	"encoding/json"
	"net/url"
	"testing"

	"github.com/TruStory/truchain/x/staking"
	"github.com/stretchr/testify/assert"
)

func TestHandle_SlashArgument(t *testing.T) {
	ctx, k := mockDB()
	handler := NewHandler(k)

	staker := k.GetParams(ctx).SlashAdmins[1]
	body := "Blockchains have the power to fund grassroots communities to solve specific problems."
	communityID := "crypto"
	claim, err := k.claimKeeper.SubmitClaim(ctx, body, communityID, staker, url.URL{})
	assert.NoError(t, err)
	arg, err := k.stakingKeeper.SubmitArgument(ctx, "arg1", "summary1", staker, claim.ID, staking.StakeChallenge)
	assert.NoError(t, err)

	slashDetailedReason := "adsfadsf"
	msg := NewMsgSlashArgument(arg.ID, SlashTypeUnhelpful, SlashReasonFocusedOnPerson, slashDetailedReason, staker)
	res := handler(ctx, msg)

	assert.True(t, res.IsOK())
}

func TestHandleMsgAddAdmin(t *testing.T) {
	ctx, keeper := mockDB()
	handler := NewHandler(keeper)
	assert.NotNil(t, handler) // assert handler is present

	_, _, admin, _ := getFakeAppAccountParams()
	creator := keeper.GetParams(ctx).SlashAdmins[0]
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

	_, _, admin, _ := getFakeAppAccountParams()
	creator := keeper.GetParams(ctx).SlashAdmins[0]
	msg := NewMsgRemoveAdmin(admin, creator)
	assert.NotNil(t, msg) // assert msgs can be created

	result := handler(ctx, msg)
	var success bool
	err := json.Unmarshal(result.Data, &success)
	assert.NoError(t, err)
	assert.Equal(t, success, true)
}
