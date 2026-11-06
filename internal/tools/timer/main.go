package timer

import (
	"fmt"
	"log"
	"os"
	"time"

	tea "charm.land/bubbletea/v2"
)

// 毎秒呼び出されるメッセージ（時間経過の通知）
type tickMsg time.Time

// アプリケーションの状態（Model）
type model struct {
	duration int // タイマーの総時間（秒）
	seconds  int // 現在のカウントダウン（秒）
	quitting bool
	msg      string
}

// 1秒ごとにtickMsgを発行するコマンド
func tick() tea.Cmd {
	sendMsg := func(t time.Time) tea.Msg {
		return tickMsg(t)
	}
	return tea.Tick(time.Second, sendMsg)
}

// 初期化処理
func (m *model) Init() tea.Cmd {
	return tick()
}

// 状態更新（Update）
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tickMsg:
		m.seconds--
		if m.seconds <= 0 {
			m.quitting = true
			return m, tea.Quit
		}
		// 再度1秒タイマーをセット
		return m, tick()
	}

	return m, nil
}

// 画面表示（View）
func (m *model) View() tea.View {
	content := fmt.Sprintf(m.msg, m.seconds) + "\n"
	return tea.NewView(content)
}

// example:Seconds("wait%ds",10)
func FatalExit(msg string, inputTime int, inputErr error) {

	fmt.Printf("%s: %v\n", msg, inputErr)
	log.Printf("%v", inputErr)
	initialModel := model{
		msg:      "\n%d秒後に閉じます...\n",
		duration: inputTime,
		seconds:  inputTime,
	}

	p := tea.NewProgram(&initialModel)
	if _, err := p.Run(); err != nil {
		log.Printf("タイマーもエラーしました: %v\n", err)
		os.Exit(1)
	}

	//delay := time.Duration(input)
	//time.Sleep(delay * time.Second)
	os.Exit(1)

}
