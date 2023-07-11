package main

import (
	"context"
	"fmt"
	"strconv"

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
	sb.session.AddHandlerOnce(func(s *discordgo.Session, i *discordgo.Ready) {
		count, err := sb.cache.Get(context.Background(), sb.dailyKey)
		if err != nil {
			sb.log.Warn("could not get daily counter from cache : ", err.Error())
			return
		}

		countString, ok := count.(string)
		if !ok {
			sb.log.Warn("could not get daily counter from cache : conversion error")
			return
		}
		countInt, err := strconv.Atoi(countString)
		if err != nil {
			sb.log.Warn("could not get daily counter from cache : ", err.Error())
			return
		}

		sb.setDailyCounter(int64(countInt))
	})

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
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{getHelp()},
			},
		})
		if err != nil {
			sb.log.Error("could not answer to user help command: ", err.Error())
		}
		return
	case "swalls":
		response = getShibesWallpaper()
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})

	if err != nil {
		sb.log.Error("could not answer to user help command: ", err.Error())
		return
	}

	sb.updateDailyCounter()
}

func (sb *Shibesbot) updateDailyCounter() {
	sb.mtx.RLock()
	defer sb.mtx.RUnlock()
	count, err := sb.cache.Incr(context.Background(), sb.dailyKey)
	if err != nil {
		sb.log.Warn("could not get daily counter from cache : ", err.Error())
		return
	}
	countInt, ok := count.(int64)
	if !ok {
		sb.log.Warn("could not get daily counter from cache")
		return
	}

	sb.setDailyCounter(countInt)
}

func (sb *Shibesbot) setDailyCounter(count int64) {
	err := sb.session.UpdateGameStatus(0, fmt.Sprintf("used %d times today", count))
	if err != nil {
		sb.log.Warn("could not update daily counter: ", err.Error())
	}
}
