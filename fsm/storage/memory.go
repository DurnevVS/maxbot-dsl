package fsm

import (
	"context"
	"sync"
)

type MemoryStorage struct {
	mu     sync.RWMutex
	states map[FSMKey]string
	data   map[FSMKey][]byte
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		states: make(map[FSMKey]string),
		data:   make(map[FSMKey][]byte),
	}
}

func (m *MemoryStorage) SetState(ctx context.Context, key FSMKey, state string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.states[key] = state
	return nil
}

func (m *MemoryStorage) GetState(ctx context.Context, key FSMKey) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	state, ok := m.states[key]
	if !ok {
		return AnyState, nil
	}
	return state, nil
}

func (m *MemoryStorage) SetData(ctx context.Context, key FSMKey, data []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[key] = data
	return nil
}

func (m *MemoryStorage) GetData(ctx context.Context, key FSMKey) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	b, ok := m.data[key]
	if !ok {
		return nil, nil
	}
	return b, nil
}

func (m *MemoryStorage) Delete(ctx context.Context, key FSMKey) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.states, key)
	delete(m.data, key)
	return nil
}
