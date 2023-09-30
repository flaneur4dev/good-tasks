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
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

type ServerConf struct {
	HTTP HTTPServerConf `yaml:"http"`
}

type LogsConf struct {
	FilePath string `yaml:"file_path"`
}

type Config struct {
	Logger   LoggerConf   `yaml:"logger"`
	Database DatabaseConf `yaml:"database"`
	Server   ServerConf   `yaml:"server"`
	Logs     LogsConf     `yaml:"logs"`
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
