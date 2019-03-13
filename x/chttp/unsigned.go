package chttp

import (
	"encoding/hex"

	"github.com/TruStory/truchain/x/db"
	"github.com/cosmos/cosmos-sdk/x/auth"
	secp "github.com/tendermint/tendermint/crypto/secp256k1"
	tcmn "github.com/tendermint/tendermint/libs/common"
)

// UnsignedRequest represents a JSON request body from a user wishing to broadcast a transaction that they've signed locally
type UnsignedRequest struct {
	MsgTypes []string      `json:"msg_types"`
	Tx       tcmn.HexBytes `json:"tx"`
	TxHash   string        `json:"tx_hash"`
}

// NewUnsignedStdTx parses an `UnsignedRequest` into an `auth.StdTx`
func (a *API) NewUnsignedStdTx(r UnsignedRequest, keyPair db.KeyPair) (auth.StdTx, error) {

	// Signing the hash
	txHashBytes := []byte(r.TxHash)
	privateKey := GetPrivateKeyObject(keyPair)
	signature, err := privateKey.Sign(txHashBytes)
	if err != nil {
		panic(err)
	}

	presignedRequest := &PresignedRequest{
		MsgTypes:   r.MsgTypes,
		Tx:         r.Tx,
		PubKeyAlgo: "secp256k1",
		PubKey:     privateKey.PubKey().Bytes(),
		Signature:  signature,
	}

	return a.NewPresignedStdTx(*presignedRequest)
}

// GetPrivateKeyObject returns the secp's object encapsulating the private key
func GetPrivateKeyObject(keyPair db.KeyPair) secp.PrivKeySecp256k1 {
	privateKey32Bytes := [32]byte{}
	privateKeyBytes, err := hex.DecodeString(keyPair.PrivateKey)
	if err != nil {
		panic(err)
	}

	// make it of the fixed length of 32 bytes
	copy(privateKey32Bytes[:], privateKeyBytes)

	return secp.PrivKeySecp256k1(privateKey32Bytes)
}
