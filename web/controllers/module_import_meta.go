package controllers

import (
	"github.com/RivenZoo/govanityimport/proto/apidef"
	"github.com/RivenZoo/govanityimport/config"
	"github.com/RivenZoo/govanityimport/errorcode"
	"github.com/RivenZoo/govanityimport/zaplog"
	"github.com/RivenZoo/govanityimport/headers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc/metadata"
	"fmt"
)

var (
	instance *controller
)

type controller struct {
	queryClient apidef.VanityImportServiceClient
	conn        *grpc.ClientConn
}

func InitControllers() error {
	cfg := config.GetConfig()
	conn, err := grpc.Dial(cfg.Web.ModuleServiceAddress,
		grpc.WithBalancerName(roundrobin.Name),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor()))
	if err != nil {
		return err
	}
	instance = &controller{
		conn:        conn,
		queryClient: apidef.NewVanityImportServiceClient(conn),
	}
	return nil
}

func GetController() *controller {
	return instance
}

func Close() error {
	if instance != nil {
		instance.conn.Close()
	}
	return nil
}

func (c *controller) addBearerAuthToken(ctx context.Context) (context.Context) {
	cfg := config.GetConfig()
	return metadata.AppendToOutgoingContext(ctx, headers.ContextKeyAuthorization,
		fmt.Sprintf("bearer %s", cfg.Web.AuthToken))
}

func (c *controller) GetModuleMetaInfo(ctx context.Context, importPath string) (*apidef.ModuleMetaInfo, error) {
	log := zaplog.GetSugarLogger()
	ctx = c.addBearerAuthToken(ctx)

	resp, err := c.queryClient.QueryImportMetaInfo(ctx, &apidef.ImportMetaInfoReq{
		ImportPath: importPath,
	}, grpc_retry.WithMax(3))
	if err != nil {
		log.Errorw("query module meta error", "error", err)
		return nil, err
	}
	if resp.Ret != int32(errorcode.OK.Ret) {
		log.Errorw("query module meta fail", "ret", resp.Ret,
			"msg", resp.Msg, "traceid", resp.TraceId)
		return nil, errorcode.ErrInnerQueryModuleError
	}
	return resp.MetaInfo, nil
}
