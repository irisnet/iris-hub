package bank

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/irisnet/irishub/client/context"
	"github.com/tendermint/tendermint/crypto"
)

type BaseAccount struct {
	Address       sdk.AccAddress `json:"address"`
	Coins         []string       `json:"coins"`
	PubKey        crypto.PubKey  `json:"public_key"`
	AccountNumber int64          `json:"account_number"`
	Sequence      int64          `json:"sequence"`
}

func ConvertAccountCoin(cliCtx context.CLIContext, acc auth.Account) (BaseAccount, error) {
	coinsString, err := cliCtx.ConvertCoinToMainUnit(acc.GetCoins().String())
	if err != nil {
		return BaseAccount{}, err
	}
	return BaseAccount{
		Address:       acc.GetAddress(),
		Coins:         coinsString,
		PubKey:        acc.GetPubKey(),
		AccountNumber: acc.GetAccountNumber(),
		Sequence:      acc.GetSequence(),
	}, nil
}

// build the sendTx msg
func BuildMsg(from sdk.AccAddress, to sdk.AccAddress, coins sdk.Coins) sdk.Msg {
	input := bank.NewInput(from, coins)
	output := bank.NewOutput(to, coins)
	msg := bank.NewMsgSend([]bank.Input{input}, []bank.Output{output})
	return msg
}
