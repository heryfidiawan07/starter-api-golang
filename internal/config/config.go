package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	DB       DBConfig
	JWT      JWTConfig
	Mail     MailConfig
	Google   OAuthConfig
	Facebook OAuthConfig
	Storage  StorageConfig

	EmailVerificationRequired bool
	ResetTokenExpire          int
}

type AppConfig struct {
	Env  string
	Port string
	URL  string
}

type DBConfig struct {
	Driver  string
	Host    string
	Port    string
	User    string
	Pass    string
	Name    string
	SSLMode string
}

type JWTConfig struct {
	Secret        string
	AccessExpire  int
	RefreshExpire int
}

type MailConfig struct {
	Host     string
	Port     int
	User     string
	Pass     string
	From     string
	FromName string
}

type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type StorageConfig struct {
	Path string
	URL  string
}

var cfg *Config

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg = &Config{
		App: AppConfig{
			Env:  getEnv("APP_ENV", "development"),
			Port: getEnv("APP_PORT", "8080"),
			URL:  getEnv("APP_URL", "http://localhost:8080"),
		},
		DB: DBConfig{
			Driver:  getEnv("DB_DRIVER", "mysql"),
			Host:    getEnv("DB_HOST", "127.0.0.1"),
			Port:    getEnv("DB_PORT", "3306"),
			User:    getEnv("DB_USER", "root"),
			Pass:    getEnv("DB_PASS", ""),
			Name:    getEnv("DB_NAME", "starter_db"),
			SSLMode: getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:        getEnv("JWT_SECRET", "default-secret"),
			AccessExpire:  getEnvInt("JWT_ACCESS_EXPIRE", 15),
			RefreshExpire: getEnvInt("JWT_REFRESH_EXPIRE", 10080),
		},
		Mail: MailConfig{
			Host:     getEnv("MAIL_HOST", ""),
			Port:     getEnvInt("MAIL_PORT", 587),
			User:     getEnv("MAIL_USER", ""),
			Pass:     getEnv("MAIL_PASS", ""),
			From:     getEnv("MAIL_FROM", ""),
			FromName: getEnv("MAIL_FROM_NAME", "Starter API"),
		},
		Google: OAuthConfig{
			ClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
			ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
			RedirectURL:  getEnv("GOOGLE_REDIRECT_URL", ""),
		},
		Facebook: OAuthConfig{
			ClientID:     getEnv("FACEBOOK_CLIENT_ID", ""),
			ClientSecret: getEnv("FACEBOOK_CLIENT_SECRET", ""),
			RedirectURL:  getEnv("FACEBOOK_REDIRECT_URL", ""),
		},
		Storage: StorageConfig{
			Path: getEnv("STORAGE_PATH", "./storage/photos"),
			URL:  getEnv("STORAGE_URL", "http://localhost:8080/storage/photos"),
		},
		EmailVerificationRequired: getEnvBool("EMAIL_VERIFICATION_REQUIRED", true),
		ResetTokenExpire:          getEnvInt("RESET_TOKEN_EXPIRE", 60),
	}

	return cfg
}

func Get() *Config {
	if cfg == nil {
		return Load()
	}
	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if v := os.Getenv(key); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return fallback
}
