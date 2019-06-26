package account

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

// interface conformance check
var _ BankKeeper = bankKeeper{}

type transaction struct {
	Address sdk.Address
	Coin    sdk.Coin
}
type bankKeeper struct {
	Transactions []transaction
}

// AddCoin mock for bank keeper
func (bk bankKeeper) AddCoin(ctx sdk.Context, to sdk.AccAddress, coin sdk.Coin,
	referenceID uint64, txType int) (sdk.Coins, sdk.Error) {

	txn := transaction{to, coin}
	bk.Transactions = append(bk.Transactions, txn)
	return sdk.Coins{coin}, nil
}

func mockDB() (sdk.Context, Keeper) {
	db := dbm.NewMemDB()

	authKey := sdk.NewKVStoreKey(ModuleName)
	accountKey := sdk.NewKVStoreKey(auth.StoreKey)
	paramsKey := sdk.NewKVStoreKey(params.StoreKey)
	transientParamsKey := sdk.NewTransientStoreKey(params.TStoreKey)

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(accountKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(paramsKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(transientParamsKey, sdk.StoreTypeTransient, db)
	ms.LoadLatestVersion()

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())

	codec := codec.New()
	cryptoAmino.RegisterAmino(codec)
	RegisterCodec(codec)
	codec.RegisterInterface((*auth.Account)(nil), nil)

	paramsKeeper := params.NewKeeper(codec, paramsKey, transientParamsKey, params.DefaultCodespace)
	bankKeeper := bankKeeper{
		Transactions: []transaction{},
	}
	accountKeeper := auth.NewAccountKeeper(codec, accountKey, paramsKeeper.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	authKeeper := NewKeeper(authKey, paramsKeeper.Subspace(ModuleName), codec, bankKeeper, accountKeeper)

	InitGenesis(ctx, authKeeper, DefaultGenesisState())

	// setting registrar
	params := authKeeper.GetParams(ctx)
	params.Registrar = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()) // creating a new key
	authKeeper.SetParams(ctx, params)

	return ctx, authKeeper
}

func getFakeAppAccountParams() (privateKey crypto.PrivKey, publicKey crypto.PubKey, address sdk.AccAddress, coins sdk.Coins) {
	privateKey, publicKey, address = getFakeKeyPubAddr()
	coins = getFakeCoins()

	return
}

func getFakeCoins() sdk.Coins {
	return sdk.Coins{
		sdk.NewInt64Coin("fake", 10000000),
	}
}

func getFakeKeyPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	key := secp256k1.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}
