package server

import (
	"fmt"

	"github.com/irisnet/irishub/codec"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagHeight = "height"
)

// ResetCmd reset app state to particular height
func ResetCmd(ctx *Context, cdc *codec.Codec, appReset AppReset) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reset",
		Short: "Reset app state to the specified height",
		RunE: func(cmd *cobra.Command, args []string) error {
			home := viper.GetString("home")
			traceWriterFile := viper.GetString(flagTraceStore)
			emptyState, err := isEmptyState(home)
			if err != nil {
				return err
			}

			if emptyState {
				fmt.Println("WARNING: State is not initialized.")
				return nil
			}

			db, err := openDB(home)
			if err != nil {
				return err
			}
			traceWriter, err := openTraceWriter(traceWriterFile)
			if err != nil {
				return err
			}
			height := viper.GetInt64(flagHeight)
			if height < 0 {
				return errors.Errorf("Height must greater than or equal to zero")
			}
			err = appReset(ctx, ctx.Logger, db, traceWriter, height)
			if err != nil {
				return errors.Errorf("Error reset state: %v\n", err)
			}

			fmt.Printf("Reset app state to height %d successfully\n", height)
			return nil
		},
	}
	cmd.Flags().Int64(flagHeight, 0, "Reset state from a particular height (0 or greater than latest height means latest height)")
	cmd.MarkFlagRequired(flagHeight)
	return cmd
}
