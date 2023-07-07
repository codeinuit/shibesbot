package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/P147x/shibesbot/pkg/cache"
	"github.com/P147x/shibesbot/pkg/cache/redis"
	"github.com/P147x/shibesbot/pkg/logger"
	"github.com/P147x/shibesbot/pkg/logger/logrus"
	"github.com/bwmarrin/discordgo"
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

func main() {
	sb := initConfiguration()
	sb.initRequests()
	sb.log.Info("Starting Shibesbot")

	if len(sb.apiConfigurations.discordToken) <= 0 {
		sb.log.Error("Environnement variable SHIBESBOT_TOKEN is not provided")
		return
	}

	if err := sb.initDiscord(); err != nil {
		sb.log.Error("Connexion error: ", err.Error())
		return
	}
	defer func() {
		if err := sb.session.Close(); err != nil {
			sb.log.Error("Discord session could not close properly:", err.Error())
			return
		}

		sb.log.Info("Discord session closed successfully")
	}()

	sb.log.Info("Shibesbot OK, ready to nicely bork on people")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	sb.log.Info("Stop signal has been received, stopping Shibesbot..")
}
