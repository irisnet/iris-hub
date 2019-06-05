package bank

import (
	"fmt"

	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/app/v1/bank"
	"github.com/irisnet/irishub/client/context"
	sdk "github.com/irisnet/irishub/types"
	"github.com/tendermint/tendermint/crypto"
)

type BaseAccount struct {
	Address       sdk.AccAddress `json:"address"`
	Coins         []string       `json:"coins"`
	PubKey        crypto.PubKey  `json:"public_key"`
	AccountNumber uint64         `json:"account_number"`
	Sequence      uint64         `json:"sequence"`
}

func ConvertAccountCoin(cliCtx context.CLIContext, acc auth.Account) (BaseAccount, error) {
	var accCoins []string
	for _, coin := range acc.GetCoins() {
		coinString, err := cliCtx.ConvertCoinToMainUnit(coin.String())
		if err == nil {
			accCoins = append(accCoins, coinString[0])
		} else {
			accCoins = append(accCoins, coin.String())
		}

	}
	return BaseAccount{
		Address:       acc.GetAddress(),
		Coins:         accCoins,
		PubKey:        acc.GetPubKey(),
		AccountNumber: acc.GetAccountNumber(),
		Sequence:      acc.GetSequence(),
	}, nil
}

// BuildBankSendMsg builds the sending coins msg
func BuildBankSendMsg(from sdk.AccAddress, to sdk.AccAddress, coins sdk.Coins) sdk.Msg {
	input := bank.NewInput(from, coins)
	output := bank.NewOutput(to, coins)
	msg := bank.NewMsgSend([]bank.Input{input}, []bank.Output{output})
	return msg
}

// BuildBankBurnMsg builds the burning coin msg
func BuildBankBurnMsg(from sdk.AccAddress, coins sdk.Coins) sdk.Msg {
	msg := bank.NewMsgBurn(from, coins)
	return msg
}

// BuildBankFreezeMsg builds the freeze coin msg
func BuildBankFreezeMsg(owner sdk.AccAddress, coin sdk.Coin) sdk.Msg {
	msg := bank.NewMsgFreeze(owner, coin)
	return msg
}

// BuildBankUnfreezeMsg builds the unfreeze coin msg
func BuildBankUnfreezeMsg(owner sdk.AccAddress, coin sdk.Coin) sdk.Msg {
	msg := bank.NewMsgUnfreeze(owner, coin)
	return msg
}

type TokenStats struct {
	FrozenTokens sdk.Coins `json:"frozen_tokens"`
	LooseTokens  sdk.Coins `json:"loose_tokens"`
	BurnedTokens sdk.Coins `json:"burned_tokens"`
	BondedTokens sdk.Coins `json:"bonded_tokens"`
}

// String implements fmt.Stringer
func (ts TokenStats) String() string {
	return fmt.Sprintf(`TokenStats:
  Loose Tokens:  %s
  Burned Tokens:  %s
  Bonded Tokens:  %s`,
		ts.LooseTokens.MainUnitString(), ts.BurnedTokens.MainUnitString(), ts.BondedTokens.MainUnitString(),
	)
}
