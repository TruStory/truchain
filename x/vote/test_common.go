package vote

import (
	"crypto/rand"
	"net/url"
	"time"

	"github.com/TruStory/truchain/x/backing"

	"github.com/TruStory/truchain/x/challenge"

	"github.com/TruStory/truchain/x/game"

	c "github.com/TruStory/truchain/x/category"
	"github.com/TruStory/truchain/x/story"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	amino "github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

func mockDB() (sdk.Context, Keeper, c.Keeper) {

	db := dbm.NewMemDB()

	accKey := sdk.NewKVStoreKey("acc")
	storyKey := sdk.NewKVStoreKey("stories")
	catKey := sdk.NewKVStoreKey("categories")
	challengeKey := sdk.NewKVStoreKey("challenges")
	gameKey := sdk.NewKVStoreKey("games")
	pendingGameQueueKey := sdk.NewKVStoreKey("pendingGameQueue")
	gameQueueKey := sdk.NewKVStoreKey("gameQueue")
	voteKey := sdk.NewKVStoreKey("vote")
	backingKey := sdk.NewKVStoreKey("backing")

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(accKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(storyKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(catKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(challengeKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(gameKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(pendingGameQueueKey, sdk.StoreTypeIAVL, db)
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
	sk := story.NewKeeper(storyKey, ck, codec)
	backingKeeper := backing.NewKeeper(backingKey, sk, bankKeeper, ck, codec)
	gameKeeper := game.NewKeeper(gameKey, gameQueueKey, gameQueueKey, sk, backingKeeper, bankKeeper, codec)
	challengeKeeper := challenge.NewKeeper(challengeKey, pendingGameQueueKey, bankKeeper, gameKeeper, sk, codec)

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

	return ctx, k, ck
}

func createFakeStory(ctx sdk.Context, sk story.WriteKeeper, ck c.WriteKeeper) int64 {
	body := "Body of story."
	cat := createFakeCategory(ctx, ck)
	creator := sdk.AccAddress([]byte{1, 2})
	storyType := story.Default
	source := url.URL{}
	argument := "this is a fake argument"

	storyID, _ := sk.Create(ctx, argument, body, cat.ID, creator, source, storyType)

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

func fakeFundedCreator(ctx sdk.Context, k bank.Keeper) sdk.AccAddress {
	bz := make([]byte, 4)
	rand.Read(bz)
	creator := sdk.AccAddress(bz)

	// give user some category coins
	amount := sdk.NewCoin("trudex", sdk.NewInt(2000))
	k.AddCoins(ctx, creator, sdk.Coins{amount})

	return creator
}

func fakeValidationGame() (ctx sdk.Context, votes poll, k Keeper) {

	ctx, k, ck := mockDB()

	storyID := createFakeStory(ctx, k.storyKeeper, ck)
	amount := sdk.NewCoin("trudex", sdk.NewInt(1000))
	trustake := sdk.NewCoin("trusteak", sdk.NewInt(1000))
	argument := "test argument"
	testURL, _ := url.Parse("http://www.trustory.io")
	evidence := []url.URL{*testURL}

	creator1 := fakeFundedCreator(ctx, k.bankKeeper)
	// remove cat coins to simulate backing conversion from trusteak
	k.bankKeeper.SubtractCoins(ctx, creator1, sdk.Coins{amount})
	// add trustake
	k.bankKeeper.AddCoins(ctx, creator1, sdk.Coins{trustake})

	creator2 := fakeFundedCreator(ctx, k.bankKeeper)
	creator3 := fakeFundedCreator(ctx, k.bankKeeper)
	creator4 := fakeFundedCreator(ctx, k.bankKeeper)
	creator5 := fakeFundedCreator(ctx, k.bankKeeper)
	creator6 := fakeFundedCreator(ctx, k.bankKeeper)
	creator7 := fakeFundedCreator(ctx, k.bankKeeper)
	creator8 := fakeFundedCreator(ctx, k.bankKeeper)
	creator9 := fakeFundedCreator(ctx, k.bankKeeper)

	// fake backings
	duration := 1 * time.Hour
	b1id, _ := k.backingKeeper.Create(ctx, storyID, trustake, argument, creator1, duration, evidence)
	b2id, _ := k.backingKeeper.Create(ctx, storyID, amount, argument, creator2, duration, evidence)
	b3id, _ := k.backingKeeper.Create(ctx, storyID, amount, argument, creator3, duration, evidence)
	b4id, _ := k.backingKeeper.Create(ctx, storyID, amount, argument, creator4, duration, evidence)

	// fake challenges
	c1id, _ := k.challengeKeeper.Create(ctx, storyID, amount, argument, creator5, evidence)
	c2id, _ := k.challengeKeeper.Create(ctx, storyID, amount, argument, creator6, evidence)

	// fake votes
	v1id, _ := k.Create(ctx, storyID, amount, true, argument, creator7, evidence)
	v2id, _ := k.Create(ctx, storyID, amount, true, argument, creator8, evidence)
	v3id, _ := k.Create(ctx, storyID, amount, false, argument, creator9, evidence)

	b1, _ := k.backingKeeper.Backing(ctx, b1id)
	// fake an interest
	b1.Interest = sdk.NewCoin("trudex", sdk.NewInt(500))
	k.backingKeeper.Update(ctx, b1)

	b2, _ := k.backingKeeper.Backing(ctx, b2id)
	b2.Interest = sdk.NewCoin("trudex", sdk.NewInt(500))
	k.backingKeeper.Update(ctx, b2)

	b3, _ := k.backingKeeper.Backing(ctx, b3id)
	b3.Interest = sdk.NewCoin("trudex", sdk.NewInt(500))
	k.backingKeeper.Update(ctx, b3)

	b4, _ := k.backingKeeper.Backing(ctx, b4id)
	b4.Interest = sdk.NewCoin("trudex", sdk.NewInt(500))
	k.backingKeeper.Update(ctx, b4)
	// change last backing vote to FALSE
	k.backingKeeper.ToggleVote(ctx, b4.ID())

	c1, _ := k.challengeKeeper.Challenge(ctx, c1id)
	c2, _ := k.challengeKeeper.Challenge(ctx, c2id)

	v1, _ := k.TokenVote(ctx, v1id)
	v2, _ := k.TokenVote(ctx, v2id)
	v3, _ := k.TokenVote(ctx, v3id)

	votes.trueVotes = append(votes.trueVotes, b1, b2, b3, v1, v2)
	votes.falseVotes = append(votes.falseVotes, b4, c1, c2, v3)

	return
}
