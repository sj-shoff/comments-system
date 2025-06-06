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

func (p Postgres) DSN() string {
	return "postgres://" + p.Username + ":" + p.Password + "@" + p.Host + ":" + p.Port + "/" + p.DBName + "?sslmode=" + p.SSLMode
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		configPath = "config.yaml"
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		slog.Error("Config file does not exist", slog.String("path", configPath))
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		slog.Error("Cannot read config", sl.Err(err))
		panic("cannot read config: " + err.Error())
	}

	cfg.Database.Password = os.Getenv("POSTGRES_PASSWORD")

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		return os.Getenv("CONFIG_PATH")
	}

	return res
}
