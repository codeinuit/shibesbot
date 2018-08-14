package main

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/ivolo/go-giphy"
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

func init() {
	req, err := http.NewRequest("GET", "https://wall.alphacoders.com/api2.0/get.php", nil)
	if err != nil {
	}
	q := req.URL.Query()
	q.Add("auth", "c8b66cee6ef7022a615da5cbba315f3c")
	q.Add("method", "search")
	q.Add("term", "Shiba")
	req.URL.RawQuery = q.Encode()
	resp, _ := http.Get(req.URL.String())
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var res AlphacodersData
	json.Unmarshal(body, &res)
	Shibes.Wallpapers.Shibes = make([]WallpaperData, len(res.Wallpapers))
	Shibes.Wallpapers.Shibes = res.Wallpapers
	Shibes.Wallpapers.Total = len(res.Wallpapers)

	gp := giphy.New("PcVZFoFsmh2vhFHqSKjhvbnwq74N7JSi")
	Shibes.Gifs.Shibes, _ = gp.Search("shiba")
	Shibes.Gifs.Total = len(Shibes.Gifs.Shibes)
	Shibes.Gifs.Cursor = 0
}

func getShibes() string {
	if Shibes.Images.Cursor >= Shibes.Images.Total {
		Shibes.Images.Cursor = 0
		Shibes.Images.Total = 10
		resp, err := http.Get("http://shibe.online/api/shibes?count=10")
		if err == nil {
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			json.Unmarshal(body, &Shibes.Images.Shibes)
		}
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
	return Shibes.Gifs.Shibes[rand.Int()%Shibes.Gifs.Total].URL
}

func getShibesWallpaper() string {
	return string(Shibes.Wallpapers.Shibes[rand.Int()%Shibes.Wallpapers.Total].Url_Image)
}
