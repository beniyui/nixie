package config

//this is config files reader

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// 初期化処理で一括で読み込む
type envValues struct {
	Server     string
	Version    string
	ArtVersion string
	Token      string
	WebHookURL string
	DensuID    string
	ExeDir     string

	AdminIDs []string
}

var envReqKeys = []string{
	"SERVER",
	"VERSION",
	"ART_VERSION",
	"TOKEN",
	"WEBHOOK_URL_1",
	"GEMINI_KEY",
	"TSKEY",
	"DENSU_ID",
}

type artValues struct {
	TitleBanner string
	UpBanner    string
}

type instValues struct {
	Nixie  string
	NixieT string
	Lily   string
	Gemini string
}

//containers

// 外部から環境変数にアクセス
var EnvCont envValues

// 外部からアートにアクセス
var ArtsCont artValues

// 外部からinstruction promptにアクセス
var PromptCont instValues

// the handler for setup these containeres
func Init() error {
	var err error

	err = EnvCont.envSetup()
	if err != nil {
		return fmt.Errorf("envSetup: %v", err)
	}

	err = ArtsCont.artSetup()
	if err != nil {
		return fmt.Errorf("artsetup:%v", err)
	}

	return nil
}

func (e *envValues) envSetup() error {

	//get executable path
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path. お前どういうトチ狂ったところからこいつを起動してんの!?: %v", err)
	}
	exedir := filepath.Dir(exePath)

	//load the env file with godotenv
	envPath := filepath.Join(exedir, ".env")
	if err := godotenv.Load(envPath); err != nil {
		return fmt.Errorf("godotenv.Load failed. envファイルはどこだよ: %v", err)
	}

	//that's all?
	err = envCheck()
	if err != nil {
		return fmt.Errorf("env Missing: %v", err)
	}

	*e = envValues{
		Server:     os.Getenv("SERVER"),
		ArtVersion: os.Getenv("ART_VERSION"),
		Version:    os.Getenv("VERSION"),
		Token:      os.Getenv("TOKEN"),
		WebHookURL: os.Getenv("WEBHOOK_URL_1"),
		DensuID:    os.Getenv("DENSU_ID"),
		ExeDir:     exedir,
	}
	return nil
}

func (a *artValues) artSetup() error {
	ver := EnvCont.ArtVersion

	ti, err := getArt(ver, "title")
	if err != nil {
		return err
	}
	up, err := getArt(ver, "up")
	if err != nil {
		return err
	}

	*a = artValues{
		TitleBanner: ti,
		UpBanner:    up,
	}

	return nil
}

func (e *envValues) IdsSetup() error {
	trgPath := filepath.Join(e.ExeDir, "data", "IDlist.txt")

	file, err := os.Open(trgPath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var ids []string
	for scanner.Scan() {
		id := scanner.Text()
		ids = append(ids, id)
	}
	e.AdminIDs = ids
	return scanner.Err()
}

func envCheck() error {
	var missings []string
	for _, key := range envReqKeys {
		_, ok := os.LookupEnv(key)
		if !ok {
			missings = append(missings, key)
		}
	}
	if missings != nil {
		return fmt.Errorf("missing env:%v", missings)
	}
	return nil
}
