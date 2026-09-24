package main

import (
	"WalletTopUp/internal/config"
	"WalletTopUp/internal/infra/database"
	"WalletTopUp/internal/infra/database/repository"
	httpiface "WalletTopUp/internal/interface/http"
	"WalletTopUp/internal/interface/http/handler"
	"WalletTopUp/internal/service"
	"WalletTopUp/pkg/lib"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Config load and Init logger.
	conf := config.Load()

	// Init infra engines.
	logger := lib.NewLogger(conf, log.Default())
	db := database.Connect(conf.Database)

	// Dependencies injection.
	txRepo := repository.NewTxRepo(logger, db)
	txSvc := service.NewTxSvc(logger, txRepo)
	handler := handler.NewWalletHandler(logger, txSvc)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Init gin server.
	gin.SetMode(conf.Server.Mode)
	router := httpiface.NewRouter(logger, &handler)

	srv := &http.Server{
		Addr:              fmt.Sprintf("%s:%s", conf.Server.Host, conf.Server.Port),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Run gin server
	errCh := make(chan error, 1)
	go func() {
		logger.Info(ctx, lib.Meta{
			Event: "start http server",
			Msg:   fmt.Sprintf("Server start on %s", srv.Addr),
		})
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	// Graceful shutdown
	select {
	case err := <-errCh:
		logger.Error(ctx, lib.Meta{
			Event: "http server in-flight",
			Msg:   "http server shutdown unexpectedly",
			Error: err,
		})
		os.Exit(1)
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	event := "shuting down server"
	logger.Info(shutdownCtx, lib.Meta{
		Event: event,
		Msg:   "http server is shutdown",
	})
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Warn(shutdownCtx, lib.Meta{
			Event: event,
			Msg:   "Server shutdown ungracefully.",
			Error: err,
		})
	}
	logger.Info(shutdownCtx, lib.Meta{
		Event: event,
		Msg:   "http server was shutdown gracefully",
	})
}
