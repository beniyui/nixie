package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

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