package service

import (
	"NotifBot/internal/bot/sender"
	"NotifBot/internal/botkit"
	"context"
	"fmt"
	"github.com/vvv9912/sddb"
)

type Service struct {
	Bot *botkit.Bot
}

func (s *Service) SendOrder(ctx context.Context, order []sddb.Orders) error {
	txt := ""
	for _, v := range order {
		txt += fmt.Sprintf("New order: \nID: %d\n TgId: %d,\n UserName: @%s,\n firstName: %s,\nLastName: %s,\n Status: %s,\n Order: %s,\n CreatedAt: %s,\n UpdateAt: %s\n",
			v.ID, v.TgID, v.UserName, v.FirstName, v.LastName, v.StatusOrder, v.Order, v.CreatedAt, v.UpdateAt)
	}
	err := sender.Send(ctx, s.Bot.Api, txt)
	if err != nil {
		return err
	}
	return nil

}
