package story

import (
	"net/url"
	"time"

	c "github.com/TruStory/truchain/x/category"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	amino "github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/encoding/amino"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

func mockDB() (sdk.Context, Keeper, c.Keeper) {
	db := dbm.NewMemDB()

	storyKey := sdk.NewKVStoreKey("stories")
	catKey := sdk.NewKVStoreKey("categories")
	challengeKey := sdk.NewKVStoreKey("challenges")

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(storyKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(catKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(challengeKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())

	codec := amino.NewCodec()
	cryptoAmino.RegisterAmino(codec)
	RegisterAmino(codec)

	ck := c.NewKeeper(catKey, codec)
	sk := NewKeeper(storyKey, ck, codec)

	return ctx, sk, ck
}

func createFakeStory(ctx sdk.Context, sk Keeper, ck c.WriteKeeper) int64 {
	body := "Body of story."

	ctx = ctx.WithBlockHeader(abci.Header{Time: time.Now().UTC()})

	cat := createFakeCategory(ctx, ck)
	creator := sdk.AccAddress([]byte{1, 2})
	storyType := Default
	source := url.URL{}

	storyID, _ := sk.NewStory(ctx, body, cat.ID, creator, source, storyType)

	return storyID
}

func createFakeCategory(ctx sdk.Context, ck c.WriteKeeper) c.Category {
	existing, err := ck.GetCategory(ctx, 1)
	if err == nil {
		return existing
	}
	id, _ := ck.NewCategory(ctx, "decentralized exchanges", sdk.AccAddress([]byte{1, 2}), "trudex", "category for experts in decentralized exchanges")
	cat, _ := ck.GetCategory(ctx, id)
	return cat
}
