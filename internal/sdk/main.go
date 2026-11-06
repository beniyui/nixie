package sdk_test

import (
	"fmt"
	"nixie/internal/ui"
	"strings"

	"github.com/eslider/go-ollama"
)

func Query() {
	// 1. クライアントの初期化（ローカルOllamaの場合はToken不要）
	client := ollama.NewOpenWebUiClient(&ollama.DSN{
		URL: "http://localhost:11434/api/generate",
	})

	var fullResponse strings.Builder

	// 2. リクエストの送信とストリーミング受信
	err := client.Query(ollama.Request{
		Model:  "gemma4:e2b", // 手元の環境にインストールされているモデル名に変更してください
		Prompt: "3行で世界の美しさについて教えて。",
		OnJson: func(res ollama.Response) error {
			if res.Response != nil {
				// トークンが生成されるたびにリアルタイムで標準出力に出力
				fmt.Print(*res.Response)
				// 全体のレスポンスを組み立てるためにバッファに蓄積
				fullResponse.WriteString(*res.Response)
			}
			return nil
		},
	})

	if err != nil {
		msg := fmt.Sprintf("Error querying Ollama:%v",err)
		ui.FatalExit(msg)
	}

	fmt.Println("\n--- 完了 ---")
}
