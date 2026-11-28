package command

func Listener(rawStr string) (Context, error) {
	var Contextbuf Context

	has, trig := hasTrigger(rawStr)
	if !has {
		Contextbuf.CmdIsDetected = false
		return Contextbuf, nil
	}

	buf, err := reader(rawStr, trig)
	if err != nil {
		return Context{}, err
	}

	Contextbuf = Context{
		CmdIsDetected: true,
		Options:       buf,
		TriggerInfo:   trig,
	}

	return Contextbuf, nil
}
