package fsm

import (
	"context"
	"encoding/json"
)

const AnyState = "*"

type FSMKey struct {
	ChatID int64
	UserID int64
}

type Storage interface {
	SetState(ctx context.Context, key FSMKey, state string) error
	GetState(ctx context.Context, key FSMKey) (string, error)

	SetData(ctx context.Context, key FSMKey, data []byte) error
	GetData(ctx context.Context, key FSMKey) ([]byte, error)

	Delete(ctx context.Context, key FSMKey) error
}

type FSMContext struct {
	storage Storage
	key     FSMKey
}

func NewFSMContext(storage Storage, key FSMKey) *FSMContext {
	return &FSMContext{
		storage: storage,
		key:     key,
	}
}

func (c *FSMContext) SetState(ctx context.Context, state string) error {
	return c.storage.SetState(ctx, c.key, state)
}

func (c *FSMContext) GetState(ctx context.Context) (string, error) {
	return c.storage.GetState(ctx, c.key)
}

func (c *FSMContext) Clear(ctx context.Context) error {
	return c.storage.Delete(ctx, c.key)
}

func (c *FSMContext) SetData(ctx context.Context, v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return c.storage.SetData(ctx, c.key, b)
}

func (c *FSMContext) GetData(ctx context.Context, v any) error {
	b, err := c.storage.GetData(ctx, c.key)
	if err != nil {
		return err
	}
	if len(b) == 0 {
		return nil
	}
	return json.Unmarshal(b, v)
}

func (c *FSMContext) UpdateData(ctx context.Context, m map[string]any) error {
	var data map[string]any

	_ = c.GetData(ctx, &data)

	if data == nil {
		data = map[string]any{}
	}

	for k, v := range m {
		data[k] = v
	}

	return c.SetData(ctx, data)
}

func (c *FSMContext) Is(ctx context.Context, state string) bool {
	s, err := c.GetState(ctx)
	if err != nil {
		return false
	}

	if state == AnyState {
		return true
	}

	return s == state
}
