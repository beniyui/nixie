package general

import (
	"fmt"
)

func fatalError(input any) {
	st := fmt.Sprintf("%v", input)
	fmt.Printf("\n==================\n\n致命的えらー\n\n\n%s\n\nENTERなどを押して閉じる", st)
	fmt.Scanln()
}
