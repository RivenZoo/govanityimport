/*
 * user define interceptor
 */
package interceptor

import (
	"google.golang.org/grpc"
	"context"
	"google.golang.org/grpc/metadata"
	"github.com/lithammer/shortuuid"
	"govanityimport/headers"
	"govanityimport/zaplog"
	"fmt"
	"time"
)

type BeforeHandlerFunc func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo) context.Context
type AfterHandlerFunc func(ctx context.Context, req interface{}, resp interface{}, err error)

// Default BeforeHandlerFunc
// Extract request information from metadata then set them to context
func RequestInfoExtract(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo) context.Context {
	md, _ := metadata.FromIncomingContext(ctx)

	var ts, requestID string
	arr := md.Get(headers.HeaderTimestamp)
	if len(arr) == 0 || arr[0] == "" {
		ts = fmt.Sprintf("%d", time.Now().Unix())
	} else {
		ts = arr[0]
	}
	ctx = context.WithValue(ctx, headers.HeaderTimestamp, ts)

	arr = md.Get(headers.HeaderRequestID)
	if len(arr) == 0 || arr[0] == "" {
		requestID = shortuuid.New()
	} else {
		requestID = arr[0]
	}
	ctx = context.WithValue(ctx, headers.HeaderRequestID, requestID)
	ctx = context.WithValue(ctx, headers.ContextKeyRPCMethod, info.FullMethod)
	return ctx
}

// Default AfterHandlerFunc
// Log request and response information if error not nil
func PostHandler(ctx context.Context, req interface{}, resp interface{}, err error) {
	requestID := ctx.Value(headers.HeaderRequestID).(string)
	ts := ctx.Value(headers.HeaderTimestamp).(string)
	method := ctx.Value(headers.ContextKeyRPCMethod).(string)
	log := zaplog.GetSugarLogger()
	if err != nil {
		log.Errorw("request fail", "request_id", requestID, "method", method, "request_ts", ts,
			"request", req, "response", resp)
	} else {
		log.Infow("request ok", "request_id", requestID, "method", method, "request_ts", ts)
	}
	return
}

// UnaryRequestInfoInterceptor returns a unary interceptor to add before/post handler for request.
// handlerFunc: if nil use default RequestInfoExtract
// afterHandlerFunc: if nil use default PostHandler
func UnaryRequestInfoInterceptor(handlerFunc BeforeHandlerFunc, afterHandlerFunc AfterHandlerFunc) grpc.UnaryServerInterceptor {
	if handlerFunc == nil {
		handlerFunc = RequestInfoExtract
	}
	if afterHandlerFunc == nil {
		afterHandlerFunc = PostHandler
	}
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		newCtx := handlerFunc(ctx, req, info)

		resp, err := handler(newCtx, req)

		afterHandlerFunc(newCtx, req, resp, err)
		return resp, err
	}
}

