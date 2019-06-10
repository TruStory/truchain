package auth

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

const (
	// TypeMsgRegisterKey represents the type of the message for registering the key
	TypeMsgRegisterKey = "register_key"
)

// MsgRegisterKey defines the message to register a new key
type MsgRegisterKey struct {
	Address    sdk.AccAddress `json:"address"`
	PubKey     crypto.PubKey  `json:"public_key"`
	PubKeyAlgo string         `json:"public_key_algo"`
	Coins      sdk.Coins      `json:"coins"`
}

// NewMsgRegisterKey returns the messages to register a new key
func NewMsgRegisterKey(address sdk.AccAddress, publicKey crypto.PubKey, publicKeyAlgo string, coins sdk.Coins) MsgRegisterKey {
	return MsgRegisterKey{
		Address:    address,
		PubKey:     publicKey,
		PubKeyAlgo: publicKeyAlgo,
		Coins:      coins,
	}
}

// ValidateBasic implements Msg
func (msg MsgRegisterKey) ValidateBasic() sdk.Error {
	if len(msg.Address) == 0 {
		return sdk.ErrInvalidAddress(fmt.Sprintf("Invalid address: %s", msg.Address.String()))
	}

	return nil
}

// Route implements Msg
func (msg MsgRegisterKey) Route() string { return RouterKey }

// Type implements Msg
func (msg MsgRegisterKey) Type() string { return TypeMsgRegisterKey }

// GetSignBytes implements Msg
func (msg MsgRegisterKey) GetSignBytes() []byte {
	msgBytes := moduleCodec.MustMarshalJSON(msg)
	return sdk.MustSortJSON(msgBytes)
}

// GetSigners implements Msg. Returns the creator as the signer.
func (msg MsgRegisterKey) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(DefaultGenesisState().Registrar)}
}
