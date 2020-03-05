package cli

import (
	"os"

	"github.com/irisnet/irishub/app/v3/oracle"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GetCmdCreateFeed implements defining a feed command
func GetCmdCreateFeed(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: `create a new feed,the feed will be in "paused" state`,
		Example: `iriscli oracle create --chain-id="irishub-test" --from=node0 --fee=0.3iris  --commit ` +
			`--feed-name="test-feed" ` +
			`--latest-history=10 ` +
			`--service-name="test-service" ` +
			`--input={request-data} ` +
			`--providers="faa1hp29kuh22vpjjlnctmyml5s75evsnsd8r4x0mm,faa15rurzhkemsgfm42dnwhafjdv5s8e2pce0ku8ya" ` +
			`--service-fee-cap=1iris ` +
			`--timeout=2 ` +
			`--frequency=10 ` +
			`--total=10 ` +
			`--threshold=1 ` +
			`--aggregate-func="avg" ` +
			`--value-json-path="high"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			creator, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			var providers []sdk.AccAddress
			for _, addr := range viper.GetStringSlice(FlagProviders) {
				provider, err := sdk.AccAddressFromBech32(addr)
				if err != nil {
					return err
				}
				providers = append(providers, provider)
			}

			serviceFeeCap, err := cliCtx.ParseCoins(viper.GetString(FlagServiceFeeCap))
			if err != nil {
				return err
			}

			msg := oracle.MsgCreateFeed{
				FeedName:          viper.GetString(FlagFeedName),
				AggregateFunc:     viper.GetString(FlagAggregateFunc),
				ValueJsonPath:     viper.GetString(FlagValueJsonPath),
				LatestHistory:     uint64(viper.GetInt64(FlagLatestHistory)),
				Description:       viper.GetString(FlagDescription),
				ServiceName:       viper.GetString(FlagServiceName),
				Providers:         providers,
				Input:             viper.GetString(FlagInput),
				Timeout:           viper.GetInt64(FlagTimeout),
				ServiceFeeCap:     serviceFeeCap,
				RepeatedFrequency: uint64(viper.GetInt64(FlagFrequency)),
				RepeatedTotal:     viper.GetInt64(FlagTotal),
				ResponseThreshold: uint16(viper.GetInt(FlagThreshold)),
				Creator:           creator,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(FsCreateFeed)
	_ = cmd.MarkFlagRequired(FlagFeedName)
	_ = cmd.MarkFlagRequired(FlagAggregateFunc)
	_ = cmd.MarkFlagRequired(FlagValueJsonPath)
	_ = cmd.MarkFlagRequired(FlagLatestHistory)
	_ = cmd.MarkFlagRequired(FlagServiceName)
	_ = cmd.MarkFlagRequired(FlagProviders)
	_ = cmd.MarkFlagRequired(FlagServiceFeeCap)
	_ = cmd.MarkFlagRequired(FlagTimeout)
	return cmd
}

// GetCmdStartFeed implements start a feed command
func GetCmdStartFeed(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "start",
		Short:   `Start a feed in "paused" state`,
		Args:    cobra.ExactArgs(1),
		Example: `iriscli oracle start <feed-name> --chain-id="irishub-test" --from=<creator> --fee=0.3iris --commit `,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			creator, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}
			msg := oracle.MsgStartFeed{
				FeedName: args[0],
				Creator:  creator,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(FsStartFeed)
	return cmd
}

// GetCmdPauseFeed implements pause a running feed command
func GetCmdPauseFeed(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pause",
		Short:   `Pause a feed in "running" state`,
		Args:    cobra.ExactArgs(1),
		Example: `iriscli oracle pause <feed-name> --chain-id="irishub-test" --from=<creator> --fee=0.3iris --commit  `,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			creator, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}
			msg := oracle.MsgPauseFeed{
				FeedName: args[0],
				Creator:  creator,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(FsStartFeed)
	return cmd
}

// GetCmdEditFeed implements edit a feed command
func GetCmdEditFeed(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit",
		Short: "Modify the feed information and update service invocation parameters by feed creator",
		Args:  cobra.ExactArgs(1),
		Example: `iriscli oracle edit <feed-name> --chain-id="irishub-test" --from=<creator> --fee=0.3iris --commit  ` +
			`--latest-history=10 ` +
			`--providers="faa1r3tyupskwlh07dmhjw70frxzaaaufta37y25yr,faa1ydahnhrhkjh9j9u0jn8p3s272l0ecqj40vra8h"` +
			`--service-fee-cap=1iris ` +
			`--timeout=2 ` +
			`--frequency=10 ` +
			`--threshold=5 ` +
			`--total=-1 ` +
			`--threshold=1`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			creator, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}
			var providers []sdk.AccAddress
			for _, addr := range viper.GetStringSlice(FlagProviders) {
				provider, err := sdk.AccAddressFromBech32(addr)
				if err != nil {
					return err
				}
				providers = append(providers, provider)
			}

			serviceFeeCap, err := cliCtx.ParseCoins(viper.GetString(FlagServiceFeeCap))
			if err != nil {
				return err
			}

			msg := oracle.MsgEditFeed{
				FeedName:          args[0],
				Description:       viper.GetString(FlagDescription),
				LatestHistory:     uint64(viper.GetInt64(FlagLatestHistory)),
				Providers:         providers,
				Timeout:           viper.GetInt64(FlagTimeout),
				ServiceFeeCap:     serviceFeeCap,
				RepeatedFrequency: uint64(viper.GetInt64(FlagFrequency)),
				RepeatedTotal:     viper.GetInt64(FlagTotal),
				ResponseThreshold: uint16(viper.GetInt(FlagThreshold)),
				Creator:           creator,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(FsEditFeed)
	return cmd
}
