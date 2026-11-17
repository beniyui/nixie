package bot

import (
	"fmt"
	"nixie/internal/config"

	"github.com/bwmarrin/discordgo"
)

func MSG(s *discordgo.Session, m *discordgo.MessageCreate) {

	//init
	content := m.Content
	channel := m.ChannelID
	author := fmt.Sprintf("%v", m.Author.ID)
	densu := config.EnvCont.DensuID

	if author == densu {
		s.ChannelMessageSend(channel, content)
	}
}
