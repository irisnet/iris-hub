package cli

import (
	"fmt"
	"os"

	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/modules/gov"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GetCmdSubmitProposal implements submitting a proposal transaction command.
func GetCmdSubmitProposal(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-proposal",
		Short: "Submit a proposal along with an initial deposit",
		RunE: func(cmd *cobra.Command, args []string) error {
			title := viper.GetString(flagTitle)
			description := viper.GetString(flagDescription)
			strProposalType := viper.GetString(flagProposalType)
			initialDeposit := viper.GetString(flagDeposit)
			paramsStr := viper.GetString(flagParams)

			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))
			txCtx := context.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			fromAddr, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			amount, err := cliCtx.ParseCoins(initialDeposit)
			if err != nil {
				return err
			}

			proposalType, err := gov.ProposalTypeFromString(strProposalType)
			if err != nil {
				return err
			}

			var params gov.Params
			if proposalType == gov.ProposalTypeParameterChange {
				if err := json.Unmarshal([]byte(paramsStr), &params); err != nil {
					fmt.Println(err.Error())
					return nil
				}
			}

			msg := gov.NewMsgSubmitProposal(title, description, proposalType, fromAddr, amount, params)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			if cliCtx.GenerateOnly {
				return utils.PrintUnsignedStdTx(txCtx, cliCtx, []sdk.Msg{msg})
			}
			// Build and sign the transaction, then broadcast to Tendermint
			// proposalID must be returned, and it is a part of response.
			cliCtx.PrintResponse = true

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagTitle, "", "title of proposal")
	cmd.Flags().String(flagDescription, "", "description of proposal")
	cmd.Flags().String(flagProposalType, "", "proposalType of proposal,eg:Text/ParameterChange/SoftwareUpgrade")
	cmd.Flags().String(flagDeposit, "", "deposit of proposal")
	cmd.Flags().String(flagParams, "", "parameter of proposal,eg. [{key:key,value:value,op:update}]")

	return cmd
}

// GetCmdDeposit implements depositing tokens for an active proposal.
func GetCmdDeposit(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit",
		Short: "deposit tokens for activing proposal",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))
			txCtx := context.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			depositerAddr, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			proposalID := viper.GetInt64(flagProposalID)

			amount, err := cliCtx.ParseCoins(viper.GetString(flagDeposit))
			if err != nil {
				return err
			}

			msg := gov.NewMsgDeposit(depositerAddr, proposalID, amount)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}
			if cliCtx.GenerateOnly {
				return utils.PrintUnsignedStdTx(txCtx, cliCtx, []sdk.Msg{msg})
			}
			// Build and sign the transaction, then broadcast to a Tendermint
			// node.
			cliCtx.PrintResponse = true

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagProposalID, "", "proposalID of proposal depositing on")
	cmd.Flags().String(flagDeposit, "", "amount of deposit")

	return cmd
}

// GetCmdVote implements creating a new vote command.
func GetCmdVote(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote",
		Short: "vote for an active proposal, options: Yes/No/NoWithVeto/Abstain",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))
			txCtx := context.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			voterAddr, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			proposalID := viper.GetInt64(flagProposalID)
			option := viper.GetString(flagOption)

			byteVoteOption, err := gov.VoteOptionFromString(option)
			if err != nil {
				return err
			}

			msg := gov.NewMsgVote(voterAddr, proposalID, byteVoteOption)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			fmt.Printf("Vote[Voter:%s,ProposalID:%d,Option:%s]",
				voterAddr.String(), msg.ProposalID, msg.Option.String(),
			)

			if cliCtx.GenerateOnly {
				return utils.PrintUnsignedStdTx(txCtx, cliCtx, []sdk.Msg{msg})
			}
			// Build and sign the transaction, then broadcast to a Tendermint
			// node.
			cliCtx.PrintResponse = true

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagProposalID, "", "proposalID of proposal voting on")
	cmd.Flags().String(flagOption, "", "vote option {Yes, No, NoWithVeto, Abstain}")

	return cmd
}
