package ui

import (
	"os"

	"github.com/sqweek/dialog"
)

//log.Fatalf()の代替。エラー表示してから終了する
func FatalExit(input interface{}) {
	dialog.Message("%s", input).
		Title("致命的なえらー!").
		Error()

	// 2. ダイアログが閉じられたあとにプログラムを終了
	os.Exit(1)
}
