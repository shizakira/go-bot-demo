package telegram

import (
	"context"
	"errors"
	"github.com/shizakira/daily-tg-bot/internal/dto"
	"regexp"
	"strconv"
)

func (tb *Bot) getValueFromQueryByRe(re string, query string) (string, error) {
	r, _ := regexp.Compile(re)
	matches := r.FindStringSubmatch(query)
	if len(matches) > 1 {
		return matches[1], nil
	}
	return "", errors.New("not found by re")
}

func (tb *Bot) handleTaskClosure(ctx context.Context, query string, isDone bool) error {
	rawID, err := tb.getValueFromQueryByRe("ID: ([0-9]+)", query)
	if err != nil {
		return err
	}

	taskID, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		return err
	}

	if err = tb.taskService.CloseTask(ctx, dto.CloseTaskInput{TaskID: taskID, IsDone: isDone}); err != nil {
		return err
	}

	return nil
}
