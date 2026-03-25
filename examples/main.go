package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/DurnevVS/maxbot-dsl/filters"
	fsm "github.com/DurnevVS/maxbot-dsl/fsm/storage"
	"github.com/DurnevVS/maxbot-dsl/routers"

	maxbot "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer stop()

	token := "BOT_TOKEN"

	api, err := maxbot.New(token)
	if err != nil {
		panic(err)
	}

	dispatcher := routers.
		NewDispatcher(fsm.NewMemoryStorage()).
		AddRouters(StartBot(), SayMyName(), StartCommand())

	for update := range api.GetUpdates(ctx) {
		go dispatcher.Dispatch(api, update, ctx)
	}
}

func StartBot() *routers.Router {
	router := routers.NewRouter()
	router.OnBotStarted(func(rb *routers.RouteBuilder[*schemes.BotStartedUpdate]) {
		rb.Handle(func(api *maxbot.Api, update *schemes.BotStartedUpdate, ctx context.Context, fsm *fsm.FSMContext) error {
			message := maxbot.NewMessage().
				SetChat(update.GetChatID()).
				SetText("Hello, world!")

			api.Messages.Send(ctx, message)
			return nil
		})
	})

	return router
}

func StartCommand() *routers.Router {
	router := routers.NewRouter()
	router.OnMessage(func(rb *routers.RouteBuilder[*schemes.MessageCreatedUpdate]) {
		rb.Filter(filters.IsCommand("/start")).
			Handle(func(api *maxbot.Api, update *schemes.MessageCreatedUpdate, ctx context.Context, fsm *fsm.FSMContext) error {
				message := maxbot.NewMessage().
					SetChat(update.GetChatID()).
					SetText("Processing /start")

				api.Messages.Send(ctx, message)

				fsm.SetState(ctx, "SayMyName")

				return nil
			})
	})

	return router
}

func SayMyName() *routers.Router {
	router := routers.NewRouter()
	router.OnMessage(func(rb *routers.RouteBuilder[*schemes.MessageCreatedUpdate]) {
		rb.Filter(filters.IsCommand("/start")).Filter(filters.StateFilter[*schemes.MessageCreatedUpdate]("SayMyName")).
			Handle(func(api *maxbot.Api, update *schemes.MessageCreatedUpdate, ctx context.Context, fsm *fsm.FSMContext) error {
				message := maxbot.NewMessage().
					SetChat(update.GetChatID()).
					SetText("Heinzerberg!")

				api.Messages.Send(ctx, message)

				fsm.Clear(ctx)

				return nil
			})
	})

	return router
}
