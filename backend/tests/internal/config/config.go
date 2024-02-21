package config

import (
	"errors"
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env        string     `yaml:"env" env-required:"true"`
	HTTPServer HTTPServer `yaml:"http-server"`
	MongoDB    MongoDB    `yaml:"mongodb"`
	Service    Service    `yaml:"service" env-required:"true"`
}

type HTTPServer struct {
	Port           string        `yaml:"port" env-default:"8080"`
	Host           string        `yaml:"host" env-default:"localhost"`
	Timeout        time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout    time.Duration `yaml:"idle-timeout" env-default:"60s"`
	MaxConnections int           `yaml:"max-connections" env-default:"100"`
	MaxHeaderSize  int64         `yaml:"max-header-size" env-default:"1048576"`
	MaxBodySize    int64         `yaml:"max-body-size" env-default:"1048576"`
}

type MongoDB struct {
	ConnectionURI string `env:"MONGO_CONNECTION_URI" env-required:"true"`
	DatabaseName  string `yaml:"database-name" env-default:"quizzify-tests"`
}

func MustLoad() *Config {
	var cfgPath string
	flag.StringVar(&cfgPath, "config", "", "path to config file")
	flag.Parse()

	if cfgPath == "" {
		cfgPath = os.Getenv("CONFIG_PATH")
	}

	cfg, err := loadFromPath(cfgPath)
	if err != nil {
		panic(err)
	}

	return cfg
}

func loadFromPath(path string) (*Config, error) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return nil, errors.New("config file not found in " + path)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
