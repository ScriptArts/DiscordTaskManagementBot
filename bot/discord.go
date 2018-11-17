package bot

import (
	"github.com/bwmarrin/discordgo"
	"os"
)

// Discord Client取得
func GetDiscordClient() (*discordgo.Session, error) {
	d, err := discordgo.New("Bot " + os.Getenv("DISCORD_BOT_TOKEN"))
	if err != nil {
		return nil, err
	}

	return d, nil
}
