package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/KarinaBotova/Normalization/pkg/server"
)

func main() {
	// Graceful shutdown
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM)
	signal.Notify(sigint, syscall.SIGINT)
	done := make(chan interface{})

	go server.Start(done)

	// Wait interrupt signal
	<-sigint
	close(done)
}
