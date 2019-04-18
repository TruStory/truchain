package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/chttp"
	"github.com/TruStory/truchain/x/truapi"
	"github.com/TruStory/truchain/x/users"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/oklog/ulid"
	abci "github.com/tendermint/tendermint/abci/types"
	tcmn "github.com/tendermint/tendermint/libs/common"
	trpc "github.com/tendermint/tendermint/rpc/core"
	trpctypes "github.com/tendermint/tendermint/rpc/core/types"
)

type truChainContextKey string

const keepersContextKey truChainContextKey = "keepers"
const storyKeeperKey = "storyKeeper"

func (app *TruChain) makeAPI() *truapi.TruAPI {
	aa := chttp.App(app)
	return truapi.NewTruAPI(&aa)
}

func (app *TruChain) startAPI() {
	app.api.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			keepers := map[string]interface{}{
				storyKeeperKey: app.storyKeeper,
			}

			ctxWithKeepers := context.WithValue(r.Context(), keepersContextKey, keepers)

			next.ServeHTTP(w, r.WithContext(ctxWithKeepers))
		})
	})

	app.api.RegisterResolvers()
	app.api.RegisterRoutes()
	app.api.GraphQLClient.GenerateSchema()
	log.Fatal(app.api.ListenAndServe(net.JoinHostPort(types.Hostname, types.Portname)))
}

// RegisterKey generates a new address/account for a public key
// Implements chttp.App
func (app *TruChain) RegisterKey(k tcmn.HexBytes, algo string) (sdk.AccAddress, uint64, sdk.Coins, error) {
	var addr []byte

	if string(algo[0]) == "*" {
		addr = []byte("cosmostestingaddress")
		algo = algo[1:]
	} else {
		addr = GenerateAddress()
	}

	tx, err := app.signedRegistrationTx(addr, k, algo)

	if err != nil {
		fmt.Println("TX Parse error: ", err, tx)
		return sdk.AccAddress{}, 0, sdk.Coins{}, err
	}

	res, err := app.DeliverPresigned(tx)

	if !res.CheckTx.IsOK() {
		fmt.Println("TX Broadcast CheckTx error: ", res.CheckTx.Log)
		return sdk.AccAddress{}, 0, sdk.Coins{}, errors.New(res.CheckTx.Log)
	}

	if !res.DeliverTx.IsOK() {
		fmt.Println("TX Broadcast DeliverTx error: ", res.DeliverTx.Log)
		return sdk.AccAddress{}, 0, sdk.Coins{}, errors.New(res.DeliverTx.Log)
	}

	if err != nil {
		fmt.Println("TX Broadcast error: ", err, res)
		return sdk.AccAddress{}, 0, sdk.Coins{}, err
	}

	accaddr := sdk.AccAddress(addr)
	stored := app.accountKeeper.GetAccount(*(app.blockCtx), accaddr)

	if stored == nil {
		return sdk.AccAddress{}, 0, sdk.Coins{}, errors.New("Unable to locate account " + string(addr))
	}

	coins := stored.GetCoins()

	return accaddr, stored.GetAccountNumber(), coins, nil
}

// DeliverPresigned broadcasts a transaction to the network and returns the result.
// Implements chttp.App
func (app *TruChain) DeliverPresigned(tx auth.StdTx) (*trpctypes.ResultBroadcastTxCommit, error) {
	bz := app.codec.MustMarshalBinaryLengthPrefixed(tx)
	return trpc.BroadcastTxCommit(bz)
}

// RunQuery takes a querier path string and parameters map, and returns the results of the query.
// Implements chttp.App
func (app *TruChain) RunQuery(path string, params interface{}) abci.ResponseQuery {
	bz, err := json.Marshal(params)

	if err != nil {
		return abci.ResponseQuery{Log: err.Error()}
	}

	return app.Query(abci.RequestQuery{Data: bz, Path: "/custom/" + path})
}

// GenerateAddress returns the first 20 characters of a ULID (https://github.com/oklog/ulid)
func GenerateAddress() []byte {
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	ulidaddr := ulid.MustNew(ulid.Timestamp(t), entropy)
	addr := []byte(ulidaddr.String())[:20]

	return addr
}

func (app *TruChain) signedRegistrationTx(addr []byte, k tcmn.HexBytes, algo string) (auth.StdTx, error) {
	msg := users.RegisterKeyMsg{
		Address:    addr,
		PubKey:     k,
		PubKeyAlgo: algo,
		Coins:      app.initialCoins(),
	}
	chainID := app.blockHeader.ChainID
	registrarAcc := app.accountKeeper.GetAccount(*(app.blockCtx), []byte(types.RegistrarAccAddress))
	registrarNum := registrarAcc.GetAccountNumber()
	registrarSequence := registrarAcc.GetSequence()
	registrationMemo := "reg"

	// Sign tx as registrar
	bytesToSign := auth.StdSignBytes(chainID, registrarNum, registrarSequence, types.RegistrationFee, []sdk.Msg{msg}, registrationMemo)
	sigBytes, err := app.registrarKey.Sign(bytesToSign)

	if err != nil {
		return auth.StdTx{}, err
	}

	// Construct and submit signed tx
	tx := auth.StdTx{
		Msgs: []sdk.Msg{msg},
		Fee:  types.RegistrationFee,
		Signatures: []auth.StdSignature{auth.StdSignature{
			PubKey:    app.registrarKey.PubKey(),
			Signature: sigBytes,
		}},
		Memo: registrationMemo,
	}

	return tx, nil
}

func (app *TruChain) initialCoins() sdk.Coins {
	coins := sdk.Coins{}
	categories, err := app.categoryKeeper.GetAllCategories(*(app.blockCtx))
	if err != nil {
		panic(err)
	}

	for _, cat := range categories {
		coin := sdk.NewCoin(cat.Denom(), types.InitialCredAmount)
		coins = append(coins, coin)
	}

	coins = append(coins, types.InitialTruStake)

	// coins need to be sorted by denom to be valid
	coins.Sort()

	// yes we should panic if coins aren't valid
	// as it undermines the whole chain
	if !coins.IsValid() {
		panic("Initial coins are not valid.")
	}

	return coins
}
