package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

type (
	config struct {
		port int
		env string
	}

	application struct {
		config config
		logger *slog.Logger
	}
)

func main() {
	var cfg config

	/*
	Get vars passed in as flags and set on a struct to be read later. Defaults provided
	*/
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	// Create the logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Create the application. Could have embedded the config, but we want to use DI to access these things really
	app := &application{
		config: cfg,
		logger: logger,
	}

	// Mux behaves as a http handler func so is used in server
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/healthcheck", app.healthcheckHandler)

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", cfg.port),
		Handler: mux,
		IdleTimeout: time.Minute,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "addr", server.Addr, "env", cfg.env)

	err := server.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}
