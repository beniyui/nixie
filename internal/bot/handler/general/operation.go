package handler

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

// this is bot's rooter.
func Operation(s *discordgo.Session, m *discordgo.MessageCreate) {
	err := handler(s, m)
	if err != nil {
		log.Printf("[error] %v", err)
		s.ChannelMessageSend(m.ChannelID,fmt.Sprintf("えらー。```%v```",err))
	}

	//life is not forever.
}

