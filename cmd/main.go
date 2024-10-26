package main

import (
	"github.com/gin-gonic/gin"

	"github.com/vindosVP/snapigw/cmd/config"
	"github.com/vindosVP/snapigw/internal/server"
	"github.com/vindosVP/snapigw/internal/services/auth"
	"github.com/vindosVP/snapigw/pkg/logger"
)

var (
	buildCommit = "N/A"
	buildTime   = "N/A"
	version     = "N/A"
)

func main() {
	cfg := config.MustParse()
	l := logger.SetupLogger(cfg.ENV, cfg.ServiceName)

	l.Info().Str("env", cfg.ENV).
		Str("buildCommit", buildCommit).
		Str("buildTime", buildTime).
		Str("version", version).
		Msg("starting service")

	l.Info().Interface("config", cfg).Msg("configuration loaded")

	if cfg.ENV == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	pxs := server.NewProxs()
	ap, err := auth.NewProxy(cfg.Services.AuthAddr, l)
	if err != nil {
		l.Fatal().Err(err).Stack().Msg("failed to create auth proxy")
	}
	pxs.WithAuth(ap)

	s := server.NewServer(cfg.Port, l)
	s.WithProxs(pxs)
	s.SetRouter()
	s.Run()
}
