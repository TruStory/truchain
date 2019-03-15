package chttp

import (
	"github.com/TruStory/truchain/x/db"
	"github.com/btcsuite/btcd/btcec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	tmcrypto "github.com/tendermint/tendermint/crypto"
	tcmn "github.com/tendermint/tendermint/libs/common"
)

// UnsignedRequest represents a JSON request body from a user wishing to broadcast a transaction that they've signed locally
type UnsignedRequest struct {
	MsgTypes []string      `json:"msg_types"`
	Tx       tcmn.HexBytes `json:"tx"`
	TxRaw    string        `json:"tx_raw"`
}

// NewUnsignedStdTx parses an `UnsignedRequest` into an `auth.StdTx`
func (a *API) NewUnsignedStdTx(r UnsignedRequest, keyPair db.KeyPair) (auth.StdTx, error) {

	// Hashing the tx
	hash := tmcrypto.Sha256([]byte(r.TxRaw))

	// Signing the hash
	privateKey := GetPrivateKeyObject(keyPair)
	signature, err := privateKey.Sign(hash)
	if err != nil {
		panic(err)
	}
	signatureBytes := serializeSig(signature)

	presignedRequest := &PresignedRequest{
		MsgTypes:   r.MsgTypes,
		Tx:         r.Tx,
		PubKeyAlgo: "secp256k1",
		PubKey:     privateKey.PubKey().SerializeCompressed(),
		Signature:  signatureBytes,
	}

	return a.NewPresignedStdTx(*presignedRequest)
}

// GetPrivateKeyObject returns the secp's object encapsulating the private key
func GetPrivateKeyObject(keyPair db.KeyPair) *btcec.PrivateKey {
	privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), []byte(keyPair.PrivateKey))

	return privKey
}

// returns the signature in the R||S format for tendermint
func serializeSig(sig *btcec.Signature) []byte {
	rBytes := sig.R.Bytes()
	sBytes := sig.S.Bytes()
	sigBytes := make([]byte, 64)
	// 0 pad the byte arrays from the left if they aren't big enough.
	copy(sigBytes[32-len(rBytes):32], rBytes)
	copy(sigBytes[64-len(sBytes):64], sBytes)
	return sigBytes
}
