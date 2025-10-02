package main

import (
	"context"
	"github.com/go-co-op/gocron/v2"
	"github.com/shizakira/daily-tg-bot/internal/adapters/telegram"
	"github.com/sirupsen/logrus"
	"time"
)

func initScheduler(ctx context.Context, notifier *telegram.Notifier) {
	s, err := gocron.NewScheduler()
	if err != nil {
		logrus.Fatal(err)
	}

	j, err := s.NewJob(
		gocron.DurationJob(
			30*time.Minute,
		),
		gocron.NewTask(
			func(jobCtx context.Context) {
				logrus.Info("notifier run")
				if runErr := notifier.Run(jobCtx); runErr != nil {
					logrus.Error(runErr)
				}
			},
		),
	)
	if err != nil {
		logrus.Error(err)
	}

	logrus.Infof("scheduler initialized id %v", j.ID)

	s.Start()

	<-ctx.Done()

	err = s.Shutdown()
	if err != nil {
		logrus.Error(err)
	}
}
