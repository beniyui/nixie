package general

import (
	"fmt"
	"os"
)

func fatalError(input any) {
	st := fmt.Sprintf("%v", input)
	fmt.Printf("\n==================\n\n致命的えらー\nMSG:\n\n%s\n\nENTERなどを押して閉じる", st)
	fmt.Scanln()
	os.Exit(1)
}
