package models

import (
	"gorm.io/datatypes"
	"time"
)

type Post struct {
	Model

	Content     string                      `json:"content"`
	Comments    []Post                      `json:"comments" gorm:"foreignKey:BelongID"`
	Tags        datatypes.JSONSlice[string] `json:"tags"`
	Attachments datatypes.JSONSlice[string] `json:"attachments"`
	Type        string                      `json:"type"`
	IpAddress   string                      `json:"ip_address"`
	PublishedAt time.Time                   `json:"published_at"`
	BelongID    *uint                       `json:"belong_id"`
	AccountID   uint                        `json:"account_id"`
}
