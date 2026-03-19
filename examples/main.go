package examples

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/DurnevVS/maxbot-dsl/filters"
	"github.com/DurnevVS/maxbot-dsl/handlers"

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

	dispatcher := handlers.
		NewDispatcher().
		AddRouter(StartBot()).
		AddRouter(StartCommand())

	for update := range api.GetUpdates(ctx) {
		go dispatcher.Dispatch(api, update, ctx)
	}
}

func StartBot() *handlers.Router {
	router := handlers.NewRouter()
	router.OnBotStarted(func(rb *handlers.RouteBuilder[*schemes.BotStartedUpdate]) {
		rb.Handle(func(api *maxbot.Api, update *schemes.BotStartedUpdate, ctx context.Context) error {
			message := maxbot.NewMessage().
				SetChat(update.GetChatID()).
				SetText("Hello, world!")

			api.Messages.Send(ctx, message)
			return nil
		})
	})

	return router
}

func StartCommand() *handlers.Router {
	router := handlers.NewRouter()
	router.OnMessage(func(rb *handlers.RouteBuilder[*schemes.MessageCreatedUpdate]) {
		rb.Filter(filters.IsCommand("/start")).
			Handle(func(api *maxbot.Api, update *schemes.MessageCreatedUpdate, ctx context.Context) error {
				message := maxbot.NewMessage().
					SetChat(update.GetChatID()).
					SetText("Processing /start")

				api.Messages.Send(ctx, message)
				return nil
			})
	})

	return router
}
