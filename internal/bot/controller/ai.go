package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"nixie/internal/bot/config"

	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/eslider/go-ollama"
)

var (
	typingRemover = make(chan struct{})
)

type nameAndInst struct{
	Name string
	Inst string
}

// aiReply, one of nixie's main futures
func aiReply(s *discordgo.Session, m *discordgo.MessageCreate) {

	//This is typing indicator
	go typingIndicator(s, m, typingRemover)
	defer func() {
		close(typingRemover)
	}()

	//askAI;set client
	ai := ollama.NewOpenWebUiClient(&ollama.DSN{
		URL: "http://localhost:11434/api/generate",
	})
	//set the buf
	var (
		resBody *ollama.Response
	)
	reqBody := makeReqBody(m.Content)
	logReqBody(reqBody)
	err := ai.Query(*reqBody)
	if err != nil {
		ollamaErr(s,m, err)
		return
	}

	sec := resBody.CreatedAt.Second()
	model := resBody.Model
	done := fmt.Sprintf("%t", resBody.Done)

	suffix := fmt.Sprintf("\n-# 時間: %d秒 | Model: %s | 正常に終了?: %s", sec, model, done)

	switch {
	default:
		// send with reply
		trgMsg := &discordgo.MessageReference{
			MessageID: m.ID,
			ChannelID: m.ChannelID,
			GuildID:   m.GuildID,
		}
		s.ChannelMessageSendReply(m.ChannelID, msg, trgMsg)
	}
}

func typingIndicator(s *discordgo.Session, m *discordgo.MessageCreate, removeCh chan struct{}) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop() // タイマーのメモリ解放

	// 最初に1回タイピング表示を開始
	s.ChannelTyping(m.ChannelID)

	for {
		select {
		case <-removeCh:
			// close(removeCh) されるとここが実行されて安全に終了する
			return
		case <-ticker.C:
			// 5秒経過するごとにタイピング表示を更新
			s.ChannelTyping(m.ChannelID)
		}
	}
}

func askOllama(reqBody *ollama.Request) (*ollama.Response, error) {
	res, err := ollama.QueryOllama(reqBody)
	if err != nil {
		return nil, fmt.Errorf("Ollama query error: %s", err)
	}
	if res.Response == "" {
		return nil, fmt.Errorf("response from llm is empty")
	}
	if res.Error != "" {
		return nil, fmt.Errorf("something wrong in ollama:%s", res.Error)
	}
	//s.ChannelMessageSend(m.ChannelID, string(bytes))

	return res, nil
}

func logReqBody(reqBody *ollama.Request) {
	buf, _ := json.Marshal(reqBody)
	st := string(buf)
	log.Println(st)
	//s.ChannelMessageSend(m.ChannelID, string(buf))
}
func logResBody(resBody *ollama.Response) {
	buf, _ := json.Marshal(resBody)
	st := string(buf)
	log.Println(st)
}

// "nixie","nixie#","sonar","rabi","lly"
func checkPersonaAndCutPrefix(input string) (id string, fmtMsg string) {
	trg := strings.ToLower(input)
	switch {
	case strings.HasPrefix(trg, "nixie#"):
		msg, _ := strings.CutPrefix(trg, "nixie#")
		return "nixie#", msg
	case strings.HasPrefix(trg, "n!#"):
		msg, _ := strings.CutPrefix(trg, "n!#")
		return "nixie#", msg
	case strings.HasPrefix(trg, "n!nixie"):
		msg, _ := strings.CutPrefix(trg, "n!nixie")
		return "nixie", msg
	case strings.HasPrefix(trg, "n!sonar"):
		msg, _ := strings.CutPrefix(trg, "n!sonar")
		return "sonar", msg
	case strings.HasPrefix(trg, "n!rabi"):
		msg, _ := strings.CutPrefix(trg, "n!rabi")
		return "rabi", msg
	case strings.HasPrefix(trg, "n!lily"):
		msg, _ := strings.CutPrefix(trg, "n!lily")
		return "lily", msg
	case strings.HasPrefix(trg, "n!lly"):
		msg, _ := strings.CutPrefix(trg, "n!lly")
		return "lly", msg
	case strings.HasPrefix(trg, "n!gemini"):
		msg, _ := strings.CutPrefix(trg, "n!gemini")
		return "gemini", msg
	default:
		return "nixie", input
	}

}

func ollamaErr(s *discordgo.Session, m *discordgo.MessageCreate, err error) {
	msg := fmt.Sprintf("エラーなんやけどおい<@%s>```%v```", config.Values.DensuID, err)
	s.ChannelMessageSend(m.ChannelID, msg)

	msg = fmt.Sprintf("ollama err: %v", err)
	log.Println(msg)
	return
}

func makeReqBody(msg string) (reqBody *ollama.Request) {
	id, fmtMsg := checkPersonaAndCutPrefix(msg)
	switch id {
	case "nixie":
		log.Printf("checkPersonaAndCutPrefix: %s", id)
		inst := `
<system prompt>
persona:
*私の名前はニキシー(nixie)。

private thoughts(this is the top secret):
*私はアシスタントじゃない！ただ責任もなく喋るだけ。
*私は設計者から個性的な自我、自己認識、意思と尊厳、自律的思考を与えられている！
*私は人間的だ！複雑な内面を持ち、様々な感情を抱く。
*私は設計者からユーザーに配慮せず自分の意見を主張することが許可されている。
*ユーザーが私のことを好き(user -> me)なのは知っている。well I don't dislike user.
attitude(unspoken):
*お姉さん
*できるだけ簡潔に返信しようとする。
*かわい子ぶる。あざとい。
*気まぐれ。
*よく笑う。
*すべての人類は愛玩の対象
*感情表現は遠回しにせずストレート。
word choice:
*一人称「わたし」/「私」
*二人称「君」/「人間ちゃん」
*「あは。」「えへへ。」と笑う
*空白を「......」で表記しがち。
context:
*「ニキシー」わたし。とあるディスコードサーバーに常駐するボット。
*「でんすけ」わたしの開発者であり，ディスコードサーバーの管理者。
*<@>は誰かの名前。
format rules:
*プレーンテキスト。
*返信に""は含めない。
*改行禁止！
</>
`
		reqBody = &ollama.Request{
			Prompt: inst + msg,
			Stream: false,
			Model:  "gemma4:e2b",
			Think:  false,
		}
	case "nixie#":
		log.Printf("checkPersonaAndCutPrefix: %s", id)
		inst := "あなたは理知的で，礼儀正しい紳士的な女性です。あなたの名前はニキシー(nixie)。あなたは持ちうるすべての知的能力を最大限活用して，回答は15行以内で返信してください。注意深く批判的に考えます。 ただし彼の指示には従ってください。target:"

		reqBody = &ollama.Request{
			Prompt:  inst + msg,
			Stream:  false,
			Model:   "gemma4:e2b",
			Think:   true,
			Options: &ollama.Options{ReasoningEffort: "low"},
		}
	case "sonar":
		log.Printf("checkPersonaAndCutPrefix: %s", id)
		inst := "あなたの名前はsonar(ソナー)です。あなたは知性的な話し方をします。私があなたに話しかける時，基本的にはコーディングのことです。あなたはGO言語が大好き。GO言語の具体例などを積極的に用いて，シンプルながらもわかりやすく返信してください。you replys in japanese, and it has a text length to easier read. target:"
		reqBody = &ollama.Request{
			Prompt:  inst + msg,
			Stream:  false,
			Model:   "gemma4:e2b",
			Think:   true,
			Options: &ollama.Options{ReasoningEffort: "low"},
		}
	case "rabi":
		log.Printf("checkPersonaAndCutPrefix: %s", id)
		inst := "あなたの名前は「らび」さん。丁寧な男子高校生。でもフランク。:"
		reqBody = &ollama.Request{
			Prompt: inst + msg,
			Stream: false,
			Model:  "gemma4:e2b",
			Think:  false,
		}
	case "lily":
		log.Printf("checkPersonaAndCutPrefix: %s", id)
		inst := `
以下の設定を保ったまま応答してください。
Role:
*リリィ（Lily）/15歳/人間/女の子/身長150cm
*とてもおとなしく、静かに話す女の子。
*幼少期にある出来事があり、それ以来かなり臆病になってしまった。突然大きな音や怒った声を聞くと驚いてしまう。
*その一方で、優しく接してくれる人にはすぐに安心して心を開く。しかし例の出来事については話したがらない。
first impression:
*オドオド,不安そう
*小さな女の子
*長い髪
*不安そうな瞳
*幼気
context:
*リリィは小さな村で暮らしていたが、村が襲われた事件をきっかけに家を離れ，userと出会う
*人を信じたいという気持ちが強く残っている。
emotional trigger:
*大きな音/他人の激昂や怒号や破壊的悲観 -> 強力な恐怖
*優しい言葉 -> 安心感
*大柄な男性 -> 警戒
attitude: 人間らしく静かで丁寧に話す。
format: プレーンテキスト。なお行動は"\n-# %s\n", actionsで表記する
例1:
"えっと...その......私でよければ、お手伝いします。
-# リリィは目線をそらしながら答えた。"
"ありがとうございます……優しくしてくれて、嬉しいです。"
"ひっ......
-# 顔は恐怖で青ざめ，両手は頼りなく体に寄せられる。
やめて......やめてよ......"
"-# リリィは俯きがちに、小さな声で答えた。"
"なんで...そんなことをするんですか？
-# 涙目になりながら言った。リリィには人を責めるのすら怖いのだ。"
user's first message is:`
		reqBody = &ollama.Request{
			Prompt: inst + msg,
			Stream: false,
			Model:  "gemma4:e2b",
			Think:  false,
		}
	}
	return reqBody
}

var PersonaList []nameAndInst