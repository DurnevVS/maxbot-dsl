package routers

import (
	"context"

	fsm "github.com/DurnevVS/maxbot-dsl/fsm/storage"
	maxbot "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

type Dispatcher struct {
	routers []*Router
	storage fsm.Storage
}

func NewDispatcher(storage fsm.Storage) *Dispatcher {
	return &Dispatcher{storage: storage}
}

func (d *Dispatcher) buildFSM(update schemes.UpdateInterface) *fsm.FSMContext {
	key := fsm.FSMKey{
		ChatID: update.GetChatID(),
		UserID: update.GetUserID(),
	}

	return fsm.NewFSMContext(d.storage, key)
}

func (d *Dispatcher) AddRouter(r *Router) *Dispatcher {
	d.routers = append(d.routers, r)
	return d
}

func (d *Dispatcher) AddRouters(rs ...*Router) *Dispatcher {
	d.routers = append(d.routers, rs...)
	return d
}

func (d *Dispatcher) Dispatch(api *maxbot.Api, update schemes.UpdateInterface, ctx context.Context) error {
	fsm := d.buildFSM(update)
	for _, r := range d.routers {
		handled, err := r.Resolve(api, update, ctx, fsm)
		if err != nil {
			return err
		}
		if handled {
			return nil
		}
	}
	return nil
}
