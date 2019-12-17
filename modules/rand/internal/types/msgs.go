package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// TypeMsgRequestRand represents the type of MsgRequestRand
	TypeMsgRequestRand = "request_rand"

	// DefaultBlockInterval is the default block interval
	DefaultBlockInterval = uint64(10)
)

var _ sdk.Msg = &MsgRequestRand{}

// MsgRequestRand represents a msg for requesting a random number
type MsgRequestRand struct {
	Consumer      sdk.AccAddress `json:"consumer" yaml:"consumer"`             // request address
	BlockInterval uint64         `json:"block_interval" yaml:"block_interval"` // block interval after which the requested random number will be generated
}

// NewMsgRequestRand constructs a MsgRequestRand
func NewMsgRequestRand(consumer sdk.AccAddress, blockInterval uint64) MsgRequestRand {
	return MsgRequestRand{
		Consumer:      consumer,
		BlockInterval: blockInterval,
	}
}

// Route implements Msg
func (msg MsgRequestRand) Route() string { return RouterKey }

// Type implements Msg
func (msg MsgRequestRand) Type() string { return TypeMsgRequestRand }

// ValidateBasic implements Msg
func (msg MsgRequestRand) ValidateBasic() sdk.Error {
	if len(msg.Consumer) == 0 {
		return ErrInvalidConsumer(DefaultCodespace, "the consumer address must be specified")
	}

	return nil
}

// GetSignBytes implements Msg
func (msg MsgRequestRand) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg
func (msg MsgRequestRand) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Consumer}
}
