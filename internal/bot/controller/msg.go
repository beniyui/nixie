package controller

import (
	"bytes"
	"fmt"
	"math/rand/v2"
	"nixie/internal/bot/config"
	"nixie/internal/bot/service"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var breaker = make(chan struct{})

func done() chan struct{} {
	return breaker
}

// この関数はメッセージが送信されるたびに実行される。
func MessageCheck(s *discordgo.Session, m *discordgo.MessageCreate) {
	if eventListen(s, m) {
		return
	}

	//test commands
	if m.Author.ID == config.Values.DensuID {
		switch m.Content {
		case "greetings":
			s.ChannelMessageSend(m.ChannelID, "う！")
		case "list":
			list := memberList(s, m)
			sendTextFile(s, m.ChannelID, list)
		case "list id":
			list := service.MemberIDlist(s, m)
			sendTextFile(s, m.ChannelID, list)

		case "dmtest":
			err := service.Dmtest(s, m)
			if err != nil {
				msg := fmt.Sprintf("%v", err)
				s.ChannelMessageSend(m.ChannelID, msg)
			}
		case "roletest":
			s.ChannelMessageSend(m.ChannelID, "roletest")
			service.Rolelisttest(s, m)
		}
	}

	switch {
	case pingCheck(m.Content):
		service.Ping(s, m)
	case triggerExist(m.Content):
		if strings.HasPrefix(m.Content, "nixie down") {
			service.Down(s, m, breaker)
			return
		}
		aiReply(s, m)

	default:
		if m.GuildID == "" {
			service.AnonSys(s, m)
		} else {
			replyRandom(s, m)
		}
	}
}

func randomPick(list []string) string {
	var selector int = rand.N(len(list))
	result := list[selector]
	return result
}

func triggerExist(target string) bool { //check trigger exist in input
	trg := strings.ToLower(target)
	switch {
	case strings.HasPrefix(trg, "n!"),
		strings.Contains(trg, "nixie"),
		strings.Contains(trg, "ニキシー"),
		strings.Contains(trg, "にきし"):
		return true
	}
	return false
}

func pingCheck(target string) bool {
	trg := strings.ToLower(target)
	switch trg {
	case "ping", "nix", "n!":
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
	list, err := service.GetHumanMembers(s, m.GuildID)
	if err != nil {
		msg := fmt.Sprintf("list err:%s", err)
		s.ChannelMessageSend(m.ChannelID, msg)
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
	case strings.Contains(m.Content, "うんこ"):
		i := []string{"うう", "こ", "こ。", "ぶっちっぱ！"}
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

func eventListen(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	switch {
	case s.State.User.ID == m.Author.ID:
		return false
	case m.Author.Bot == true:
		return false
	default:
		return true
	}
}
