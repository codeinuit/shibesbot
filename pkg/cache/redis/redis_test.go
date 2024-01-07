package redis

import (
	"context"
	"errors"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func NewClientMock(t *testing.T) (*RedisDB, redismock.ClientMock) {
	t.Helper()
	db, mock := redismock.NewClientMock()

	return &RedisDB{
		client: db,
	}, mock
}

func TestGet(t *testing.T) {
	cli, mock := NewClientMock(t)
	key := "key"
	value := "value"
	ctx := context.TODO()

	mock.ExpectGet(key).SetVal(value)
	resp, err := cli.Get(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, value, resp)
}

func TestSet(t *testing.T) {
	cli, mock := NewClientMock(t)
	key := "key"
	value := "value"
	ctx := context.TODO()

	mock.ExpectSet(key, value, 0).SetVal(value)
	setR, err := cli.Set(ctx, key, value)

	assert.Nil(t, err)
	assert.Equal(t, value, setR)

	mock.ExpectSet(key, value, 0).SetErr(errors.New("could not get value"))
	_, err = cli.Set(ctx, key, value)
	assert.NotNil(t, err)
}
