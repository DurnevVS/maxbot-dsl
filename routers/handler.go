package routers

import (
	"context"

	fsm "github.com/DurnevVS/maxbot-dsl/fsm/storage"
	maxbot "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

type FilterFunc[T schemes.UpdateInterface] func(update T, ctx context.Context, fsm *fsm.FSMContext) (bool, error)
type HandlerFunc[T schemes.UpdateInterface] func(api *maxbot.Api, update T, ctx context.Context, fsm *fsm.FSMContext) error

type Handler[T schemes.UpdateInterface] struct {
	Filters []FilterFunc[T]
	Handle  HandlerFunc[T]
}

func (h *Handler[T]) Run(api *maxbot.Api, update T, ctx context.Context, fsm *fsm.FSMContext) (bool, error) {
	for _, f := range h.Filters {
		ok, err := f(update, ctx, fsm)
		if err != nil || !ok {
			return false, err
		}
	}
	return true, h.Handle(api, update, ctx, fsm)
}
