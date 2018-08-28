package rpcserver

import (
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"net"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"govanityimport/authorize"
	"govanityimport/zaplog"
	"govanityimport/interceptor"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RPCRegistryFunc func(s *grpc.Server)

func StartRPCServer(ctx context.Context, listen string, registryFunc RPCRegistryFunc) error {
	log := zaplog.GetSugarLogger()
	server := grpc.NewServer(grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			log.Errorw("server panic", "error", p)
			return status.Errorf(codes.Internal, "%s", p)
		})),
		grpc_auth.UnaryServerInterceptor(authorize.AuthorizeToken),
		interceptor.UnaryRequestInfoInterceptor(nil, nil),
	)))
	// registry rpc handler
	registryFunc(server)

	l, err := net.Listen("tcp", listen)
	if err != nil {
		return err
	}
	go func() {
		<-ctx.Done()
		server.GracefulStop()
	}()
	return server.Serve(l)
}
