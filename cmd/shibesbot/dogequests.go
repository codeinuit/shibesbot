package main

import (
	"math/rand"
	"net/http"

	"github.com/ivolo/go-giphy"

	"encoding/json"

	"github.com/bwmarrin/discordgo"
)

var (
	Shibes ShibesData
)

type WallpaperData struct {
	Id        int
	Width     int
	Height    int
	Url_Image string
}

type AlphacodersData struct {
	Success     bool
	Wallpapers  []WallpaperData
	Total_Match int
}

type ShibesPictures struct {
	Shibes []string
	Total  int
	Cursor int
}

type ShibesGifs struct {
	Shibes []giphy.Gif
	Total  int
	Cursor int
}

type ShibesWallpapers struct {
	Shibes []WallpaperData
	Total  int
	Cursor int
}

type ShibesData struct {
	Images     ShibesPictures
	Gifs       ShibesGifs
	Wallpapers ShibesWallpapers
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
					"- *sgifs* to get a random gif of shiba !\n" +
					"- *shelp* to get help\n" +
					"- *swalls* to get an amazing shibe wallpaper",
				Inline: false,
			},
		},
		Title: "Hello shibes !",
	}
	return test
}

func getShibesGifs() string {
	if Shibes.Gifs.Total <= 0 {
		return "no gifs available, sorry. :("
	}
	return Shibes.Gifs.Shibes[rand.Int()%Shibes.Gifs.Total].URL
}

func getShibesWallpaper() string {
	if Shibes.Wallpapers.Total <= 0 {
		return "no wallpapers available, sorry. :("
	}
	return string(Shibes.Wallpapers.Shibes[rand.Int()%Shibes.Wallpapers.Total].Url_Image)
}
