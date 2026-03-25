package filters

import (
	"context"
	"strings"

	fsm "github.com/DurnevVS/maxbot-dsl/fsm/storage"
	"github.com/DurnevVS/maxbot-dsl/routers"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

func StateFilter[T schemes.UpdateInterface](state string) routers.FilterFunc[T] {
	return func(update T, ctx context.Context, fsm *fsm.FSMContext) (bool, error) {
		return fsm.Is(ctx, state), nil
	}
}

func IsCommand(cmd string) routers.FilterFunc[*schemes.MessageCreatedUpdate] {
	return func(update *schemes.MessageCreatedUpdate, ctx context.Context, fsm *fsm.FSMContext) (bool, error) {
		text := update.Message.Body.Text
		if strings.HasPrefix(text, cmd) {
			return true, nil
		}
		return false, nil
	}
}
