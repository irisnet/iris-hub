package cli

import (
	"encoding/hex"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/irisnet/irishub/client"
	"github.com/irisnet/irishub/client/context"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/common"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// QueryTxCmd implements the default command for a tx query.
func QueryTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tx [hash]",
		Short: "Matches this txhash over all committed blocks",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// find the key to look up the account
			hashHexStr := args[0]

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			output, err := queryTx(cdc, cliCtx, hashHexStr)
			if err != nil {
				return err
			}

			fmt.Println(string(output))
			return nil
		},
	}

	cmd.Flags().StringP(client.FlagNode, "n", "tcp://localhost:26657", "Node to connect to")
	cmd.Flags().Bool(client.FlagTrustNode, false, "Trust connected full node (don't verify proofs for responses)")
	cmd.Flags().String(client.FlagChainID, "", "Chain ID of Tendermint node")
	return cmd
}

func queryTx(cdc *wire.Codec, cliCtx context.CLIContext, hashHexStr string) ([]byte, error) {
	hash, err := hex.DecodeString(hashHexStr)
	if err != nil {
		return nil, err
	}

	node, err := cliCtx.GetNode()
	if err != nil {
		return nil, err
	}

	res, err := node.Tx(hash, !cliCtx.TrustNode)
	if err != nil {
		return nil, err
	}

	if !cliCtx.TrustNode {
		err := ValidateTxResult(cliCtx, res)
		if err != nil {
			return nil, err
		}
	}

	info, err := formatTxResult(cdc, res)
	if err != nil {
		return nil, err
	}

	return cdc.MarshalJSONIndent(info, "", "  ")
}

// ValidateTxResult performs transaction verification
func ValidateTxResult(cliCtx context.CLIContext, res *ctypes.ResultTx) error {
	check, err := cliCtx.Certify(res.Height)
	if err != nil {
		return err
	}

	err = res.Proof.Validate(check.Header.DataHash)
	if err != nil {
		return err
	}
	return nil
}

func formatTxResult(cdc *wire.Codec, res *ctypes.ResultTx) (Info, error) {
	tx, err := parseTx(cdc, res.Tx)
	if err != nil {
		return Info{}, err
	}

	return Info{
		Hash:   res.Hash,
		Height: res.Height,
		Tx:     tx,
		Result: res.TxResult,
	}, nil
}

// Info is used to prepare info to display
type Info struct {
	Hash   common.HexBytes        `json:"hash"`
	Height int64                  `json:"height"`
	Tx     sdk.Tx                 `json:"tx"`
	Result abci.ResponseDeliverTx `json:"result"`
}

func parseTx(cdc *wire.Codec, txBytes []byte) (sdk.Tx, error) {
	var tx auth.StdTx

	err := cdc.UnmarshalBinary(txBytes, &tx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
