package db

import (
	"time"

	ts "github.com/TruStory/trucoin/x/trustory/types"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

	amino "github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"

	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/tendermint/tendermint/crypto/encoding/amino"
)

// Helper functions used for keeper tests

func createFakeStory(ms sdk.MultiStore, k TruKeeper) int64 {
	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	body := "Body of story."
	category := ts.DEX
	creator := sdk.AccAddress([]byte{1, 2})
	escrow := sdk.AccAddress([]byte{3, 4})
	storyType := ts.Default
	t := time.Date(2018, time.September, 13, 23, 0, 0, 0, time.UTC)

	storyID, _ := k.AddStory(ctx, body, category, creator, escrow, storyType, t, t)
	return storyID
}

func setupMultiStore() (sdk.MultiStore, *sdk.KVStoreKey, *sdk.KVStoreKey, *sdk.KVStoreKey) {
	db := dbm.NewMemDB()
	accKey := sdk.NewKVStoreKey("acc")
	storyKey := sdk.NewKVStoreKey("stories")
	voteKey := sdk.NewKVStoreKey("votes")
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(accKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(storyKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(voteKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()
	return ms, accKey, storyKey, voteKey
}

func makeCodec() *amino.Codec {
	cdc := amino.NewCodec()
	cryptoAmino.RegisterAmino(cdc)
	ts.RegisterAmino(cdc)
	cdc.RegisterInterface((*auth.Account)(nil), nil)
	cdc.RegisterConcrete(&auth.BaseAccount{}, "cosmos-sdk/BaseAccount", nil)
	return cdc
}
