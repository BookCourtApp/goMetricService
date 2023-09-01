package config

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env string `yaml:"env" env-default:"local" env-required:"true"`
	//StoragePath string `yaml:"storage_path" env-required:"true"`
	Clickhouse `yaml:"clickhouse" env-required:"true"`
	Redis      `yaml:"redis" env-required:"true"`
	HTTPServer `yaml:"http_server"`
}

type Clickhouse struct {
	Host     string `yaml:"host" env-default:"localhost:9000"`
	Name     string `yaml:"name" env:"BOOKCOURT_DB" env-default:""`
	User     string `yaml:"user" env:"BOOKCOURT_USER" env-required:"true"`
	Password string `yaml:"password" env:"BOOKCOURT_PASSWORD" env-required:"true"`
}

type Redis struct {
	Host     string `yaml:"host" env-default:"localhost6379"`
	Password string `yaml:"password" env-default:""`
	Name     int    `yaml:"name" env-default:"0"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:5078"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	//configPath := os.Getenv("CONFIG_PATH")

	configPath := flag.String(
		"config",
		"",
		"Path to a config file",
	)

	flag.Parse()

	if *configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(*configPath); os.IsNotExist(err) {
		log.Fatalf("config file '%s' does not exist", *configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(*configPath, &cfg); err != nil {
		log.Fatalf("can't read config: %s", err)
	}

	return &cfg
}
