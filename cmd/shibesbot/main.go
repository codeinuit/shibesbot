package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/P147x/shibesbot/pkg/logger"
	"github.com/P147x/shibesbot/pkg/logger/logrus"
	"github.com/bwmarrin/discordgo"
)

const (
	DISCORD_TOKEN      = "SHIBESBOT_TOKEN"
	ALPHACODERS_TOKEN  = "ALPHACODERS_TOKEN"
	SHIBESONLINE_TOKEN = "SHIBESONLINE_TOKEN"
	GIPHY_TOKEN        = "GIPHY_TOKEN"
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
}

func initConfiguration() *Shibesbot {
	return &Shibesbot{
		log: logrus.NewLogrusLogger(),
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
