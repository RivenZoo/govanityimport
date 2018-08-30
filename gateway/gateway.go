package gateway

import "time"

import (
	"context"
	"net/http"
	"fmt"
	"google.golang.org/grpc"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc/balancer/roundrobin"
	"github.com/lithammer/shortuuid"
	"google.golang.org/grpc/metadata"
	"github.com/RivenZoo/govanityimport/headers"
	"github.com/RivenZoo/govanityimport/zaplog"
)

type EndpointRegistryFunc func(ctx context.Context, mux *runtime.ServeMux, opts []grpc.DialOption) error

func headerAnnotator(ctx context.Context, request *http.Request) metadata.MD {
	requestID := request.Header.Get(headers.HeaderRequestID)
	if requestID == "" {
		requestID = shortuuid.New()
	}
	return metadata.Pairs(headers.HeaderTimestamp, fmt.Sprintf("%d", time.Now().Unix()),
		headers.HeaderRequestID, requestID)
}

func StartGateway(ctx context.Context, listen string, endpointFunc EndpointRegistryFunc) error {
	marshaler := &runtime.JSONPb{
		OrigName:     true,
		EmitDefaults: true,
	}
	mux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, marshaler),
		runtime.WithMetadata(headerAnnotator))
	dialOpt := []grpc.DialOption{grpc.WithInsecure(), grpc.WithBalancerName(roundrobin.Name)}
	// registry grpc endpoint
	err := endpointFunc(ctx, mux, dialOpt)
	if err != nil {
		zaplog.GetSugarLogger().Errorw("registry endpoint fail", "error", err)
		return err
	}

	server := http.Server{
		Addr:    listen,
		Handler: mux,
	}
	go func() {
		select {
		case <-ctx.Done():
			func() {
				tmCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
				defer cancel()
				server.Shutdown(tmCtx)
			}()
		}
	}()
	return server.ListenAndServe()
}
