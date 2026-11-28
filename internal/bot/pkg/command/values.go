package command

var triggerWords = []string{"nixie", "ニキシ", "にきし", "ねえさん"}
var triggerPrefixes = []string{"n!"}

// command.Context
type Context struct {
    CmdIsDetected bool
	Options     []*optionBody
	TriggerInfo triggerBody
}

// option
type optionBody struct {
	ID  optionID
	Arg string
}

type optionID int

const(
    AI optionID = iota
	PING 
	THINK
	NOTFOUND
)

// trigger
type triggerBody struct {
	Type TriggerType
	Key  string
}

type TriggerType int

const (
	//0
	WORD TriggerType = iota + 1
	PREFIX
)
