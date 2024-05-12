package main

import (
	"net/http"

	"encoding/json"

	"github.com/bwmarrin/discordgo"
)

var (
	Shibes ShibesData
)

type ShibesPictures struct {
	Shibes []string
	Total  int
	Cursor int
}

type ShibesData struct {
	Images ShibesPictures
}

func (sb *Shibesbot) getShibes() string {
	if Shibes.Images.Cursor >= Shibes.Images.Total {
		Shibes.Images.Cursor = 0
		Shibes.Images.Total = 10
		resp, err := http.Get("http://shibe.online/api/shibes?count=10")
		if err != nil {
			sb.log.Warn("could not get images from shibes.online: ", err.Error())
			return ""
		}
		defer resp.Body.Close()
		err = json.NewDecoder(resp.Body).Decode(&Shibes.Images.Shibes)
		if err != nil {
			sb.log.Warn("could not get images from shibes.online: ", err.Error())
			return ""
		}
		sb.log.Info("Updated ", Shibes.Images.Total, " pictures from shibes.online")
	}
	Shibes.Images.Cursor++
	return Shibes.Images.Shibes[Shibes.Images.Cursor-1]
}

func getHelp() *discordgo.MessageEmbed {
	test := &discordgo.MessageEmbed{
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "http://img.over-blog-kiwi.com/1/47/73/14/20160709/ob_bcc896_chiot-shiba-inu-a-vendre-2016.jpg",
		},
		Description: "Thanks for using Shibesbot on your Discord server !\n\n" +
			"Our purpose is to distribute many **shibes** on your server, using http://shibe.online/ as puppy distributor.\n\n" +
			"Are you enjoying this bot ? You can help us spread the doge ! https://github.com/P147x/discord-shibesbot",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name: "Available commands",
				Value: "- *shibes* to get a random shibe !\n" +
					"- *shelp* to get help",
				Inline: false,
			},
		},
		Title: "Hello shibes !",
	}
	return test
}
