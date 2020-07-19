package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	Users int
	Time  int
)

func initDiscord(t string) {
	dg, err := discordgo.New("Bot " + t)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(shibesHandler)
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Shibes ready for duty. Press CTRL+C for no shibes.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	dg.Close()
}

func commandPicker(s *discordgo.Session, m *discordgo.MessageCreate) error {
	var err error
	if strings.HasPrefix(m.Content, "p") {
		switch m.Content {
		case "papainperdu":
			_, err = s.ChannelMessageSend(m.ChannelID, getShibes())
			presenceUpdate(s)
			break
		case "papainperdumaisanime":
			_, err = s.ChannelMessageSend(m.ChannelID, getShibesGifs())
			break
		case "painperdualaide":
			_, err = s.ChannelMessageSendEmbed(m.ChannelID, getHelp())
			break
		case "painperdumur":
			_, err = s.ChannelMessageSend(m.ChannelID, getShibesWallpaper())
			break
		case "painperduquiestreslongaecriremaiscavautpeutetrelecoup":
			_, _ = s.ChannelMessageSend(m.ChannelID, "https://giphy.com/stickers/imoji-wtf-l4FGAwUT9YmqqAZyg")
			break
		}
	}
	return err
}

func presenceUpdate(s *discordgo.Session) {
	if Time != int(time.Now().Day()) {
		Time = int(time.Now().Day())
		Users = 0
	}
	Users++
	s.UpdateStatus(Users, "shelp for help || "+
		strconv.Itoa(Users)+" usages this day.")
}

func shibesHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if commandPicker(s, m) != nil {
		s.ChannelMessageSend(m.ChannelID, "Oops, something wrong happened. :<")
	}
}
