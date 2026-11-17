package general

import (
	"fmt"
	"nixie/internal/config"
	"nixie/internal/bot"

	"github.com/bwmarrin/discordgo"
)

func general() error {

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
	// dg.AddHandler(bot.RootHandler)

	if err := dg.Open(); err != nil {
		return fmt.Errorf("discordgo conection error\nたぶんネットが無いだけ:%v", err)
	}
	defer dg.Close()

	//welcomeMSG
	fmt.Println("成功！")
	fmt.Printf("welcome to:\n%s\nめんどくさがりだけど優しいボット\n", config.ArtsCont.TitleBanner)
	fmt.Printf("version: %s | exeDir: %s\n", config.EnvCont.Version, config.EnvCont.ExeDir)
	fmt.Printf("now nixie UP. logs:\n=================\n\n")

	//stay alive
	stay()

	//exit
	fmt.Println("bye")
	return nil
}
