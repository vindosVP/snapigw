package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/vindosVP/snapigw/internal/middleware"
	"github.com/vindosVP/snapigw/internal/services/auth"
)

type Server struct {
	l      zerolog.Logger
	port   int
	router *gin.Engine
	proxs  *Proxs
}

func (s *Server) WithProxs(proxs *Proxs) *Server {
	s.proxs = proxs
	return s
}

func (s *Server) Run() {

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: s.router,
	}

	s.l.Info().Str("addr", srv.Addr).Msg("starting server")
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.l.Fatal().Err(err).Stack().Msg("failed to start server")
		}
	}()

	quit := make(chan os.Signal, 2)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	s.l.Info().Msg("shutting down gracefully")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		s.l.Fatal().Err(err).Stack().Msg("failed to shut down gracefully")
	}

	select {
	case <-ctx.Done():
		s.l.Info().Msg("server stopped")
	}
}

type Proxs struct {
	auth *auth.Proxy
}

func (p *Proxs) WithAuth(auth *auth.Proxy) *Proxs {
	p.auth = auth
	return p
}

func NewProxs() *Proxs {
	return &Proxs{}
}

func NewServer(port int, l zerolog.Logger) *Server {
	s := &Server{port: port, l: l}
	return s
}

func (s *Server) SetRouter() {
	r := gin.Default()
	r.Use(middleware.RequestId())
	r.POST("/api/users/register", s.proxs.auth.RegisterHandler())
	r.POST("/api/users/login", s.proxs.auth.LoginHandler())
	r.POST("/api/users/refresh", s.proxs.auth.RefreshHandler())
	s.router = r
}
