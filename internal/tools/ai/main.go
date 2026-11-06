package ai

import (
	"encoding/json"
	"fmt"

	"nixie/internal/tools/tools"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type ModelID string
type ModelExplanation string
type PromptName string

// モデル情報をまとめる構造体
type ModelInfo struct {
	ID          ModelID
	Explanation ModelExplanation
}

// パッケージ全体で共有するモデルリスト（一元管理）
var Models = []ModelInfo{
	{
		ID:          "gemma4:e2b",
		Explanation: "Gemma4 e2b | 考える。",
	},
	{
		ID:          "hf.co/DuoNeural/Gemma-4-E2B-Heretic-GGUF:Q4_K_M",
		Explanation: "Gemma4 e2b Heretic | 検閲なし。nixie#非対応。",
	},
}

type UseInfoPayload struct {
	Normal_model string
	Think_model  string
	Prompt_name  string
}

func Reply(s *discordgo.Session, m *discordgo.MessageCreate, typingStart func(), u UseInfoPayload) {
	//sync.Mutex prevent from memory race(=メモリ競合。複数の処理が同時に同じメモリ領域にアクセスするとクラッシュを起こす現象。)
	//discordgoはいつでも複数の処理が動作する可能性がある構造であり，かつ，今回はマップを操作するので必要である
	// たとえば，このmessageCheckが同時に5回くらい呼び出されていてもおかしくない
	//詳しくは「go言語 排他制御」で調べよう
	//mu.Lock
	//mu.Unlock

	typingStart()
	reply, repJson, err := reply(s, m, u)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "AIにアクセスできねえんだけど！どうしてくれんのこれ<@1409508397811892278>")
		msg := fmt.Sprintf("ollama err: %v", err)
		tools.LogToDiscord(s, msg)
		return
	}

	sec := int(time.Duration(repJson.Time).Seconds())

	bytes, _ := json.Marshal(repJson)
	jsonSt := string(bytes)
	tools.LogToDiscord(s,jsonSt)
	fmtReply := func(raw string) string {
		if raw == "" {
			return "<返答が空白です>"
		}
		return raw
	}

	finalMsg := fmt.Sprintf("%s\n-# %ds", fmtReply(reply), sec)

	//timeDist := int(math.Floor(respTime))
	//replyDist := fmt.Sprintf("%s\n<%ds>", reply, timeDist)
	trgMsg := &discordgo.MessageReference{
		MessageID: m.ID,
		ChannelID: m.ChannelID,
		GuildID:   m.GuildID,
	}
	s.ChannelMessageSendReply(m.ChannelID, finalMsg, trgMsg)
}

func reply(s *discordgo.Session, m *discordgo.MessageCreate, u UseInfoPayload) (string, *Response, error) {

	//add some prompt before  user prompt
	userPrompt := m.Content
	trgpath := fmt.Sprintf("./config/prompt/%s", u.Prompt_name)
	context, err := tools.GetTxt(trgpath)

	payloadPrompt := fmt.Sprintf("%s\n\nMy Message is '%s'", context, userPrompt)

	thinkBOOL := strings.HasPrefix(userPrompt, "nixie#")
	if thinkBOOL {
		s.MessageReactionAdd(m.ChannelID, m.ID, "think:1411655541490450532") //think
	}
	reqBody := returnReqBody(thinkBOOL, u, payloadPrompt)

	buf, _ := json.Marshal(reqBody)

	tools.LogToDiscord(s, string(buf))
	//s.ChannelMessageSend(m.ChannelID, string(buf))

	res, err := QueryOllama(reqBody)
	if err != nil {
		return "", nil, fmt.Errorf("Ollama query error: %s", err)
	}

	//s.ChannelMessageSend(m.ChannelID, string(bytes))

	return res.Response, res, nil
}

func returnReqBody(thinkBOOL bool, u UseInfoPayload, prompt string) Request {
	if thinkBOOL {
		reqBody := Request{
			Prompt: prompt,
			Stream: false,

			Model:   u.Think_model,
			Think:   true,
			Options: &Options{ReasoningEffort: "max"},
		}
		return reqBody
	} else {
		reqBody := Request{
			Prompt: prompt,
			Stream: false,

			Model: u.Normal_model,
			Think: false,
		}
		return reqBody
	}
}
