package cli

import (
	"os"

	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/slashing"
	sdk "github.com/irisnet/irishub/types"
	"github.com/spf13/cobra"
)

// GetCmdUnrevoke implements the create unrevoke validator command.
func GetCmdUnrevoke(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "unjail",
		Args:    cobra.ExactArgs(0),
		Short:   "unjail validator previously jailed for downtime",
		Example: "iriscli stake unjail --from <key name> --fee=0.4iris --chain-id=<chain-id>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			validatorAddr, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			msg := slashing.NewMsgUnjail(sdk.ValAddress(validatorAddr))

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	return cmd
}
