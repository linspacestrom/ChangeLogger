package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string
	DbConfig
}

type DbConfig struct {
	Host     string
	Port     string
	DbName   string
	User     string
	Password string
}

// getEnv функция для загрузки переменной из .env
// если переменной нет в .env, то возвращается дефолтное значение
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// LoadConfig функция для загрузки всех переменных окружения из .env
func LoadConfig() (*Config, error) {
	// загружаем .env
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	return &Config{
		ServerPort: getEnv("PORT", "8000"),
		DbConfig: DbConfig{
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnv("POSTGRES_PORT", "5432"),
			DbName:   getEnv("POSTGRES_NAME", "changelogger_db"),
			User:     getEnv("POSTGRES_USER", "postgres"),
			Password: getEnv("POSTGRES_PASSWORD", "postgres"),
		},
	}, nil
}
