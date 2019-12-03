package types

import (
	"fmt"
	"regexp"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/tools/protoidl"
)

const (
	// types for msgs
	TypeMsgSvcDef           = "define_service"
	TypeMsgSvcBind          = "bind_service"
	TypeMsgSvcBindingUpdate = "update_service_binding"
	TypeMsgSvcDisable       = "disable_service"
	TypeMsgSvcEnable        = "enable_service"
	TypeMsgSvcRefundDeposit = "refund_service_deposit"
	TypeMsgSvcRequest       = "call_service"
	TypeMsgSvcResponse      = "respond_service"
	TypeMsgSvcRefundFees    = "refund_service_fees"
	TypeMsgSvcWithdrawFees  = "withdraw_service_fees"
	TypeMsgSvcWithdrawTax   = "withdraw_service_tax"

	// name to idetify transaction types
	outputPrivacy = "output_privacy"
	outputCached  = "output_cached"
	description   = "description"
)

var _, _, _, _, _, _, _, _, _, _, _ sdk.Msg = MsgSvcDef{}, MsgSvcBind{}, MsgSvcBindingUpdate{}, MsgSvcDisable{}, MsgSvcEnable{}, MsgSvcRefundDeposit{}, MsgSvcRequest{}, MsgSvcResponse{}, MsgSvcRefundFees{}, MsgSvcWithdrawFees{}, MsgSvcWithdrawTax{}

//______________________________________________________________________

// MsgSvcDef - struct for define a service
type MsgSvcDef struct {
	SvcDef
}

func NewMsgSvcDef(name, chainId, description string, tags []string, author sdk.AccAddress, authorDescription, idlContent string) MsgSvcDef {
	return MsgSvcDef{
		SvcDef{
			Name:              name,
			ChainId:           chainId,
			Description:       description,
			Tags:              tags,
			Author:            author,
			AuthorDescription: authorDescription,
			IDLContent:        idlContent,
		},
	}
}

func (msg MsgSvcDef) Route() string { return RouterKey }
func (msg MsgSvcDef) Type() string  { return "define_service" }

func (msg MsgSvcDef) GetSignBytes() []byte {
	if len(msg.Tags) == 0 {
		msg.Tags = nil
	}
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgSvcDef) ValidateBasic() sdk.Error {
	if len(msg.ChainId) == 0 {
		return ErrInvalidChainId(DefaultCodespace)
	}
	if !validServiceName(msg.Name) {
		return ErrInvalidServiceName(DefaultCodespace, msg.Name)
	}
	if len(msg.Author) == 0 {
		return ErrInvalidAuthor(DefaultCodespace)
	}
	if len(msg.IDLContent) == 0 {
		return ErrInvalidIDL(DefaultCodespace, "content is empty")
	}
	if err := msg.EnsureLength(); err != nil {
		return err
	}
	methods, err := protoidl.GetMethods(msg.IDLContent)
	if err != nil {
		return ErrInvalidIDL(DefaultCodespace, err.Error())
	}
	if valid, err := validateMethods(methods); !valid {
		return err
	}
	return nil
}

func (msg MsgSvcDef) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Author}
}

func validateMethods(methods []protoidl.Method) (bool, sdk.Error) {
	for _, method := range methods {
		if len(method.Name) == 0 {
			return false, ErrInvalidMethodName(DefaultCodespace)
		}
		if _, ok := method.Attributes[outputPrivacy]; ok {
			_, err := OutputPrivacyEnumFromString(method.Attributes[outputPrivacy])
			if err != nil {
				return false, ErrInvalidOutputPrivacyEnum(DefaultCodespace, method.Attributes[outputPrivacy])
			}
		}
		if _, ok := method.Attributes[outputCached]; ok {
			_, err := OutputCachedEnumFromString(method.Attributes[outputCached])
			if err != nil {
				return false, ErrInvalidOutputCachedEnum(DefaultCodespace, method.Attributes[outputCached])
			}
		}
	}
	return true, nil
}

func methodToMethodProperty(index int, method protoidl.Method) (methodProperty MethodProperty, err sdk.Error) {
	// set default value
	opp := NoPrivacy
	opc := NoCached

	var err1 error
	if _, ok := method.Attributes[outputPrivacy]; ok {
		opp, err1 = OutputPrivacyEnumFromString(method.Attributes[outputPrivacy])
		if err1 != nil {
			return methodProperty, ErrInvalidOutputPrivacyEnum(DefaultCodespace, method.Attributes[outputPrivacy])
		}
	}
	if _, ok := method.Attributes[outputCached]; ok {
		opc, err1 = OutputCachedEnumFromString(method.Attributes[outputCached])
		if err != nil {
			return methodProperty, ErrInvalidOutputCachedEnum(DefaultCodespace, method.Attributes[outputCached])
		}
	}
	methodProperty = MethodProperty{
		ID:            int16(index),
		Name:          method.Name,
		Description:   method.Attributes[description],
		OutputPrivacy: opp,
		OutputCached:  opc,
	}
	return
}

//______________________________________________________________________

// MsgSvcBinding - struct for bind a service
type MsgSvcBind struct {
	DefName     string         `json:"def_name"`
	DefChainID  string         `json:"def_chain_id"`
	BindChainID string         `json:"bind_chain_id"`
	Provider    sdk.AccAddress `json:"provider"`
	BindingType BindingType    `json:"binding_type"`
	Deposit     sdk.Coins      `json:"deposit"`
	Prices      []sdk.Coin     `json:"price"`
	Level       Level          `json:"level"`
}

func NewMsgSvcBind(defChainID, defName, bindChainID string, provider sdk.AccAddress, bindingType BindingType, deposit sdk.Coins, prices []sdk.Coin, level Level) MsgSvcBind {
	return MsgSvcBind{
		DefChainID:  defChainID,
		DefName:     defName,
		BindChainID: bindChainID,
		Provider:    provider,
		BindingType: bindingType,
		Deposit:     deposit,
		Prices:      prices,
		Level:       level,
	}
}

func (msg MsgSvcBind) Route() string { return RouterKey }
func (msg MsgSvcBind) Type() string  { return TypeMsgSvcBind }

func (msg MsgSvcBind) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgSvcBind) ValidateBasic() sdk.Error {
	if len(msg.DefChainID) == 0 {
		return ErrInvalidDefChainId(DefaultCodespace)
	}
	if len(msg.BindChainID) == 0 {
		return ErrInvalidChainId(DefaultCodespace)
	}
	if err := ensureChainIdLength(msg.DefChainID, "def_chain_id"); err != nil {
		return err
	}
	if err := ensureChainIdLength(msg.BindChainID, "bind_chain_id"); err != nil {
		return err
	}
	if !validServiceName(msg.DefName) {
		return ErrInvalidServiceName(DefaultCodespace, msg.DefName)
	}
	if err := ensureNameLength(msg.DefName); err != nil {
		return err
	}
	if !validBindingType(msg.BindingType) {
		return ErrInvalidBindingType(DefaultCodespace, msg.BindingType)
	}
	if len(msg.Provider) == 0 {
		return sdk.ErrInvalidAddress(msg.Provider.String())
	}
	if !msg.Deposit.IsValidIrisAtto() {
		return sdk.ErrInvalidCoins(fmt.Sprintf("invalid deposit [%s]", msg.Deposit))
	}
	for _, price := range msg.Prices {
		if !price.IsValidIrisAtto() {
			return sdk.ErrInvalidCoins(fmt.Sprintf("invalid price [%s]", price))
		}
	}
	if !validLevel(msg.Level) {
		return ErrInvalidLevel(DefaultCodespace, msg.Level)
	}
	return nil
}

func (msg MsgSvcBind) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgSvcBindingUpdate - struct for update a service binding
type MsgSvcBindingUpdate struct {
	DefName     string         `json:"def_name"`
	DefChainID  string         `json:"def_chain_id"`
	BindChainID string         `json:"bind_chain_id"`
	Provider    sdk.AccAddress `json:"provider"`
	BindingType BindingType    `json:"binding_type"`
	Deposit     sdk.Coins      `json:"deposit"`
	Prices      []sdk.Coin     `json:"price"`
	Level       Level          `json:"level"`
}

func NewMsgSvcBindingUpdate(defChainID, defName, bindChainID string, provider sdk.AccAddress, bindingType BindingType, deposit sdk.Coins, prices []sdk.Coin, level Level) MsgSvcBindingUpdate {
	return MsgSvcBindingUpdate{
		DefChainID:  defChainID,
		DefName:     defName,
		BindChainID: bindChainID,
		Provider:    provider,
		BindingType: bindingType,
		Deposit:     deposit,
		Prices:      prices,
		Level:       level,
	}
}
func (msg MsgSvcBindingUpdate) Route() string { return RouterKey }
func (msg MsgSvcBindingUpdate) Type() string  { return MsgSvcBindingUpdate }

func (msg MsgSvcBindingUpdate) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgSvcBindingUpdate) ValidateBasic() sdk.Error {
	if len(msg.DefChainID) == 0 {
		return ErrInvalidDefChainId(DefaultCodespace)
	}
	if len(msg.BindChainID) == 0 {
		return ErrInvalidChainId(DefaultCodespace)
	}
	if err := ensureChainIdLength(msg.DefChainID, "def_chain_id"); err != nil {
		return err
	}
	if err := ensureChainIdLength(msg.BindChainID, "bind_chain_id"); err != nil {
		return err
	}
	if !validServiceName(msg.DefName) {
		return ErrInvalidServiceName(DefaultCodespace, msg.DefName)
	}
	if err := ensureNameLength(msg.DefName); err != nil {
		return err
	}
	if len(msg.Provider) == 0 {
		return sdk.ErrInvalidAddress(msg.Provider.String())
	}
	if msg.BindingType != 0x00 && !validBindingType(msg.BindingType) {
		return ErrInvalidBindingType(DefaultCodespace, msg.BindingType)
	}
	if !msg.Deposit.IsValidIrisAtto() {
		return sdk.ErrInvalidCoins(fmt.Sprintf("invalid deposit [%s]", msg.Deposit))
	}
	for _, price := range msg.Prices {
		if !price.IsValidIrisAtto() {
			return sdk.ErrInvalidCoins(fmt.Sprintf("invalid price [%s]", price))
		}
	}
	if !validUpdateLevel(msg.Level) {
		return ErrInvalidLevel(DefaultCodespace, msg.Level)
	}
	return nil
}

func (msg MsgSvcBindingUpdate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgSvcDisable - struct for disable a service binding
type MsgSvcDisable struct {
	DefName     string         `json:"def_name"`
	DefChainID  string         `json:"def_chain_id"`
	BindChainID string         `json:"bind_chain_id"`
	Provider    sdk.AccAddress `json:"provider"`
}

func NewMsgSvcDisable(defChainID, defName, bindChainID string, provider sdk.AccAddress) MsgSvcDisable {
	return MsgSvcDisable{
		DefChainID:  defChainID,
		DefName:     defName,
		BindChainID: bindChainID,
		Provider:    provider,
	}
}

func (msg MsgSvcDisable) Route() string { return RouterKey }
func (msg MsgSvcDisable) Type() string  { return TypeMsgSvcDisable }

func (msg MsgSvcDisable) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgSvcDisable) ValidateBasic() sdk.Error {
	if len(msg.DefChainID) == 0 {
		return ErrInvalidDefChainId(DefaultCodespace)
	}
	if len(msg.BindChainID) == 0 {
		return ErrInvalidChainId(DefaultCodespace)
	}
	if err := ensureChainIdLength(msg.DefChainID, "def_chain_id"); err != nil {
		return err
	}
	if err := ensureChainIdLength(msg.BindChainID, "bind_chain_id"); err != nil {
		return err
	}
	if !validServiceName(msg.DefName) {
		return ErrInvalidServiceName(DefaultCodespace, msg.DefName)
	}
	if err := ensureNameLength(msg.DefName); err != nil {
		return err
	}
	if len(msg.Provider) == 0 {
		return sdk.ErrInvalidAddress(msg.Provider.String())
	}
	return nil
}

func (msg MsgSvcDisable) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgSvcEnable - struct for enable a service binding
type MsgSvcEnable struct {
	DefName     string         `json:"def_name"`
	DefChainID  string         `json:"def_chain_id"`
	BindChainID string         `json:"bind_chain_id"`
	Provider    sdk.AccAddress `json:"provider"`
	Deposit     sdk.Coins      `json:"deposit"`
}

func NewMsgSvcEnable(defChainID, defName, bindChainID string, provider sdk.AccAddress, deposit sdk.Coins) MsgSvcEnable {
	return MsgSvcEnable{
		DefChainID:  defChainID,
		DefName:     defName,
		BindChainID: bindChainID,
		Provider:    provider,
		Deposit:     deposit,
	}
}

func (msg MsgSvcEnable) Route() string { return RouterKey }
func (msg MsgSvcEnable) Type() string  { return TypeMsgSvcEnable }

func (msg MsgSvcEnable) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgSvcEnable) ValidateBasic() sdk.Error {
	if len(msg.DefChainID) == 0 {
		return ErrInvalidDefChainId(DefaultCodespace)
	}
	if len(msg.BindChainID) == 0 {
		return ErrInvalidChainId(DefaultCodespace)
	}
	if err := ensureChainIdLength(msg.DefChainID, "def_chain_id"); err != nil {
		return err
	}
	if err := ensureChainIdLength(msg.BindChainID, "bind_chain_id"); err != nil {
		return err
	}
	if !validServiceName(msg.DefName) {
		return ErrInvalidServiceName(DefaultCodespace, msg.DefName)
	}
	if err := ensureNameLength(msg.DefName); err != nil {
		return err
	}
	if !msg.Deposit.IsValidIrisAtto() {
		return sdk.ErrInvalidCoins(fmt.Sprintf("invalid deposit [%s]", msg.Deposit))
	}
	if len(msg.Provider) == 0 {
		return sdk.ErrInvalidAddress(msg.Provider.String())
	}
	return nil
}

func (msg MsgSvcEnable) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgSvcRefundDeposit - struct for refund deposit from a service binding
type MsgSvcRefundDeposit struct {
	DefName     string         `json:"def_name"`
	DefChainID  string         `json:"def_chain_id"`
	BindChainID string         `json:"bind_chain_id"`
	Provider    sdk.AccAddress `json:"provider"`
}

func NewMsgSvcRefundDeposit(defChainID, defName, bindChainID string, provider sdk.AccAddress) MsgSvcRefundDeposit {
	return MsgSvcRefundDeposit{
		DefChainID:  defChainID,
		DefName:     defName,
		BindChainID: bindChainID,
		Provider:    provider,
	}
}

func (msg MsgSvcRefundDeposit) Route() string { return RouterKey }
func (msg MsgSvcRefundDeposit) Type() string  { return TypeMsgSvcRefundDeposit }

func (msg MsgSvcRefundDeposit) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgSvcRefundDeposit) ValidateBasic() sdk.Error {
	if len(msg.DefChainID) == 0 {
		return ErrInvalidDefChainId(DefaultCodespace)
	}
	if len(msg.BindChainID) == 0 {
		return ErrInvalidChainId(DefaultCodespace)
	}
	if err := ensureChainIdLength(msg.DefChainID, "def_chain_id"); err != nil {
		return err
	}
	if err := ensureChainIdLength(msg.BindChainID, "bind_chain_id"); err != nil {
		return err
	}
	if !validServiceName(msg.DefName) {
		return ErrInvalidServiceName(DefaultCodespace, msg.DefName)
	}
	if err := ensureNameLength(msg.DefName); err != nil {
		return err
	}
	if len(msg.Provider) == 0 {
		return sdk.ErrInvalidAddress(msg.Provider.String())
	}
	return nil
}

func (msg MsgSvcRefundDeposit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgSvcRequest - struct for call a service
type MsgSvcRequest struct {
	DefChainID  string         `json:"def_chain_id"`
	DefName     string         `json:"def_name"`
	BindChainID string         `json:"bind_chain_id"`
	ReqChainID  string         `json:"req_chain_id"`
	MethodID    int16          `json:"method_id"`
	Provider    sdk.AccAddress `json:"provider"`
	Consumer    sdk.AccAddress `json:"consumer"`
	Input       []byte         `json:"input"`
	ServiceFee  sdk.Coins      `json:"service_fee"`
	Profiling   bool           `json:"profiling"`
}

func NewMsgSvcRequest(defChainID, defName, bindChainID, reqChainID string, consumer, provider sdk.AccAddress, methodID int16, input []byte, serviceFee sdk.Coins, profiling bool) MsgSvcRequest {
	return MsgSvcRequest{
		DefChainID:  defChainID,
		DefName:     defName,
		BindChainID: bindChainID,
		ReqChainID:  reqChainID,
		Consumer:    consumer,
		Provider:    provider,
		MethodID:    methodID,
		Input:       input,
		ServiceFee:  serviceFee,
		Profiling:   profiling,
	}
}

func (msg MsgSvcRequest) Route() string { return RouterKey }
func (msg MsgSvcRequest) Type() string  { return TypeMsgSvcRequest }

func (msg MsgSvcRequest) GetSignBytes() []byte {
	if len(msg.Input) == 0 {
		msg.Input = nil
	}
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgSvcRequest) ValidateBasic() sdk.Error {
	if len(msg.DefChainID) == 0 {
		return ErrInvalidDefChainId(DefaultCodespace)
	}
	if len(msg.BindChainID) == 0 {
		return ErrInvalidBindChainId(DefaultCodespace)
	}
	if len(msg.ReqChainID) == 0 {
		return ErrInvalidChainId(DefaultCodespace)
	}
	if err := ensureChainIdLength(msg.DefChainID, "def_chain_id"); err != nil {
		return err
	}
	if err := ensureChainIdLength(msg.BindChainID, "bind_chain_id"); err != nil {
		return err
	}
	if err := ensureChainIdLength(msg.ReqChainID, "req_chain_id"); err != nil {
		return err
	}
	if !validServiceName(msg.DefName) {
		return ErrInvalidServiceName(DefaultCodespace, msg.DefName)
	}
	if err := ensureNameLength(msg.DefName); err != nil {
		return err
	}
	if len(msg.Provider) == 0 {
		return sdk.ErrInvalidAddress(msg.Provider.String())
	}
	if len(msg.Consumer) == 0 {
		return sdk.ErrInvalidAddress(msg.Consumer.String())
	}
	if !msg.ServiceFee.IsValidIrisAtto() {
		return sdk.ErrInvalidCoins(fmt.Sprintf("invalid service fee [%s]", msg.ServiceFee))
	}
	return nil
}

func (msg MsgSvcRequest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Consumer}
}

//______________________________________________________________________

// MsgSvcResponse - struct for respond a service call
type MsgSvcResponse struct {
	ReqChainID string         `json:"req_chain_id"`
	RequestID  string         `json:"request_id"`
	Provider   sdk.AccAddress `json:"provider"`
	Output     []byte         `json:"output"`
	ErrorMsg   []byte         `json:"error_msg"`
}

func NewMsgSvcResponse(reqChainID string, requestId string, provider sdk.AccAddress, output, errorMsg []byte) MsgSvcResponse {
	return MsgSvcResponse{
		ReqChainID: reqChainID,
		RequestID:  requestId,
		Provider:   provider,
		Output:     output,
		ErrorMsg:   errorMsg,
	}
}

func (msg MsgSvcResponse) Route() string { return RouterKey }
func (msg MsgSvcResponse) Type() string  { return TypeMsgSvcResponse }

func (msg MsgSvcResponse) GetSignBytes() []byte {
	if len(msg.Output) == 0 {
		msg.Output = nil
	}
	if len(msg.ErrorMsg) == 0 {
		msg.ErrorMsg = nil
	}
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgSvcResponse) ValidateBasic() sdk.Error {
	if len(msg.ReqChainID) == 0 {
		return ErrInvalidReqChainId(DefaultCodespace)
	}
	if err := ensureChainIdLength(msg.ReqChainID, "req_chain_id"); err != nil {
		return err
	}
	if len(msg.Provider) == 0 {
		return sdk.ErrInvalidAddress(msg.Provider.String())
	}
	_, _, _, err := ConvertRequestID(msg.RequestID)
	if err != nil {
		return ErrInvalidReqId(DefaultCodespace, msg.RequestID)
	}

	return nil
}

func (msg MsgSvcResponse) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgSvcRefundFees - struct for refund fees
type MsgSvcRefundFees struct {
	Consumer sdk.AccAddress `json:"consumer"`
}

func NewMsgSvcRefundFees(consumer sdk.AccAddress) MsgSvcRefundFees {
	return MsgSvcRefundFees{
		Consumer: consumer,
	}
}

func (msg MsgSvcRefundFees) Route() string { return RouterKey }
func (msg MsgSvcRefundFees) Type() string  { return TypeMsgSvcRefundFees }

func (msg MsgSvcRefundFees) GetSignBytes() []byte {
	b := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(b)
}

func (msg MsgSvcRefundFees) ValidateBasic() sdk.Error {
	if len(msg.Consumer) == 0 {
		return sdk.ErrInvalidAddress(msg.Consumer.String())
	}
	return nil
}

func (msg MsgSvcRefundFees) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Consumer}
}

//______________________________________________________________________

// MsgSvcWithdrawFees - struct for withdraw fees
type MsgSvcWithdrawFees struct {
	Provider sdk.AccAddress `json:"provider"`
}

func NewMsgSvcWithdrawFees(provider sdk.AccAddress) MsgSvcWithdrawFees {
	return MsgSvcWithdrawFees{
		Provider: provider,
	}
}

func (msg MsgSvcWithdrawFees) Route() string { return RouterKey }
func (msg MsgSvcWithdrawFees) Type() string  { return TypeMsgSvcWithdrawFees }

func (msg MsgSvcWithdrawFees) GetSignBytes() []byte {
	b := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(b)
}

func (msg MsgSvcWithdrawFees) ValidateBasic() sdk.Error {
	if len(msg.Provider) == 0 {
		return sdk.ErrInvalidAddress(msg.Provider.String())
	}
	return nil
}

func (msg MsgSvcWithdrawFees) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgSvcWithdrawTax - struct for withdraw tax
type MsgSvcWithdrawTax struct {
	Trustee     sdk.AccAddress `json:"trustee"`
	DestAddress sdk.AccAddress `json:"dest_address"`
	Amount      sdk.Coins      `json:"amount"`
}

func NewMsgSvcWithdrawTax(trustee, destAddress sdk.AccAddress, amount sdk.Coins) MsgSvcWithdrawTax {
	return MsgSvcWithdrawTax{
		Trustee:     trustee,
		DestAddress: destAddress,
		Amount:      amount,
	}
}

func (msg MsgSvcWithdrawTax) Route() string { return RouterKey }
func (msg MsgSvcWithdrawTax) Type() string  { return TypeMsgSvcWithdrawTax }

func (msg MsgSvcWithdrawTax) GetSignBytes() []byte {
	b := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(b)
}

func (msg MsgSvcWithdrawTax) ValidateBasic() sdk.Error {
	if len(msg.Trustee) == 0 {
		return sdk.ErrInvalidAddress(msg.Trustee.String())
	}
	if len(msg.DestAddress) == 0 {
		return sdk.ErrInvalidAddress(msg.DestAddress.String())
	}
	if !msg.Amount.IsValidIrisAtto() {
		return sdk.ErrInvalidCoins(fmt.Sprintf("invalid withdraw amount [%s]", msg.Amount))
	}
	return nil
}

func (msg MsgSvcWithdrawTax) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Trustee}
}

//______________________________________________________________________

func validServiceName(name string) bool {
	if len(name) == 0 || len(name) > 128 {
		return false
	}

	// Must contain alphanumeric characters, _ and - only
	reg := regexp.MustCompile(`[^a-zA-Z0-9_-]`)
	return !reg.Match([]byte(name))
}

func (msg MsgSvcDef) EnsureLength() sdk.Error {
	if err := ensureNameLength(msg.Name); err != nil {
		return err
	}

	if err := ensureChainIdLength(msg.ChainId, "chain_id"); err != nil {
		return err
	}

	if len(msg.Description) > 280 {
		return sdk.ErrInvalidLength(DefaultCodespace, CodeInvalidInput, "description", len(msg.Description), 280)
	}
	if len(msg.Tags) > 10 {
		return sdk.ErrInvalidLength(DefaultCodespace, CodeInvalidInput, "tags", len(msg.Tags), 10)
	} else {
		for i, tag := range msg.Tags {
			if len(tag) > 70 {
				return sdk.ErrInvalidLength(DefaultCodespace, CodeInvalidInput, fmt.Sprintf("tags[%d]", i), len(tag), 70)
			}
		}
	}
	if len(msg.AuthorDescription) > 280 {
		return sdk.ErrInvalidLength(DefaultCodespace, CodeInvalidInput, "author_description", len(msg.AuthorDescription), 280)
	}
	return nil
}

func ensureNameLength(name string) sdk.Error {
	if len(name) > 70 {
		return sdk.ErrInvalidLength(DefaultCodespace, CodeInvalidInput, "name", len(name), 70)
	}
	return nil
}

func ensureChainIdLength(chainId, fieldNm string) sdk.Error {
	if len(chainId) > 50 {
		return sdk.ErrInvalidLength(DefaultCodespace, CodeInvalidInput, fieldNm, len(chainId), 50)
	}
	return nil
}
