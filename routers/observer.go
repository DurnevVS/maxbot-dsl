package routers

import (
	"context"

	maxbot "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

type EventObserver[T schemes.UpdateInterface] struct {
	handlers []*Handler[T]
}

func (o *EventObserver[T]) Register(h *Handler[T]) {
	o.handlers = append(o.handlers, h)
}

func (o *EventObserver[T]) Trigger(api *maxbot.Api, update T, ctx context.Context) error {
	for _, h := range o.handlers {
		if err := h.Run(api, update, ctx); err != nil {
			return err
		}
	}
	return nil
}
