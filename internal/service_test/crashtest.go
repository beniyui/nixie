package service_test

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

type crashLog struct {
	LastRun  time.Time `json:"last_run"`
	LastDate string    `json:"last_date"`
}

func CrashTest(s *discordgo.Session, m *discordgo.MessageCreate) {

	//idList := MemberIDlist(s, m)
	//ix := randomIndex(idList)
	//trgID := idList[ix]

	//s.ChannelMessageSend(m.ChannelID, trgID)

}

func save() string {
	t := time.Now()

	year, month, day := t.Date()
	date := fmt.Sprintf("%d-%d-%d", year, month, day)

	c := crashLog{
		LastRun:  t,
		LastDate: date,
	}

	// 1. 構造体をJSON形式（見やすく整形）のバイト配列に変換
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		log.Println("Marshal Error:", err)
		return ""
	}

	// 2. ファイルに書き出し (ファイル名, データ, 権限)
	// 0644 は「所有者は読み書き可能、その他は読み込みのみ」のパーミッション
	err = os.WriteFile("./config/save_file/crash_log.json", data, 0644)
	if err != nil {
		log.Println("WriteFile Error:", err)
		return ""
	}
	return string(data)
}

func randomIndex(list []string) int {
	var selector int = rand.N(len(list))
	return selector
}

func GetHumanMembers(s *discordgo.Session, guildID string) ([]*discordgo.Member, error) {
	var humanMembers []*discordgo.Member
	var lastID string

	for {
		// 1回につき最大1000件取得
		members, err := s.GuildMembers(guildID, lastID, 1000)
		if err != nil {
			return nil, err
		}

		if len(members) == 0 {
			break
		}

		for _, m := range members {
			// Botフラグがfalseのユーザーのみ抽出
			if m.User != nil && !m.User.Bot {
				humanMembers = append(humanMembers, m)
			}
			lastID = m.User.ID
		}

		// 取得したメンバー数が1000未満なら全件取得完了
		if len(members) < 1000 {
			break
		}
	}

	return humanMembers, nil
}

func MemberIDlist(s *discordgo.Session, m *discordgo.MessageCreate) []string {
	list, err := GetHumanMembers(s, m.GuildID)
	if err != nil {
		log.Printf("list err:%v\n", err)
		s.ChannelMessageSend(m.ChannelID, "err")
	}
	lines := make([]string, 0, len(list))

	for _, m2 := range list {
		newline := fmt.Sprintf("%s", m2.User.ID)

		lines = append(lines, newline)
	}
	return lines
}

func SendDM(s *discordgo.Session, m *discordgo.MessageCreate, trgID string, content string) error {
	dmChannel, err := s.UserChannelCreate(os.Getenv("DENSU_ID"))
	if err != nil {
		// ユーザーがDMを閉じていたり、ブロックされている場合はエラーになります
		log.Println("DMチャンネルの作成に失敗しました:", err)
		return err
	}
	_, err = s.ChannelMessageSend(dmChannel.ID, content)
	if err != nil {
		log.Println("DMの送信に失敗しました:", err)
		return err
	}
	return nil
}

func getUserRoles(s *discordgo.Session, guildID, userID string) ([]string, error) {
	// キャッシュ（State）からメンバー情報を取得を試みる
	member, err := s.State.Member(guildID, userID)
	if err != nil {
		// キャッシュにない場合はAPIリクエストを送信して取得
		member, err = s.GuildMember(guildID, userID)
		if err != nil {
			return nil, fmt.Errorf("メンバー情報の取得に失敗しました: %w", err)
		}
	}

	// member.Roles にロールIDの配列([]string)が含まれています
	return member.Roles, nil
}
