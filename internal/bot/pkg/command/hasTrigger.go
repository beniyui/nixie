package command

import "strings"

//call this func from parser
func hasTrigger(rawString string) (hasTrigger bool, trig triggerBody) {
	s := strings.ToLower(rawString)
	if b, trig := hasTriggerWord(s); b == true {
		return true, triggerBody{Type: WORD, Key: trig}
	}
	if b, trig := hasTriggerPrefix(s); b == true {
		return true, triggerBody{Type: PREFIX, Key: trig}
	}
	return false, triggerBody{Type: 0, Key: ""}
}

func hasTriggerWord(trg string) (exist bool, trigger string) {
	for _, key := range triggerWords {
		if strings.Contains(trg, key) {
			// exit = true
			return true, key
		}
	}
	return false, ""
}

func hasTriggerPrefix(trg string) (exist bool, trigger string) {
	for i, key := range triggerPrefixes {
		if strings.HasPrefix(trg, key) {
			return true, triggerPrefixes[i]
		}
	}
	return false, ""
}
