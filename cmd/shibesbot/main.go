package main

import (
	"os"

	"github.com/P147x/shibesbot/pkg/logger"
	"github.com/P147x/shibesbot/pkg/logger/logrus"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Starting bot.")
	var logger logger.Logger

	logger = logrus.NewLogrusLogger()

	var token string
	if token = os.Getenv("SHIBESBOT_TOKEN"); len(token) == 0 {
		log.Error("Environnement variable SHIBESBOT_TOKEN is not provided")
		return
	}

	initDiscord(token)
}
