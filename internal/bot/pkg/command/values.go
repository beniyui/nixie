package command

var triggerWords = []string{"nixie", "ニキシ", "にきし", "ねえさん"}
var triggerPrefixes = []string{"n!"}


//command.Context
type Context struct {
	CmdIsDetected bool
	TriggerInfo triggerCont
	CmdInfo cmdCont
}

//trigger context
type triggerCont struct {
	Type   TriggerType
	String string
}
type TriggerType int
const (
	//0
	WORD TriggerType = iota + 1
	PREFIX
)

//cmd context
type cmdCont struct{
	Main cmdbody
	Options []cmdbody
}


// cmdbody
type cmdbody struct{
	ID id
	Arg []string
}

type id int
const(
	// = 0
	TEST id = iota + 1
	PING
	AI

)
