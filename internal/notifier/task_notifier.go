package notifier

import "github.com/shizakira/daily-tg-bot/internal/ports"

type TaskNotifier struct {
	taskRepo *ports.TaskRepository
}

func (tn *TaskNotifier) SendNotifyForExpiredTasks() {

}

func (tn *TaskNotifier) SendNotifyFortNearExpiredTasks() {

}
