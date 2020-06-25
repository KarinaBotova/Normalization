package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/kolya59/easy_normalization/pkg/client"
)

func main() {
	// создание канала
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM)//когда появится сигнал от системы
	signal.Notify(sigint, syscall.SIGINT)
	done := make(chan interface{}) //для внутренней логики
    // запускаем функцию в асинхронном режиме
	go client.Start(done)

	// ждем сообщения от системы
	<-sigint
	close(done)
}
