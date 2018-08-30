package authorize

import (
	"context"
	"google.golang.org/grpc/metadata"
	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/RivenZoo/govanityimport/zaplog"
	"github.com/RivenZoo/govanityimport/headers"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
)

var (
	staticTokens = map[string]bool{}
)

func InitToken(tokens []string) {

	for _, token := range tokens {
		staticTokens[token] = true
	}
}

func AuthorizeToken(ctx context.Context) (context.Context, error) {
	md, _ := metadata.FromIncomingContext(ctx)

	log := zaplog.GetSugarLogger()
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		log.Errorw("no auth token", "request-id", md.Get(headers.HeaderRequestID),
			"request-timestamp", md.Get(headers.HeaderTimestamp))
		return nil, status.Errorf(codes.Unauthenticated, "lack of auth token")
	}
	_, ok := staticTokens[token]
	if !ok {
		log.Errorw("invalid token", "request-id", md.Get(headers.HeaderRequestID),
			"request-timestamp", md.Get(headers.HeaderTimestamp))
		return nil, status.Errorf(codes.PermissionDenied, "invalid auth token")
	}
	return ctx, nil
}
