package main

import (
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"

	log "github.com/sirupsen/logrus"
)

var (
	Users           int
	UsageResetTimer time.Timer
	Mtx             sync.Mutex
	BotSession      *discordgo.Session
)

func initDiscord(t string) {
	var err error
	BotSession, err = discordgo.New("Bot " + t)
	if err != nil {
		log.Error("Connexion error: ", err.Error())
		return
	}
	BotSession.AddHandler(shibesHandler)
	err = BotSession.Open()
	if err != nil {
		log.Error("Connexion error : ", err.Error())
		return
	}
	defer BotSession.Close()

	resetUsageCounter()
	log.Info("Bot OK, ready to bork on people")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	log.Info("Signal received, stopping in progress")
}

func commandPicker(s *discordgo.Session, m *discordgo.MessageCreate) (err error) {
	if strings.HasPrefix(m.Content, "s") {
		switch m.Content {
		case "shibes":
			_, err = s.ChannelMessageSend(m.ChannelID, getShibes())
			incrementPresenceUpdate()
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

func resetUsageCounter() {
	t := time.Now()
	n := time.Date(t.Year(), t.Month(), t.Day(), 24, 0, 0, 0, t.Location())
	resetPresenceUpdate()
	log.Info("Reset programmed in ", n.Sub(t).String())
	time.AfterFunc(n.Sub(t), resetUsageCounter)
}

func resetPresenceUpdate() {
	Mtx.Lock()
	defer Mtx.Unlock()
	Users = 0
	updatePresenceUpdate(BotSession)
}

func incrementPresenceUpdate() {
	Mtx.Lock()
	defer Mtx.Unlock()
	Users++
	updatePresenceUpdate(BotSession)
}

func updatePresenceUpdate(s *discordgo.Session) {
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
