package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/Fe4p3b/gophermart/internal/api/handler"
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

	r := chi.NewRouter()

	db, err := pg.New("postgres://postgres:12345@localhost:5432/gophermart?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	as := authService.NewAuth(logger.Sugar(), db)
	os := orderService.New(logger.Sugar(), db)
	bs := balanceService.New(logger.Sugar(), db, db)

	ah := handler.NewAuth(logger.Sugar(), as)
	oh := handler.NewOrder(logger.Sugar(), os)
	bh := handler.NewBalance(logger.Sugar(), bs)

	h := handler.New(logger.Sugar())
	h.SetupRouting(r, ah, oh, bh)

	srv := http.Server{Addr: ":8080", Handler: r}

	errgroup, ctx := errgroup.WithContext(context.Background())
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	errgroup.Go(func() error {
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
		log.Fatal(err)
	}
}
