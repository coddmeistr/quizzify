package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string        `yaFml:"env" env-default:"local"`
	PostgresUrl string        `env:"POSTGRES_URL"`
	TokenTTL    time.Duration `yaml:"token_ttl" env-default:"1h"`
	GRPC        GRPCConfig    `yaml:"grpc"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port" env-default:"8000"`
	Timeout time.Duration `yaml:"timeout" env-default:"5s"`
}

func MustLoad() *Config {
	cfgPath := fetchConfigPath()
	if cfgPath == "" {
		panic("config file path is not defined")
	}

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		panic("config file does not exist: " + cfgPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
		panic("Failed read config: " + err.Error())
	}

	return &cfg
}

func MustLoadByPath(path string) *Config {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("Failed read config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var cfgPath string

	flag.StringVar(&cfgPath, "config", "", "defined app's config path")
	flag.Parse()
	if cfgPath == "" {
		return os.Getenv("CONFIG_PATH")
	}

	return cfgPath
}
