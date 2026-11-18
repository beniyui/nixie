package bot

import (
	"nixie/internal/config"

	"github.com/bwmarrin/discordgo"
)

func MSG(s *discordgo.Session, m *discordgo.MessageCreate) {

	//init
	content := m.Content
	channel := m.ChannelID
	author := m.Author.ID
	densu := config.EnvCont.DensuID

	if author == densu {
		s.ChannelMessageSend(channel, content)
	}
}
