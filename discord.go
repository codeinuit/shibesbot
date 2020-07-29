package main

import (
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	
	log "github.com/sirupsen/logrus"
)

var (
	Users int
	Time  int
)

func initDiscord(t string) {
	dg, err := discordgo.New("Bot " + t)
	if err != nil {
		log.Error("Connexion error: ", err.Error())
		return
	}

	dg.AddHandler(shibesHandler)
	err = dg.Open()
	if err != nil {
		log.Error("Connexion error : ", err.Error())
		return
	}

	log.Info("Bot OK, ready to bork on people")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	log.Info("Signal received, stopping in progress")
	dg.Close()
}

func commandPicker(s *discordgo.Session, m *discordgo.MessageCreate) (err error) {
	if strings.HasPrefix(m.Content, "s") {
		switch m.Content {
		case "shibes":
			_, err = s.ChannelMessageSend(m.ChannelID, getShibes())
			presenceUpdate(s)
			break
		case "sgifs":
			_, err = s.ChannelMessageSend(m.ChannelID, getShibesGifs())
			break
		case "shelp":
			_, err = s.ChannelMessageSendEmbed(m.ChannelID, getHelp())
			break
		case "swalls":
			_, err = s.ChannelMessageSend(m.ChannelID, getShibesWallpaper())
			break
		}
	}
	return err
}

func presenceUpdate(s *discordgo.Session) {
	if Time != int(time.Now().Day()) {
		Time = int(time.Now().Day())
		Users = 0
	}
	Users++
	s.UpdateStatus(Users, "shelp for help || "+
		strconv.Itoa(Users)+" usages this day.")
}

func shibesHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if commandPicker(s, m) != nil {
		s.ChannelMessageSend(m.ChannelID, "Oops, something wrong happened. :<")
	}
}
