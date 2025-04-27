package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort          string
	PostgresDSN       string
	RedisAddr         string
	KafkaBrokers      []string
	JWTSecret         []byte
	AccessTTLSeconds  int64
	RefreshTTLSeconds int64
}

func Load() Config {
	_ = godotenv.Load("../infrastructure/auth.env", "./auth.env")

	atoi := func(key string, def int64) int64 {
		if v, err := strconv.ParseInt(os.Getenv(key), 10, 64); err == nil {
			return v
		}
		return def
	}

	cfg := Config{
		HTTPPort:          env("AUTH_HTTP_PORT", "8080"),
		PostgresDSN:       env("POSTGRES_DSN", "postgres://auth:auth@postgres:5432/auth"),
		RedisAddr:         env("REDIS_ADDR", "redis:6379"),
		KafkaBrokers:      split(env("KAFKA_BROKERS", "kafka:9092")),
		JWTSecret:         []byte(env("JWT_SECRET", "super-secret-key")),
		AccessTTLSeconds:  atoi("ACCESS_TTL", 30*60),
		RefreshTTLSeconds: atoi("REFRESH_TTL", 7*24*60*60),
	}
	log.Printf("config loaded: %+v", cfg)
	return cfg
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
func split(s string) []string { return strings.Split(s, ",") }
