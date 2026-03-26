# maxbot-dsl

The lightweight Go DSL for creating event-driven MAX messenger bots with easy routing and handler management.

## Requirements

maxbot-dsl is a framework 🙂 built **on top of the official Go MAX client**: [max-bot-api-client-go](https://github.com/max-messenger/max-bot-api-client-go) licensed under Apache License 2.0.

## Contributing
Open to contributors!

## Installation:

```bash
go get github.com/maxbot-dsl/maxbot-dsl/v0.2.2
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

	dispatcher := routers.
		NewDispatcher(fsm.NewMemoryStorage()).
		AddRouters(StartRouter(), BackToMainMenu())

	for update := range api.GetUpdates(ctx) {
		go dispatcher.Dispatch(api, update, ctx)
	}
}

func StartRouter() *routers.Router {
	router := routers.NewRouter()

	router.OnBotStarted(func(rb *routers.RouteBuilder[*schemes.BotStartedUpdate]) {
		rb.Handle(func(api *maxbot.Api, update *schemes.BotStartedUpdate, ctx context.Context, fsm *fsm.FSMContext) error {
			return SendMainMenu(api, update, ctx, fsm)
		})
	})

	router.OnMessage(func(rb *routers.RouteBuilder[*schemes.MessageCreatedUpdate]) {
		rb.Filter(filters.IsCommand("/start")).
			Handle(func(api *maxbot.Api, update *schemes.MessageCreatedUpdate, ctx context.Context, fsm *fsm.FSMContext) error {
				return SendMainMenu(api, update, ctx, fsm)
			})
	})

	router.OnCallback(func(rb *routers.RouteBuilder[*schemes.MessageCallbackUpdate]) {
		rb.Filter(filters.Callback("Idle")).
			Handle(func(api *maxbot.Api, update *schemes.MessageCallbackUpdate, ctx context.Context, fsm *fsm.FSMContext) error {
				return EditToMainMenu(api, update, ctx, fsm)
			})
	})

	return router
}

const menuMsg = "Processing main menu handler. Press The Button!"

func SendMainMenu(api *maxbot.Api, update schemes.UpdateInterface, ctx context.Context, fsm *fsm.FSMContext) error {

	kb := maxbot.Keyboard{}
	kb.AddRow().AddCallback("The Button", schemes.DEFAULT, "Button")

	message := maxbot.NewMessage().
		SetChat(update.GetChatID()).
		AddKeyboard(&kb).
		SetText(menuMsg)

	api.Messages.Send(ctx, message)
	return nil
}

func EditToMainMenu(api *maxbot.Api, update *schemes.MessageCallbackUpdate, ctx context.Context, fsm *fsm.FSMContext) error {

	kb := maxbot.Keyboard{}
	kb.AddRow().AddCallback("The Button", schemes.DEFAULT, "Button")

	message := maxbot.NewMessage().
		SetChat(update.GetChatID()).
		AddKeyboard(&kb).
		SetText(menuMsg)

	api.Messages.EditMessage(ctx, update.Message.Body.Mid, message)
	return nil
}

func BackToMainMenu() *routers.Router {
	router := routers.NewRouter()

	router.OnCallback(func(rb *routers.RouteBuilder[*schemes.MessageCallbackUpdate]) {
		rb.Filter(filters.Callback("Button")).
			Handle(func(api *maxbot.Api, update *schemes.MessageCallbackUpdate, ctx context.Context, fsm *fsm.FSMContext) error {
				kb := maxbot.Keyboard{}
				kb.AddRow().AddCallback("Back", schemes.DEFAULT, "Idle")
				message := maxbot.NewMessage().
					SetChat(update.GetChatID()).
					AddKeyboard(&kb).
					SetText("Back to main menu")

				api.Messages.EditMessage(ctx, update.Message.Body.Mid, message)

				return nil
			})
	})

	return router
}
```