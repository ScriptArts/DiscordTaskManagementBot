package bot

import (
	"errors"
	"fmt"
	"github.com/ScriptArts/DiscordTaskManagementBot/models"
	"github.com/ScriptArts/DiscordTaskManagementBot/utils"
	"github.com/bwmarrin/discordgo"
	"time"
)

type CreatorCommand struct {
	BaseCommand
}

// $creator
// $creator list # クリエイターをすべて取得
// $creator add <@xxxxxx> # クリエイターを追加
// $creator remove <uuid> #クリエイターを削除（紐づく依頼があれば削除に失敗する）
// $creator remove <uuid> -f #クリエイターを削除（紐づく依頼もすべて削除する）
func (c *CreatorCommand) Run(params []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	c.session = s
	c.messageCreate = m

	// コマンドヘルプ
	if len(params) == 1 {
		return c.Help()
	}

	var err error

	p := params[1]
	switch p {
	case "list":
		err = c.List()
	case "add":
		err = c.Add(m.Mentions)
		break
	case "remove":
		if len(params) < 3 {
			return errors.New("正しいコマンドを入力してください")
		}

		uid := params[2]
		forceDelete := false
		if len(params) > 3 && params[3] == "-f" {
			forceDelete = true
		}

		err = c.Remove(uid, forceDelete)
		break
	}

	return err
}

func (c *CreatorCommand) Help() error {
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name: "タスク管理ボット",
		},
		Color: 0x00ff00,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "$creator list",
				Value: "クリエイター一覧を取得",
			},
			{
				Name:  "$creator add @xxxx",
				Value: "クリエイターを追加",
			},
			{
				Name:  "$creator remove <uuid>",
				Value: "クリエイターを削除（依頼がまだ存在している場合は失敗）",
			},
			{
				Name:  "$creator remove <uuid> -f",
				Value: "クリエイターを削除（関連する依頼も削除）",
			},
		},
		Description: "クリエイターコマンドのヘルプです",
		Timestamp:   time.Now().Format(time.RFC3339),
		Title:       "クリエイターコマンド",
	}

	_, err := c.session.ChannelMessageSendEmbed(c.messageCreate.ChannelID, embed)
	return err
}

func (c *CreatorCommand) List() error {
	r := new(models.CreatorRepository)
	creators, err := r.GetAll(c.messageCreate.GuildID)
	if err != nil {
		return err
	}

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name: "タスク管理ボット",
		},
		Color:       0x00ff00,
		Fields:      []*discordgo.MessageEmbedField{},
		Description: "クリエイターの一覧を表示します",
		Timestamp:   time.Now().Format(time.RFC3339),
		Title:       "クリエイター",
	}

	for _, creator := range creators {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  creator.UUID,
			Value: fmt.Sprintf("<@%s>", creator.DiscordID),
		})
	}

	c.session.ChannelMessageSendEmbed(c.messageCreate.ChannelID, embed)

	return nil
}

func (c *CreatorCommand) Add(mentions []*discordgo.User) error {
	if !utils.IsCanUseOpCommand(c.messageCreate.Author) {
		return errors.New("使用できないコマンドです")
	}

	if len(mentions) == 0 {
		return errors.New("メンションが設定されていません")
	}

	user := mentions[0]

	r := new(models.CreatorRepository)
	err := r.Create(user.ID, c.messageCreate.GuildID)
	if err != nil {
		return err
	}

	str := fmt.Sprintf("%s <@%s> を追加しました", user.Username, user.ID)
	c.session.ChannelMessageSend(c.messageCreate.ChannelID, str)

	return nil
}

func (c *CreatorCommand) Remove(uid string, forceDelete bool) error {
	if !utils.IsCanUseOpCommand(c.messageCreate.Author) {
		return errors.New(fmt.Sprintf("<@%s> 使用できないコマンドです", c.messageCreate.Author.ID))
	}

	r := new(models.CreatorRepository)
	creator, err := r.Get(uid, c.messageCreate.GuildID)
	if err != nil {
		return err
	}

	err = r.Remove(creator.DiscordID, c.messageCreate.GuildID, forceDelete)
	if err != nil {
		return err
	}

	str := fmt.Sprintf("<@%s> クリエイター「<@%s>」 を削除しました", c.messageCreate.Author.ID, creator.DiscordID)
	c.session.ChannelMessageSend(c.messageCreate.ChannelID, str)

	return nil
}
