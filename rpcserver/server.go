package rpcserver

import (
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"net"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"govanityimport/authorize"
)

type RPCRegistryFunc func(s *grpc.Server)

func StartRPCServer(ctx context.Context, listen string, registryFunc RPCRegistryFunc) error {
	server := grpc.NewServer(grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		grpc_recovery.UnaryServerInterceptor(),
		grpc_auth.UnaryServerInterceptor(authorize.AuthorizeToken),
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
