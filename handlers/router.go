package handlers

import (
	"context"

	maxbot "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

type RouteBuilder[T schemes.UpdateInterface] struct {
	filters []FilterFunc[T]
	handle  HandlerFunc[T]
}

func (b *RouteBuilder[T]) Filter(f FilterFunc[T]) *RouteBuilder[T] {
	b.filters = append(b.filters, f)
	return b
}

func (b *RouteBuilder[T]) Handle(h HandlerFunc[T]) *RouteBuilder[T] {
	b.handle = h
	return b
}

type Router struct {
	messageObserver  *EventObserver[*schemes.MessageCreatedUpdate]
	callbackObserver *EventObserver[*schemes.MessageCallbackUpdate]
	botStartObserver *EventObserver[*schemes.BotStartedUpdate]
	botEndObserver   *EventObserver[*schemes.BotStopedFromChatUpdate]
}

func NewRouter() *Router {
	return &Router{
		messageObserver:  &EventObserver[*schemes.MessageCreatedUpdate]{},
		callbackObserver: &EventObserver[*schemes.MessageCallbackUpdate]{},
		botStartObserver: &EventObserver[*schemes.BotStartedUpdate]{},
		botEndObserver:   &EventObserver[*schemes.BotStopedFromChatUpdate]{},
	}
}

func (r *Router) Resolve(
	api *maxbot.Api,
	update schemes.UpdateInterface,
	ctx context.Context,
) error {

	switch upd := update.(type) {

	case *schemes.BotStartedUpdate:
		return r.botStartObserver.Trigger(api, upd, ctx)

	case *schemes.MessageCreatedUpdate:
		return r.messageObserver.Trigger(api, upd, ctx)

	case *schemes.MessageCallbackUpdate:
		return r.callbackObserver.Trigger(api, upd, ctx)

	case *schemes.BotStopedFromChatUpdate:
		return r.botEndObserver.Trigger(api, upd, ctx)
	}

	return nil
}

func (r *Router) OnBotStarted(build func(*RouteBuilder[*schemes.BotStartedUpdate])) {
	b := &RouteBuilder[*schemes.BotStartedUpdate]{}
	build(b)

	r.botStartObserver.Register(&Handler[*schemes.BotStartedUpdate]{
		Filters: b.filters,
		Handle:  b.handle,
	})
}

func (r *Router) OnMessage(build func(*RouteBuilder[*schemes.MessageCreatedUpdate])) {
	b := &RouteBuilder[*schemes.MessageCreatedUpdate]{}
	build(b)

	r.messageObserver.Register(&Handler[*schemes.MessageCreatedUpdate]{
		Filters: b.filters,
		Handle:  b.handle,
	})
}

func (r *Router) OnCallback(build func(*RouteBuilder[*schemes.MessageCallbackUpdate])) {
	b := &RouteBuilder[*schemes.MessageCallbackUpdate]{}
	build(b)

	r.callbackObserver.Register(&Handler[*schemes.MessageCallbackUpdate]{
		Filters: b.filters,
		Handle:  b.handle,
	})
}

func (r *Router) OnBotEnd(build func(*RouteBuilder[*schemes.BotStopedFromChatUpdate])) {
	b := &RouteBuilder[*schemes.BotStopedFromChatUpdate]{}
	build(b)

	r.botEndObserver.Register(&Handler[*schemes.BotStopedFromChatUpdate]{
		Filters: b.filters,
		Handle:  b.handle,
	})
}
