package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func RootHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if err := MSG(s, m); err != nil {
		log.Println(err)
	}
}
