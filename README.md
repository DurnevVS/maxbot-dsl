# maxbot-dsl

The lightweight Go DSL for creating event-driven MAX messenger bots with easy routing and handler management.

## Requirements

maxbot-dsl is a framework 🙂 built **on top of the official Go MAX client**: [max-bot-api-client-go](https://github.com/max-messenger/max-bot-api-client-go) licensed under Apache License 2.0.

## Contributing
Open to contributors!

## Installation:

```bash
go get github.com/maxbot-dsl/maxbot-dsl/v0.1.1
```

## Example:

```go
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
```