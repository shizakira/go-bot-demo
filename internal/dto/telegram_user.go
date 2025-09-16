package dto

type CreateTelegramUserInput struct {
	ChatID     int64  `json:"chat_id"`
	TelegramID int64  `json:"telegram_id"`
	Username   string `json:"username"`
}

type CreateTelegramUserOutput struct {
	UserID int64
}
