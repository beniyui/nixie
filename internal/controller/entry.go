package controller

import (
	"context"
	"fmt"
	"log"
	"nixie/internal/tools/ai"
	"nixie/internal/tools/timer"
	"nixie/internal/tools/tools"
	"nixie/internal/tools/wizard"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/briandowns/spinner"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	//"time"
)

type Worker struct {
	BreakSwitch <-chan string // 受信専用チャンネルとして定義するのが安全
}

func NewWorker(jobChan <-chan string) *Worker {
	return &Worker{BreakSwitch: jobChan}
}

var s = makeSpinner()

func GeneralController() {
	//スピナーの設定

	//.envをロード
	if err := godotenv.Load("./.env"); err != nil {
		timer.FatalExit("envの読み込みに失敗", 5, err)
	}

	log.Println("welcome to nixie")
	fmt.Println("welcome to:")

	fmt.Println(
		tools.GetArt("title", os.Getenv("ART_VERSION")),
	)
	version := os.Getenv("VERSION")
	fmt.Println("めんどくさがりだけど優しいボット")
	fmt.Printf("version: %s\n", version)

	fmt.Println("-")

	//the banner is displayed

	//起動時設定ウィザード
	norm, thi, trgPrompt := wizard.Setup()
	prName := string(trgPrompt)
	normName := string(norm)
	thinkName := string(thi)
	//保存
	payload := ai.UseInfoPayload{
		Normal_model: normName,
		Think_model:  thinkName,
		Prompt_name:  prName,
	}

	//事前読み込みウィザード
	if err := wizard.Preload(s); err != nil {
		timer.FatalExit("preload failed", 5, err)
	}

	fmt.Print("\n")
	fmt.Print("\n")

	dg, dgErr := discordgo.New("Bot " + os.Getenv("TOKEN"))
	if dgErr != nil {
		timer.FatalExit("discordgo: new session error:", 5, dgErr)
		return
	}

	s.Prefix = "discordgo: connecting"
	s.FinalMSG = "connection OK\n"
	s.Start()
	if err := dg.Open(); err != nil {
		timer.FatalExit("discordgo: conection err", 10, err)
		return
	}
	s.Stop()

	//discordgo identify intentsの設定はがちでよくわからん
	dg.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentDirectMessages | discordgo.IntentMessageContent
	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {

		MessageCheck(s, m, payload)

	}) //messageCheck()が繰り返される msg.go を参照
	fmt.Println("注意：モデルはPCが再起動されるまでずっとメモリに保持されます。\n消したい場合はmodel_unloader.exeを実行してください。")
	s.Prefix = "nixie稼働中"
	s.FinalMSG = "session end"
	s.Start()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	breaker := Done()
	select {
	case <-ctx.Done():
	case <-breaker:
		log.Println("nixieDOWNをうけとりました。")
	}
	s.Stop()

	dg.Close()
	os.Exit(1)
}

func ModelUnloader() error {
	fmt.Println("モデルはPCが再起動されるまでずっとメモリに保持されます。\nアンロードしますか？")
	yn := wizard.AskYesOrNo()
	if yn {
		err := wizard.Unload(s)
		if err != nil {
			return err
		}
		return nil
	}else{
		return fmt.Errorf("canceled")
	}
}

// helper func

func makeSpinner() *spinner.Spinner {
	s := spinner.New(spinner.CharSets[43], 100*time.Millisecond)
	s.Prefix = ""
	return s
}
