開発指示

このコードはディスコードサーバー「ゲイリー」にて動作する私的な利用のためのボットのソースコードである
ボットの名前は"nixie（ニキシー）"である

go言語を使用する
ライブラリはbwmarrin/discordgoを使用する

機能：
/core/func/ping.go ニキシーに"ping"とDMすると"pong"と返信する メンテナンス用機能
/core/func/anon.go ニキシーはDMを受け取るとそれをウェブフックで送信する(anon = 匿名)

将来：
/core/func/clush.go ニキシーは定期的に無差別にサーバーの人間をキックする
