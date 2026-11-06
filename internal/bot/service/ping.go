package service

import (
	"fmt"
	"nixie/internal/bot/config"

	"github.com/bwmarrin/discordgo"
)

func Ping(s *discordgo.Session, m *discordgo.MessageCreate) {
	//トリガーの後ろが空白
	//ping
	statusMsg := fmt.Sprintf("nixieUP  server: %s  version: %s", config.ENV_values.Server, config.ENV_values.Version)
	payload := fmt.Sprintf("```%s```\n%s", config.ENV_values.UpBanner, statusMsg)
	s.ChannelMessageSend(m.ChannelID, payload)
}
