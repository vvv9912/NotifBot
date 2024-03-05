package main

import (
	"NotifBot/internal/botkit"
	"NotifBot/internal/config"
	"NotifBot/internal/logger"
	"NotifBot/internal/notifier"
	"NotifBot/internal/service"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
	"github.com/vvv9912/sddb"
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
	//	"time"
)

func run() error {
	flagLogLevel := "info"
	if err := logger.Initialize2(flagLogLevel, config.Get().PathFileLog, "NotifBot"); err != nil {
		return err
	}
	return nil
}

func main() {

	token := config.Get().TelegramBotNotifToken

	if err := run(); err != nil {
		logger.Log.Error("run error", zap.Error(err))
	}

	botApi, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		logger.Log.AddOriginalError(err).CustomError("failed to create bot", nil)
		return
	}

	db, err := sqlx.Connect("postgres", config.Get().DatabaseForLogDSN)
	if err != nil {
		logger.Log.Fatal("dbLog error", zap.Error(err))
	}

	bot := botkit.New(botApi)

	storage := sddb.NewStorageOrder(db)
	s := service.Service{Bot: bot}
	notifier.NewNotifierOrder(storage, s, 5*time.Second).StartNotifierOrder(context.Background())

	if err := bot.Run(context.Background()); err != nil {
		logger.Log.Error("run error", zap.Error(err))
		return
	}

}
