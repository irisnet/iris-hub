package guardian

import (
	"github.com/irisnet/irishub/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgAddProfiler{}, "iris-hub/guardian/MsgAddProfiler", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
