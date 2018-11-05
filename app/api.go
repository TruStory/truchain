package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	params "github.com/TruStory/truchain/parameters"
	"github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/chttp"
	"github.com/TruStory/truchain/x/registration"
	"github.com/TruStory/truchain/x/truapi"
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
	log.Fatal(app.api.ListenAndServe(params.Hostname + ":" + params.Portname))
}

// RegisterKey generates a new address/account for a public key
// Implements chttp.App
func (app *TruChain) RegisterKey(k tcmn.HexBytes, algo string) (sdk.AccAddress, int64, sdk.Coins, error) {
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

	if err != nil {
		fmt.Println("TX Broadcast error: ", err, res)
		return sdk.AccAddress{}, 0, sdk.Coins{}, err
	}

	stored := app.accountMapper.GetAccount(*(app.blockCtx), sdk.AccAddress(addr))
	accaddr := sdk.AccAddress(addr)
	coins := stored.GetCoins()

	return accaddr, stored.GetAccountNumber(), coins, nil
}

// DeliverPresigned broadcasts a transaction to the network and returns the result.
// Implements chttp.App
func (app *TruChain) DeliverPresigned(tx auth.StdTx) (*trpctypes.ResultBroadcastTxCommit, error) {
	bz := app.codec.MustMarshalBinary(tx)
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

// GenerateAddress returns the first 20 characters of a ULID generated with github.com/oklog/ulid
func GenerateAddress() []byte {
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	ulidaddr := ulid.MustNew(ulid.Timestamp(t), entropy)
	addr := []byte(ulidaddr.String())[:20]

	return addr
}

func (app *TruChain) signedRegistrationTx(addr []byte, k tcmn.HexBytes, algo string) (auth.StdTx, error) {
	msg := registration.RegisterKeyMsg{Address: addr, PubKey: k, PubKeyAlgo: algo, Coins: params.InitialCoins}
	chainID := app.blockHeader.ChainID
	registrarAcc := app.accountMapper.GetAccount(*(app.blockCtx), []byte(types.RegistrarAccAddress))
	registrarNum := registrarAcc.GetAccountNumber()
	registrarSequence := registrarAcc.GetSequence()
	registrationMemo := "reg"

	// Sign tx as registrar
	bytesToSign := auth.StdSignBytes(chainID, registrarNum, registrarSequence, params.RegistrationFee, []sdk.Msg{msg}, registrationMemo)
	sigBytes, err := app.registrarKey.Sign(bytesToSign)

	if err != nil {
		return auth.StdTx{}, err
	}

	// Construct and submit signed tx
	tx := auth.StdTx{
		Msgs: []sdk.Msg{msg},
		Fee:  params.RegistrationFee,
		Signatures: []auth.StdSignature{auth.StdSignature{
			PubKey:        app.registrarKey.PubKey(),
			Signature:     sigBytes,
			AccountNumber: 1,
			Sequence:      registrarSequence,
		}},
		Memo: registrationMemo,
	}

	return tx, nil
}
