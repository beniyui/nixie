package config

//this is config files reader

import (
	"bufio"
	"fmt"
	"nixie/internal/art"
	"os"
	"path/filepath"
	"reflect"

	"github.com/joho/godotenv"
)

type (
	ModelID string
)

// 初期化処理で一括で読み込む
type envContainer struct {
	Server     string
	Version    string
	Token      string
	WebHookURL string
	DensuID    string
	ExeDir     string

	AdminIDs []string
}

type artContainer struct {
	TitleBanner string
	UpBanner    string
}

type promptContainer struct {
	ID   string
	Inst string
}

// 環境変数にアクセスできる
var Values *envContainer

// artにアクセスできる
var Arts *artContainer

// プロンプトにアクセスする
var Prompts *[]promptContainer

func Init() error {

	//作業ディレクトリを習得しないと始まらない。
	err := Values.dirSetp()
	if err != nil {
		return fmt.Errorf("dirSetup: %v Values: %+v", err, Values)
	}

	err = Values.envSetup()
	if err != nil {
		return fmt.Errorf("envSetup: %v Values: %+v", err, Values)
	}
	err = Values.IdsSetup()
	if err != nil {
		return fmt.Errorf("IdsSetup: %v Values: %+v", err, Values)
	}
	err = Arts.artSetup()
	if err != nil {
		return fmt.Errorf("artSetup: %v Values: %+v", err, Arts)
	}
	return nil
}

func (e *envContainer) dirSetp() error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %v", err)
	}
	dir := filepath.Dir(exePath)
	e.ExeDir = dir
	return nil
}

// .envをロードし，環境変数やアートをコンテナに格納。
func (e *envContainer) envSetup() error {

	//load .env
	envPath := filepath.Join(e.ExeDir, ".env")
	if err := godotenv.Load(envPath); err != nil {
		return fmt.Errorf("godotenv.Load failed. env file path is wrong???:%d", err)
	}

	e = &envContainer{
		Server:     os.Getenv("SERVER"),
		Version:    os.Getenv("VERSION"),
		Token:      os.Getenv("TOKEN"),
		WebHookURL: os.Getenv("WEBHOOK_URL_1"),
		DensuID:    os.Getenv("DENSU_ID"),
	}

	if !isAllFieldsHasValue(e) {
		return fmt.Errorf("envにはアクセスできたけど一部envの読み取りに失敗した。")
	}
	return nil
}

func (a *artContainer) artSetup() error {
	ver := os.Getenv("ART_VERSION")
	if ver == "" {
		return fmt.Errorf("failed to get art version")
	}

	getTrgArt := func(kind string) (string, error) {
		name := fmt.Sprintf("%s.txt", ver)
		path := filepath.Join("assets", kind, name)
		art, err := art.GetArt(path)
		if err != nil {
			return "", err
		}
		return art, nil
	}
	ti, err := getTrgArt("title")
	if err != nil {
		return err
	}
	up, err := getTrgArt("up")
	if err != nil {
		return err
	}
	a = &artContainer{
		TitleBanner: ti,
		UpBanner:    up,
	}
	if !isAllFieldsHasValue(a) {
		return fmt.Errorf("ん？Artが取得できてなくね？")
	}
	return nil
}

func (e *envContainer) IdsSetup() error {
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

func promptSetup() error {
	trgDir := filepath.Join(Values.ExeDir,"prompt")
	NameList, err := os.ReadDir(trgDir)
	if err != nil {
		return err
	}
	var contentList []string
	for range NameList{
		
	}
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
