package models

import (
	"errors"
	"github.com/satori/go.uuid"
	"time"
)

type Creator struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DiscordID string    `json:"discord_id" gorm:"unique_index:idx_creator_discord_id_guild_id"`
	UUID      string    `json:"uuid" gorm:"unique"`
	Requests  []Request `json:"requests"`
	GuildID   string    `json:"guild_id" gorm:"unique_index:idx_creator_discord_id_guild_id"`
}

type CreatorRepository struct{}

func (r *CreatorRepository) GetAll(guildID string) ([]*Creator, error) {
	db, err := GetDatabase()
	if err != nil {
		return nil, err
	}

	var creators []*Creator
	if err := db.Find(&creators, "guild_id = ?", guildID).Error; err != nil {
		return nil, err
	}

	return creators, nil
}

func (r *CreatorRepository) Get(uid, guildID string) (*Creator, error) {
	db, err := GetDatabase()
	if err != nil {
		return nil, err
	}

	var creator Creator
	if err := db.Find(&creator, "uuid = ? AND guild_id = ?", uid, guildID).Error; err != nil {
		return nil, err
	}

	return &creator, nil
}

func (r *CreatorRepository) GetByDiscordID(discordID, guildID string) (*Creator, error) {
	db, err := GetDatabase()
	if err != nil {
		return nil, err
	}

	var creator Creator
	if err := db.Find(&creator, "discord_id = ? AND guild_id = ?", discordID, guildID).Error; err != nil {
		return nil, err
	}

	return &creator, nil
}

func (r *CreatorRepository) Create(discordID, guildID string) error {
	db, err := GetDatabase()
	if err != nil {
		return err
	}

	var exist Creator
	if err := db.Find(&exist, "discord_id = ? AND guild_id = ?", discordID, guildID).Error; err == nil {
		// 存在する
		return errors.New("すでに存在するクリエイターです")
	}

	uid := uuid.NewV4()
	c := &Creator{
		DiscordID: discordID,
		UUID:      uid.String(),
		GuildID:   guildID,
	}

	return db.Save(c).Error
}

func (r *CreatorRepository) Remove(discordID, guildID string, forceDelete bool) error {
	db, err := GetDatabase()
	if err != nil {
		return err
	}

	tx := db.Begin()
	var creator Creator
	if err := db.Find(&creator, "discord_id = ? AND guild_id", discordID, guildID).Error; err != nil {
		tx.Rollback()
		return err
	}

	var requests []*Request
	db.Find(&requests, "creator_id = ? AND guild_id", creator.ID, guildID)
	if forceDelete == false && len(requests) > 0 {
		tx.Rollback()
		return errors.New("クリエイターには依頼が存在します")
	}

	for _, request := range requests {
		if err := tx.Delete(&request).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Delete(&creator).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
