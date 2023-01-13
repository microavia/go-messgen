package main

import (
	"os"

	"github.com/powerman/structlog"

	"github.com/microavia/go-messgen/internal/config"
)

func main() {
	structlog.DefaultLogger.SetLogLevel(structlog.INF)

	cfg, err := config.Parse(os.Args[1:])
	if err != nil {
		structlog.DefaultLogger.Fatal(err)
	}

	if *cfg.Verbose {
		structlog.DefaultLogger.SetLogLevel(structlog.DBG)
	}

	structlog.DefaultLogger.Debug("started", "config", cfg)
}
