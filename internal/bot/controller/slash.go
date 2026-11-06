package controller

import "github.com/bwmarrin/discordgo"

// スラッシュコマンドの登録
func registerCommands(s *discordgo.Session) error {
	_, err := s.ApplicationCommandCreate(s.State.User.ID, "", &discordgo.ApplicationCommand{
		Name:        "controller",
		Description: "ディスコード上でコントロールパネルを開く",
	})
	return err
}

// intent handler
func onInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		handleCommand(s, i)
	case discordgo.InteractionMessageComponent:
		handleComponent(s, i)
	}
}

// スラッシュコマンドの処理
func handleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.ApplicationCommandData().Name != "control" {
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    "🎛️ **ボット コントロールパネル**\n下のボタンから操作を選択してください。",
			Components: getPanelComponents(),
			Flags:      discordgo.MessageFlagsEphemeral,
		},
	})
}

// ボタンなどのコンポーネント操作の処理
func handleComponent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	customID := i.MessageComponentData().CustomID
	responseText := getResponseText(customID)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content:    responseText,
			Components: getPanelComponents(),
		},
	})
}

// ボタンの CustomID に応じたテキストを返す
func getResponseText(customID string) string {
	switch customID {
	case "btn_status":
		return "🟢 システム状態: 正常稼働中 (CPU: 12%, Mem: 45%)"
	case "btn_reboot":
		return "🔄 パネルを再起動しています..."
	default:
		return "未知のアクションです。"
	}
}

// コントロールパネルのUI（ボタン）定義
func getPanelComponents() []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "ステータス確認",
					Style:    discordgo.PrimaryButton,
					CustomID: "btn_status",
				},
				discordgo.Button{
					Label:    "システム再起動",
					Style:    discordgo.DangerButton,
					CustomID: "btn_reboot",
				},
			},
		},
	}
}
