package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Embedの構造体定義


type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

type EmbedFooter struct {
	Text string `json:"text"`
}

type Embed struct {
	Title       string         `json:"title,omitempty"`
	Description string         `json:"description,omitempty"`
	URL         string         `json:"url,omitempty"`
	Color       int            `json:"color,omitempty"` // example: 0x00FF00 -> 65280
	Fields      []EmbedField   `json:"fields,omitempty"`
	Footer      *EmbedFooter   `json:"footer,omitempty"`
	Timestamp   string         `json:"timestamp,omitempty"`
}

// Discord Webhookのリクエスト構造体
type WebhookPayload struct {
	Username  string  `json:"username,omitempty"`
	AvatarURL string  `json:"avatar_url,omitempty"`
	Embeds    []Embed `json:"embeds"`
}

func PostWebhook(targetURL string) {
	webhookURL := targetURL

	// Embedデータの作成
	payload := WebhookPayload{
		Username: "Go Bot",
		Embeds: []Embed{
			{
				Title:       "test title",
				Description: "test description",
				URL:         "https://golang.org",
				Color:       0x3498db, // 青色
				Fields: []EmbedField{
					{
						Name:   "field1",
						Value:  "Horizontal",
						Inline: true,
					},
					{
						Name:   "field2",
						Value:  "Horizontal",
						Inline: true,
					},
				},
				Footer: &EmbedFooter{
					Text: "footer ready",
				},
				Timestamp: time.Now().Format(time.RFC3339),
			},
		},
	}

	// JSONに変換
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("JSON convert error:", err)
		return
	}

	// HTTP POSTリクエストの作成と送信
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		fmt.Println("post request error", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusOK {
		fmt.Println("succes!")
	} else {
		fmt.Printf("failed (Status: %s)\n", resp.Status)
	}
}
