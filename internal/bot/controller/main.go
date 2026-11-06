package controller

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"nixie/internal/bot/config"
	"nixie/internal/ui"

	"github.com/briandowns/spinner"
	"github.com/bwmarrin/discordgo"
)

type Worker struct {
	BreakSwitch <-chan string
}

func NewWorker(jobChan <-chan string) *Worker {
	return &Worker{BreakSwitch: jobChan}
}

// スピナーの設定
var spi = spinner.New(spinner.CharSets[43], 100*time.Millisecond)

func Entry() {
	boot()
	welcomeMSG()
	beet()

	os.Exit(1)
}

func boot() {
	//env init
	if err := config.Init(); err != nil {
		msg := fmt.Sprintf("%+v", config.Values)
		msg = fmt.Sprintf("%s_values:%s", err, msg)
		ui.FatalExit(msg)
	}
	fmt.Println("env ok")

	//art
	if err := config.ArtInit(); err != nil {
		log.Println(err)
	}

	//make new session;the client
	token := config.Values.Token
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		msg := fmt.Sprintf("discordgo: new session error:%v", err)
		ui.FatalExit(msg)
		return
	}

	//discordgo identify intentsの設定はがちでよくわからん コピペ。
	dg.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentDirectMessages | discordgo.IntentMessageContent

	//add handlers
	dg.AddHandler(MessageCheck)
	dg.AddHandler(onInteractionCreate)

	//registerSlashCommands is still test. work in progress
	if err := registerCommands(dg); err != nil {
		log.Printf("コマンドの登録に失敗しました: %v", err)
	}

	//open
	spi.Prefix = "discordgo opening"
	spi.FinalMSG = "discordgo opened.\n"
	spi.Start()
	if err := dg.Open(); err != nil {
		msg := fmt.Sprintf("discordgo conection error\nたぶんネットが無いだけ:%v", err)
		ui.FatalExit(msg)
		return
	}
	defer dg.Close()
	spi.Stop()
}

func beet() {
	//プログラムの終了を検知するクライアント
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	//nixie down commandのブレーカー
	breaker := done()

	//ここでストップする
	select {
	case <-ctx.Done():
	case <-breaker:
		log.Println("nixieDOWN done")
	}
}

// display banner and wellcomeMSG
func welcomeMSG() {
	fmt.Printf("welcome to:\n%s\nめんどくさがりだけど優しいボット\n", config.Arts.TitleBanner)
	fmt.Printf("version: %s | exeDir: %s\n", config.Values.Version, config.Values.ExeDir)
	fmt.Println("now nixie UP. 以下ログ\n___\n")
}
