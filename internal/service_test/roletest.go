package service_test

import (
	"fmt"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func Rolelisttest(s *discordgo.Session, m *discordgo.MessageCreate) {
	rolelist, err := getUserRoles(s, m.GuildID, os.Getenv("DENSU_ID"))
	if err != nil {
		msg := fmt.Sprintf("%v", err)
		s.ChannelMessageSend(m.ChannelID, msg)
	}
	rolelistTxt := strings.Join(rolelist, "\n")

	s.ChannelMessageSend(m.ChannelID, rolelistTxt)
	jsonTxt := save()
	s.ChannelMessageSend(m.ChannelID, jsonTxt)
}
