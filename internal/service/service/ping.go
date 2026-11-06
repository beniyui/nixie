package service

import (
	"fmt"
	"nixie/internal/tools/tools"
	"os"

	"github.com/bwmarrin/discordgo"
)

func Ping(s *discordgo.Session, m *discordgo.MessageCreate) {
	//トリガーの後ろが空白
	//ping

	statusMsg := fmt.Sprintf("nixieUP  server: %s  version: %s", os.Getenv("SERVER"), os.Getenv("VERSION"))
	pingMsg := fmt.Sprintf("```%s```\n%s", tools.GetArt("up", os.Getenv("ART_VERSION")), statusMsg)
	s.ChannelMessageSend(m.ChannelID, pingMsg)
}
