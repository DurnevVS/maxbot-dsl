package routers

import (
	"context"

	maxbot "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

type Dispatcher struct {
	routers []*Router
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{}
}

func (d *Dispatcher) AddRouter(r *Router) *Dispatcher {
	d.routers = append(d.routers, r)
	return d
}

func (d *Dispatcher) Dispatch(api *maxbot.Api, update schemes.UpdateInterface, ctx context.Context) error {
	var err error
	for _, r := range d.routers {
		rErr := r.Resolve(api, update, ctx)
		if rErr != nil {
			err = rErr
		}
	}
	return err
}
