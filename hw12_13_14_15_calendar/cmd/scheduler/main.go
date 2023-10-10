package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	rmq "github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/brokers/rabbitmq"
	"github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/config"
	"github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/logger"
	sdb "github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/storage/scheduler"
	"github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/utils"
)

var configFile = flag.String("config", "./configs/scheduler.yaml", "path to configuration file")

func main() {
	flag.Parse()

	cfg, err := config.Scheduler(*configFile)
	if err != nil {
		fmt.Printf("invalid config: %s", err)
		os.Exit(1)
	}

	logg := logger.New(cfg.Logger.Level)

	p, err := rmq.NewProducer(
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
	)
	if err != nil {
		logg.Error("failed to create producer: " + err.Error())
		os.Exit(1)
	}
	logg.Info("producer connected")

	s, err := sdb.New(cfg.Database.Dsn)
	if err != nil {
		logg.Error("failed to connection to storage: " + err.Error())
		os.Exit(1)
	}
	logg.Info("database connected")

	sch, err := newScheduler(cfg.Interval, s, p)
	if err != nil {
		logg.Error("failed to create scheduler: " + err.Error())
		s.Close()
		os.Exit(1)
	}
	logg.Info("scheduler is starting")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	execute(ctx, cancel, logg, sch, time.Now())
	for {
		select {
		case <-ctx.Done():
			err := sch.stop()
			if err != nil {
				logg.Error("failed to stop scheduler: " + err.Error())
				os.Exit(1) //nolint:gocritic
			}

			logg.Info("scheduler is stopped")
			return

		case t := <-sch.ticker.C:
			execute(ctx, cancel, logg, sch, t)
		}
	}
}

func execute(ctx context.Context, cf context.CancelFunc, l *logger.Logger, s *scheduler, t time.Time) {
	l.Info("running scheduler iteration at: " + t.String())

	if err := s.execute(ctx, utils.ZeroTime(t)); err != nil {
		l.Error("failed to execute scheduler: " + err.Error())

		err := s.stop()
		if err != nil {
			l.Error("failed to stop scheduler: " + err.Error())
		}

		cf()
		os.Exit(1)
	}
}
