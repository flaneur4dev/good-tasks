package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/app"
	"github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/config"
	cs "github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/contracts"
	"github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/logger"
	httpserver "github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/server/http"
	"github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/storage/pgdb"
)

type AppStore interface {
	Events(ctx context.Context, start, end time.Time) ([]cs.Event, error)
	CreateEvent(ctx context.Context, ne cs.Event) error
	UpdateEvent(ctx context.Context, id string, ne cs.Event) error
	DeleteEvent(ctx context.Context, id string) error
	Check(ctx context.Context) error
	Close() error
}

var configFile = flag.String("config", "./configs/config.yaml", "Path to configuration file")

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

	cfg, err := config.New(*configFile)
	if err != nil {
		fmt.Printf("invalid config: %s", err)
		os.Exit(1)
	}

	logg := logger.New(cfg.Logger.Level)

	var s AppStore
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

	server, err := httpserver.New(
		logg, app.New(logg, s),
		cfg.Logs.FilePath,
		cfg.Server.HTTP.Address,
		cfg.Server.HTTP.Timeout,
		cfg.Server.HTTP.IdleTimeout,
	)
	if err != nil {
		logg.Error("failed to create server: " + err.Error())
		s.Close()
		os.Exit(1) //nolint:gocritic
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.HTTP.StopTimeout)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			logg.Info("calendar is stopped")
			return
		}
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		s.Close()
		os.Exit(1) //nolint:gocritic
	}
}
