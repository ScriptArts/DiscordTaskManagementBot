package models

import (
	"errors"
	"github.com/satori/go.uuid"
	"os"
	"time"
)

type Request struct {
	ID        uint          `gorm:"primary_key" json:"id"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	CreatorID uint          `json:"creator_id"`
	ClientID  uint          `json:"client_id"`
	Content   string        `json:"content"`
	Status    RequestStatus `json:"status"`
	UUID      string        `json:"uuid" gorm:"unique"`
	GuildID   string        `json:"guild_id"`
}

type RequestData struct {
	ID               uint      `json:"id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	CreatorDiscordID string    `json:"creator_discord_id"`
	ClientDiscordID  string    `json:"client_discord_id"`
	Content          string    `json:"content"`
	Status           int       `json:"status"`
	UUID             string    `json:"uuid" gorm:"unique"`
}

type RequestStatus int

const (
	StatusRequest RequestStatus = iota
	StatusAccept
	StatusDoing
	StatusDone
	StatusCancel
)

type RequestUserType int

const (
	UserTypeCreator RequestUserType = iota
	UserTypeClient
)

type RequestRepository struct{}

func (r *RequestRepository) GetAll(uid, guildID string) ([]*RequestData, error) {
	db, err := GetDatabase()
	if err != nil {
		return nil, err
	}

	var data []*RequestData
	err = db.Table("requests as r").
		Select("r.id, r.created_at, r.updated_at, c.discord_id as creator_discord_id, c2.discord_id as client_discord_id, r.content, r.status, r.uuid").
		Joins("LEFT JOIN creators c ON c.id = r.creator_id").
		Joins("LEFT JOIN clients c2 ON c2.id = r.client_id").
		Where("r.guild_id = ? AND (c.discord_id = ? OR c2.discord_id = ?)", guildID, uid, uid).Scan(&data).Error

	return data, err
}

func (r *RequestRepository) Get(uid, guildID string) (*RequestData, error) {
	db, err := GetDatabase()
	if err != nil {
		return nil, err
	}

	var data RequestData
	err = db.Table("requests as r").
		Select("r.id, r.created_at, r.updated_at, c.discord_id as creator_discord_id, c2.discord_id as client_discord_id, r.content, r.status, r.uuid").
		Joins("LEFT JOIN creators c ON c.id = r.creator_id").
		Joins("LEFT JOIN clients c2 ON c2.id = r.client_id").
		Where("r.uuid = ? AND r.guild_id = ?", uid, guildID).Scan(&data).Error

	return &data, err
}

func (r *RequestRepository) Create(creatorID uint, clientDiscordID, guildID, content string) error {
	db, err := GetDatabase()
	if err != nil {
		return err
	}

	tx := db.Begin()
	var client Client
	if err := tx.Find(&client, "discord_id = ? AND guild_id = ?", clientDiscordID, guildID).Error; err != nil {
		// 依頼者データが存在しないので作成する
		c := &Client{
			DiscordID: clientDiscordID,
			UUID:      uuid.NewV4().String(),
			GuildID:   guildID,
		}

		if err := tx.Save(c).Error; err != nil {
			tx.Rollback()
			return err
		}

		client = *c
	}

	// 連続リクエスト防止処理
	if os.Getenv("DISCORD_BOT_DEBUG") != "true" {
		if !client.LastRequestAt.IsZero() {
			if client.LastRequestAt.Add(time.Minute).After(time.Now()) {
				tx.Rollback()
				return errors.New("エラー")
			}
		}
	}

	req := &Request{
		CreatorID: creatorID,
		ClientID:  client.ID,
		Content:   content,
		Status:    StatusRequest,
		UUID:      uuid.NewV4().String(),
		GuildID:   guildID,
	}

	if err := tx.Save(req).Error; err != nil {
		tx.Rollback()
		return err
	}

	client.LastRequestAt = time.Now()
	if err := tx.Save(client).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *RequestRepository) UpdateContent(requestID uint, guildID, content string) error {
	db, err := GetDatabase()
	if err != nil {
		return err
	}

	var req Request
	if err := db.Find(&req, "id = ? AND guild_id = ?", requestID, guildID).Error; err != nil {
		return err
	}

	req.Content = content
	return db.Save(&req).Error
}

func (r *RequestRepository) UpdateStatus(requestID uint, guildID string, status RequestStatus) error {
	db, err := GetDatabase()
	if err != nil {
		return err
	}

	var req Request
	if err := db.Find(&req, "id = ? AND guild_id = ?", requestID, guildID).Error; err != nil {
		return err
	}

	req.Status = status
	return db.Save(&req).Error
}

func (r *RequestRepository) Remove(id uint) error {
	db, err := GetDatabase()
	if err != nil {
		return err
	}

	var req Request
	if err := db.Find(&req, "id = ?", id).Error; err != nil {
		return err
	}

	return db.Delete(&req).Error
}

func (r *RequestRepository) GetRequestUserType(requestID uint, guildID, discordID string) (RequestUserType, error) {
	db, err := GetDatabase()
	if err != nil {
		return -1, err
	}

	var data RequestData
	err = db.Table("requests as r").
		Select("r.id, r.created_at, r.updated_at, c.discord_id as creator_discord_id, c2.discord_id as client_discord_id, r.content, r.status, r.uuid").
		Joins("LEFT JOIN creators c ON c.id = r.creator_id").
		Joins("LEFT JOIN clients c2 ON c2.id = r.client_id").
		Where("r.id = ? AND r.guild_id = ?", requestID, guildID).Scan(&data).Error
	if err != nil {
		return -1, err
	}

	if data.CreatorDiscordID == discordID {
		return UserTypeCreator, nil
	}

	if data.ClientDiscordID == discordID {
		return UserTypeClient, nil
	}

	return -1, errors.New("")
}
