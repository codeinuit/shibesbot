package cache_test

import (
	"context"
	"testing"

	"github.com/codeinuit/shibesbot/pkg/cache"
	"github.com/codeinuit/shibesbot/pkg/cache/localstorage"
	"github.com/stretchr/testify/assert"
)

func TestLocalStorage(t *testing.T) {
	var cache cache.Cache

	const key string = "test"
	const value string = "test"

	cache = localstorage.NewLocalStorageCache()
	setValue, error := cache.Set(context.Background(), key, value)
	assert.Nil(t, error)
	assert.Equal(t, value, setValue)

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

	setValue, error := cache.Set(context.Background(), key, value)
	assert.Nil(t, error)
	assert.Equal(t, value, setValue)

	awaitedValue, error := cache.Incr(context.Background(), key)
	convertedAwaitedValue, ok := awaitedValue.(int)

	assert.True(t, ok)
	assert.Equal(t, nil, error)
	assert.Equal(t, newAwaitedValue, convertedAwaitedValue)
}
