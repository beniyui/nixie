package service

import (
	"fmt"
	"nixie/internal/bot/config"

	"github.com/bwmarrin/discordgo"
)

func Ping(s *discordgo.Session, m *discordgo.MessageCreate) {
	//トリガーの後ろが空白
	//ping
	statusMsg := fmt.Sprintf(
		"nixieUP  server: %s  version: %s",
		config.Values.Server,config.Values.Version,
	)
	payload := fmt.Sprintf("```%s```\n%s", config.Arts.UpBanner, statusMsg)
	s.ChannelMessageSend(m.ChannelID, payload)
}
