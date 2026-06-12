package localstorage

import (
	"context"
	"errors"
	"sync"
)

type LocalStorage struct {
	array map[string]any
	rwm   sync.RWMutex
}

func NewLocalStorageCache() *LocalStorage {
	return &LocalStorage{
		array: make(map[string]any),
	}
}

func (ls *LocalStorage) Get(ctx context.Context, k string) (any, error) {
	ls.rwm.RLock()
	defer ls.rwm.RUnlock()

	value, ok := ls.array[k]
	if !ok {
		return nil, errors.New("could not get value")
	}

	return value, nil
}

func (ls *LocalStorage) Set(ctx context.Context, k string, v any) (any, error) {
	ls.rwm.Lock()
	defer ls.rwm.Unlock()

	ls.array[k] = v

	return v, nil
}

func (ls *LocalStorage) Incr(ctx context.Context, k string) (any, error) {
	ls.rwm.Lock()
	defer ls.rwm.Unlock()

	if ls.array[k] == nil {
		ls.array[k] = 1
		return 1, nil
	}
	convertedValue, ok := ls.array[k].(int)
	if !ok {
		return nil, errors.New("could not increment value")
	}

	ls.array[k] = convertedValue + 1
	return ls.array[k], nil
}

func (ls *LocalStorage) SetNX(ctx context.Context, k string, v any) (bool, error) {
	ls.rwm.Lock()
	defer ls.rwm.Unlock()

	if ls.array[k] != nil {
		ls.array[k] = v

		return false, nil
	}

	return true, nil
}
