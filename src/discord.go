package main

import (
	"github.com/bwmarrin/discordgo"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"strings"
	"time"
	"strconv"
)

var (
	Users int
	Time int
)

func initDiscord(t string) {
	dg, err := discordgo.New("Bot " + Token)
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
	if strings.HasPrefix(m.Content, "s") {
		presenceUpdate(s)
		switch m.Content {
		case "shibes":
			_, err = s.ChannelMessageSend(m.ChannelID, getShibes())
			break
		case "sgifs":
			_, err = s.ChannelMessageSend(m.ChannelID, getShibesGifs())
			break
		case "shelp":
			_, err = s.ChannelMessageSendEmbed(m.ChannelID, getHelp())
			break
		case "swalls":
			_, err = s.ChannelMessageSend(m.ChannelID, getShibesWallpaper())
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
	s.UpdateStatus(Users,"shelp for help || " +
		strconv.Itoa(Users) + " usages this day.")
}

func shibesHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if commandPicker(s, m) != nil {
		s.ChannelMessageSend(m.ChannelID, "Oops, something wrong happened. :<")
	}
}
