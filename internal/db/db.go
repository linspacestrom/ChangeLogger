package db

import (
	"ChangeLogger/internal/config"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// buildPostgresDSN функция для создания DatabaseSourceName
func buildPostgresDSN(cfg config.DbConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DbName)
}

// NewPoolPostgres функция для подключения к Postgres
func NewPoolPostgres(ctx context.Context, cfg config.DbConfig) (*pgxpool.Pool, error) {
	// создаем DatabaseSourceName
	dsn := buildPostgresDSN(cfg)

	// создаем пул соединений для Postgres
	pool, err := pgxpool.New(ctx, dsn)
	// если не удается создать пул соединений, то подключение к Postgres не получится
	if err != nil {
		return nil, fmt.Errorf("create pgx pool: %w", err)
	}

	// создаем контекст с тайм-аутом для пинга Postgres
	pgxCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	// делаем пинг для проверки, что удалось подключиться к Postgres
	if err := pool.Ping(pgxCtx); err != nil {
		// если не удалось подключиться к Postgres, то закрываем пул соединений
		pool.Close()
		return nil, fmt.Errorf("ping to postgres: %w", err)
	}

	return pool, nil

}
