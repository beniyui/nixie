package ollama

import(
	"encoding/json"
	"net/http"
	"bytes"
	"io"
	
)

const (
	TargetModel  = "gemma4:e2b"
	OllamaURL    = "http://localhost:11434/api/generate"
)

type response struct {
	Response string `json:"response"`
}
type options struct {
    ReasoningEffort string `json:"reasoning_effort"`
}
type request struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
	Think bool   `json:"think"`
	Options options `json:"options"`
}

func QueryOllama(prompt string) (string, error) {
	Body := request{
		Model: TargetModel,
		Prompt: prompt,
		Stream: false,
		Think: false,
		Options: options{ ReasoningEffort: "low",},
	}

	jsonData, err := json.Marshal(Body)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(OllamaURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var ollamaResp response
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return "", err
	}

	return ollamaResp.Response, nil
}
