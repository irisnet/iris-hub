package keys

import (
	"github.com/irisnet/irishub/client/keys"
	"github.com/spf13/cobra"
)

// CMD

// listKeysCmd represents the list command
var listKeysCmd = &cobra.Command{
	Use:   "list",
	Short: "List all keys",
	Long: `Return a list of all public keys stored by this key manager
along with their associated name and address.`,
	Example: "iriscli keys list",
	RunE: runListCmd,
}

func runListCmd(cmd *cobra.Command, args []string) error {
	kb, err := keys.GetKeyBase()
	if err != nil {
		return err
	}

	infos, err := kb.List()
	if err == nil {
		keys.PrintInfos(cdc, infos)
	}
	return err
}
