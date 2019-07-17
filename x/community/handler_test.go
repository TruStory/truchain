package community

import (
	"encoding/json"
	"testing"

	"github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestHandleMsgNewCommunity(t *testing.T) {
	ctx, keeper := mockDB()
	handler := NewHandler(keeper)
	assert.NotNil(t, handler) // assert handler is present

	name, id, description := getFakeCommunityParams()
	creator := keeper.GetParams(ctx).CommunityAdmins[0]
	msg := NewMsgNewCommunity(id, name, description, creator)
	assert.NotNil(t, msg) // assert msgs can be created

	result := handler(ctx, msg)
	idresult := new(types.IDStringResult)
	err := json.Unmarshal(result.Data, &idresult)
	assert.NoError(t, err)

	// TODO: if same community is created twice, it should actually throw an error
	result2 := handler(ctx, msg)
	idresult2 := new(types.IDStringResult)
	err = json.Unmarshal(result2.Data, &idresult2)
	assert.NoError(t, err)
}

func TestByzantineMsg(t *testing.T) {
	ctx, keeper := mockDB()

	handler := NewHandler(keeper)
	assert.NotNil(t, handler)

	res := handler(ctx, nil)
	assert.Equal(t, sdk.CodeUnknownRequest, res.Code)
	assert.Equal(t, sdk.CodespaceRoot, res.Codespace)
}
