package main

import (
	"os"
	
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Starting bot.")
	
	var token string
	if token = os.Getenv("SHIBESBOT_TOKEN"); len(token) == 0 {
		log.Error("Environnement variable SHIBESBOT_TOKEN is not provided")
		return
	}

	initDiscord(token)
}
