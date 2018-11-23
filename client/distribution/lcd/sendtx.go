package lcd

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
)

type setWithdrawAddressBody struct {
	WithdrawAddress sdk.AccAddress `json:"withdraw_address"`
	BaseTx          context.BaseTx `json:"base_tx"`
}

// SetWithdrawAddressHandlerFn - http request handler to set withdraw address
// nolint: gocyclo
func SetWithdrawAddressHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Init context and read request parameters
		vars := mux.Vars(r)
		bech32addr := vars["delegatorAddr"]
		delegatorAddress, err := sdk.AccAddressFromBech32(bech32addr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx = utils.InitReqCliCtx(cliCtx, r)
		var m setWithdrawAddressBody
		err = utils.ReadPostBody(w, r, cdc, &m)
		if err != nil {
			return
		}
		baseReq := m.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w, cliCtx) {
			return
		}
		// Build message
		msg := types.NewMsgSetWithdrawAddress(delegatorAddress, m.WithdrawAddress)
		// Broadcast or return unsigned transaction
		utils.SendOrReturnUnsignedTx(w, cliCtx, m.BaseTx, []sdk.Msg{msg})
	}
}

type withdrawRewardsBody struct {
	ValidatorAddress sdk.ValAddress `json:"validator_address"`
	IsValidator      bool           `json:"is_validator"`
	BaseTx           context.BaseTx `json:"base_tx"`
}

func WithdrawRewardsHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Init context and read request parameters
		vars := mux.Vars(r)
		bech32addr := vars["delegatorAddr"]
		delegatorAddress, err := sdk.AccAddressFromBech32(bech32addr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx = utils.InitReqCliCtx(cliCtx, r)
		var m withdrawRewardsBody
		err = utils.ReadPostBody(w, r, cdc, &m)
		if err != nil {
			return
		}
		baseReq := m.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w, cliCtx) {
			return
		}
		// Build message
		onlyFromVal := m.ValidatorAddress
		isVal := m.IsValidator
		if onlyFromVal != nil && isVal {
			utils.WriteErrorResponse(w, http.StatusBadRequest, "if is_validator is true, validator_address should not be specified")
			return
		}

		var msg sdk.Msg
		switch {
		case isVal:
			valAddr := sdk.ValAddress(delegatorAddress)
			msg = types.NewMsgWithdrawValidatorRewardsAll(valAddr)
		case onlyFromVal != nil:
			msg = types.NewMsgWithdrawDelegatorReward(delegatorAddress, m.ValidatorAddress)
		default:
			msg = types.NewMsgWithdrawDelegatorRewardsAll(delegatorAddress)
		}
		// Broadcast or return unsigned transaction
		utils.SendOrReturnUnsignedTx(w, cliCtx, m.BaseTx, []sdk.Msg{msg})
	}
}
