package handler

import (
	"fmt"
	"nixie/internal/config"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type cmdID int

const (
	WRONG cmdID = iota
	N           //n, none, do nothing
	AI
)

func handler(s *discordgo.Session, m *discordgo.MessageCreate) error {
	id := msgListener(m)
	switch id {
	case WRONG:
		return fmt.Errorf("id is something wrong")
	case N:
		return nil
	case AI:

	}

	return nil
}

// return id
func msgListener(m *discordgo.MessageCreate) (id cmdID) {
	if m.Author.Bot {
		return N
	}
	if m.Author.ID != config.EnvCont.DensuID {
		return N
	}

	switch hasTrigger(m.Content) {
	case true:
		return AI
	case false:
		return N
	default:
		return WRONG

	}

	return
}

var triggerWords = []string{"nixie", "ニキシ", "にきし", "ねえさん"}
var triggerPrefixes = []string{"n!"}

func hasTrigger(s string) bool {
	trg := strings.ToLower(s)

	for _, str := range triggerWords {
		if strings.Contains(trg, str) {
			return true
		}
	}

	for _, str := range triggerPrefixes {
		if strings.HasPrefix(trg, str) {
			return true
		}
	}
	return false
}
