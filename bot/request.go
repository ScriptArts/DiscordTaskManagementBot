package bot

import (
	"errors"
	"fmt"
	"github.com/ScriptArts/DiscordTaskManagementBot/models"
	"github.com/bwmarrin/discordgo"
	"github.com/satori/go.uuid"
	"log"
	"regexp"
	"time"
)

type RequestCommand struct {
	BaseCommand
}

// $request
// $request list # 自分に関係する依頼を取得
// $request add <クリエイターUUID> <内容> # クリエイターに依頼をする
// $request update <依頼UUID> <内容> # 依頼の内容を更新（依頼者のみ）
// $request update_status <依頼UUID> <request/accept/doing/done/cancel> <コメント> # 依頼のステータスを更新（依頼者はキャンセルにのみ更新可能）
// $request remove <依頼UUID> <コメント> #依頼を削除
func (c *RequestCommand) Run(params []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
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
		if len(params) < 4 {
			return errors.New("正しいコマンドを入力してください")
		}

		err = c.Add(params[2], params[3])
	case "update":
		if len(params) < 4 {
			return errors.New("正しいコマンドを入力してください")
		}
		err = c.Update(params[2], params[3])
	case "update_status":
		if len(params) < 4 {
			return errors.New("正しいコマンドを入力してください")
		}

		comment := ""
		if len(params) > 4 {
			comment = params[4]
		}

		err = c.UpdateStatus(params[2], params[3], comment)
	case "remove":
		if len(params) < 4 {
			return errors.New("正しいコマンドを入力してください")
		}
		err = c.Remove(params[2], params[3])
	}

	return err
}

// 依頼コマンドのヘルプ
func (c *RequestCommand) Help() error {
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name: "タスク管理ボット",
		},
		Color: 0x00ff00,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "$request list",
				Value: "自分に関係ある一覧を取得",
			},
			{
				Name:  "$request add <uuid> <内容>",
				Value: "依頼を行う",
			},
			{
				Name:  "$request update <uuid> <内容>",
				Value: "依頼内容を変更（依頼者のみ可能）",
			},
			{
				Name:  "$request update_status <uuid> <request/accept/doing/done/cancel> <コメント>",
				Value: "依頼の進行状況を更新（依頼者は request のものを cancel に更新可能。クリエイターは自由更新。）",
			},
			{
				Name:  "$request remove <uuid> <コメント>",
				Value: "依頼を削除。クリエイターのみ使用可能。",
			},
		},
		Description: "クリエイターコマンドのヘルプです",
		Timestamp:   time.Now().Format(time.RFC3339),
		Title:       "クリエイターコマンド",
	}

	_, err := c.session.ChannelMessageSendEmbed(c.messageCreate.ChannelID, embed)
	return err
}

// 依頼リストコマンド
func (c *RequestCommand) List() error {
	// コマンド実行者
	discordId := c.messageCreate.Author.ID

	r := new(models.RequestRepository)
	requests, err := r.GetAll(discordId, c.messageCreate.GuildID)
	if err != nil {
		return err
	}

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name: "タスク管理ボット",
		},
		Color:       0x00ff00,
		Fields:      []*discordgo.MessageEmbedField{},
		Description: "依頼の一覧を表示します",
		Timestamp:   time.Now().Format(time.RFC3339),
		Title:       "依頼",
	}

	for _, request := range requests {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name: request.UUID,
			Value: fmt.Sprintf(
				"クリエイター：<@%s>\n依頼者：<@%s>\n内容：%s\nステータス：%s",
				request.CreatorDiscordID, request.ClientDiscordID, request.Content, c.getStatus(request.Status),
			),
		})
	}

	c.session.ChannelMessageSendEmbed(c.messageCreate.ChannelID, embed)

	return nil
}

func (c *RequestCommand) Add(id, content string) error {
	creator, err := c.getCreator(id, c.messageCreate.GuildID)
	if err != nil {
		return err
	}

	log.Println(creator)
	clientID := c.messageCreate.Author.ID

	// 作成処理
	r := new(models.RequestRepository)
	err = r.Create(creator.ID, clientID, c.messageCreate.GuildID, content)
	if err != nil {
		return err
	}

	// コマンド実行結果を送信
	str := fmt.Sprintf("<@%s> 依頼を作成しました。\nクリエイター：<@%s>\n依頼内容：%s", clientID, creator.DiscordID, content)
	c.session.ChannelMessageSend(c.messageCreate.ChannelID, str)

	// 依頼者に対して、プライベートメッセージで依頼内容を送信
	channel, _ := c.session.UserChannelCreate(creator.DiscordID)
	str = fmt.Sprintf("<@%s> から依頼が届きました。\n依頼内容：%s", clientID, content)
	c.session.ChannelMessageSend(channel.ID, str)

	return nil
}

func (c *RequestCommand) Update(uid, content string) error {
	user := c.messageCreate.Author

	requestRepo := new(models.RequestRepository)
	request, err := requestRepo.Get(uid, c.messageCreate.GuildID)
	if err != nil {
		return err
	}

	if request.ClientDiscordID != user.ID {
		return errors.New("xxx")
	}

	err = requestRepo.UpdateContent(request.ID, c.messageCreate.GuildID, content)
	if err != nil {
		return err
	}

	str := fmt.Sprintf("<@%s> 依頼が更新されました。\nクリエイター：<@%s>\n依頼内容：%s", user.ID, request.CreatorDiscordID, content)
	c.session.ChannelMessageSend(c.messageCreate.ChannelID, str)

	channel, _ := c.session.UserChannelCreate(request.CreatorDiscordID)
	str = fmt.Sprintf("<@%s> によって、依頼が更新されました。\n依頼内容：%s\nステータス：%s", user.ID, content, c.getStatus(int(request.Status)))
	c.session.ChannelMessageSend(channel.ID, str)

	return nil
}

func (c *RequestCommand) UpdateStatus(uid, status, comment string) error {
	user := c.messageCreate.Author

	s, err := c.convertStatusStrToInt(status)
	if err != nil {
		return err
	}

	requestRepo := new(models.RequestRepository)
	request, err := requestRepo.Get(uid, c.messageCreate.GuildID)
	if err != nil {
		return err
	}

	// 依頼者が更新する場合
	if request.CreatorDiscordID == user.ID {
		err = c.updateStatusByCreator(request, s)
	} else if request.ClientDiscordID == user.ID {
		err = c.updateStatusByClient(request, s)
	} else {
		err = errors.New("")
	}

	if err != nil {
		return err
	}

	str := fmt.Sprintf("<@%s> 依頼が更新されました。\nクリエイター：<@%s>\n依頼内容：%s", user.ID, request.CreatorDiscordID, request.Content)
	c.session.ChannelMessageSend(c.messageCreate.ChannelID, str)

	channel, _ := c.session.UserChannelCreate(request.CreatorDiscordID)
	str = fmt.Sprintf("<@%s> によって、依頼が更新されました。\n依頼内容：%s\nステータス：%s", user.ID, request.Content, c.getStatus(int(s)))
	c.session.ChannelMessageSend(channel.ID, str)

	return nil
}

func (c *RequestCommand) Remove(uid, content string) error {
	user := c.messageCreate.Author

	requestRepo := new(models.RequestRepository)
	request, err := requestRepo.Get(uid, c.messageCreate.GuildID)
	if err != nil {
		return err
	}

	if request.ClientDiscordID != user.ID {
		return errors.New("削除できない依頼です")
	}

	err = requestRepo.Remove(request.ID)
	if err != nil {
		return err
	}

	str := fmt.Sprintf("<@%s> 依頼を削除しました", user.ID)
	c.session.ChannelMessageSend(c.messageCreate.ChannelID, str)

	channel, _ := c.session.UserChannelCreate(request.ClientDiscordID)
	str = fmt.Sprintf("<@%s> によって、依頼が削除されました。\n依頼内容：%s\nステータス：%s", user.ID, request.Content, c.getStatus(int(request.Status)))
	c.session.ChannelMessageSend(channel.ID, str)

	return nil
}

func (c *RequestCommand) getStatus(status int) string {
	switch status {
	case 0:
		return "依頼中"
	case 1:
		return "受付済"
	case 2:
		return "作業中"
	case 3:
		return "完了"
	case 4:
		return "キャンセル"
	default:
		return ""
	}
}

func (c *RequestCommand) convertStatusStrToInt(status string) (models.RequestStatus, error) {
	switch status {
	case "request":
		return models.StatusRequest, nil
	case "accept":
		return models.StatusAccept, nil
	case "doing":
		return models.StatusDoing, nil
	case "done":
		return models.StatusDone, nil
	case "cancel":
		return models.StatusCancel, nil
	default:
		return -1, errors.New("不正なステータスです")
	}
}

func (c *RequestCommand) getCreator(id, guildID string) (*models.Creator, error) {
	creatorRepo := new(models.CreatorRepository)

	_, err := uuid.FromString(id)
	if err == nil {
		log.Println(id)
		creator, err := creatorRepo.Get(id, c.messageCreate.GuildID)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}

		return creator, nil
	}

	r := regexp.MustCompile("<@(.*)>")
	result := r.FindStringSubmatch(id)
	if len(result) != 2 {
		return nil, errors.New("")
	}

	user, err := c.session.User(result[1])
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	log.Println(user.ID)

	return creatorRepo.GetByDiscordID(user.ID, c.messageCreate.GuildID)
}

func (c *RequestCommand) updateStatusByClient(request *models.RequestData, status models.RequestStatus) error {
	if request.Status != int(models.StatusRequest) {
		return errors.New("")
	}

	if request.Status != int(models.StatusCancel) {
		return errors.New("")
	}

	requestRepo := new(models.RequestRepository)
	return requestRepo.UpdateStatus(request.ID, c.messageCreate.GuildID, status)
}

func (c *RequestCommand) updateStatusByCreator(request *models.RequestData, status models.RequestStatus) error {
	requestRepo := new(models.RequestRepository)
	return requestRepo.UpdateStatus(request.ID, c.messageCreate.GuildID, status)
}
