package command

import (
	"log"
	"strings"
)

func cmdReader(rawStr string, trig triggerCont) ([]cmdBody, error) {
	var (
		bufBody cmdBody
	)

	if trig.Type == WORD {
		bufBody.ID = AI
		bufBody.Arg = rawStr
		return []cmdBody{bufBody}, nil
	}

	prefixTrimedStr := strings.TrimPrefix(rawStr, trig.Key)
	spaceStrSlice := strings.Split(prefixTrimedStr, " ")
	firstBlock := spaceStrSlice[0]
	// LoweredFirstBlock := strings.ToLower(firstBlock)
	optionsStrSlice := strings.Split(firstBlock, "--")

	switch Options {
	case "ping", "":
		bufCont.Options.ID = PING
		return bufCont, nil
	case "think":
		bufCont.Options.ID = AI
		msg := strings.TrimPrefix(prefixTrimedStr, optionsStrSlice[0])
		log.Println(msg)
		bufCont.Options.Arg = []string{msg}
		return bufCont, nil
	default:
		bufCont.Options.ID = AI
		bufCont.Options.Arg = []string{prefixTrimedStr}
		return bufCont, nil
	}

}

func optionsReader(str []string) []cmdBody {
	var bufSlice []cmdBody
	for _, value := range str {

		buf := strings.Split(value, "=")
		elemsReader(buf)

		append(bufSlice)
	}
}

func elemsReader(elem []string) cmdBody {
	var bufBody cmdBody
	if len(elem) == 2 {
		bufBody.Arg = {}
	}
}
