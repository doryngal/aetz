package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Config struct {
	Server struct {
		Port string
	}
	Database struct {
		Host     string
		Port     int
		User     string
		Password string
		DbName   string
	}
	Migrations struct {
		Path   string
		DbName string
	}
	SMTP struct {
		Host     string
		Port     int
		User     string
		Password string
		From     string
	}
	OutputDir string
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return fallback
}

func InitConfig(envPath string) (*Config, error) {
	err := godotenv.Load(envPath)
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	cfg := &Config{}
	cfg.Server.Port = getEnv("SERVER_PORT", "8080")

	cfg.Database.Host = getEnv("DB_HOST", "localhost")
	cfg.Database.Port = getEnvAsInt("DB_PORT", 5432)
	cfg.Database.User = getEnv("DB_USER", "postgres")
	cfg.Database.Password = getEnv("DB_PASSWORD", "349349")
	cfg.Database.DbName = getEnv("DB_NAME", "TamaqQr")

	cfg.Migrations.Path = getEnv("MIGRATIONS_PATH", "./migrations")
	cfg.Migrations.DbName = getEnv("MIGRATIONS_DB_NAME", "binai")

	cfg.SMTP.Host = getEnv("SMTP_HOST", "smtp.gmail.com")
	cfg.SMTP.Port = getEnvAsInt("SMTP_PORT", 587)
	cfg.SMTP.User = getEnv("SMTP_USER", "")
	cfg.SMTP.Password = getEnv("SMTP_PASSWORD", "")
	cfg.SMTP.From = getEnv("SMTP_FROM", "example@gmail.com")

	cfg.OutputDir = getEnv("OUTPUT_DIR", "./output")
	return cfg, nil
}

func ConnectDB(cfg *Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.DbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Println("Successfully connected to the database")
	return db, nil
}
