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

	name, slug, description := getFakeCommunityParams()
	creator := sdk.AccAddress([]byte{1, 2})
	msg := NewMsgNewCommunity(name, slug, description, creator)
	assert.NotNil(t, msg) // assert msgs can be created

	result := handler(ctx, msg)
	idresult := new(types.IDResult)
	err := json.Unmarshal(result.Data, &idresult)
	assert.NoError(t, err)

	// TODO: if same community is created twice, it should actually throw an error
	result2 := handler(ctx, msg)
	idresult2 := new(types.IDResult)
	err = json.Unmarshal(result2.Data, &idresult2)
	assert.NoError(t, err)

	assert.Equal(t, int64(1), idresult.ID, "incorrect result data")
	assert.Equal(t, int64(2), idresult2.ID, "incorrect result data")
}

func TestByzantineMsg(t *testing.T) {
	ctx, keeper := mockDB()

	handler := NewHandler(keeper)
	assert.NotNil(t, handler)

	res := handler(ctx, nil)
	assert.Equal(t, sdk.CodeUnknownRequest, res.Code)
	assert.Equal(t, sdk.CodespaceRoot, res.Codespace)
}
