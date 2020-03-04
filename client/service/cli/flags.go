// nolint
package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagName              = "name"
	FlagDescription       = "description"
	FlagTags              = "tags"
	FlagAuthorDescription = "author-description"
	FlagSchemas           = "schemas"
	FlagServiceName       = "service-name"
	FlagDeposit           = "deposit"
	FlagPricing           = "pricing"
	FlagProviders         = "providers"
	FlagServiceFeeCap     = "service-fee-cap"
	FlagTimeout           = "timeout"
	FlagData              = "data"
	FlagSuperMode         = "super-mode"
	FlagRepeated          = "repeated"
	FlagFrequency         = "frequency"
	FlagTotal             = "total"
	FlagRequestID         = "request-id"
	FlagError             = "error"
)

// common flagsets to add to various functions
var (
	FsServiceDefine               = flag.NewFlagSet("", flag.ContinueOnError)
	FsServiceBind                 = flag.NewFlagSet("", flag.ContinueOnError)
	FsServiceUpdateBinding        = flag.NewFlagSet("", flag.ContinueOnError)
	FsServiceEnable               = flag.NewFlagSet("", flag.ContinueOnError)
	FsServiceRequest              = flag.NewFlagSet("", flag.ContinueOnError)
	FsServiceRespond              = flag.NewFlagSet("", flag.ContinueOnError)
	FsServiceUpdateRequestContext = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsServiceDefine.String(FlagName, "", "service name")
	FsServiceDefine.String(FlagDescription, "", "service description")
	FsServiceDefine.StringSlice(FlagTags, []string{}, "service tags")
	FsServiceDefine.String(FlagAuthorDescription, "", "service author description")
	FsServiceDefine.String(FlagSchemas, "", "interface schemas content or path")

	FsServiceBind.String(FlagServiceName, "", "service name")
	FsServiceBind.String(FlagDeposit, "", "deposit of the binding")
	FsServiceBind.String(FlagPricing, "", "pricing content or path")

	FsServiceUpdateBinding.String(FlagDeposit, "", "added deposit for the binding")
	FsServiceUpdateBinding.String(FlagPricing, "", "new pricing of the binding, which is a Pricing JSON Schema instance")

	FsServiceEnable.String(FlagDeposit, "", "added deposit for enabling the binding")

	FsServiceRequest.String(FlagServiceName, "", "service name")
	FsServiceRequest.StringSlice(FlagProviders, []string{}, "providers to request")
	FsServiceRequest.String(FlagServiceFeeCap, "", "maximal fee to pay for a service request")
	FsServiceRequest.String(FlagData, "", "input of the service request, which is an Input JSON schema instance")
	FsServiceRequest.Uint64(FlagTimeout, 0, "request timeout")
	FsServiceRequest.Bool(FlagSuperMode, false, "indicate if the signer is a super user")
	FsServiceRequest.Bool(FlagRepeated, false, "indicate if the request is repetitive")
	FsServiceRequest.Uint64(FlagFrequency, 0, "request frequency when repeated, default to timeout")
	FsServiceRequest.Int64(FlagTotal, 0, "request count when repeated, -1 means unlimited")

	FsServiceUpdateRequestContext.StringSlice(FlagProviders, []string{}, "new providers to request, not updated if empty")
	FsServiceUpdateRequestContext.String(FlagServiceFeeCap, "", "maximal fee to pay for a service request, empty means not updated")
	FsServiceUpdateRequestContext.Uint64(FlagTimeout, 0, "new request timeout, not updated if set to 0")
	FsServiceUpdateRequestContext.Uint64(FlagFrequency, 0, "new request frequency, not updated if set to 0")
	FsServiceUpdateRequestContext.Int64(FlagTotal, 0, "new request count, not updated if set to 0")

	FsServiceRespond.String(FlagRequestID, "", "ID of the request to respond")
	FsServiceRespond.String(FlagData, "", "content or path of the output of the service response, which is an Output JSON schema instance")
	FsServiceRespond.String(FlagError, "", "content or path of the err msg of the service response, which is an Error JSON schema instance")
}
