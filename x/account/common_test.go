package account

import (
	"testing"

	app "github.com/TruStory/truchain/types"
	bankexported "github.com/TruStory/truchain/x/bank/exported"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
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
	referenceID uint64, txType bankexported.TransactionType, setters ...bankexported.TransactionSetter) (sdk.Coins, sdk.Error) {

	txn := transaction{to, coin}
	bk.Transactions = append(bk.Transactions, txn)
	return sdk.Coins{coin}, nil
}

func (bk bankKeeper) IterateUserTransactions(ctx sdk.Context, creator sdk.AccAddress, reverse bool, cb func(transaction bankexported.Transaction) (stop bool)) {

}

func mockDB(t *testing.T) (sdk.Context, Keeper) {
	db := dbm.NewMemDB()

	authKey := sdk.NewKVStoreKey(ModuleName)
	accountKey := sdk.NewKVStoreKey(auth.StoreKey)
	paramsKey := sdk.NewKVStoreKey(params.StoreKey)
	transientParamsKey := sdk.NewTransientStoreKey(params.TStoreKey)
	supplyKey := sdk.NewKVStoreKey(supply.StoreKey)

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(accountKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(paramsKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(transientParamsKey, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(supplyKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())

	codec := codec.New()
	cryptoAmino.RegisterAmino(codec)
	RegisterCodec(codec)
	codec.RegisterInterface((*auth.Account)(nil), nil)
	codec.RegisterConcrete(&auth.BaseAccount{}, "auth/Account", nil)
	supply.RegisterCodec(codec)

	maccPerms := map[string][]string{
		auth.FeeCollectorName:     nil,
		mint.ModuleName:           {supply.Minter},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		gov.ModuleName:            {supply.Burner},
		UserGrowthPoolName:        {supply.Burner, supply.Staking},
	}

	paramsKeeper := params.NewKeeper(codec, paramsKey, transientParamsKey, params.DefaultCodespace)
	accountKeeper := auth.NewAccountKeeper(codec, accountKey, paramsKeeper.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	cosmosBankKeeper := bank.NewBaseKeeper(accountKeeper, paramsKeeper.Subspace(bank.DefaultParamspace), bank.DefaultCodespace, nil)
	supplyKeeper := supply.NewKeeper(codec, supplyKey, accountKeeper, cosmosBankKeeper, maccPerms)

	userGrowthAcc := supply.NewEmptyModuleAccount(UserGrowthPoolName, supply.Burner, supply.Staking)
	initCoins := sdk.NewCoins(sdk.NewCoin(app.StakeDenom, sdk.NewInt(10000000000)))
	err := userGrowthAcc.SetCoins(initCoins)
	require.NoError(t, err)
	supplyKeeper.SetModuleAccount(ctx, userGrowthAcc)

	totalSupply := initCoins
	supplyKeeper.SetSupply(ctx, supply.NewSupply(totalSupply))

	bankKeeper := bankKeeper{
		Transactions: []transaction{},
	}
	authKeeper := NewKeeper(authKey, paramsKeeper.Subspace(ModuleName), codec, bankKeeper, accountKeeper, supplyKeeper)

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
		sdk.NewInt64Coin(app.StakeDenom, 10000000),
	}
}

func getFakeKeyPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	key := secp256k1.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}
