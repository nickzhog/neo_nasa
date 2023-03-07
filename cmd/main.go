package main

import (
	"context"
	"os"
	"os/signal"
	"sync"

	"github.com/nickzhog/neo_nasa/internal/config"
	"github.com/nickzhog/neo_nasa/internal/migration"
	"github.com/nickzhog/neo_nasa/internal/service/neo/db"
	"github.com/nickzhog/neo_nasa/internal/service/processer"
	"github.com/nickzhog/neo_nasa/internal/web"
	"github.com/nickzhog/neo_nasa/pkg/logging"
	"github.com/nickzhog/neo_nasa/pkg/postgres"
)

func main() {
	cfg := config.GetConfig()
	logger := logging.GetLogger()
	logger.Tracef("config: %+v", cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		oscall := <-c
		logger.Tracef("system call:%+v", oscall)
		cancel()
	}()

	postgresClient, err := postgres.NewClient(ctx, 2, cfg.Settings.PostgresStorage.DatabaseDSN)
	if err != nil {
		logger.Fatalf("db error: %s", err.Error())
	}

	err = migration.Migrate(cfg.Settings.PostgresStorage.DatabaseDSN)
	if err != nil {
		logger.Fatalf("migrate error: %s", err.Error())
	}

	storage := db.NewRepository(postgresClient, logger)

	wg := new(sync.WaitGroup)
	wg.Add(2)
	go func() {
		srv := web.PrepareServer(logger, cfg, storage)
		if err := web.Serve(ctx, logger, srv); err != nil {
			logger.Errorf("failed to serve: %s", err.Error())
		}
		wg.Done()
	}()

	go func() {
		err := processer.NewProcesser(logger, cfg, storage).StartScan(ctx)
		if err != nil {
			logger.Errorf("neo processer error: %s", err.Error())
		}
		wg.Done()
	}()

	wg.Wait()
}
