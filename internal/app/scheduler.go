package app

import (
	"context"
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
	tg "github.com/shizakira/daily-tg-bot/internal/adapters/telegram"
	"github.com/sirupsen/logrus"
)

type NotifierScheduler struct {
	scheduler gocron.Scheduler
	notifier  *tg.Notifier
	interval  time.Duration
}

func NewNotifierScheduler(notifier *tg.Notifier, interval time.Duration) (*NotifierScheduler, error) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	return &NotifierScheduler{
		scheduler: scheduler,
		notifier:  notifier,
		interval:  interval,
	}, nil
}

func (s *NotifierScheduler) Start(ctx context.Context) error {
	if s == nil {
		return nil
	}

	_, err := s.scheduler.NewJob(
		gocron.DurationJob(s.interval),
		gocron.NewTask(func(jobCtx context.Context) {
			logrus.Info("notifier run")
			if err := s.notifier.Run(jobCtx); err != nil {
				logrus.WithError(err).Error("notifier job failed")
			}
		}),
	)
	if err != nil {
		return fmt.Errorf("schedule notifier job: %w", err)
	}

	s.scheduler.Start()
	<-ctx.Done()

	return s.scheduler.Shutdown()
}
