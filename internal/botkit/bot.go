package botkit

import (
	"NotifBot/internal/logger"
	"context"
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"time"
)

//логика работы с ботом
import (
	"runtime/debug"
)

type BotCommand struct {
	Cmd  string `json:"cmd,omitempty"`
	Data string `json:"data,omitempty"` //в дата упоквано в зависимости от сообщения модель
}
type Bot struct {
	Api           *tgbotapi.BotAPI
	cmdViews      map[string]ViewFunc // комманды тг бота
	textViews     map[string]ViewFunc
	callbackViews map[string]ViewFunc
}
type BotInfo struct {
	TgId     int64
	IdStatus int
	IdState  int
}

type ViewFunc func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update, botInfo BotInfo) error

func New(api *tgbotapi.BotAPI) *Bot {
	return &Bot{Api: api}
}

func (b *Bot) RegisterTextView(cmd string, view ViewFunc) {
	if b.textViews == nil {
		b.textViews = make(map[string]ViewFunc)
	}
	b.textViews[cmd] = view
}
func (b *Bot) RegisterCallbackView(cmd string, view ViewFunc) {
	if b.callbackViews == nil {
		b.callbackViews = make(map[string]ViewFunc)
	}
	b.callbackViews[cmd] = view
}
func (b *Bot) RegisterCmdView(cmd string, view ViewFunc) {
	if b.cmdViews == nil {
		b.cmdViews = make(map[string]ViewFunc)
	}
	b.cmdViews[cmd] = view
}

func (b *Bot) handleUpdate(ctx context.Context, update tgbotapi.Update) {
	defer func() { //перехват паники
		if p := recover(); p != nil {
			logger.Log.CustomError("panic recovered", map[string]interface{}{
				"panic": p,
				"stack": string(debug.Stack()),
			})
		}
	}()
	//проверка авторизации пользователя
	//проверка из бд или кэша пользователя
	//
	var view ViewFunc
	if update.Message == nil {
		if update.CallbackQuery != nil {
			//какая то сложная логика, нужно рефакторингом заняться todo
			var Data BotCommand

			err := json.Unmarshal([]byte(update.CallbackQuery.Data), &Data)
			if err != nil {
				logger.Log.Error("Json преобразование callback", zap.Error(err))
				return
			}

			callbackView, ok := b.callbackViews[Data.Cmd]

			if !ok {
				return
			}
			view = callbackView
		} else if update.InlineQuery != nil {
			//тут  чтото
		} else {

			return
		}
	} else {
		if update.Message.IsCommand() {
			cmd := update.Message.Command()

			cmdView, ok := b.cmdViews[cmd]
			if !ok {
				return
			}

			view = cmdView

		} else if update.Message.Document != nil {

			cmd := update.Message.Caption
			cmdView, ok := b.cmdViews[cmd]
			if !ok {
				return
			}

			view = cmdView
		} else {
			//Если текстовая команда
			text := update.Message.Text
			if text == "" {
				return
			}
			textView, ok := b.textViews[text]
			if !ok {
				return
			}
			view = textView

		}
	}
	var botInfo BotInfo
	botInfo.TgId = update.FromChat().ID
	if err := view(ctx, b.Api, update, botInfo); err != nil {
		logger.Log.Error("failed to handle update, view error", zap.Error(err))
		//b.Api.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintln(err)))
		if _, err := b.Api.Send(
			tgbotapi.NewMessage(update.Message.Chat.ID, "internal error"),
		); err != nil {
			logger.Log.Error("failed to handle update, send error to tg", zap.Error(err))
		}
	}
}

func (b *Bot) Run(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 120

	updates := b.Api.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			go func(update tgbotapi.Update) {
				updateCtx, updateCancel := context.WithTimeout(ctx, 60*time.Second)
				defer updateCancel()
				b.handleUpdate(updateCtx, update)

			}(update)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
