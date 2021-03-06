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
	withdrawalService "github.com/Fe4p3b/gophermart/internal/service/withdrawal"
	"github.com/Fe4p3b/gophermart/internal/storage/pg"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
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

	r.Use(chiMiddleware.Compress(5), middleware.ContextMiddleware)

	db, err := pg.New(cfg.DatabaseURI, cfg.MigrationsFolder)
	if err != nil {
		sugaredLogger.Fatalw("error initializing db", "error", err)
	}

	accrual := accrual.New(sugaredLogger, cfg.AccrualURL)

	as, err := authService.NewAuth(sugaredLogger, pg.NewUserStorage(db), 14, []byte(cfg.Secret))
	if err != nil {
		sugaredLogger.Fatalw("error creating auth service", "error", err)
	}

	bs := balanceService.New(sugaredLogger, pg.NewBalanceStorage(db))

	ah := handler.NewAuth(sugaredLogger, as)
	oh := handler.NewOrder(sugaredLogger, orderService.New(sugaredLogger, pg.NewOrderStorage(db), accrual))
	bh := handler.NewBalance(sugaredLogger, bs)
	wh := handler.NewWithdrawal(sugaredLogger, withdrawalService.New(sugaredLogger, pg.NewWithdrawalStorage(db), bs))

	m := middleware.NewAuthMiddleware(as)

	h := handler.New(sugaredLogger)
	h.SetupRouting(r, m, ah, oh, bh, wh)

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
