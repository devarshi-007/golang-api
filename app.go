// Golang API.
//
//	Schemes: https
//	Host: localhost
//	BasePath: /api/v1
//	Version: 0.0.1-alpha
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package main

import (
	"time"

	"github.com/Improwised/golang-api/cli"
	"github.com/Improwised/golang-api/config"
	"github.com/Improwised/golang-api/logger"
	"github.com/Improwised/golang-api/monitoring"
	"github.com/Improwised/golang-api/routinewrapper"
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
)

func main() {
	// Collecting config from env or file or flag
	cfg := config.GetConfig()

	logger, err := logger.NewRootLogger(cfg.Debug, cfg.IsDevelopment)
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)

	err = monitoring.InitSentry(cfg.Sentry, logger)
	if err != nil {
		logger.Error("Failed to initialize Sentry", zap.Error(err))
	}

	if cfg.Sentry.IsEnabled {
		// Sentry Go routine initialization
		sentryLoggedFunc := func() {
			err := recover()
			if err != nil {
				sentry.CurrentHub().Recover(err)
				sentry.Flush(time.Second * 2)
			}
		}

		routinewrapper.Init(sentryLoggedFunc)
		defer sentryLoggedFunc()

		defer monitoring.CloseSentry(cfg.Sentry, logger)
	}

	err = cli.Init(cfg, logger)
	if err != nil {
		panic(err)
	}

}
