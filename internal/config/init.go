package config

//this is config files reader

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/joho/godotenv"
)

type (
	ModelID string
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

type artValues struct {
	TitleBanner string
	UpBanner    string
}

type instValues struct {
	Nixie  string
	NixieS string
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

// setup these containeres
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

	return fmt.Errorf("error test")
}

func (e *envValues) envSetup() error {

	//get executable path
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path. お前どういうトチ狂ったところからこいつを起動してんの!?: %v", err)
	}
	dir := filepath.Dir(exePath)
	e.ExeDir = dir

	//load the env file with godotenv
	envPath := filepath.Join(e.ExeDir, ".env")
	if err := godotenv.Load(envPath); err != nil {
		return fmt.Errorf("godotenv.Load failed. envファイルはどこだよ: %d", err)
	}

	*e = envValues{
		Server:     os.Getenv("SERVER"),
		ArtVersion: os.Getenv("ART_VERSION"),
		Version:    os.Getenv("VERSION"),
		Token:      os.Getenv("TOKEN"),
		WebHookURL: os.Getenv("WEBHOOK_URL_1"),
		DensuID:    os.Getenv("DENSU_ID"),
	}

	//that's all?
	if !isAllFieldsHasValue(e) {
		return fmt.Errorf("envにはアクセスできたけど一部envの読み取りができてないな")
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

	if !isAllFieldsHasValue(a) {
		return fmt.Errorf("ん？Artが取得できてなくね？")
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

//helper func:

// すべてのフィールドが値を持つかを判定する
func isAllFieldsHasValue(s any) bool {
	v := reflect.ValueOf(s)
	// for文でフィールドを1つずつチェック
	for _, field := range v.Fields() {
		field := field

		// フィールドが初期値の場合
		if field.IsZero() {
			return false
		}
	}
	return true
}
