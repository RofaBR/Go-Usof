package cli

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/RofaBR/Go-Usof/internal/app"
)

func Run(args []string) bool {
	application, err := app.New()
	if err != nil {
		println("Failed to initialize application:", err.Error())
		return false
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := application.Run(); err != nil {
			application.Logger().Error("server error", "error", err)
			cancel()
		}
	}()

	select {
	case <-quit:
		application.Logger().Info("received shutdown signal")
	case <-ctx.Done():
		application.Logger().Info("context cancelled")
	}

	if err := application.Shutdown(context.Background()); err != nil {
		application.Logger().Error("error during shutdown", "error", err)
		return false
	}

	return true
}
