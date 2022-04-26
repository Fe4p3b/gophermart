package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/Fe4p3b/gophermart/internal/api/handlers"
	"github.com/Fe4p3b/gophermart/internal/service"
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

	logger.Sugar().Info("asdfdf")

	as := service.NewAuth(logger.Sugar())

	ah := handlers.NewAuth(logger.Sugar(), as)
	oh := handlers.NewOrders(logger.Sugar())
	bh := handlers.NewBalance(logger.Sugar())

	h := handlers.New(logger.Sugar())
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
