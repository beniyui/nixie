package handler

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)



func handler(s *discordgo.Session, m *discordgo.MessageCreate) error {
	id := msgListener(m)
	switch id {
	case WRONG:
		return fmt.Errorf("id is something wrong")
	case NONE:
		return nil
	case AI:

	}

	return nil
}
