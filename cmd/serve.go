package cmd

import (
	"github.com/spf13/cobra"
	"govanityimport/config"
	"govanityimport/zaplog"
	"govanityimport/tracing"
	"govanityimport/rpcserver"
	"govanityimport/controllers"
	"govanityimport/authorize"
	"govanityimport/proto/apidef"
	sigutil "govanityimport/signal"
	"os"
	"fmt"
	"google.golang.org/grpc"
)

var serveCommand = &cobra.Command{
	Use:   "serve",
	Short: "Start go vanity import service",
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

		authorize.InitToken(cfg.Authorize.RawTokens)
		err = rpcserver.StartRPCServer(ctx, cfg.Listen, func(s *grpc.Server) {
			apidef.RegisterVanityImportServiceServer(s, controllers.GetController())
		})
		if err != nil {
			zaplog.GetSugarLogger().Errorw("start rpc server fail", "error", err)
			os.Exit(-1)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCommand)
}
