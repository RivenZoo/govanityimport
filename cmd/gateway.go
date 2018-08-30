package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/RivenZoo/govanityimport/gateway"
	"github.com/RivenZoo/govanityimport/config"
	sigutil "github.com/RivenZoo/govanityimport/signal"
	"github.com/RivenZoo/govanityimport/tracing"
	"github.com/RivenZoo/govanityimport/zaplog"
	"github.com/RivenZoo/govanityimport/proto/apidef"
	"os"
	"fmt"
	"context"
	"google.golang.org/grpc"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

var BackendEndpoint *string

var gatewayCmd = &cobra.Command{
	Use:   "gateway",
	Short: "Start grpc gateway",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.GetConfig()

		err := zaplog.InitLogger(cfg.DeployEnv)
		if err != nil {
			fmt.Fprintln(os.Stderr, "init logger error: ", err)
			os.Exit(-1)
		}
		defer zaplog.Close()

		ctx := sigutil.RegisterDoneSignal()
		if cfg.Debug.Trace {
			go tracing.StartGRPCTraceHTTPServer(ctx, cfg.Debug.GRPCTraceAddress)
		}
		gateway.StartGateway(ctx, cfg.Listen,
			func(ctx context.Context, mux *runtime.ServeMux, opts []grpc.DialOption) error {
				return apidef.RegisterVanityImportServiceHandlerFromEndpoint(ctx, mux, cfg.Gateway.Endpoint, opts)
			})
	},
}

func init() {
	// bind subcommand arguments to config
	gatewayCmd.Flags().StringP("endpoint", "E", "", "Backend grpc endpoint")
	viper.BindPFlag("gateway.endpoint", gatewayCmd.Flags().Lookup("endpoint"))

	rootCmd.AddCommand(gatewayCmd)
}
