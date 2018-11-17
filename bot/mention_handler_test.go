package bot

import (
	"github.com/bwmarrin/discordgo"
	"testing"
)

func TestMentionHandler(t *testing.T) {
	t.Helper()
	s := &discordgo.Session{}

	var mentions []*discordgo.User
	mentions = append(mentions, &discordgo.User{})

	mc := &discordgo.MessageCreate{
		Message: &discordgo.Message{
			ID:        "messageid",
			Content:   "!ping",
			ChannelID: "channelid",
			Author: &discordgo.User{
				Username: "sample",
			},
			Mentions: mentions,
		},
	}

	MentionHandler(s, mc)
}

func TestMentionHandlerMultiple(t *testing.T) {
	t.Helper()
	s := &discordgo.Session{}

	var mentions []*discordgo.User
	mentions = append(mentions, &discordgo.User{})
	mentions = append(mentions, &discordgo.User{})

	mc := &discordgo.MessageCreate{
		Message: &discordgo.Message{
			ID:        "messageid",
			Content:   "!ping",
			ChannelID: "channelid",
			Author: &discordgo.User{
				Username: "sample",
			},
			Mentions: mentions,
		},
	}

	MentionHandler(s, mc)
}
