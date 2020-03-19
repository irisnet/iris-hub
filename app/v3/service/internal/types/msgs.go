package types

import (
	"encoding/json"
	"fmt"
	"regexp"

	sdk "github.com/irisnet/irishub/types"
)

const (
	MsgRoute = "service" // route for service msgs

	TypeMsgDefineService        = "define_service"         // type for MsgDefineService
	TypeMsgBindService          = "bind_service"           // type for MsgBindService
	TypeMsgUpdateServiceBinding = "update_service_binding" // type for MsgUpdateServiceBinding
	TypeMsgSetWithdrawAddress   = "set_withdraw_address"   // type for MsgSetWithdrawAddress
	TypeMsgDisableService       = "disable_service"        // type for MsgDisableService
	TypeMsgEnableService        = "enable_service"         // type for MsgEnableService
	TypeMsgRefundServiceDeposit = "refund_service_deposit" // type for MsgRefundServiceDeposit
	TypeMsgRequestService       = "request_service"        // type for MsgRequestService
	TypeMsgRespondService       = "respond_service"        // type for MsgRespondService
	TypeMsgPauseRequestContext  = "pause_request_context"  // type for MsgPauseRequestContext
	TypeMsgStartRequestContext  = "start_request_context"  // type for MsgStartRequestContext
	TypeMsgKillRequestContext   = "kill_request_context"   // type for MsgKillRequestContext
	TypeMsgUpdateRequestContext = "update_request_context" // type for MsgUpdateRequestContext
	TypeMsgWithdrawEarnedFees   = "withdraw_earned_fees"   // type for MsgWithdrawEarnedFees
	TypeMsgWithdrawTax          = "withdraw_tax"           // type for MsgWithdrawTax

	MaxNameLength        = 70  // max length of the service name
	MaxDescriptionLength = 280 // max length of the service and author description
	MaxTagsNum           = 10  // max total number of the tags
	MaxTagLength         = 70  // max length of the tag

	MaxProvidersNum     = 10 // max total number of the providers to request
	RequestContextIDLen = 32 // length of the request context ID in bytes
)

// the service name only accepts alphanumeric characters, _ and -
var reServiceName = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)

var (
	_ sdk.Msg = MsgDefineService{}
	_ sdk.Msg = MsgBindService{}
	_ sdk.Msg = MsgUpdateServiceBinding{}
	_ sdk.Msg = MsgSetWithdrawAddress{}
	_ sdk.Msg = MsgDisableService{}
	_ sdk.Msg = MsgEnableService{}
	_ sdk.Msg = MsgRefundServiceDeposit{}
	_ sdk.Msg = MsgRequestService{}
	_ sdk.Msg = MsgRespondService{}
	_ sdk.Msg = MsgPauseRequestContext{}
	_ sdk.Msg = MsgStartRequestContext{}
	_ sdk.Msg = MsgKillRequestContext{}
	_ sdk.Msg = MsgUpdateRequestContext{}
	_ sdk.Msg = MsgWithdrawEarnedFees{}
	_ sdk.Msg = MsgWithdrawTax{}
)

//______________________________________________________________________

// MsgDefineService defines a message to define a service
type MsgDefineService struct {
	Name              string         `json:"name"`
	Description       string         `json:"description"`
	Tags              []string       `json:"tags"`
	Author            sdk.AccAddress `json:"author"`
	AuthorDescription string         `json:"author_description"`
	Schemas           string         `json:"schemas"`
}

// NewMsgDefineService creates a new MsgDefineService instance
func NewMsgDefineService(name, description string, tags []string, author sdk.AccAddress, authorDescription, schemas string) MsgDefineService {
	return MsgDefineService{
		Name:              name,
		Description:       description,
		Tags:              tags,
		Author:            author,
		AuthorDescription: authorDescription,
		Schemas:           schemas,
	}
}

// Route implements Msg
func (msg MsgDefineService) Route() string { return MsgRoute }

// Type implements Msg
func (msg MsgDefineService) Type() string { return TypeMsgDefineService }

// GetSignBytes implements Msg
func (msg MsgDefineService) GetSignBytes() []byte {
	if len(msg.Tags) == 0 {
		msg.Tags = nil
	}

	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg
func (msg MsgDefineService) ValidateBasic() sdk.Error {
	if len(msg.Author) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "author missing")
	}

	if !validServiceName(msg.Name) {
		return ErrInvalidServiceName(DefaultCodespace, msg.Name)
	}

	if err := ensureServiceDefLength(msg); err != nil {
		return err
	}

	if sdk.HasDuplicate(msg.Tags) {
		return ErrDuplicateTags(DefaultCodespace)
	}

	if len(msg.Schemas) == 0 {
		return ErrInvalidSchemas(DefaultCodespace, "schemas missing")
	}

	return ValidateServiceSchemas(msg.Schemas)
}

// GetSigners implements Msg
func (msg MsgDefineService) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Author}
}

//______________________________________________________________________

// MsgBindService defines a message to bind a service
type MsgBindService struct {
	ServiceName string         `json:"service_name"`
	Provider    sdk.AccAddress `json:"provider"`
	Deposit     sdk.Coins      `json:"deposit"`
	Pricing     string         `json:"pricing"`
}

// NewMsgBindService creates a new MsgBindService instance
func NewMsgBindService(serviceName string, provider sdk.AccAddress, deposit sdk.Coins, pricing string) MsgBindService {
	return MsgBindService{
		ServiceName: serviceName,
		Provider:    provider,
		Deposit:     deposit,
		Pricing:     pricing,
	}
}

// Route implements Msg.
func (msg MsgBindService) Route() string { return MsgRoute }

// Type implements Msg.
func (msg MsgBindService) Type() string { return TypeMsgBindService }

// GetSignBytes implements Msg.
func (msg MsgBindService) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgBindService) ValidateBasic() sdk.Error {
	if len(msg.Provider) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "provider missing")
	}

	if err := ValidateServiceName(msg.ServiceName); err != nil {
		return err
	}

	if !validServiceCoins(msg.Deposit) {
		return ErrInvalidDeposit(DefaultCodespace, fmt.Sprintf("invalid deposit: %s", msg.Deposit))
	}

	if len(msg.Pricing) == 0 {
		return ErrInvalidPricing(DefaultCodespace, "pricing missing")
	}

	return validatePricing(msg.Pricing)
}

// GetSigners implements Msg.
func (msg MsgBindService) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgUpdateServiceBinding defines a message to update a service binding
type MsgUpdateServiceBinding struct {
	ServiceName string         `json:"service_name"`
	Provider    sdk.AccAddress `json:"provider"`
	Deposit     sdk.Coins      `json:"deposit"`
	Pricing     string         `json:"pricing"`
}

// NewMsgUpdateServiceBinding creates a new MsgUpdateServiceBinding instance
func NewMsgUpdateServiceBinding(serviceName string, provider sdk.AccAddress, deposit sdk.Coins, pricing string) MsgUpdateServiceBinding {
	return MsgUpdateServiceBinding{
		ServiceName: serviceName,
		Provider:    provider,
		Deposit:     deposit,
		Pricing:     pricing,
	}
}

// Route implements Msg.
func (msg MsgUpdateServiceBinding) Route() string { return MsgRoute }

// Type implements Msg.
func (msg MsgUpdateServiceBinding) Type() string { return TypeMsgUpdateServiceBinding }

// GetSignBytes implements Msg.
func (msg MsgUpdateServiceBinding) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgUpdateServiceBinding) ValidateBasic() sdk.Error {
	if len(msg.Provider) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "provider missing")
	}

	if err := ValidateServiceName(msg.ServiceName); err != nil {
		return err
	}

	if !msg.Deposit.Empty() && !validServiceCoins(msg.Deposit) {
		return ErrInvalidDeposit(DefaultCodespace, fmt.Sprintf("invalid deposit: %s", msg.Deposit))
	}

	if len(msg.Pricing) != 0 {
		return validatePricing(msg.Pricing)
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgUpdateServiceBinding) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgSetWithdrawAddress defines a message to set the withdrawal address for a provider
type MsgSetWithdrawAddress struct {
	Provider        sdk.AccAddress `json:"provider"`
	WithdrawAddress sdk.AccAddress `json:"withdraw_address"`
}

// NewMsgSetWithdrawAddress creates a new MsgSetWithdrawAddress instance
func NewMsgSetWithdrawAddress(provider sdk.AccAddress, withdrawAddr sdk.AccAddress) MsgSetWithdrawAddress {
	return MsgSetWithdrawAddress{
		Provider:        provider,
		WithdrawAddress: withdrawAddr,
	}
}

// Route implements Msg.
func (msg MsgSetWithdrawAddress) Route() string { return MsgRoute }

// Type implements Msg.
func (msg MsgSetWithdrawAddress) Type() string { return TypeMsgSetWithdrawAddress }

// GetSignBytes implements Msg.
func (msg MsgSetWithdrawAddress) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgSetWithdrawAddress) ValidateBasic() sdk.Error {
	if len(msg.Provider) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "provider missing")
	}

	if len(msg.WithdrawAddress) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "withdrawal address missing")
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgSetWithdrawAddress) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgDisableService defines a message to disable a service binding
type MsgDisableService struct {
	ServiceName string         `json:"service_name"`
	Provider    sdk.AccAddress `json:"provider"`
}

// NewMsgDisableService creates a new MsgDisableService instance
func NewMsgDisableService(serviceName string, provider sdk.AccAddress) MsgDisableService {
	return MsgDisableService{
		ServiceName: serviceName,
		Provider:    provider,
	}
}

// Route implements Msg.
func (msg MsgDisableService) Route() string { return MsgRoute }

// Type implements Msg.
func (msg MsgDisableService) Type() string { return TypeMsgDisableService }

// GetSignBytes implements Msg.
func (msg MsgDisableService) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgDisableService) ValidateBasic() sdk.Error {
	if len(msg.Provider) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "provider missing")
	}

	return ValidateServiceName(msg.ServiceName)
}

// GetSigners implements Msg.
func (msg MsgDisableService) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgEnableService defines a message to enable a service binding
type MsgEnableService struct {
	ServiceName string         `json:"service_name"`
	Provider    sdk.AccAddress `json:"provider"`
	Deposit     sdk.Coins      `json:"deposit"`
}

// NewMsgEnableService creates a new MsgEnableService instance
func NewMsgEnableService(serviceName string, provider sdk.AccAddress, deposit sdk.Coins) MsgEnableService {
	return MsgEnableService{
		ServiceName: serviceName,
		Provider:    provider,
		Deposit:     deposit,
	}
}

// Route implements Msg.
func (msg MsgEnableService) Route() string { return MsgRoute }

// Type implements Msg.
func (msg MsgEnableService) Type() string { return TypeMsgEnableService }

// GetSignBytes implements Msg.
func (msg MsgEnableService) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgEnableService) ValidateBasic() sdk.Error {
	if len(msg.Provider) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "provider missing")
	}

	if err := ValidateServiceName(msg.ServiceName); err != nil {
		return err
	}

	if !msg.Deposit.Empty() && !validServiceCoins(msg.Deposit) {
		return ErrInvalidDeposit(DefaultCodespace, fmt.Sprintf("invalid deposit: %s", msg.Deposit))
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgEnableService) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgRefundServiceDeposit defines a message to refund deposit from a service binding
type MsgRefundServiceDeposit struct {
	ServiceName string         `json:"service_name"`
	Provider    sdk.AccAddress `json:"provider"`
}

// NewMsgRefundServiceDeposit creates a new MsgRefundServiceDeposit instance
func NewMsgRefundServiceDeposit(serviceName string, provider sdk.AccAddress) MsgRefundServiceDeposit {
	return MsgRefundServiceDeposit{
		ServiceName: serviceName,
		Provider:    provider,
	}
}

// Route implements Msg.
func (msg MsgRefundServiceDeposit) Route() string { return MsgRoute }

// Type implements Msg.
func (msg MsgRefundServiceDeposit) Type() string { return TypeMsgRefundServiceDeposit }

// GetSignBytes implements Msg.
func (msg MsgRefundServiceDeposit) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgRefundServiceDeposit) ValidateBasic() sdk.Error {
	if len(msg.Provider) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "provider missing")
	}

	return ValidateServiceName(msg.ServiceName)
}

// GetSigners implements Msg.
func (msg MsgRefundServiceDeposit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgRequestService defines a message to request a service
type MsgRequestService struct {
	ServiceName       string           `json:"service_name"`
	Providers         []sdk.AccAddress `json:"providers"`
	Consumer          sdk.AccAddress   `json:"consumer"`
	Input             string           `json:"input"`
	ServiceFeeCap     sdk.Coins        `json:"service_fee_cap"`
	Timeout           int64            `json:"timeout"`
	SuperMode         bool             `json:"super_mode"`
	Repeated          bool             `json:"repeated"`
	RepeatedFrequency uint64           `json:"repeated_frequency"`
	RepeatedTotal     int64            `json:"repeated_total"`
}

// NewMsgRequestService creates a new MsgRequestService instance
func NewMsgRequestService(
	serviceName string,
	providers []sdk.AccAddress,
	consumer sdk.AccAddress,
	input string,
	serviceFeeCap sdk.Coins,
	timeout int64,
	superMode bool,
	repeated bool,
	repeatedFrequency uint64,
	repeatedTotal int64,
) MsgRequestService {
	return MsgRequestService{
		ServiceName:       serviceName,
		Providers:         providers,
		Consumer:          consumer,
		Input:             input,
		ServiceFeeCap:     serviceFeeCap,
		Timeout:           timeout,
		SuperMode:         superMode,
		Repeated:          repeated,
		RepeatedFrequency: repeatedFrequency,
		RepeatedTotal:     repeatedTotal,
	}
}

// Route implements Msg.
func (msg MsgRequestService) Route() string { return MsgRoute }

// Type implements Msg.
func (msg MsgRequestService) Type() string { return TypeMsgRequestService }

// GetSignBytes implements Msg.
func (msg MsgRequestService) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgRequestService) ValidateBasic() sdk.Error {
	if len(msg.Consumer) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "consumer missing")
	}

	return ValidateRequest(
		msg.ServiceName, msg.ServiceFeeCap, msg.Providers, msg.Input, msg.Timeout,
		msg.Repeated, msg.RepeatedFrequency, msg.RepeatedTotal,
	)
}

// GetSigners implements Msg.
func (msg MsgRequestService) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Consumer}
}

//______________________________________________________________________

// MsgRespondService defines a message to respond to a service request
type MsgRespondService struct {
	RequestID string         `json:"request_id"`
	Provider  sdk.AccAddress `json:"provider"`
	Result    string         `json:"result"`
	Output    string         `json:"output"`
}

// NewMsgRespondService creates a new MsgRespondService instance
func NewMsgRespondService(
	requestID string,
	provider sdk.AccAddress,
	result,
	output string,
) MsgRespondService {
	return MsgRespondService{
		RequestID: requestID,
		Provider:  provider,
		Result:    result,
		Output:    output,
	}
}

// Route implements Msg.
func (msg MsgRespondService) Route() string { return MsgRoute }

// Type implements Msg.
func (msg MsgRespondService) Type() string { return TypeMsgRespondService }

// GetSignBytes implements Msg.
func (msg MsgRespondService) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgRespondService) ValidateBasic() sdk.Error {
	if len(msg.Provider) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "provider missing")
	}

	if _, err := ConvertRequestID(msg.RequestID); err != nil {
		return ErrInvalidRequestID(DefaultCodespace, fmt.Sprintf("failed to parse %s", msg.RequestID))
	}

	if len(msg.Result) == 0 {
		return ErrInvalidResponseResult(DefaultCodespace, "result missing")
	}

	if err := ValidateResponseResult(msg.Result); err != nil {
		return err
	}

	result, err := ParseResult(msg.Result)
	if err != nil {
		return err
	}

	if result.Code == 200 && len(msg.Output) == 0 {
		return ErrInvalidResponse(DefaultCodespace, "output must be specified when the result code is 200")
	}

	if result.Code != 200 && len(msg.Output) != 0 {
		return ErrInvalidResponse(DefaultCodespace, "output should not be specified when the result code is not 200")
	}

	if len(msg.Output) > 0 {
		if !json.Valid([]byte(msg.Output)) {
			return ErrInvalidResponseOutput(DefaultCodespace, "output is not valid JSON")
		}
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgRespondService) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgPauseRequestContext defines a message to suspend a request context
type MsgPauseRequestContext struct {
	RequestContextID []byte         `json:"request_context_id"`
	Consumer         sdk.AccAddress `json:"consumer"`
}

// NewMsgPauseRequestContext creates a new MsgPauseRequestContext instance
func NewMsgPauseRequestContext(requestContextID []byte, consumer sdk.AccAddress) MsgPauseRequestContext {
	return MsgPauseRequestContext{
		RequestContextID: requestContextID,
		Consumer:         consumer,
	}
}

// Route implements Msg.
func (msg MsgPauseRequestContext) Route() string { return MsgRoute }

// Type implements Msg.
func (msg MsgPauseRequestContext) Type() string { return TypeMsgPauseRequestContext }

// GetSignBytes implements Msg.
func (msg MsgPauseRequestContext) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgPauseRequestContext) ValidateBasic() sdk.Error {
	if len(msg.Consumer) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "consumer missing")
	}

	if len(msg.RequestContextID) != RequestContextIDLen {
		return ErrInvalidRequestContextID(DefaultCodespace, fmt.Sprintf("length of the request context ID must be %d in bytes", RequestContextIDLen))
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgPauseRequestContext) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Consumer}
}

//______________________________________________________________________

// MsgStartRequestContext defines a message to resume a request context
type MsgStartRequestContext struct {
	RequestContextID []byte         `json:"request_context_id"`
	Consumer         sdk.AccAddress `json:"consumer"`
}

// NewMsgStartRequestContext creates a new MsgStartRequestContext instance
func NewMsgStartRequestContext(requestContextID []byte, consumer sdk.AccAddress) MsgStartRequestContext {
	return MsgStartRequestContext{
		RequestContextID: requestContextID,
		Consumer:         consumer,
	}
}

// Route implements Msg.
func (msg MsgStartRequestContext) Route() string { return MsgRoute }

// Type implements Msg.
func (msg MsgStartRequestContext) Type() string { return TypeMsgStartRequestContext }

// GetSignBytes implements Msg.
func (msg MsgStartRequestContext) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgStartRequestContext) ValidateBasic() sdk.Error {
	if len(msg.Consumer) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "consumer missing")
	}

	if len(msg.RequestContextID) != RequestContextIDLen {
		return ErrInvalidRequestContextID(DefaultCodespace, fmt.Sprintf("length of the request context ID must be %d in bytes", RequestContextIDLen))
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgStartRequestContext) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Consumer}
}

//______________________________________________________________________

// MsgKillRequestContext defines a message to terminate a request context
type MsgKillRequestContext struct {
	RequestContextID []byte         `json:"request_context_id"`
	Consumer         sdk.AccAddress `json:"consumer"`
}

// NewMsgKillRequestContext creates a new MsgKillRequestContext instance
func NewMsgKillRequestContext(requestContextID []byte, consumer sdk.AccAddress) MsgKillRequestContext {
	return MsgKillRequestContext{
		RequestContextID: requestContextID,
		Consumer:         consumer,
	}
}

// Route implements Msg.
func (msg MsgKillRequestContext) Route() string { return MsgRoute }

// Type implements Msg.
func (msg MsgKillRequestContext) Type() string { return TypeMsgKillRequestContext }

// GetSignBytes implements Msg.
func (msg MsgKillRequestContext) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgKillRequestContext) ValidateBasic() sdk.Error {
	if len(msg.Consumer) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "consumer missing")
	}

	if len(msg.RequestContextID) != RequestContextIDLen {
		return ErrInvalidRequestContextID(DefaultCodespace, fmt.Sprintf("length of the request context ID must be %d in bytes", RequestContextIDLen))
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgKillRequestContext) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Consumer}
}

//______________________________________________________________________

// MsgUpdateRequestContext defines a message to update a request context
type MsgUpdateRequestContext struct {
	RequestContextID  []byte           `json:"request_context_id"`
	Providers         []sdk.AccAddress `json:"providers"`
	ServiceFeeCap     sdk.Coins        `json:"service_fee_cap"`
	Timeout           int64            `json:"timeout"`
	RepeatedFrequency uint64           `json:"repeated_frequency"`
	RepeatedTotal     int64            `json:"repeated_total"`
	Consumer          sdk.AccAddress   `json:"consumer"`
}

// NewMsgUpdateRequestContext creates a new MsgUpdateRequestContext instance
func NewMsgUpdateRequestContext(
	requestContextID []byte,
	providers []sdk.AccAddress,
	serviceFeeCap sdk.Coins,
	timeout int64,
	repeatedFrequency uint64,
	repeatedTotal int64,
	consumer sdk.AccAddress,
) MsgUpdateRequestContext {
	return MsgUpdateRequestContext{
		RequestContextID:  requestContextID,
		Providers:         providers,
		ServiceFeeCap:     serviceFeeCap,
		Timeout:           timeout,
		RepeatedFrequency: repeatedFrequency,
		RepeatedTotal:     repeatedTotal,
		Consumer:          consumer,
	}
}

// Route implements Msg.
func (msg MsgUpdateRequestContext) Route() string { return MsgRoute }

// Type implements Msg.
func (msg MsgUpdateRequestContext) Type() string { return TypeMsgUpdateRequestContext }

// GetSignBytes implements Msg.
func (msg MsgUpdateRequestContext) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgUpdateRequestContext) ValidateBasic() sdk.Error {
	if len(msg.Consumer) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "consumer missing")
	}

	if len(msg.RequestContextID) != RequestContextIDLen {
		return ErrInvalidRequestContextID(DefaultCodespace, fmt.Sprintf("length of the request context ID must be %d in bytes", RequestContextIDLen))
	}

	return ValidateRequestContextUpdating(msg.Providers, msg.ServiceFeeCap, msg.Timeout, msg.RepeatedFrequency, msg.RepeatedTotal)
}

// GetSigners implements Msg.
func (msg MsgUpdateRequestContext) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Consumer}
}

//______________________________________________________________________

// MsgWithdrawEarnedFees defines a message to withdraw the fees earned by the provider
type MsgWithdrawEarnedFees struct {
	Provider sdk.AccAddress `json:"provider"`
}

// NewMsgWithdrawEarnedFees creates a new MsgWithdrawEarnedFees instance
func NewMsgWithdrawEarnedFees(provider sdk.AccAddress) MsgWithdrawEarnedFees {
	return MsgWithdrawEarnedFees{
		Provider: provider,
	}
}

// Route implements Msg.
func (msg MsgWithdrawEarnedFees) Route() string { return MsgRoute }

// Type implements Msg.
func (msg MsgWithdrawEarnedFees) Type() string { return TypeMsgWithdrawEarnedFees }

// GetSignBytes implements Msg.
func (msg MsgWithdrawEarnedFees) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgWithdrawEarnedFees) ValidateBasic() sdk.Error {
	if len(msg.Provider) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "provider missing")
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgWithdrawEarnedFees) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgWithdrawTax defines a message to withdraw the service tax
type MsgWithdrawTax struct {
	Trustee     sdk.AccAddress `json:"trustee"`
	DestAddress sdk.AccAddress `json:"dest_address"`
	Amount      sdk.Coins      `json:"amount"`
}

// NewMsgWithdrawTax creates a new MsgWithdrawTax instance
func NewMsgWithdrawTax(trustee, destAddress sdk.AccAddress, amount sdk.Coins) MsgWithdrawTax {
	return MsgWithdrawTax{
		Trustee:     trustee,
		DestAddress: destAddress,
		Amount:      amount,
	}
}

// Route implements Msg.
func (msg MsgWithdrawTax) Route() string { return MsgRoute }

// Type implements Msg.
func (msg MsgWithdrawTax) Type() string { return TypeMsgWithdrawTax }

// GetSignBytes implements Msg.
func (msg MsgWithdrawTax) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgWithdrawTax) ValidateBasic() sdk.Error {
	if len(msg.Trustee) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "trustee missing")
	}

	if len(msg.DestAddress) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "destination address missing")
	}

	if !validServiceCoins(msg.Amount) {
		return sdk.ErrInvalidCoins(fmt.Sprintf("invalid withdrawal amount: %s", msg.Amount))
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgWithdrawTax) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Trustee}
}

//______________________________________________________________________

// ValidateServiceName validates the service name
func ValidateServiceName(name string) sdk.Error {
	if !validServiceName(name) {
		return ErrInvalidServiceName(DefaultCodespace, name)
	}

	return ensureServiceNameLength(name)
}

func validServiceName(name string) bool {
	return reServiceName.MatchString(name)
}

func ensureServiceNameLength(name string) sdk.Error {
	if len(name) > MaxNameLength {
		return ErrInvalidServiceName(DefaultCodespace, name)
	}

	return nil
}

func ensureServiceDefLength(msg MsgDefineService) sdk.Error {
	if err := ensureServiceNameLength(msg.Name); err != nil {
		return err
	}

	if len(msg.Description) > MaxDescriptionLength {
		return ErrInvalidLength(DefaultCodespace, fmt.Sprintf("invalid description length; got: %d, max: %d", len(msg.Description), MaxDescriptionLength))
	}

	if len(msg.Tags) > MaxTagsNum {
		return ErrInvalidLength(DefaultCodespace, fmt.Sprintf("invalid tags size; got: %d, max: %d", len(msg.Tags), MaxTagsNum))
	}

	for i, tag := range msg.Tags {
		if len(tag) == 0 {
			return ErrInvalidLength(DefaultCodespace, fmt.Sprintf("invalid tag[%d] length: tag must not be empty", i))
		}

		if len(tag) > MaxTagLength {
			return ErrInvalidLength(DefaultCodespace, fmt.Sprintf("invalid tag[%d] length; got: %d, max: %d", i, len(tag), MaxTagLength))
		}
	}

	if len(msg.AuthorDescription) > MaxDescriptionLength {
		return ErrInvalidLength(DefaultCodespace, fmt.Sprintf("invalid author description length; got: %d, max: %d", len(msg.AuthorDescription), MaxDescriptionLength))
	}

	return nil
}

func validatePricing(pricing string) sdk.Error {
	if err := ValidateBindingPricing(pricing); err != nil {
		return err
	}

	p, err := ParsePricing(pricing)
	if err != nil {
		return err
	}

	if !validServiceCoins(p.Price) {
		return ErrInvalidPricing(DefaultCodespace, fmt.Sprintf("invalid pricing coins: %s", p.Price))
	}

	return nil
}

// ValidateRequest validates the request params
func ValidateRequest(
	serviceName string,
	serviceFeeCap sdk.Coins,
	providers []sdk.AccAddress,
	input string,
	timeout int64,
	repeated bool,
	repeatedFrequency uint64,
	repeatedTotal int64,
) sdk.Error {
	if err := ValidateServiceName(serviceName); err != nil {
		return err
	}

	if !validServiceCoins(serviceFeeCap) {
		return ErrInvalidServiceFee(DefaultCodespace, fmt.Sprintf("invalid service fee: %s", serviceFeeCap))
	}

	if len(providers) == 0 {
		return ErrInvalidProviders(DefaultCodespace, "providers missing")
	}

	if len(providers) > MaxProvidersNum {
		return ErrInvalidProviders(DefaultCodespace, fmt.Sprintf("total number of the providers must not be greater than %d", MaxProvidersNum))
	}

	if err := checkDuplicateProviders(providers); err != nil {
		return err
	}

	if len(input) == 0 {
		return ErrInvalidRequestInput(DefaultCodespace, "input missing")
	}

	if !json.Valid([]byte(input)) {
		return ErrInvalidRequestInput(DefaultCodespace, "input is not valid JSON")
	}

	if timeout <= 0 {
		return ErrInvalidTimeout(DefaultCodespace, fmt.Sprintf("timeout [%d] must be greater than 0", timeout))
	}

	if repeated {
		if repeatedFrequency > 0 && repeatedFrequency < uint64(timeout) {
			return ErrInvalidRepeatedFreq(DefaultCodespace, fmt.Sprintf("repeated frequency [%d] must not be less than timeout [%d]", repeatedFrequency, timeout))
		}

		if repeatedTotal < -1 || repeatedTotal == 0 {
			return ErrInvalidRepeatedTotal(DefaultCodespace, fmt.Sprintf("repeated total number [%d] must be greater than 0 or equal to -1", repeatedTotal))
		}
	}

	return nil
}

// ValidateRequestContextUpdating validates the request context updating operation
func ValidateRequestContextUpdating(
	providers []sdk.AccAddress,
	serviceFeeCap sdk.Coins,
	timeout int64,
	repeatedFrequency uint64,
	repeatedTotal int64,
) sdk.Error {
	if err := checkDuplicateProviders(providers); err != nil {
		return err
	}

	if !serviceFeeCap.Empty() && !validServiceCoins(serviceFeeCap) {
		return ErrInvalidServiceFee(DefaultCodespace, fmt.Sprintf("invalid service fee: %s", serviceFeeCap))
	}

	if timeout < 0 {
		return ErrInvalidTimeout(DefaultCodespace, fmt.Sprintf("timeout must not be less than 0: %d", timeout))
	}

	if timeout != 0 && repeatedFrequency != 0 && repeatedFrequency < uint64(timeout) {
		return ErrInvalidRepeatedFreq(DefaultCodespace, fmt.Sprintf("frequency [%d] must not be less than timeout [%d]", repeatedFrequency, timeout))
	}

	if repeatedTotal < -1 {
		return ErrInvalidRepeatedTotal(DefaultCodespace, fmt.Sprintf("repeated total number must not be less than -1: %d", repeatedTotal))
	}

	return nil
}

func checkDuplicateProviders(providers []sdk.AccAddress) sdk.Error {
	providerArr := make([]string, len(providers))

	for i, provider := range providers {
		providerArr[i] = provider.String()
	}

	if sdk.HasDuplicate(providerArr) {
		return ErrInvalidProviders(DefaultCodespace, "there exists duplicate providers")
	}

	return nil
}

func validServiceCoins(coins sdk.Coins) bool {
	return coins.IsValidIrisAtto()
}
