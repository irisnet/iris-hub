package utils

import (
	"net/http"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/client/context"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

const (
	Async = "async"
	GenerateOnly  = "generate-only"
)

type BaseTx struct {
	LocalAccountName string    `json:"name"`
	Password         string    `json:"password"`
	ChainID          string    `json:"chain_id"`
	AccountNumber    int64     `json:"account_number"`
	Sequence         int64     `json:"sequence"`
	Gas              int64     `json:"gas"`
	Fees             string    `json:"fee"`
}

// WriteErrorResponse prepares and writes a HTTP error
// given a status code and an error message.
func WriteErrorResponse(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	w.Write([]byte(msg))
}

func SendOrReturnUnsignedTx(w http.ResponseWriter, cliCtx context.CLIContext, txCtx context.TxContext, baseTx BaseTx, msgs []sdk.Msg)  {

	if cliCtx.GenerateOnly {
		WriteGenerateStdTxResponse(w, txCtx, msgs)
		return
	}

	txBytes, err := txCtx.BuildAndSign(baseTx.LocalAccountName, baseTx.Password, msgs)
	if err != nil {
		WriteErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	var res interface{}
	if cliCtx.Async {
		res, err = cliCtx.BroadcastTxAsync(txBytes)
	} else {
		res, err = cliCtx.BroadcastTx(txBytes)
	}

	output, err := txCtx.Codec.MarshalJSONIndent(res, "", "  ")
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Write(output)
}

// WriteGenerateStdTxResponse writes response for the generate_only mode.
func WriteGenerateStdTxResponse(w http.ResponseWriter, txCtx context.TxContext, msgs []sdk.Msg) {
	stdMsg, err := txCtx.Build(msgs)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	output, err := txCtx.Codec.MarshalJSON(auth.NewStdTx(stdMsg.Msgs, stdMsg.Fee, nil, stdMsg.Memo))
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Write(output)
	return
}
