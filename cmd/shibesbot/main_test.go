package main

import (
	"testing"
	"time"

	"github.com/P147x/shibesbot/pkg/cache/localstorage"
	"github.com/P147x/shibesbot/pkg/logger/logrus"
	"github.com/stretchr/testify/assert"
)

func TestResetDailyCounterKey(t *testing.T) {
	bot := &Shibesbot{
		cache: localstorage.NewLocalStorageCache(),
		log:   logrus.NewLogrusLogger(),
	}

	date := time.Unix(1690212965, 0)
	bot.setDailyKey(date)
	assert.Equal(t, "usage:2472023", bot.dailyKey)

	date = time.Unix(1690236000, 0)
	bot.setDailyKey(date)
	assert.Equal(t, "usage:2572023", bot.dailyKey)
}
