package types

import (
	sdk "github.com/irisnet/irishub/types"
)

const (
	DefaultCodespace sdk.CodespaceType = "oracle"

	CodeUnknownFeedName          sdk.CodeType = 100
	CodeInvalidFeedName          sdk.CodeType = 101
	CodeExistedFeedName          sdk.CodeType = 102
	CodeUnauthorized             sdk.CodeType = 103
	CodeInvalidServiceName       sdk.CodeType = 104
	CodeInvalidLatestHistory     sdk.CodeType = 105
	CodeEmptyProviders           sdk.CodeType = 106
	CodeInvalidServiceFeeCap     sdk.CodeType = 107
	CodeInvalidResponseThreshold sdk.CodeType = 108
	CodeInvalidAddress           sdk.CodeType = 109
	CodeInvalidAggregateFunc     sdk.CodeType = 110
	CodeInvalidValueJsonPath     sdk.CodeType = 111
	CodeUnknownRequestContextID  sdk.CodeType = 112
	CodeNotRegisterFunc          sdk.CodeType = 113
	CodeInvalidFeedState         sdk.CodeType = 114
	CodeNotProfiler              sdk.CodeType = 115
	CodeInvalidDescription       sdk.CodeType = 116
	CodeInvalidTimeout           sdk.CodeType = 117
)

func ErrUnknownFeedName(codespace sdk.CodespaceType, feedName string) sdk.Error {
	return sdk.NewError(codespace, CodeUnknownFeedName, "feed name %s does not exist", feedName)
}

func ErrInvalidFeedName(codespace sdk.CodespaceType, feedName string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidFeedName, "invalid feed name %s; only alphanumeric characters, _ and - accepted, the length ranges in (0,%d]", feedName, MaxNameLen)
}

func ErrExistedFeedName(codespace sdk.CodespaceType, feedName string) sdk.Error {
	return sdk.NewError(codespace, CodeExistedFeedName, "feed name %s already exists", feedName)
}

func ErrUnauthorized(codespace sdk.CodespaceType, feedName string, owner sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeUnauthorized, "feed %s does not belong to %s", feedName, owner.String())
}

func ErrInvalidServiceName(codespace sdk.CodespaceType, serviceName string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidServiceName, "service name %s; only alphanumeric characters, _ and - accepted, the length ranges in (0,%d]", serviceName, MaxNameLen)
}

func ErrEmptyProviders(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyProviders, "provider can not be empty")
}

func ErrInvalidLatestHistory(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidLatestHistory, "latest history is invalid, should be between 1 and %d", MaxLatestHistory)
}

func ErrInvalidServiceFeeCap(codespace sdk.CodespaceType, fees sdk.Coins) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidServiceFeeCap, "service fee %s is invalid", fees.String())
}

func ErrInvalidResponseThreshold(codespace sdk.CodespaceType, limit int) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidResponseThreshold, "response threshold should be between 1 and %d", limit)
}

func ErrInvalidAddress(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAddress, msg)
}

func ErrInvalidAggregateFunc(codespace sdk.CodespaceType, len int) sdk.Error {
	return sdk.ErrInvalidLength(codespace, CodeInvalidAggregateFunc, "aggregateFunc", len, MaxNameLen)
}

func ErrInvalidValueJsonPath(codespace sdk.CodespaceType, len int) sdk.Error {
	return sdk.ErrInvalidLength(codespace, CodeInvalidValueJsonPath, "valueJsonPath", len, MaxNameLen)
}

func ErrUnknownRequestContextID(codespace sdk.CodespaceType, reqCtxID []byte) sdk.Error {
	return sdk.NewError(codespace, CodeUnknownRequestContextID, "request context ID %s does not exist", string(reqCtxID))
}

func ErrNotRegisterFunc(codespace sdk.CodespaceType, methodName string) sdk.Error {
	return sdk.NewError(codespace, CodeNotRegisterFunc, "method %s don't register", methodName)
}

func ErrInvalidFeedState(codespace sdk.CodespaceType, feedName string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidFeedState, "feed %s may be a invalid state", feedName)
}

func ErrNotProfiler(codespace sdk.CodespaceType, profiler sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeNotProfiler, "%s is not a profiler address", profiler)
}

func ErrInvalidDescription(codespace sdk.CodespaceType, descLen int) sdk.Error {
	return sdk.ErrInvalidLength(codespace, CodeInvalidDescription, "description", descLen, MaxDescriptionLen)
}

func ErrInvalidTimeout(codespace sdk.CodespaceType, timeout int64, frequency uint64) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidTimeout, "timeout[%d] should be no more than frequency[%d]", timeout, frequency)
}
