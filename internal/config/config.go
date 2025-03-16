package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	Database     `env-prefix:"DB_"`
	Server       `env-prefix:"SERVER_"`
	JWTSecretKey string `env:"JWT_SECRET_KEY" env-required:"true"`
	BotToken     string `env:"TELEGRAM_BOT_TOKEN" env-required:"true"`
}

type Database struct {
	Host     string `env:"HOST" env-required:"true"`
	Port     string `env:"PORT" env-required:"true"`
	User     string `env:"USER" env-required:"true"`
	Password string `env:"PASSWORD" env-required:"true"`
	Name     string `env:"NAME" env-required:"true"`
}

type Server struct {
	Port string `env:"PORT" env-required:"true"`
}

func Setup() *Config {
	var cfg Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatalf("Error loading config: %w", err)
	}
	return &cfg
}
