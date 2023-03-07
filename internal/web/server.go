package web

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/nickzhog/neo_nasa/internal/config"
	"github.com/nickzhog/neo_nasa/internal/service/neo"
	"github.com/nickzhog/neo_nasa/pkg/logging"
)

func PrepareServer(logger *logging.Logger, cfg *config.Config, storage neo.Repository) *http.Server {
	handler := NewHandler(logger, cfg, storage)

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Route("/neo", func(r chi.Router) {
		r.Get("/count", handler.getNeoCount)
		r.Post("/count", handler.upsertNeoCount)
	})

	return &http.Server{
		Addr:    cfg.Settings.Address,
		Handler: r,
	}
}

func Serve(ctx context.Context, logger *logging.Logger, srv *http.Server) (err error) {
	go func() {
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen:%+s\n", err)
		}
	}()

	logger.Tracef("server started")

	<-ctx.Done()

	logger.Tracef("server stopped")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err = srv.Shutdown(ctxShutDown); err != nil {
		logger.Fatalf("server Shutdown Failed:%+s", err)
	}

	logger.Tracef("server exited properly")

	if err == http.ErrServerClosed {
		err = nil
	}

	return
}
