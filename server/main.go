package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"redis-from-scratch/server"

	"github.com/joho/godotenv"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	port := os.Getenv("PORT")
	if port == "" {
		port = "6379"
	}

	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	level := os.Getenv("LOG_LEVEL")
	setLogLoggerLevel(level)

	s := server.NewServer(
		server.Config{
			ListenAddr: ":" + port,
		})

	go gracefulShutdown(cancel)

	err = s.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func gracefulShutdown(cancel context.CancelFunc) {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)

	<-sc
	cancel()
}

func setLogLoggerLevel(level string) {
	switch strings.ToLower(level) {
	case "debug":
		slog.SetLogLoggerLevel(slog.LevelDebug)
	case "info":
		slog.SetLogLoggerLevel(slog.LevelInfo)
	case "error":
		slog.SetLogLoggerLevel(slog.LevelError)
	case "warn":
		slog.SetLogLoggerLevel(slog.LevelWarn)
	default:
		slog.SetLogLoggerLevel(slog.LevelInfo)
	}
}
