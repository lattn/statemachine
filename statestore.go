package statemachine

import (
	"context"
	"fmt"
	"sync"
)

var _ StateStore[string, int] = (*MemoryStateStore[string, int])(nil)

type StateStore[K comparable, T any] interface {
	Has(ctx context.Context, key K) (bool, error)
	Get(ctx context.Context, key K) (T, error)
	Set(ctx context.Context, key K, state T) error
}

type MemoryStateStore[K comparable, T any] struct {
	m sync.Map
}

func (m *MemoryStateStore[K, T]) Has(ctx context.Context, key K) (bool, error) {
	_, ok := m.m.Load(key)
	return ok, nil
}

func (m *MemoryStateStore[K, T]) Get(ctx context.Context, key K) (T, error) {
	value, ok := m.m.Load(key)
	if !ok {
		var zero T
		return zero, fmt.Errorf("key not found: %v", key)
	}
	return value.(T), nil
}

func (m *MemoryStateStore[K, T]) Set(ctx context.Context, key K, state T) error {
	m.m.Store(key, state)
	return nil
}

type storedState[K comparable, T any] struct {
	key K
	ss  StateStore[K, T]
}

func (s *storedState[K, T]) Get(ctx context.Context) (T, error) {
	return s.ss.Get(ctx, s.key)
}

func (s *storedState[K, T]) Set(ctx context.Context, state T) error {
	return s.ss.Set(ctx, s.key, state)
}
