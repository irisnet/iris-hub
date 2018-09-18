package utils

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/keys"
)

func SendOrPrintTx(txCtx context.TxContext, cliCtx context.CLIContext, msgs []sdk.Msg) error {
	if cliCtx.GenerateOnly {
		return PrintUnsignedStdTx(txCtx, cliCtx, msgs)
	}
	// Build and sign the transaction, then broadcast to a Tendermint
	// node.
	cliCtx.PrintResponse = true

	txCtx, err := prepareTxContext(txCtx, cliCtx)
	if err != nil {
		return err
	}

	passphrase, err := keys.GetPassphrase(cliCtx.FromAddressName)
	if err != nil {
		return err
	}

	// build and sign the transaction
	txBytes, err := txCtx.BuildAndSign(cliCtx.FromAddressName, passphrase, msgs)
	if err != nil {
		return err
	}

	// broadcast to a Tendermint node
	return cliCtx.EnsureBroadcastTx(txBytes)
}

// PrintUnsignedStdTx builds an unsigned StdTx and prints it to os.Stdout.
func PrintUnsignedStdTx(txCtx context.TxContext, cliCtx context.CLIContext, msgs []sdk.Msg) (err error) {
	stdTx, err := buildUnsignedStdTx(txCtx, cliCtx, msgs)
	if err != nil {
		return
	}
	json, err := txCtx.Codec.MarshalJSON(stdTx)
	if err == nil {
		fmt.Printf("%s\n", json)
	}
	return
}

// buildUnsignedStdTx builds a StdTx as per the parameters passed in the
// contexts. Gas is automatically estimated if gas wanted is set to 0.
func buildUnsignedStdTx(txCtx context.TxContext, cliCtx context.CLIContext, msgs []sdk.Msg) (stdTx auth.StdTx, err error) {
	txCtx, err = prepareTxContext(txCtx, cliCtx)
	if err != nil {
		return
	}
	stdSignMsg, err := txCtx.Build(msgs)
	if err != nil {
		return
	}
	return auth.NewStdTx(stdSignMsg.Msgs, stdSignMsg.Fee, nil, stdSignMsg.Memo), nil
}

func prepareTxContext(txCtx context.TxContext, cliCtx context.CLIContext) (context.TxContext, error) {
	if err := cliCtx.EnsureAccountExists(); err != nil {
		return txCtx, err
	}

	from, err := cliCtx.GetFromAddress()
	if err != nil {
		return txCtx, err
	}

	// automatically doing a manual lookup.
	if txCtx.AccountNumber == 0 {
		accNum, err := cliCtx.GetAccountNumber(from)
		if err != nil {
			return txCtx, err
		}
		txCtx = txCtx.WithAccountNumber(accNum)
	}

	// automatically doing a manual lookup.
	if txCtx.Sequence == 0 {
		accSeq, err := cliCtx.GetAccountSequence(from)
		if err != nil {
			return txCtx, err
		}
		txCtx = txCtx.WithSequence(accSeq)
	}
	return txCtx, nil
}
