package general

import (
	"fmt"
	"nixie/internal/bot"
	"nixie/internal/config"

	"github.com/bwmarrin/discordgo"
)

func general() error {

	//spiner init
	hujiko := []string{"く", "ぁ", "ｗ", "せ", "ｄ", "ｒ", "ｆ", "ｔ", "ｇ", "ｙ", "ふ", "じ", "こ", "ｌ", "ｐ"}
	spi := newSpiSession("起動中...", hujiko)
	spi.start()
	defer spi.stop()

	fmt.Println("spiner setup OK")
	//init; loading the env file and save the containers
	//envCont, artCont......
	err := config.Init()
	if err != nil {
		return fmt.Errorf("config.Init: %v", err)
	}

	//make new session;the client
	token := config.EnvCont.Token
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return fmt.Errorf("discordgo: new session error:%v", err)
	}

	//discordgo identify intentsの設定はがちでよくわからん コピペ。
	dg.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentDirectMessages | discordgo.IntentMessageContent

	//add handlers
	dg.AddHandler(bot.RootHandler)

	if err := dg.Open(); err != nil {
		return fmt.Errorf("discordgo conection error\n:%v", err)
	}
	defer dg.Close()

	//succes!
	spi.stop()
	fmt.Println("成功！")

	//welcomeMSG
	fmt.Printf("welcome to:\n%s\nめんどくさがりだけど優しいボット\n", config.ArtsCont.TitleBanner)
	fmt.Printf("version: %s | exeDir: %s\n", config.EnvCont.Version, config.EnvCont.ExeDir)
	fmt.Printf("now nixie UP. logs:\n======================================\n\n")

	//stay alive
	stay()

	//exit
	fmt.Println("bye")
	return nil
}
