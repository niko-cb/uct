package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/niko-cb/uct/internal/infrastructure/monitor/log"
	"github.com/niko-cb/uct/internal/infrastructure/web/config"
)

type server struct {
	*echo.Echo
}

func NewServer() *server {
	e := echo.New()
	var a []echo.MiddlewareFunc
	e.GET("/", nil, a...)
	return &server{
		Echo: e,
	}
}

// Run starts the server with graceful shutdown in mind
func (s *server) Run() {
	cfg := s.getServerConfig()
	s.routing()
	s.CORS()
	s.Auth(cfg.JwtSecret)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		if err := s.Start(fmt.Sprintf(":%s", cfg.Port)); err != nil {
			log.Fatal(ctx, fmt.Errorf("server start error: %+v", err))
		}
	}()
	<-ctx.Done()

	s.GracefulShutdown(ctx)
}

// GracefulShutdown handles the graceful shutdown process
func (s *server) GracefulShutdown(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatal(ctx, fmt.Errorf("server shutdown error: %+v", err))
	}
}

func (s *server) getServerConfig() *config.Config {
	err := config.Parse()
	if err != nil {
		log.Fatal(context.Background(), err)

	}
	return &config.Cfg
}
