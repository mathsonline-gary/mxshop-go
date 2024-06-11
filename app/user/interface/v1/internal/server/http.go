package server

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/zycgary/mxshop-go/app/user/interface/v1/internal/config"
	"github.com/zycgary/mxshop-go/app/user/interface/v1/internal/service"
	"github.com/zycgary/mxshop-go/pkg/log"
	"github.com/zycgary/mxshop-go/pkg/transport"
)

var _ transport.Server = (*httpServer)(nil)

type httpServer struct {
	engine *gin.Engine
	Host   string
	Port   uint32
	us     service.UserService
	logger *log.Sugar
}

func NewHttpServer(conf *config.Config, us *service.UserService, logger log.Logger) transport.Server {
	s := &httpServer{
		Host:   conf.Server.HTTP.Host,
		Port:   conf.Server.HTTP.Port,
		us:     *us,
		logger: log.NewSugar(logger),
	}

	if conf.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})
	rg := r.Group("/v1")
	users := rg.Group("/users")
	users.GET("/", s.us.Index)

	s.engine = r

	return s
}

func (s *httpServer) Start(_ context.Context) error {
	s.logger.Debugf("Starting HTTP server at %s:%d", s.Host, s.Port)
	return s.engine.Run(fmt.Sprintf("%s:%d", s.Host, s.Port))
}

func (s *httpServer) Stop(_ context.Context) error {
	s.logger.Debugf("Stopping HTTP server")

	// TODO: Graceful shutdown
	return nil
}
