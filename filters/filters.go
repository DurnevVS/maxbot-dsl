package filters

import (
	"context"
	"strings"

	"github.com/DurnevVS/maxbot-dsl/handlers"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

func IsCommand(cmd string) handlers.FilterFunc[*schemes.MessageCreatedUpdate] {
	return func(update *schemes.MessageCreatedUpdate, ctx context.Context) (bool, error) {
		text := update.Message.Body.Text
		if strings.HasPrefix(text, cmd) {
			return true, nil
		}
		return false, nil
	}
}
