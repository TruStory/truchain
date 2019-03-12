package voting

import (
	"crypto/rand"
	"net/url"
	"time"

	"github.com/TruStory/truchain/x/stake"
	"github.com/TruStory/truchain/x/trubank"
	"github.com/TruStory/truchain/x/vote"

	"github.com/TruStory/truchain/x/backing"

	"github.com/TruStory/truchain/x/challenge"

	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/category"
	"github.com/TruStory/truchain/x/story"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	sdkparams "github.com/cosmos/cosmos-sdk/x/params"
	amino "github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

func mockDB() (sdk.Context, Keeper, category.Keeper) {

	db := dbm.NewMemDB()
	accKey := sdk.NewKVStoreKey("acc")
	storyKey := sdk.NewKVStoreKey("stories")
	storyQueueKey := sdk.NewKVStoreKey(story.PendingQueueStoreKey)
	expiredStoryQueueKey := sdk.NewKVStoreKey(story.ExpiringQueueStoreKey)
	catKey := sdk.NewKVStoreKey("categories")
	challengeKey := sdk.NewKVStoreKey("challenges")
	gameKey := sdk.NewKVStoreKey("games")
	pendingGameListKey := sdk.NewKVStoreKey("pendingGameList")
	votingStoryQueueKey := sdk.NewKVStoreKey("gameQueue")
	voteKey := sdk.NewKVStoreKey("vote")
	votingKey := sdk.NewKVStoreKey(StoreKey)
	backingKey := sdk.NewKVStoreKey("backing")
	paramsKey := sdk.NewKVStoreKey(sdkparams.StoreKey)
	transientParamsKey := sdk.NewTransientStoreKey(sdkparams.TStoreKey)

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(accKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(storyKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(storyQueueKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(expiredStoryQueueKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(catKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(challengeKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(gameKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(pendingGameListKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(votingStoryQueueKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(voteKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(votingKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(backingKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(paramsKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(transientParamsKey, sdk.StoreTypeTransient, db)
	ms.LoadLatestVersion()

	header := abci.Header{Time: time.Now()}
	ctx := sdk.NewContext(ms, header, false, log.NewNopLogger())

	codec := amino.NewCodec()
	cryptoAmino.RegisterAmino(codec)
	codec.RegisterInterface((*auth.Account)(nil), nil)
	codec.RegisterConcrete(&auth.BaseAccount{}, "auth/Account", nil)

	pk := sdkparams.NewKeeper(codec, paramsKey, transientParamsKey)
	am := auth.NewAccountKeeper(codec, accKey, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(am,
		pk.Subspace(bank.DefaultParamspace),
		bank.DefaultCodespace,
	)
	ck := category.NewKeeper(catKey, codec)
	category.InitGenesis(ctx, ck, category.DefaultCategories())

	sk := story.NewKeeper(
		storyKey,
		storyQueueKey,
		expiredStoryQueueKey,
		votingStoryQueueKey,
		ck,
		pk.Subspace(story.StoreKey),
		codec)

	story.InitGenesis(ctx, sk, story.DefaultGenesisState())

	truBankKey := sdk.NewKVStoreKey(trubank.StoreKey)
	ms.MountStoreWithDB(truBankKey, sdk.StoreTypeIAVL, db)
	truBankKeeper := trubank.NewKeeper(
		truBankKey,
		bankKeeper,
		ck,
		codec)

	stakeKeeper := stake.NewKeeper(
		sk,
		truBankKeeper,
		pk.Subspace(stake.StoreKey),
	)
	stake.InitGenesis(ctx, stakeKeeper, stake.DefaultGenesisState())

	backingKeeper := backing.NewKeeper(
		backingKey,
		stakeKeeper,
		sk,
		bankKeeper,
		ck,
		codec,
	)

	challengeKeeper := challenge.NewKeeper(
		challengeKey,
		stakeKeeper,
		backingKeeper,
		bankKeeper,
		sk,
		pk.Subspace(challenge.StoreKey),
		codec,
	)
	challenge.InitGenesis(ctx, challengeKeeper, challenge.DefaultGenesisState())

	voteKeeper := vote.NewKeeper(
		voteKey,
		votingStoryQueueKey,
		stakeKeeper,
		am,
		backingKeeper,
		challengeKeeper,
		sk,
		bankKeeper,
		pk.Subspace(vote.StoreKey),
		codec,
	)
	vote.InitGenesis(ctx, voteKeeper, vote.DefaultGenesisState())

	k := NewKeeper(
		votingKey,
		votingStoryQueueKey,
		am,
		backingKeeper,
		challengeKeeper,
		stakeKeeper,
		sk,
		voteKeeper,
		bankKeeper,
		truBankKeeper,
		pk.Subspace(StoreKey),
		codec)
	InitGenesis(ctx, k, DefaultGenesisState())

	return ctx, k, ck
}

func createFakeStory(ctx sdk.Context, sk story.WriteKeeper, ck category.WriteKeeper) int64 {
	body := "TruStory validators can be bootstrapped with a single genesis file."
	creator := sdk.AccAddress([]byte{1, 2})
	storyType := story.Default
	source := url.URL{}
	categoryID := int64(1)

	storyID, _ := sk.Create(ctx, body, categoryID, creator, source, storyType)

	return storyID
}

func fakeFundedCreator(ctx sdk.Context, k bank.Keeper) sdk.AccAddress {
	bz := make([]byte, 4)
	rand.Read(bz)
	creator := sdk.AccAddress(bz)

	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(2000000000000))
	k.AddCoins(ctx, creator, sdk.Coins{amount})

	return creator
}

func fakeConfirmedGame() (ctx sdk.Context, votes poll, k Keeper) {
	ctx, k, ck := mockDB()

	storyID := createFakeStory(ctx, k.storyKeeper, ck)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(1000000000000))
	largeAmount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(2000000000000))
	argument := "test argument"

	// GAME SCENARIO (STORY WILL BE CONFIRMED)
	// 4 backers @ 1000, 2 challengers @ 1000, 1 challenger @ 2000
	// game meets 100% challenge threshold
	// total reward pool = 4000 (backers) + 4000 (challengers)

	// 4 backers, interest each for 24 hours = 66,733,300,000 = 66.73
	// total interest = 266.92 shanev

	// 2 challengers, interest each for 24 hours = 66,733,300,000 = 66.73
	// 1 challengers, interest each for 24 hours = 66,733,300,000 = 66.73 *2 (double amount)
	// total interest = 266.92 shanev

	// VOTING BEGINS (New Voters)
	// 3 True Votes @ 1000 each, 1 False Vote @ 1000 each

	// GAME ENDS
	// 7 TRUE VOTES (4 Backers, 3 True Voters)
	// 4 FALSE VOTES (3 Challengers, 1 False Voter)

	// True Total = 4000 from Backers + 3000 from Voters = 7000 (before interest)
	// False Total = 4000 from Challengers + 1000 from Voter = 5000 (before interest)

	// Total Reward Pool = False Total (5000) + False Interest (266) = 5266
	// 75% of pool = 3949
	// 25% of pool = 1317

	creator1 := fakeFundedCreator(ctx, k.bankKeeper)
	creator2 := fakeFundedCreator(ctx, k.bankKeeper)
	creator3 := fakeFundedCreator(ctx, k.bankKeeper)
	creator4 := fakeFundedCreator(ctx, k.bankKeeper)
	creator5 := fakeFundedCreator(ctx, k.bankKeeper)
	creator6 := fakeFundedCreator(ctx, k.bankKeeper)
	creator7 := fakeFundedCreator(ctx, k.bankKeeper)
	creator8 := fakeFundedCreator(ctx, k.bankKeeper)
	creator9 := fakeFundedCreator(ctx, k.bankKeeper)
	creator10 := fakeFundedCreator(ctx, k.bankKeeper) // largeAmountChallenger
	creator11 := fakeFundedCreator(ctx, k.bankKeeper)

	// fake backings
	// each should end up with: 2000trusteak, 1054cred (from false voters)
	// 1054cred = (3949 / 4) + 66.73
	b1id, _ := k.backingKeeper.Create(ctx, storyID, amount, argument, creator1, false)
	b2id, _ := k.backingKeeper.Create(ctx, storyID, amount, argument, creator2, false)
	b3id, _ := k.backingKeeper.Create(ctx, storyID, amount, argument, creator3, false)
	b4id, _ := k.backingKeeper.Create(ctx, storyID, amount, argument, creator4, false)

	// fake challenges
	// c1 & c2 should end up with: 1000trusteak (they lost 1000trusteak in staking)
	// and no interest (they lost)
	// c3 should end up with 0trusteak (since they staked 2000trusteak)
	c1id, _ := k.challengeKeeper.Create(ctx, storyID, amount, argument, creator5, false)
	c2id, _ := k.challengeKeeper.Create(ctx, storyID, amount, argument, creator6, false)
	c3id, _ := k.challengeKeeper.Create(ctx, storyID, largeAmount, argument, creator10, false)

	// fake votes (true)
	// each should end up with: 2000trusteak, 439cred (from false voters)
	// 439cred = (1317 / 3)
	v1id, _ := k.voteKeeper.Create(ctx, storyID, amount, true, argument, creator7)
	v2id, _ := k.voteKeeper.Create(ctx, storyID, amount, true, argument, creator8)
	v4id, _ := k.voteKeeper.Create(ctx, storyID, amount, true, argument, creator11)

	// fake votes(false)
	// ends up with: 1000trusteak (lost 1000trusteak)
	v3id, _ := k.voteKeeper.Create(ctx, storyID, amount, false, argument, creator9)

	b1, _ := k.backingKeeper.Backing(ctx, b1id)
	b2, _ := k.backingKeeper.Backing(ctx, b2id)
	b3, _ := k.backingKeeper.Backing(ctx, b3id)
	b4, _ := k.backingKeeper.Backing(ctx, b4id)

	c1, _ := k.challengeKeeper.Challenge(ctx, c1id)
	c2, _ := k.challengeKeeper.Challenge(ctx, c2id)
	c3, _ := k.challengeKeeper.Challenge(ctx, c3id)

	v1, _ := k.voteKeeper.TokenVote(ctx, v1id)
	v2, _ := k.voteKeeper.TokenVote(ctx, v2id)
	v3, _ := k.voteKeeper.TokenVote(ctx, v3id)
	v4, _ := k.voteKeeper.TokenVote(ctx, v4id)

	votes.trueVotes = append(votes.trueVotes, b1, b2, b3, b4, v1, v2, v4)
	votes.falseVotes = append(votes.falseVotes, c1, c2, c3, v3)

	return
}

func fakeConfirmedGameNoStakers() (ctx sdk.Context, votes poll, k Keeper) {
	ctx, k, ck := mockDB()

	storyID := createFakeStory(ctx, k.storyKeeper, ck)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(1000000000000))
	argument := "test argument"

	creator1 := fakeFundedCreator(ctx, k.bankKeeper)
	creator2 := fakeFundedCreator(ctx, k.bankKeeper)
	creator3 := fakeFundedCreator(ctx, k.bankKeeper)

	c1id, _ := k.challengeKeeper.Create(ctx, storyID, amount, argument, creator1, false)

	v1id, _ := k.voteKeeper.Create(ctx, storyID, amount, true, argument, creator2)
	v2id, _ := k.voteKeeper.Create(ctx, storyID, amount, true, argument, creator3)

	c1, _ := k.challengeKeeper.Challenge(ctx, c1id)

	v1, _ := k.voteKeeper.TokenVote(ctx, v1id)
	v2, _ := k.voteKeeper.TokenVote(ctx, v2id)

	votes.trueVotes = append(votes.trueVotes, v1, v2)
	votes.falseVotes = append(votes.falseVotes, c1)

	return
}

func fakeRejectedGame() (ctx sdk.Context, votes poll, k Keeper) {
	ctx, k, ck := mockDB()

	storyID := createFakeStory(ctx, k.storyKeeper, ck)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(1000000000000))
	argument := "test argument"

	// GAME SCENARIO (STORY WILL BE REJECTED)
	// 1 challenger, interest for 24 hours = 66,733,300,000 = 66.73

	// GAME ENDS
	// 0 TRUE VOTES
	// 1 FALSE VOTES (1 Challenger)

	// True Total = 0 (before interest)
	// False Total = 1000 from Challenger = 1000 (before interest)

	// Total Reward Pool = True Total (0) + True Interest (0) = 0
	// 75% of pool = 0
	// 25% of pool = 0

	creator1 := fakeFundedCreator(ctx, k.bankKeeper)

	// fake challenges
	c1id, _ := k.challengeKeeper.Create(ctx, storyID, amount, argument, creator1, false)
	c1, _ := k.challengeKeeper.Challenge(ctx, c1id)

	votes.falseVotes = append(votes.falseVotes, c1)

	return
}
