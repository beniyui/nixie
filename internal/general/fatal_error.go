package general

import (
	"fmt"
)

func fatalError(input any) {
	st := fmt.Sprintf("%v", input)

	fmt.Print("\n\n致命的えらー\n\n")
	fmt.Println("スクショしてとりあえずおれに見せたほうがいい。原因の例：\nネットがない、プログラムのミスなど")
	fmt.Printf("%s\n\nENTERなどを押して閉じる\n", st)
	fmt.Scanln()
}
