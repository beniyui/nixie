package service

import (
	"fmt"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func Down(s *discordgo.Session, m *discordgo.MessageCreate, breaker chan struct{}) {
	serverName := os.Getenv("SERVER")
	densuID := os.Getenv("DENSU_ID")
	cbaID := os.Getenv("CBA_ID")
	inpName := strings.TrimPrefix(m.Content, "nixie down ")
	trgName := strings.ToUpper(inpName)

	if m.Author.ID != densuID && m.Author.ID != cbaID {
		return
	}
	if trgName != "DENSU" && trgName != "CBA" {
		s.ChannelMessageSend(m.ChannelID, "無効：densu?cba?")
		return
	}

	if trgName == serverName {
		msg := fmt.Sprintf("nixieDOWN! server:%s のセッションは終了しました", serverName)
		s.ChannelMessageSend(m.ChannelID, msg)
		close(breaker)
	}
}
