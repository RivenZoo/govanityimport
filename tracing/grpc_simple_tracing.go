package tracing

import (
	"net/http"
	"context"
	"time"
	"golang.org/x/net/trace"
	"google.golang.org/grpc"
)

func StartGRPCTraceHTTPServer(ctx context.Context, listen string) {
	grpc.EnableTracing = true
	trace.AuthRequest = func(req *http.Request) (any, sensitive bool) {
		return true, true
	}
	s := &http.Server{
		Addr: listen,
	}
	ch := make(chan struct{})
	go func() {
		close(ch)
		select {
		case <-ctx.Done():
			stopCtx, fn := context.WithTimeout(context.Background(), time.Second*5)
			defer fn()
			s.Shutdown(stopCtx)
		}
	}()
	<-ch
	s.ListenAndServe()
}
