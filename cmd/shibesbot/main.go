package main

import (
	"os"

	"github.com/P147x/shibesbot/pkg/logger"
	"github.com/P147x/shibesbot/pkg/logger/logrus"
)

func main() {
	var log logger.Logger

	log = logrus.NewLogrusLogger()
	log.Info("Starting Shibesbot")

	var token string
	if token = os.Getenv("SHIBESBOT_TOKEN"); len(token) == 0 {
		log.Error("Environnement variable SHIBESBOT_TOKEN is not provided")
		return
	}

	initDiscord(log, token)
}
