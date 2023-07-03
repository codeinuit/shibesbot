package main

import (
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	Users           int
	UsageResetTimer time.Timer
	Mtx             sync.Mutex

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

func (sb *Shibesbot) initDiscord() error {
	var err error

	sb.session, err = discordgo.New("Bot " + sb.apiConfigurations.discordToken)
	if err != nil {
		return err
	}

	sb.session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) { sb.commandPicker(s, i) })
	if err = sb.session.Open(); err != nil {
		return err
	}

	for _, cmd := range commands {
		_, err := sb.session.ApplicationCommandCreate(sb.session.State.User.ID, "", cmd)
		if err != nil {
			return err
		}
	}

	// resetUsageCounter(log)
	return nil
}

func (sb *Shibesbot) commandPicker(s *discordgo.Session, i *discordgo.InteractionCreate) {
	sb.log.Info("command received: " + i.ApplicationCommandData().Name)
	var response string
	switch i.ApplicationCommandData().Name {
	case "shibes":
		response = sb.getShibes()
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

// func resetUsageCounter(log logger.Logger) {
// 	t := time.Now()
// 	n := time.Date(t.Year(), t.Month(), t.Day(), 24, 0, 0, 0, t.Location())
// 	resetPresenceUpdate()
// 	log.Info("Reset programmed in ", n.Sub(t).String())
// 	f := func() { resetUsageCounter(log) }
// 	time.AfterFunc(n.Sub(t), f)
// }

// func resetPresenceUpdate() {
// 	Mtx.Lock()
// 	defer Mtx.Unlock()
// 	Users = 0
// 	updatePresenceUpdate(BotSession)
// }
//
// func incrementPresenceUpdate() {
// 	Mtx.Lock()
// 	defer Mtx.Unlock()
// 	Users++
// 	updatePresenceUpdate(BotSession)
// }

// func updatePresenceUpdate(s *discordgo.Session) {
// 	s.UpdateGameStatus(0, "shelp for help || "+
// 		strconv.Itoa(Users)+" usages this day.")
// }
