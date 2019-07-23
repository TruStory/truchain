package staking

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// verify interface at compile time
var _ sdk.Msg = &MsgSubmitArgument{}
var _ sdk.Msg = &MsgSubmitUpvote{}
var _ sdk.Msg = &MsgDeleteArgument{}
var _ sdk.Msg = &MsgEditArgument{}
var _ sdk.Msg = &MsgAddAdmin{}
var _ sdk.Msg = &MsgRemoveAdmin{}
var _ sdk.Msg = &MsgUpdateParams{}

const (
	TypeMsgSubmitArgument = "submit_argument"
	TypeMsgSubmitUpvote   = "submit_upvote"
	TypeMsgDeleteArgument = "delete_argument"
	TypeMsgEditArgument   = "edit_argument"
	TypeMsgAddAdmin       = "add_admin"
	TypeMsgRemoveAdmin    = "remove_admin"
	TypeMsgUpdateParams   = "update_params"
)

// MsgSubmitArgument msg for creating an argument.
type MsgSubmitArgument struct {
	ClaimID   uint64         `json:"claim_id"`
	Summary   string         `json:"summary"`
	Body      string         `json:"body"`
	StakeType StakeType      `json:"stake_type"`
	Creator   sdk.AccAddress `json:"creator"`
}

// NewMsgSubmitArgument returns a new submit argument message.
func NewMsgSubmitArgument(creator sdk.AccAddress, claimID uint64, summary, body string, stakeType StakeType) MsgSubmitArgument {
	return MsgSubmitArgument{
		ClaimID:   claimID,
		Summary:   summary,
		Body:      body,
		StakeType: stakeType,
		Creator:   creator,
	}
}
func (MsgSubmitArgument) Route() string {
	return RouterKey
}

func (MsgSubmitArgument) Type() string {
	return TypeMsgSubmitArgument
}

func (msg MsgSubmitArgument) ValidateBasic() sdk.Error {
	if !msg.StakeType.ValidForArgument() {
		return ErrCodeInvalidStakeType(msg.StakeType)
	}

	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress("Must provide a valid address")
	}

	if len(msg.Body) == 0 {
		return ErrCodeInvalidBodyLength()
	}

	if len(msg.Summary) == 0 {
		return ErrCodeInvalidSummaryLength()
	}
	return nil
}

// GetSignBytes gets the bytes for Msg signer to sign on
func (msg MsgSubmitArgument) GetSignBytes() []byte {
	msgBytes := ModuleCodec.MustMarshalJSON(msg)
	return sdk.MustSortJSON(msgBytes)
}

// GetSigners gets the signs of the Msg
func (msg MsgSubmitArgument) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
}

// MsgSubmitUpvote msg for upvoting an argument.
type MsgSubmitUpvote struct {
	ArgumentID uint64         `json:"argument_id"`
	Creator    sdk.AccAddress `json:"creator"`
}

func NewMsgSubmitUpvote(creator sdk.AccAddress, argumentID uint64) MsgSubmitUpvote {
	return MsgSubmitUpvote{
		ArgumentID: argumentID,
		Creator:    creator,
	}
}

func (MsgSubmitUpvote) Route() string {
	return RouterKey
}

func (MsgSubmitUpvote) Type() string {
	return TypeMsgSubmitUpvote
}

func (msg MsgSubmitUpvote) ValidateBasic() sdk.Error {
	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress("Must provide a valid address")
	}
	return nil
}

// GetSignBytes gets the bytes for Msg signer to sign on
func (msg MsgSubmitUpvote) GetSignBytes() []byte {
	msgBytes := ModuleCodec.MustMarshalJSON(msg)
	return sdk.MustSortJSON(msgBytes)
}

// GetSigners gets the signs of the Msg
func (msg MsgSubmitUpvote) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
}

// MsgDeleteArgument msg for deleting an argument.
type MsgDeleteArgument struct {
	ArgumentID uint64         `json:"argument_id"`
	Creator    sdk.AccAddress `json:"creator"`
}

func (MsgDeleteArgument) Route() string {
	return RouterKey
}

func (MsgDeleteArgument) Type() string {
	return TypeMsgDeleteArgument
}

func (msg MsgDeleteArgument) ValidateBasic() sdk.Error {
	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress("Must provide a valid address")
	}
	return nil
}

// GetSignBytes gets the bytes for Msg signer to sign on
func (msg MsgDeleteArgument) GetSignBytes() []byte {
	msgBytes := ModuleCodec.MustMarshalJSON(msg)
	return sdk.MustSortJSON(msgBytes)
}

// GetSigners gets the signs of the Msg
func (msg MsgDeleteArgument) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
}

// MsgEditArgument msg for creating an argument.
type MsgEditArgument struct {
	Creator    sdk.AccAddress `json:"creator"`
	ArgumentID uint64         `json:"argument_id"`
	Summary    string         `json:"summary"`
	Body       string         `json:"body"`
}

// NewMsgEditArgument returns a new edit argument message.
func NewMsgEditArgument(creator sdk.AccAddress, argumentID uint64, summary string, body string) MsgEditArgument {
	return MsgEditArgument{
		Creator:    creator,
		ArgumentID: argumentID,
		Summary:    summary,
		Body:       body,
	}
}

func (MsgEditArgument) Route() string {
	return RouterKey
}

func (MsgEditArgument) Type() string {
	return TypeMsgEditArgument
}

func (msg MsgEditArgument) ValidateBasic() sdk.Error {

	if len(msg.Body) == 0 {
		return ErrCodeInvalidBodyLength()
	}

	if len(msg.Summary) == 0 {
		return ErrCodeInvalidSummaryLength()
	}

	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress("Must provide a valid address")
	}

	return nil
}

// GetSignBytes gets the bytes for Msg signer to sign on
func (msg MsgEditArgument) GetSignBytes() []byte {
	msgBytes := ModuleCodec.MustMarshalJSON(msg)
	return sdk.MustSortJSON(msgBytes)
}

// GetSigners gets the address of the signer of the Msg
func (msg MsgEditArgument) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
}

// MsgAddAdmin defines the message to add a new admin
type MsgAddAdmin struct {
	Admin   sdk.AccAddress `json:"admin"`
	Creator sdk.AccAddress `json:"creator"`
}

// NewMsgAddAdmin returns the messages to add a new admin
func NewMsgAddAdmin(admin, creator sdk.AccAddress) MsgAddAdmin {
	return MsgAddAdmin{
		Admin:   admin,
		Creator: creator,
	}
}

// ValidateBasic implements Msg
func (msg MsgAddAdmin) ValidateBasic() sdk.Error {
	if len(msg.Admin) == 0 {
		return sdk.ErrInvalidAddress(fmt.Sprintf("Invalid address: %s", msg.Admin.String()))
	}

	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress(fmt.Sprintf("Invalid address: %s", msg.Creator.String()))
	}

	return nil
}

// Route implements Msg
func (msg MsgAddAdmin) Route() string { return RouterKey }

// Type implements Msg
func (msg MsgAddAdmin) Type() string { return TypeMsgAddAdmin }

// GetSignBytes implements Msg
func (msg MsgAddAdmin) GetSignBytes() []byte {
	msgBytes := ModuleCodec.MustMarshalJSON(msg)
	return sdk.MustSortJSON(msgBytes)
}

// GetSigners implements Msg. Returns the creator as the signer.
func (msg MsgAddAdmin) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Creator)}
}

// MsgRemoveAdmin defines the message to remove an admin
type MsgRemoveAdmin struct {
	Admin   sdk.AccAddress `json:"admin"`
	Remover sdk.AccAddress `json:"remover"`
}

// NewMsgRemoveAdmin returns the messages to remove an admin
func NewMsgRemoveAdmin(admin, remover sdk.AccAddress) MsgRemoveAdmin {
	return MsgRemoveAdmin{
		Admin:   admin,
		Remover: remover,
	}
}

// ValidateBasic implements Msg
func (msg MsgRemoveAdmin) ValidateBasic() sdk.Error {
	if len(msg.Admin) == 0 {
		return sdk.ErrInvalidAddress(fmt.Sprintf("Invalid address: %s", msg.Admin.String()))
	}

	if len(msg.Remover) == 0 {
		return sdk.ErrInvalidAddress(fmt.Sprintf("Invalid address: %s", msg.Remover.String()))
	}

	return nil
}

// Route implements Msg
func (msg MsgRemoveAdmin) Route() string { return RouterKey }

// Type implements Msg
func (msg MsgRemoveAdmin) Type() string { return TypeMsgRemoveAdmin }

// GetSignBytes implements Msg
func (msg MsgRemoveAdmin) GetSignBytes() []byte {
	msgBytes := ModuleCodec.MustMarshalJSON(msg)
	return sdk.MustSortJSON(msgBytes)
}

// GetSigners implements Msg. Returns the remover as the signer.
func (msg MsgRemoveAdmin) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Remover)}
}

// MsgUpdateParams defines the message to remove an admin
type MsgUpdateParams struct {
	Updates       Params         `json:"updates"`
	UpdatedFields []string       `json:"updated_fields"`
	Updater       sdk.AccAddress `json:"updater"`
}

// NewMsgUpdateParams returns the message to update the params
func NewMsgUpdateParams(updates Params, updatedFields []string, updater sdk.AccAddress) MsgUpdateParams {
	return MsgUpdateParams{
		Updates:       updates,
		UpdatedFields: updatedFields,
		Updater:       updater,
	}
}

// ValidateBasic implements Msg
func (msg MsgUpdateParams) ValidateBasic() sdk.Error {
	return nil
}

// Route implements Msg
func (msg MsgUpdateParams) Route() string { return RouterKey }

// Type implements Msg
func (msg MsgUpdateParams) Type() string { return TypeMsgUpdateParams }

// GetSignBytes implements Msg
func (msg MsgUpdateParams) GetSignBytes() []byte {
	msgBytes := ModuleCodec.MustMarshalJSON(msg)
	return sdk.MustSortJSON(msgBytes)
}

// GetSigners implements Msg. Returns the remover as the signer.
func (msg MsgUpdateParams) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Updater)}
}
