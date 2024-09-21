package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env                 string   `yaml:"env" env-default:"development"`
	HealthCheckInterval string   `yaml:"healthCheckInterval" env-default:"5s"`
	Servers             []string `yaml:"servers"`
	ListenPort          string   `yaml:"listenPort" env-default:":8080"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()

	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable is not set")
	}

	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("Error while opening config file: %s", err)
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("error while reading config file: %s", err)
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
