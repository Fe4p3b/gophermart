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

	as := authService.NewAuth(logger.Sugar())
	os := orderService.New(logger.Sugar())
	bs := balanceService.New(logger.Sugar())

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
