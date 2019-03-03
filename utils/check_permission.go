package utils

import (
	"github.com/bwmarrin/discordgo"
	"os"
	"strings"
)

func IsCanUseOpCommand(user *discordgo.User) bool {
	l := os.Getenv("DISCORD_BOT_OP_LIST")
	list := strings.Split(l, ",")
	for _, v := range list {
		if v == user.ID {
			return true
		}
	}

	return false
}
