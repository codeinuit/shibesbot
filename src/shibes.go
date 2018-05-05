package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"github.com/bwmarrin/discordgo"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strings"
	"github.com/ivolo/go-giphy"
	"math/rand"
	"unicode"
	"strconv"
)

var (
	Token string
)

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func getShibesWallpaper(s string) string {
	req, err := http.NewRequest("GET", "https://wall.alphacoders.com/api2.0/get.php", nil)
    if err != nil {
    }

    q := req.URL.Query()
    q.Add("auth", "c8b66cee6ef7022a615da5cbba315f3c")
    q.Add("method", "search")
		q.Add("term", s)
    req.URL.RawQuery = q.Encode()
		resp, err := http.Get(req.URL.String())
		if err != nil {
			fmt.Print("ERR")
}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Print("ERR")
		}

		var res AlphacodersData
		json.Unmarshal(body, &res)
		return string(res.Wallpapers[rand.Int()%len(res.Wallpapers)].Url_Image)

}

type WallpaperData struct {
	Id int
	Width int
	Height int
	Url_Image string
}

type AlphacodersData struct {
	Success bool
	Wallpapers []WallpaperData
	Total_Match int
}


func isInt(s string) bool {
    for _, c := range s {
        if !unicode.IsDigit(c) {
            return false
        }
    }
    return true
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "shibes") {
		if (strings.Split(m.Content, " "))[0] == "shibes" {
			nbr, _ := strconv.Atoi(strings.Split(m.Content, " ")[1])
			fmt.Print(isInt(strings.Split(m.Content, " ")[1]))
			if isInt(strings.Split(m.Content, " ")[1]) == true && nbr < 10 {
				resp, err := http.Get("http://shibe.online/api/shibes?count=" + strings.Split(m.Content, " ")[1])

				if err == nil {
					defer resp.Body.Close()
					body, _ := ioutil.ReadAll(resp.Body)
					var u []string
					json.Unmarshal(body, &u)
					s.ChannelMessageSend(m.ChannelID, strings.Join(u, " "))
				}
			} else if isInt(strings.Split(m.Content, " ")[1]) == true && nbr != 0{
				s.ChannelMessageSend(m.ChannelID, "You can't call more than 10 shibes at the time.")
			} else if nbr <= 0 {
				s.ChannelMessageSend(m.ChannelID, "Dude, what. You can't do that.")
			} else {
				resp, err := http.Get("http://shibe.online/api/shibes")
				if err == nil {
					defer resp.Body.Close()
					body, err := ioutil.ReadAll(resp.Body)
					if err == nil {

					}
					var u []string
					json.Unmarshal(body, &u)
					s.ChannelMessageSend(m.ChannelID, u[0])
				}
			}

		} else if m.Content == "shibeshelp" {
			test := &discordgo.MessageEmbed{
				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: "http://img.over-blog-kiwi.com/1/47/73/14/20160709/ob_bcc896_chiot-shiba-inu-a-vendre-2016.jpg",
				},
				Description: "Thanks for using Shibesbot on your Discord server !\n\n" +
					"Our purpose is to distribute many **shibes** on your server, using http://shibe.online/ as puppy distributor.\n\n",
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Available commands",
						Value:  "- *shibes* to get a random shibe !\n" +
							"- *shibesgif* to get a random gif of shiba !\n" +
							"- *shibeshelp* to get help",
						Inline: false,
					},
				},
				Title: "Hello shibes !",
			}
			s.ChannelMessageSendEmbed(m.ChannelID, test)
		} else if m.Content == "shibesgif" {
			gp := giphy.New("PcVZFoFsmh2vhFHqSKjhvbnwq74N7JSi")
			gifs, err := gp.Search("shiba")
			if err != nil {
				return
			}
			s.ChannelMessageSend(m.ChannelID, gifs[rand.Int()%len(gifs)].URL)
		} else if m.Content == "shibeswallpaper" {
				s.ChannelMessageSend(m.ChannelID, string(getShibesWallpaper("Shiba")))
			}

		}
	}
