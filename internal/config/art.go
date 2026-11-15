package config

import (
	"embed"
	"fmt"
	"path"
)

//go:embed assets
var files embed.FS

// アート取得
func getArt(artVer string, kind string) (string, error) {
	path := path.Join("assets", artVer, fmt.Sprintf("%s.txt", kind))
	buf, err := files.ReadFile(path)
	if err != nil {
		return "", err
	}
	msg := string(buf)
	return string(msg), nil
}
