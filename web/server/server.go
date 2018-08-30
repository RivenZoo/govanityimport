package server

import (
	"context"
	"net/http"
	"github.com/gin-gonic/gin"
	"net"
	"time"
	"github.com/RivenZoo/govanityimport/version"
	"github.com/RivenZoo/govanityimport/web/views"
	"github.com/RivenZoo/govanityimport/web/controllers"
	"github.com/gin-contrib/location"
)

func StartWebServer(ctx context.Context, listen string) error {
	err := controllers.InitControllers()
	if err != nil {
		return err
	}
	defer controllers.Close()

	locationCfg := location.Config{
		Scheme: "http",
		Host:   "localhost",
	}
	engine := gin.New()
	engine.Use(gin.Recovery(), gin.Logger(), location.New(locationCfg))
	engine.NoRoute(views.ModuleImportMetaView)
	engine.Any("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, version.GetVersion())
	})

	server := &http.Server{
		Addr:    listen,
		Handler: engine,
	}
	l, err := net.Listen("tcp", listen)
	if err != nil {
		return err
	}
	go func() {
		<-ctx.Done()
		stopCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		server.Shutdown(stopCtx)
	}()
	server.Serve(l)
	return nil
}
