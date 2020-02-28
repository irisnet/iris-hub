package service

import (
	"github.com/irisnet/irishub/app/v3/service/internal/keeper"
	"github.com/irisnet/irishub/app/v3/service/internal/types"
)

// nolint
const (
	DefaultParamSpace     = types.DefaultParamSpace
	DefaultCodespace      = types.DefaultCodespace
	MsgRoute              = types.MsgRoute
	MetricsSubsystem      = types.MetricsSubsystem
	QueryDefinition       = types.QueryDefinition
	QueryBinding          = types.QueryBinding
	QueryBindings         = types.QueryBindings
	QueryWithdrawAddress  = types.QueryWithdrawAddress
	QueryRequest          = types.QueryRequest
	QueryRequests         = types.QueryRequests
	QueryResponse         = types.QueryResponse
	QueryRequestContext   = types.QueryRequestContext
	QueryRequestsByReqCtx = types.QueryRequestsByReqCtx
	QueryResponses        = types.QueryResponses
	QueryFees             = types.QueryFees
	RUNNING               = types.RUNNING
	PAUSED                = types.PAUSED
	COMPLETED             = types.COMPLETED
	BATCHRUNNING          = types.BATCHRUNNING
	BATCHCOMPLETED        = types.BATCHCOMPLETED
)

// nolint
var (
	// variables and functions aliases
	TagServiceName             = types.TagServiceName
	TagProvider                = types.TagProvider
	TagConsumer                = types.TagConsumer
	TagRequestContextID        = types.TagRequestContextID
	TagRequestID               = types.TagRequestID
	TagSlashedCoins            = types.TagSlashedCoins
	KeyTxSizeLimit             = types.KeyTxSizeLimit
	NewKeeper                  = keeper.NewKeeper
	NewQuerier                 = keeper.NewQuerier
	RegisterCodec              = types.RegisterCodec
	NewGenesisState            = types.NewGenesisState
	DefaultGenesisState        = types.DefaultGenesisState
	DefaultGenesisStateForTest = types.DefaultGenesisStateForTest
	ValidateGenesis            = types.ValidateGenesis
	DefaultParams              = types.DefaultParams
	PrometheusMetrics          = types.PrometheusMetrics
	NewServiceDefinition       = types.NewServiceDefinition
	NewServiceBinding          = types.NewServiceBinding
	NewRequest                 = types.NewRequest
	NewResponse                = types.NewResponse
	NewEarnedFees              = types.NewEarnedFees
	NewMsgDefineService        = types.NewMsgDefineService
	NewMsgBindService          = types.NewMsgBindService
	NewMsgUpdateServiceBinding = types.NewMsgUpdateServiceBinding
	NewMsgSetWithdrawAddress   = types.NewMsgSetWithdrawAddress
	NewMsgDisableService       = types.NewMsgDisableService
	NewMsgEnableService        = types.NewMsgEnableService
	NewMsgRefundServiceDeposit = types.NewMsgRefundServiceDeposit
	NewMsgRequestService       = types.NewMsgRequestService
	NewMsgRespondService       = types.NewMsgRespondService
	NewMsgPauseRequestContext  = types.NewMsgPauseRequestContext
	NewMsgStartRequestContext  = types.NewMsgStartRequestContext
	NewMsgKillRequestContext   = types.NewMsgKillRequestContext
	NewMsgUpdateRequestContext = types.NewMsgUpdateRequestContext
	NewMsgWithdrawEarnedFees   = types.NewMsgWithdrawEarnedFees
	NewMsgWithdrawTax          = types.NewMsgWithdrawTax
	ConvertRequestID           = types.ConvertRequestID
	RequestIDToString          = types.RequestIDToString
)

// nolint
type (
	Keeper                      = keeper.Keeper
	GenesisState                = types.GenesisState
	Params                      = types.Params
	ServiceDefinition           = types.ServiceDefinition
	ServiceBinding              = types.ServiceBinding
	ServiceBindings             = types.ServiceBindings
	RequestContext              = types.RequestContext
	Request                     = types.Request
	Requests                    = types.Requests
	CompactRequest              = types.CompactRequest
	Response                    = types.Response
	Responses                   = types.Responses
	EarnedFees                  = types.EarnedFees
	Metrics                     = types.Metrics
	MsgDefineService            = types.MsgDefineService
	MsgBindService              = types.MsgBindService
	MsgUpdateServiceBinding     = types.MsgUpdateServiceBinding
	MsgSetWithdrawAddress       = types.MsgSetWithdrawAddress
	MsgDisableService           = types.MsgDisableService
	MsgEnableService            = types.MsgEnableService
	MsgRefundServiceDeposit     = types.MsgRefundServiceDeposit
	MsgRequestService           = types.MsgRequestService
	MsgRespondService           = types.MsgRespondService
	MsgPauseRequestContext      = types.MsgPauseRequestContext
	MsgStartRequestContext      = types.MsgStartRequestContext
	MsgKillRequestContext       = types.MsgKillRequestContext
	MsgUpdateRequestContext     = types.MsgUpdateRequestContext
	MsgWithdrawEarnedFees       = types.MsgWithdrawEarnedFees
	MsgWithdrawTax              = types.MsgWithdrawTax
	QueryDefinitionParams       = types.QueryDefinitionParams
	QueryBindingParams          = types.QueryBindingParams
	QueryBindingsParams         = types.QueryBindingsParams
	QueryWithdrawAddressParams  = types.QueryWithdrawAddressParams
	QueryRequestParams          = types.QueryRequestParams
	QueryRequestsParams         = types.QueryRequestsParams
	QueryResponseParams         = types.QueryResponseParams
	QueryRequestContextParams   = types.QueryRequestContextParams
	QueryRequestsByReqCtxParams = types.QueryRequestsByReqCtxParams
	QueryResponsesParams        = types.QueryResponsesParams
	QueryFeesParams             = types.QueryFeesParams
)
