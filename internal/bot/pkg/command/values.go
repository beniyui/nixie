package command

var triggerWords = []string{"nixie", "ニキシ", "にきし", "ねえさん"}
var triggerPrefixes = []string{"n!"}

// command.Context
type Context struct {
	CmdIsDetected bool
	Cmd cmdBody
	Options     []cmdBody
	TriggerInfo triggerCont
}

// trigger context
type triggerCont struct {
	Type TriggerType
	Key  string
}
type TriggerType int

const (
	//0
	WORD TriggerType = iota + 1
	PREFIX
)

// option

// cmdBody
type cmdBody struct {
	ID  id
	Arg string
}

type id int

const (
	// = 0
	TEST id = iota + 1
	PING
	AI
)
