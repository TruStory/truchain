package chttp

import (
	"encoding/json"
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	crypto "github.com/tendermint/tendermint/crypto"
	ed "github.com/tendermint/tendermint/crypto/ed25519"
	secp "github.com/tendermint/tendermint/crypto/secp256k1"
	tcmn "github.com/tendermint/tendermint/libs/common"
)

// PresignedRequest represents a JSON request body from a user wishing to broadcast a transaction that they've signed locally
type PresignedRequest struct {
	MsgTypes   []string      `json:"msg_types"`
	Tx         tcmn.HexBytes `json:"tx"`
	PubKeyAlgo string        `json:"pubkey_algo"`
	PubKey     tcmn.HexBytes `json:"pubkey"`
	Signature  tcmn.HexBytes `json:"signature"`
}

// NewPresignedStdTx parses a `PresignedRequest` into an `auth.StdTx`
func (a *API) NewPresignedStdTx(r PresignedRequest) (auth.StdTx, error) {
	msgs, fee, signatures, doc, err := a.stdTxFragments(r)

	if err != nil {
		return auth.StdTx{}, err
	}

	tx := auth.NewStdTx(msgs, fee, signatures, doc.Memo)

	return tx, nil
}

func (a *API) stdTxFragments(r PresignedRequest) ([]sdk.Msg, auth.StdFee, []auth.StdSignature, auth.StdSignDoc, error) {
	doc, err := a.stdSignDoc(r.Tx.Bytes())

	if err != nil {
		fmt.Println("Error decoding StdSignDoc: ", err)
		return []sdk.Msg{}, auth.StdFee{}, []auth.StdSignature{}, auth.StdSignDoc{}, err
	}

	msgs, err := a.stdMsgs(r.MsgTypes, doc.Msgs)

	if err != nil {
		fmt.Println("Error decoding StdMsgs: ", err)
		return []sdk.Msg{}, auth.StdFee{}, []auth.StdSignature{}, auth.StdSignDoc{}, err
	}

	fee, err := a.stdFee(doc.Fee)

	if err != nil {
		fmt.Println("Error decoding StdFee: ", err)
		return []sdk.Msg{}, auth.StdFee{}, []auth.StdSignature{}, auth.StdSignDoc{}, err
	}

	signatures, err := a.stdSignatures(r, doc)

	if err != nil {
		fmt.Println("Error decoding StdSignatures: ", err)
		return []sdk.Msg{}, auth.StdFee{}, []auth.StdSignature{}, auth.StdSignDoc{}, err
	}

	return msgs, fee, signatures, doc, nil
}

func (a *API) stdSignDoc(bs []byte) (auth.StdSignDoc, error) {
	doc := new(auth.StdSignDoc)
	err := json.Unmarshal(bs, &doc)

	if err != nil {
		return auth.StdSignDoc{}, err
	}

	return *doc, nil
}

func (a *API) stdMsgs(msgTypes []string, msgBodies []json.RawMessage) ([]sdk.Msg, error) {
	msgs := []sdk.Msg{}

	if len(msgTypes) != len(msgBodies) {
		return msgs, expectedMessagesError(len(msgBodies), msgTypes)
	}

	for i, body := range msgBodies {
		typeName := msgTypes[i]
		msg, err := a.stdMsg(typeName, body)

		if err != nil {
			return msgs, err
		}

		msgs = append(msgs, msg)
	}

	return msgs, nil
}

func (a *API) stdMsg(name string, raw json.RawMessage) (sdk.Msg, error) {
	t := reflect.TypeOf(a.Supported[name])

	if t == nil {
		return bank.MsgSend{}, unsupportedMsgTypeError(name, a.supportedMsgTypeNames())
	}

	obj := reflect.New(t).Interface()
	err := json.Unmarshal(raw, &obj)

	if err != nil {
		fmt.Println("Error unmarshaling msg JSON: ", obj, string(raw), err)
		return bank.MsgSend{}, internalDecodingError(err.Error())
	}

	msg, ok := obj.(sdk.Msg)

	if !ok {
		return bank.MsgSend{}, internalDecodingError(name + " does not implement sdk.Msg")
	}

	return msg, nil
}

func (a *API) stdFee(fragment json.RawMessage) (auth.StdFee, error) {
	fee := new(auth.StdFee)
	err := json.Unmarshal(fragment, fee)

	if err != nil {
		fmt.Println(err, string(fragment))
		return auth.StdFee{}, err
	}

	return *fee, nil
}

func (a *API) stdSignatures(r PresignedRequest, d auth.StdSignDoc) ([]auth.StdSignature, error) {
	key, err := StdKey(r.PubKeyAlgo, r.PubKey)

	if err != nil {
		return []auth.StdSignature{}, err
	}

	stdSig := auth.StdSignature{
		PubKey:    key,
		Signature: r.Signature.Bytes(),
	}

	sigs := []auth.StdSignature{stdSig}

	return sigs, nil
}

func (a *API) supportedMsgTypeNames() []string {
	types := []string{}

	for k := range a.Supported {
		types = append(types, k)
	}

	return types
}

// StdKey returns an instance of `crypto.PubKey` using the given algorithm
func StdKey(algo string, bytes []byte) (crypto.PubKey, error) {
	switch algo {
	case "ed25519":
		ek := ed.PubKeyEd25519{}
		copy(ek[:], bytes)
		return ek, nil
	case "secp256k1":
		sk := secp.PubKeySecp256k1{}
		copy(sk[:], bytes)
		return sk, nil
	default:
		return secp.PubKeySecp256k1{}, unsupportedAlgoError(algo, []string{"ed25519", "secp256k1"})
	}
}
