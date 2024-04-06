package botkit

import (
	"NotifBot/internal/logger"
	"context"
	"github.com/vvv9912/sddb"
)

type UsersStorager interface {
	GetStatusUserByTgID(ctx context.Context, tgID int64) (int, int, error)
	AddUser(ctx context.Context, users sddb.Users) error
	UpdateStateByTgID(ctx context.Context, tgId int64, state int) error
	//GetCorzinaByTgID(ctx context.Context, tgID int64) ([]int64, error)
	//UpdateShopCartByTgId(ctx context.Context, tgId int64, corzina []int64) error
}
type LogStorager interface {
	WriteEvent(ctx context.Context, event logger.LoggerMsg) error
	ReadLastEvent(ctx context.Context) (*logger.LoggerMsg, error)
	ReadLastEventMicroservice(ctx context.Context, nameMicroservice string) (*logger.LoggerMsg, error)
}

type OrderStorager interface {
	sddb.OrderStorager
}
