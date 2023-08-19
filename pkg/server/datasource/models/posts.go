package models

import (
	"gorm.io/datatypes"
	"time"
)

const (
	PostTypeText  = "text"
	PostTypeImage = "image"
	PostTypeVideo = "video"
)

type Post struct {
	Model

	Content     string                      `json:"content"`
	Comments    []*Post                     `json:"comments" gorm:"many2many:post_comments"`
	Tags        datatypes.JSONSlice[string] `json:"tags"`
	Attachments datatypes.JSONSlice[string] `json:"attachments"`
	Type        string                      `json:"type"`
	IpAddress   string                      `json:"ip_address"`
	PublishedAt time.Time                   `json:"published_at"`
	AccountID   uint                        `json:"account_id"`
}
