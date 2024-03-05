package mw

import (
	"NotifBot/internal/botkit"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Middleware struct {
	UserStorage botkit.UsersStorager
}

func NewMiddleware(userStorage botkit.UsersStorager) *Middleware {
	return &Middleware{UserStorage: userStorage}
}

// сюда кэш можно
func (m *Middleware) MwAdminOnly(next botkit.ViewFunc) botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update, botInfo botkit.BotInfo) error {
		//проверка на админа
		next(ctx, bot, update, botInfo)

		return nil
	}
}
