package lcd

import (
	"os"

	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/client"
	bankHandler "github.com/irisnet/irishub/client/bank/lcd"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/keys"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	tmserver "github.com/tendermint/tendermint/rpc/lib/server"
	"net/http"
)

// ServeLCDStartCommand will start irislcd node, which provides rest APIs with swagger-ui
func ServeLCDStartCommand(cdc *wire.Codec) *cobra.Command {
	flagListenAddr := "laddr"
	flagCORS := "cors"
	flagMaxOpenConnections := "max-open"

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start irislcd (irishub light-client daemon), a local REST server with swagger-ui: http://localhost:1317/swagger-ui/",
		RunE: func(cmd *cobra.Command, args []string) error {
			listenAddr := viper.GetString(flagListenAddr)
			router := createHandler(cdc)

			statikFS, err := fs.New()
			if err != nil {
				panic(err)
			}
			staticServer := http.FileServer(statikFS)
			router.PathPrefix("/swagger-ui/").Handler(http.StripPrefix("/swagger-ui/", staticServer))

			logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "irislcd")
			maxOpen := viper.GetInt(flagMaxOpenConnections)

			listener, err := tmserver.StartHTTPServer(
				listenAddr, router, logger,
				tmserver.Config{MaxOpenConnections: maxOpen},
			)
			if err != nil {
				return err
			}

			logger.Info("irislcd server started")

			// wait forever and cleanup
			cmn.TrapSignal(func() {
				err := listener.Close()
				logger.Error("error closing listener", "err", err)
			})

			return nil
		},
	}

	cmd.Flags().String(flagListenAddr, "tcp://localhost:1317", "The address for the server to listen on")
	cmd.Flags().String(flagCORS, "", "Set the domains that can make CORS requests (* for all)")
	cmd.Flags().String(client.FlagChainID, "", "The chain ID to connect to")
	cmd.Flags().String(client.FlagNode, "tcp://localhost:26657", "Address of the node to connect to")
	cmd.Flags().Int(flagMaxOpenConnections, 1000, "The number of maximum open connections")
	cmd.Flags().Bool(client.FlagTrustNode, false, "Don't verify proofs for responses")

	return cmd
}

func createHandler(cdc *wire.Codec) *mux.Router {
	r := mux.NewRouter()
	kb, err := keys.GetKeyBase()
	if err != nil {
		panic(err)
	}
	cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout)

	r.HandleFunc("/version", CLIVersionRequestHandler).Methods("GET")
	r.HandleFunc("/node_version", NodeVersionRequestHandler(cliCtx)).Methods("GET")

	bankHandler.RegisterRoutes(cliCtx, r, cdc, kb)

	/*
		keys.RegisterRoutes(r)
		rpc.RegisterRoutes(cliCtx, r)
		tx.RegisterRoutes(cliCtx, r, cdc)

		auth.RegisterRoutes(cliCtx, r, cdc, "acc")
		bank.RegisterRoutes(cliCtx, r, cdc, kb)
		ibc.RegisterRoutes(cliCtx, r, cdc, kb)
		stake.RegisterRoutes(cliCtx, r, cdc, kb)
		slashing.RegisterRoutes(cliCtx, r, cdc, kb)
		gov.RegisterRoutes(cliCtx, r, cdc)
	*/
	return r
}
