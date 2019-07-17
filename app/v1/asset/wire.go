package asset

import (
	"github.com/irisnet/irishub/codec"
)

// Register concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgIssueToken{}, "irishub/asset/MsgIssueToken", nil)
	cdc.RegisterConcrete(MsgCreateGateway{}, "irishub/asset/MsgCreateGateway", nil)
	cdc.RegisterConcrete(MsgEditGateway{}, "irishub/asset/MsgEditGateway", nil)
	cdc.RegisterConcrete(MsgEditToken{}, "irishub/asset/MsgEditToken", nil)
	cdc.RegisterConcrete(MsgTransferGatewayOwner{}, "irishub/asset/MsgTransferGatewayOwner", nil)
	cdc.RegisterConcrete(MsgMintToken{}, "irishub/asset/MsgMintToken", nil)
	cdc.RegisterConcrete(MsgTransferTokenOwner{}, "irishub/asset/MsgTransferTokenOwner", nil)

	cdc.RegisterConcrete(BaseToken{}, "irishub/asset/BaseToken", nil)
	cdc.RegisterConcrete(FungibleToken{}, "irishub/asset/FungibleToken", nil)

	cdc.RegisterConcrete(&Params{}, "irishub/asset/Params", nil)
	cdc.RegisterConcrete(&Gateway{}, "irishub/asset/Gateway", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
