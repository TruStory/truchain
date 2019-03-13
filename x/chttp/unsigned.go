package chttp

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/TruStory/truchain/x/db"
	"github.com/cosmos/cosmos-sdk/x/auth"
	ethsecp "github.com/ethereum/go-ethereum/crypto/secp256k1"
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

	// Hashing the tx
	hasher := sha256.New()
	hasher.Write([]byte(r.TxHash))
	hash := hasher.Sum(nil)

	// Signing the hash
	privateKey := GetPrivateKeyObject(keyPair)
	privateKeyBytes, _ := hex.DecodeString(fmt.Sprintf("%x", privateKey.D))
	publicKeyBytes, _ := hex.DecodeString(fmt.Sprintf("%x", ethsecp.CompressPubkey(privateKey.PublicKey.X, privateKey.PublicKey.Y)))
	signature, err := ethsecp.Sign(hash, privateKeyBytes)
	if err != nil {
		panic(err)
	}

	presignedRequest := &PresignedRequest{
		MsgTypes:   r.MsgTypes,
		Tx:         r.Tx,
		PubKeyAlgo: "secp256k1",
		PubKey:     publicKeyBytes,
		Signature:  signature,
	}

	return a.NewPresignedStdTx(*presignedRequest)
}

// GetPrivateKeyObject returns the secp's object encapsulating the private key
func GetPrivateKeyObject(keyPair db.KeyPair) *ecdsa.PrivateKey {
	privateKeyHex, _ := hex.DecodeString(keyPair.PrivateKey)
	privateKeyInt := big.NewInt(0)
	privateKeyInt.SetBytes(privateKeyHex)

	privateKeyObj := new(ecdsa.PrivateKey)
	privateKeyObj.PublicKey.Curve = ethsecp.S256()
	privateKeyObj.D = privateKeyInt
	privateKeyObj.PublicKey.X, privateKeyObj.PublicKey.Y = privateKeyObj.PublicKey.Curve.ScalarBaseMult(privateKeyInt.Bytes())

	return privateKeyObj
}

// func GetPrivateKeyObject(keyPair db.KeyPair) secp.PrivKeySecp256k1 {
// 	privateKey32Bytes := [32]byte{}
// 	privateKeyBytes, err := hex.DecodeString(keyPair.PrivateKey)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// make it of the fixed length of 32 bytes
// 	copy(privateKey32Bytes[:], privateKeyBytes)

// 	return secp.PrivKeySecp256k1(privateKey32Bytes)
// }
