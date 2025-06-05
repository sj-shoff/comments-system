package config

import (
	"log/slog"
	"os"

	"comments-system/pkg/logger/sl"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server     ServerConfig `yaml:"server"`
	Database   Postgres     `yaml:"postgres"`
	Storage    string       `yaml:"storage"`
	Env        string       `yaml:"env" env-default:"local"`
	Migrations string       `yaml:"migrations" env-default:"./migrations"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type Postgres struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `env:"POSTGRES_PASSWORD"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

func (p Postgres) DSN() string {
	return "postgres://" + p.Username + ":" + p.Password + "@" + p.Host + ":" + p.Port + "/" + p.DBName + "?sslmode=" + p.SSLMode
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		slog.Error("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		slog.Error("Config file does not exist: %s", configPath, sl.Err(err))
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		slog.Error("Cannot read config: %s", sl.Err(err))
	}

	cfg.Database.Password = os.Getenv("POSTGRES_PASSWORD")

	return &cfg
}
