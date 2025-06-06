package config

import (
	"flag"
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

func MustLoad() *Config {
	configPath := flag.String("config", "", "path to config file")
	flag.Parse()

	path := ""
	if *configPath != "" {
		path = *configPath
	} else if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		path = envPath
	} else {
		path = "config.yaml"
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		slog.Error("Config file does not exist", slog.String("path", path))
		panic("config file does not exist: " + path)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		slog.Error("Cannot read config", sl.Err(err))
		panic("cannot read config: " + err.Error())
	}

	if cfg.Storage == "postgres" {
		if password := os.Getenv("POSTGRES_PASSWORD"); password != "" {
			cfg.Database.Password = password
		} else {
			slog.Error("POSTGRES_PASSWORD environment variable not set")
			os.Exit(1)
		}
	}

	return &cfg
}
