package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/linspacestrom/ChangeLogger/internal/config"
	"github.com/linspacestrom/ChangeLogger/internal/db"
	"github.com/linspacestrom/ChangeLogger/internal/handlers"
	"github.com/linspacestrom/ChangeLogger/internal/repositories"
	"github.com/linspacestrom/ChangeLogger/internal/services"
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

	repo := repositories.NewPoolProjectRepository(pgxPool)
	svc := services.NewProjectService(repo)
	hs := handlers.NewRoutes(svc)

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: hs.Handler(),
	}

	go func() {
		log.Printf("HTTP server listening on :%s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %s", err)
		}
	}()

	// обработчик сигналов
	select {
	case <-ctx.Done(): // завершаем работу по контексту
		log.Println("Context canceling, shutting down")
	case val := <-signChan: // завершаем работу по graceful shutdown (ctrl+c)
		log.Printf("\nRecieved signal %s, shutting down", val)
		cancel()
	}

	_ = srv.Close()

	log.Println("ChangeLogger stopped")
}
