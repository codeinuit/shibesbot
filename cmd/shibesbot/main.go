package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/P147x/shibesbot/pkg/cache"
	"github.com/P147x/shibesbot/pkg/cache/redis"
	"github.com/P147x/shibesbot/pkg/logger"
	"github.com/P147x/shibesbot/pkg/logger/logrus"
	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
)

// ENV variables
const (
	DISCORD_TOKEN      = "SHIBESBOT_TOKEN"
	ALPHACODERS_TOKEN  = "ALPHACODERS_TOKEN"
	SHIBESONLINE_TOKEN = "SHIBESONLINE_TOKEN"
	GIPHY_TOKEN        = "GIPHY_TOKEN"

	// Redis configuration
	REDIS_ADDR = "REDIS_ADDR"
	REDIS_PORT = "REDIS_PORT"
	REDIS_PASS = "REDIS_PASS"
	REDIS_DB   = "REDIS_DB"
)

type ApiConfigurations struct {
	discordToken     string
	alphacodersToken string
	shibesolineToken string
	giphyToken       string
}

type Shibesbot struct {
	session *discordgo.Session

	dailyKey string
	mtx      sync.RWMutex

	apiConfigurations ApiConfigurations
	log               logger.Logger
	cache             cache.Cache
}

func initConfiguration() *Shibesbot {
	port, err := strconv.Atoi(os.Getenv(REDIS_PORT))
	if err != nil {
		port = 6379
	}

	r, err := redis.NewRedisCache(redis.RedisOptions{
		Address:  os.Getenv(REDIS_ADDR),
		Port:     int32(port),
		Password: os.Getenv(REDIS_PASS),
	})

	if err != nil {
		log.Fatal(err.Error())
	}

	return &Shibesbot{
		cache: r,
		log:   logrus.NewLogrusLogger(),
		apiConfigurations: ApiConfigurations{
			discordToken:     os.Getenv(DISCORD_TOKEN),
			alphacodersToken: os.Getenv(ALPHACODERS_TOKEN),
			shibesolineToken: os.Getenv(SHIBESONLINE_TOKEN),
			giphyToken:       os.Getenv(GIPHY_TOKEN),
		},
	}
}

func (sb *Shibesbot) setDailyKey(t time.Time) {
	sb.mtx.Lock()
	defer sb.mtx.Unlock()

	key := fmt.Sprintf("usage:%d%d%d", t.Day(), t.Month(), t.Year())

	isUnset, err := sb.cache.SetNX(context.Background(), key, 0)
	if err != nil {
		sb.log.Warn("could not update and retrieve usage count: ", err.Error())
		return
	}

	if isUnset {
		sb.dailyKey = key

		return
	}
}

func main() {
	sb := initConfiguration()
	sb.initRequests()
	sb.log.Info("starting Shibesbot")
	c := cron.New()

	if len(sb.apiConfigurations.discordToken) <= 0 {
		sb.log.Error("environnement variable SHIBESBOT_TOKEN is not provided")
		return
	}

	if err := sb.initDiscord(); err != nil {
		sb.log.Error("connexion error: ", err.Error())
		return
	}
	defer func() {
		if err := sb.session.Close(); err != nil {
			sb.log.Error("discord session could not close properly:", err.Error())
			return
		}

		sb.log.Info("discord session closed successfully")
	}()

	_, err := c.AddFunc("0 0 * * *", func() {
		sb.log.Info("updating usage count status")
		sb.setDailyKey(time.Now())
		sb.setDailyCounter(0)
	})

	if err != nil {
		sb.log.Error("could not create cronjob: ", err.Error())
		return
	}

	c.Start()
	defer c.Stop()

	sb.log.Info("shibesbot OK, ready to nicely bork on people")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	sb.log.Info("stop signal has been received, stopping Shibesbot..")
}
