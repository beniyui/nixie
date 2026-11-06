package controller

import (
	"bytes"
	"fmt"
	"os"

	"math/rand/v2"
	"nixie/internal/tools/ai"
	"nixie/internal/service/service"
	"nixie/internal/service_test"
	"nixie/internal/tools/tools"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var breaker = make(chan struct{})

func Done() chan struct{} {
	return breaker
}

// この関数はメッセージが送信されるたびに実行される。
func MessageCheck(s *discordgo.Session, m *discordgo.MessageCreate, u ai.UseInfoPayload) {
	if s.State.User.ID == m.Author.ID { //nixieの投稿
		return
	}

	//This is typing indicator. after called, nixie type forever. send typingRemover true to stop.
	typingRemover := make(chan struct{})
	typingStart := func() { go typingIndicator(s, m, typingRemover) }

	//after all, remove TypingIndicator
	defer func() {
		close(typingRemover)
	}()

	//test commands
	if m.Author.ID == os.Getenv("DENSU_ID") {
		switch m.Content {
		case "list":
			list := memberList(s, m)
			sendTextFile(s, m.ChannelID, list)
		case "list id":
			list := service_test.MemberIDlist(s, m)
			sendTextFile(s, m.ChannelID, list)

		case "dmtest":
			err := service_test.Dmtest(s, m)
			if err != nil {
				msg := fmt.Sprintf("%v", err)
				s.ChannelMessageSend(m.ChannelID, msg)
			}
		case "roletest":
			s.ChannelMessageSend(m.ChannelID, "roletest")
			service_test.Rolelisttest(s, m)
		}
	}

	switch {
	case triggerExist(m.Content):
		if strings.HasPrefix(m.Content, "nixie down") {
			service.Down(s, m, breaker)
			return
		}
		if afterTriggerEmpty(m.Content) {
			service.Ping(s, m)
		} else {
			ai.Reply(s, m, typingStart, u)
		}
	case !triggerExist(m.Content):
		if m.GuildID != "" {
			replyRandom(s, m)
		} else {
			service.AnonSys(s, m)
		}
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

func randomPick(list []string) string {
	var selector int = rand.N(len(list))
	result := list[selector]
	return result
}

func triggerExist(target string) bool { //check trigger exist in input
	if strings.HasPrefix(target, "nixie") {
		return true
	}
	if strings.HasPrefix(target, "Nixie") {
		return true
	}
	if strings.HasPrefix(target, "NIXIE") {
		return true
	}
	if strings.HasPrefix(target, "ニキシー") {
		return true
	}
	if strings.HasPrefix(target, "nix") {
		return true
	}
	return false
}

func afterTriggerEmpty(target string) bool {
	// check only trigger is in input
	if target == "nixie" {
		return true
	}
	if target == "Nixie" {
		return true
	}
	if target == "NIXIE" {
		return true
	}
	if target == "ニキシー" {
		return true
	}
	if target == "nix" {
		return true
	}
	return false
}

func attachExist(m *discordgo.MessageCreate) bool {
	if len(m.Attachments) > 0 {
		return true
	}
	return false
}

func memberList(s *discordgo.Session, m *discordgo.MessageCreate) []string {
	list, err := service_test.GetHumanMembers(s, m.GuildID)
	if err != nil {
		msg := fmt.Sprintf("list err:%s", err)
		tools.LogToDiscord(s, msg)
		s.ChannelMessageSend(m.ChannelID, "err")
	}
	lines := make([]string, 0, len(list))

	for _, m2 := range list {
		newline := fmt.Sprintf("%s | %s", m2.DisplayName(), m2.User.ID)

		lines = append(lines, newline)
	}
	return lines
}

func sendTextFile(s *discordgo.Session, channelID string, list []string) error {
	// リストを改行で結合
	result := strings.Join(list, "\n")

	// 文字列をio.Reader（バッファ）に変換
	fileReader := bytes.NewBufferString(result)

	// Discordにファイルをアップロードして送信
	// 第2引数: ファイル名, 第3引数: io.Reader
	_, err := s.ChannelFileSend(channelID, "excluded_bots.txt", fileReader)
	if err != nil {
		return fmt.Errorf("ファイルの送信に失敗しました: %w", err)
	}

	return nil
}

func replyRandom(s *discordgo.Session, m *discordgo.MessageCreate) {
	//DM以外の場所（サーバー）で，トリガーがない
	switch {
	case strings.Contains(m.Content, "でんすけ"):
		s.ChannelMessageSend(m.ChannelID, "<@1409508397811892278>")
	case m.Content == "う", strings.Contains(m.Content, "うう"), strings.Contains(m.Content, "うー"):
		i := []string{"うう", "ううう", "う", "ううあああええう"}
		s.ChannelMessageSend(m.ChannelID, randomPick(i))
	case strings.Contains(m.Content, "みろ"), strings.Contains(m.Content, "見ろ"):
		i := []string{"見とけよ見とけよ～", "見ろよ見ろよ",
			"ホラホラホラホラ", "ホラ，見ろよ"}
		s.ChannelMessageSend(m.ChannelID, randomPick(i))
	case strings.Contains(m.Content, "そうだよ"), strings.HasSuffix(m.Content, "ゾ"):
		i := []string{"当たり前だよなあ",
			"おっそうだな（適当）", "そうだよ（便乗）"}
		s.ChannelMessageSend(m.ChannelID, randomPick(i))
	case strings.Contains(m.Content, "おお"):
		i := []string{"おお", "おお。", "おおじゃないが", "おおかな",
			"これはおお", "お？", "これはおおか", "これはおおだな",
			"おおか？", "o", "oo", "Ooh", "oh", "おっほおﾞおﾞお～～ッ♥",
			"おっおっおっおっおっ"}
		s.ChannelMessageSend(m.ChannelID, randomPick(i))
	case strings.Contains(m.Content, "ちんちん"):
		i := []string{"じゅるじゅるぐっぽぐっぽ", "うん，おいしい！", ""}
		s.ChannelMessageSend(m.ChannelID, randomPick(i))
	case strings.Contains(m.Content, "行き"), strings.Contains(m.Content, "行く"), strings.Contains(m.Content, "いけ"):
		i := []string{"イキスギ？", "いくいくいくいく"}
		s.ChannelMessageSend(m.ChannelID, randomPick(i))
	case strings.Contains(m.Content, "スギ"):
		s.ChannelMessageSend(m.ChannelID, "お前やりますねえスギ！")
	case strings.HasSuffix(m.Content, "やん"):
		s.ChannelMessageSend(m.ChannelID, "どうしてくれんのこれ")
	}
}
