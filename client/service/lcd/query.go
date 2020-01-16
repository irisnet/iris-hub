package lcd

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/app/v3/service"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	// query a service definition
	r.HandleFunc(
		fmt.Sprintf("/service/definitions/{%s}", ServiceName),
		queryDefinitionHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// query a service binding
	r.HandleFunc(
		fmt.Sprintf("/service/bindings/{%s}/{%s}", ServiceName, Provider),
		queryBindingHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// query all bindings of a service definition
	r.HandleFunc(
		fmt.Sprintf("/service/bindings/{%s}", ServiceName),
		queryBindingsHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// get all active requests of a binding
	r.HandleFunc(
		fmt.Sprintf("/service/requests/{%s}/{%s}/{%s}/{%s}", DefChainID, ServiceName, BindChainID, Provider),
		requestsHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// get a single response
	r.HandleFunc(
		fmt.Sprintf("/service/responses/{%s}/{%s}", ReqChainID, ReqID),
		responseGetHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// get return fee and incoming fee of a account
	r.HandleFunc(
		fmt.Sprintf("/service/fees/{%s}", Address),
		feesHandlerFn(cliCtx, cdc),
	).Methods("GET")
}

func queryDefinitionHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		serviceName := vars[ServiceName]

		params := service.QueryDefinitionParams{
			ServiceName: serviceName,
		}

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", protocol.ServiceRoute, service.QueryDefinition)
		res, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cliCtx.Codec, res, cliCtx.Indent)
	}
}

func queryBindingHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		serviceName := vars[ServiceName]
		providerStr := vars[Provider]

		provider, err := sdk.AccAddressFromBech32(providerStr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := service.QueryBindingParams{
			ServiceName: serviceName,
			Provider:    provider,
		}

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", protocol.ServiceRoute, service.QueryBinding)
		res, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cliCtx.Codec, res, cliCtx.Indent)
	}
}

func queryBindingsHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		serviceName := vars[ServiceName]

		params := service.QueryBindingsParams{
			ServiceName: serviceName,
		}

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", protocol.ServiceRoute, service.QueryBindings)
		res, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cliCtx.Codec, res, cliCtx.Indent)
	}
}

func requestsHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		defChainID := vars[DefChainID]
		serviceName := vars[ServiceName]
		bindChainID := vars[BindChainID]
		bechProviderAddr := vars[Provider]

		provider, err := sdk.AccAddressFromBech32(bechProviderAddr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := service.QueryRequestsParams{
			DefChainID:  defChainID,
			ServiceName: serviceName,
			BindChainID: bindChainID,
			Provider:    provider,
		}

		bz, err := cdc.MarshalJSON(params)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", protocol.ServiceRoute, service.QueryRequests)
		res, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func responseGetHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		reqChainID := vars[ReqChainID]
		reqID := vars[ReqID]

		params := service.QueryResponseParams{
			ReqChainID: reqChainID,
			RequestID:  reqID,
		}

		bz, err := cdc.MarshalJSON(params)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", protocol.ServiceRoute, service.QueryResponse)
		res, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func feesHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bechAddress := vars[Address]

		address, err := sdk.AccAddressFromBech32(bechAddress)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := service.QueryFeesParams{
			Address: address,
		}

		bz, err := cdc.MarshalJSON(params)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", protocol.ServiceRoute, service.QueryFees)
		res, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}
