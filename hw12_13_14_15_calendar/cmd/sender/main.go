package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	rmq "github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/brokers/rabbitmq"
	"github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/config"
	"github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/logger"
)

var configFile = flag.String("config", "./configs/sender.yaml", "path to configuration file")

func main() {
	flag.Parse()

	cfg, err := config.Sender(*configFile)
	if err != nil {
		fmt.Printf("invalid config: %s", err)
		os.Exit(1)
	}

	logg := logger.New(cfg.Logger.Level)

	c, err := rmq.NewConsumer(
		rmq.MQOptions{
			URL:         cfg.RabbitMQ.URL,
			RoutingKey:  cfg.RabbitMQ.RoutingKey,
			Ename:       cfg.RabbitMQ.Exchange.Name,
			Etype:       cfg.RabbitMQ.Exchange.Type,
			Edurable:    cfg.RabbitMQ.Exchange.Durable,
			EautoDelete: cfg.RabbitMQ.Exchange.AutoDelete,
			Einternal:   cfg.RabbitMQ.Exchange.Internal,
			EnoWait:     cfg.RabbitMQ.Exchange.NoWait,
			Qname:       cfg.RabbitMQ.Queue.Name,
			Qdurable:    cfg.RabbitMQ.Queue.Durable,
			QautoDelete: cfg.RabbitMQ.Queue.AutoDelete,
			Qexclusive:  cfg.RabbitMQ.Queue.Exclusive,
			QnoWait:     cfg.RabbitMQ.Queue.NoWait,
		},
		cfg.ConsumerTag,
	)
	if err != nil {
		logg.Error("failed to create consumer: " + err.Error())
		os.Exit(1)
	}
	logg.Info("consumer connected")

	s := newSender(c, os.Stdout)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		if err := s.stop(); err != nil {
			logg.Error("failed to stop sender: " + err.Error())
		}
	}()

	logg.Info("sender is starting")
	if err := s.start(ctx); err != nil {
		logg.Error("failed to start sender: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}

	logg.Info("sender is stopped")
}
