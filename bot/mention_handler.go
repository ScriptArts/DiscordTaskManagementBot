package bot

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

// 誰かがメッセージを投稿したときのハンドラー
func MentionHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	mentions := m.Mentions
	if len(mentions) != 1 {
		return
	}

	log.Println("[Single mention]", m.Author.Username, m.Content)
}
