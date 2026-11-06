package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	OllamaURL = "http://localhost:11434/api/generate"
)

type Response struct {
	Response string `json:"response"`
	Time     int64  `json:"total_duration"`
	Error    string `json:"error,omitempty"`
	DoneReason string `json:"done_reason"`
}

type Options struct {
	ReasoningEffort string `json:"reasoning_effort,omitempty"`
}

type Request struct {
	Model     string  `json:"model"`
	Prompt    string  `json:"prompt"` // omit empty= 空で送信されたときは省略する
	Stream    bool    `json:"stream"`
	Think     bool    `json:"think"`
	KeepAlive any     `json:"keep_alive,omitempty"` // -1で常駐、"10m"や"1h"などで時間指定
	Options   *Options `json:"options,omitempty"`
	Format string `json:"format,omitempty"`
}

func QueryOllama(body Request) (*Response, error) {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// 1. NewRequest を使ってリクエストオブジェクトを作成
	req, err := http.NewRequest("POST", OllamaURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	// 2. 明示的に Content-Type ヘッダーを設定
	req.Header.Set("Content-Type", "application/json")

	// 3. Client 経由で送信
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to request: error code: %d", resp.StatusCode)
	}

	replyBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var r *Response
	if err := json.Unmarshal(replyBody, &r); err != nil {
		return nil, err
	}
	return r, nil
}

func PreloadModel(model ModelID) error {
	reqBody := Request{
		Model:     string(model),
		Prompt:    "",
		KeepAlive: -1,
	}
	_, err := QueryOllama(reqBody)
	return err
}

func UnloadModel(model ModelID) error {
	reqBody := Request{
		Model:     string(model),
		KeepAlive: 0,
	}
	_, err := QueryOllama(reqBody)
	return err
}
