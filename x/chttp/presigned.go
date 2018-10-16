package chttp

import (
  "encoding/json"
  "fmt"
  "reflect"
  
  auth "github.com/cosmos/cosmos-sdk/x/auth"
  crypto "github.com/tendermint/tendermint/crypto"
  ed "github.com/tendermint/tendermint/crypto/ed25519"
  sdk "github.com/cosmos/cosmos-sdk/types"
  secp "github.com/tendermint/tendermint/crypto/secp256k1"
  tcmn "github.com/tendermint/tendermint/libs/common"
)

type PresignedRequest struct {
  MsgTypes []string `json:"msgTypes"`
  Tx tcmn.HexBytes `json:"tx"`
  PubKeyAlgo string `json:"pubkeyAlgo"`
  PubKey tcmn.HexBytes `json:"pubkey"`
  Signature tcmn.HexBytes `json:"signature"`
}

func (a *Api) NewPresignedStdTx(r PresignedRequest) (*auth.StdTx, error) {
  msgs, fee, signatures, doc, err := a.stdTxFragments(r)

  if err != nil {
    return nil, err
  }

  tx := auth.NewStdTx(*msgs, *fee, *signatures, doc.Memo)

  return &tx, nil
}

func (a *Api) stdTxFragments(r PresignedRequest) (*[]sdk.Msg, *auth.StdFee, *[]auth.StdSignature, *auth.StdSignDoc, error) {
  doc, err := a.stdSignDoc(r.Tx.Bytes())

  if err != nil {
    fmt.Println("Error decoding StdSignDoc: ", err)
    return nil, nil, nil, nil, err
  }
  
  msgs, err := a.stdMsgs(r.MsgTypes, doc.Msgs)

  if err != nil {
    fmt.Println("Error decoding StdMsgs: ", err)
    return nil, nil, nil, nil, err
  }

  fee, err := a.stdFee(doc.Fee)

  if err != nil {
    fmt.Println("Error decoding StdFee: ", err)
    return nil, nil, nil, nil, err
  }

  signatures, err := a.stdSignatures(r, *doc)

  if err != nil {
    fmt.Println("Error decoding StdSignatures: ", err)
    return nil, nil, nil, nil, err
  }

  return msgs, fee, signatures, doc, nil
}

func (a *Api) stdSignDoc(bs []byte) (*auth.StdSignDoc, error) {
  doc := new(auth.StdSignDoc)
  err := json.Unmarshal(bs, &doc)

  if err != nil {
    return nil, err
  }

  return doc, nil
}

func (a *Api) stdMsgs(msgTypes []string, msgBodies []json.RawMessage) (*[]sdk.Msg, error) {
  msgs := []sdk.Msg{}
  finalTypeIndex := len(msgTypes) - 1

  if len(msgTypes) > len(msgBodies) {
    return nil, expectedMoreMessagesError(len(msgBodies), msgTypes)
  }

  for i, body := range msgBodies {
    var typeName string

    if i > finalTypeIndex {
      typeName = msgTypes[finalTypeIndex]
    } else {
      typeName = msgTypes[i]
    }

    msg, err := a.stdMsg(typeName, body)

    if err != nil {
      return nil, err
    }
    
    msgs = append(msgs, *msg)
  }

  return &msgs, nil
}

func (a *Api) stdMsg(name string, raw json.RawMessage) (*sdk.Msg, error) {
  t := reflect.TypeOf(a.Supported[name])

  if t == nil {
    return nil, unsupportedMsgTypeError(name, a.supportedMsgTypeNames())
  }
  
  obj := reflect.New(t).Interface()
  err := json.Unmarshal(raw, &obj)

  if err != nil {
    fmt.Println("Error unmarshaling msg JSON: ", obj, string(raw), err)
    return nil, internalDecodingError(err.Error())
  }

  msg, ok := obj.(sdk.Msg)

  if !ok {
    return nil, internalDecodingError(name + " does not implement sdk.Msg")
  }

  return &msg, nil
}

func (a *Api) stdFee(fragment json.RawMessage) (*auth.StdFee, error) {
  fee := new(auth.StdFee)
  fmt.Println(string(fragment))
  err := json.Unmarshal(fragment, fee)

  if err != nil {
    fmt.Println(err, string(fragment))
    return nil, err
  }

  return fee, nil
}

func (a *Api) stdSignatures(r PresignedRequest, d auth.StdSignDoc) (*[]auth.StdSignature, error) {
  key, err := StdKey(r.PubKeyAlgo, r.PubKey)
 
  if err != nil {
    return nil, err
  }

  stdSig := auth.StdSignature{
    PubKey: *key,
    Signature: r.Signature.Bytes(),
    AccountNumber: d.AccountNumber,
    Sequence: d.Sequence,
  }

  sigs := []auth.StdSignature{stdSig}

  return &sigs, nil
}

func (a *Api) supportedMsgTypeNames() []string {
  types := []string{}

  for k, _ := range a.Supported {
    types = append(types, string(k))
  }

  return types
}

func StdKey(algo string, bytes []byte) (*crypto.PubKey, error) {
  switch algo {
  case "ed25519":
    ek := ed.PubKeyEd25519{}
    copy(ek[:], bytes)
    key := crypto.PubKey(ek)
    return &key, nil
  case "secp256k1":
    sk := secp.PubKeySecp256k1{}
    copy(sk[:], bytes)
    key := crypto.PubKey(sk)
    fmt.Println("Got key from bytes", algo, bytes, sk)
    return &key, nil
  default:
    return nil, unsupportedAlgoError(algo, []string{"ed25519", "secp256k1"})
  }
}