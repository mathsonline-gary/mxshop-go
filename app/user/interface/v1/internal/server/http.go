package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zycgary/mxshop-go/app/user/interface/v1/internal/config"
	"github.com/zycgary/mxshop-go/app/user/interface/v1/internal/service"
	"github.com/zycgary/mxshop-go/pkg/log"
	"github.com/zycgary/mxshop-go/pkg/transport"
)

var _ transport.Server = (*httpServer)(nil)

type httpServer struct {
	engine *gin.Engine
	server *http.Server
	Host   string
	Port   uint32
	us     service.UserService
	as     service.AuthService
	logger *log.Sugar
}

func NewHttpServer(conf *config.Config, us *service.UserService, as *service.AuthService, logger log.Logger) transport.Server {
	s := &httpServer{
		Host:   conf.Server.HTTP.Host,
		Port:   conf.Server.HTTP.Port,
		us:     *us,
		as:     *as,
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

	// Login route
	rg.POST("/login", s.as.Login)

	// User routes
	users := rg.Group("/users")
	users.GET("/", s.us.Index)

	s.engine = r

	return s
}

func (s *httpServer) Start(_ context.Context) error {
	s.logger.Infof("Starting HTTP server at %s:%d", s.Host, s.Port)

	s.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.Host, s.Port),
		Handler: s.engine.Handler(),
	}

	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *httpServer) Stop(ctx context.Context) error {
	s.logger.Infof("Stopping HTTP server")

	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
