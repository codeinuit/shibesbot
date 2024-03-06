package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/codeinuit/shibesbot/cmd/shibesbot/monitoring"
	"github.com/codeinuit/shibesbot/pkg/cache"
	"github.com/codeinuit/shibesbot/pkg/cache/localstorage"
	"github.com/codeinuit/shibesbot/pkg/cache/redis"
	"github.com/codeinuit/shibesbot/pkg/logger"
	"github.com/codeinuit/shibesbot/pkg/logger/logrus"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
)

// ENV variables
const (
	// Token configuration
	DISCORD_TOKEN      = "SHIBESBOT_TOKEN"
	SHIBESONLINE_TOKEN = "SHIBESONLINE_TOKEN"

	// Flags
	ENV_CACHE = "CACHE"

	// Redis configuration
	REDIS_ADDR = "REDIS_ADDR"
	REDIS_PORT = "REDIS_PORT"
	REDIS_PASS = "REDIS_PASS"
	REDIS_DB   = "REDIS_DB"
)

type ApiConfigurations struct {
	discordToken     string
	shibesolineToken string
}

type Shibesbot struct {
	session *discordgo.Session

	dailyKey string
	mtx      sync.RWMutex

	apiConfigurations ApiConfigurations
	log               logger.Logger
	cache             cache.Cache
}

func NewShibesbot() (*Shibesbot, error) {
	var cache cache.Cache = localstorage.NewLocalStorageCache()
	var log logger.Logger = logrus.NewLogrusLogger()
	var err error

	// check if Redis is enabled; otherwise fallback to local storage
	if c := os.Getenv(ENV_CACHE); strings.ToUpper(c) == "REDIS" {
		var port int

		log.Info("Redis enabled")

		address := os.Getenv(REDIS_ADDR)
		if port, err = strconv.Atoi(os.Getenv(REDIS_PORT)); err != nil {
			log.Warnf("environnement variable %s is undefined; using default value", REDIS_PORT)
			port = 6379
		}
		log.Infof("using Redis on %s with port %d", address, port)

		cache, err = redis.NewRedisCache(redis.RedisOptions{
			Address:  address,
			Port:     int32(port),
			Password: os.Getenv(REDIS_PASS),
		})
	}

	return &Shibesbot{
		cache: cache,
		log:   log,
		apiConfigurations: ApiConfigurations{
			discordToken:     os.Getenv(DISCORD_TOKEN),
			shibesolineToken: os.Getenv(SHIBESONLINE_TOKEN),
		},
	}, err
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
	sb, err := NewShibesbot()
	if err != nil {
		fmt.Printf("could not initialize bot : %s", err.Error())
		os.Exit(1)
	}

	sb.log.Info("starting bot")
	c := cron.New()
	monitor := monitoring.NewHTTPMonitorServer(sb.log)

	if len(sb.apiConfigurations.discordToken) <= 0 {
		sb.log.Errorf("environnement variable %s is not provided", SHIBESONLINE_TOKEN)
		return
	}

	if err := sb.initDiscord(); err != nil {
		sb.log.Error("connexion error: ", err.Error())
		os.Exit(1)
	}
	defer func() {
		if err := sb.session.Close(); err != nil {
			sb.log.Error("discord session could not close properly:", err.Error())
			os.Exit(1)
		}

		sb.log.Info("discord session closed successfully")
	}()

	_, err = c.AddFunc("0 0 * * *", func() {
		sb.log.Info("updating usage count status")
		sb.setDailyKey(time.Now())
		sb.setDailyCounter(0)
	})

	if err != nil {
		sb.log.Error("could not create cronjob: ", err.Error())
		os.Exit(1)
	}

	c.Start()
	go monitor.Run()
	defer func() {
		monitor.Stop()
		c.Stop()
	}()

	sb.log.Info("shibesbot OK, ready to nicely bork on people")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	sb.log.Info("stop signal has been received, stopping Shibesbot..")
}
