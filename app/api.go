package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

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

func (app *TruChain) makeApi() *truapi.TruApi {
	aa := chttp.App(app)
	return truapi.NewTruApi(&aa)
}

func (app *TruChain) startApi() {
	app.api.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "keepers", map[string]interface{}{
				"storyKeeper": app.storyKeeper,
			})))
		})
	})

	app.api.RegisterRoutes()
	app.api.RegisterResolvers()
	log.Fatal(app.api.ListenAndServe("0.0.0.0:8080")) // TODO: Make port configurable [notduncansmith]
}

// Implements chttp.App
func (app *TruChain) RegisterKey(k tcmn.HexBytes, algo string) (*sdk.AccAddress, int64, *sdk.Coins, error) {
	addr := getUlid()
	tx := app.signedRegistrationTx(addr, k, algo)
	res, err := app.DeliverPresigned(tx)

	if err != nil {
		fmt.Println("TX Broadcast error: ", err, res)
		return nil, 0, nil, err
	}

	stored := app.accountMapper.GetAccount(*(app.blockCtx), sdk.AccAddress(addr))
	accaddr := sdk.AccAddress(addr)
	coins := stored.GetCoins()

	return &accaddr, stored.GetAccountNumber(), &coins, nil
}

// Implements chttp.App
func (app *TruChain) DeliverPresigned(tx auth.StdTx) (*trpctypes.ResultBroadcastTxCommit, error) {
	bz := app.codec.MustMarshalBinary(tx)
	return trpc.BroadcastTxCommit(bz)
}

// Implements chttp.App
func (app *TruChain) RunQuery(path string, params interface{}) abci.ResponseQuery {
	bz, err := json.Marshal(params)

	if err != nil {
		return abci.ResponseQuery{Log: err.Error()}
	}

	return app.Query(abci.RequestQuery{Data: bz, Path: "/custom/" + path})
}

func (app *TruChain) signedRegistrationTx(addr []byte, k tcmn.HexBytes, algo string) auth.StdTx {
	msg := registration.RegisterKeyMsg{Address: addr, PubKey: k, PubKeyAlgo: algo, Coins: initialCoins}
	chainId := app.blockHeader.ChainID
	registrarAcc := app.accountMapper.GetAccount(*(app.blockCtx), []byte(types.RegistrarAccAddress))
	registrarNum := registrarAcc.GetAccountNumber()
	registrarSequence := registrarAcc.GetSequence()
	registrationMemo := "reg"

	// Sign tx as registrar
	bytesToSign := auth.StdSignBytes(chainId, registrarNum, registrarSequence, registrationFee, []sdk.Msg{msg}, registrationMemo)
	sigBytes, err := app.registrarKey.Sign(bytesToSign)

	if err != nil {
		panic(err)
	}

	// Construct and submit signed tx
	tx := auth.StdTx{
		Msgs: []sdk.Msg{msg},
		Fee:  registrationFee,
		Signatures: []auth.StdSignature{auth.StdSignature{
			PubKey:        app.registrarKey.PubKey(),
			Signature:     sigBytes,
			AccountNumber: 1,
			Sequence:      registrarSequence,
		}},
		Memo: registrationMemo,
	}

	return tx
}

func getUlid() []byte {
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	ulidaddr := ulid.MustNew(ulid.Timestamp(t), entropy)
	addr := []byte(ulidaddr.String())[:20]

	return addr
}
