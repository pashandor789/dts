package main

import (
	"context"
	pkghttp "dts/pkg/http"
	"dts/tabasco/http/api"
	storage "dts/tabasco/storage"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"dts/pkg/config"

	"dts/pkg/log"

	_ "dts/tabasco/http/docs"
	"golang.org/x/sync/errgroup"
)

type HttpTabascoConfig struct {
	Port    string         `yaml:"port"`
	Storage storage.Config `yaml:"storage"`
	Log     log.Config     `yaml:"log"`
}

// @title           Task Batching Storage Coordinator API
// @version         1.0
// @description     HTTP Tabasco
// @termsOfService  http://swagger.io/terms/

// @host      localhost:8000
// @BasePath  /
func main() {
	cfg := appconfig.MustParseAppConfig[*HttpTabascoConfig]()

	logger := log.New(os.Stdout, cfg.Log)
	defer func() { _ = logger.Sync() }()
	logger.Infof("Service started with config: %+v", cfg)

	g, ctx := errgroup.WithContext(context.Background())

	defer LogShutdownDuration(ctx, logger)()

	g.Go(func() error {
		return ListenSignal(ctx, logger)
	})

	s, err := storage.New(cfg.Storage, logger)
	if err != nil {
		logger.Fatalf("Storage error : %w", err)
	}

	h := api.NewPublicHandler(
		logger,
		s,
	)

	publicHandler := pkghttp.NewHandler("/",
		pkghttp.WithLogging(logger),
		pkghttp.WithSwagger(),
		h.WithAPIHandlers(),
	)

	g.Go(func() error {
		return log.IfError(
			logger,
			pkghttp.RunServer(ctx, cfg.Port, logger, publicHandler),
			"Public server error",
		)
	})

	err = g.Wait()
	if err != nil && !errors.Is(err, ErrSignalReceived) {
		logger.WithError(err).Errorf("Exit reason")
	}
}

func LogShutdownDuration(ctx context.Context, logger log.Logger) func() {
	var shutdownTime time.Time
	go func() {
		<-ctx.Done()
		shutdownTime = time.Now()
	}()
	return func() {
		log.Elapsed(logger, shutdownTime, "Shutdown duration")
	}
}

var ErrSignalReceived = errors.New("operating system signal")

func ListenSignal(ctx context.Context, logger log.Logger) error {
	sigquit := make(chan os.Signal, 1)
	signal.Ignore(syscall.SIGHUP, syscall.SIGPIPE)
	signal.Notify(sigquit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		return nil
	case sig := <-sigquit:
		logger.WithField("signal", sig).Infof("Captured signal")
		logger.Infof("Gracefully shutting down server...")
		return errors.New("operating system signal")
	}
}
