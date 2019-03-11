package types_test

import (
	"fmt"
	"testing"
	"time"

	app "github.com/TruStory/truchain/types"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	amino "github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

func Test_List(t *testing.T) {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)

	header := abci.Header{Time: time.Now().Add(50 * 24 * time.Hour)}

	keeperStoreKey := sdk.NewKVStoreKey("keeper")
	listStoreKey := sdk.NewKVStoreKey("list")

	ms.MountStoreWithDB(keeperStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(listStoreKey, sdk.StoreTypeIAVL, db)
	err := ms.LoadLatestVersion()
	assert.NoError(t, err)

	codec := amino.NewCodec()
	ctx := sdk.NewContext(ms, header, false, log.NewNopLogger())
	address := sdk.AccAddress([]byte("address"))
	keeper := app.NewKeeper(codec, keeperStoreKey)
	list := app.NewUserList(listStoreKey)

	storyID := int64(1)
	list.Append(ctx, keeper, storyID, address, 5)

	// test directly from the store
	val := keeper.GetStore(ctx).Get(getKey("list", "keeper", address.String(), 1))
	var result int64
	codec.MustUnmarshalBinaryBare(val, &result)
	assert.Equal(t, int64(5), result)

	// test Get
	assert.Equal(t, int64(5), list.Get(ctx, keeper, storyID, address))

	list.Delete(ctx, keeper, storyID, address)

	// test directly from the store
	val = keeper.GetStore(ctx).Get(getKey("list", "keeper", address.String(), 1))
	assert.Nil(t, val, "should return nil bytes")
	// test Get
	assert.Equal(t, int64(0), list.Get(ctx, keeper, storyID, address))

}

// "[foreignStoreKey]:id:[keyID]:[storeKey]:users:[user]"
func getKey(foreignStoreKey, storeKey, address string, keyID int64) []byte {
	key := fmt.Sprintf(
		"%s:id:%d:%s:user:%s",
		foreignStoreKey,
		keyID,
		storeKey,
		address)
	return []byte(key)
}
