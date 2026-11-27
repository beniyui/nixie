package command

import "strings"

//call this func from parser
func hasTrigger(inputString string) (hasTrigger bool, trig triggerCont) {
	s := strings.ToLower(inputString)
	if b, trig := hasTriggerWord(s); b == true {
		return true, triggerCont{Type: WORD, Key: trig}
	}
	if b, trig := hasTriggerPrefix(s); b == true {
		return true, triggerCont{Type: PREFIX, Key: trig}
	}
	return false, triggerCont{Type: 0, Key: ""}
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
