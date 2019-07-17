package gov

import (
	"github.com/irisnet/irishub/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSubmitProposal{}, "irishub/gov/MsgSubmitProposal", nil)
	cdc.RegisterConcrete(MsgSubmitTxTaxUsageProposal{}, "irishub/gov/MsgSubmitTxTaxUsageProposal", nil)
	cdc.RegisterConcrete(MsgSubmitSoftwareUpgradeProposal{}, "irishub/gov/MsgSubmitSoftwareUpgradeProposal", nil)
	cdc.RegisterConcrete(MsgSubmitAddTokenProposal{}, "irishub/gov/MsgSubmitAddTokenProposal", nil)
	cdc.RegisterConcrete(MsgDeposit{}, "irishub/gov/MsgDeposit", nil)
	cdc.RegisterConcrete(MsgVote{}, "irishub/gov/MsgVote", nil)

	cdc.RegisterInterface((*Proposal)(nil), nil)
	cdc.RegisterConcrete(&BasicProposal{}, "irishub/gov/BasicProposal", nil)
	cdc.RegisterConcrete(&ParameterProposal{}, "irishub/gov/ParameterProposal", nil)
	cdc.RegisterConcrete(&PlainTextProposal{}, "irishub/gov/PlainTextProposal", nil)
	cdc.RegisterConcrete(&AddTokenProposal{}, "irishub/gov/AddTokenProposal", nil)
	cdc.RegisterConcrete(&SoftwareUpgradeProposal{}, "irishub/gov/SoftwareUpgradeProposal", nil)
	cdc.RegisterConcrete(&SystemHaltProposal{}, "irishub/gov/SystemHaltProposal", nil)
	cdc.RegisterConcrete(&TaxUsageProposal{}, "irishub/gov/TaxUsageProposal", nil)
	cdc.RegisterConcrete(&Vote{}, "irishub/gov/Vote", nil)
}

var msgCdc = codec.New()
