package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	// TODO: Use enviromnment variables here
	slog.SetLogLoggerLevel(slog.LevelDebug)

	s := NewServer(Config{
		ListenAddr: ":6379",
	})

	go gracefulShutdown(cancel)
	err := s.Start(ctx)
	if err != nil {
		panic(err)
	}
}

func gracefulShutdown(cancel context.CancelFunc) {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)

	<-sc
	cancel()
}
