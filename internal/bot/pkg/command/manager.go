package command

func Manager(input string) (Context, error) {
	var Contextbuf Context

	has, trig := hasTrigger(input)
	if !has {
		Contextbuf.CmdIsDetected = false
		return Contextbuf,nil
	}
	Contextbuf.CmdIsDetected = true

	buf := cmdReader(trig)
	Contextbuf.CmdInfo = cmdCont{ID: 0, Arg: []BadExpr, Options: []id{}}

}
