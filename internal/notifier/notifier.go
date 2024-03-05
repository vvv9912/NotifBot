package notifier

import (
	"NotifBot/internal/botkit"
	"NotifBot/internal/logger"
	"NotifBot/internal/service"
	"context"
	"github.com/vvv9912/sddb"
	"go.uber.org/zap"
	"time"
)

type NotifierOrder struct {
	storageOrder botkit.OrderStorager
	TimerSend    time.Duration
	service.Service
}

func NewNotifierOrder(storageOrder botkit.OrderStorager, service service.Service, timerSend time.Duration) *NotifierOrder {
	return &NotifierOrder{storageOrder: storageOrder, TimerSend: timerSend, Service: service}
}

// Проверка состояния заказа
func (n *NotifierOrder) NotifierPending(ctx context.Context) error {

	orders, err := n.storageOrder.GetOrderByStatus(ctx, sddb.StatusOrderNew)
	if err != nil {
		logger.Log.Error("NotifierPending error", zap.Error(err))
		return err
	}
	if len(orders) == 0 {
		return nil
	}

	err = n.Service.SendOrder(ctx, orders)
	if err != nil {
		logger.Log.Error("NotifierPending error", zap.Error(err))
		return err
	}
	for _, v := range orders {
		err = n.storageOrder.UpdateOrderByStatus(ctx, sddb.StatusOrderRead, v.ID)

		if err != nil {
			logger.Log.Error("UpdateOrderByStatus error", zap.Error(err))
			return err
		}

	}
	return nil
}

func (n *NotifierOrder) StartNotifierOrder(ctx context.Context) {
	if n.TimerSend == 0 {
		return
	}
	go func() {
		ticker := time.NewTicker(n.TimerSend)
		for {
			select {
			case <-ticker.C:
				n.NotifierPending(ctx)

			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()

	return

}
