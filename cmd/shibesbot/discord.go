package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
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

type Bot interface {
	Run() error
	Stop() error
	SetShardOptions(shardID int, shardCount int)
}

type DiscordBot struct {
	Session *discordgo.Session
}

func newBot(token string) (Bot, error) {
	newSession, err := discordgo.New("Bot " + token)

	newSession.AddHandler(commandPicker)
	//for _, cmd := range commands {
	//	newSession.ApplicationCommandCreate(newSession.State.User.ID, "", cmd)
	//}

	return &DiscordBot{
		Session: newSession,
	}, err
}

func (b *DiscordBot) Run() error {
	log.Info("Running bot")
	return b.Session.Open()
}

func (b *DiscordBot) Stop() error {
	log.Info("Closing bot connexion")
	return b.Session.Close()
}

func (b *DiscordBot) SetShardOptions(shardID, shardCount int) {
	log.WithFields(log.Fields{
		"ShardCount": shardCount,
		"ShardID":    shardID,
	}).Info("Updating shard settings")
}

func initDiscord(t string) {
	var err error
	BotSession, err = discordgo.New("Bot " + t)
	if err != nil {
		log.Error("Connexion error: ", err.Error())
		return
	}

	shardCount, err := strconv.Atoi(os.Getenv("SHIBESBOT_SHARD_COUNT"))
	shardID, err := strconv.Atoi(os.Getenv("SHIBESBOT_SHARD_ID"))
	if err != nil {
		shardCount = 1
	}

	BotSession.ShardCount = shardCount
	BotSession.ShardID = shardID
	log.WithFields(log.Fields{
		"ShardCount": BotSession.ShardCount,
		"ShardID":    BotSession.ShardID,
	}).Info("Bot using Shard Mode")
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

	resetUsageCounter()
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
		// incrementPresenceUpdate()
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
	s.UpdateGameStatus(0, "shelp for help || "+
		strconv.Itoa(Users)+" usages this day.")
}
