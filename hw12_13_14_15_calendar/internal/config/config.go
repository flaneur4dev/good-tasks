package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type LoggerConf struct {
	Level string `yaml:"level"`
}

type DatabaseConf struct {
	Use bool   `yaml:"use"`
	Dsn string `yaml:"dsn"`
}

type HTTPServerConf struct {
	Address     string        `yaml:"address"`
	LogFile     string        `yaml:"log_file"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

type GRPCServerConf struct {
	Port    string `yaml:"port"`
	LogFile string `yaml:"log_file"`
}

type ServerConf struct {
	HTTP HTTPServerConf `yaml:"http"`
	GRPC GRPCServerConf `yaml:"grpc"`
}

type Config struct {
	Logger   LoggerConf   `yaml:"logger"`
	Database DatabaseConf `yaml:"database"`
	Server   ServerConf   `yaml:"server"`
}

func New(path string) (Config, error) {
	var cfg Config
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
