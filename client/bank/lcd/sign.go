package lcd

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/keyerror"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
)

// SignBody defines the properties of a sign request's body.
type SignBody struct {
	Tx            auth.StdTx `json:"tx"`
	Name          string     `json:"name"`
	Password      string     `json:"password"`
	ChainID       string     `json:"chain_id"`
	AccountNumber int64      `json:"account_number"`
	Sequence      int64      `json:"sequence"`
	AppendSig     bool       `json:"append_sig"`
}

// nolint: unparam
// SignTxRequestHandlerFn sign tx REST handler
func SignTxRequestHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m SignBody
		err := utils.ReadPostBody(w, r, cdc, &m)
		if err != nil {
			return
		}

		txCtx := context.TxContext{
			Codec:         cliCtx.Codec,
			ChainID:       m.ChainID,
			AccountNumber: m.AccountNumber,
			Sequence:      m.Sequence,
		}

		signedTx, err := txCtx.SignStdTx(m.Name, m.Password, m.Tx, m.AppendSig)
		if keyerror.IsErrKeyNotFound(err) {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		} else if keyerror.IsErrWrongPassword(err) {
			utils.WriteErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		} else if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cdc, signedTx, cliCtx.Indent)
	}
}
