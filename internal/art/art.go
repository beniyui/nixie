package art

import (
	"embed"
)

//go:embed assets
var files embed.FS

// GetFS 他のパッケージからこの関数を呼んで埋め込みファイルを取得します
func GetFS() embed.FS {
	return files
}

func GetArt(path string) (string, error) {
	buf, err := GetFS().ReadFile(path)
	if err != nil {
		return "", err
	}
	msg := string(buf)
	return string(msg), nil
}