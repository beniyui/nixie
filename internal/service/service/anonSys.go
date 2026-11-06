package service

import (
	"fmt"
	"nixie/internal/tools/tools"
	"os"
	"sync"

	"github.com/bwmarrin/discordgo"
)

var (
	AnonDmList = make(map[string]string)
	Mu         sync.Mutex
)

func AnonSys(s *discordgo.Session, m *discordgo.MessageCreate) {
	// means it has no trigger but it is in nixie's DM
	//匿名投稿システム Anon anonymous

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

			tools.PostWebhook(os.Getenv("WEBHOOK_URL_1"), anonMsg)
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
