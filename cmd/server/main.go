package main

import (
	"ChangeLogger/internal/config"
	"ChangeLogger/internal/db"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Println("ChangeLogger started")

	// Обработка контекста с отменой
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Обработка сигналов по graceful shutdown
	signChan := make(chan os.Signal, 1)
	signal.Notify(signChan, syscall.SIGINT, syscall.SIGINT)
	defer close(signChan)

	// загрузка конфига из файла .env
	cfg, errConfig := config.LoadConfig()
	if errConfig != nil {
		log.Printf("Config %s\n", errConfig)
	}

	log.Println("Config successes loaded")

	// устанавливаем пул соединений с базой данных
	pgxPool, errPgxPool := db.NewPoolPostgres(ctx, cfg.DbConfig)
	if errPgxPool != nil {
		log.Printf("%s\n", errPgxPool)
		cancel()
	} else {
		// закрываем пул соединений с базой данных
		defer pgxPool.Close()
	}
	log.Println("Connection to the database is established")

	// обработчик сигналов
	select {
	case <-ctx.Done(): // завершаем работу по контексту
		log.Println("Context canceling, shutting down")
	case val := <-signChan: // завершаем работу по graceful shutdown (ctrl+c)
		log.Printf("\nRecieved signal %s, shutting down", val)
		cancel()
	}

	log.Println("ChangeLogger stopped")
}
