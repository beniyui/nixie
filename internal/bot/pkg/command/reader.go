package command

import (
	"fmt"
	"strings"
)

func reader(rawStr string, trig triggerBody) (options []optionBody, err error) {
	var bufBody optionBody
	switch trig.Type {
	case WORD:
		bufBody.ID = AI
		bufBody.Arg = rawStr
		return []optionBody{bufBody}, nil
	case 0:
		return []optionBody{}, fmt.Errorf("something strange. triggerType is 0, means it has no trigger")
	}

	//it has Prefix...

	PrefixTrimedStr := strings.TrimPrefix(rawStr, trig.Key)
	spaceStr := strings.Split(PrefixTrimedStr, " ")
	FistString := strings.ToLower(spaceStr[0])
	//FirstString without prefix, lowered.

	hasOptions := strings.Contains(FistString, "--")
	if !hasOptions {
		//means no options
		bufBody.ID = AI
		bufBody.Arg = rawStr
		return []optionBody{bufBody}, nil
	}

	//it has also Options...

	optionsStrSlice := strings.Split(FistString, "--")
	opt, err := optionsBuilder(optionsStrSlice,PrefixTrimedStr)
	if err != nil {
		return []optionBody{}, err
	}
	return opt, nil
}

func optionsBuilder(optionStr []string, prefixTrimedInput string) ([]optionBody, error) {
	//they are options!

	var bufSlice []optionBody

	for _, value := range optionStr {
		var bufBody optionBody
		elems := strings.Split(value, "=")
		switch len(elems) {
		case 1:
			switch elems[0] {
			case "ping":
				bufBody.ID = PING
				return []optionBody{bufBody}, nil
			case "ai":
				bufBody.ID = AI
				bufBody.Arg = prefixTrimedInput
			case "think":
				bufBody.ID = THINK
			default:
				bufBody.ID = NOTFOUND
				
			}
			bufSlice = append(bufSlice, bufBody)
		case 2:
		//オプションの引数はまだ非対応。
		default:
			return []optionBody{}, fmt.Errorf("こちらのオプションは２個以上の引数を持っています。: %s", value)
		}
	}

	return bufSlice, nil
}

// func optionParse(elem []string) cmdBody {
// 	var bufBody cmdBody
// 	if len(elem) == 2 {
//
// 	}
// }
// // // // // //
