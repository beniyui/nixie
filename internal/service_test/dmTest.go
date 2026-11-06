package service_test

import (
	"os"

	"github.com/bwmarrin/discordgo"
)

func Dmtest(s *discordgo.Session,m *discordgo.MessageCreate)error {
	err := SendDM(s, m, os.Getenv("DENSU_ID"), "https://discord.gg/QCGCfqxEC3")
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "キャンセル")
		return err
	}
	return nil
}
