package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type (
	loggerConf struct {
		Level string `yaml:"level"`
	}

	databaseConf struct {
		Use bool   `yaml:"use"`
		Dsn string `yaml:"dsn"`
	}

	httpServerConf struct {
		Address     string        `yaml:"address"`
		LogFile     string        `yaml:"log_file"`
		Timeout     time.Duration `yaml:"timeout"`
		IdleTimeout time.Duration `yaml:"idle_timeout"`
	}

	grpcServerConf struct {
		Port    string `yaml:"port"`
		LogFile string `yaml:"log_file"`
	}

	serverConf struct {
		HTTP httpServerConf `yaml:"http"`
		GRPC grpcServerConf `yaml:"grpc"`
	}

	exchangeConf struct {
		Name       string `yaml:"name"`
		Type       string `yaml:"type"`
		Durable    bool   `yaml:"durable"`
		AutoDelete bool   `yaml:"auto_delete"`
		Internal   bool   `yaml:"internal"`
		NoWait     bool   `yaml:"no_wait"`
	}

	queueConf struct {
		Name       string `yaml:"name"`
		Durable    bool   `yaml:"durable"`
		AutoDelete bool   `yaml:"auto_delete"`
		Exclusive  bool   `yaml:"exclusive"`
		NoWait     bool   `yaml:"no_wait"`
	}

	rabbitConf struct {
		URL        string       `yaml:"url"`
		RoutingKey string       `yaml:"routing_key"`
		Exchange   exchangeConf `yaml:"exchange"`
		Queue      queueConf    `yaml:"queue"`
	}

	CalendarConf struct {
		Logger   loggerConf   `yaml:"logger"`
		Database databaseConf `yaml:"database"`
		Server   serverConf   `yaml:"server"`
	}

	SchedulerConf struct {
		Logger   loggerConf   `yaml:"logger"`
		Database databaseConf `yaml:"database"`
		RabbitMQ rabbitConf   `yaml:"rabbitmq"`
		Interval string       `yaml:"interval"`
	}

	SenderConf struct {
		Logger      loggerConf `yaml:"logger"`
		RabbitMQ    rabbitConf `yaml:"rabbitmq"`
		ConsumerTag string     `yaml:"consumer_tag"`
	}
)

func Calendar(path string) (CalendarConf, error) {
	var cfg CalendarConf
	err := config(path, &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func Scheduler(path string) (SchedulerConf, error) {
	var cfg SchedulerConf
	err := config(path, &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func Sender(path string) (SenderConf, error) {
	var cfg SenderConf
	err := config(path, &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func config(path string, out any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, out)
	if err != nil {
		return err
	}

	return nil
}
