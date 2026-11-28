package handler

import (
	"nixie/internal/bot/pkg/command"

	"github.com/bwmarrin/discordgo"
)

func handler(s *discordgo.Session, m *discordgo.MessageCreate) error {
	cmdCont, err := command.Listener(m.Content)

	return nil
}
