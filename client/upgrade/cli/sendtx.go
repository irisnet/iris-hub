package cli

import (
	"os"

	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/upgrade"
	sdk "github.com/irisnet/irishub/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagProposalID = "proposal-id"
	flagTitle      = "title"
	flagVoter      = "voter"
)

// submit switch msg
func GetCmdSubmitSwitch(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "submit-switch",
		Short:   "Submit a switch msg for a upgrade propsal",
		Example: "iriscli upgrade submit-switch --chain-id=<chain-id> --from=<key name> --fee=0.004iris --proposal-id 1 --title <title>",
		RunE: func(cmd *cobra.Command, args []string) error {
			title := viper.GetString(flagTitle)
			proposalID := uint64(viper.GetInt64(flagProposalID))

			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := context.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			// get the from/to address
			from, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			msg := upgrade.NewMsgSwitch(title, proposalID, from)
			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			// Build and sign the transaction, then broadcast to a Tendermint
			// node.
			cliCtx.PrintResponse = true
			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagTitle, "", "title of switch")
	cmd.Flags().String(flagProposalID, "", "proposalID of upgrade proposal")

	return cmd
}
