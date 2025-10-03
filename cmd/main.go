package main

import (
	"context"
	"errors"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"

	"github.com/shizakira/daily-tg-bot/config"
	"github.com/shizakira/daily-tg-bot/internal/app"
)

func main() {
	logrus.SetOutput(os.Stdout)

	cfg, err := config.Load()
	if err != nil {
		logrus.Fatalf("load config: %v", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	application, err := app.New(ctx, cfg)
	if err != nil {
		logrus.Fatalf("init app: %v", err)
	}

	defer func() {
		if closeErr := application.Close(); closeErr != nil {
			logrus.WithError(closeErr).Error("failed to close application")
		}
	}()

	logrus.Info("starting tg bot")
	if err := application.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		logrus.WithError(err).Error("application stopped with error")
	}
}
