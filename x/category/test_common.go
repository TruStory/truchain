package category

import (
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	amino "github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

func mockDB() (sdk.Context, Keeper) {
	db := dbm.NewMemDB()

	catKey := sdk.NewKVStoreKey("categories")

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(catKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())

	cdc := amino.NewCodec()
	cryptoAmino.RegisterAmino(cdc)
	RegisterAmino(cdc)

	ck := NewKeeper(catKey, cdc)

	return ctx, ck
}

func createFakeCategory(ctx sdk.Context, ck WriteKeeper) Category {
	existing, err := ck.GetCategory(ctx, 1)
	if err == nil {
		return existing
	}
	id := ck.NewCategory(ctx, "decentralized exchanges", "trudex", "category for experts in decentralized exchanges")
	cat, _ := ck.GetCategory(ctx, id)
	return cat
}
