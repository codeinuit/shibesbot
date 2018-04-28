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
)

var (
	Token string
)

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "shibes") {
		if m.Content == "shibes" {
			resp, err := http.Get("http://shibe.online/api/shibes/")
			if err == nil {
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				if err == nil {

				}
				var u []string
				json.Unmarshal(body, &u)
				s.ChannelMessageSend(m.ChannelID, u[0])
			}
		} else if m.Content == "shibeshelp" {
			test := &discordgo.MessageEmbed{
				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: "http://img.over-blog-kiwi.com/1/47/73/14/20160709/ob_bcc896_chiot-shiba-inu-a-vendre-2016.jpg",
				},
				Description: "Thanks for using Shibesbot on your Discord server !\n\n" +
					"Our purpose is to distribute many **shibes** on your server, using http://shibes.online/ as puppy distributor.",
				Title: "Hello shibes !",
			}
			s.ChannelMessageSendEmbed(m.ChannelID, test)
		} else if m.Content == "shibesgif" {
			gp := giphy.New("PcVZFoFsmh2vhFHqSKjhvbnwq74N7JSi")
			gifs, err := gp.Search("shiba")
			if err != nil {
				return
			}
			s.ChannelMessageSend(m.ChannelID, gifs[rand.Int() % len(gifs)].URL)
		}
	}
}
