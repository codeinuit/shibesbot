package main

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Starting bot.")

	var token string
	if token = os.Getenv("SHIBESBOT_TOKEN"); len(token) == 0 {
		log.Error("Environnement variable SHIBESBOT_TOKEN is not provided")
		return
	}

	bot, err := newBot(token)
	if err != nil {
		log.Error("Error while creating bot instance", err.Error())
	}

	if shardMode := os.Getenv("SHARD_MODE"); shardMode == "true" {
		log.Info("Shard mode enabled")
		// Run Shard manager
	} else {
		log.Info("Shard mode disabled")
		if err = bot.Run(); err != nil {
			log.Error("Error occurred while opening connexion", err.Error())
		}
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
		<-sc
		log.Info("Signal received, stopping in progress")
	}

}
