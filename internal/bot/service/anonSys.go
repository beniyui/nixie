package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	AnonDmList = make(map[string]string)
	Mu         sync.Mutex
)

func AnonSys(s *discordgo.Session, m *discordgo.MessageCreate) {
	// when it has no trigger but it is in nixie's DM
	//匿名投稿システム Anon anonymous

	//sync.Mutex prevent from memory race(=メモリ競合。複数の処理が同時に同じメモリ領域にアクセスするとクラッシュを起こす現象。)
	//discordgoはいつでも複数の処理が動作する可能性がある構造であり，かつ，今回はマップを操作するので必要である
	// たとえば，このmessageCheckが同時に5回くらい呼び出されていてもおかしくない
	//詳しくは「go言語 排他制御」で調べよう
	//mu.Lock
	//mu.Unlock

	//AnonDmList = make(map[chanID]string)
	Mu.Lock()
	_, existInDMlist := AnonDmList[m.ChannelID]
	Mu.Unlock()

	if !existInDMlist { //新規の受付
		anonMsg := m.Content
		Mu.Lock()
		AnonDmList[m.ChannelID] = anonMsg
		Mu.Unlock()
		checkMsg := fmt.Sprintf("以下の内容で匿名投稿をする？(はい/いいえ)\n```%s```\nはい なら`y`か`はい`, いいえ なら`n`か`いいえ`と送信して（全角，大文字可）", anonMsg)
		s.ChannelMessageSend(m.ChannelID, checkMsg)
	} else {
		//existInDM listがtrueということは新規受付を済ませて，今yかnかの返答を待っているということだ
		//送ったあとリストから消すから
		switch m.Content {
		case "y", "ｙ", "Y", "はい":

			Mu.Lock()
			anonMsg := AnonDmList[m.ChannelID]
			delete(AnonDmList, m.ChannelID)
			Mu.Unlock()

			PostWebhook(os.Getenv("WEBHOOK_URL_1"), anonMsg)
			s.ChannelMessageSend(m.ChannelID, "投稿を送信しました")
		case "n", "ｎ", "N", "いいえ":
			Mu.Lock()
			delete(AnonDmList, m.ChannelID)
			Mu.Unlock()
			s.ChannelMessageSend(m.ChannelID, "投稿をやめました")
		default:
			s.ChannelMessageSend(m.ChannelID, "不明: はい なら`y`, いいえ なら`n`と送信して（大文字，全角可）")
		}
	}
}


//以下webhook関連。

// Embed 構造体定義
type Embed struct {
	Description string `json:"description,omitempty"`
	Timestamp   string `json:"timestamp,omitempty"`
}

// WebhookPayload Discord Webhookのリクエスト構造体
type WebhookPayload struct {
	Username  string  `json:"username,omitempty"`
	AvatarURL string  `json:"avatar_url,omitempty"`
	Embeds    []Embed `json:"embeds"`
}

// PostWebhook 指定したURLへDiscordメッセージを送信します
func PostWebhook(targetURL string, msg string) error {
	payload := WebhookPayload{
		Embeds: []Embed{
			{
				Description: msg,
				Timestamp:   time.Now().UTC().Format(time.RFC3339),
			},
		},
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	// タイムアウト付きのHTTPクライアントを作成
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, targetURL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	// レスポンスボディの破棄とクローズを確実に行う
	defer func() {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	// HTTPステータスコードのチェック（200 OK または 204 No Content）
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d (%s)", resp.StatusCode, resp.Status)
	}

	return nil
}
