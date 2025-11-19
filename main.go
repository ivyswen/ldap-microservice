package main

import (
	"context"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	cfg := LoadConfigFromEnv()

	// Configure logger based on config
	logLevel := parseLogLevel(cfg.LogLevel)
	zerolog.SetGlobalLevel(logLevel)

	// Setup multi-writer: console + file (if configured)
	writers := []io.Writer{zerolog.ConsoleWriter{Out: os.Stdout}}

	if cfg.LogFile != "" {
		fileWriter := &lumberjack.Logger{
			Filename:   cfg.LogFile,
			MaxSize:    100, // megabytes
			MaxBackups: 3,
			MaxAge:     28, // days
			Compress:   true,
		}
		writers = append(writers, fileWriter)
	}

	multi := zerolog.MultiLevelWriter(writers...)
	log.Logger = zerolog.New(multi).With().Timestamp().Logger()

	log.Info().Interface("config", cfg.ToMap()).Msg("starting ldap microservice")

	router := mux.NewRouter()
	// routes
	// 使用配置的 BasePath 前缀注册路由
	basePath := cfg.BasePath
	router.HandleFunc(basePath+"/v1/auth", AuthHandler(cfg)).Methods("POST")
	router.HandleFunc(basePath+"/v1/healthz", HealthHandler).Methods("GET")
	router.HandleFunc(basePath+"/v1/readyz", ReadyHandler).Methods("GET")

	srv := &http.Server{
		Addr:         ":" + cfg.ServicePort,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// graceful shutdown
	go func() {
		log.Info().Msgf("listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server failed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Info().Msg("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("server shutdown error")
	}
	log.Info().Msg("server exited")
}

func parseLogLevel(level string) zerolog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	default:
		return zerolog.InfoLevel
	}
}
