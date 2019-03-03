package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/txgruppi/parseargs-go"
	"log"
	"strings"
)

type BaseCommand struct {
	session       *discordgo.Session
	messageCreate *discordgo.MessageCreate
}

type CommandInterface interface {
	Run(params []string, s *discordgo.Session, m *discordgo.MessageCreate) error
}

// 誰かがメッセージを投稿したときのハンドラー
func MentionHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	c := m.Content
	if !strings.HasPrefix(c, "$") {
		return
	}

	log.Println(m.GuildID)

	parsed, err := parseargs.Parse(c)
	if err != nil {
		log.Fatal(err)
	}

	var cmd CommandInterface

	switch parsed[0] {
	case "$creator":
		cmd = &CreatorCommand{}
	case "$request":
		cmd = &RequestCommand{}
		break
	}

	if cmd != nil {
		err = cmd.Run(parsed, s, m)
	}

	if err != nil {
		log.Println(err.Error())
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> %s", m.Author.ID, err.Error()))
	}
	//log.Println("[Single mention]", m.Author.Username, m.Content)
}
