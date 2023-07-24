package cache_test

import (
	"context"
	"testing"

	"github.com/P147x/shibesbot/pkg/cache"
	"github.com/P147x/shibesbot/pkg/cache/localstorage"
	"github.com/stretchr/testify/assert"
)

func TestLocalStorage(t *testing.T) {
	var cache cache.Cache

	const key string = "test"
	const value string = "test"

	cache = localstorage.NewLocalStorageCache()
	cache.Set(context.Background(), key, value)
	awaitedValue, error := cache.Get(context.Background(), key)

	assert.Equal(t, nil, error)
	assert.Equal(t, awaitedValue, value)
}

func TestLocalStorageIncr(t *testing.T) {
	var cache cache.Cache

	const key string = "test"
	const value int = 41
	const newAwaitedValue int = 42

	cache = localstorage.NewLocalStorageCache()

	cache.Set(context.Background(), key, value)
	awaitedValue, error := cache.Incr(context.Background(), key)
	convertedAwaitedValue, ok := awaitedValue.(int)

	assert.True(t, ok)
	assert.Equal(t, nil, error)
	assert.Equal(t, newAwaitedValue, convertedAwaitedValue)
}
