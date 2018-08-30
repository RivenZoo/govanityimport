package cmd

import (
	"github.com/spf13/cobra"
	"github.com/RivenZoo/govanityimport/config"
	"github.com/RivenZoo/govanityimport/zaplog"
	"github.com/RivenZoo/govanityimport/tracing"
	"github.com/RivenZoo/govanityimport/rpcserver"
	"github.com/RivenZoo/govanityimport/controllers"
	"github.com/RivenZoo/govanityimport/authorize"
	"github.com/RivenZoo/govanityimport/proto/apidef"
	"github.com/RivenZoo/govanityimport/model"
	sigutil "github.com/RivenZoo/govanityimport/signal"
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
		model.InitModel()
		defer model.Close()

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
