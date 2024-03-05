package sender

import (
	"NotifBot/internal/logger"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

// id группы, куда отправляем уведомления
const idTgGroup = -4190783806

func Send(ctx context.Context, bot *tgbotapi.BotAPI, txt string) error {
	if _, err := bot.Send(tgbotapi.NewMessage(idTgGroup, txt)); err != nil {
		logger.Log.Error("Send error", zap.Error(err))
		return err
	}
	return nil
}
