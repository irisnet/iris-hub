package lcd

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/app/v1/asset"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
)

// queryGateway queries a gateway of the given moniker from the specified endpoint
func queryGateway(cliCtx context.CLIContext, cdc *codec.Codec, endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		moniker := vars["Moniker"]
		if len(moniker) == 0 || uint32(len(moniker)) > asset.MaximumGatewayMonikerSize {
			utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("the length of the moniker must be (0,%d]", asset.MaximumGatewayMonikerSize))
			return
		}

		params := asset.QueryGatewayParams{
			Moniker: moniker,
		}

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/%s", protocol.AssetRoute, asset.QueryGateway), bz)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cliCtx.Codec, res, cliCtx.Indent)
	}
}

// queryGateways queries all gateways of an owner from the specified endpoint
func queryGateways(cliCtx context.CLIContext, cdc *codec.Codec, endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		ownerStr := vars["owner"]
		owner, err := sdk.AccAddressFromBech32(ownerStr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := asset.QueryGatewaysParams{
			Owner: owner,
		}

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/%s", protocol.AssetRoute, asset.QueryGateways), bz)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cliCtx.Codec, res, cliCtx.Indent)
	}
}
