package handler

import (
	"fmt"
	"slices"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type cmdID int

const (
	wrong cmdID = iota
	n           //n, none, do nothing
	test
)

func handler(s *discordgo.Session, m *discordgo.MessageCreate) error {
	id := msgListener(s, m)
	switch id {
	case wrong:
		return fmt.Errorf("id is something wrong")
	case test:
		s.ChannelMessageSend(m.ChannelID, "nixie triggered.")
	}

	return nil
}

// return id
func msgListener(s *discordgo.Session, m *discordgo.MessageCreate) (id cmdID) {
	if m.Author.Bot {
		return n
	}
	switch hasTrigger(m.Content) {
	case true:
	case false:
		return n
	default:
		return wrong

	}

	return
}

var triggerWords = []string{"nixie", "ニキシ", "にきし", "ねえさん"}

func hasTrigger(s string) bool {
	trg := strings.ToLower(s)

	f := func(w string) bool {
		return strings.Contains(trg, w)
	}
	return slices.ContainsFunc(triggerWords, f)

	// for _, str := range triggerWords {
	// 	if strings.Contains(trg, str) {
	// 		return true
	// 	}
	// }
	// return false
}
