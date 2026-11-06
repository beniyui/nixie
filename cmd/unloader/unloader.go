package main

import (
	"fmt"
	"nixie/internal/controller"
	"nixie/internal/tools/timer"
	"os"
	"time"
)

func main() {
	err := controller.ModelUnloader()
	if err != nil {
		timer.FatalExit("アンロード失敗：", 5, err)
	}
	fmt.Println("成功。")
	time.Sleep(2 * time.Second)
	os.Exit(1)
}
