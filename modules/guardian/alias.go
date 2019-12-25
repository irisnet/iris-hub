// nolint
// autogenerated code using github.com/rigelrozanski/multitool
// aliases generated for the following subdirectories:
// ALIASGEN: github.com/irisnet/irishub/modules/guardian/internal/keeper
// ALIASGEN: github.com/irisnet/irishub/modules/guardian/internal/types
package guardian

import (
	"github.com/irisnet/irishub/modules/guardian/internal/keeper"
	"github.com/irisnet/irishub/modules/guardian/internal/types"
)

const (
	DefaultCodespace            = types.DefaultCodespace
	CodeInvalidOperator         = types.CodeInvalidOperator
	CodeProfilerExists          = types.CodeProfilerExists
	CodeProfilerNotExists       = types.CodeProfilerNotExists
	CodeTrusteeExists           = types.CodeTrusteeExists
	CodeTrusteeNotExists        = types.CodeTrusteeNotExists
	CodeInvalidDescription      = types.CodeInvalidDescription
	CodeDeleteGenesisProfiler   = types.CodeDeleteGenesisProfiler
	CodeDeleteGenesisTrustee    = types.CodeDeleteGenesisTrustee
	CodeInvalidGuardian         = types.CodeInvalidGuardian
	EventTypeAddProfiler        = types.EventTypeAddProfiler
	EventTypeAddTrustee         = types.EventTypeAddTrustee
	EventTypeDeleteProfiler     = types.EventTypeDeleteProfiler
	EventTypeDeleteTrustee      = types.EventTypeDeleteTrustee
	AttributeKeyProfilerAddress = types.AttributeKeyProfilerAddress
	AttributeKeyTrusteeAddress  = types.AttributeKeyTrusteeAddress
	AttributeKeyAddedBy         = types.AttributeKeyAddedBy
	AttributeKeyDeletedBy       = types.AttributeKeyDeletedBy
	AttributeValueCategory      = types.AttributeValueCategory
	ModuleName                  = types.ModuleName
	StoreKey                    = types.StoreKey
	RouterKey                   = types.RouterKey
	TypeMsgAddProfiler          = types.TypeMsgAddProfiler
	TypeMsgDeleteProfiler       = types.TypeMsgDeleteProfiler
	TypeMsgAddTrustee           = types.TypeMsgAddTrustee
	TypeMsgDeleteTrustee        = types.TypeMsgDeleteTrustee
	QuerierRoute                = types.QuerierRoute
	QueryProfilers              = types.QueryProfilers
	QueryTrustees               = types.QueryTrustees
	Genesis                     = types.Genesis
	Ordinary                    = types.Ordinary
)

var (
	// functions aliases
	NewKeeper                = keeper.NewKeeper
	NewQuerier               = keeper.NewQuerier
	RegisterCodec            = types.RegisterCodec
	ErrInvalidOperator       = types.ErrInvalidOperator
	ErrProfilerNotExists     = types.ErrProfilerNotExists
	ErrDeleteGenesisProfiler = types.ErrDeleteGenesisProfiler
	ErrProfilerExists        = types.ErrProfilerExists
	ErrTrusteeExists         = types.ErrTrusteeExists
	ErrTrusteeNotExists      = types.ErrTrusteeNotExists
	ErrDeleteGenesisTrustee  = types.ErrDeleteGenesisTrustee
	ErrInvalidDescription    = types.ErrInvalidDescription
	NewGenesisState          = types.NewGenesisState
	DefaultGenesisState      = types.DefaultGenesisState
	ProfilerKey              = types.ProfilerKey
	TrusteeKey               = types.TrusteeKey
	GetProfilerKey           = types.GetProfilerKey
	GetTrusteeKey            = types.GetTrusteeKey
	GetProfilersSubspaceKey  = types.GetProfilersSubspaceKey
	GetTrusteesSubspaceKey   = types.GetTrusteesSubspaceKey
	NewMsgAddProfiler        = types.NewMsgAddProfiler
	NewMsgDeleteProfiler     = types.NewMsgDeleteProfiler
	NewMsgAddTrustee         = types.NewMsgAddTrustee
	NewMsgDeleteTrustee      = types.NewMsgDeleteTrustee
	NewGuardian              = types.NewGuardian
	AccountTypeFromString    = types.AccountTypeFromString

	// variable aliases
	ModuleCdc = types.ModuleCdc
)

type (
	Keeper            = keeper.Keeper
	GenesisState      = types.GenesisState
	MsgAddProfiler    = types.MsgAddProfiler
	MsgDeleteProfiler = types.MsgDeleteProfiler
	MsgAddTrustee     = types.MsgAddTrustee
	MsgDeleteTrustee  = types.MsgDeleteTrustee
	AddGuardian       = types.AddGuardian
	DeleteGuardian    = types.DeleteGuardian
	Guardian          = types.Guardian
	Profilers         = types.Profilers
	Trustees          = types.Trustees
	AccountType       = types.AccountType
)
