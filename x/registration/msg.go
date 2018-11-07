package registration

import (
	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tcmn "github.com/tendermint/tendermint/libs/common"
)

// RegisterKeyMsg describes a public key that should be registered with some initial coins
type RegisterKeyMsg struct {
	Address    sdk.AccAddress `json:"address"`
	PubKey     tcmn.HexBytes  `json:"pubkey"`
	PubKeyAlgo string         `json:"pubkey_algo"`
	Coins      sdk.Coins      `json:"coins"`
}

// Type implements Msg
func (msg RegisterKeyMsg) Type() string { return app.GetType(msg) }

// Route implements Msg
func (msg RegisterKeyMsg) Route() string { return app.GetName(msg) }

// GetSignBytes implements Msg
func (msg RegisterKeyMsg) GetSignBytes() []byte {
	return app.MustGetSignBytes(msg)
}

// ValidateBasic implements Msg
func (msg RegisterKeyMsg) ValidateBasic() sdk.Error {
	// TODO
	return nil
}

// GetSigners implements Msg
func (msg RegisterKeyMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress([]byte(app.RegistrarAccAddress))}
}
