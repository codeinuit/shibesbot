package main

import (
	"github.com/bwmarrin/discordgo"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"strings"
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

func shibesHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if strings.HasPrefix(m.Content, "shibes") {
		s.ChannelMessageSend(m.ChannelID, getShibes())
	}
}
