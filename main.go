package main

import (
	"WalletTopUp/internal/config"
	"WalletTopUp/internal/infra/cache"
	"WalletTopUp/internal/infra/database"
	"WalletTopUp/internal/infra/database/repository"
	httpiface "WalletTopUp/internal/interface/http"
	"WalletTopUp/internal/interface/http/handler"
	"WalletTopUp/internal/service"
	"WalletTopUp/pkg/lib"
	"context"
	"errors"
	"flag"
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
	doMigrate := flag.String("db", "", "Execute auto migration on run.")
	flag.Parse()
	// Config load and Init logger.
	conf := config.Load()

	// Init infra engines.
	logger := lib.NewLogger(conf, log.Default())
	db := database.Connect(conf.Database)
	if doMigrate != nil && *doMigrate == "migrate" {
		if err := database.Migrate(db); err != nil {
			log.Panicf("failed to migrate: %v", err)
		}
	}

	client := cache.Connect(context.Background(), conf.CacheClient)

	// Dependencies injection.
	txRepo := repository.NewTxRepo(logger, db)
	userRepo := repository.NewUserRepo(logger, db)
	txManager := repository.NewTxManager(db)
	walletRepo := repository.NewWalletRepo(logger, db)

	txCache := cache.NewNoopTxCache()
	if conf.CacheClient.Enabled {
		txCache = cache.NewTxCache(client)
	}

	walletSvc := service.NewWalletSvc(logger, userRepo, txRepo, txManager, walletRepo, txCache)

	handler := handler.NewWalletHandler(logger, walletSvc)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Init gin server.
	gin.SetMode(conf.Server.Mode)
	router := httpiface.NewRouter(logger, handler)

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

	event := "shuting down server"
	logger.Info(ctx, lib.Meta{
		Event: event,
		Msg:   "shutdown signal received.",
	})

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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
