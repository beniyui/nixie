package handler

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

// this is bot's rooter.
func Operation(s *discordgo.Session, m *discordgo.MessageCreate) {
	err := handler(s, m)
	if err != nil {
		log.Printf("[error] %v", err)
	}

	//life is not forever.
}
