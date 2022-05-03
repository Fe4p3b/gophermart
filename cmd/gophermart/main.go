package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/Fe4p3b/gophermart/cmd/gophermart/config"
	"github.com/Fe4p3b/gophermart/internal/api/accrual"
	"github.com/Fe4p3b/gophermart/internal/api/handler"
	"github.com/Fe4p3b/gophermart/internal/api/middleware"
	authService "github.com/Fe4p3b/gophermart/internal/service/auth"
	balanceService "github.com/Fe4p3b/gophermart/internal/service/balance"
	orderService "github.com/Fe4p3b/gophermart/internal/service/order"
	"github.com/Fe4p3b/gophermart/internal/storage/pg"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()
	sugaredLogger := logger.Sugar()

	cfg, err := config.SetConfig()
	if err != nil {
		sugaredLogger.Fatal(err)
	}

	sugaredLogger.Infow("Initialized config", "config", cfg)

	r := chi.NewRouter()

	db, err := pg.New(cfg.DatabaseURI)
	if err != nil {
		sugaredLogger.Fatalw("error initializing db", "error", err)
	}

	if err := db.InitialMigration(); err != nil {
		sugaredLogger.Fatalw("error applying migration", "error", err)
	}

	accrual := accrual.New(sugaredLogger, cfg.AccrualURL)

	us := pg.NewUserStorage(db)
	os := pg.NewOrderStorage(db)
	bs := pg.NewBalanceStorage(db)

	as, err := authService.NewAuth(sugaredLogger, us, 14, []byte(cfg.Secret))
	if err != nil {
		sugaredLogger.Fatalw("error creating auth service", "error", err)
	}

	ah := handler.NewAuth(sugaredLogger, as)
	oh := handler.NewOrder(sugaredLogger, orderService.New(sugaredLogger, os, accrual))
	bh := handler.NewBalance(sugaredLogger, balanceService.New(sugaredLogger, bs, db))

	m := middleware.NewAuthMiddleware(as)

	h := handler.New(sugaredLogger)
	h.SetupRouting(r, m, ah, oh, bh)

	srv := http.Server{Addr: cfg.Address, Handler: r}

	errgroup, ctx := errgroup.WithContext(context.Background())
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	errgroup.Go(func() error {
		sugaredLogger.Infof("started server at address - %s", cfg.Address)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			return err
		}
		return nil
	})

	errgroup.Go(func() error {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		return srv.Shutdown(ctx)
	})

	if err := errgroup.Wait(); err != nil {
		sugaredLogger.Fatalw("error while running server", "error", err)
	}
}
