package main

import (
	"github.com/bwmarrin/discordgo"
)

var (
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

	return nil
}

func (sb *Shibesbot) commandPicker(s *discordgo.Session, i *discordgo.InteractionCreate) {
	sb.log.Info("command received: " + i.ApplicationCommandData().Name)
	var response string
	switch i.ApplicationCommandData().Name {
	case "shibes":
		response = sb.getShibes()
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
