package cli

import (
	"os"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/modules/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/irisnet/irishub/client/context"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	cmn "github.com/irisnet/irishub/client/service"
)

func GetCmdQueryScvDef(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "definition",
		Short:   "Query service definition",
		Example: "iriscli service definition --def-chain-id=<chain-id> --service-name=<service name>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			name := viper.GetString(FlagServiceName)
			defChainId := viper.GetString(FlagDefChainID)

			res, err := cliCtx.QueryStore(service.GetServiceDefinitionKey(defChainId, name), storeName)
			if err != nil {
				return err
			}
			if len(res) == 0 {
				return fmt.Errorf("chain-id [%s] service [%s] is not existed", defChainId, name)
			}

			var svcDef service.SvcDef
			cdc.MustUnmarshalBinaryLengthPrefixed(res, &svcDef)

			res2, err := cliCtx.QuerySubspace(service.GetMethodsSubspaceKey(defChainId, name), storeName)
			if err != nil {
				return err
			}

			var methods []service.MethodProperty
			for i := 0; i < len(res2); i++ {
				var method service.MethodProperty
				cdc.MustUnmarshalBinaryLengthPrefixed(res2[i].Value, &method)
				methods = append(methods, method)
			}

			service := cmn.ServiceOutput{SvcDef: svcDef, Methods: methods}
			output, err := codec.MarshalJSONIndent(cdc, service)
			if err != nil {
				return err
			}

			fmt.Println(string(output))
			return nil
		},
	}
	cmd.Flags().AddFlagSet(FsDefChainID)
	cmd.Flags().AddFlagSet(FsServiceName)

	return cmd
}

func GetCmdQueryScvBind(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "binding",
		Short:   "Query service binding",
		Example: "iriscli service binding --def-chain-id=<chain-id> --service-name=<service name> --bind-chain-id=<chain-id> --provider=<provider>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			name := viper.GetString(FlagServiceName)
			defChainId := viper.GetString(FlagDefChainID)
			bindChainId := viper.GetString(FlagBindChainID)
			providerStr := viper.GetString(FlagProvider)

			provider, err := sdk.AccAddressFromBech32(providerStr)
			if err != nil {
				return err
			}
			res, err := cliCtx.QueryStore(service.GetServiceBindingKey(defChainId, name, bindChainId, provider), storeName)
			if err != nil {
				return err
			}
			if len(res) == 0 {
				return fmt.Errorf("def-chain-id [%s] service [%s] bind-chain-id [%s] provider [%s] is not existed", defChainId, name, bindChainId, provider)
			}

			var svcBinding service.SvcBinding
			cdc.MustUnmarshalBinaryLengthPrefixed(res, &svcBinding)
			output, err := codec.MarshalJSONIndent(cdc, svcBinding)
			fmt.Println(string(output))
			return nil
		},
	}
	cmd.Flags().AddFlagSet(FsDefChainID)
	cmd.Flags().AddFlagSet(FsServiceName)
	cmd.Flags().AddFlagSet(FsBindChainID)
	cmd.Flags().AddFlagSet(FsProvider)

	return cmd
}

func GetCmdQueryScvBinds(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "bindings",
		Short:   "Query service bindings",
		Example: "iriscli service bindings --def-chain-id=<chain-id> --service-name=<service name>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			name := viper.GetString(FlagServiceName)
			defChainId := viper.GetString(FlagDefChainID)

			res, err := cliCtx.QuerySubspace(service.GetBindingsSubspaceKey(defChainId, name), storeName)
			if err != nil {
				return err
			}

			var bindings []service.SvcBinding
			for i := 0; i < len(res); i++ {

				var binding service.SvcBinding
				cdc.MustUnmarshalBinaryLengthPrefixed(res[i].Value, &binding)
				bindings = append(bindings, binding)
			}

			output, err := cdc.MarshalJSONIndent(bindings, "", "")
			fmt.Println(string(output))
			return nil
		},
	}
	cmd.Flags().AddFlagSet(FsDefChainID)
	cmd.Flags().AddFlagSet(FsServiceName)

	return cmd
}
