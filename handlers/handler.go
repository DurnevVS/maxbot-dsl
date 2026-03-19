package handlers

import (
	"context"

	maxbot "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

type FilterFunc[T schemes.UpdateInterface] func(update T, ctx context.Context) (bool, error)

type HandlerFunc[T schemes.UpdateInterface] func(api *maxbot.Api, update T, ctx context.Context) error

type Handler[T schemes.UpdateInterface] struct {
	Filters []FilterFunc[T]
	Handle  HandlerFunc[T]
}

func (h *Handler[T]) Run(api *maxbot.Api, update T, ctx context.Context) error {
	for _, f := range h.Filters {
		ok, err := f(update, ctx)
		if err != nil || !ok {
			return nil
		}
	}
	return h.Handle(api, update, ctx)
}
