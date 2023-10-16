package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/app"
	"github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/config"
	cs "github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/contracts"
	"github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/logger"
	grpcserver "github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/server/grpc"
	httpserver "github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/server/http"
	"github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/storage/pgdb"
)

type appStore interface {
	Events(ctx context.Context, start, end time.Time) ([]cs.Event, error)
	CreateEvent(ctx context.Context, ne cs.Event) error
	UpdateEvent(ctx context.Context, id string, ne cs.Event) error
	DeleteEvent(ctx context.Context, id string) error
	Check(ctx context.Context) error
	Close() error
}

var configFile = flag.String("config", "./configs/calendar.yaml", "path to configuration file")

func main() {
	flag.Parse()

	switch flag.Arg(0) {
	case "help":
		printHelp()
		return
	case "version":
		printVersion()
		return
	}

	cfg, err := config.Calendar(*configFile)
	if err != nil {
		fmt.Printf("invalid config: %s", err)
		os.Exit(1)
	}

	logg := logger.New(cfg.Logger.Level)

	var s appStore
	if cfg.Database.Use {
		s, err = pgdb.New(cfg.Database.Dsn)
		if err != nil {
			logg.Error("failed to connection to storage: " + err.Error())
			os.Exit(1)
		}
		logg.Info("database connected...")
	} else {
		s = memory.New()
	}
	defer s.Close()

	calendar := app.New(logg, s)

	hsrv, err := httpserver.New(
		logg, calendar,
		cfg.Server.HTTP.LogFile,
		cfg.Server.HTTP.Address,
		cfg.Server.HTTP.Timeout,
		cfg.Server.HTTP.IdleTimeout,
	)
	if err != nil {
		logg.Error("failed to create http server: " + err.Error())
		s.Close()
		os.Exit(1) //nolint:gocritic
	}

	gsrv, err := grpcserver.New(
		calendar,
		cfg.Server.GRPC.LogFile,
		cfg.Server.GRPC.Port,
	)
	if err != nil {
		logg.Error("failed to create grpc server: " + err.Error())
		s.Close()
		os.Exit(1) //nolint:gocritic
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := hsrv.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}

		if err := gsrv.Stop(); err != nil {
			logg.Error("failed to stop grpc server: " + err.Error())
		}
	}()

	go func() {
		defer wg.Done()

		logg.Info("grpc server is running...")

		if err := gsrv.Start(); err != nil {
			logg.Error("failed to start grpc server: " + err.Error())
			cancel()
			return
		}

		logg.Info("grpc server is stopped")
	}()

	logg.Info("http server is running...")

	if err := hsrv.Start(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			logg.Info("http server is stopped")
		} else {
			logg.Error("failed to start http server: " + err.Error())
			cancel()
			s.Close()
			os.Exit(1) //nolint:gocritic
		}
	}

	wg.Wait()
}
