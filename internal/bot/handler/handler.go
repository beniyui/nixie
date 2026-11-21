package handler

import (
	"nixie/internal/config"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type cmdID int

const (
	n cmdID = iota //n, none, do nothing
	b
)

func handler(s *discordgo.Session, m *discordgo.MessageCreate) error {
	id, err := msgListener(s, m)
	if err != nil {
		return err
	}
	switch id {

	case n:
		return nil

	}
	return nil
}

// return id
func msgListener(s *discordgo.Session, m *discordgo.MessageCreate) (id cmdID, err error) {

	//init
	content := m.Content
	channel := m.ChannelID
	authorID := m.Author.ID

	if authorID == config.EnvCont.DensuID {
		if strings.HasPrefix(content, "nixie") {
			s.ChannelMessageSend(channel, "triggered.")

		}
	}
	return n, nil
}
