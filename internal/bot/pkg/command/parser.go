package command

import "strings"

func cmdReader(rawStr string, trig triggerCont) (cmdCont, error) {
	var bufcmdCont cmdCont
	switch trig.Type {
	case WORD:
		bufcmdCont.Main.ID = AI
		bufcmdCont.Main.Arg = []string{rawStr}
		return bufcmdCont, nil
	case PREFIX:
		buf := strings.TrimPrefix(rawStr, trig.String)
		if buf == "" ||strings.HasPrefix(buf,"ping"){
			bufcmdCont.Main.ID = PING
			return bufcmdCont, nil
		}
		if strings.Contains(buf," --")||strings.Contains(buf," -")||strings.HasPrefix(buf,"-"){
			h
		}

	}

}
