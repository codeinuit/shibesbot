package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/P147x/shibesbot/pkg/logger"
	"github.com/bwmarrin/discordgo"
)

var (
	Users           int
	UsageResetTimer time.Timer
	Mtx             sync.Mutex
	BotSession      *discordgo.Session

	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "shibes",
			Description: "Returns an image of a Shiba",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "count",
					Description: "Ask for more pictures in one request",
					Required:    false,
				},
			},
		},
		{
			Name:        "swalls",
			Description: "Returns a wallpaper with a Shiba ina !",
		},
		{
			Name:        "sgifs",
			Description: "Returns a gif with a Shiba ina !"},
		{
			Name:        "shelp",
			Description: "Returns helper",
		},
	}
)

func initDiscord(log logger.Logger, t string) {
	var err error
	BotSession, err = discordgo.New("Bot " + t)
	if err != nil {
		log.Error("Connexion error: ", err.Error())
		return
	}
	BotSession.AddHandler(commandPicker)
	err = BotSession.Open()
	if err != nil {
		log.Error("Connexion error : ", err.Error())
		return
	}
	defer BotSession.Close()
	for _, cmd := range commands {
		BotSession.ApplicationCommandCreate(BotSession.State.User.ID, "", cmd)
	}

	resetUsageCounter(log)
	log.Info("Bot OK, ready to bork on people")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	log.Info("Signal received, stopping in progress")
}

func commandPicker(s *discordgo.Session, i *discordgo.InteractionCreate) {
	fmt.Println("command received: " + i.ApplicationCommandData().Name)
	var response string
	switch i.ApplicationCommandData().Name {
	case "shibes":
		response = getShibes()
		incrementPresenceUpdate()
	case "sgifs":
		response = getShibesGifs()
	case "shelp":
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{getHelp()},
			},
		})
		return
	case "swalls":
		response = getShibesWallpaper()
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
}

func resetUsageCounter(log logger.Logger) {
	t := time.Now()
	n := time.Date(t.Year(), t.Month(), t.Day(), 24, 0, 0, 0, t.Location())
	resetPresenceUpdate()
	log.Info("Reset programmed in ", n.Sub(t).String())
	f := func() { resetUsageCounter(log) }
	time.AfterFunc(n.Sub(t), f)
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
	s.UpdateGameStatus(0, "shelp for help || "+
		strconv.Itoa(Users)+" usages this day.")
}
