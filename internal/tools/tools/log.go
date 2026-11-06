package tools

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func LogToDiscord(s *discordgo.Session, msg string) {
	logChanId := "1551040385172905985"
	finalMsg := fmt.Sprintf("```%s```",msg)
	s.ChannelMessageSend(logChanId, finalMsg)
}
