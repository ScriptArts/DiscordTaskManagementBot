package models

import "time"

type Client struct {
	ID            uint      `gorm:"primary_key" json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DiscordID     string    `json:"discord_id" gorm:"unique_index:idx_client_discord_id_guild_id"`
	UUID          string    `json:"uuid" gorm:"unique"`
	LastRequestAt time.Time `json:"last_request_at"`
	GuildID       string    `json:"guild_id" gorm:"unique_index:idx_client_discord_id_guild_id"`
}

type ClientRepository struct{}

func (r *ClientRepository) Get(discordId string) (*Client, error) {
	db, err := GetDatabase()
	if err != nil {
		return nil, err
	}

	var c Client
	err = db.Find(&c, "discord_id = ?", discordId).Error

	return &c, err
}
