package telegram

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/shizakira/daily-tg-bot/internal/domain"
	"github.com/shizakira/daily-tg-bot/internal/ports"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"time"
)

type NotifyData struct {
	Task         *domain.Task
	TelegramUser *domain.TelegramUser
	Message      string
}

type Notifier struct {
	bot          *bot.Bot
	telegramRepo ports.TelegramUserRepository
	taskRepo     ports.TaskRepository
}

func NewNotifier(bot *bot.Bot, telegramRepo ports.TelegramUserRepository, taskRepo ports.TaskRepository) *Notifier {
	return &Notifier{bot: bot, telegramRepo: telegramRepo, taskRepo: taskRepo}
}

func (n *Notifier) Run(ctx context.Context) error {
	var expired []*domain.Task
	var soonExpired []*domain.Task

	g, fetchCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		tasks, err := n.taskRepo.GetExpiredTasks(fetchCtx)
		if err != nil {
			return fmt.Errorf("get expired tasks: %w", err)
		}
		expired = tasks
		return nil
	})

	g.Go(func() error {
		tasks, err := n.taskRepo.GetSoonExpiredTasks(fetchCtx)
		if err != nil {
			return fmt.Errorf("get soon expired tasks: %w", err)
		}
		soonExpired = tasks
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}

	userIds := make([]int64, 0, len(expired)+len(soonExpired))
	for _, task := range append(expired, soonExpired...) {
		userIds = append(userIds, task.UserID)
	}

	if len(userIds) == 0 {
		return nil
	}

	telegramUsers, err := n.telegramRepo.FindByUserIDs(ctx, userIds)
	if err != nil {
		return err
	}

	telegramUsersMap := make(map[int64]*domain.TelegramUser)
	for _, telegramUser := range telegramUsers {
		telegramUsersMap[telegramUser.UserID] = telegramUser
	}

	g2, notifyCtx := errgroup.WithContext(ctx)
	for _, task := range expired {
		task := task
		g2.Go(func() error {
			return n.Notify(notifyCtx, task, telegramUsersMap[task.UserID], "Task expired")
		})
	}

	for _, task := range soonExpired {
		task := task
		g2.Go(func() error {
			return n.Notify(notifyCtx, task, telegramUsersMap[task.UserID], "Task will be expire soon")
		})
	}

	if err = g2.Wait(); err != nil {
		return err
	}

	return nil
}

func (n *Notifier) Notify(ctx context.Context, task *domain.Task, tgUser *domain.TelegramUser, msg string) error {
	loc, _ := time.LoadLocation("Asia/Yekaterinburg")
	taskMsg := fmt.Sprintf(
		"ID: %d\nTitle: %s\nDescription: %s\nDeadline %s\n\n",
		task.ID, task.Title, task.Description, task.Deadline.In(loc).Format("2006-01-02 15:04"),
	)
	text := fmt.Sprintf("%s, %s", msg, taskMsg)
	logrus.WithFields(logrus.Fields{"task_message": taskMsg}).Infof("notify user %d", tgUser.UserID)
	_, err := n.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: tgUser.ChatID,
		Text:   text,
	})

	return err
}
