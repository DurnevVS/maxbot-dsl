package routers

import (
	"context"

	fsm "github.com/DurnevVS/maxbot-dsl/fsm/storage"
	maxbot "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

type EventObserver[T schemes.UpdateInterface] struct {
	handlers []*Handler[T]
}

func (o *EventObserver[T]) Register(h *Handler[T]) {
	o.handlers = append(o.handlers, h)
}

func (o *EventObserver[T]) Trigger(api *maxbot.Api, update T, ctx context.Context, fsm *fsm.FSMContext) (bool, error) {
	for _, h := range o.handlers {
		handled, err := h.Run(api, update, ctx, fsm)
		if err != nil {
			return false, err
		}
		if handled {
			return true, nil
		}
	}
	return false, nil
}
