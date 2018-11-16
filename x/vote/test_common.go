package vote

import (
	"time"

	"github.com/TruStory/truchain/x/backing"

	"github.com/TruStory/truchain/x/challenge"

	"github.com/TruStory/truchain/x/game"

	c "github.com/TruStory/truchain/x/category"
	s "github.com/TruStory/truchain/x/story"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	amino "github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/encoding/amino"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

func mockDB() (
	sdk.Context,
	Keeper,
	s.Keeper,
	c.Keeper,
	challenge.Keeper,
	bank.Keeper,
	backing.Keeper, auth.AccountKeeper) {

	db := dbm.NewMemDB()

	accKey := sdk.NewKVStoreKey("acc")
	storyKey := sdk.NewKVStoreKey("stories")
	catKey := sdk.NewKVStoreKey("categories")
	challengeKey := sdk.NewKVStoreKey("challenges")
	gameKey := sdk.NewKVStoreKey("games")
	gameQueueKey := sdk.NewKVStoreKey("game_queue")
	voteKey := sdk.NewKVStoreKey("vote")
	backingKey := sdk.NewKVStoreKey("backing")

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(accKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(storyKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(catKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(challengeKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(gameKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(gameQueueKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(voteKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(backingKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	header := abci.Header{Time: time.Now().Add(50 * 24 * time.Hour)}
	ctx := sdk.NewContext(ms, header, false, log.NewNopLogger())

	codec := amino.NewCodec()
	cryptoAmino.RegisterAmino(codec)
	codec.RegisterInterface((*auth.Account)(nil), nil)
	codec.RegisterConcrete(&auth.BaseAccount{}, "auth/Account", nil)

	am := auth.NewAccountKeeper(codec, accKey, auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(am)
	ck := c.NewKeeper(catKey, codec)
	sk := s.NewKeeper(storyKey, ck, codec)
	backingKeeper := backing.NewKeeper(backingKey, sk, bankKeeper, ck, codec)
	gameKeeper := game.NewKeeper(gameKey, gameQueueKey, gameQueueKey, sk, bankKeeper, codec)
	challengeKeeper := challenge.NewKeeper(challengeKey, gameQueueKey, bankKeeper, gameKeeper, sk, codec)

	k := NewKeeper(
		voteKey,
		gameQueueKey,
		am,
		backingKeeper,
		challengeKeeper,
		sk,
		gameKeeper,
		bankKeeper,
		codec)

	return ctx, k, sk, ck, challengeKeeper, bankKeeper, backingKeeper, am
}

func createFakeStory(ctx sdk.Context, sk s.Keeper, ck c.WriteKeeper) int64 {
	body := "Body of story."
	cat := createFakeCategory(ctx, ck)
	creator := sdk.AccAddress([]byte{1, 2})
	storyType := s.Default

	storyID, _ := sk.NewStory(ctx, body, cat.ID, creator, storyType)

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
