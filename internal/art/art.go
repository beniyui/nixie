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
